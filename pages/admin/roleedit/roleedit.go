package roleedit

import (
	"strconv"

	"sunnyvaleserv.org/portal/pages/admin/rolelist"
	"sunnyvaleserv.org/portal/pages/admin/roleselect"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /admin/roles/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user        *person.Person
		rl          *role.Role
		ur          *role.Updater
		f           form.Form
		prioBefore  uint
		prioOptions []uint
		prioLabels  map[uint]string
		impliers    map[role.ID]struct{}
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	f.Buttons = []*form.Button{{
		Label:   "Save",
		OnClick: func() bool { return saveRole(r, user, rl, ur, prioBefore) },
	}}
	if idstr == "NEW" {
		ur = new(role.Updater)
		f.Title = "New Role"
	} else {
		if rl = role.WithID(r, role.ID(util.ParseID(idstr)), role.UpdaterFields); rl == nil {
			errpage.NotFound(r, user)
			return
		}
		ur = rl.Updater()
		f.Title = "Edit Role"
		if personrole.PeopleCountForRole(r, rl.ID()) == 0 {
			f.Buttons = append(f.Buttons, &form.Button{
				Name: "delete", Label: "Delete", Style: "danger",
				OnClick: func() bool { return deleteRole(r, user, rl) },
			})
		}
	}
	prioBefore, prioOptions, prioLabels = setupPriorityList(r, ur)
	if ur.ID != 0 {
		impliers = role.AllThatImply(r, ur.ID)
	}
	f.Rows = []form.Row{
		&nameRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "roleeditName",
				Label: "Name",
				Help:  "Collective name for people who hold this role",
			},
			Name:   "name",
			ValueP: &ur.Name,
		}, ur},
		&form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "roleeditTitle",
				Label: "Title",
				Help:  "Title for a single person who holds this role, or empty if this role should not be called out in people lists",
			},
			Name:     "title",
			ValueP:   &ur.Title,
			Validate: form.NoValidate,
		},
		&orgRow{form.RadioGroupRow[enum.Org]{
			LabeledRow: form.LabeledRow{
				RowID: "roleeditOrg",
				Label: "Organization",
				Help:  "Organization for which this role grants privileges",
			},
			Name:      "org",
			ValueP:    &ur.Org,
			Options:   enum.AllOrgs,
			LabelFunc: func(_ *request.Request, org enum.Org) string { return org.Label() },
		}, ur},
		&form.RadioGroupRow[enum.PrivLevel]{
			LabeledRow: form.LabeledRow{
				RowID: "roleeditPriv",
				Label: "Privilege",
				Help:  "Privilege level granted to people holding this role",
			},
			Name:    "priv",
			ValueP:  &ur.PrivLevel,
			Options: []enum.PrivLevel{0, enum.PrivStudent, enum.PrivMember, enum.PrivLeader},
			LabelFunc: func(_ *request.Request, priv enum.PrivLevel) string {
				if priv == 0 {
					return "None"
				}
				return priv.String()
			},
		},
		&form.SelectRow[uint]{
			LabeledRow: form.LabeledRow{
				RowID: "roleeditPriority",
				Label: "Sort before",
			},
			Name:      "priority",
			ValueP:    &prioBefore,
			Options:   prioOptions,
			ValueFunc: func(u uint) string { return strconv.Itoa(int(u)) },
			LabelFunc: func(_ *request.Request, v uint) string { return prioLabels[v] },
		},
		&form.FlagsRow[role.Flags]{
			CheckboxesRow: form.CheckboxesRow{
				LabeledRow: form.LabeledRow{Label: "Flags"},
				Validate:   form.NoValidate,
			},
			ValueP: &ur.Flags,
			Flags:  []role.Flags{role.Filter, role.ImplicitOnly, role.Archived},
			LabelFunc: func(_ *request.Request, v role.Flags) string {
				return map[role.Flags]string{
					role.Filter:       "Available choice on People list page",
					role.ImplicitOnly: "Role can only be implied, not assigned",
					role.Archived:     "Archived (no longer in use)",
				}[v]
			},
		},
		roleselect.NewRoleSelectRow(r, 0, func(rl *role.Role) bool {
			_, ok := impliers[rl.ID()]
			return !ok && rl.ID() != ur.ID
		}, "Implies", "implies", &ur.Implies, false),
		&listsRow{form.LabeledRow{Label: "Lists"}, ur},
	}
	f.Handle(r)
}

func setupPriorityList(r *request.Request, ur *role.Updater) (before uint, options []uint, labels map[uint]string) {
	labels = make(map[uint]string)
	role.All(r, role.FID|role.FName|role.FPriority, func(rl *role.Role) {
		if rl.ID() != ur.ID {
			if before == 0 && rl.Priority() > ur.Priority {
				before = rl.Priority()
			}
			options = append(options, rl.Priority())
			labels[rl.Priority()] = rl.Name()
		}
	})
	options = append(options, 0)
	labels[0] = "(at end)"
	return
}

type nameRow struct {
	form.TextInputRow
	ur *role.Updater
}

func (nr *nameRow) Read(r *request.Request) bool {
	if !nr.TextInputRow.Read(r) {
		return false
	}
	if nr.ur.Name == "" {
		nr.Error = "The role name is required."
		return false
	} else if nr.ur.DuplicateName(r) {
		nr.Error = "Another role has this name."
		return false
	}
	return true
}

type orgRow struct {
	form.RadioGroupRow[enum.Org]
	ur *role.Updater
}

func (or *orgRow) Read(r *request.Request) bool {
	if !or.RadioGroupRow.Read(r) {
		return false
	}
	if or.ur.Org == 0 {
		or.Error = "The organization is required."
		return false
	}
	return true
}

type listsRow struct {
	form.LabeledRow
	ur *role.Updater
}

func (lr *listsRow) Read(r *request.Request) bool { return true }

func (lr *listsRow) ShouldEmit(vl request.ValidationList) bool {
	return lr.ur.ID != 0 && !vl.Enabled()
}

func (lr *listsRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	var found bool

	row := lr.EmitPrefix(r, parent, "")
	box := row.E("div class=formInput")
	listrole.AllListsForRole(r, lr.ur.ID, func(l *list.List, sender bool, submodel listrole.SubscriptionModel) {
		found = true
		var name string
		if l.Type == list.SMS {
			name = "SMS: " + l.Name
		} else {
			name = l.Name + "@sunnyvaleserv.org"
		}
		div := box.E("div")
		div.E("a href=/admin/lists/%d up-target=main>%s", l.ID, name)
		div.R(" (")
		if submodel != 0 {
			div.R(submodel.String())
		}
		if submodel != 0 && sender {
			div.R(", ")
		}
		if sender {
			div.R("sender")
		}
		div.R(")")
	})
	if !found {
		box.R("None")
	}
	lr.EmitSuffix(r, row)
}

func saveRole(r *request.Request, user *person.Person, rl *role.Role, ur *role.Updater, prioBefore uint) bool {
	r.Transaction(func() {
		if rl == nil {
			rl = role.Create(r, ur)
		} else {
			rl.Update(r, ur)
		}
		switch {
		case prioBefore == 0:
			rl.Reorder(r, 0)
		case prioBefore < rl.Priority():
			rl.Reorder(r, prioBefore)
		case prioBefore > rl.Priority()+1:
			rl.Reorder(r, prioBefore-1)
		}
		role.Recalculate(r)
	})
	rolelist.Render(r, user)
	return true
}

func deleteRole(r *request.Request, user *person.Person, rl *role.Role) bool {
	r.Transaction(func() {
		rl.Delete(r)
		role.Recalculate(r)
	})
	rolelist.Render(r, user)
	return true
}
