package eventview

import (
	"fmt"
	"net/http"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

const EventFields = event.FID | event.FStart | event.FVenue | identEventFields | detailsEventFields | taskEventFields

// Handle handles /events/${id} requests.
func Handle(r *request.Request, idstr string) {
	const personFields = taskPersonFields
	var (
		user *person.Person
		e    *event.Event
	)
	if user = auth.SessionUser(r, personFields, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if e = event.WithID(r, event.ID(util.ParseID(idstr)), EventFields); e == nil {
		errpage.NotFound(r, user)
		return
	}
	state.SetEventsMonth(r, e.Start()[0:7])
	if r.Method == http.MethodPost {
		if handleDelete(r, e) {
			return
		}
		handleSignup(r, user, e)
		handleHours(r, user, e)
	}
	Render(r, user, e, "")
}

// Render renders the event view page, or a particular section of it.  It is
// called by Get, above, and also by the edit dialogs after accepting a change
// to an event.
func Render(r *request.Request, user *person.Person, e *event.Event, section string) {
	const taskFields = task.FID | task.FOrg | identTaskFields | detailsTaskFields | taskTaskFields
	var ts []*task.Task
	var canDelete = !shiftperson.EventHasSignups(r, e.ID()) && !taskperson.ExistsForEvent(r, e.ID())
	var canAddTask = user.HasPrivLevel(0, enum.PrivLeader)
	var canCopy = user.HasPrivLevel(0, enum.PrivLeader)

	task.AllForEvent(r, e.ID(), taskFields, func(t *task.Task) {
		var clone = *t
		ts = append(ts, &clone)
		if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
			canDelete, canCopy = false, false
		}
	})
	opts := ui.PageOpts{
		Title:    e.Start()[:10],
		Banner:   e.Start()[:10] + " " + e.Name(),
		MenuItem: "events",
		Tabs: []ui.PageTab{
			{Name: r.Loc("Calendar"), URL: "/events/calendar/" + e.Start()[0:7], Target: ".pageCanvas"},
			{Name: r.Loc("List"), URL: "/events/list/" + e.Start()[0:4], Target: ".pageCanvas"},
			{Name: r.Loc("Details"), URL: fmt.Sprintf("/events/%d", e.ID()), Target: "main", Active: true},
		},
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		box := main.E("div class=eventview")
		if section == "" || section == "details" {
			showIdent(r, box, e, ts)
			showDetails(r, box, user, e, ts)
		}
		for _, t := range ts {
			if section == "" || section == fmt.Sprintf("task%d", t.ID()) {
				showTask(r, box, user, e, t, len(ts) > 1)
			}
		}
		if section == "" && (canAddTask || canDelete || canCopy) {
			buttons := main.E("form class=eventviewButtons method=POST")
			buttons.E("input type=hidden name=csrf value=%s", r.CSRF)
			if canAddTask {
				buttons.E("a href=/events/edtask/NEW?eid=%d up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add Task", e.ID())
			}
			if canCopy {
				buttons.E("a href=/events/%d/copy up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Copy Event", e.ID())
			}
			if canDelete {
				buttons.E("input name=delete type=submit class='sbtn sbtn-danger' value='Delete Event'")
			}
		}
	})
}

func handleDelete(r *request.Request, e *event.Event) bool {
	if r.FormValue("delete") != "" && !shiftperson.EventHasSignups(r, e.ID()) && !taskperson.ExistsForEvent(r, e.ID()) {
		r.Transaction(func() {
			e.Delete(r)
		})
		http.Redirect(r, r.Request, state.GetEventsURL(r), http.StatusSeeOther)
		return true
	}
	return false
}
