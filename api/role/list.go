package role

import (
	"errors"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetRoles handles GET /api/roles requests.
func GetRoles(r *util.Request) error {
	var out jwriter.Writer

	if !r.Person.Roles[model.Webmaster] {
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

// GetRoles2 handles GET /api/roles2 requests.
func GetRoles2(r *util.Request) error {
	var out jwriter.Writer

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	out.RawByte('[')
	for i, role := range r.Tx.FetchRoles() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		if role.Org != model.OrgNone2 {
			out.RawString(`,"org":`)
			out.String(model.OrgNames[role.Org])
		}
		if role.PrivLevel != model.PrivNone {
			out.RawString(`,"privLevel":`)
			out.String(model.PrivLevelNames[role.PrivLevel])
		}
		if role.ImplicitOnly {
			out.RawString(`,"implicitOnly":true`)
		}
		out.RawString(`,"people":`)
		out.Int(len(role.People))
		out.RawByte('}')
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRoles2 handles POST /api/roles2 requests (which re-order the existing
// roles).
func PostRoles2(r *util.Request) error {
	var priorities = make(map[*model.Role2]int)

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	r.ParseMultipartForm(1048576)
	if len(r.Form["role"]) != len(r.Tx.FetchRoles()) {
		return errors.New("wrong number of roles")
	}
	for i, idstr := range r.Form["role"] {
		if role := r.Tx.FetchRole(model.Role2ID(util.ParseID(idstr))); role == nil {
			return errors.New("invalid role")
		} else if priorities[role] != 0 {
			return errors.New("role appears twice")
		} else {
			priorities[role] = i + 1
		}
	}
	for role, prio := range priorities {
		r.Tx.WillUpdateRole(role)
		role.Priority = prio
		r.Tx.UpdateRole(role)
	}
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}
