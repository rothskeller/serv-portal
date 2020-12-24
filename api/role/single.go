package role

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetRole handles GET /api/roles/${id} requests.
func GetRole(r *util.Request, idstr string) error {
	var (
		role *model.Role
		out  jwriter.Writer
	)
	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	if idstr == "NEW" {
		role = &model.Role{
			Implies: make(map[model.RoleID]bool),
			Lists:   make(map[model.ListID]model.RoleToList),
		}
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
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
	out.String(role.Org.String())
	out.RawString(`,"privLevel":`)
	out.String(role.PrivLevel.String())
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
		if ir == role || ir.ID == model.Webmaster {
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

// PostRole handles POST /api/roles/${id} requests.
func PostRole(r *util.Request, idstr string) error {
	var role *model.Role
	var err error

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	if idstr == "NEW" {
		role = &model.Role{
			Implies: make(map[model.RoleID]bool),
			Lists:   make(map[model.ListID]model.RoleToList),
		}
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
		r.Tx.WillUpdateRole(role)
		*role = model.Role{
			ID:       role.ID,
			Priority: role.Priority,
			Implies:  make(map[model.RoleID]bool),
			Lists:    make(map[model.ListID]model.RoleToList),
		}
	}
	role.Name = r.FormValue("name")
	role.Title = r.FormValue("title")
	if role.Org, err = model.ParseOrg(r.FormValue("org")); err != nil {
		return err
	}
	if str := r.FormValue("privLevel"); str != "" {
		if role.PrivLevel, err = model.ParsePrivLevel(str); err != nil {
			return err
		}
	}
	role.ShowRoster, _ = strconv.ParseBool(r.FormValue("showRoster"))
	role.ImplicitOnly, _ = strconv.ParseBool(r.FormValue("implicitOnly"))
	for _, iridstr := range r.Form["implies"] {
		if irid := model.RoleID(util.ParseID(iridstr)); irid > 0 {
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
	if err := ValidateRole(r.Tx, role); err != nil {
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

// DeleteRole handles DELETE /api/roles/${id} requests.
func DeleteRole(r *util.Request, idstr string) error {
	var role *model.Role

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
		return util.NotFound
	}
	for _, e := range r.Tx.FetchEvents("2001-01-01", "2099-12-31") {
		j := 0
		for _, r := range e.Roles {
			if r != role.ID {
				e.Roles[j] = r
				j++
			}
		}
		if j < len(e.Roles) {
			e.Roles = e.Roles[:j]
			r.Tx.UpdateEvent(e)
		}
	}
	r.Tx.DeleteRole(role)
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}
