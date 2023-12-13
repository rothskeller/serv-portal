package eventcopy

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskrole"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

/* EVENT COPY DIALOG

This dialog creates one or more copies of an event.  It has the following UI.

Repeat every:  [COUNT] [DAY|WEEK|MONTH]
Repeat on:     [REPEATON]
Stop on:       [STOPDATE]
                      [Cancel] [[Copy]]

[COUNT] is a positive number, defaulting to 1.
[DAY|WEEK|MONTH] is a dropdown with "day", "week", and "month" options,
defaulting to "week".  When [COUNT] is something other than 0, the options are
pluralized.

[REPEATON] varies depending on the [DAY|WEEK|MONTH] choice.  When it is set to
"day", the whole "Repeat on" row does not appear.  When it is set to "week",
[REPEATON] contains checkboxes for the seven days of the week, with only the one
representing the DOW of the source event initially selected.  When it is set to
"month", [REPEATON] is a set of radio buttons with two or three choices:
  - "the ##th of the month"
  - "the xxxrd XXXday of the month"
  - "the last XXXday of the month"
The phrasing of these reflects the date of the source event.  The last one
appears only if the source event is on the last XXXday of its month.

[STOPDATE] is the date of the last copy.  It is updated whenever [COUNT] or
[DAY|WEEK|MONTH] is changed, such that it would result in a single copy.

Copying an event copies all of its tasks and shifts.  Everything gets new IDs,
of course, and the dates change, but nothing else.  The set of people signed up
for, or declining, shifts is not carried over; neither are attendance records.
The whole operation will fail and no copies will be created if any copy would be
invalid, which basically can only happen on an event name conflict.

On success, the event edit page for the final copy is displayed.
*/

type copyData struct {
	e           *event.Event
	ts          []*task.Task
	roles       [][]*role.Role
	ss          [][]*shift.Shift
	vs          [][]*venue.Venue
	everyCount  int
	everyType   int // 1, 7, or 31
	repeatOn    int // for 7: bitmask of weekdays; for 31: 0=day, week number, or 5=last
	weekday     time.Weekday
	weeknum     int
	lastweek    bool
	stopOn      string
	everyError  string
	repeatError string
	stopError   string
}

// Handle handles /events/$eid/copy requests.
func Handle(r *request.Request, idstr string) {
	const eventFields = event.UpdaterFields
	const taskFields = task.UpdaterFields
	var (
		user    *person.Person
		allowed bool
		cd      copyData
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if cd.e = event.WithID(r, event.ID(util.ParseID(idstr)), eventFields); cd.e == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(0, enum.PrivLeader) || cd.e.Flags()&event.OtherHours != 0 {
		errpage.Forbidden(r, user)
		return
	}
	allowed = true
	task.AllForEvent(r, cd.e.ID(), taskFields, func(t *task.Task) {
		cd.ts = append(cd.ts, t.Clone())
		if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
			allowed = false
		}
	})
	if !allowed {
		errpage.Forbidden(r, user)
		return
	}
	cd.setDefaults()
	if cd.handlePost(r, user) {
		return
	}
	r.HTMLNoCache()
	if r.Method == http.MethodPost {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	cd.writeForm(r)
}

// Set defaults sets the default repeat pattern: a single copy, one week later
// than the source event.
func (cd *copyData) setDefaults() {
	cd.everyCount = 1
	cd.everyType = 7
	date, _ := time.Parse("2006-01-02", cd.e.Start()[:10])
	nextweek := date.AddDate(0, 0, 7)
	cd.repeatOn = 1 << date.Weekday()
	cd.stopOn = nextweek.Format("2006-01-02")
	cd.weekday = date.Weekday()
	cd.weeknum = (date.Day()-1)/7 + 1
	cd.lastweek = nextweek.Month() != date.Month()
}

func (cd *copyData) writeForm(r *request.Request) {
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	form.E("div class='formTitle formTitle-primary'>Copy Event")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	row := form.E("div class=formRow")
	row.E("label for=eventcopyEveryCount>Repeat every")
	box := row.E("div class='formInput eventcopyEvery'")
	box.E("input type=number id=eventcopyEveryCount name=everyCount class=formInput min=1 value=%d up-validate", cd.everyCount)
	sel := box.E("select name=everyType class=formInput up-validate")
	if cd.everyCount == 1 {
		sel.E("option value=1", cd.everyType == 1, "selected").R("day")
		sel.E("option value=7", cd.everyType == 7, "selected").R("week")
		sel.E("option value=31", cd.everyType == 31, "selected").R("month")
	} else {
		sel.E("option value=1", cd.everyType == 1, "selected").R("days")
		sel.E("option value=7", cd.everyType == 7, "selected").R("weeks")
		sel.E("option value=31", cd.everyType == 31, "selected").R("months")
	}
	if cd.everyError != "" {
		row.E("div class=formError>%s", cd.everyError)
	}
	row = form.E("div class=formRow")
	switch cd.everyType {
	case 7:
		row.E("label for=eventcopyRepeatSunday>Repeat on")
		box = row.E("div class='formInput eventcopyRepeat'")
		box.E("div>S")
		box.E("div>M")
		box.E("div>T")
		box.E("div>W")
		box.E("div>T")
		box.E("div>F")
		box.E("div>S")
		box.E("input type=checkbox id=eventcopyRepeatSunday name=repeat class=s-check value=0 up-validate", cd.repeatOn&(1<<time.Sunday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=1 up-validate", cd.repeatOn&(1<<time.Monday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=2 up-validate", cd.repeatOn&(1<<time.Tuesday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=3 up-validate", cd.repeatOn&(1<<time.Wednesday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=4 up-validate", cd.repeatOn&(1<<time.Thursday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=5 up-validate", cd.repeatOn&(1<<time.Friday) != 0, "checked")
		box.E("input type=checkbox name=repeat class=s-check value=6 up-validate", cd.repeatOn&(1<<time.Saturday) != 0, "checked")
	case 31:
		row.E("label for=eventcopyRepeatNth>Repeat on")
		box = row.E("div class=formInput")
		box.E("s-radio id=eventcopyRepeatNth name=repeat value=0 label='On the %s of each month' up-validate",
			ordinalDate(cd.e.Start()[8:10]), cd.repeatOn == 0, "checked")
		if cd.weeknum != 5 {
			box.E("s-radio name=repeat value=%d label='On the %s %s of each month' up-validate",
				cd.weeknum, ordinalWeek[cd.weeknum], cd.weekday.String(), cd.repeatOn == cd.weeknum, "checked")
		}
		if cd.lastweek {
			box.E("s-radio name=repeat value=5 label='On the last %s of each month' up-validate",
				cd.weekday.String(), cd.repeatOn == 5, "checked")
		}
	}
	if cd.repeatError != "" {
		row.E("div class=formError>%s", cd.repeatError)
	}
	row = form.E("div class=formRow")
	row.E("label for=eventcopyStop>Stop on")
	row.E("input type=date id=eventcopyStop name=stop class=formInput value=%s", cd.stopOn)
	// up-validate doesn't work right on type=date.  If I add
	// up-watch-event=blur, then the submit button doesn't work right.  Best
	// to not validate it and allow the form submission to do so.
	if cd.stopError != "" {
		row.E("div class=formError>%s", cd.stopError)
	}
	box = row.E("div class=formButtons")
	box.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	box.E("input type=submit class='sbtn sbtn-primary' value=Copy")
}

func (cd *copyData) handlePost(r *request.Request, user *person.Person) bool {
	var (
		ues   []*event.Updater
		laste *event.Event
	)
	if r.Method != http.MethodPost {
		return false
	}
	cd.readForm(r)
	if r.Request.Header.Get("X-Up-Validate") != "" || cd.everyError != "" || cd.repeatError != "" || cd.stopError != "" {
		return false
	}
	if ues = cd.buildEventUpdaters(r); len(ues) == 0 {
		return false
	}
	cd.readDetails(r)
	r.Transaction(func() {
		for _, ue := range ues {
			laste = cd.doCopy(r, ue)
		}
	})
	r.Header().Set("X-Up-Location", fmt.Sprintf("/events/%d", laste.ID()))
	eventview.Render(r, user, laste, "")
	return true
}

func (cd *copyData) readForm(r *request.Request) {
	cd.everyCount, _ = strconv.Atoi(r.FormValue("everyCount"))
	cd.everyType, _ = strconv.Atoi(r.FormValue("everyType"))
	if cd.everyCount < 1 {
		cd.everyError = "Count must be a positive number."
	}
	switch cd.everyType {
	case 1:
		// nothing
	case 7:
		cd.repeatOn = 0
		for _, str := range r.Form["repeat"] {
			if val := util.ParseID(str); val >= 0 && val < 7 {
				cd.repeatOn |= 1 << val
			}
		}
		if cd.repeatOn == 0 {
			cd.repeatError = "At least one day must be selected."
		}
	case 31:
		cd.repeatOn = util.ParseID(r.FormValue("repeat"))
		switch cd.repeatOn {
		case 0:
			// nothing
		case 1, 2, 3, 4:
			cd.repeatOn = cd.weeknum
		case 5:
			if !cd.lastweek {
				cd.repeatOn = cd.weeknum
			}
		default:
			cd.repeatOn = 0
		}
	default:
		cd.everyType, cd.repeatOn = 7, 1<<cd.weekday
	}
	next := cd.increment(cd.e.Start()[:10])
	if r.Request.Header.Get("X-Up-Validate") != "" {
		cd.stopOn = next
	} else {
		cd.stopOn = r.FormValue("stop")
	}
	if stop, err := time.Parse("2006-01-02", cd.stopOn); err != nil || stop.Format("2006-01-02") != cd.stopOn {
		cd.stopError = "A valid stop date is required."
	} else if cd.stopOn < next {
		cd.stopError = "With this stop date, no copies will be created."
	}
}

func (cd *copyData) buildEventUpdaters(r *request.Request) (ues []*event.Updater) {
	var v = venue.WithID(r, cd.e.Venue(), venue.FID|venue.FName)
	var next = cd.increment(cd.e.Start()[:10])
	for next <= cd.stopOn {
		var ue = cd.e.Updater(r, v)
		ue.ID = 0 // change to create
		ue.Start = next + ue.Start[10:]
		ue.End = next + ue.End[10:]
		if ue.DuplicateName(r) {
			cd.stopError = fmt.Sprintf("A conflicting event with the same name exists on %s.", next)
			return nil
		}
		ues = append(ues, ue)
		next = cd.increment(next)
	}
	return ues
}

func (cd *copyData) readDetails(r *request.Request) {
	for _, t := range cd.ts {
		var roles []*role.Role
		taskrole.Get(r, t.ID(), role.FID|role.FName, func(rl *role.Role) {
			roles = append(roles, rl.Clone())
		})
		cd.roles = append(cd.roles, roles)
		var shifts []*shift.Shift
		var venues []*venue.Venue
		shift.AllForTask(r, t.ID(), shift.UpdaterFields, venue.FID|venue.FName, func(s *shift.Shift, v *venue.Venue) {
			shifts = append(shifts, s.Clone())
			venues = append(venues, v.Clone())
		})
		cd.ss = append(cd.ss, shifts)
		cd.vs = append(cd.vs, venues)
	}
}

func (cd *copyData) increment(date string) string {
	dt, _ := time.Parse("2006-01-02", date)
	switch cd.everyType {
	case 1:
		dt = dt.AddDate(0, 0, cd.everyCount)
	case 7:
		var repeat = cd.repeatOn
		if repeat == 0 { // could be validating invalid form
			repeat = 1 << cd.weekday
		}
		dt = dt.AddDate(0, 0, 1)
		for {
			if dt.Weekday() == time.Sunday && cd.everyCount > 1 {
				dt = dt.AddDate(0, 0, 7*cd.everyCount-7)
			}
			if repeat&(1<<int(dt.Weekday())) != 0 {
				break
			}
			dt = dt.AddDate(0, 0, 1)
		}
	case 31:
		if cd.repeatOn == 0 {
			dnum, _ := strconv.Atoi(cd.e.Start()[8:10])
			dt = time.Date(dt.Year(), dt.Month()+time.Month(cd.everyCount), dnum, 0, 0, 0, 0, time.Local)
			if dt.Day() != dnum { // wrapped into the next month
				dt = dt.AddDate(0, 0, -dnum) // use the last day of the month instead
			}
			break
		}
		if cd.repeatOn < 5 {
			dt = time.Date(dt.Year(), dt.Month()+time.Month(cd.everyCount), 7*cd.repeatOn-6, 0, 0, 0, 0, time.Local)
		} else {
			dt = time.Date(dt.Year(), dt.Month()+time.Month(cd.everyCount)+1, -6, 0, 0, 0, 0, time.Local)
		}
		for dt.Weekday() != cd.weekday {
			dt = dt.AddDate(0, 0, 1)
		}
	}
	return dt.Format("2006-01-02")
}

func (cd *copyData) doCopy(r *request.Request, ue *event.Updater) (e *event.Event) {
	e = event.Create(r, ue)
	for ti, t := range cd.ts {
		var ut = t.Updater(r, e)
		ut.ID = 0
		ut.Flags &^= task.HasAttended | task.HasCredited
		var nt = task.Create(r, ut)
		taskrole.Set(r, e, nt, cd.roles[ti], []*role.Role{})
		for si, s := range cd.ss[ti] {
			var us = s.Updater(r, e, nt, cd.vs[ti][si])
			us.ID = 0
			us.Start = ue.Start[:10] + us.Start[10:]
			us.End = ue.End[:10] + us.End[10:]
			shift.Create(r, us)
		}
	}
	return e
}

var ordinalWeek = map[int]string{
	1: "first", 2: "second", 3: "third", 4: "fourth",
}

func ordinalDate(date string) string {
	date = strings.TrimLeft(date, "0")
	switch date[len(date)-1] {
	case '1':
		date += "st"
	case '2':
		date += "nd"
	case '3':
		date += "rd"
	default:
		date += "th"
	}
	return date
}
