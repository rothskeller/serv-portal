package eventedit

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
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type shiftEditor struct {
	user       *person.Person
	e          *event.Event
	t          *task.Task
	s          *shift.Shift
	us         *shift.Updater
	hasSignups bool
	timesError string
	limitError string
	hasError   bool
	op         string
	validate   request.ValidationList
}

// HandleShift handles requests for /events/edshift/$id.  $id may be "NEW", in
// which case a tid=$tid parameter must be part of the request.
func HandleShift(r *request.Request, sidstr string) {
	var se *shiftEditor

	if se = getShiftEditor(r, sidstr); se == nil {
		return
	}
	if se.op == "delete" {
		se.handleDelete(r)
		return
	}
	if se.op != "get" {
		readShiftVenue(r, se.us)
		se.timesError = readShiftTimes(r, se.us)
		se.limitError = readShiftLimits(r, se.us)
		se.hasError = se.timesError != "" || se.limitError != ""
	}
	if !se.hasError && (se.op == "save" || (se.op == "copy" && se.s == nil)) {
		r.Transaction(func() {
			if se.s == nil {
				se.s = shift.Create(r, se.us)
			} else {
				se.s.Update(r, se.us)
			}
		})
		if se.op == "save" {
			eventview.Render(r, se.user, se.e, fmt.Sprintf("task%d", se.t.ID()))
			return
		}
	}
	if se.op == "copy" && se.s != nil && se.us.Start == se.s.Start() && se.us.End == se.s.End() && se.us.Venue.ID() == se.s.Venue() {
		// We're copying something that has not been modified.  Shift
		// the times forward to try to prevent overlap.
		start, _ := time.ParseInLocation("2006-01-02T15:04", se.us.Start, time.Local)
		end, _ := time.ParseInLocation("2006-01-02T15:04", se.us.End, time.Local)
		newend := end.Add(end.Sub(start))
		if newend.Day() != start.Day() {
			newend = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 0, 0, time.Local)
		}
		se.us.Start = se.us.End
		se.us.End = newend.Format("2006-01-02T15:04")
	}
	if se.op == "copy" {
		se.s = nil
		se.us.ID = 0
	}
	r.HTMLNoCache()
	if se.op != "get" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	se.writeForm(r, html)
}

func getShiftEditor(r *request.Request, sidstr string) (se *shiftEditor) {
	const eventFields = event.FID | event.FName | event.FStart | event.FFlags
	const taskFields = task.FID | task.FEvent | task.FName | task.FOrg | task.FFlags
	var tid task.ID

	// Get a valid user.
	se = new(shiftEditor)
	if se.user = auth.SessionUser(r, 0, true); se.user == nil {
		return nil
	}
	if !auth.CheckCSRF(r, se.user) {
		return nil
	}
	// Get and check the event, task, and shift.
	if sidstr == "NEW" {
		tid = task.ID(util.ParseID(r.FormValue("tid")))
	} else {
		if se.s = shift.WithID(r, shift.ID(util.ParseID(sidstr)), shift.UpdaterFields); se.s == nil {
			errpage.NotFound(r, se.user)
			return nil
		}
		tid = se.s.Task()
	}
	if se.t = task.WithID(r, tid, taskFields); se.t == nil {
		errpage.NotFound(r, se.user)
		return nil
	}
	if se.e = event.WithID(r, se.t.Event(), eventFields); se.e == nil {
		errpage.NotFound(r, se.user)
		return nil
	}
	if !se.user.HasPrivLevel(se.t.Org(), enum.PrivLeader) || se.e.Flags()&event.OtherHours != 0 {
		errpage.Forbidden(r, se.user)
		return nil
	}
	se.hasSignups = shiftperson.HasSignups(r, se.s.ID())
	// Get an updater.
	if se.s == nil {
		se.us = &shift.Updater{Event: se.e, Task: se.t}
	} else {
		se.us = se.s.Updater(r, se.e, se.t, nil)
	}
	// Get the operation.
	if r.Method != http.MethodPost {
		se.op = "get"
	} else if se.validate = r.ValidationList(); se.validate.Enabled() {
		se.op = "validate"
	} else if r.FormValue("delete") != "" && se.s != nil && !se.hasSignups {
		se.op = "delete"
	} else if r.FormValue("copy") != "" {
		se.op = "copy"
	} else {
		se.op = "save"
	}
	return se
}

func (se *shiftEditor) writeForm(r *request.Request, html *htmlb.Element) {
	form := html.E("form class='form form-2col eventeditShiftForm' method=POST up-main up-layer=parent up-target=#eventviewTask%d", se.t.ID())
	if se.s == nil {
		form.Attr("action=/events/edshift/NEW")
		form.E("input type=hidden name=tid value=%d", se.t.ID())
		form.E("div class='formTitle formTitle-primary'>New Shift")
	} else {
		form.Attr("action=/events/edshift/%d", se.s.ID())
		form.E("div class='formTitle formTitle-primary'>Edit Shift")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if se.validate.ValidatingAny("start", "end", "venue") {
		emitShiftTimes(form, se.us, se.timesError != "" || !se.hasError, se.timesError)
	}
	if se.validate.Validating("venue") {
		emitShiftVenue(form, se.us)
	}
	if se.validate.ValidatingAny("min", "max") {
		emitShiftLimits(form, se.us, se.limitError != "", se.limitError)
	}
	if !se.validate.Enabled() {
		emitShiftButtons(form, se.us, se.s != nil && !se.hasSignups)
	}
}

func readShiftTimes(r *request.Request, us *shift.Updater) string {
	var start, end = r.FormValue("start"), r.FormValue("end")
	us.Start = us.Event.Start()[:10] + "T" + start
	us.End = us.Event.Start()[:10] + "T" + end
	if start == "" || end == "" {
		return "The start and end times are required."
	}
	if t, err := time.Parse("2006-01-02T15:04", us.Start); err != nil || us.Start != t.Format("2006-01-02T15:04") {
		return "The start time is not a valid time."
	}
	if t, err := time.Parse("2006-01-02T15:04", us.End); err != nil || us.End != t.Format("2006-01-02T15:04") {
		return "The end time is not a valid time."
	}
	if us.End < us.Start {
		return "The end time must not be before the start time."
	}
	if us.Venue != nil && us.Venue.Flags()&venue.CanOverlap == 0 && us.OverlappingShift(r) {
		return "Another shift is happening at the same place and an overlapping time."
	}
	return ""
}
func emitShiftTimes(form *htmlb.Element, us *shift.Updater, focus bool, err string) {
	var start, end string
	if parts := strings.Split(us.Start, "T"); len(parts) > 1 {
		start = parts[1]
	}
	if parts := strings.Split(us.End, "T"); len(parts) > 1 {
		end = parts[1]
	}
	row := form.E("div id=eventeditShiftTimesRow class=formRow")
	row.E("label for=eventeditShiftStart>Time")
	box := row.E("div class='formInput eventeditTimes'")
	box.E("input type=time id=eventeditShiftStart name=start class=formInput s-validate value=%s", start, focus, "autofocus")
	box.R("to")
	box.E("input type=time name=end class=formInput s-validate value=%s", end)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readShiftVenue(r *request.Request, us *shift.Updater) string {
	us.Venue = nil
	var vkey = r.FormValue("venue")
	if vkey == "" {
		return ""
	}
	if strings.HasPrefix(vkey, "V") {
		us.Venue = venue.WithID(r, venue.ID(util.ParseID(vkey[1:])), venue.FID|venue.FName|venue.FFlags)
	}
	if us.Venue == nil {
		return "The venue is not recognized."
	}
	return ""
}
func emitShiftVenue(form *htmlb.Element, us *shift.Updater) {
	var vkey, vname string

	if us.Venue != nil {
		vkey = us.Venue.IndexKey(nil)
		vname = us.Venue.Name()
	}
	row := form.E("div id=eventeditShiftVenueRow class=formRow")
	row.E("label for=eventeditShiftVenue>Venue")
	row.E("s-searchcombo id=eventeditShiftVenue name=venue class=formInput value=%s valuelabel=%s facet=type:Venue placeholder=TBD", vkey, vname)
}

func readShiftLimits(r *request.Request, us *shift.Updater) string {
	var minstr = r.FormValue("min")
	if minstr == "" {
		us.Min = 0
	} else if min, err := strconv.Atoi(minstr); err != nil || min < 0 {
		us.Min = 0
		return "The minimum value is not valid."
	} else {
		us.Min = uint(min)
	}
	var maxstr = r.FormValue("max")
	if maxstr == "" {
		us.Max = 0
	} else if max, err := strconv.Atoi(maxstr); err != nil || max < 0 {
		us.Max = 0
		return "The maximum value is not valid."
	} else {
		us.Max = uint(max)
	}
	if us.Max != 0 && us.Max < us.Min {
		return "The maximum value cannot be lower than the minimum value."
	}
	return ""
}
func emitShiftLimits(form *htmlb.Element, us *shift.Updater, focus bool, err string) {
	var minstr, maxstr string

	if us.Min != 0 {
		minstr = strconv.Itoa(int(us.Min))
	}
	if us.Max != 0 {
		maxstr = strconv.Itoa(int(us.Max))
	}
	row := form.E("div id=eventeditShiftCapacityRow class='formRow eventeditShiftCapacityRow'")
	row.E("label for=eventeditShiftMin>Capacity")
	in := row.E("div class='formInput eventeditShiftCapacity'")
	in.E("span>Need")
	in.E("input type=number id=eventeditShiftMin name=min class=formInput s-validate min=0 value=%s", minstr)
	in.E("span>Limit")
	in.E("input type=number id=eventeditShiftMax name=max class=formInput s-validate min=0 value=%s", maxstr)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func emitShiftButtons(form *htmlb.Element, us *shift.Updater, canDelete bool) {
	buttons := form.E("div class=formButtons")
	buttons.E("div class=formButtonSpace")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	// These buttons must appear lexically after Save, even though they
	// appear visually before it, so that Save is the default button when
	// the user presses Enter.  The formButton-beforeAll class implements
	// that.
	if canDelete {
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
	buttons.E("input type=submit name=copy class='sbtn sbtn-secondary formButton-beforeAll' value=Copy")
}

func (se *shiftEditor) handleDelete(r *request.Request) {
	r.Transaction(func() {
		se.s.Delete(r, se.e, se.t)
		if !shift.ExistsForTask(r, se.t.ID()) && se.t.Flags()&task.SignupsOpen != 0 {
			se.t = task.WithID(r, se.s.Task(), task.UpdaterFields)
			var ut = se.t.Updater(r, se.e)
			ut.Flags &^= task.SignupsOpen
			se.t.Update(r, ut)
		}
	})
	eventview.Render(r, se.user, se.e, fmt.Sprintf("task%d", se.t.ID()))
}
