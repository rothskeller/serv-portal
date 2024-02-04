package listedit

import (
	"regexp"
	"slices"
	"sort"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/listlist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type roleData struct {
	rl       *role.Role
	sender   bool
	submodel listrole.SubscriptionModel
}

// Handle handles /admin/lists/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user  *person.Person
		l     *list.List
		ul    *list.Updater
		roles []*roleData
		f     form.Form
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	f.Buttons = []*form.Button{{
		Label:   "Save",
		OnClick: func() bool { return saveList(r, user, l, ul, roles) },
	}}
	if idstr == "NEW" {
		ul = new(list.Updater)
		f.Title = "New List"
	} else {
		if l = list.WithID(r, list.ID(util.ParseID(idstr))); l == nil {
			errpage.NotFound(r, user)
			return
		}
		f.Title = "Edit List"
		ul = l.Updater()
		listrole.AllRolesForList(r, ul.ID, role.FID|role.FName, func(rl *role.Role, sender bool, submodel listrole.SubscriptionModel) {
			roles = append(roles, &roleData{rl.Clone(), sender, submodel})
		})
		f.Buttons = append(f.Buttons, &form.Button{
			Name: "delete", Label: "Delete", Style: "danger",
			OnClick: func() bool { return deleteList(r, user, l) },
		})
	}
	f.Rows = []form.Row{
		&typeRow{form.RadioGroupRow[list.Type]{
			LabeledRow: form.LabeledRow{
				RowID: "listeditType",
				Label: "Type",
			},
			Name:    "type",
			ValueP:  &ul.Type,
			Options: list.AllTypes,
		}},
		&nameRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "listeditName",
				Label: "Name",
			},
			Name:   "name",
			ValueP: &ul.Name,
		}, ul},
		&rolesRow{form.LabeledRow{Label: "Roles"}, ul, &roles},
	}
	f.Handle(r)
}

type typeRow struct{ form.RadioGroupRow[list.Type] }

func (tr *typeRow) Read(r *request.Request) bool {
	if !tr.RadioGroupRow.Read(r) {
		return false
	}
	if *tr.ValueP == 0 {
		tr.Error = "The list type is required."
		return false
	}
	return true
}

type nameRow struct {
	form.TextInputRow
	ul *list.Updater
}

func (nr *nameRow) ShouldEmit(vl request.ValidationList) bool {
	return vl.ValidatingAny("type", "name")
}

var emailNameRE = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

func (nr *nameRow) Read(r *request.Request) bool {
	if !nr.TextInputRow.Read(r) {
		return false
	}
	if nr.ul.Name == "" {
		nr.Error = "The list name is required."
		return false
	}
	if nr.ul.DuplicateName(r) {
		nr.Error = "Another list has this name."
		return false
	}
	if nr.ul.Type == list.Email && !emailNameRE.MatchString(nr.ul.Name) {
		nr.Error = "The list name is not valid as the first part of an @sunnyvaleserv.org email address."
		return false
	}
	return true
}

type rolesRow struct {
	form.LabeledRow
	ul    *list.Updater
	roles *[]*roleData
}

func (rr *rolesRow) Read(r *request.Request) bool {
	var roles []*roleData

	for key := range r.Form {
		var rd roleData
		if !strings.HasPrefix(key, "role") {
			continue
		}
		if rd.rl = role.WithID(r, role.ID(util.ParseID(key[4:])), role.FID|role.FName); rd.rl == nil {
			continue
		}
		parts := strings.Split(r.FormValue(key), ":")
		if len(parts) != 2 {
			continue
		}
		if rd.submodel = listrole.SubscriptionModel(util.ParseID(parts[0])); !rd.submodel.Valid() {
			continue
		}
		rd.sender = parts[1] == "true"
		roles = append(roles, &rd)
	}
	if len(roles) == 0 {
		rr.Error = "At least one role must have privileges."
		return false
	}
	for _, or := range *rr.roles {
		if !slices.ContainsFunc(roles, func(rd *roleData) bool { return rd.rl.ID() == or.rl.ID() }) {
			or.submodel, or.sender = 0, false
			roles = append(roles, or)
		}
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].rl.Name() < roles[j].rl.Name() })
	*rr.roles = roles
	rr.Error = ""
	return true
}

func (rr *rolesRow) ShouldEmit(vl request.ValidationList) bool {
	return vl.Validating("csrf") // this happens when listrole dialog closes
}

func (rr *rolesRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	row := rr.EmitPrefix(r, parent, "")
	box := row.E("div class=formInput")
	for _, rd := range *rr.roles {
		if rd.submodel == 0 && !rd.sender {
			continue
		}
		div := box.E("div")
		div.E("a href=# class=listeditRoleEdit data-list=%d data-role=%d>%s", rr.ul.ID, rd.rl.ID(), rd.rl.Name())
		div.R(" (")
		if rd.submodel != 0 {
			div.R(rd.submodel.String())
		}
		if rd.submodel != 0 && rd.sender {
			div.R(", ")
		}
		if rd.sender {
			div.R("sender")
		}
		div.R(")")
		div.E("input type=hidden id=listeditRole%d name=role%d value=%d:%v", rd.rl.ID(), rd.rl.ID(), rd.submodel, rd.sender)
	}
	box.E("a href=# class='sbtn sbtn-small sbtn-primary listeditRoleEdit' data-list=%d>Add", rr.ul.ID)
	rr.EmitSuffix(r, row)
}

func saveList(r *request.Request, user *person.Person, l *list.List, ul *list.Updater, roles []*roleData) bool {
	r.Transaction(func() {
		if l == nil {
			l = list.Create(r, ul)
		} else {
			l.Update(r, ul)
		}
		for _, rd := range roles {
			listrole.SetListRole(r, l, rd.rl, rd.sender, rd.submodel)
		}
		role.Recalculate(r)
	})
	listlist.Render(r, user)
	return true
}

func deleteList(r *request.Request, user *person.Person, l *list.List) bool {
	r.Transaction(func() {
		l.Delete(r)
		role.Recalculate(r)
	})
	listlist.Render(r, user)
	return true
}
