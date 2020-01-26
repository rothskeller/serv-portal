package role

import (
	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/util"
)

// GetRoles handles GET /api/roles requests.
func GetRoles(r *util.Request) error {
	var out jwriter.Writer

	if !auth.IsWebmaster(r) {
		return util.Forbidden
	}
	roles := r.Tx.FetchRoles()
	r.Tx.Commit()
	out.RawByte('[')
	for i, r := range roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(r.ID))
		out.RawString(`,"name":`)
		out.String(r.Name)
		out.RawByte('}')
	}
	out.RawByte(']')
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}
