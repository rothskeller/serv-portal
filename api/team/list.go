package team

import (
	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetTeams handles GET /api/teams requests.
func GetTeams(r *util.Request) error {
	var out jwriter.Writer

	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	teams := r.Tx.FetchTeams()
	r.Tx.Commit()
	teams = util.SortTeamsHierarchically(teams)
	out.RawByte('[')
	for i, t := range teams {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"indent":`)
		out.Int(t.Depth())
		out.RawString(`,"id":`)
		out.Int(int(t.ID))
		out.RawString(`,"name":`)
		out.String(t.Name)
		out.RawString(`,"email":`)
		out.String(t.Email)
		switch t.Type {
		case model.TeamNormal:
			out.RawString(`,"type":"normal"`)
		case model.TeamAncestor:
			out.RawString(`,"type":"ancestor"`)
		case model.TeamTiedRoles:
			out.RawString(`,"type":"tiedRoles"`)
		}
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
	out.RawByte(']')
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}
