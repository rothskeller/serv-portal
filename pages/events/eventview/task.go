package eventview

import (
	"fmt"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/events/signups"
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
	"sunnyvaleserv.org/portal/ui/orgdot"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const taskEventFields = event.FID | event.FStart | event.FDetails | taskperson.SetEventFields
const taskTaskFields = task.FID | task.FName | task.FOrg | task.FFlags | task.FDetails
const taskPersonFields = shiftperson.EligibilityCheckerPersonFields | taskperson.SetPersonFields

func showTask(r *request.Request, main *htmlb.Element, user *person.Person, e *event.Event, t *task.Task, multiple bool) {
	const shiftFields = shift.FID | shift.FStart | shift.FEnd | shift.FVenue | shift.FMin | shift.FMax
	var (
		roles    []string
		shifts   []*shift.Shift
		signedUp []bool
	)
	// Determine the viewer's involvement in the task.
	minutes, flags := taskperson.Get(r, t.ID(), user.ID())
	var hasrole, signedUpAny bool
	taskrole.Get(r, t.ID(), role.FID|role.FName, func(rl *role.Role) {
		if held, _ := personrole.PersonHasRole(r, user.ID(), rl.ID()); held {
			hasrole = true
			roles = append(roles, rl.Name())
		}
	})
	shift.AllForTask(r, t.ID(), shiftFields, 0, func(s *shift.Shift, _ *venue.Venue) {
		sclone := *s
		shifts = append(shifts, &sclone)
		signedUp = append(signedUp, shiftperson.Get(r, s.ID(), user.ID()) > 0)
		if signedUp[len(signedUp)-1] {
			signedUpAny = true
		}
	})
	var editable = user.HasPrivLevel(t.Org(), enum.PrivLeader)
	// If the viewer isn't involved in any way, don't show the task.
	if !editable && !hasrole && !signedUpAny && minutes == 0 && flags == 0 {
		return
	}
	// Display the task header.
	section := main.E("div id=eventviewTask%d class=eventviewSection", t.ID())
	sheader := section.E("div class=eventviewSectionHeader")
	title := sheader.E("div class=eventviewSectionHeaderText>%s", t.Name())
	orgdot.OrgDot(r, title.E("span class=eventviewTaskOrg"), t.Org())
	if t.Flags()&task.CoveredByDSW != 0 {
		title.E("span class=eventviewTaskDSW>DSW")
	}
	if editable {
		sheader.E("div class=eventviewSectionHeaderEdit").
			E("a href=/events/edtask/%d up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Edit", t.ID())
	}
	bdiv := section.E("div class=eventviewTask")
	if t.Details() != "" {
		bdiv.E("div class=eventviewTaskDetails").R(t.Details())
	}
	// Display signups if there are any shifts.
	if len(shifts) != 0 {
		showTaskSignups(r, bdiv, user, t, shifts, signedUp, roles, editable, hasrole, signedUpAny)
	}
	// Display tracking if the user signed in, got credit, can edit, has an
	// associated role and anyone got credit, or has an associated role in a
	// non-signup task and anyone signed in.  Also display if volunteer
	// hours recording is enabled and the user has hours recorded or has an
	// associated role.
	attended, credited := flags&taskperson.Attended != 0, flags&taskperson.Credited != 0
	anyAttended, anyCredited := t.Flags()&task.HasAttended != 0, t.Flags()&task.HasCredited != 0
	hoursTracked := t.Flags()&task.RecordHours != 0
	date, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
	hoursCutoff := time.Date(date.Year(), date.Month()+1, 11, 0, 0, 0, 0, time.Local)
	canRecordHours := time.Now().Before(hoursCutoff) || user.IsWebmaster()
	if editable || attended || credited || (hasrole && len(shifts) == 0 && anyAttended) || (hasrole && anyCredited) ||
		(hoursTracked && (minutes != 0 || hasrole)) {
		showTaskTracking(r, bdiv, t, editable, len(shifts) != 0, hasrole, attended, credited, anyAttended, anyCredited, hoursTracked, canRecordHours, minutes)
	}
}

// showTaskSignups shows the signups for shifts.
func showTaskSignups(r *request.Request, body *htmlb.Element, user *person.Person, t *task.Task, shifts []*shift.Shift, signedUp []bool, roles []string, editable, hasrole, signedUpAny bool) {
	form := body.E("form method=POST up-target=.eventview")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("input type=hidden name=shift")
	form.E("input type=hidden name=signedup")
	form.E("input type=hidden name=person")
	heading := form.E("div class=eventviewTaskHeading").R(r.Loc("Signups"))
	if editable {
		heading.E("a href=/events/proxysignup/%d up-layer=new up-size=grow up-dismissable=false up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Proxy Signup", t.ID())
	}
	if !signedUpAny && !hasrole && !editable {
		var conjoin string

		switch len(roles) {
		case 0:
			form.E("div").R(r.Loc("No one can sign up right now."))
			return
		case 1:
			conjoin = roles[0]
		case 2:
			conjoin = roles[0] + " " + r.Loc("and") + " " + roles[1]
		default:
			conjoin = strings.Join(roles[:len(roles)-1], ", ") + ", " + r.Loc("and") + " " + roles[len(roles)-1]
		}
		form.E("div").TF(r.Loc("Only %s can sign up."), conjoin)
		return
	}
	if t.Flags()&task.RequiresBGCheck != 0 && !editable && (!user.BGChecks().DOJ.Valid() || !user.BGChecks().FBI.Valid()) {
		form.E("div").R(r.Loc("Signups for this task require a completed background check."))
		return
	}
	if t.Flags()&task.CoveredByDSW != 0 && !editable {
		if reg, _ := user.DSWRegistrationForOrg(t.Org()); !reg.Valid() {
			form.E("div").R(r.Loc("Signups for this task require current DSW registration."))
			return
		}
	}
	signups.ShowTaskSignups(r, form, t, user, editable, true)
	if editable {
		buttons := form.E("div class=eventviewTaskSignupsEdit")
		buttons.E("a href=/events/edshift/NEW?tid=%d up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Add Shift", t.ID())
	}
}

func handleSignup(r *request.Request, user *person.Person, e *event.Event) {
	signups.HandleShiftSignup(r, user, user)
}

// showTaskTracking displays the tracking information for the Task.
func showTaskTracking(r *request.Request, body *htmlb.Element, t *task.Task, editable, hasshifts, hasrole, attended, credited, anyAttended, anyCredited, hoursTracked, canRecordHours bool, minutes uint) {
	heading := body.E("div class=eventviewTaskHeading").R(r.Loc("Attendance"))
	if editable {
		heading.E("a href=/events/attendance/%d up-layer=new up-size=grow up-dismissable=false up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Record Attendance", t.ID())
	}
	box := body.E("form class=eventviewTaskAttendance method=POST up-target=.eventviewTaskAttendance")
	box.E("input type=hidden name=csrf value=%s", r.CSRF)
	if attended || (anyAttended && hasrole && !hasshifts) {
		box.E("div class=eventviewTaskStatusIcon", attended, "class=attended", !attended, "class=false").
			E("s-icon icon=signature")
		if attended {
			box.E("div class=eventviewTaskStatusText").R(r.Loc("You signed in."))
		} else {
			box.E("div class=eventviewTaskStatusText").R(r.Loc("You did not sign in."))
		}
	}
	if credited || (anyCredited && hasrole) {
		if credited {
			box.E("div class='eventviewTaskStatusIcon credited'").E("s-icon icon=star-solid")
			box.E("div class=eventviewTaskStatusText").R(r.Loc("You were credited for this session."))
		} else {
			box.E("div class='eventviewTaskStatusIcon false'").E("s-icon icon=star")
			box.E("div class=eventviewTaskStatusText").R(r.Loc("You were not credited for this session."))
		}
	}
	if hoursTracked && (minutes != 0 || hasrole) {
		box.E("div class=eventviewTaskStatusIcon", minutes != 0, "class=minutes", minutes == 0, "class=false").
			E("s-icon icon=clock")
		if canRecordHours {
			line := box.E("div class=eventviewTaskHours").R(r.Loc("Volunteer hours") + ": ")
			line.E("s-hours name=hours value=%s", ui.MinutesToHours(minutes))
			line.E("input type=submit name=edhours%d class='sbtn sbtn-warning eventviewTaskHoursSave' hidden value=%s>", t.ID(), r.Loc("Save"))
		} else if minutes != 0 {
			var m = ui.MinutesToHours(minutes)
			if m != "½" && m != "1" {
				box.E("div").TF(r.Loc("You spent %s volunteer hours."), m)
			} else {
				box.E("div").TF(r.Loc("You spent %s volunteer hour."), m)
			}
		} else {
			box.E("div").R(r.Loc("You did not record volunteer hours."))
		}
	}
}

func handleHours(r *request.Request, user *person.Person, e *event.Event) {
	const taskFields = task.FID | task.FFlags | taskperson.SetTaskFields
	date, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
	hoursCutoff := time.Date(date.Year(), date.Month()+1, 11, 0, 0, 0, 0, time.Local)
	if !time.Now().Before(hoursCutoff) && !user.IsWebmaster() {
		return // past the hours recording period for this event
	}
	task.AllForEvent(r, e.ID(), taskFields, func(t *task.Task) {
		if t.Flags()&task.RecordHours == 0 {
			return // not recording hours for this task
		}
		if r.FormValue(fmt.Sprintf("edhours%d", t.ID())) == "" {
			return // didn't hit Save on hours for this task
		}
		wantmin, ok := ui.SHoursValue(r.FormValue("hours"))
		if !ok {
			return // input doesn't contain a valid value
		}
		havemin, haveflags := taskperson.Get(r, t.ID(), user.ID())
		if wantmin == havemin {
			return // no change
		}
		if havemin == 0 {
			if !taskperson.HasRoleForTask(r, t.ID(), user.ID()) {
				return // not eligible to record hours
			}
		}
		r.Transaction(func() {
			taskperson.Set(r, e, t, user, wantmin, haveflags)
		})
	})
}
