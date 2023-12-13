package eventattend

import (
	"cmp"
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/store/taskrole"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

/* RECORD ATTENDANCE Dialog

The dialog shows a list of people and the hours, signed-in, and/or credited
flags for each.  The people are listed in alphabetical order by sortname, except
for any people added to the list while the dialog is open, who are added at the
bottom.

Each line of the list starts with the hours entryfield, but only if the task
allows recording of hours.  After that are signed-in and credited icons, which
are gray outlines for false and filled/colored for true, and which can be
toggled by clicking on them.  After that is the person's sortname.

Above the list is a row containing default values, set off from the list with
background coloration and margin.  The name column of this row says "Default
values for task\n(Click name to apply)".  When the dialog opens, these controls
are set as follows:
  - The hours entry (if present) is set to the duration of the event.
  - The signed-in flag is on.
  - The credited flag is on if hours are not being collected or the task already
    has people with credit.

The set of people initially shown includes:
  - Everyone who already has hours, signed-in, or credited for the task.
  - If the task has shifts, everyone signed up for any shift.
  - If the task has no shifts, everyone holding one of the task roles, but only
    if that list is <= 100 people.

On each person's line:
  - Entering a number of hours sets that person's hours.
  - Entering a timesheet pair in the hours box sets that person's hours and sets
    the signed-in flag.
  - Clicking on the signed-in or credited flags toggles them.
  - Clicking on the person's name does the following:
    - If the person's hours, signed-in, and credited flags exactly match the
      default line, the hours are zeroed and the signed-in and credited flags
      are turned off.  Otherwise:
    - The person's signed-in and credited flags are set to match the default
      line.
    - The person's hours are increased, but not decreased, to match the default
      line.

The last line of the list shows a search entry field in place of the person
name.  If a person is selected in that entry field, the line for that person is
set to match the defaults, and a new line is added with another empty entry
field.

No changes are committed until the Save button is pressed.  No entry errors are
possible, so there is no validation.
*/

// Handle handles /events/attendance/$tid requests.
func Handle(r *request.Request, tidstr string) {
	const eventFields = event.FID | event.FName | event.FStart | event.FEnd | event.FFlags | eventview.EventFields
	const taskFields = task.FID | task.FName | task.FEvent | task.FOrg | task.FFlags
	var (
		user *person.Person
		e    *event.Event
		t    *task.Task
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if t = task.WithID(r, task.ID(util.ParseID(tidstr)), taskFields); t == nil {
		errpage.NotFound(r, user)
		return
	}
	e = event.WithID(r, t.Event(), eventFields)
	if !user.HasPrivLevel(t.Org(), enum.PrivLeader) || e.Flags()&event.OtherHours != 0 {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		postAttendance(r, user, e, t)
	} else {
		getAttendance(r, user, e, t)
	}
}

func getAttendance(r *request.Request, user *person.Person, e *event.Event, t *task.Task) {
	type adata struct {
		id      person.ID
		name    string
		minutes uint
		flags   taskperson.Flag
	}
	var (
		hasShifts       bool
		defaultCredited bool
		eventStart      time.Time
		eventEnd        time.Time
		eventDuration   uint
		plist           []*adata
		people          = make(map[person.ID]*adata)
	)
	// Show all people who are signed up for the task's shifts if any.
	shift.AllForTask(r, t.ID(), shift.FID, 0, func(s *shift.Shift, _ *venue.Venue) {
		hasShifts = true
		shiftperson.PeopleForShift(r, s.ID(), person.FID|person.FSortName, func(p *person.Person) {
			people[p.ID()] = &adata{id: p.ID(), name: p.SortName()}
		})
	})
	// If we didn't find any shifts, then show all people who have eligible
	// roles, but not if that's more than 100 people.
	if !hasShifts {
		taskrole.Get(r, t.ID(), role.FID, func(rl *role.Role) {
			personrole.PeopleForRole(r, rl.ID(), person.FID|person.FSortName, func(p *person.Person, _ bool) {
				people[p.ID()] = &adata{id: p.ID(), name: p.SortName()}
			})
		})
		if len(people) > 100 {
			clear(people)
		}
	}
	// In either case, always show the people who already have attendance
	// recorded.
	taskperson.PeopleForTask(r, t.ID(), person.FID|person.FSortName, func(p *person.Person, minutes uint, flags taskperson.Flag) {
		people[p.ID()] = &adata{p.ID(), p.SortName(), minutes, flags}
		if flags&taskperson.Credited != 0 {
			defaultCredited = true
		}
	})
	plist = make([]*adata, 0, len(people))
	for _, p := range people {
		plist = append(plist, p)
	}
	slices.SortFunc(plist, func(a, b *adata) int { return cmp.Compare(a.name, b.name) })
	// Determine the default hours for attendance recording.
	eventStart, _ = time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
	eventEnd, _ = time.ParseInLocation("2006-01-02T15:04", e.End(), time.Local)
	eventDuration = uint(eventEnd.Sub(eventStart) / time.Minute)
	eventDuration = ((eventDuration + 19) / 30) * 30
	// Start the form.
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=#eventviewTask%d", t.ID())
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("div class='formTitle formTitle-primary'>Record Attendance")
	outer := form.E("div class=formRow-3col")
	grid := outer.E("div class=attendanceDefault")
	// Show the defaults row.
	row := grid.E("div class=attendanceRow")
	if t.Flags()&task.RecordHours != 0 {
		row.E("s-hours class=attendanceHours value=%s", ui.MinutesToHours(eventDuration))
	} else {
		defaultCredited = true
	}
	box := row.E("div class='attendanceSignedIn true'")
	box.E("input type=hidden value=true")
	box.E("s-icon icon=signature")
	box = row.E("div class=attendanceCredited", defaultCredited, "class=true")
	box.E("s-icon", defaultCredited, "icon=star-solid", !defaultCredited, "icon=star")
	box.E("input type=hidden value=%v", defaultCredited)
	row.E("div class=attendanceName><b>Defaults for this task</b><br>(Click name to apply)")
	// Show the grid with each person's row.
	grid = outer.E("div class=attendanceGrid")
	for _, p := range plist {
		row = grid.E("div class=attendanceRow")
		if t.Flags()&task.RecordHours != 0 {
			row.E("s-hours class=attendanceHours name=hours%d value=%s", p.id, ui.MinutesToHours(p.minutes))
		}
		box := row.E("div class=attendanceSignedIn", p.flags&taskperson.Attended != 0, "class=true")
		box.E("s-icon icon=signature")
		box.E("input type=hidden name=signedin%d value=%v", p.id, p.flags&taskperson.Attended != 0)
		box = row.E("div class=attendanceCredited", p.flags&taskperson.Credited != 0, "class=true")
		box.E("s-icon", p.flags&taskperson.Credited != 0, "icon=star-solid", p.flags&taskperson.Credited == 0, "icon=star")
		box.E("input type=hidden name=credited%d value=%v", p.id, p.flags&taskperson.Credited != 0)
		row.E("div class=attendanceName data-key=P%d>%s", p.id, p.name)
	}
	// Show the last row with a person search box.
	grid.E("div class=attendanceNew").
		E("input class='s-search formInput' s-filter=type:Person placeholder='(add person)'")
	// Create a template for an empty row.
	tmpl := grid.E("template class=attendanceTemplate")
	row = tmpl.E("div class=attendanceRow")
	if t.Flags()&task.RecordHours != 0 {
		row.E("s-hours class=attendanceHours")
	}
	box = row.E("div class=attendanceSignedIn")
	box.E("s-icon icon=signature")
	box.E("input type=hidden")
	box = row.E("div class=attendanceCredited")
	box.E("s-icon")
	box.E("input type=hidden")
	row.E("div class=attendanceName")
	// Show the Cancel and Save buttons.
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit class='sbtn sbtn-primary' value=Save")
}

func postAttendance(r *request.Request, user *person.Person, e *event.Event, t *task.Task) {
	const personFields = person.FID | person.FInformalName
	r.Transaction(func() {
		for key := range r.Form {
			if !strings.HasPrefix(key, "signedin") {
				continue
			}
			var p = person.WithID(r, person.ID(util.ParseID(key[8:])), personFields)
			if p == nil {
				continue
			}
			minutes, _ := ui.SHoursValue(r.FormValue("hours" + key[8:]))
			var flags taskperson.Flag
			if r.FormValue("signedin"+key[8:]) == "true" {
				flags |= taskperson.Attended
			}
			if r.FormValue("credited"+key[8:]) == "true" {
				flags |= taskperson.Credited
			}
			taskperson.Set(r, e, t, p, minutes, flags)
		}
	})
	eventview.Render(r, user, e, fmt.Sprintf("task%d", t.ID()))
}
