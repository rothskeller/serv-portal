package regedit

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/pages/classes"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

var statusLabels = map[string]string{
	"registered": "Registered",
	"waitlist":   "On wait list",
	"inclass":    "Accepted, student role granted",
}

// Handle handles /classes/regedit/$id requests.
func Handle(r *request.Request, ridstr string) {
	const classFields = class.FID | class.FStart | class.FLimit | class.FReferrals | class.FType | class.FRegURL | class.FRole
	var (
		user      *person.Person
		cr        *classreg.ClassReg
		c         *class.Class
		p         *person.Person
		regby     *person.Person
		ur        *classreg.Updater
		personID  person.ID
		status    string
		nameError string
		f         form.Form
		statuses  = []string{"registered", "waitlist"} // may add "inclass" below
	)
	// Get the user information.
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	// Get the registration and class information.
	if cr = classreg.WithID(r, classreg.ID(util.ParseID(ridstr)), classreg.UpdaterFields); cr == nil {
		errpage.NotFound(r, user)
	}
	if c = class.WithID(r, cr.Class(), classFields); c == nil || c.RegURL() != "" {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(c.Type().Org(), enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	personID = cr.Person()
	p = person.WithID(r, personID, person.FInformalName|person.FEmail|person.FEmail2|person.FCellPhone|person.FHomePhone|person.FWorkPhone)
	regby = person.WithID(r, cr.RegisteredBy(), person.FInformalName|person.FEmail|person.FCellPhone)
	ur = cr.Updater(r, c, p, regby)
	if cr.Waitlist() {
		status = "waitlist"
	} else {
		status = "registered"
	}
	if c.Role() != 0 {
		statuses = append(statuses, "inclass")
		if cr.Person() != 0 {
			if held, _ := personrole.PersonHasRole(r, cr.Person(), c.Role()); held {
				status = "inclass"
			}
		}
	}
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	f.Title = "Edit Registration"
	f.Rows = []form.Row{
		&nameRow{form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditFirstname",
				Label: "First Name",
			},
			Name:     "firstname",
			ValueP:   &ur.FirstName,
			Validate: "#regeditFirstname,#regeditPerson",
		}, nil},
		&nameRow{form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditLastname",
				Label: "Last Name",
			},
			Name:     "lastname",
			ValueP:   &ur.LastName,
			Validate: "#regeditLastname,#regeditPerson",
		}, &nameError},
		&emailRow{form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditEmail",
				Label: "Email",
			},
			Name:     "email",
			ValueP:   &ur.Email,
			Validate: "#regeditEmail,#regeditPerson",
		}},
		&cellPhoneRow{form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditCellphone",
				Label: "Cell Phone",
			},
			Name:     "cellphone",
			ValueP:   &ur.CellPhone,
			Validate: "#regeditCellphone,#regeditPerson",
		}, ur},
		&form.MessageRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditRegdby",
				Label: "Registered By",
			},
			HTML: fmt.Sprintf("%s %s %s", regby.InformalName(), regby.Email(), regby.CellPhone()),
		},
		&havePersonRow{form.MessageRow{
			LabeledRow: form.LabeledRow{
				RowID: "regeditPerson",
				Label: "Person",
			},
		}, ur.Person},
		&wantPersonRow{
			RadioGroupRow: form.RadioGroupRow[person.ID]{
				LabeledRow: form.LabeledRow{
					RowID: "regeditPerson",
					Label: "Person",
				},
				Name:      "person",
				ValueP:    &personID,
				ValueFunc: func(v person.ID) string { return strconv.Itoa(int(v)) },
				Validate:  "#regeditPerson,#regeditStatus",
				Wide:      true,
			},
			class:    c,
			classreg: ur,
			personID: &personID,
		},
		&statusRow{form.RadioGroupRow[string]{
			LabeledRow: form.LabeledRow{
				RowID: "regeditStatus",
				Label: "Status",
			},
			Name:      "status",
			ValueP:    &status,
			Options:   statuses,
			LabelFunc: func(_ *request.Request, v string) string { return statusLabels[v] },
		}, &personID},
	}
	f.Buttons = []*form.Button{{
		Label: "Save",
		OnClick: func() (ok bool) {
			ok, nameError = saveRegistration(r, user, c, cr, ur, personID, status)
			return ok
		},
	}}
	if status != "inclass" {
		f.Buttons = append(f.Buttons, &form.Button{
			Name: "delete", Label: "Delete", Style: "danger",
			OnClick: func() bool { return deleteRegistration(r, user, c, cr) },
		})
	}
	f.Handle(r)
}

func saveRegistration(r *request.Request, user *person.Person, c *class.Class, cr *classreg.ClassReg, ur *classreg.Updater, personID person.ID, status string) (ok bool, nameError string) {
	r.Transaction(func() {
		if personID == -1 { // new person
			p := new(person.Updater)
			p.FormalName = ur.FirstName + " " + ur.LastName
			p.InformalName = p.FormalName
			p.SortName = ur.LastName + ", " + ur.FirstName
			if p.DuplicateSortName(r) {
				ok, nameError = false, "Another person exists with the same name."
				r.DoNotCommit()
				return
			}
			p.Email = ur.Email
			if p.DuplicateEmail(r) {
				p.Email, p.Email2 = "", p.Email
			}
			p.CellPhone = cr.CellPhone()
			ur.Person = person.Create(r, p)
		} else if personID != 0 && ur.Person == nil {
			ur.Person = person.WithID(r, personID, person.FID|person.FInformalName)
		}
		if c.Role() != 0 && ur.Person != nil {
			rl := role.WithID(r, c.Role(), role.FID|role.FName)
			if status == "inclass" {
				personrole.AddRole(r, ur.Person, rl)
			} else {
				personrole.RemoveRole(r, ur.Person, rl)
			}
		}
		ur.Waitlist = status == "waitlist"
		cr.Update(r, ur)
		ok = true
	})
	if ok {
		classes.RenderRegList(r, user, c)
	}
	return ok, nameError
}

func deleteRegistration(r *request.Request, user *person.Person, c *class.Class, cr *classreg.ClassReg) bool {
	r.Transaction(func() {
		cr.Delete(r, c)
	})
	classes.RenderRegList(r, user, c)
	return true
}

type nameRow struct {
	form.InputRow
	saveError *string
}

func (nr *nameRow) Read(r *request.Request) bool {
	if !nr.InputRow.Read(r) {
		return false
	}
	if *nr.ValueP == "" {
		nr.Error = "Both first and last name are required."
		return false
	}
	if nr.saveError != nil && *nr.saveError != "" {
		nr.Error = *nr.saveError
		*nr.saveError = ""
		return false
	}
	return true
}

type emailRow struct{ form.InputRow }

func (er *emailRow) Read(r *request.Request) bool {
	if !er.InputRow.Read(r) {
		return false
	}
	if *er.ValueP != "" && !emailRE.MatchString(*er.ValueP) {
		er.Error = "This is not a valid email address."
		return false
	}
	return true
}

type cellPhoneRow struct {
	form.InputRow
	ucr *classreg.Updater
}

func (cr *cellPhoneRow) Read(r *request.Request) bool {
	if !cr.InputRow.Read(r) {
		return false
	}
	if !fmtPhone(cr.ValueP) {
		cr.Error = "This is not a valid phone number."
		return false
	}
	if pid := cr.ucr.Person.ID(); pid != 0 {
		if p := person.WithCellPhone(r, *cr.ValueP, person.FID); p != nil && p.ID() != pid {
			cr.Error = "This cell phone number is in use by another person."
			return false
		}
	}
	return true
}

type havePersonRow struct {
	form.MessageRow
	p *person.Person
}

func (hr *havePersonRow) ShouldEmit(_ request.ValidationList) bool { return hr.p != nil }

func (hr *havePersonRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if hr.HTML == "" {
		hr.HTML = fmt.Sprintf("%s %s %s", hr.p.InformalName(), hr.p.Email(), hr.p.CellPhone())
	}
	hr.MessageRow.Emit(r, parent, focus)
}

type wantPersonRow struct {
	form.RadioGroupRow[person.ID]
	class    *class.Class
	classreg *classreg.Updater
	personID *person.ID
	labels   map[person.ID]string
}

func (wr *wantPersonRow) ShouldEmit(_ request.ValidationList) bool {
	return wr.classreg.Person.ID() == 0
}

func (wr *wantPersonRow) Get(r *request.Request) {
	const personFields = person.FCellPhone | person.FEmail | person.FEmail2 | person.FHomePhone | person.FID | person.FInformalName | person.FSortName | person.FWorkPhone
	var used = sets.New[person.ID]()
	wr.Options = []person.ID{0}
	wr.labels = map[person.ID]string{0: "(none)"}
	classreg.AllForClass(r, wr.class.ID(), classreg.FPerson, func(cr *classreg.ClassReg) {
		if cr.Person() != 0 {
			used.Insert(cr.Person())
		}
	})
	person.All(r, personFields, func(p *person.Person) {
		var registrar string
		if used.Has(p.ID()) {
			return
		}
		lastname, _, _ := strings.Cut(p.SortName(), ",")
		switch {
		case wr.classreg.Email != "" && (strings.EqualFold(wr.classreg.Email, p.Email()) || strings.EqualFold(wr.classreg.Email, p.Email2())):
		case wr.classreg.CellPhone != "" && (wr.classreg.CellPhone == p.CellPhone() || wr.classreg.CellPhone == p.HomePhone() || wr.classreg.CellPhone == p.WorkPhone()):
		case strings.EqualFold(wr.classreg.LastName, lastname):
		default:
			return
		}
		wr.Options = append(wr.Options, p.ID())
		if p.ID() == wr.classreg.RegisteredBy.ID() {
			registrar = " (registrar)"
		}
		wr.labels[p.ID()] = fmt.Sprintf("#%d %s %s %s %s %s %s%s", p.ID(), p.InformalName(), p.Email(), p.Email2(), p.CellPhone(), p.HomePhone(), p.WorkPhone(), registrar)
	})
	wr.Options = append(wr.Options, -1)
	wr.labels[-1] = "Create new person"
	wr.LabelFunc = func(_ *request.Request, v person.ID) string { return wr.labels[v] }
}

func (wr *wantPersonRow) Read(r *request.Request) bool {
	wr.Get(r)
	if wr.ShouldEmit(nil) {
		return wr.RadioGroupRow.Read(r)
	}
	return true
}

type statusRow struct {
	form.RadioGroupRow[string]
	personID *person.ID
}

func (sr *statusRow) Read(r *request.Request) bool {
	if !sr.RadioGroupRow.Read(r) {
		return false
	}
	if *sr.ValueP == "inclass" && *sr.personID == 0 {
		sr.Error = "The 'Accepted' option is not valid unless a Person is selected."
		return false
	}
	return true
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func fmtPhone(p *string) bool {
	digits := strings.Map(func(r rune) rune {
		if r < '0' || r > '9' {
			return -1
		}
		return r
	}, *p)
	if len(digits) == 11 && digits[0] == '1' {
		digits = digits[1:]
	}
	switch len(digits) {
	case 0:
		*p = ""
		return true
	case 10:
		*p = digits[0:3] + "-" + digits[3:6] + "-" + digits[6:10]
		return true
	}
	return false
}
