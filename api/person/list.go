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
		focus *model.Role
		out   jwriter.Writer
	)
	focus = r.Tx.FetchRole(model.RoleID(util.ParseID(r.FormValue("role"))))
	out.RawString(`{"people":[`)
	first := true
	for _, p := range r.Tx.FetchPeople() {
		if !auth.CanViewPerson(r, p) {
			continue
		}
		if focus != nil && focus.Tag == model.RoleDisabled && auth.IsEnabled(r, p) {
			// Special case because the lack of *any* role also
			// means disabled.
			continue
		} else if focus != nil && !auth.HasRole(p, focus) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(p.ID))
		out.RawString(`,"firstName":`)
		out.String(p.FirstName)
		out.RawString(`,"lastName":`)
		out.String(p.LastName)
		out.RawString(`,"nickname":`)
		out.String(p.Nickname)
		out.RawString(`,"suffix":`)
		out.String(p.Suffix)
		out.RawString(`,"email":`)
		out.String(p.Email)
		out.RawString(`,"phone":`)
		out.String(p.Phone)
		out.RawString(`,"roles":[`)
		first := true
		for _, r := range p.Roles {
			if r.MemberLabel == "" {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(r.MemberLabel)
		}
		out.RawString(`]}`)
	}
	out.RawString(`],"viewableRoles":[`)
	first = true
	for _, role := range r.Tx.FetchRoles() {
		if !auth.CanViewRole(r, role) {
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
	out.Bool(auth.CanCreatePeople(r))
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
