package eventslist

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

// Get handles GET /events/list/${year} requests.
func Get(r *request.Request, yearstr string) {
	const eventFields = event.FID | event.FName | event.FStart | event.FVenue | event.FFlags
	var (
		user  *person.Person
		opts  ui.PageOpts
		month string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if year, err := strconv.Atoi(yearstr); err != nil || year < 2000 || year >= 2100 {
		errpage.NotFound(r, user)
		return
	}
	state.SetEventsPage(r, "list")
	if month = state.GetEventsMonth(r); month[0:4] != yearstr {
		if month = time.Now().Format("2006-01"); month[0:4] != yearstr {
			month = yearstr + "-01"
		}
	}
	state.SetEventsMonth(r, month)
	opts = ui.PageOpts{
		Title:    r.LangString("Events", "Eventos"),
		MenuItem: "events",
		Tabs: []ui.PageTab{
			{Name: r.LangString("Calendar", "Calendario"), URL: "/events/calendar/" + month, Target: "main"},
			{Name: r.LangString("List", "Lista"), URL: "/events/list/" + month[0:4], Target: "main", Active: true},
			{Name: r.LangString("Signups", "Inscripciones"), URL: "/events/signups", Target: "main"},
		},
	}
	if user.HasPrivLevel(0, enum.PrivLeader) {
		opts.Tabs = append(opts.Tabs, ui.PageTab{Name: "Add Event", URL: "/events/create", Target: "main"})
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		var (
			lastDate   string
			dateEvents []event.Event
			venueCache = make(map[venue.ID]*venue.Venue)
		)
		main.A("class=eventslist")
		main.E("div class=eventslistTitle").E("s-year id=eventslistYear value=%s", month[0:4])
		table := main.E("div class=eventslistTable")
		table.E("div class='eventslistHeading eventslistDate'").R(r.LangString("Date", "Fecha"))
		table.E("div class='eventslistHeading eventslistEvent'").R(r.LangString("Event", "Evento"))
		table.E("div class='eventslistHeading eventslistLocation'").R(r.LangString("Location", "Sitio"))
		// Walk through all of the tasks in the year, gathering them up
		// by date and emitting all of them on the same date together.
		event.AllBetween(r, yearstr+"-01-01", yearstr+"-12-32", eventFields, 0, func(e *event.Event, _ *venue.Venue) {
			if e.Flags()&event.OtherHours != 0 {
				return
			}
			if e.Start()[:10] != lastDate && lastDate != "" {
				emitDateEvents(r, table, dateEvents, venueCache)
				dateEvents = dateEvents[:0]
			}
			dateEvents = append(dateEvents, *e)
			lastDate = e.Start()[:10]
		})
		if len(dateEvents) != 0 {
			emitDateEvents(r, table, dateEvents, venueCache)
		}
	})
}

// emitDateEvents emits all of the events on a given date.  This function
// updates the venueCache map.
func emitDateEvents(r *request.Request, table *htmlb.Element, events []event.Event, venueCache map[venue.ID]*venue.Venue) {
	for i := range events {
		var orgs = make([]bool, enum.NumOrgs)

		task.AllForEvent(r, events[i].ID(), task.FOrg, func(t *task.Task) {
			orgs[t.Org()] = true
		})
		date := table.E("div class=eventslistDate")
		date.E("span class=eventslistYear").R(events[i].Start()[:5])
		date.E("span").R(events[i].Start()[5:10])
		date.E("span class=eventslistStart").R(events[i].Start()[11:])
		ediv := table.E("div class=eventslistEvent")
		for org := enum.Org(0); org < enum.NumOrgs; org++ {
			if orgs[org] {
				orgdot.OrgDot(r, ediv, org)
			}
		}
		ediv.E("a up-target=.pageCanvas href=/events/%d", events[i].ID()).T(events[i].Name())
		loc := table.E("div class=eventslistLocation")
		switch vid := events[i].Venue(); vid {
		case 0:
			loc.R(r.LangString("TBD", "Por determinar"))
		default:
			var v *venue.Venue
			if v = venueCache[vid]; v == nil {
				venueCache[vid] = venue.WithID(r, vid, venue.FName|venue.FURL)
				v = venueCache[vid]
			}
			if v.URL() != "" {
				loc.E("a href=%s target=_blank", v.URL()).T(v.Name())
			} else {
				loc.T(v.Name())
			}
		}
	}
}
