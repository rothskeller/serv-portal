package person

import (
	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetPeople handles GET /api/people requests.
func GetPeople(r *util.Request) error {
	var (
		focus *model.Group
		out   jwriter.Writer
	)
	focus = r.Tx.FetchGroup(model.GroupID(util.ParseID(r.FormValue("group"))))
	out.RawString(`{"people":[`)
	first := true
	for _, p := range r.Tx.FetchPeople() {
		if !auth.CanViewPerson(r, p) {
			continue
		}
		if focus != nil && focus.Tag == model.GroupDisabled && auth.IsEnabled(r, p) {
			// Special case because the lack of *any* role also
			// means disabled.
			continue
		} else if focus != nil && !auth.IsMember(p, focus) {
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
		if auth.CanViewContactInfo(r, p) {
			out.RawString(`,"emails":[`)
			for i, e := range p.Emails {
				if i != 0 {
					out.RawByte(',')
				}
				e.MarshalEasyJSON(&out)
			}
			out.RawString(`],"homeAddress":`)
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
		for i, role := range p.Roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(r.Tx.FetchRole(role).Name)
		}
		out.RawString(`]}`)
	}
	out.RawString(`],"viewableGroups":[`)
	first = true
	for _, group := range r.Tx.FetchGroups() {
		if !auth.CanViewGroup(r, group) {
			continue
		}
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
	out.Bool(auth.CanCreatePeople(r))
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
