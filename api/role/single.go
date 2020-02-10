package role

import (
	"fmt"
	"strconv"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetRole handles GET /api/roles/$id requests (where $id may be "NEW").
func GetRole(r *util.Request, idstr string) error {
	var (
		role *model.Role
		out  jwriter.Writer
	)
	if idstr == "NEW" {
		role = r.Auth.CreateRole() // but we won't save it
	} else {
		if role = r.Auth.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
	}
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	out.RawString(`{"role":{"id":`)
	out.Int(int(role.ID))
	out.RawString(`,"name":`)
	out.String(role.Name)
	out.RawString(`,"individual":`)
	out.Bool(role.Individual)
	out.RawString(`},"canDelete":`)
	out.Bool(role.Tag == "")
	out.RawString(`,"privs":[`)
	for i, g := range r.Auth.FetchGroups(r.Auth.AllGroups()) {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g.ID))
		out.RawString(`,"name":`)
		out.String(g.Name)
		for _, priv := range model.AllPrivileges {
			out.RawString(`,"` + model.PrivilegeNames[priv] + `":`)
			out.Bool(r.Auth.CanRAG(role.ID, priv, g.ID))
		}
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRole handles POST /api/roles/$id requests (where $id may be "NEW").
func PostRole(r *util.Request, idstr string) error {
	var role *model.Role

	if idstr == "NEW" {
		role = r.Auth.CreateRole()
	} else {
		if role = r.Auth.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
		r.Auth.WillUpdateRole(role)
	}
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && role.Tag == "" && idstr != "NEW" {
		return deleteRole(r, role)
	}
	role.Name = r.FormValue("name")
	role.Individual, _ = strconv.ParseBool(r.FormValue("individual"))
	if err := ValidateRole(r.Auth, role); err != nil {
		if err.Error() == "duplicate name" {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		}
		return err
	}
	for _, gid := range r.Auth.AllGroups() {
		var privs model.Privilege
		for _, p := range model.AllPrivileges {
			key := fmt.Sprintf("%s:%d", model.PrivilegeNames[p], gid)
			if len(r.Form[key]) != 0 {
				privs |= p
			}
		}
		privs = validatePrivileges(privs)
		r.Auth.SetPrivileges(role.ID, privs, gid)
	}
	r.Auth.UpdateRole(role)
	r.Auth.Save()
	r.Tx.Commit()
	return nil
}

func deleteRole(r *util.Request, role *model.Role) error {
	r.Auth.DeleteRole(role.ID)
	r.Auth.Save()
	r.Tx.Commit()
	return nil
}

func validatePrivileges(privs model.Privilege) model.Privilege {
	if privs&(model.PrivManageEvents|model.PrivManageMembers|model.PrivSendTextMessages|model.PrivViewContactInfo) != 0 {
		privs |= model.PrivViewMembers
	}
	if privs&model.PrivSendTextMessages != 0 {
		privs |= model.PrivViewContactInfo
	}
	return privs
}
