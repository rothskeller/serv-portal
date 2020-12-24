package report

import (
	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

type clearanceParameters struct {
	role      *model.Role
	with      string
	without   string
	renderCSV bool
	// Not really a parameter, but cached here for convenience:
	allowedRoles []*model.Role
}

func readClearanceParameters(r *util.Request) (params clearanceParameters) {
	if params.allowedRoles = allowedRoles(r); len(params.allowedRoles) == 0 {
		return params // caller will raise error
	}
	if rolestr := r.FormValue("role"); rolestr != "0" {
		params.role = r.Tx.FetchRole(model.RoleID(util.ParseID(rolestr)))
		if params.role == nil || !allowedRole(r, params.role) {
			params.role = params.allowedRoles[0]
		}
	}
	if params.with = r.FormValue("with"); !validClearanceRestriction(r.Person, params.with) {
		params.with = ""
	}
	if params.without = r.FormValue("without"); !validClearanceRestriction(r.Person, params.without) {
		params.without = ""
	}
	if r.FormValue("format") == "csv" {
		params.renderCSV = true
	}
	return params
}

func clearancePersonMatch(p *model.Person, params clearanceParameters) bool {
	if params.role != nil {
		_, ok := p.Roles[params.role.ID]
		return ok
	}
	for _, role := range params.allowedRoles {
		if _, ok := p.Roles[role.ID]; ok {
			return true
		}
	}
	return false
}

func clearanceRenderParams(r *util.Request, out *jwriter.Writer, params clearanceParameters) {
	out.RawString(`"parameters":{"role":`)
	if params.role != nil {
		out.Int(int(params.role.ID))
	} else {
		out.RawByte('0')
	}
	out.RawString(`,"with":`)
	out.String(params.with)
	out.RawString(`,"without":`)
	out.String(params.without)
	out.RawString(`,"allowedRoles":[`)
	for i, role := range params.allowedRoles {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		out.RawByte('}')
	}
	out.RawString(`],"allowedRestrictions":`)
	emitValidClearanceRestrictions(r.Person, out)
	out.RawString(`}`)
}

func allowedRole(r *util.Request, role *model.Role) bool {
	return role.ShowRoster && r.Person.Orgs[role.Org].PrivLevel >= model.PrivLeader
}
func allowedRoles(r *util.Request) (roles []*model.Role) {
	for _, role := range r.Tx.FetchRoles() {
		if allowedRole(r, role) {
			roles = append(roles, role)
		}
	}
	return roles
}
