package role

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/api/authz"
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
	out.RawString(`,"detail":`)
	out.Bool(role.Detail)
	for _, p := range model.AllPermissions {
		out.RawString(`,"`)
		out.RawString(model.PermissionNames[p])
		out.RawString(`":`)
		out.Bool(role.Permissions&p != 0)
	}
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
	role.Detail, _ = strconv.ParseBool(r.FormValue("detail"))
	role.Permissions = 0
	for _, p := range model.AllPermissions {
		if r.FormValue(model.PermissionNames[p]) == "true" {
			role.Permissions |= p
		}
	}
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

// GetRole2 handles GET /api/roles2/${id} requests.
func GetRole2(r *util.Request, idstr string) error {
	var (
		role *model.Role2
		out  jwriter.Writer
	)
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if idstr == "NEW" {
		role = &model.Role2{
			Implies: make(map[model.Role2ID]bool),
			Lists:   make(map[model.ListID]model.RoleToList),
		}
	} else {
		if role = r.Tx.FetchRole(model.Role2ID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
	}
	out.RawString(`{"id":`)
	out.Int(int(role.ID))
	out.RawString(`,"name":`)
	out.String(role.Name)
	out.RawString(`,"title":`)
	out.String(role.Title)
	out.RawString(`,"org":`)
	out.String(model.OrgNames[role.Org])
	out.RawString(`,"privLevel":`)
	out.String(model.PrivLevelNames[role.PrivLevel])
	out.RawString(`,"showRoster":`)
	out.Bool(role.ShowRoster)
	out.RawString(`,"implicitOnly":`)
	out.Bool(role.ImplicitOnly)
	out.RawString(`,"implies":[`)
	var first = true
	for irid, direct := range role.Implies {
		if !direct {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.Int(int(irid))
	}
	out.RawString(`],"impliable":[`)
	first = true
	for _, ir := range r.Tx.FetchRoles() {
		if ir == role {
			continue
		}
		if _, ok := ir.Implies[role.ID]; ok {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(ir.ID))
		out.RawString(`,"name":`)
		out.String(ir.Name)
		out.RawByte('}')
	}
	out.RawString(`],"lists":[`)
	for i, l := range r.Tx.FetchLists() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(l.ID))
		out.RawString(`,"type":`)
		out.String(model.ListTypeNames[l.Type])
		out.RawString(`,"name":`)
		out.String(l.Name)
		rtl := role.Lists[l.ID]
		out.RawString(`,"subModel":`)
		out.String(model.ListSubModelNames[rtl.SubModel()])
		out.RawString(`,"sender":`)
		out.Bool(rtl.Sender())
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRole2 handles POST /api/roles2/${id} requests.
func PostRole2(r *util.Request, idstr string) error {
	var role *model.Role2

	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if idstr == "NEW" {
		role = &model.Role2{
			Implies: make(map[model.Role2ID]bool),
			Lists:   make(map[model.ListID]model.RoleToList),
		}
	} else {
		if role = r.Tx.FetchRole(model.Role2ID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
		r.Tx.WillUpdateRole(role)
		*role = model.Role2{
			ID:       role.ID,
			Priority: role.Priority,
			Implies:  make(map[model.Role2ID]bool),
			Lists:    make(map[model.ListID]model.RoleToList),
		}
	}
	role.Name = r.FormValue("name")
	role.Title = r.FormValue("title")
	if str := r.FormValue("org"); str != "" {
		var found = false
		for v, s := range model.OrgNames {
			if s == str {
				role.Org = v
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid org")
		}
	}
	if str := r.FormValue("privLevel"); str != "" {
		var found = false
		for v, s := range model.PrivLevelNames {
			if s == str {
				role.PrivLevel = v
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid privLevel")
		}
	}
	role.ShowRoster, _ = strconv.ParseBool(r.FormValue("showRoster"))
	role.ImplicitOnly, _ = strconv.ParseBool(r.FormValue("implicitOnly"))
	for _, iridstr := range r.Form["implies"] {
		if irid := model.Role2ID(util.ParseID(iridstr)); irid > 0 {
			role.Implies[irid] = true
		} else {
			return errors.New("invalid implies")
		}
	}
	for _, liststr := range r.Form["lists"] {
		var triplet = strings.Split(liststr, ":")
		var rtl model.RoleToList
		if len(triplet) != 3 {
			return errors.New("invalid lists")
		}
		var list = r.Tx.FetchList(model.ListID(util.ParseID(triplet[0])))
		if list == nil {
			return errors.New("invalid lists.id")
		}
		if triplet[1] != "" {
			var found = false
			for v, s := range model.ListSubModelNames {
				if triplet[1] == s {
					rtl.SetSubModel(v)
					found = true
					break
				}
			}
			if !found {
				return errors.New("invalid lists.subModel")
			}
		}
		sender, _ := strconv.ParseBool(triplet[2])
		rtl.SetSender(sender)
		role.Lists[list.ID] = rtl
	}
	if err := ValidateRole2(r.Tx, role); err != nil {
		switch err.Error() {
		case "duplicate name":
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		case "duplicate title":
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateTitle":true}`))
			return nil
		default:
			return err
		}
	}
	if idstr == "NEW" {
		r.Tx.CreateRole(role)
	} else {
		r.Tx.UpdateRole(role)
	}
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}

// DeleteRole2 handles DELETE /api/roles2/${id} requests.
func DeleteRole2(r *util.Request, idstr string) error {
	var role *model.Role2

	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if role = r.Tx.FetchRole(model.Role2ID(util.ParseID(idstr))); role == nil {
		return util.NotFound
	}
	r.Tx.DeleteRole(role)
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}
