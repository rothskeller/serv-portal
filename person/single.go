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
		individualHeld map[model.RoleID]model.PersonID
		roles          = map[model.RoleID]bool{}
		wantEdit       = r.FormValue("edit") != ""
	)
	if idstr == "NEW" {
		if !r.Auth.CanA(model.PrivManageMembers) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditDetails = true
		canViewContact = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if !r.Auth.CanAP(model.PrivViewMembers, person.ID) {
			return util.Forbidden
		}
		canEditDetails = r.Person == person || r.Auth.IsWebmaster()
		canViewContact = canEditDetails || r.Auth.CanAP(model.PrivViewContactInfo, person.ID)
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
		out.RawString(`,"email":`)
		out.String(person.Email)
		out.RawString(`,"email2":`)
		out.String(person.Email2)
		out.RawString(`,"homeAddress":`)
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
	for _, r := range r.Auth.RolesP(person.ID) {
		roles[r] = true
	}
	individualHeld = cacheIndividuallyHeldRoles(r.Auth, person.ID)
	out.RawString(`,"roles":[`)
	first := true
	for _, role := range r.Auth.FetchRoles(r.Auth.AllRoles()) {
		var canAssign = r.Auth.CanAR(model.PrivManageMembers, role.ID)
		canEditRoles = canEditRoles || canAssign
		if individualHeld[role.ID] != 0 {
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
			if r.Person == person {
				attended = append(attended, event)
			} else {
				for _, group := range event.Groups {
					if r.Auth.CanAG(model.PrivManageEvents, group) {
						attended = append(attended, event)
						break
					}
				}
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
		out.RawString(`,"notes":[`)
		first := true
		for _, n := range person.Notes {
			if n.Privilege == 0 && !r.Auth.IsWebmaster() {
				continue
			}
			if n.Privilege != 0 && !r.Auth.CanAP(n.Privilege, person.ID) {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"date":`)
			out.String(n.Date)
			out.RawString(`,"note":`)
			out.String(n.Note)
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	out.RawString(`,"canEdit":`)
	out.Bool(canEditDetails || canEditRoles)
	out.RawString(`,"noEmail":`)
	out.Bool(person.NoEmail)
	out.RawString(`,"noText":`)
	out.Bool(person.NoText)
	out.RawByte('}')
	if wantEdit {
		out.RawString(`,"canEditRoles":`)
		out.Bool(canEditRoles)
		out.RawString(`,"canEditDetails":`)
		out.Bool(canEditDetails)
		out.RawString(`,"allowBadPassword":`)
		out.Bool(r.Auth.IsWebmaster())
		out.RawString(`,"canEditUsername":`)
		out.Bool(r.Auth.IsWebmaster())
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
		roles          []model.RoleID
		err            error
	)
	if idstr == "NEW" {
		if !r.Auth.CanA(model.PrivManageMembers) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditDetails = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		r.Tx.WillUpdatePerson(person)
		canEditDetails = r.Person == person || r.Auth.IsWebmaster()
	}
	if !canEditDetails && !r.Auth.CanA(model.PrivManageMembers) {
		return util.Forbidden
	}
	// Remove all roles that the user is allowed to change; keep the ones
	// that they aren't.
	roles = r.Auth.RolesP(person.ID)
	j := 0
	for _, role := range roles {
		if !r.Auth.CanAR(model.PrivManageMembers, role) {
			roles[j] = role
			j++
		}
	}
	roles = roles[:j]
	// Add roles that the user requested.
	r.ParseMultipartForm(1048576)
	for _, ridstr := range r.Form["role"] {
		if role := r.Auth.FetchRole(model.RoleID(util.ParseID(ridstr))); role != nil && r.Auth.CanAR(model.PrivManageMembers, role.ID) {
			roles = append(roles, role.ID)
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
		person.Email = r.FormValue("email")
		person.Email2 = r.FormValue("email2")
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
	if err = ValidatePerson(r.Tx, person, roles); err != nil {
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
		if !r.Auth.IsWebmaster() {
			if !auth.StrongPassword(person, password) {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"weakPassword":true}`))
				return nil
			}
		}
		auth.SetPassword(r, person, password)
	}
	if person.ID == 0 {
		r.Tx.CreatePerson(person)
	} else {
		r.Tx.UpdatePerson(person)
	}
	r.Auth.SetPersonRoles(person.ID, roles)
	r.Auth.Save()
	r.Tx.Commit()
	return nil
}
