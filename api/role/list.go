package role

import (
	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
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
		out.RawString(`,"servGroups":[`)
		first := true
		for _, g := range model.AllSERVGroups {
			if r.SERVGroup&g != 0 {
				if first {
					first = false
				} else {
					out.RawByte(',')
				}
				out.String(servGroupNames[g])
			}
		}
		out.RawString(`]}`)
	}
	out.RawByte(']')
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

var servGroupNames = map[model.SERVGroup]string{
	model.GroupSERVAdmin:      "Admin",
	model.GroupCERTDeployment: "CERT-D",
	model.GroupCERTTraining:   "CERT-T",
	model.GroupListos:         "Listos",
	model.GroupOutreach:       "Outreach",
	model.GroupPEP:            "PEP",
	model.GroupSARES:          "SARES",
	model.GroupSNAP:           "SNAP",
}
