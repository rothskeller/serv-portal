package event

import (
	"errors"
	"html/template"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

var dateRE = regexp.MustCompile(`^20\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
var yearRE = regexp.MustCompile(`^20\d\d$`)

var eventTypeNames = map[model.EventType]string{
	model.EventTraining: "Training",
	model.EventCivic:    "Civic Event",
	model.EventClass:    "Class",
	model.EventContEd:   "Continuing Ed",
	model.EventDrill:    "Drill",
	model.EventIncident: "Incident",
	model.EventMeeting:  "Meeting",
}

// EditEvent handles GET and POST /events/$id requests (where $id may be "NEW").
func EditEvent(r *util.Request, idstr string) error {
	var (
		eed   editEventData
		title string
		err   error
	)
	eed.Teams = r.Person.SchedulableTeams()
	if idstr == "NEW" {
		if eed.Teams == nil {
			return util.Forbidden
		}
		eed.Event = &model.Event{Hours: 1.0}
		title = "New Event"
	} else {
		if eed.Event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); eed.Event == nil {
			return util.NotFound
		}
		if !r.Person.CanViewEvent(eed.Event) {
			return util.Forbidden
		}
		if !r.Person.CanManageEvent(eed.Event) {
			return showEvent(r, eed.Event)
		}
		title = eed.Event.Date + " " + eed.Event.Name
	}
	if eed.Year = r.FormValue("year"); eed.Year != "" && !yearRE.MatchString(eed.Year) {
		return errors.New("bad year")
	}
	if r.Method == http.MethodPost {
		if r.FormValue("delete") != "" && eed.Event.ID != 0 {
			r.Tx.DeleteEvent(eed.Event)
			http.Redirect(r, r.Request, "/events?year="+eed.Year, http.StatusSeeOther)
			r.Tx.Commit()
			return nil
		}
		eed.Event.Date = r.FormValue("date")
		if eed.Event.Date == "" {
			eed.DateError = "The event date is required."
		} else if !dateRE.MatchString(eed.Event.Date) {
			return errors.New("invalid date (YYYY-MM-DD)")
		}
		eed.Year = eed.Event.Date[:4]
		if eed.Event.Name = strings.TrimSpace(r.FormValue("name")); eed.Event.Name == "" {
			eed.NameError = "The event name is required."
		}
		eed.Event.Hours, err = strconv.ParseFloat(r.FormValue("hours"), 64)
		if err != nil || eed.Event.Hours < 0.0 || eed.Event.Hours > 24.0 || math.IsNaN(eed.Event.Hours) {
			eed.HoursError = "The event hours are not valid."
		}
		eed.Event.Hours = math.Round(eed.Event.Hours*2.0) / 2.0
		if eed.Event.Type = model.EventType(r.FormValue("type")); eed.Event.Type == "" {
			eed.TypeError = "The event type is required."
		} else {
			found := false
			for _, t := range model.AllEventTypes {
				if eed.Event.Type == t {
					found = true
					break
				}
			}
			if !found {
				return errors.New("invalid event type")
			}
		}
		eed.Event.Teams = eed.Event.Teams[:0]
		for _, idstr := range r.Form["team"] {
			team := r.Tx.FetchTeam(model.TeamID(util.ParseID(idstr)))
			if team == nil {
				return errors.New("invalid team")
			}
			found := false
			for _, t := range eed.Teams {
				if t == team {
					found = true
					break
				}
			}
			if !found {
				return util.Forbidden
			}
			eed.Event.Teams = append(eed.Event.Teams, team)
		}
		if len(eed.Event.Teams) == 0 {
			eed.TeamsError = "At least one team must be selected."
		}
		if eed.DateError != "" || eed.NameError != "" || eed.HoursError != "" || eed.TypeError != "" || eed.TeamsError != "" {
			goto SHOWFORM
		}
		for _, e := range r.Tx.FetchEvents(eed.Event.Date, eed.Event.Date) {
			if e.ID != eed.Event.ID && e.Name == eed.Event.Name {
				eed.NameError = "Another event on this date already has this name."
				goto SHOWFORM
			}
		}
		r.Tx.SaveEvent(eed.Event)
		r.Tx.Commit()
		http.Redirect(r, r.Request, "/events?year="+eed.Year, http.StatusSeeOther)
		return nil
	}
SHOWFORM:
	r.Tx.Commit()
	eed.Types = model.AllEventTypes
	eed.TeamMap = make(map[model.TeamID]bool)
	eed.EventTypeNames = eventTypeNames
	for _, t := range eed.Event.Teams {
		eed.TeamMap[t.ID] = true
	}
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "events",
		BodyData: &eed,
	}, template.Must(template.New("editEvent").Funcs(map[string]interface{}{
		"formatHours": formatHours,
	}).Parse(editEventTemplate)))
	return nil
}

type editEventData struct {
	DateError      string
	NameError      string
	HoursError     string
	TypeError      string
	TeamsError     string
	Event          *model.Event
	Year           string
	TeamMap        map[model.TeamID]bool
	Types          []model.EventType
	Teams          []*model.Team
	EventTypeNames map[model.EventType]string
}

const editEventTemplate = `{{ define "body" -}}
<div id="editEvent">
  <div class="pageTitle">{{ if .Event.ID }}Edit Event{{ else }}Create Event{{ end }}</div>
  <form method="POST">
    <input type="hidden" name="year" value="{{ .Year }}">
    <div id="editEvent-date-row">
      <label for="editEvent-date" id="editEvent-date-label">Event date</label>
      <input id="editEvent-date" name="date" type="date" class="form-control" autofocus value="{{ .Event.Date }}">
    </div>
    {{- if .DateError }}<div id="editEvent-date-error">{{ .DateError }}</div>{{ end }}
    <div id="editEvent-name-row">
      <label for="editEvent-name" id="editEvent-name-label">Event name</label>
      <input id="editEvent-name" name="name" class="form-control" value="{{ .Event.Name }}">
    </div>
    {{- if .NameError }}<div id="editEvent-name-error">{{ .NameError }}</div>{{ end }}
    <div id="editEvent-hours-row">
      <label for="editEvent-hours" id="editEvent-hours-label">Event hours</label>
      <input id="editEvent-hours" name="hours" type="number" min="0.0" max="24.0" step="0.5" class="form-control" value="{{ .Event.Hours }}">
    </div>
    {{- if .HoursError }}<div id="editEvent-hours-error">{{ .HoursError }}</div>{{ end }}
    {{- $et := .Event.Type }}
    {{- $etn := .EventTypeNames }}
    <div class="editEvent-group-label">Event type:</div>
    {{- range $id, $type := .Types }}
      <div class="editEvent-group-row">
        <input id="editEvent-type-{{ $id }}" name="type" type="radio" value="{{ . }}"{{ if eq . $et }} checked{{ end }}>
	<label for="editEvent-type-{{ $id }}">{{ index $etn . }}</label>
      </div>
    {{ end }}
    {{- if .TypeError }}<div id="editEvent-type-error">{{ .TypeError }}</div>{{ end }}
    {{- $tmap := .TeamMap }}
    <div class="editEvent-group-label">Event is for these teams:</div>
    {{- range .Teams }}
      <div class="editEvent-group-row">
        <input id="editEvent-team-{{ .ID }}" name="team" type="checkbox" value="{{ .ID }}"{{ if index $tmap .ID }} checked{{ end }}>
	  <label for="editEvent-team-{{ .ID }}">{{ .Name }}</label>
      </div>
    {{ end }}
    {{- if .TeamsError }}<div id="editEvent-teams-error">{{ .TeamsError }}</div>{{ end }}
    <div id="editEvent-submit-row">
      <button type="submit" class="btn btn-primary">{{ if .Event.ID }}Save Event{{ else }}Create Event{{ end }}</button>
      <a class="btn btn-secondary" href="/events?year={{ .Year }}">Cancel</a>
      {{- if .Event.ID }}
        <button id="editEvent-delete" name="delete" type="submit" class="btn btn-danger">Delete Event</button>
      {{ end }}
    </div>
  </form>
</div>
{{- end }}`
