package person

import (
	"errors"
	"strconv"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/api/authn"
	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPerson handles GET /api/people/$id requests (where $id may be "NEW").
func GetPerson(r *util.Request, idstr string) error {
	var (
		person         *model.Person
		canEditDetails bool
		canViewContact bool
		out            jwriter.Writer
	)
	if idstr == "NEW" {
		if !r.Person.HasPrivLevel(model.PrivLeader) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditDetails = true
		canViewContact = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		var canView bool
		if canView, canViewContact = canViewPerson(r.Person, person); !canView {
			return util.Forbidden
		}
		canEditDetails = r.Person == person || r.Person.HasPrivLevel(model.PrivLeader)
		canViewContact = canEditDetails || canViewContact
	}
	out.RawString(`{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"informalName":`)
	out.String(person.InformalName)
	out.RawString(`,"formalName":`)
	out.String(person.FormalName)
	out.RawString(`,"sortName":`)
	out.String(person.SortName)
	out.RawString(`,"callSign":`)
	out.String(person.CallSign)
	if canViewContact {
		out.RawString(`,"contact":{"email":`)
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
		out.RawByte('}')
	}
	out.RawString(`,"roles":[`)
	var first = true
	for _, role := range r.Tx.FetchRoles() {
		if !person.Roles[role.ID] || role.Title == "" {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"title":`)
		out.String(role.Title)
		out.RawString(`,"org":`)
		out.String(role.Org.String())
		out.RawByte('}')
	}
	out.RawByte(']')
	if r.Person == person || r.Person.HasPrivLevel(model.PrivLeader) {
		out.RawString(`,"lists":[`)
		var first = true
		for _, l := range r.Tx.FetchLists() {
			if lps, ok := l.People[person.ID]; !ok || lps&model.ListSubscribed == 0 {
				continue
			}
			if l.Type == model.ListEmail && person.NoEmail {
				continue
			}
			if l.Type == model.ListSMS && person.NoText {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			if l.Type == model.ListEmail {
				out.String(l.Name + "@SunnyvaleSERV.org")
			} else {
				out.String("SMS: " + l.Name)
			}
		}
		out.RawByte(']')
	}
	if person == r.Person || r.Person.HasPrivLevel(model.PrivLeader) {
		switch {
		case r.Person.IsAdminLeader():
			out.RawString(`,"status":{"canEdit":true,"level":"admin","volgistics":{"needed":`)
		case r.Person.HasPrivLevel(model.PrivLeader):
			out.RawString(`,"status":{"canEdit":true,"level":"leader","volgistics":{"needed":`)
		default:
			out.RawString(`,"status":{"level":"self","volgistics":{"needed":`)
		}
		out.Bool(person.HasPrivLevel(model.PrivMember))
		out.RawString(`,"id":`)
		out.Int(int(person.VolgisticsID))
		out.RawByte('}')
		for _, c := range model.AllDSWClasses {
			out.RawString(`,"dsw`)
			out.RawString(model.DSWClassNames[c][:4])
			out.RawString(`":{"needed":`)
			out.Bool(needDSW(r, person, c, nil))
			if person.DSWRegistrations != nil && !person.DSWRegistrations[c].IsZero() {
				out.RawString(`,"registered":`)
				out.String(person.DSWRegistrations[c].Format("2006-01-02"))
				out.RawString(`,"expires":`)
				out.String(person.DSWUntil[c].Format("2006-01-02"))
				if person.DSWUntil[c].Before(time.Now()) {
					out.RawString(`,"expired":true`)
				}
			}
			out.RawByte('}')
		}
		out.RawString(`,"backgroundCheck":{"needed":`)
		out.Bool(person.HasPrivLevel(model.PrivMember))
		switch person.BackgroundCheck {
		case "":
			break
		case "true":
			out.RawString(`,"cleared":true`)
		default:
			out.RawString(`,"cleared":`)
			out.String(person.BackgroundCheck)
		}
		first := true
		out.RawString(`},"identification":[`)
		for _, t := range model.AllIdentTypes {
			if person.Identification&t == 0 {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(model.IdentTypeNames[t])
		}
		out.RawString(`]}`)
	}
	out.RawString(`,"notes":[`)
	first = true
	for _, n := range person.Notes {
		// Need to redo how notes are privileged.  A few had
		// a "contact" privilege and should be more visible than
		// this.  TODO
		if !r.Person.IsAdminLeader() {
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
	out.RawString(`],"canEdit":`)
	out.Bool(canEditDetails)
	out.RawString(`,"canEditRoles":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader) && person.ID != model.AdminPersonID)
	out.RawString(`,"canEditLists":`)
	out.Bool(person == r.Person || r.Person.Roles[model.Webmaster])
	out.RawString(`,"canChangePassword":`)
	out.Bool(person == r.Person || r.Person.Roles[model.Webmaster])
	out.RawString(`,"canHours":`)
	out.Bool(person.ID == r.Person.ID || r.Person.IsAdminLeader())
	out.RawString(`,"noEmail":`)
	out.Bool(person.NoEmail)
	out.RawString(`,"noText":`)
	out.Bool(person.NoText)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// needVolgisticsID returns whether the person is in a group from which
// volunteer hours are requested.
func needVolgisticsID(r *util.Request, p *model.Person, role *model.Role) bool {
	if role != nil {
		return p.Orgs[role.Org].PrivLevel >= model.PrivMember
	}
	return p.HasPrivLevel(model.PrivMember)
}

// needDSW returns whether the person is in a group that requires DSW clearance
// for the specified class.
func needDSW(r *util.Request, p *model.Person, c model.DSWClass, role *model.Role) bool {
	if role != nil {
		return role.Org.DSWClass() == c
	}
	for o, om := range p.Orgs {
		if om.PrivLevel >= model.PrivMember && model.Org(o).DSWClass() == c {
			return true
		}
	}
	return false
}

// GetPersonContact handles GET /api/people/$id/contact requests.
func GetPersonContact(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	out.RawString(`{"id":`)
	out.Int(int(person.ID))
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
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonContact handles POST /api/people/$id/contact requests.
func PostPersonContact(r *util.Request, idstr string) error {
	var (
		person *model.Person
		err    error
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
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
	switch err = ValidatePerson(r.Tx, person); err {
	case nil:
		break
	case errDuplicateEmail:
		return util.SendConflict(r, "email")
	case errDuplicateCellPhone:
		return util.SendConflict(r, "cellPhone")
	default:
		return err
	}
	r.Tx.UpdatePerson(person)
	r.Tx.Commit()
	return nil
}

// GetPersonLists handles GET /api/people/${id}/lists requests.
func GetPersonLists(r *util.Request, idstr string) error {
	var (
		person       *model.Person
		canSubscribe map[model.ListID]bool
		out          jwriter.Writer
	)
	// idstr may be either a person ID, in string format, or an unsubscribe
	// token.  We can distinguish between them by the length.
	if len(idstr) <= 5 {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if person != r.Person && !r.Person.Roles[model.Webmaster] {
			return util.Forbidden
		}
	} else {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	}
	canSubscribe = authz.CanSubscribe(r.Tx, person)
	out.RawString(`{"noEmail":`)
	out.Bool(person.NoEmail)
	out.RawString(`,"noText":`)
	out.Bool(person.NoText)
	out.RawString(`,"lists":[`)
	var first = true
	for _, list := range r.Tx.FetchLists() {
		if !canSubscribe[list.ID] {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(list.ID))
		out.RawString(`,"type":`)
		out.String(model.ListTypeNames[list.Type])
		out.RawString(`,"name":`)
		if list.Type == model.ListEmail {
			out.String(list.Name + "@SunnyvaleSERV.org")
		} else {
			out.String("SMS: " + list.Name)
		}
		lps := list.People[person.ID]
		out.RawString(`,"subscribed":`)
		out.Bool(lps&model.ListSubscribed != 0)
		var firstWarn = true
		out.RawString(`,"subWarn":[`)
		for rid := range person.Roles {
			role := r.Tx.FetchRole(rid)
			if role.Lists[list.ID].SubModel() == model.ListWarnUnsub {
				if firstWarn {
					firstWarn = false
				} else {
					out.RawByte(',')
				}
				if role.Title != "" {
					out.String(role.Title)
				} else {
					out.String(role.Name)
				}
			}
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonLists handles POST /api/people/${id}/lists requests.
func PostPersonLists(r *util.Request, idstr string) error {
	var (
		person       *model.Person
		noEmail      bool
		noText       bool
		canSubscribe map[model.ListID]bool
		subscribed   = make(map[model.ListID]bool)
	)
	// idstr may be either a person ID, in string format, or an unsubscribe
	// token.  We can distinguish between them by the length.
	if len(idstr) <= 5 {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if person != r.Person && !r.Person.Roles[model.Webmaster] {
			return util.Forbidden
		}
	} else {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	}
	noEmail, _ = strconv.ParseBool(r.FormValue("noEmail"))
	noText, _ = strconv.ParseBool(r.FormValue("noText"))
	if noEmail != person.NoEmail || noText != person.NoText {
		r.Tx.WillUpdatePerson(person)
		person.NoEmail, person.NoText = noEmail, noText
		r.Tx.UpdatePerson(person)
	}
	canSubscribe = authz.CanSubscribe(r.Tx, person)
	for _, lidstr := range r.Form["list"] {
		list := r.Tx.FetchList(model.ListID(util.ParseID(lidstr)))
		if list == nil {
			return errors.New("nonexistent list")
		}
		if !canSubscribe[list.ID] {
			return errors.New("forbidden list")
		}
		subscribed[list.ID] = true
	}
	for lid := range canSubscribe {
		list := r.Tx.FetchList(lid)
		if subscribed[list.ID] && list.People[person.ID]&model.ListSubscribed == 0 {
			r.Tx.WillUpdateList(list)
			list.People[person.ID] = (list.People[person.ID] &^ model.ListUnsubscribed) | model.ListSubscribed
			r.Tx.UpdateList(list)
		}
		if !subscribed[list.ID] && list.People[person.ID]&model.ListUnsubscribed == 0 {
			r.Tx.WillUpdateList(list)
			list.People[person.ID] = (list.People[person.ID] &^ model.ListSubscribed) | model.ListUnsubscribed
			r.Tx.UpdateList(list)
		}
	}
	r.Tx.Commit()
	return nil
}

// GetPersonNames handles GET /api/people/$id/names requests.
func GetPersonNames(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	out.RawString(`{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"informalName":`)
	out.String(person.InformalName)
	out.RawString(`,"formalName":`)
	out.String(person.FormalName)
	out.RawString(`,"sortName":`)
	out.String(person.SortName)
	out.RawString(`,"callSign":`)
	out.String(person.CallSign)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonNames handles POST /api/people/$id/names requests.
func PostPersonNames(r *util.Request, idstr string) error {
	var (
		person *model.Person
		err    error
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
	person.InformalName = r.FormValue("informalName")
	person.FormalName = r.FormValue("formalName")
	person.SortName = r.FormValue("sortName")
	person.CallSign = r.FormValue("callSign")
	switch err = ValidatePerson(r.Tx, person); err {
	case nil:
		break
	case errDuplicateSortName:
		return util.SendConflict(r, "sortName")
	case errDuplicateCallSign:
		return util.SendConflict(r, "callSign")
	default:
		return err
	}
	r.Tx.UpdatePerson(person)
	r.Tx.Commit()
	return nil
}

// GetPersonPassword handles GET /api/people/$id/password requests.
func GetPersonPassword(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	out.RawString(`[`)
	for i, h := range authn.SERVPasswordHints {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(h)
	}
	out.RawByte(',')
	out.String(person.InformalName)
	out.RawByte(',')
	out.String(person.FormalName)
	out.RawByte(',')
	out.String(person.CallSign)
	out.RawByte(',')
	out.String(person.Email)
	out.RawByte(',')
	out.String(person.Email2)
	out.RawByte(',')
	out.String(person.HomeAddress.Address)
	out.RawByte(',')
	out.String(person.MailAddress.Address)
	out.RawByte(',')
	out.String(person.WorkAddress.Address)
	out.RawByte(',')
	out.String(person.CellPhone)
	out.RawByte(',')
	out.String(person.HomePhone)
	out.RawByte(',')
	out.String(person.WorkPhone)
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonPassword handles POST /api/people/$id/password requests.
func PostPersonPassword(r *util.Request, idstr string) error {
	var person *model.Person

	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if r.Person != person && !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
	if password := r.FormValue("password"); password != "" {
		if !r.Person.Roles[model.Webmaster] {
			if oldPassword := r.FormValue("oldPassword"); !authn.CheckPassword(person, oldPassword) {
				return util.SendConflict(r, "wrongOldPassword")
			}
			if !authn.StrongPassword(person, password) {
				return util.SendConflict(r, "weakPassword")
			}
		}
		authn.SetPassword(r, person, password)
	}
	r.Tx.UpdatePerson(person)
	r.Tx.Commit()
	return nil
}

// GetPersonRoles handles GET /api/people/${id}/roles requests.
func GetPersonRoles(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if person.ID == model.AdminPersonID {
		return util.Forbidden
	}
	out.RawByte('[')
	var first = true
	for _, org := range model.AllOrgs {
		if r.Person.Orgs[org].PrivLevel != model.PrivLeader &&
			r.Person.Orgs[model.OrgAdmin].PrivLevel != model.PrivLeader {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"org":`)
		out.String(org.String())
		out.RawString(`,"roles":[`)
		var firstRole = true
		for _, role := range r.Tx.FetchRoles() {
			if role.Org != org {
				continue
			}
			if firstRole {
				firstRole = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(role.ID))
			out.RawString(`,"name":`)
			if role.Title != "" {
				out.String(role.Title)
			} else {
				out.String(role.Name)
			}
			direct, held := person.Roles[role.ID]
			out.RawString(`,"held":`)
			out.Bool(held)
			out.RawString(`,"direct":`)
			out.Bool(direct)
			if role.ImplicitOnly {
				out.RawString(`,"implicitOnly":true`)
			}
			out.RawString(`,"implies":[`)
			var firstImplies = true
			for irid := range role.Implies {
				if firstImplies {
					firstImplies = false
				} else {
					out.RawByte(',')
				}
				out.Int(int(irid))
			}
			out.RawString(`]}`)
		}
		out.RawString(`]}`)
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonRoles handles POST /api/people/${id}/roles requests.
func PostPersonRoles(r *util.Request, idstr string) error {
	var person *model.Person

	if !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if person.ID == model.AdminPersonID {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
	for rid, direct := range person.Roles {
		if !direct {
			delete(person.Roles, rid)
			continue
		}
		role := r.Tx.FetchRole(rid)
		if r.Person.Orgs[role.Org].PrivLevel < model.PrivLeader && r.Person.Orgs[model.OrgAdmin].PrivLevel < model.PrivLeader {
			continue
		}
		delete(person.Roles, rid)
	}
	r.ParseMultipartForm(1048576)
	for _, ridstr := range r.Form["role"] {
		role := r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr)))
		if role == nil {
			return errors.New("invalid role")
		}
		if r.Person.Orgs[role.Org].PrivLevel < model.PrivLeader && r.Person.Orgs[model.OrgAdmin].PrivLevel < model.PrivLeader {
			return errors.New("forbidden role")
		}
		person.Roles[role.ID] = true
	}
	if err := ValidatePerson(r.Tx, person); err != nil {
		return err
	}
	r.Tx.UpdatePerson(person)
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}

// GetPersonStatus handles GET /api/people/$id/status requests.
func GetPersonStatus(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if !r.Person.IsAdminLeader() {
		return util.Forbidden
	}
	out.RawString(`{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"volgistics":`)
	out.Int(person.VolgisticsID)
	for _, c := range model.AllDSWClasses {
		out.RawString(`,"dsw`)
		out.RawString(model.DSWClassNames[c][:4])
		out.RawString(`":{"needed":`)
		out.Bool(needDSW(r, person, c, nil))
		if person.DSWRegistrations != nil && !person.DSWRegistrations[c].IsZero() {
			out.RawString(`,"registered":`)
			out.String(person.DSWRegistrations[c].Format("2006-01-02"))
			out.RawString(`,"expires":`)
			out.String(person.DSWUntil[c].Format("2006-01-02"))
			if person.DSWUntil[c].Before(time.Now()) {
				out.RawString(`,"expired":true`)
			}
		} else {
			out.RawString(`,"registered":""`)
		}
		out.RawByte('}')
	}
	out.RawString(`,"backgroundCheck":{"needed":`)
	out.Bool(person.HasPrivLevel(model.PrivMember))
	switch person.BackgroundCheck {
	case "":
		out.RawString(`,"cleared":""`)
	case "true":
		out.RawString(`,"cleared":true`)
	default:
		out.RawString(`,"cleared":`)
		out.String(person.BackgroundCheck)
	}
	out.RawString(`},"identification":[`)
	for i, t := range model.AllIdentTypes {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"type":`)
		out.String(model.IdentTypeNames[t])
		out.RawString(`,"held":`)
		out.Bool(person.Identification&t != 0)
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonStatus handles POST /api/people/$id/status requests.
func PostPersonStatus(r *util.Request, idstr string) error {
	var (
		person *model.Person
		err    error
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if !r.Person.IsAdminLeader() {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
	person.VolgisticsID, _ = strconv.Atoi(r.FormValue("volgistics"))
	if person.DSWRegistrations == nil {
		person.DSWRegistrations = make(map[model.DSWClass]time.Time)
	}
	for _, c := range model.AllDSWClasses {
		date := r.FormValue("dsw" + model.DSWClassNames[c][:4])
		if date == "" {
			delete(person.DSWRegistrations, c)
		} else if person.DSWRegistrations[c], err = time.ParseInLocation("2006-01-02", date, time.Local); err != nil {
			return errors.New("invalid DSW date")
		}
	}
	person.BackgroundCheck = r.FormValue("backgroundCheck")
	person.Identification = 0
IDENTS:
	for _, n := range r.Form["identification"] {
		for t, tn := range model.IdentTypeNames {
			if n == tn {
				person.Identification |= t
				continue IDENTS
			}
		}
		return errors.New("invalid identification type")
	}
	if err = ValidatePerson(r.Tx, person); err != nil {
		return err
	}
	r.Tx.UpdatePerson(person)
	r.Tx.Commit()
	return nil
}
