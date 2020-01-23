package person

import (
	"errors"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
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
	out.RawString(`{"canEditDetails":`)
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
	out.RawString(`,"person":{"id":`)
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
		out.RawString(`,"phones":[`)
		for i, p := range person.Phones {
			if i != 0 {
				out.RawByte(',')
			}
			p.MarshalEasyJSON(&out)
		}
		out.RawByte(']')
	}
	for _, r := range person.Roles {
		roles[r] = true
	}
	individualHeld = cacheIndividuallyHeldRoles(r, person)
	out.RawString(`,"roles":[`)
	first := true
	for _, role := range r.Tx.FetchRoles() {
		canAssign := auth.CanAssignRole(r, role)
		canEditRoles = canEditRoles || canAssign
		if individualHeld[role] {
			canAssign = false
		}
		if !roles[role.ID] && !canAssign {
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
		out.RawString(`,"canAssign":`)
		out.Bool(canAssign)
		out.RawString(`,"held":`)
		out.Bool(roles[role.ID])
		out.RawByte('}')
	}
	out.RawByte(']')
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
	out.RawString(`},"canEditRoles":`)
	out.Bool(canEditRoles)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// PostPerson handles POST /api/people/$id requests (where $id may be "NEW").
func PostPerson(r *util.Request, idstr string) error {
	var (
		person         *model.Person
		canEditDetails bool
		individualHeld map[*model.Role]bool
		err            error
		previousRoles  = map[model.RoleID]bool{}
		requestedRoles = map[model.RoleID]bool{}
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
	individualHeld = cacheIndividuallyHeldRoles(r, person)
	if canEditDetails {
		if person.InformalName = strings.TrimSpace(r.FormValue("informalName")); person.InformalName == "" {
			return errors.New("missing informalName")
		}
		if person.FormalName = strings.TrimSpace(r.FormValue("formalName")); person.FormalName == "" {
			return errors.New("missing formalName")
		}
		if person.SortName = strings.TrimSpace(r.FormValue("sortName")); person.SortName == "" {
			return errors.New("missing sortName")
		}
		person.Username = strings.TrimSpace(r.FormValue("userName"))
		person.CallSign = strings.TrimSpace(r.FormValue("callSign"))
		for _, p := range r.Tx.FetchPeople() {
			if p.ID == person.ID {
				continue
			}
			if p.SortName == person.SortName {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"duplicateSortName":true}`))
				return nil
			}
			if p.Username != "" && p.Username == person.Username {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"duplicateUsername":true}`))
				return nil
			}
			if p.CallSign != "" && p.CallSign == person.CallSign {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"duplicateCallSign":true}`))
				return nil
			}
		}
		person.Emails = person.Emails[:0]
		for i, e := range r.Form["email"] {
			var email model.PersonEmail
			email.Email = strings.TrimSpace(e)
			if !emailRE.MatchString(email.Email) {
				return errors.New("invalid email")
			}
			if len(r.Form["emailLabel"]) > i {
				email.Label = strings.TrimSpace(r.Form["emailLabel"][i])
			}
			person.Emails = append(person.Emails, &email)
		}
		person.Phones = person.Phones[:0]
		for i, p := range r.Form["phone"] {
			var phone model.PersonPhone
			phone.Phone = strings.Map(keepDigits, p)
			if len(phone.Phone) != 10 {
				return errors.New("invalid phone")
			}
			phone.Phone = phone.Phone[0:3] + "-" + phone.Phone[3:6] + "-" + phone.Phone[6:10]
			if len(r.Form["phoneLabel"]) > i {
				phone.Label = strings.TrimSpace(r.Form["phoneLabel"][i])
			}
			if len(r.Form["phoneSMS"]) > i {
				phone.SMS, _ = strconv.ParseBool(r.Form["phoneSMS"][i])
			}
			person.Phones = append(person.Phones, &phone)
		}
		person.HomeAddress = model.Address{}
		if person.HomeAddress.Address = strings.TrimSpace(r.FormValue("homeAddress")); person.HomeAddress.Address != "" {
			person.HomeAddress.Latitude, err = strconv.ParseFloat(r.FormValue("homeAddressLatitude"), 64)
			if err != nil || person.HomeAddress.Latitude < -90 || person.HomeAddress.Latitude > 90 {
				return errors.New("invalid latitude")
			}
			person.HomeAddress.Longitude, err = strconv.ParseFloat(r.FormValue("homeAddressLongitude"), 64)
			if err != nil || person.HomeAddress.Longitude < -180 || person.HomeAddress.Longitude > 180 {
				return errors.New("invalid longitude")
			}
		}
		person.MailAddress = model.Address{}
		person.MailAddress.Address = strings.TrimSpace(r.FormValue("mailAddress"))
		if sameAsHome, _ := strconv.ParseBool(r.FormValue("mailAddressSameAsHome")); sameAsHome && person.MailAddress.Address == "" {
			person.MailAddress.SameAsHome = true
		}
		person.WorkAddress = model.Address{}
		if person.WorkAddress.Address = strings.TrimSpace(r.FormValue("workAddress")); person.WorkAddress.Address != "" {
			person.WorkAddress.Latitude, err = strconv.ParseFloat(r.FormValue("workAddressLatitude"), 64)
			if err != nil || person.WorkAddress.Latitude < -90 || person.WorkAddress.Latitude > 90 {
				return errors.New("invalid latitude")
			}
			person.WorkAddress.Longitude, err = strconv.ParseFloat(r.FormValue("workAddressLongitude"), 64)
			if err != nil || person.WorkAddress.Longitude < -180 || person.WorkAddress.Longitude > 180 {
				return errors.New("invalid longitude")
			}
		} else if sameAsHome, _ := strconv.ParseBool(r.FormValue("workAddressSameAsHome")); sameAsHome {
			person.WorkAddress.SameAsHome = true
		}
		if password := r.FormValue("password"); password != "" {
			if !auth.IsWebmaster(r) {
				if !auth.StrongPassword(r, person, password) {
					r.Header().Set("Content-Type", "application/json; charset=utf-8")
					r.Write([]byte(`{"weakPassword":true}`))
					return nil
				}
			}
			auth.SetPassword(r, person, password)
		}
	}
	for _, role := range person.Roles {
		previousRoles[role] = true
	}
	for _, ridstr := range r.Form["role"] {
		if role := r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); role != nil {
			requestedRoles[role.ID] = true
		} else {
			return errors.New("bad role")
		}
	}
	person.Roles = person.Roles[:0]
	for _, role := range r.Tx.FetchRoles() {
		if !auth.CanAssignRole(r, role) {
			if previousRoles[role.ID] {
				person.Roles = append(person.Roles, role.ID)
			}
		} else if !individualHeld[role] {
			if requestedRoles[role.ID] {
				person.Roles = append(person.Roles, role.ID)
			}
		}
	}
	if len(person.Roles) == 0 && person.ID == 0 {
		return errors.New("new user with no roles")
	}
	r.Tx.SavePerson(person)
	r.Tx.Commit()
	return nil
}

func cacheIndividuallyHeldRoles(r *util.Request, except *model.Person) (held map[*model.Role]bool) {
	held = make(map[*model.Role]bool)
	for _, p := range r.Tx.FetchPeople() {
		if p.ID == except.ID {
			continue
		}
		for _, role := range r.Tx.FetchRoles() {
			if !role.Individual {
				continue
			}
			if auth.HasRole(p, role) {
				held[role] = true
			}
		}
	}
	return held
}

func usernameInUse(r *util.Request, person *model.Person) bool {
	for _, p := range r.Tx.FetchPeople() {
		if p.ID == person.ID {
			continue
		}
		if strings.EqualFold(p.Username, person.Username) {
			return true
		}
	}
	return false
}

func keepDigits(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
