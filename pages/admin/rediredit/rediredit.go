package rediredit

import (
	"net/url"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/redirlist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/redirect"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /admin/redirects/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user *person.Person
		rd   *redirect.Redirect
		ur   *redirect.Updater
		f    form.Form
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	f.Buttons = []*form.Button{{
		Label:   "Save",
		OnClick: func() bool { return saveRedirect(r, user, rd, ur) },
	}}
	if idstr == "NEW" {
		ur = new(redirect.Updater)
		f.Title = "New Redirect"
	} else {
		if rd = redirect.WithID(r, redirect.ID(util.ParseID(idstr))); rd == nil {
			errpage.NotFound(r, user)
			return
		}
		ur = rd.Updater()
		f.Title = "Edit Redirect"
		f.Buttons = append(f.Buttons, &form.Button{
			Name: "delete", Label: "Delete", Style: "danger",
			OnClick: func() bool { return deleteRedirect(r, user, rd) },
		})
	}
	f.Rows = []form.Row{
		&entryRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "redireditEntry",
				Label: "Entry",
			},
			Name:   "entry",
			ValueP: &ur.Entry,
		}, ur},
		&targetRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "redireditTarget",
				Label: "Target",
			},
			Name:   "target",
			ValueP: &ur.Target,
		}, ur},
	}
	f.Handle(r)
}

type entryRow struct {
	form.TextInputRow
	ur *redirect.Updater
}

func (er *entryRow) Read(r *request.Request) bool {
	if !er.TextInputRow.Read(r) {
		return false
	}
	if er.ur.Entry == "" {
		er.Error = "The redirect entry URL is required."
		return false
	} else if er.ur.DuplicateEntry(r) {
		er.Error = "Another redirect has this entry URL."
		return false
	} else if _, err := url.Parse(er.ur.Entry); err != nil || !strings.HasPrefix(er.ur.Entry, "/") {
		er.Error = "The entry URL must be a valid URL starting with '/'."
		return false
	}
	return true
}

type targetRow struct {
	form.TextInputRow
	ur *redirect.Updater
}

func (tr *targetRow) Read(r *request.Request) bool {
	if !tr.TextInputRow.Read(r) {
		return false
	}
	if _, err := url.Parse(tr.ur.Target); err != nil {
		tr.Error = "The map URL must be a valid URL."
		return false
	}
	return true
}

func saveRedirect(r *request.Request, user *person.Person, rd *redirect.Redirect, ur *redirect.Updater) bool {
	r.Transaction(func() {
		if rd == nil {
			rd = redirect.Create(r, ur)
		} else {
			rd.Update(r, ur)
		}
	})
	redirlist.Render(r, user)
	return true
}

func deleteRedirect(r *request.Request, user *person.Person, rd *redirect.Redirect) bool {
	r.Transaction(func() {
		rd.Delete(r)
	})
	redirlist.Render(r, user)
	return true
}
