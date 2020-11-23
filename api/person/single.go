package person

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
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
		if !r.Auth.CanAP(model.PrivViewMembers, person.ID) && person != r.Person {
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
	out.RawString(`],"roles2":[`)
	first = true
	for _, role := range r.Tx.FetchRoles() {
		if !person.Roles[role.ID] || role.Title == "" {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.String(role.Title)
	}
	out.RawByte(']')
	if r.Person == person || r.Auth.IsWebmaster() {
		canSubscribe := authz.CanSubscribe(r.Tx, person)
		if len(canSubscribe) != 0 {
			out.RawString(`,"lists":[`)
			var first = true
			for _, l := range r.Tx.FetchLists() {
				if lps, ok := l.People[person.ID]; !ok || lps&model.ListSubscribed == 0 {
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
	}
	if r.Auth.May(model.PermViewClearances) {
		out.RawString(`,"identification":[`)
		first = true
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
		out.RawByte(']')
	}
	if wantEdit {
		if r.Auth.May(model.PermEditClearances) {
			out.RawString(`,"volgistics":`)
			out.Int(person.VolgisticsID)
			out.RawString(`,"dsw":{`)
			for i, c := range model.AllDSWClasses {
				if i != 0 {
					out.RawByte(',')
				}
				out.String(model.DSWClassNames[c])
				out.RawByte(':')
				if person.DSWRegistrations == nil || person.DSWRegistrations[c].IsZero() {
					out.RawString(`""`)
				} else {
					out.String(person.DSWRegistrations[c].Format("2006-01-02"))
				}
			}
			out.RawString(`},"backgroundCheck":`)
			out.String(person.BackgroundCheck)
		}
	} else {
		if r.Auth.May(model.PermViewClearances) {
			if person.VolgisticsID != 0 {
				out.RawString(`,"volgistics":`)
				out.Int(person.VolgisticsID)
			} else if needVolgisticsID(r, person, nil) {
				out.RawString(`,"volgistics":false`)
			}
			if person.BackgroundCheck != "" {
				out.RawString(`,"backgroundCheck":`)
				out.String(person.BackgroundCheck)
			}
		} else if r.Person == person {
			if person.VolgisticsID != 0 {
				out.RawString(`,"volgistics":true`)
			} else if needVolgisticsID(r, person, nil) {
				out.RawString(`,"volgistics":false`)
			}
		}
		if r.Person == person || r.Auth.May(model.PermViewClearances) {
			out.RawString(`,"dsw":{`)
			var first = true
			for _, c := range model.AllDSWClasses {
				needed := needDSW(r, person, c, nil)
				if (person.DSWRegistrations == nil || person.DSWRegistrations[c].IsZero()) && !needed {
					continue
				}
				if first {
					first = false
				} else {
					out.RawByte(',')
				}
				out.String(model.DSWClassNames[c])
				out.RawString(`:{`)
				if person.DSWRegistrations == nil || person.DSWRegistrations[c].IsZero() {
					out.RawString(`"needed":true`)
				} else {
					out.RawString(`"registered":`)
					out.String(person.DSWRegistrations[c].Format("2006-01-02"))
					out.RawString(`,"expires":`)
					out.String(person.DSWUntil[c].Format("2006-01-02"))
					if person.DSWUntil[c].Before(time.Now()) {
						out.RawString(`,"expired":true`)
					}
				}
				out.RawByte('}')
			}
			out.RawByte('}')
			switch person.BackgroundCheck {
			case "":
				if needBackgroundCheck(r, person, nil) && r.Auth.IsWebmaster() {
					// Setting this to webmaster only until we have accurate BG check data.
					out.RawString(`,"backgroundCheck":false`)
				}
			case "true":
				out.RawString(`,"backgroundCheck":true`)
			default:
				out.RawString(`,"backgroundCheck":`)
				out.String(person.BackgroundCheck)
			}
		}
		attendmap = r.Tx.FetchAttendanceByPerson(person)
		for eid := range attendmap {
			if attendmap[eid].Type == model.AttendAsAbsent {
				continue
			}
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
	out.RawString(`,"canEditRoles2":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader) && person.Username != "admin")
	out.RawString(`,"canHours":`)
	out.Bool(person.ID == r.Person.ID || r.Auth.IsWebmaster())
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
		out.RawString(`,"canEditClearances":`)
		out.Bool(r.Auth.May(model.PermEditClearances))
		out.RawString(`,"allowBadPassword":`)
		out.Bool(r.Auth.IsWebmaster())
		out.RawString(`,"canEditUsername":`)
		out.Bool(r.Auth.IsWebmaster())
		out.RawString(`,"identTypes":[`)
		for i, t := range model.AllIdentTypes {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(model.IdentTypeNames[t])
		}
		out.RawByte(']')
		if canEditDetails {
			out.RawString(`,"passwordHints":[`)
			for i, h := range authn.SERVPasswordHints {
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
		person.UnsubscribeToken = util.RandomToken()
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
	// If there are no resulting roles, add the disabled role.
	if len(roles) == 0 {
		roles = append(roles, r.Auth.FetchRoleByTag(model.RoleDisabled).ID)
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
	if r.Auth.May(model.PermEditClearances) {
		person.VolgisticsID, _ = strconv.Atoi(r.FormValue("volgistics"))
		if person.DSWRegistrations == nil {
			person.DSWRegistrations = make(map[model.DSWClass]time.Time)
		}
		for _, c := range model.AllDSWClasses {
			date := r.FormValue("dsw-" + model.DSWClassNames[c])
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
			if oldPassword := r.FormValue("oldPassword"); !authn.CheckPassword(person, oldPassword) {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"wrongOldPassword":true}`))
				return nil
			}
			if !authn.StrongPassword(person, password) {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"weakPassword":true}`))
				return nil
			}
		}
		authn.SetPassword(r, person, password)
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

// needVolgisticsID returns whether the person is in a group from which
// volunteer hours are requested.
func needVolgisticsID(r *util.Request, p *model.Person, g *model.Group) bool {
	if g != nil {
		return g.GetHours
	}
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsP(p.ID)) {
		if g.GetHours {
			return true
		}
	}
	return false
}

// needDSW returns whether the person is in a group that requires DSW clearance
// for the specified class.
func needDSW(r *util.Request, p *model.Person, c model.DSWClass, g *model.Group) bool {
	if g != nil {
		return model.OrganizationToDSWClass[g.Organization] == c
	}
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsP(p.ID)) {
		if model.OrganizationToDSWClass[g.Organization] == c {
			return true
		}
	}
	return false
}

// needBackgroundCheck returns whether the person is in a group that requires a
// background check.
func needBackgroundCheck(r *util.Request, p *model.Person, g *model.Group) bool {
	if g != nil {
		return g.BackgroundCheckRequired
	}
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsP(p.ID)) {
		if g.BackgroundCheckRequired {
			return true
		}
	}
	return false
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
		if person != r.Person && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
	} else {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	}
	if canSubscribe = authz.CanSubscribe(r.Tx, person); len(canSubscribe) == 0 {
		return util.Forbidden
	}
	out.RawByte('[')
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
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPersonLists handles POST /api/people/${id}/lists requests.
func PostPersonLists(r *util.Request, idstr string) error {
	var (
		person       *model.Person
		canSubscribe map[model.ListID]bool
		subscribed   = make(map[model.ListID]bool)
	)
	// idstr may be either a person ID, in string format, or an unsubscribe
	// token.  We can distinguish between them by the length.
	if len(idstr) <= 5 {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if person != r.Person && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
	} else {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	}
	if canSubscribe = authz.CanSubscribe(r.Tx, person); len(canSubscribe) == 0 {
		return util.Forbidden
	}
	r.ParseMultipartForm(1048576)
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
	if person.Username == "admin" {
		return util.Forbidden
	}
	out.RawByte('[')
	var first = true
	for _, org := range model.AllOrgs {
		if r.Person.Orgs[org].PrivLevel != model.PrivLeader &&
			r.Person.Orgs[model.OrgAdmin2].PrivLevel != model.PrivLeader {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"org":`)
		out.String(model.OrgNames[org])
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
	if person.Username == "admin" {
		return util.Forbidden
	}
	r.Tx.WillUpdatePerson(person)
	for rid, direct := range person.Roles {
		if !direct {
			delete(person.Roles, rid)
			continue
		}
		role := r.Tx.FetchRole(rid)
		if r.Person.Orgs[role.Org].PrivLevel < model.PrivLeader && r.Person.Orgs[model.OrgAdmin2].PrivLevel < model.PrivLeader {
			continue
		}
		delete(person.Roles, rid)
	}
	r.ParseMultipartForm(1048576)
	for _, ridstr := range r.Form["role"] {
		role := r.Tx.FetchRole(model.Role2ID(util.ParseID(ridstr)))
		if role == nil {
			return errors.New("invalid role")
		}
		if r.Person.Orgs[role.Org].PrivLevel < model.PrivLeader && r.Person.Orgs[model.OrgAdmin2].PrivLevel < model.PrivLeader {
			return errors.New("forbidden role")
		}
		person.Roles[role.ID] = true
	}
	if err := ValidatePerson(r.Tx, person, r.Auth.RolesP(person.ID)); err != nil {
		return err
	}
	r.Tx.UpdatePerson(person)
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}
