package group

import (
	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetGroups handles GET /api/groups requests.
func GetGroups(r *util.Request) error {
	var out jwriter.Writer

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	groups := r.Auth.FetchGroups(r.Auth.AllGroups())
	out.RawByte('[')
	for i, group := range groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(group.ID))
		out.RawString(`,"name":`)
		out.String(group.Name)
		out.RawString(`,"roles":[`)
		for i, r := range r.Auth.FetchRoles(r.Auth.RolesG(group.ID)) {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(r.Name)
		}
		out.RawString(`]}`)
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}
