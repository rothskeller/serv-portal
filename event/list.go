package event

import (
	"fmt"
	"html/template"
	"strconv"
	"strings"
	"time"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// ListEvents handles GET /events requests.
func ListEvents(r *util.Request) error {
	var led listEventsData

	if led.Year, _ = strconv.Atoi(r.FormValue("year")); led.Year < 2000 || led.Year > 2099 {
		led.Year = time.Now().Year()
	}
	led.Events = r.Tx.FetchEvents(fmt.Sprintf("%d-01-01", led.Year), fmt.Sprintf("%d-12-31", led.Year))
	r.Tx.Commit()
	j := 0
	for _, e := range led.Events {
		if r.Person.CanViewEvent(e) {
			led.Events[j] = e
			j++
		}
	}
	led.Events = led.Events[:j]
	for year := 2019; year <= time.Now().Year()+1; year++ {
		led.Years = append(led.Years, year)
	}
	led.CanAdd = r.Person.CanCreateEvents()
	led.CanAttendance = make(map[*model.Event]bool)
	for _, e := range led.Events {
		led.CanAttendance[e] = r.Person.CanRecordAttendanceAtEvent(e)
	}
	util.RenderPage(r, &util.Page{
		Title:    "Events",
		MenuItem: "events",
		BodyData: &led,
	}, template.Must(template.New("listEvents").Funcs(map[string]interface{}{
		"formatHours": formatHours,
	}).Parse(listEventsTemplate)))
	return nil
}

func formatHours(hours float64) (s string) {
	return strings.Replace(strings.Replace(fmt.Sprintf("%.1f", hours), ".5", "½", -1), ".0", "", -1)
}

type listEventsData struct {
	Year          int
	Events        []*model.Event
	Years         []int
	CanAdd        bool
	CanAttendance map[*model.Event]bool
}

const listEventsTemplate = `{{ define "body" -}}
<div id="listEvents">
  <form id="listEvents-title" method="GET" class="pageTitle">
    Events
    <select id="listEvents-year" name="year">
      {{- $year := .Year }}
      {{- range .Years }}
	<option value="{{ . }}"{{ if eq . $year }} selected{{ end }}>{{ . }}</option>
      {{ end }}
    </select>
  </form>
  <table id="listEvents-table">
    <thead>
      <tr>
        <th>Date</th>
	<th>Event</th>
	<th></th>
      </tr>
    </thead>
    <tbody>
      {{- $cra := .CanAttendance }}
      {{- range .Events }}
        <tr>
	  <td>{{ .Date }}</td>
	  <td><a href="/events/{{ .ID }}">{{ .Name }}</a></td>
	  <td>
	    {{- if index $cra . }}
	      <a href="/events/{{ .ID }}/attendance">Attendance</a>
	    {{ end }}
	  </td>
	</tr>
      {{ end }}
    </tbody>
  </table>
  {{- if .CanAdd }}
    <div id="listEvents-buttons">
      <a class="btn btn-secondary" href="/events/NEW?year={{ .Year }}">Add Event</a>
    </div>
  {{ end }}
</div>
{{- end }}`
