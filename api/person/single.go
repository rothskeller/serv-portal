package person

import (
	"errors"
	"regexp"
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
		canEditInfo    bool
		out            jwriter.Writer
		individualHeld map[*model.Role]bool
	)
	if idstr == "NEW" {
		if !auth.CanCreatePeople(r) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditInfo = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if !auth.CanViewPerson(r, person) {
			return util.Forbidden
		}
		canEditInfo = r.Person == person || auth.IsWebmaster(r)
	}
	individualHeld = cacheIndividuallyHeldRoles(r, person)
	r.Tx.Commit()
	out.RawString(`{"canEditInfo":`)
	out.Bool(canEditInfo)
	out.RawString(`,"allowBadPassword":`)
	out.Bool(auth.IsWebmaster(r))
	out.RawString(`,"person":{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"firstName":`)
	out.String(person.FirstName)
	out.RawString(`,"lastName":`)
	out.String(person.LastName)
	out.RawString(`,"email":`)
	out.String(person.Email)
	out.RawString(`,"phone":`)
	out.String(person.Phone)
	out.RawString(`},"roles":[`)
	first := true
	for _, role := range r.Tx.FetchRoles() {
		var enabled = true
		if !auth.CanAssignRole(r, role) || individualHeld[role] || role.ImplyOnly {
			enabled = false
		}
		if !enabled {
			var found bool
			for _, r := range person.Roles {
				if role == r {
					found = true
				}
			}
			if !found {
				continue
			}
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
		out.RawString(`,"memberLabel":`)
		out.String(role.MemberLabel)
		out.RawString(`,"held":`)
		out.Bool(auth.HasRole(person, role))
		out.RawString(`,"enabled":`)
		out.Bool(enabled)
		out.RawByte('}')
	}
	out.RawString(`],"passwordHints":[`)
	for i, h := range auth.SERVPasswordHints {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(h)
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// PostPerson handles POST /api/people/$id requests (where $id may be "NEW").
func PostPerson(r *util.Request, idstr string) error {
	var (
		person         *model.Person
		canEditInfo    bool
		individualHeld map[*model.Role]bool
		previousRoles  = map[*model.Role]bool{}
		requestedRoles = map[*model.Role]bool{}
	)
	if idstr == "NEW" {
		if !auth.CanCreatePeople(r) {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditInfo = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		canEditInfo = r.Person == person || auth.IsWebmaster(r)
	}
	if !canEditInfo && !auth.CanAssignAnyRole(r) {
		return util.Forbidden
	}
	individualHeld = cacheIndividuallyHeldRoles(r, person)
	if canEditInfo {
		if person.FirstName = strings.TrimSpace(r.FormValue("firstName")); person.FirstName == "" {
			return errors.New("missing firstName")
		}
		if person.LastName = strings.TrimSpace(r.FormValue("lastName")); person.LastName == "" {
			return errors.New("missing lastName")
		}
		for _, p := range r.Tx.FetchPeople() {
			if p.ID == person.ID {
				continue
			}
			if p.FirstName == person.FirstName && p.LastName == person.LastName {
				r.Header().Set("Content-Type", "application/json; charset=utf-8")
				r.Write([]byte(`{"duplicateName":true}`))
				return nil
			}
		}
		if person.Email = strings.TrimSpace(r.FormValue("email")); person.Email == "" {
			return errors.New("missing email")
		} else if !emailRE.MatchString(person.Email) {
			return errors.New("invalid email")
		} else if emailInUse(r, person) {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateEmail":true}`))
			return nil
		}
		if person.Phone = strings.TrimSpace(r.FormValue("phone")); person.Phone != "" {
			ph := strings.Map(keepDigits, person.Phone)
			if len(ph) != 10 {
				return errors.New("invalid phone")
			}
			person.Phone = ph[0:3] + "-" + ph[3:6] + "-" + ph[6:10]
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
			requestedRoles[role] = true
		} else {
			return errors.New("bad role")
		}
	}
	person.Roles = person.Roles[:0]
	for _, role := range r.Tx.FetchRoles() {
		if !auth.CanAssignRole(r, role) {
			if previousRoles[role] {
				person.Roles = append(person.Roles, role)
			}
		} else if !individualHeld[role] && !role.ImplyOnly {
			if requestedRoles[role] {
				person.Roles = append(person.Roles, role)
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
			if p.PrivMap.Has(role.ID, model.PrivHoldsRole) {
				held[role] = true
			}
		}
	}
	return held
}

func emailInUse(r *util.Request, person *model.Person) bool {
	for _, p := range r.Tx.FetchPeople() {
		if p.ID == person.ID {
			continue
		}
		if strings.EqualFold(p.Email, person.Email) {
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
