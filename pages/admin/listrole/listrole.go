package listrole

import (
	"sunnyvaleserv.org/portal/pages/admin/roleselect"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles /admin/lists/$lid/roleedit/$rid requests, where either ID may be
// "NEW".
func Get(r *request.Request, lidstr, ridstr string) {
	var (
		user     *person.Person
		ridlist  []role.ID
		sender   bool
		submodel listrole.SubscriptionModel
		f        form.Form
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if rid := role.ID(util.ParseID(ridstr)); rid > 0 {
		ridlist = append(ridlist, rid)
	}
	sender, submodel = listrole.Get(r, list.ID(util.ParseID(lidstr)), role.ID(util.ParseID(ridstr)))

	f.Attrs = "class=listeditRoleForm"
	f.Dialog, f.NoSubmit, f.TwoCol = true, true, true
	f.Title = "Role Privileges"
	f.Buttons = []*form.Button{{Label: "OK"}}
	f.Rows = []form.Row{
		roleselect.NewRoleSelectRow(r, 0, nil, "Role(s)", "roles", &ridlist, false),
		&form.RadioGroupRow[listrole.SubscriptionModel]{
			LabeledRow: form.LabeledRow{
				RowID: "listroleSubmodel",
				Label: "Subscription",
			},
			Name:    "submodel",
			ValueP:  &submodel,
			Options: listrole.AllSubscriptionModels,
			LabelFunc: func(r *request.Request, sm listrole.SubscriptionModel) string {
				return r.Loc(sm.LongString())
			},
		},
		&form.CheckboxesRow{
			LabeledRow: form.LabeledRow{
				RowID: "listroleSender",
				Label: "Sender",
			},
			Boxes: []*form.Checkbox{{
				Name:     "sender",
				Label:    "Can send without moderation",
				CheckedP: &sender,
			}},
		},
	}
	f.Handle(r)
}
