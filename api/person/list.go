package person

import (
	"sort"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPeople handles GET /api/people requests.
func GetPeople(r *util.Request) error {
	var (
		focus *model.Role
		out   jwriter.Writer
	)
	if _, ok := r.Form["search"]; ok {
		return getPeopleSearch(r)
	}
	focus = r.Tx.FetchRole(model.RoleID(util.ParseID(r.FormValue("role"))))
	if focus != nil && r.Person.Orgs[focus.Org].PrivLevel < model.PrivMember {
		focus = nil
	}
	out.RawString(`{"people":[`)
	first := true
	for _, p := range r.Tx.FetchPeople() {
		canView, canViewContactInfo := canViewPerson(r.Person, p)
		if !canView {
			continue
		}
		if focus != nil {
			if _, ok := p.Roles[focus.ID]; !ok {
				continue
			}
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(p.ID))
		out.RawString(`,"informalName":`)
		out.String(p.InformalName)
		out.RawString(`,"sortName":`)
		out.String(p.SortName)
		out.RawString(`,"callSign":`)
		out.String(p.CallSign)
		if canViewContactInfo {
			out.RawString(`,"email":`)
			out.String(p.Email)
			out.RawString(`,"email2":`)
			out.String(p.Email2)
			out.RawString(`,"homeAddress":`)
			p.HomeAddress.MarshalEasyJSON(&out)
			out.RawString(`,"workAddress":`)
			p.WorkAddress.MarshalEasyJSON(&out)
			out.RawString(`,"cellPhone":`)
			out.String(p.CellPhone)
			out.RawString(`,"homePhone":`)
			out.String(p.HomePhone)
			out.RawString(`,"workPhone":`)
			out.String(p.WorkPhone)

		}
		var roles, minPrio = rolesToShow(r, p, focus)
		out.RawString(`,"priority":`)
		out.Int(minPrio)
		out.RawString(`,"roles":[`)
		for i, role := range roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(role)
		}
		out.RawString(`]}`)
	}
	out.RawString(`],"viewableRoles":[`)
	first = true
	for _, role := range r.Tx.FetchRoles() {
		if r.Person.Orgs[role.Org].PrivLevel < model.PrivMember {
			continue
		}
		if !role.ShowRoster {
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
		out.RawByte('}')
	}
	out.RawString(`],"canAdd":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader))
	out.RawString(`,"showCallSign":`)
	out.Bool(focus != nil && focus.Org == model.OrgSARES)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// canViewPerson returns whether the specified viewer is allowed to see the
// specified viewee.  It returns two flags: one indicating whether the viewee
// can be seen in the roster at all; the other indicating whether the viewee's
// contact information can be seen.
func canViewPerson(viewer, viewee *model.Person) (roster, contact bool) {
	if viewer == viewee || viewer.HasPrivLevel(model.PrivLeader) {
		return true, true
	}
	for o, om := range viewer.Orgs {
		if om.PrivLevel < model.PrivMember {
			continue
		}
		if viewee.Orgs[o].PrivLevel == model.PrivNone {
			continue
		}
		roster = true
		if om.PrivLevel < model.PrivMember {
			continue
		}
		if om.PrivLevel == model.PrivMember && viewee.Orgs[o].PrivLevel < model.PrivLeader && !model.Org(o).MembersCanViewContactInfo() {
			continue
		}
		contact = true
		return
	}
	return
}

// rolesToShow returns a list of role titles to show for a particular person,
// when viewing a particular focus group.  It also returns the highest role
// priority (numerically lowest number) of any of those roles.
func rolesToShow(r *util.Request, person *model.Person, focus *model.Role) (list []string, minPrio int) {
	// We want to show all roles that:
	//   - the person holds directly
	//   - have titles
	//   - imply, directly or indirectly, the focus role if any
	//   - are not the focus role itself
	// If there is more than one such role, they are sorted by priority.
	var roles []*model.Role
	minPrio = 99999
	for rid, direct := range person.Roles {
		if !direct || (focus != nil && rid == focus.ID) {
			continue
		}
		role := r.Tx.FetchRole(rid)
		if role.Title == "" {
			continue
		}
		if focus != nil {
			if _, ok := role.Implies[focus.ID]; !ok {
				continue
			}
		}
		roles = append(roles, role)
		if role.Priority < minPrio {
			minPrio = role.Priority
		}
	}
	if len(roles) == 0 {
		return nil, minPrio
	}
	sort.Sort(model.Roles{Roles: roles})
	list = make([]string, len(roles))
	for i, role := range roles {
		list[i] = role.Title
	}
	return list, minPrio
}

// getPeopleSearch handles GET /api/people?search=XXX requests.
func getPeopleSearch(r *util.Request) error {
	var (
		out    jwriter.Writer
		count  int
		search = strings.ToLower(strings.TrimSpace(r.FormValue("search")))
	)
	if !r.Person.HasPrivLevel(model.PrivLeader) {
		return util.Forbidden
	}
	out.RawByte('[')
	for _, p := range r.Tx.FetchPeople() {
		if !strings.Contains(strings.ToLower(p.SortName), search) &&
			!strings.Contains(strings.ToLower(p.FormalName), search) &&
			!strings.Contains(strings.ToLower(p.CallSign), search) {
			continue
		}
		if count != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(p.ID))
		out.RawString(`,"sortName":`)
		out.String(p.SortName)
		out.RawByte('}')
		count++
		if count > 10 {
			break
		}
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
