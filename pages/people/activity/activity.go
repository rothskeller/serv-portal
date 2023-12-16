package activity

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/ui/orgdot"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const userFields = person.FID | person.FInformalName | person.FPrivLevels
const personFields = person.FID | person.FInformalName | person.FFlags

// HandleVolunteerHours handles /volunteer-hours/$token requests.
func HandleVolunteerHours(r *request.Request, token string) {
	var p = person.WithHoursToken(r, token, personFields)
	if p == nil {
		errpage.NotFound(r, nil)
		return
	}
	cy, cm := currentPeriod()
	handleCommon(r, p, p, cy, cm, cy, cm, nil)
}

// HandleActivity handles /people/$personid/activity/$period requests.
func HandleActivity(r *request.Request, pidstr, period string) {
	var (
		user   *person.Person
		p      *person.Person
		tabs   []ui.PageTab
		ny, nm int
		cy, cm int
		y, m   int
	)
	// Validate the user and target person.
	if user = auth.SessionUser(r, userFields, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(pidstr)), personFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() != p.ID() && !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	// Determine the period to be edited/viewed.
	cy, cm = currentPeriod()
	if period == "current" {
		y, m = cy, cm
	} else if y, m = parsePeriod(period); y == 0 {
		errpage.NotFound(r, user)
		return
	}
	// If the period is a year greater than the current, switch to month
	// view of the current month.
	if m == 0 && y > cy {
		http.Redirect(r, r.Request, fmt.Sprintf("/people/%d/activity/current", p.ID()), http.StatusSeeOther)
		return
	}
	// If the period is a month greater than now, switch to now.
	ny, nm = time.Now().Year(), int(time.Now().Month())
	if y > ny || (y == ny && m > nm) {
		http.Redirect(r, r.Request, fmt.Sprintf("/people/%d/activity/%d-%02d", p.ID(), ny, nm), http.StatusSeeOther)
		return
	}
	// If the period is a month prior to the current one, switch to year
	// view of the year containing it.  Special exception for webmaster, who
	// can add ?edit=true to edit past months.
	if m != 0 && (y < cy || (y == cy && m < cm)) && (r.FormValue("edit") == "" || !user.IsWebmaster()) {
		http.Redirect(r, r.Request, fmt.Sprintf("/people/%d/activity/%d", p.ID(), y), http.StatusSeeOther)
		return
	}
	tabs = []ui.PageTab{
		{Name: r.LangString("List", "Lista"), URL: "/people", Target: ".pageCanvas"},
		{Name: r.LangString("Map", "Mapa"), URL: "/people/map", Target: ".pageCanvas"},
		{Name: r.LangString("Details", "Detalles"), URL: fmt.Sprintf("/people/%d", p.ID()), Target: "main"},
		{Name: r.LangString("Activity", "Actividad"), URL: fmt.Sprintf("/people/%d/activity/%s", p.ID(), period), Target: "main", Active: true},
	}
	handleCommon(r, user, p, cy, cm, y, m, tabs)
}

func currentPeriod() (year, month int) {
	now := time.Now()
	if now.Day() > 10 {
		return now.Year(), int(now.Month())
	} else {
		lastmonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local)
		return lastmonth.Year(), int(lastmonth.Month())
	}
}

func parsePeriod(period string) (year, month int) {
	var err error

	if year, err = strconv.Atoi(period); err == nil && year >= 2000 && year <= 2099 {
		return year, 0
	}
	if date, err := time.ParseInLocation("2006-01", period, time.Local); err == nil && period == date.Format("2006-01") && date.Year() >= 2000 && date.Year() <= 2099 {
		return date.Year(), int(date.Month())
	}
	return 0, 0
}

func handleCommon(r *request.Request, user, p *person.Person, cy, cm, y, m int, tabs []ui.PageTab) {
	var opts ui.PageOpts

	// Just visiting the page is enough to clear the hours reminder.
	if user.ID() == p.ID() && p.Flags()&person.HoursReminder != 0 && y == cy && m == cm {
		var up = p.Updater()
		up.Flags &^= person.HoursReminder
		r.Transaction(func() {
			p.Update(r, up, person.FFlags)
		})
	}
	// Save any changed data.
	if r.Method == http.MethodPost && m != 0 {
		saveHours(r, user, p, y, m)
	}
	// Set the page options.
	opts.Title = r.LangString("Activity", "Actividad")
	if user.ID() == p.ID() {
		opts.Banner = r.LangString("Volunteer Activity", "Actividad de voluntariado")
	} else {
		opts.Banner = r.LangString(p.InformalName()+" Activity", "Actividad de "+p.InformalName())
	}
	if user.ID() == p.ID() {
		opts.MenuItem = "profile"
	} else {
		opts.MenuItem = "people"
	}
	opts.Tabs = tabs // which may be nil
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		if m == 0 {
			showYearView(r, main, p, y)
		} else {
			showMonthView(r, main, user, p, y, m)
		}
	})
}

func saveHours(r *request.Request, user, p *person.Person, y, m int) {
	const eventFields = event.FID | event.FName | event.FStart | event.FFlags
	const taskFields = task.FID | task.FName | task.FEvent | task.FFlags
	type change struct {
		e       *event.Event
		t       *task.Task
		minutes uint
		flags   taskperson.Flag
	}
	var changes []change

	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02T00:00")
	event.AllBetween(r, fmt.Sprintf("%d-%02d-01", y, m), fmt.Sprintf("%d-%02d-32", y, m), eventFields, 0, func(e *event.Event, _ *venue.Venue) {
		if e.Start() >= tomorrow && e.Flags()&event.OtherHours == 0 {
			return
		}
		task.AllForEvent(r, e.ID(), taskFields, func(t *task.Task) {
			if t.Flags()&task.RecordHours == 0 {
				return
			}
			if user.ID() != p.ID() && !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
				return
			}
			want, ok := ui.SHoursValue(r.FormValue(fmt.Sprintf("t%d", t.ID())))
			if !ok {
				return
			}
			have, flags := taskperson.Get(r, t.ID(), p.ID())
			if want == have {
				return
			}
			changes = append(changes, change{e.Clone(), t.Clone(), want, flags})
		})
	})
	if len(changes) != 0 {
		r.Transaction(func() {
			for _, change := range changes {
				taskperson.Set(r, change.e, change.t, p, change.minutes, change.flags)
			}
		})
	}
}

func showYearView(r *request.Request, main *htmlb.Element, p *person.Person, y int) {
	const eventFields = event.FID | event.FName | event.FStart | event.FFlags
	const taskFields = task.FID | task.FName | task.FOrg
	var grid *htmlb.Element

	main.E("s-year class=activityYear value=%d", y)
	taskperson.AllBetween(r, fmt.Sprintf("%d-01-01", y), fmt.Sprintf("%d-01-01", y+1), p.ID(), eventFields, taskFields,
		func(e *event.Event, t *task.Task, minutes uint, flags taskperson.Flag) {
			if grid == nil {
				grid = main.E("div class=activityYearGrid")
			}
			if minutes != 0 {
				var half bool

				fh := ui.MinutesToHours(minutes)
				if strings.HasSuffix(fh, "½") {
					half = true
					fh = fh[:len(fh)-2]
				}
				if fh == "" {
					fh = "0"
				}
				grid.E("div class=activityYearHours>%s", fh)
				if half {
					grid.E("div class=activityYearHalf>½")
				}
			}
			if flags&taskperson.Attended != 0 {
				grid.E("s-icon class=activityYearAttended icon=signature title=%s", r.LangString("Signed In", "Registrado"))
			}
			if flags&taskperson.Credited != 0 {
				grid.E("s-icon class=activityYearCredited icon=star-solid title=%s", r.LangString("Credited", "Acreditado"))
			}
			if e.Flags()&event.OtherHours != 0 {
				grid.E("div class=activityYearDate>%s", e.Start()[:7])
			} else {
				grid.E("div class=activityYearDate>%s", e.Start()[:10])
			}
			label := grid.E("div class=activityYearLabel")
			orgdot.OrgDot(r, label, t.Org())
			if e.Flags()&event.OtherHours != 0 {
				label.TF(r.LangString(" Other %s Hours", " Otras horas para %s"), t.Name())
			} else {
				label.TF(" %s", e.Name())
				if t.Name() != "Tracking" {
					label.E("span class=activityYearTaskName>%s", t.Name())
				}
			}
		})
	if grid == nil {
		main.E("div").R(r.LangString("No activity.", "No hay actividad."))
	}
}

func showMonthView(r *request.Request, main *htmlb.Element, user, p *person.Person, y, m int) {
	const eventFields = event.FID | event.FName | event.FStart | event.FFlags
	const taskFields = task.FID | task.FName | task.FOrg | task.FFlags

	form := main.E("form class=activity method=POST up-target=.activity")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("s-month value=%d-%02d", y, m)
	grid := form.E("div class=activityGrid")
	tomorrow := time.Now().AddDate(0, 0, 1).Format("2006-01-02T00:00")
	event.AllBetween(r, fmt.Sprintf("%d-%02d-01", y, m), fmt.Sprintf("%d-%02d-32", y, m), eventFields, 0, func(e *event.Event, _ *venue.Venue) {
		if e.Start() >= tomorrow && e.Flags()&event.OtherHours == 0 {
			return
		}
		task.AllForEvent(r, e.ID(), taskFields, func(t *task.Task) {
			if t.Flags()&task.RecordHours == 0 {
				return
			}
			editable := user.ID() == p.ID() || user.HasPrivLevel(t.Org(), enum.PrivLeader)
			minutes, flags := taskperson.Get(r, t.ID(), p.ID())
			grid.E("s-hours class=activityHours name=t%d value=%s", t.ID(), ui.MinutesToHours(minutes), !editable, "disabled")
			if flags&taskperson.Attended != 0 {
				grid.E("s-icon class=activityAttended icon=signature title=%s", r.LangString("Signed In", "Registrado"))
			}
			if flags&taskperson.Credited != 0 {
				grid.E("s-icon class=activityCredited icon=star-solid title=%s", r.LangString("Credited", "Acreditado"))
			}
			label := grid.E("div class=activityLabel")
			orgdot.OrgDot(r, label, t.Org())
			if e.Flags()&event.OtherHours != 0 {
				label.TF(r.LangString(" %s Other %s Hours", " %s Otras horas para %s"), e.Start()[:7], t.Name())
			} else {
				label.TF(" %s %s", e.Start()[:10], e.Name())
				if t.Name() != "Tracking" {
					label.E("div class=activityTaskName>%s", t.Name())
				}
			}
		})
	})
	form.E("div class=activityButtons hidden").
		E("input type=submit class='sbtn sbtn-warning' value=%s", r.LangString("Save", "Guardar"))
}
