package person

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPerson handles GET /api/people/$id requests (where $id may be "NEW").
func GetPerson(r *util.Request, idstr string) error {
	var (
		person         *model.Person
		canEditDetails bool
		canEditRoles   bool
		canViewContact bool
		out            jwriter.Writer
		attendmap      map[model.EventID]model.AttendanceInfo
		attended       []*model.Event
		individualHeld map[*model.Role]bool
		roles          = map[model.RoleID]bool{}
		wantEdit       = r.FormValue("edit") != ""
	)
	if idstr == "NEW" {
		if !auth.CanCreatePeople(r) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditDetails = true
		canViewContact = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if !auth.CanViewPerson(r, person) {
			return util.Forbidden
		}
		canEditDetails = r.Person == person || auth.IsWebmaster(r)
		canViewContact = canEditDetails || auth.CanViewContactInfo(r, person)
	}
	out.RawString(`{"person":{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"username":`)
	out.String(person.Username)
	out.RawString(`,"informalName":`)
	out.String(person.InformalName)
	out.RawString(`,"formalName":`)
	out.String(person.FormalName)
	out.RawString(`,"sortName":`)
	out.String(person.SortName)
	out.RawString(`,"callSign":`)
	out.String(person.CallSign)
	if canViewContact {
		out.RawString(`,"emails":[`)
		for i, e := range person.Emails {
			if i != 0 {
				out.RawByte(',')
			}
			e.MarshalEasyJSON(&out)
		}
		out.RawString(`],"homeAddress":`)
		person.HomeAddress.MarshalEasyJSON(&out)
		out.RawString(`,"mailAddress":`)
		person.MailAddress.MarshalEasyJSON(&out)
		out.RawString(`,"workAddress":`)
		person.WorkAddress.MarshalEasyJSON(&out)
		out.RawString(`,"cellPhone":`)
		out.String(person.CellPhone)
		out.RawString(`,"homePhone":`)
		out.String(person.HomePhone)
		out.RawString(`,"workPhone":`)
		out.String(person.WorkPhone)
	}
	for _, r := range person.Roles {
		roles[r] = true
	}
	individualHeld = cacheIndividuallyHeldRoles(r.Tx, r.Tx.FetchPeople(), person)
	out.RawString(`,"roles":[`)
	first := true
	for _, role := range r.Tx.FetchRoles() {
		canAssign := auth.CanAssignRole(r, role)
		canEditRoles = canEditRoles || canAssign
		if individualHeld[role] {
			canAssign = false
		}
		if !roles[role.ID] && (!canAssign || !wantEdit) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		if wantEdit {
			out.RawString(`,"canAssign":`)
			out.Bool(canAssign)
			out.RawString(`,"held":`)
			out.Bool(roles[role.ID])
		}
		out.RawByte('}')
	}
	out.RawByte(']')
	if !wantEdit {
		attendmap = r.Tx.FetchAttendanceByPerson(person)
		for eid := range attendmap {
			event := r.Tx.FetchEvent(eid)
			if r.Person == person || auth.CanRecordAttendanceAtEvent(r, event) {
				attended = append(attended, event)
			}
		}
		if len(attended) > 0 {
			sort.Sort(model.EventSort(attended))
			out.RawString(`,"attended":[`)
			for i := len(attended) - 1; i >= 0; i-- {
				if i != len(attended)-1 {
					out.RawByte(',')
				}
				e := attended[i]
				out.RawString(`{"id":`)
				out.Int(int(e.ID))
				out.RawString(`,"date":`)
				out.String(e.Date)
				out.RawString(`,"name":`)
				out.String(e.Name)
				out.RawString(`,"type":`)
				out.String(model.AttendanceTypeNames[attendmap[e.ID].Type])
				out.RawString(`,"minutes":`)
				out.Uint16(attendmap[e.ID].Minutes)
				out.RawByte('}')
			}
			out.RawByte(']')
		}
	}
	out.RawString(`,"canEdit":`)
	out.Bool(canEditDetails || canEditRoles)
	out.RawByte('}')
	if wantEdit {
		out.RawString(`,"canEditRoles":`)
		out.Bool(canEditRoles)
		out.RawString(`,"canEditDetails":`)
		out.Bool(canEditDetails)
		out.RawString(`,"allowBadPassword":`)
		out.Bool(auth.IsWebmaster(r))
		out.RawString(`,"canEditUsername":`)
		out.Bool(auth.IsWebmaster(r))
		if canEditDetails {
			out.RawString(`,"passwordHints":[`)
			for i, h := range auth.SERVPasswordHints {
				if i != 0 {
					out.RawByte(',')
				}
				out.String(h)
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPerson handles POST /api/people/$id requests (where $id may be "NEW").
func PostPerson(r *util.Request, idstr string) error {
	var (
		person         *model.Person
		canEditDetails bool
		err            error
	)
	if idstr == "NEW" {
		if !auth.CanCreatePeople(r) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditDetails = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		canEditDetails = r.Person == person || auth.IsWebmaster(r)
	}
	if !canEditDetails && !auth.CanAssignAnyRole(r) {
		return util.Forbidden
	}
	// Remove all roles that the user is allowed to change; keep the ones
	// that they aren't.
	j := 0
	for _, role := range person.Roles {
		if !auth.CanAssignRole(r, r.Tx.FetchRole(role)) {
			person.Roles[j] = role
			j++
		}
	}
	person.Roles = person.Roles[:j]
	// Add roles that the user requested.
	for _, ridstr := range r.Form["role"] {
		if role := r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); role != nil && auth.CanAssignRole(r, role) {
			person.Roles = append(person.Roles, role.ID)
		} else {
			return errors.New("bad role")
		}
	}
	if canEditDetails {
		person.InformalName = r.FormValue("informalName")
		person.FormalName = r.FormValue("formalName")
		person.SortName = r.FormValue("sortName")
		person.Username = r.FormValue("username")
		person.CallSign = r.FormValue("callSign")
		person.Emails = person.Emails[:0]
		for i, e := range r.Form["email"] {
			var email model.PersonEmail
			email.Email = e
			if len(r.Form["emailLabel"]) > i {
				email.Label = r.Form["emailLabel"][i]
			}
			person.Emails = append(person.Emails, &email)
		}
		person.CellPhone = r.FormValue("cellPhone")
		person.HomePhone = r.FormValue("homePhone")
		person.WorkPhone = r.FormValue("workPhone")
		person.HomeAddress.Address = r.FormValue("homeAddress")
		if l := r.FormValue("homeAddressLatitude"); l != "" {
			if person.HomeAddress.Latitude, err = strconv.ParseFloat(l, 64); err != nil {
				return errors.New("invalid latitude")
			}
		}
		if l := r.FormValue("homeAddressLongitude"); l != "" {
			if person.HomeAddress.Longitude, err = strconv.ParseFloat(l, 64); err != nil {
				return errors.New("invalid longitude")
			}
		}
		person.MailAddress.Address = r.FormValue("mailAddress")
		person.MailAddress.SameAsHome, _ = strconv.ParseBool(r.FormValue("mailAddressSameAsHome"))
		person.WorkAddress.Address = r.FormValue("workAddress")
		if l := r.FormValue("workAddressLatitude"); l != "" {
			if person.WorkAddress.Latitude, err = strconv.ParseFloat(l, 64); err != nil {
				return errors.New("invalid latitude")
			}
		}
		if l := r.FormValue("workAddressLongitude"); l != "" {
			if person.WorkAddress.Longitude, err = strconv.ParseFloat(l, 64); err != nil {
				return errors.New("invalid longitude")
			}
		}
		person.WorkAddress.SameAsHome, _ = strconv.ParseBool(r.FormValue("workAddressSameAsHome"))
	}
	if err = ValidatePerson(r.Tx, person); err != nil {
		if estr := err.Error(); strings.HasPrefix(estr, "duplicate ") {
			// These need to be sent back to the client as 200
			// responses with error details.
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			fmt.Fprintf(r, `{"duplicate%s":true}`, strings.ToUpper(estr[10:11])+estr[11:])
			return nil
		}
		return err
	}
	// We do password after validation so that it can use the other
	// fields as password hints.
	if password := r.FormValue("password"); password != "" && canEditDetails {
		if !auth.IsWebmaster(r) {
			if !auth.StrongPassword(person, password) {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"weakPassword":true}`))
				return nil
			}
		}
		auth.SetPassword(r, person, password)
	}
	r.Tx.SavePerson(person)
	r.Tx.Commit()
	return nil
}
