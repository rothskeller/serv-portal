package group

import (
	"errors"
	"fmt"
	"sort"

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
		group = r.Auth.CreateGroup() // but we won't save it
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
	if len(group.NoEmail) != 0 {
		people := make([]*model.Person, len(group.NoEmail))
		for i := range group.NoEmail {
			people[i] = r.Tx.FetchPerson(group.NoEmail[i])
		}
		sort.Sort(model.PersonSort(people))
		out.RawString(`,"noEmail":[`)
		for i, p := range people {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(p.SortName)
		}
		out.RawByte(']')
	}
	if len(group.NoText) != 0 {
		people := make([]*model.Person, len(group.NoText))
		for i := range group.NoText {
			people[i] = r.Tx.FetchPerson(group.NoText[i])
		}
		sort.Sort(model.PersonSort(people))
		out.RawString(`,"noText":[`)
		for i, p := range people {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(p.SortName)
		}
		out.RawByte(']')
	}
	switch group.DSWType {
	case model.DSWNone:
		out.RawString(`,"dswType":""`)
	case model.DSWExtended:
		out.RawString(`,"dswType":"extended"`)
	case model.DSWRequired:
		out.RawString(`,"dswType":"required"`)
	}
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
		group = r.Auth.CreateGroup()
		println(group.ID)
	} else {
		if group = r.Auth.FetchGroup(model.GroupID(util.ParseID(idstr))); group == nil {
			return util.NotFound
		}
		r.Auth.WillUpdateGroup(group)
	}
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && group.Tag == "" && idstr != "NEW" {
		return deleteGroup(r, group)
	}
	group.Name = r.FormValue("name")
	group.Email = r.FormValue("email")
	switch r.FormValue("dswType") {
	case "":
		group.DSWType = model.DSWNone
	case "extended":
		group.DSWType = model.DSWExtended
	case "required":
		group.DSWType = model.DSWRequired
	default:
		return errors.New("invalid dswType")
	}
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
	r.Auth.UpdateGroup(group)
	r.Auth.Save()
	r.Tx.Commit()
	return nil
}

func deleteGroup(r *util.Request, group *model.Group) error {
	for _, event := range r.Tx.FetchEvents("2000-01-01", "2099-12-31") {
		found := false
		j := 0
		for _, g := range event.Groups {
			if g != group.ID {
				event.Groups[j] = g
				j++
			} else {
				found = true
			}
		}
		if found {
			event.Groups = event.Groups[:j]
			r.Tx.UpdateEvent(event)
		}
	}
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
