package classes

import (
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleRegister handles /classes/$id/register requests.
func HandleRegister(r *request.Request, cidstr string) {
	const personFields = person.FID | person.FInformalName | person.FSortName | person.FEmail | person.FEmail2 | person.FCellPhone
	var (
		user        *person.Person
		c           *class.Class
		langfield   class.Fields
		regs        []*classreg.ClassReg
		uregs       []*classreg.Updater
		askReferral bool
	)
	if user = auth.SessionUser(r, personFields, false); user == nil {
		handleRegisterNotLoggedIn(r, cidstr)
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if r.Language == "es" {
		langfield = class.FEsDesc
	} else {
		langfield = class.FEnDesc
	}
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), class.FLimit|class.FReferrals|class.FStart|class.FType|langfield); c == nil {
		errpage.NotFound(r, user)
		return
	}
	classreg.AllForClass(r, c.ID(), classreg.UpdaterFields, func(cr *classreg.ClassReg) {
		if cr.RegisteredBy() != user.ID() {
			return
		}
		regs = append(regs, cr.Clone())
		uregs = append(uregs, cr.Updater(r, c, nil, user))
	})
	if len(regs) == 0 {
		uregs = append(uregs, &classreg.Updater{
			Class:        c,
			Person:       user,
			RegisteredBy: user,
			FirstName:    personFirstName(user),
			LastName:     personLastName(user),
			Email:        personEmail(user),
			CellPhone:    user.CellPhone(),
		})
		regs = append(regs, nil)
		askReferral = r.Method == http.MethodGet
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	for i := 0; i < len(uregs); i++ {
		emitRow(r, form, uregs[i], i)
	}
	emitRow(r, form, new(classreg.Updater), len(uregs))
	if askReferral {
		emitReferral(r, form)
	}
	emitButtons(r, form)
}

func handleRegisterNotLoggedIn(r *request.Request, cidstr string) {
	panic("not implemented")
}

func emitRow(r *request.Request, form *htmlb.Element, reg *classreg.Updater, idx int) {
	div := form.E("div class='formRow-3col classregDivider'", idx == 0, "class=first")
	div.E("div").TF(r.Loc("Student %d"), idx+1)
	div.E("button type=button class='sbtn sbtn-xsmall sbtn-danger classregClear' data-row=%d>%s", idx, r.Loc("Clear"))
	row := form.E("div class=formRow")
	row.E("label for=classregFirstname%d>%s", idx, r.Loc("Name"))
	names := row.E("div class='formInput classregNames'")
	names.E("input id=classregFirstname%d name=firstName class='formInput classregFirstname' placeholder=%s value=%s", idx, r.Loc("First"), reg.FirstName)
	names.E("input id=classregLastname%d name=lastName class='formInput classregLastname' placeholder=%s value=%s", idx, r.Loc("Last"), reg.LastName)
	row = form.E("div class=formRow")
	row.E("label for=classregEmail%d>%s", idx, r.Loc("Email"))
	row.E("input id=classregEmail%d name=email class='formInput classregEmail' value=%s", idx, reg.Email)
	row = form.E("div class=formRow")
	row.E("label for=classregCellPhone%d>%s", idx, r.Loc("Cell Phone"))
	row.E("input id=classregCellPhone%d name=cellPhone class='formInput classregCellPhone' value=%s", idx, reg.CellPhone)
}

func personFirstName(p *person.Person) string {
	parts := strings.SplitN(p.SortName(), ",", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func personLastName(p *person.Person) string {
	parts := strings.Split(p.SortName(), ",")
	return parts[0]
}

func personEmail(p *person.Person) string {
	if p.Email() != "" {
		return p.Email()
	}
	return p.Email2()
}

func emitReferral(r *request.Request, form *htmlb.Element) {
	row := form.E("div class='formRow-3col classregReferral'")
	row.E("label for=classregReferral>%s", r.Loc("How did you find out about this class?"))
	sel := row.E("select id=classregReferral name=referral class=formInput")
	sel.E("option value=''>%s", r.Loc("(select one)"))
	for _, ref := range class.AllReferrals {
		sel.E("option value=%d>%s", ref, r.Loc(ref.String()))
	}
}

func emitButtons(r *request.Request, form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Sign Up"))
}
