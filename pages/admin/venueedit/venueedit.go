package venueedit

import (
	"net/url"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/venuelist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /admin/venues/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user *person.Person
		v    *venue.Venue
		uv   *venue.Updater
		f    form.Form
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	f.Attrs = "method=POST up-target=main"
	f.Dialog = true
	f.Buttons = []*form.Button{{
		Label:   "Save",
		OnClick: func() bool { return saveVenue(r, user, v, uv) },
	}}
	if idstr == "NEW" {
		uv = new(venue.Updater)
		f.Title = "New Venue"
	} else {
		if v = venue.WithID(r, venue.ID(util.ParseID(idstr)), venue.UpdaterFields); v == nil {
			errpage.NotFound(r, user)
			return
		}
		uv = v.Updater()
		f.Title = "Edit Venue"
		if !event.ExistsWithVenue(r, v.ID()) && !shift.ExistsWithVenue(r, v.ID()) {
			f.Buttons = append(f.Buttons, &form.Button{
				Name: "delete", Label: "Delete", Style: "danger",
				OnClick: func() bool { return deleteVenue(r, user, v) },
			})
		}
	}
	f.Rows = []form.Row{
		&nameRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "venueeditName",
				Label: "Name",
			},
			Name:   "name",
			ValueP: &uv.Name,
		}, uv},
		&urlRow{form.TextInputRow{
			LabeledRow: form.LabeledRow{
				RowID: "venueeditURL",
				Label: "MapURL",
				Help:  "Google Maps URL for the venue.  Should be zoomed out enough to see major cross-streets or freeways.",
			},
			Name:   "url",
			ValueP: &uv.URL,
		}, uv},
		&form.FlagsRow[venue.Flag]{
			CheckboxesRow: form.CheckboxesRow{
				LabeledRow: form.LabeledRow{Label: "Flags"},
				Validate:   form.NoValidate,
				Name:       "flags",
			},
			ValueP: &uv.Flags,
			Flags:  []venue.Flag{venue.CanOverlap},
			LabelFunc: func(_ *request.Request, v venue.Flag) string {
				return map[venue.Flag]string{
					venue.CanOverlap: "Venue can have simultaneous events",
				}[v]
			},
		},
	}
	f.Handle(r)
}

type nameRow struct {
	form.TextInputRow
	uv *venue.Updater
}

func (nr *nameRow) Read(r *request.Request) bool {
	if !nr.TextInputRow.Read(r) {
		return false
	}
	if nr.uv.Name == "" {
		nr.Error = "The venue name is required."
		return false
	} else if nr.uv.DuplicateName(r) {
		nr.Error = "Another venue has this name."
		return false
	}
	return true
}

type urlRow struct {
	form.TextInputRow
	uv *venue.Updater
}

func (ur *urlRow) Read(r *request.Request) bool {
	if !ur.TextInputRow.Read(r) {
		return false
	}
	if ur.uv.URL == "" {
		return true
	} else if _, err := url.Parse(ur.uv.URL); err != nil {
		ur.Error = "The map URL must be a valid URL."
		return false
	} else if !strings.HasPrefix(ur.uv.URL, "https://www.google.com/maps/") {
		ur.Error = "The map URL must begin with https://www.google.com/maps/."
		return false
	}
	return true
}

func saveVenue(r *request.Request, user *person.Person, v *venue.Venue, uv *venue.Updater) bool {
	r.Transaction(func() {
		if v == nil {
			v = venue.Create(r, uv)
		} else {
			v.Update(r, uv)
		}
	})
	venuelist.Render(r, user)
	return true
}

func deleteVenue(r *request.Request, user *person.Person, v *venue.Venue) bool {
	r.Transaction(func() {
		v.Delete(r)
	})
	venuelist.Render(r, user)
	return true
}
