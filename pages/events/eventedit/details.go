package eventedit

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleDetails handles requests for /events/$id/eddetails.
func HandleDetails(r *request.Request, idstr string) {
	var (
		user       *person.Person
		allowed    bool
		e          *event.Event
		v          *venue.Venue
		ue         *event.Updater
		nameError  string
		dateError  string
		timesError string
		venueError string
		hasError   bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if e = event.WithID(r, event.ID(util.ParseID(idstr)), event.UpdaterFields); e == nil {
		errpage.NotFound(r, user)
		return
	}
	if allowed = user.HasPrivLevel(0, enum.PrivLeader); allowed {
		task.AllForEvent(r, e.ID(), task.FOrg, func(t *task.Task) {
			if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
				allowed = false
			}
		})
	}
	if !allowed || e.Flags()&event.OtherHours != 0 {
		errpage.Forbidden(r, user)
		return
	}
	v = venue.WithID(r, e.Venue(), venue.FID|venue.FName|venue.FURL)
	ue = e.Updater(r, v)
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		nameError = readEventName(r, ue)
		readActivation(r, ue)
		dateError = readDate(r, ue)
		timesError = readEventTimes(r, ue)
		venueError = readEventVenue(r, ue)
		hasError = nameError != "" || dateError != "" || timesError != "" || venueError != ""
		readEventDetails(r, ue)
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !hasError {
			r.Transaction(func() {
				needShiftUpdates := e.Start()[:10] != ue.Start[:10]
				e.Update(r, ue)
				if needShiftUpdates {
					updateShiftDates(r, e, ue)
				}
			})
			eventview.Render(r, user, e, "details")
			return
		}
	}
	r.HTMLNoCache()
	if hasError {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.eventviewIdent,.eventviewDetails")
	form.E("div class='formTitle formTitle-primary'>Edit Details")
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
	if len(validate) == 0 {
		emitEventDetails(form, ue)
		emitDetailsButtons(form)
	}
}

func readEventName(r *request.Request, ue *event.Updater) string {
	if ue.Name = strings.TrimSpace(r.FormValue("name")); ue.Name == "" {
		return "The event name is required."
	}
	// Temporary read of date, so that the duplicate check can run.
	if ue.Start = r.FormValue("date"); len(ue.Start) == 10 {
		ue.Start += "T00:00"
		if ue.DuplicateName(r) {
			return fmt.Sprintf("Another event on %s has the name %q.", ue.Start[:10], ue.Name)
		}
	}
	return ""
}
func emitEventName(form *htmlb.Element, ue *event.Updater, focus bool, err string) {
	row := form.E("div id=eventeditNameRow class=formRow")
	row.E("label for=eventeditName>Name")
	row.E("input id=eventeditName name=name s-validate value=%s", ue.Name, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readActivation(r *request.Request, ue *event.Updater) {
	ue.Activation = strings.ToUpper(strings.TrimSpace(r.FormValue("activation")))
}
func emitActivation(form *htmlb.Element, ue *event.Updater) {
	row := form.E("div id=eventeditActivationRow class=formRow")
	row.E("label for=eventeditActivation>Act. Number")
	row.E("input id=eventeditActivation name=activation s-validate value=%s", ue.Activation)
	row.E("div class=formHelp>Sunnyvale OES activation number for this event")
}

func readDate(r *request.Request, ue *event.Updater) string {
	if ue.Start = strings.TrimSpace(r.FormValue("date")); ue.Start == "" {
		ue.Start, ue.End = "0000-00-00", "0000-00-00"
		return "The event date is required."
	}
	ue.End = ue.Start // readTimes will fill in the times.
	if d, err := time.Parse("2006-01-02", ue.Start); err != nil || ue.Start != d.Format("2006-01-02") {
		return "The event date is not a valid date."
	}
	return ""
}
func emitDate(form *htmlb.Element, ue *event.Updater, focus bool, err string) {
	var date = strings.Split(ue.Start, "T")[0]
	if date == "0000-00-00" {
		date = ""
	}
	row := form.E("div id=eventeditDateRow class=formRow")
	row.E("label for=eventeditDate>Date")
	row.E("input type=date id=personeditDate name=date s-validate=#eventeditNameRow,#eventeditDateRow value=%s", date, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readEventTimes(r *request.Request, ue *event.Updater) string {
	var start, end = r.FormValue("start"), r.FormValue("end")
	if start == "" {
		start = "00:00"
	}
	if end == "" {
		end = start
	}
	ue.Start += "T" + start
	ue.End += "T" + end
	if t, err := time.Parse("2006-01-02T15:04", ue.Start); err != nil || ue.Start != t.Format("2006-01-02T15:04") {
		return "The start time is not a valid time."
	}
	if t, err := time.Parse("2006-01-02T15:04", ue.End); err != nil || ue.End != t.Format("2006-01-02T15:04") {
		return "The end time is not a valid time."
	}
	if ue.End < ue.Start {
		return "The end time must not be before the start time."
	}
	return ""
}
func emitEventTimes(form *htmlb.Element, ue *event.Updater, focus bool, err string) {
	var start, end string
	if parts := strings.Split(ue.Start, "T"); len(parts) > 1 {
		start = parts[1]
	}
	if parts := strings.Split(ue.End, "T"); len(parts) > 1 {
		end = parts[1]
	}
	if end == start {
		end = ""
	}
	if start == "00:00" && end == "" {
		start = ""
	}
	row := form.E("div id=eventeditTimesRow class=formRow")
	row.E("label for=eventeditStart>Time")
	box := row.E("div class='formInput eventeditTimes'")
	box.E("input type=time id=eventeditStart name=start class=formInput s-validate value=%s", start, focus, "autofocus")
	box.R("to")
	box.E("input type=time name=end class=formInput s-validate value=%s", end)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readEventVenue(r *request.Request, ue *event.Updater) string {
	ue.Venue, ue.VenueURL = nil, ""
	var vkey = r.FormValue("venue")
	if vkey == "" {
		return ""
	}
	if strings.HasPrefix(vkey, "V") {
		ue.Venue = venue.WithID(r, venue.ID(util.ParseID(vkey[1:])), venue.FID|venue.FName|venue.FURL)
	}
	if ue.Venue == nil {
		return "The venue is not recognized."
	}
	if ue.Venue.URL() == "" {
		if ue.VenueURL = r.FormValue("venueURL"); ue.VenueURL != "" {
			if uri, err := url.Parse(ue.VenueURL); err != nil || (uri.Scheme != "http" && uri.Scheme != "https") {
				return "The venue URL is not a valid http: or https: URL."
			}
		}
	}
	return ""
}
func emitEventVenue(form *htmlb.Element, ue *event.Updater, focus bool, err string) {
	var vkey, vname string

	if ue.Venue != nil {
		vkey = ue.Venue.IndexKey(nil)
		vname = ue.Venue.Name()
	}
	row := form.E("div id=eventeditVenueRow class=formRow")
	row.E("label for=eventeditVenue>Venue")
	row.E("s-searchcombo id=eventeditVenue name=venue class=formInput value=%s valuelabel=%s facet=type:Venue edit=Venue placeholder=TBD", vkey, vname)
	row = form.E("div id=eventeditVenueURLRow class=formRow")
	if ue.Venue != nil && ue.Venue.URL() == "" {
		row.E("label for=eventeditVenueURL>Venue URL")
		row.E("input id=eventeditVenueURL name=venueURL s-validate value=%s", ue.VenueURL)
		if err != "" {
			row.E("div class=formError>%s", err)
		}
	} else {
		row.Attr("style=display:none")
	}
}

var htmlSanitizer = bluemonday.NewPolicy().
	RequireParseableURLs(true).
	AllowURLSchemes("http", "https").
	RequireNoFollowOnLinks(true).
	AllowAttrs("href").OnElements("a").
	AddTargetBlankToFullyQualifiedLinks(true)

func readEventDetails(r *request.Request, ue *event.Updater) {
	ue.Details = htmlSanitizer.Sanitize(strings.TrimSpace(r.FormValue("details")))
}
func emitEventDetails(form *htmlb.Element, ue *event.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=eventeditDetails>Details")
	row.E("textarea id=eventeditDetails name=details wrap=soft rows=3").T(ue.Details)
	row.E("div class=formHelp>This may contain HTML &lt;a&gt; tags for links, but no other tags.")
}

func emitDetailsButtons(form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
}

// updateShiftDates updates the date in the start and end fields of all shifts
// belonging to the event.
func updateShiftDates(r *request.Request, e *event.Event, ue *event.Updater) {
	task.AllForEvent(r, ue.ID, task.FID|task.FName, func(t *task.Task) {
		shift.AllForTask(r, t.ID(), shift.UpdaterFields, venue.FID|venue.FName, func(s *shift.Shift, v *venue.Venue) {
			us := s.Updater(r, e, t, v)
			us.Start = ue.Start[:10] + us.Start[10:]
			us.End = ue.End[:10] + us.End[10:]
			s.Update(r, us)
		})
	})
}
