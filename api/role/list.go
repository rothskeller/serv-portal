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
	out.RawByte('[')
	for i, role := range r.Tx.FetchRoles() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		if role.Org != model.OrgNone {
			out.RawString(`,"org":`)
			out.String(role.Org.String())
		}
		if role.PrivLevel != model.PrivNone {
			out.RawString(`,"privLevel":`)
			out.String(role.PrivLevel.String())
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

// PostRoles handles POST /api/roles requests (which re-order the existing
// roles).
func PostRoles(r *util.Request) error {
	var priorities = make(map[*model.Role]int)

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	r.ParseMultipartForm(1048576)
	if len(r.Form["role"]) != len(r.Tx.FetchRoles()) {
		return errors.New("wrong number of roles")
	}
	for i, idstr := range r.Form["role"] {
		if role := r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
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
