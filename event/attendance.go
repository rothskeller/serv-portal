package event

import (
	"errors"
	"html/template"
	"net/http"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// GetEventAttendance handles GET /events/$id/attendance requests.
func GetEventAttendance(r *util.Request, idstr string) error {
	var (
		gead  getEventAttendanceData
		title string
	)
	if gead.Event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); gead.Event == nil {
		return util.NotFound
	}
	if !r.Person.CanRecordAttendanceAtEvent(gead.Event) {
		return util.Forbidden
	}
	gead.People = r.Tx.FetchPeople()
	gead.Attended = r.Tx.FetchAttendanceByEvent(gead.Event)
	r.Tx.Commit()
	gead.Year = gead.Event.Date[:4]
	title = gead.Event.Date + " " + gead.Event.Name
	j := 0
	for _, p := range gead.People {
		if p.CanViewEvent(gead.Event) {
			gead.People[j] = p
			j++
		}
	}
	gead.People = gead.People[:j]
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "events",
		BodyData: &gead,
	}, template.Must(template.New("getEventAttendance").Parse(getEventAttendanceTemplate)))
	return nil
}

// PostEventAttendance handles POST /events/$id/attendance requests.
func PostEventAttendance(r *util.Request, idstr string) error {
	var (
		event  *model.Event
		person *model.Person
		people []*model.Person
	)
	if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
		return util.NotFound
	}
	if !r.Person.CanRecordAttendanceAtEvent(event) {
		return util.Forbidden
	}
	r.ParseForm()
	for _, idstr := range r.Form["person"] {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return errors.New("invalid person")
		}
		if !person.CanViewEvent(event) {
			return errors.New("illegal person")
		}
		people = append(people, person)
	}
	r.Tx.SaveEventAttendance(event, people)
	r.Tx.Commit()
	http.Redirect(r, r.Request, "/events?year="+event.Date[:4], http.StatusSeeOther)
	return nil
}

type getEventAttendanceData struct {
	Event    *model.Event
	Year     string
	People   []*model.Person
	Attended map[model.PersonID]bool
}

const getEventAttendanceTemplate = `{{ define "body" -}}
<div id="editEvent">
  <div class="pageTitle">Event Attendance</div>
  <form method="POST">
    <div id="getEventAttendance-group-label">This event was attended by:</div>
    {{- $att := .Attended }}
    {{- range .People }}
      <div class="getEventAttendance-group-row">
        <input id="getEventAttendance-person-{{ .ID }}" name="person" type="checkbox" value="{{ .ID }}"{{ if index $att .ID }} checked{{ end }}>
	<label for="getEventAttendance-person-{{ .ID }}">{{ .LastName }}, {{ .FirstName }}</label>
      </div>
    {{ end }}
    <div id="getEventAttendance-submit-row">
      <button type="submit" class="btn btn-primary">Save Attendance</button>
      <a class="btn btn-secondary" href="/events?year={{ .Year }}">Cancel</a>
    </div>
  </form>
</div>
{{- end }}`
