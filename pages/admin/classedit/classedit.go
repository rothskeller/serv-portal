package classedit

import (
	"cmp"
	"slices"
	"strconv"

	"sunnyvaleserv.org/portal/pages/admin/classlist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /admin/classes/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user         *person.Person
		c            *class.Class
		uc           *class.Updater
		f            form.Form
		canDelete    bool
		roleID       role.ID
		studentRoles []role.ID
		roleMap      = map[role.ID]string{0: "(none)"}
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if idstr == "NEW" {
		uc = new(class.Updater)
	} else {
		if c = class.WithID(r, class.ID(util.ParseID(idstr)), class.UpdaterFields); c == nil {
			errpage.NotFound(r, user)
			return
		}
		uc = c.Updater(r, nil)
		canDelete = !classreg.ClassHasSignups(r, c.ID())
		roleID = c.Role()
	}
	role.All(r, role.FID|role.FName|role.FPrivLevel|role.FFlags, func(rl *role.Role) {
		if rl.PrivLevel() == enum.PrivStudent && rl.Flags()&(role.Archived|role.ImplicitOnly) == 0 {
			studentRoles = append(studentRoles, rl.ID())
			roleMap[rl.ID()] = rl.Name()
		}
	})
	slices.SortFunc(studentRoles, func(a, b role.ID) int { return cmp.Compare(roleMap[a], roleMap[b]) })
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	if c == nil {
		f.Title = "New Class"
	} else {
		f.Title = "Edit Class"
	}
	f.Rows = []form.Row{
		&typeRow{form.SelectRow[class.Type]{
			LabeledRow: form.LabeledRow{
				RowID: "classeditType",
				Label: "Type",
			},
			Name:        "type",
			ValueP:      &uc.Type,
			Options:     class.AllTypes,
			Placeholder: "(select type)",
			Validate:    "#classeditType,#classeditStart",
		}},
		&startRow{form.DateRow{InputRow: form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "classeditStart",
				Label: "Start Date",
				Help:  "Use 2999-12-31 for a waiting list placeholder.",
			},
			Name:   "start",
			ValueP: &uc.Start,
		}}, uc},
		&reqDescRow{form.TextAreaRow{
			LabeledRow: form.LabeledRow{
				RowID: "classeditEnDesc",
				Label: "English Desc.",
				Help:  "Include date(s), time(s), location(s), maybe language",
			},
			Name:   "enDesc",
			ValueP: &uc.EnDesc,
		}},
		&reqDescRow{form.TextAreaRow{
			LabeledRow: form.LabeledRow{
				RowID: "classeditEsDesc",
				Label: "Spanish Desc.",
			},
			Name:   "esDesc",
			ValueP: &uc.EsDesc,
		}},
		&form.IntegerRow[uint]{
			InputRow: form.InputRow{
				LabeledRow: form.LabeledRow{
					RowID: "classeditLimit",
					Label: "Enrollment Limit",
					Help:  "0 = unlimited",
				},
				Name: "limit",
			},
			ValueP: &uc.Limit,
		},
		&form.InputRow{
			LabeledRow: form.LabeledRow{
				RowID: "classeditRegURL",
				Label: "External URL",
				Help:  "URL to external registration page",
			},
			Name:   "regURL",
			ValueP: &uc.RegURL,
		},
		&form.SelectRow[role.ID]{
			LabeledRow: form.LabeledRow{
				RowID: "classeditRole",
				Label: "Student Role",
			},
			Name:      "role",
			ValueP:    &roleID,
			Options:   studentRoles,
			ValueFunc: func(id role.ID) string { return strconv.Itoa(int(id)) },
			LabelFunc: func(r *request.Request, v role.ID) string { return roleMap[v] },
		},
		&referralsRow{form.LabeledRow{Label: "Referrals"}, uc},
	}
	f.Buttons = []*form.Button{{
		Label:   "Save",
		OnClick: func() bool { return saveClass(r, user, c, uc) },
	}}
	if canDelete {
		f.Buttons = append(f.Buttons, &form.Button{
			Name: "delete", Label: "Delete", Style: "danger",
			OnClick: func() bool { return deleteClass(r, user, c) },
		})
	}
	if uc.ID != 0 {
		f.Buttons = append(f.Buttons, &form.Button{
			Name: "copy", Label: "Save Copy", Style: "secondary",
			OnClick: func() bool {
				uc.ID = 0
				return saveClass(r, user, nil, uc)
			},
		})
	}
	f.Handle(r)
}

type typeRow struct{ form.SelectRow[class.Type] }

func (tr *typeRow) Read(r *request.Request) bool {
	if !tr.SelectRow.Read(r) {
		return false
	}
	if *tr.ValueP == 0 {
		tr.Error = "The class type is required."
		return false
	}
	return true
}

type startRow struct {
	form.DateRow
	uc *class.Updater
}

func (sr *startRow) ShouldEmit(vl request.ValidationList) bool {
	return vl.ValidatingAny("type", "start")
}

func (sr *startRow) Read(r *request.Request) bool {
	if !sr.DateRow.Read(r) {
		return false
	}
	if sr.uc.Start == "" {
		sr.Error = "The class starting date is required."
		return false
	}
	if sr.uc.DuplicateStart(r) {
		sr.Error = "Another class has the same type and start date."
		return false
	}
	return true
}

type reqDescRow struct{ form.TextAreaRow }

func (rdr *reqDescRow) Read(r *request.Request) bool {
	if !rdr.TextAreaRow.Read(r) {
		return false
	}
	if *rdr.ValueP == "" {
		rdr.Error = "The description is required."
		return false
	}
	return true
}

type referralsRow struct {
	form.LabeledRow
	uc *class.Updater
}

func (rr *referralsRow) ShouldEmit(_ request.ValidationList) bool {
	return len(rr.uc.Referrals) != 0
}

func (rr *referralsRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	row := rr.EmitPrefix(r, parent, "")
	grid := row.E("div class='classeditReferrals formInput'")
	for _, ref := range class.AllReferrals {
		grid.E("div>%d", rr.uc.Referrals[ref])
		grid.E("div>%s", ref.String())
	}
	rr.EmitSuffix(r, row)
}

func (rr *referralsRow) Read(r *request.Request) bool { return true }

func saveClass(r *request.Request, user *person.Person, c *class.Class, ur *class.Updater) bool {
	r.Transaction(func() {
		if c == nil {
			c = class.Create(r, ur)
		} else {
			c.Update(r, ur)
		}
	})
	classlist.Render(r, user)
	return true
}

func deleteClass(r *request.Request, user *person.Person, c *class.Class) bool {
	r.Transaction(func() {
		c.Delete(r)
	})
	classlist.Render(r, user)
	return true
}
