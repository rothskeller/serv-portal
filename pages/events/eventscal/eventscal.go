package eventscal

import (
	"strconv"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/ui/orgdot"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

// Get handles GET /events/calendar/${month} requests.
func Get(r *request.Request, month string) {
	const eventFields = event.FID | event.FStart | event.FName | event.FFlags
	var (
		user  *person.Person
		opts  ui.PageOpts
		start time.Time
		err   error
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if start, err = time.ParseInLocation("2006-01", month, time.Local); err != nil || start.Format("2006-01") != month {
		errpage.NotFound(r, user)
		return
	}
	state.SetEventsPage(r, "calendar")
	state.SetEventsMonth(r, month)
	opts = ui.PageOpts{
		Title:    "Events",
		MenuItem: "events",
		Tabs: []ui.PageTab{
			{Name: "Calendar", URL: "/events/calendar/" + month, Target: "main", Active: true},
			{Name: "List", URL: "/events/list/" + month[0:4], Target: "main"},
			{Name: "Signups", URL: "/events/signups", Target: "main"},
		},
	}
	if user.HasPrivLevel(0, enum.PrivLeader) {
		opts.Tabs = append(opts.Tabs, ui.PageTab{Name: "Add Event", URL: "/events/create", Target: "main"})
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		var (
			lastDay = start
			events  []event.Event
		)
		main.A("class=eventscal")
		grid := main.E("div class=eventscalGrid")
		grid.E("div class=eventscalHeading").E("s-month id=eventscalMonth value=%s", month)
		grid.E("div class=eventscalWeekday>S")
		grid.E("div class=eventscalWeekday>M")
		grid.E("div class=eventscalWeekday>T")
		grid.E("div class=eventscalWeekday>W")
		grid.E("div class=eventscalWeekday>T")
		grid.E("div class=eventscalWeekday>F")
		grid.E("div class=eventscalWeekday>S")
		for i := time.Sunday; i < start.Weekday(); i++ {
			grid.E("div class='eventscalDay eventscalDay-empty'")
		}
		event.AllBetween(r, month+"-01", month+"-32", eventFields, 0, func(e *event.Event, _ *venue.Venue) {
			if e.Flags()&event.OtherHours != 0 {
				return
			}
			for lastDay.Format("2006-01-02") != e.Start()[:10] {
				emitDayCell(r, grid, lastDay, events)
				events = events[:0]
				lastDay = lastDay.AddDate(0, 0, 1)
			}
			events = append(events, *e)
		})
		for lastDay.Weekday() != time.Sunday {
			if lastDay.Format("2006-01") == month {
				emitDayCell(r, grid, lastDay, events)
				events = events[:0]
			} else {
				grid.E("div class='eventscalDay eventscalDay-empty'")
			}
			lastDay = lastDay.AddDate(0, 0, 1)
		}
		main.E("div id=eventscalFooter")
	})
}

// emitDayCell emits the calendar grid cell for the specified day, with all of
// the events on that day.
func emitDayCell(r *request.Request, grid *htmlb.Element, day time.Time, events []event.Event) {
	cell := grid.E("div class=eventscalDay data-date=%s", day.Format("Monday, January 2, 2006"))
	cell.E("div").R(strconv.Itoa(day.Day()))
	if len(events) == 0 {
		return
	}
	evs := cell.E("div class=eventscalEvents")
	for i := range events {
		var orgs = make([]bool, enum.NumOrgs)

		task.AllForEvent(r, events[i].ID(), task.FOrg, func(t *task.Task) {
			orgs[t.Org()] = true
		})
		ev := evs.E("div class=eventscalEvent")
		for org := enum.Org(0); org < enum.NumOrgs; org++ {
			if orgs[org] {
				orgdot.OrgDot(ev, org)
			}
		}
		ev.E("a href=/events/%d up-target=.pageCanvas class=eventscalEventLink title=%s",
			events[i].ID(), events[i].Name()).T(events[i].Name())
	}
}
