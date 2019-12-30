package report

import (
	"html/template"
	"time"

	"serv.rothskeller.net/portal/util"
)

// GetIndex handles GET /reports requests.
func GetIndex(r *util.Request) error {
	var (
		rid reportIndexData
		now time.Time
	)
	now = time.Now()
	rid.DateFrom = time.Date(now.Year(), now.Month()/3*3+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	rid.DateTo = time.Date(now.Year(), (now.Month()/3+1)*3+1, 1, 0, 0, 0, 0, time.Local).Add(-2 * time.Hour).Format("2006-01-02")
	util.RenderPage(r, &util.Page{
		Title:    "Reports",
		MenuItem: "reports",
		BodyData: &rid,
	}, template.Must(template.New("reportIndex").Parse(reportIndexTemplate)))
	return nil
}

type reportIndexData struct {
	DateFrom string
	DateTo   string
}

const reportIndexTemplate = `{{ define "body" -}}
<div id="reportIndex">
  <div class="pageTitle">CERT Attendance Report</div>
  <form id="reportIndex-cert-attendance-form" method="GET" action="/reports/cert-attendance">
    <div class="reportIndex-form-row">
      <div class="reportIndex-form-label">Report on team</div>
      <div>
        <input id="reportIndex-cert-attendance-team-alpha" name="team" type="radio" value="Alpha">
	<label for="reportIndex-cert-attendance-team-alpha">Alpha</label>
        <input id="reportIndex-cert-attendance-team-bravo" name="team" type="radio" value="Bravo">
	<label for="reportIndex-cert-attendance-team-bravo">Bravo</label>
        <input id="reportIndex-cert-attendance-team-both" name="team" type="radio" value="Both" checked>
	<label for="reportIndex-cert-attendance-team-both">Both</label>
      </div>
    </div>
    <div class="reportIndex-form-row">
      <div class="reportIndex-form-label">Date range</div>
      <div>
        <input name="dateFrom" type="date" value="{{ .DateFrom }}">
	through
        <input name="dateTo" type="date" value="{{ .DateTo }}">
      </div>
    </div>
    <div class="reportIndex-form-row">
      <div class="reportIndex-form-label">Statistics by</div>
      <div>
        <input id="reportIndex-cert-attendance-stats-count" name="stats" type="radio" value="count" checked>
	<label for="reportIndex-cert-attendance-stats-count">Number of Events</label>
        <input id="reportIndex-cert-attendance-stats-hours" name="stats" type="radio" value="hours">
	<label for="reportIndex-cert-attendance-stats-hours">Cumulative Hours</label>
      </div>
    </div>
    <div class="reportIndex-form-row">
      <div class="reportIndex-form-label">Show detail</div>
      <div>
        <input id="reportIndex-cert-attendance-detail-event" name="detail" type="radio" value="event">
	<label for="reportIndex-cert-attendance-detail-event">Show each event</label>
        <input id="reportIndex-cert-attendance-detail-month" name="detail" type="radio" value="month" checked>
	<label for="reportIndex-cert-attendance-detail-month">Show each month</label>
        <input id="reportIndex-cert-attendance-detail-total" name="detail" type="radio" value="total">
	<label for="reportIndex-cert-attendance-detail-total">Show totals only</label>
      </div>
    </div>
    <div class="reportIndex-form-row">
      <div class="reportIndex-form-label"></div>
      <div><button type="submit" class="btn btn-primary">Generate Report</button></div>
    </div>
  </form>
</div>
{{- end }}`
