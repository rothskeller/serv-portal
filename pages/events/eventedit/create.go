package eventedit

import (
	"fmt"
	"net/http"
	"slices"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskrole"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

// HandleCreate handles /events/create requests.
func HandleCreate(r *request.Request) {
	var (
		user       *person.Person
		ue         *event.Updater
		ut         *task.Updater
		roles      []*role.Role
		nameError  string
		dateError  string
		timesError string
		venueError string
		orgError   string
		hasError   bool
		month      string
		opts       ui.PageOpts
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	ue = new(event.Updater)
	ut = new(task.Updater)
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		nameError = readEventName(r, ue)
		readActivation(r, ue)
		dateError = readDate(r, ue)
		timesError = readEventTimes(r, ue)
		venueError = readEventVenue(r, ue)
		orgError = readOrg(r, ut)
		roles = readRoles(r, user, roles)
		readTaskFlags(r, ut)
		readEventDetails(r, ue)
		hasError = nameError != "" || dateError != "" || timesError != "" || venueError != "" || orgError != ""
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !hasError {
			r.Transaction(func() {
				ut.Event = event.Create(r, ue)
				ut.Name = "Tracking"
				t := task.Create(r, ut)
				taskrole.Set(r, ut.Event, t, roles, []*role.Role{})
				if ut.Flags&task.SignupsOpen != 0 {
					shift.Create(r, &shift.Updater{
						Event: ut.Event,
						Task:  t,
						Start: ut.Event.Start(),
						End:   ut.Event.End(),
						Venue: ue.Venue,
					})
				}
			})
			http.Redirect(r, r.Request, fmt.Sprintf("/events/%d", ut.Event.ID()), http.StatusSeeOther)
			return
		}
	}
	r.HTMLNoCache()
	if hasError {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	month = state.GetEventsMonth(r)
	opts = ui.PageOpts{
		Title:    "New Event",
		MenuItem: "events",
		Tabs: []ui.PageTab{
			{Name: "Calendar", URL: "/events/calendar/" + month, Target: ".pageCanvas"},
			{Name: "List", URL: "/events/list/" + month[:4], Target: ".pageCanvas"},
			{Name: "Signups", URL: "/events/signups", Target: ".pageCanvas"},
			{Name: "Add Event", URL: "/events/create", Target: "main", Active: true},
		},
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		form := main.E("form class=form method=POST up-main")
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		if len(validate) == 0 || slices.Contains(validate, "name") || slices.Contains(validate, "date") {
			emitEventName(form, ue, nameError != "" || !hasError, nameError)
		}
		if len(validate) == 0 || slices.Contains(validate, "activation") {
			emitActivation(form, ue)
		}
		if len(validate) == 0 || slices.Contains(validate, "date") {
			emitDate(form, ue, dateError != "", dateError)
		}
		if len(validate) == 0 || slices.Contains(validate, "start") || slices.Contains(validate, "end") {
			emitEventTimes(form, ue, timesError != "", timesError)
		}
		if len(validate) == 0 || slices.Contains(validate, "venue") || slices.Contains(validate, "venueURL") {
			emitEventVenue(form, ue, venueError != "", venueError)
		}
		if len(validate) == 0 || slices.Contains(validate, "org") {
			emitOrg(form, user, ut, orgError != "", orgError)
		}
		if len(validate) == 0 {
			emitRoles(r, form, user, roles)
			emitTaskFlags(form, ut, "event")
			emitEventDetails(form, ue)
			emitCreateButtons(form)
		}
	})
}

func emitCreateButtons(form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
}
