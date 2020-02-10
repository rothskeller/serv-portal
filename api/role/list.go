package role

import (
	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/util"
)

// GetRoles handles GET /api/roles requests.
func GetRoles(r *util.Request) error {
	var out jwriter.Writer

	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	roles := r.Auth.FetchRoles(r.Auth.AllRoles())
	out.RawByte('[')
	for i, role := range roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		out.RawString(`,"groups":[`)
		for i, g := range r.Auth.FetchGroups(r.Auth.GroupsR(role.ID)) {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(g.Name)
		}
		out.RawString(`]}`)
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}
