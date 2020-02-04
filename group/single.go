package group

import (
	"fmt"
	"strconv"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetGroup handles GET /api/groups/$id requests (where $id may be "NEW").
func GetGroup(r *util.Request, idstr string) error {
	var (
		group *model.Group
		out   jwriter.Writer
	)
	if idstr == "NEW" {
		group = new(model.Group)
		r.Auth.CreateGroup(group) // but we won't save it
	} else {
		if group = r.Auth.FetchGroup(model.GroupID(util.ParseID(idstr))); group == nil {
			return util.NotFound
		}
	}
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	out.RawString(`{"group":{"id":`)
	out.Int(int(group.ID))
	out.RawString(`,"name":`)
	out.String(group.Name)
	out.RawString(`,"email":`)
	out.String(group.Email)
	out.RawString(`,"allowTextMessages":`)
	out.Bool(group.AllowTextMessages)
	out.RawString(`},"canDelete":`)
	out.Bool(group.Tag == "")
	out.RawString(`,"privs":[`)
	for i, role := range r.Auth.FetchRoles(r.Auth.AllRoles()) {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(role.ID))
		out.RawString(`,"name":`)
		out.String(role.Name)
		for _, priv := range model.AllPrivileges {
			out.RawString(`,"` + model.PrivilegeNames[priv] + `":`)
			if role.Tag == model.RoleWebmaster && idstr == "NEW" && priv != model.PrivMember {
				out.Bool(true)
			} else {
				out.Bool(r.Auth.CanRAG(role.ID, priv, group.ID))
			}
		}
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostGroup handles POST /api/groups/$id requests (where $id may be "NEW").
func PostGroup(r *util.Request, idstr string) error {
	var group *model.Group

	if idstr == "NEW" {
		group = new(model.Group)
		r.Auth.CreateGroup(group)
	} else {
		if group = r.Auth.FetchGroup(model.GroupID(util.ParseID(idstr))); group == nil {
			return util.NotFound
		}
	}
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && group.Tag == "" && idstr != "NEW" {
		return deleteGroup(r, group)
	}
	group.Name = r.FormValue("name")
	group.Email = r.FormValue("email")
	group.AllowTextMessages, _ = strconv.ParseBool(r.FormValue("allowTextMessages"))
	if err := ValidateGroup(r.Auth, group); err != nil {
		if err.Error() == "duplicate name" {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		}
		if err.Error() == "duplicate email" {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"duplicateEmail":true}`))
			return nil
		}
		return err
	}
	for _, rid := range r.Auth.AllRoles() {
		var privs model.Privilege
		for _, p := range model.AllPrivileges {
			key := fmt.Sprintf("%s:%d", model.PrivilegeNames[p], rid)
			if len(r.Form[key]) != 0 {
				privs |= p
			}
		}
		privs = validatePrivileges(privs)
		r.Auth.SetPrivileges(rid, privs, group.ID)
	}
	r.Auth.Save()
	r.Tx.Commit()
	return nil
}

func deleteGroup(r *util.Request, group *model.Group) error {
	r.Auth.DeleteGroup(group.ID)
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
