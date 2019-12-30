package event

import (
	"html/template"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

func showEvent(r *util.Request, event *model.Event) error {
	var (
		title string
		sed   = showEventData{Event: event}
	)
	r.Tx.Commit()
	title = sed.Event.Date + " " + sed.Event.Name
	sed.Year = sed.Event.Date[:4]
	sed.EventTypeNames = eventTypeNames
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "events",
		BodyData: &sed,
	}, template.Must(template.New("showEvent").Funcs(map[string]interface{}{
		"formatHours": formatHours,
	}).Parse(showEventTemplate)))
	return nil
}

type showEventData struct {
	Event          *model.Event
	Year           string
	EventTypeNames map[model.EventType]string
}

const showEventTemplate = `{{ define "body" -}}
<div id="editEvent">
  <div class="pageTitle">Event Details</div>
  <div id="editEvent-date-row">
    <label id="editEvent-date-label">Event date</label>
    <div>{{ .Event.Date }}</div>
  </div>
  <div id="editEvent-name-row">
    <label id="editEvent-name-label">Event name</label>
    <div>{{ .Event.Name }}</div>
  </div>
  <div id="editEvent-hours-row">
    <label id="editEvent-hours-label">Event hours</label>
    <div>{{ formatHours .Event.Hours }}</div>
  </div>
  <div id="editEvent-type-row">
    <label id="editEvent-type-label">Event type</label>
    <div>{{ index .EventTypeNames .Event.Type }}</div>
  </div>
  <div id="editEvent-teams-row">
    <label id="editEvent-teams-label">Teams</label>
    <div>
      {{- range .Event.Teams }}
        <div>{{ .Name }}</div>
      {{ end }}
    </div>
  </div>
  <div id="editEvent-submit-row">
    <a class="btn btn-secondary" href="/events?year={{ .Year }}">Back</a>
  </div>
</div>
{{- end }}`
