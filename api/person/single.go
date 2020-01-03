package person

import (
	"errors"
	"regexp"
	"sort"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetPerson handles GET /api/people/$id requests (where $id may be "NEW").
func GetPerson(r *util.Request, idstr string) error {
	var (
		person      *model.Person
		canEditInfo bool
		teams       []*model.Team
		manageTeams map[*model.Team]bool
		teamMap     map[*model.Team]*model.Role
		roleMap     map[*model.Role]bool
		out         jwriter.Writer
		first       = true
	)
	teams, _, manageTeams = editPersonTeams(r)
	if idstr == "NEW" {
		if len(manageTeams) == 0 {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditInfo = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		if !r.Person.CanViewPerson(person) {
			return util.Forbidden
		}
		canEditInfo = r.Person == person || r.Person.IsWebmaster()
	}
	r.Tx.Commit()
	teamMap = make(map[*model.Team]*model.Role)
	roleMap = make(map[*model.Role]bool)
	for _, role := range person.Roles {
		roleMap[role] = true
		teamMap[role.Team] = role
	}
	out.RawString(`{"canEditInfo":`)
	out.Bool(canEditInfo)
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
	out.RawString(`},"teams":[`)
	for _, t := range teams {
		privs := r.Person.PrivMap.Get(t)
		if privs&model.PrivManage == 0 && teamMap[t] == nil {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(t.ID))
		out.RawString(`,"name":`)
		out.String(t.Name)
		out.RawString(`,"role":`)
		if teamMap[t] == nil {
			out.RawByte('0')
		} else {
			out.Int(int(teamMap[t].ID))
		}
		out.RawString(`,"canManage":`)
		out.Bool(privs&model.PrivManage != 0)
		out.RawString(`,"canAdmin":`)
		out.Bool(privs&model.PrivAdmin != 0)
		out.RawString(`,"roles":[`)
		for j, r := range t.Roles {
			if j != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(r.ID))
			out.RawString(`,"name":`)
			out.String(r.Name)
			out.RawByte('}')
		}
		out.RawString(`]}`)
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
		person      *model.Person
		canEditInfo bool
		adminTeams  map[*model.Team]bool
		manageTeams map[*model.Team]bool
	)
	_, adminTeams, manageTeams = editPersonTeams(r)
	if idstr == "NEW" {
		if len(manageTeams) == 0 {
			return util.Forbidden
		}
		person = new(model.Person)
		canEditInfo = true
	} else {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
		canEditInfo = r.Person == person || r.Person.IsWebmaster()
	}
	if !canEditInfo && len(adminTeams) == 0 && len(manageTeams) == 0 {
		return util.Forbidden
	}
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
				r.Write([]byte(`{"nameError":"An account already exists for a person with this name."}`))
				return nil
			}
		}
		if person.Email = strings.TrimSpace(r.FormValue("email")); person.Email == "" {
			return errors.New("missing email")
		} else if !emailRE.MatchString(person.Email) {
			return errors.New("invalid email")
		} else if emailInUse(r, person) {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"emailError":"This email address is in use by a different person."}`))
			return nil
		}
		if person.Phone = strings.TrimSpace(r.FormValue("phone")); person.Phone != "" {
			ph := strings.Map(keepDigits, person.Phone)
			if len(ph) != 10 {
				return errors.New("invalid phone")
			}
			person.Phone = ph[0:3] + "-" + ph[3:6] + "-" + ph[6:10]
		}
	}
	var rmap = make(map[*model.Team]*model.Role, len(person.Roles))
	for _, r := range person.Roles {
		if !manageTeams[r.Team] {
			rmap[r.Team] = r
		}
	}
	for _, ridstr := range r.Form["role"] {
		if role := r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); role != nil {
			if !adminTeams[role.Team] {
				continue
			}
			if manageTeams[role.Team] || rmap[role.Team] != nil {
				rmap[role.Team] = role
			}
		} else {
			return errors.New("invalid role")
		}
	}
	person.Roles = person.Roles[:0]
	for _, t := range r.Tx.FetchTeams() {
		if r := rmap[t]; r != nil {
			person.Roles = append(person.Roles, r)
		}
	}
	r.Tx.SavePerson(person)
	r.Tx.Commit()
	return nil
}

func editPersonTeams(r *util.Request) (teams []*model.Team, admin, manage map[*model.Team]bool) {
	admin = make(map[*model.Team]bool)
	manage = make(map[*model.Team]bool)
	for _, t := range r.Tx.FetchTeams() {
		if t.Type != model.TeamNormal {
			continue
		}
		teams = append(teams, t)
		if r.Person.PrivMap.Has(t, model.PrivManage) {
			manage[t] = true
		}
		if r.Person.PrivMap.Has(t, model.PrivAdmin) {
			admin[t] = true
		}
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams, admin, manage
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
