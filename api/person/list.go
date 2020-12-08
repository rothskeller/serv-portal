package person

import (
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPeople handles GET /api/people requests.
func GetPeople(r *util.Request) error {
	var (
		focus *model.Role
		out   jwriter.Writer
		now   = time.Now()
	)
	focus = r.Tx.FetchRole(model.RoleID(util.ParseID(r.FormValue("role"))))
	if _, ok := r.Form["search"]; ok {
		return getPeopleSearch(r)
	}
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
			out.RawString(`,"mailAddress":`)
			p.MailAddress.MarshalEasyJSON(&out)
			out.RawString(`,"workAddress":`)
			p.WorkAddress.MarshalEasyJSON(&out)
			out.RawString(`,"cellPhone":`)
			out.String(p.CellPhone)
			out.RawString(`,"homePhone":`)
			out.String(p.HomePhone)
			out.RawString(`,"workPhone":`)
			out.String(p.WorkPhone)
		}
		if r.Person.HasPrivLevel(model.PrivLeader) {
			var badges []string
			if needVolgisticsID(r, p, focus) {
				if p.VolgisticsID == 0 {
					badges = append(badges, "Not Volunteer")
				}
			} else {
				if p.VolgisticsID != 0 {
					badges = append(badges, "Volunteer")
				}
			}
			for _, c := range model.AllDSWClasses {
				if needDSW(r, p, c, focus) {
					if p.DSWUntil == nil || p.DSWUntil[c].Before(now) {
						badges = append(badges, "No DSW "+model.DSWClassNames[c][:4])
					}
				} else {
					if p.DSWUntil != nil && !p.DSWUntil[c].Before(now) {
						badges = append(badges, "DSW "+model.DSWClassNames[c][:4])
					}
				}
			}
			if (focus == nil && p.HasPrivLevel(model.PrivMember)) || (focus != nil && p.Orgs[focus.Org].PrivLevel >= model.PrivMember) {
				if p.BackgroundCheck == "" && r.Person.IsAdminLeader() {
					// Setting this to admins only until we have accurate BG check data.
					badges = append(badges, "No BG Check")
				}
			} else {
				if p.BackgroundCheck != "" {
					badges = append(badges, "BG Check")
				}
			}
			if len(badges) != 0 {
				out.RawString(`,"badges":[`)
				for i, b := range badges {
					if i != 0 {
						out.RawByte(',')
					}
					out.String(b)
				}
				out.RawByte(']')
			}
		}
		out.RawString(`,"roles":[`)
		first2 := true
		for rid, direct := range p.Roles {
			if !direct {
				continue
			}
			if first2 {
				first2 = false
			} else {
				out.RawByte(',')
			}
			out.String(r.Tx.FetchRole(rid).Name)
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
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
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
		if om.PrivLevel == model.PrivMember && viewee.Orgs[o].PrivLevel < model.PrivLeader && !model.Org(o).MembersCanViewContactInfo() {
			continue
		}
		contact = true
		return
	}
	return
}
