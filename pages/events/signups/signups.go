package signups

import (
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

// Handle handles /events/signups and /events/signups/${token} requests.
func Handle(r *request.Request, token string) {
	var (
		user      *person.Person
		opts      ui.PageOpts
		month     string
		fromToken bool
		roles     = map[role.ID]bool{}
	)
	if user = auth.SessionUser(r, ShowTaskSignupsPersonFields, false); user == nil && token != "" {
		fromToken = true
		user = person.WithUnsubscribeToken(r, token, ShowTaskSignupsPersonFields)
	}
	if user == nil {
		http.Redirect(r, r.Request, "/login"+r.Path, http.StatusSeeOther)
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	personrole.RolesForPerson(r, user.ID(), role.FID, func(rl *role.Role, _ bool) {
		roles[rl.ID()] = true
	})
	if r.Method == http.MethodPost {
		HandleShiftSignup(r, user, user)
	}
	month = state.GetEventsMonth(r)
	opts = ui.PageOpts{
		Title:    "Event Signups",
		MenuItem: "events",
	}
	if !fromToken {
		opts.Tabs = []ui.PageTab{
			{Name: "Calendar", URL: "/events/calendar/" + month, Target: "main"},
			{Name: "List", URL: "/events/list/" + month[0:4], Target: "main"},
			{Name: "Signups", URL: "/events/signups", Target: "main", Active: true},
		}
		if user.HasPrivLevel(0, enum.PrivLeader) {
			opts.Tabs = append(opts.Tabs, ui.PageTab{Name: "Add Event", URL: "/events/create", Target: "main"})
		}
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		const eventFields = event.FID | event.FStart | event.FName | event.FDetails
		const taskFields = task.FID | task.FName | task.FDetails | task.FOrg | ShowTaskSignupsTaskFields
		var (
			form  *htmlb.Element
			lastd string
			ediv  *htmlb.Element
			laste event.ID
			tdiv  *htmlb.Element
			lastt task.ID
		)
		shift.AllAfter(r, time.Now().Format("2006-01-02T15:04"), eventFields, taskFields, shift.FID, 0,
			func(e *event.Event, t *task.Task, s *shift.Shift, _ *venue.Venue) {
				if e.Start()[:10] != lastd {
					lastd = e.Start()[:10]
					form = main.E("form id=d%s class=signupForm method=POST up-submit up-target=#d%s", lastd, lastd)
					form.E("input type=hidden name=csrf value=%s", r.CSRF)
					form.E("input type=hidden name=shift")
					form.E("input type=hidden name=signedup")
				}
				if e.ID() != laste {
					date, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
					ediv = form.E("div class=signupEvent")
					ediv.E("div class=signupEventName>%s", e.Name())
					ediv.E("div class=signupEventDate>%s", date.Format("Monday, January 2, 2006"))
					if e.Details() != "" {
						form.E("div class=signupEventDetails").R(e.Details())
					}
					laste = e.ID()
				}
				if t.ID() != lastt {
					tdiv = ediv.E("div class=signupTask")
					tdiv.E("div class=signupTaskName>%s", t.Name())
					if t.Details() != "" {
						tdiv.E("div class=signupTaskDetails").R(t.Details())
					}
					ShowTaskSignups(r, tdiv, t, user, user.HasPrivLevel(t.Org(), enum.PrivLeader), false)
					lastt = t.ID()
				}
			})
		if lastd == "" {
			main.E("div>There are no upcoming events with signups.")
		}
	})
}
