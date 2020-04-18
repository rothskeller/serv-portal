package person

import (
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPeople handles GET /api/people requests.
func GetPeople(r *util.Request) error {
	var (
		focus  *model.Group
		people []*model.Person
		out    jwriter.Writer
	)
	focus = r.Auth.FetchGroup(model.GroupID(util.ParseID(r.FormValue("group"))))
	if _, ok := r.Form["search"]; ok {
		return getPeopleSearch(r)
	}
	if focus != nil && !r.Auth.CanAG(model.PrivViewMembers, focus.ID) {
		focus = nil
	}
	if focus != nil {
		people = r.Auth.FetchPeople(r.Auth.PeopleG(focus.ID))
	} else {
		people = r.Tx.FetchPeople()
	}
	out.RawString(`{"people":[`)
	first := true
	for _, p := range people {
		if !r.Auth.CanAP(model.PrivViewMembers, p.ID) && p != r.Person {
			continue
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
		if r.Auth.CanAP(model.PrivViewContactInfo, p.ID) || p == r.Person {
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
		out.RawString(`,"roles":[`)
		first2 := true
		for _, role := range r.Auth.FetchRoles(r.Auth.RolesP(p.ID)) {
			if role.Detail {
				continue
			}
			if first2 {
				first2 = false
			} else {
				out.RawByte(',')
			}
			out.String(role.Name)
		}
		out.RawString(`]}`)
	}
	out.RawString(`],"viewableGroups":[`)
	first = true
	for _, group := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivViewMembers)) {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(group.ID))
		out.RawString(`,"name":`)
		out.String(group.Name)
		out.RawByte('}')
	}
	out.RawString(`],"canAdd":`)
	out.Bool(r.Auth.CanA(model.PrivManageMembers))
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
	out.RawByte('[')
	for _, p := range r.Auth.FetchPeople(r.Auth.PeopleA(model.PrivViewMembers)) {
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
