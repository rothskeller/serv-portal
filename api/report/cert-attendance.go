package report

import (
	"encoding/csv"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

type columnKey string
type eventTypeAbbr string

var dateRE = regexp.MustCompile(`^20\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
var eventTypeAbbrs = map[model.EventType]eventTypeAbbr{
	model.EventCivic:    "Civ",
	model.EventClass:    "Cls",
	model.EventContEd:   "CE",
	model.EventDrill:    "Drl",
	model.EventIncident: "Inc",
	model.EventMeeting:  "Mtg",
	model.EventTraining: "Trn",
}

// CERTAttendanceReport handles GET /api/reports/cert-attendance requests.
func CERTAttendanceReport(r *util.Request) error {
	var (
		team     *model.Role
		events   []*model.Event
		people   []*model.Person
		rendered attendanceReport
		data     = map[model.PersonID]map[columnKey]map[eventTypeAbbr]int{}
		pmap     = map[model.PersonID]*model.Person{}
		teamStr  = r.FormValue("team")
		dateFrom = r.FormValue("dateFrom")
		dateTo   = r.FormValue("dateTo")
		stats    = r.FormValue("stats")
		detail   = r.FormValue("detail")
		format   = r.FormValue("format")
	)
	switch teamStr {
	case "Alpha":
		team = r.Tx.FetchRoleByTag("cert-team-alpha")
	case "Bravo":
		team = r.Tx.FetchRoleByTag("cert-team-bravo")
	default:
		team = r.Tx.FetchRoleByTag("cert-teams")
	}
	if !auth.CanManageEvents(r, team) {
		return util.Forbidden
	}
	if !dateRE.MatchString(dateFrom) {
		now := time.Now()
		dateFrom = time.Date(now.Year(), now.Month()/3*3+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02")
	}
	if !dateRE.MatchString(dateTo) || dateTo <= dateFrom {
		now := time.Now()
		dateTo = time.Date(now.Year(), (now.Month()/3+1)*3+1, 1, 0, 0, 0, 0, time.Local).Add(-2 * time.Hour).Format("2006-01-02")
	}
	if stats != "hours" {
		stats = "count"
	}
	if detail != "date" && detail != "total" {
		detail = "month"
	}
	// Get the events to which CERT was invited during the time range.
	events = r.Tx.FetchEvents(dateFrom, dateTo)
	j := 0
	for _, e := range events {
		found := false
		for _, t := range e.Roles {
			switch t.Tag {
			case "cert-teams", "cert-team-alpha", "cert-team-bravo":
				found = true
				break
			}
		}
		if found {
			events[j] = e
			j++
		}
	}
	events = events[:j]
	// Get all relevant people.
	people = r.Tx.FetchPeople()
	j = 0
	for _, p := range people {
		if auth.HasRole(p, team) {
			people[j] = p
			j++
			pmap[p.ID] = p
		}
	}
	people = people[:j]
	// Get the attendance data.
	for _, e := range events {
		for pid := range r.Tx.FetchAttendanceByEvent(e) {
			if pmap[pid] == nil {
				continue
			}
			addAttendance(data, e, pid, stats, detail)
			addAttendance(data, e, 0, stats, detail)
		}
		addAttendance(data, e, -1, stats, detail)
	}
	r.Tx.Commit()
	// Convert the report into output-format-independent rows and columns.
	rendered = renderAttendance(data, people, stats, detail)
	if format == "CSV" {
		attendanceCSV(r, rendered)
	} else {
		attendanceJSON(r, rendered)
	}
	return nil
}

func addAttendance(
	data map[model.PersonID]map[columnKey]map[eventTypeAbbr]int, event *model.Event, pid model.PersonID, stats, detail string,
) {
	if data[pid] == nil {
		data[pid] = make(map[columnKey]map[eventTypeAbbr]int)
	}
	switch detail {
	case "date":
		addAttendance2(data[pid], columnKey(event.Date), event, stats)
	case "month":
		addAttendance2(data[pid], columnKey(event.Date[:7]), event, stats)
	default:
	}
	addAttendance2(data[pid], "TOTALS", event, stats)
}
func addAttendance2(data map[columnKey]map[eventTypeAbbr]int, key columnKey, event *model.Event, stats string) {
	if data[key] == nil {
		data[key] = make(map[eventTypeAbbr]int)
	}
	if stats == "hours" {
		data[key][eventTypeAbbrs[event.Type]] += int(2 * event.Hours())
		data[key]["ALL"] += int(2 * event.Hours())
	} else {
		data[key][eventTypeAbbrs[event.Type]]++
		data[key]["ALL"]++
	}
}

type attendanceReport struct {
	Header     [][]attendanceReportHeadCell
	Body       [][]string
	Footer     [][]string
	SpanStarts map[int]bool
}
type attendanceReportHeadCell struct {
	Span int
	Text string
}

func renderAttendance(
	data map[model.PersonID]map[columnKey]map[eventTypeAbbr]int, people []*model.Person, stats, detail string,
) (report attendanceReport) {
	var (
		etypes []eventTypeAbbr
		keys   []columnKey
		col    int
	)
	// Get the sorted list of keys.
	keys = make([]columnKey, 0, len(data[-1])-1)
	for key := range data[-1] {
		if key != "TOTALS" {
			keys = append(keys, key)
		}
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	// Get the complete set of event types.
	{
		etmap := map[eventTypeAbbr]bool{}
		for et := range data[-1]["TOTALS"] {
			if et != "ALL" {
				etmap[et] = true
			}
		}
		etypes = make([]eventTypeAbbr, 0, len(etmap))
		for _, et := range model.AllEventTypes {
			if eta := eventTypeAbbrs[et]; etmap[eta] {
				etypes = append(etypes, eta)
			}
		}
	}
	// Create the rows and lay out the leftmost column.
	report.Header = append(report.Header, []attendanceReportHeadCell{{}})
	col++
	report.Header = append(report.Header, []attendanceReportHeadCell{{}})
	if stats == "hours" {
		report.Header = append(report.Header, []attendanceReportHeadCell{{}})
	}
	report.Body = make([][]string, len(people))
	for i, p := range people {
		report.Body[i] = []string{fmt.Sprintf("%s, %s", p.LastName, p.FirstName)}
	}
	report.Footer = [][]string{{"TOTALS"}}
	// Create and fill the non-total columns.
	report.SpanStarts = make(map[int]bool)
	if detail != "total" {
		for _, key := range keys {
			var ketypes []eventTypeAbbr
			if detail == "month" {
				// For month reports, show all of the etypes in
				// the report every month, even if their totals
				// are zero.
				ketypes = etypes
			} else {
				// For date reports, show only the etypes
				// actually used on that date.
				for _, et := range etypes {
					for used := range data[-1][key] {
						if used == et {
							ketypes = append(ketypes, et)
						}
					}
				}
			}
			// First row: key names, spanning.
			report.SpanStarts[col] = true
			col += len(ketypes)
			report.Header[0] = append(report.Header[0], attendanceReportHeadCell{Span: len(ketypes), Text: string(key)})
			for _, et := range ketypes {
				// Second row: event type.
				report.Header[1] = append(report.Header[1], attendanceReportHeadCell{Text: string(et)})
				if stats == "hours" {
					// Third row: hours in the events.
					report.Header[2] = append(report.Header[2], attendanceReportHeadCell{
						Text: renderAttendanceValue(data, -1, key, et, stats),
					})
				}
				// Rows for people.
				for i, p := range people {
					report.Body[i] = append(report.Body[i], renderAttendanceValue(data, p.ID, key, et, stats))
				}
				// Totals row.
				report.Footer[0] = append(report.Footer[0], renderAttendanceTotal(data, 0, key, et, stats))
			}
		}
	}
	// Create and fill the total columns.
	if len(keys) != 1 {
		span := len(etypes)
		if span != 1 {
			span++
		}
		// First row: "TOTALS" spanning label.
		report.SpanStarts[col] = true
		report.Header[0] = append(report.Header[0], attendanceReportHeadCell{Span: span, Text: "TOTALS"})
		for _, et := range etypes {
			// Second row: event type.
			report.Header[1] = append(report.Header[1], attendanceReportHeadCell{Text: string(et)})
			if stats == "hours" {
				// Third row: hours in the events.
				report.Header[2] = append(report.Header[2], attendanceReportHeadCell{
					Text: renderAttendanceValue(data, -1, "TOTALS", et, stats),
				})
			}
			// Rows for people.
			for i, p := range people {
				report.Body[i] = append(report.Body[i], renderAttendanceValue(data, p.ID, "TOTALS", et, stats))
			}
			// Totals row.
			report.Footer[0] = append(report.Footer[0], renderAttendanceTotal(data, 0, "TOTALS", et, stats))
		}
		if len(etypes) != 1 {
			// Second row: event type.
			report.Header[1] = append(report.Header[1], attendanceReportHeadCell{Text: "ALL"})
			if stats == "hours" {
				// Third row: hours in the events.
				report.Header[2] = append(report.Header[2], attendanceReportHeadCell{
					Text: renderAttendanceTotal(data, -1, "TOTALS", "ALL", stats),
				})
			}
			// Rows for people.
			for i, p := range people {
				report.Body[i] = append(report.Body[i], renderAttendanceTotal(data, p.ID, "TOTALS", "ALL", stats))
			}
			// Totals row.
			report.Footer[0] = append(report.Footer[0], renderAttendanceTotal(data, 0, "TOTALS", "ALL", stats))
		}
	}
	return report
}

func attendanceJSON(r *util.Request, report attendanceReport) {
	var out jwriter.Writer

	out.RawString(`{"header":[`)
	for i, h := range report.Header {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawByte('[')
		for j, c := range h {
			if j != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"span":`)
			out.Int(c.Span)
			out.RawString(`,"text":`)
			out.String(c.Text)
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	out.RawString(`],"body":[`)
	for i, b := range report.Body {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawByte('[')
		for j, c := range b {
			if j != 0 {
				out.RawByte(',')
			}
			out.String(c)
		}
		out.RawByte(']')
	}
	out.RawString(`],"footer":[`)
	for i, f := range report.Footer {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawByte('[')
		for j, c := range f {
			if j != 0 {
				out.RawByte(',')
			}
			out.String(c)
		}
		out.RawByte(']')
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
}

func attendanceCSV(r *util.Request, report attendanceReport) {
	var (
		cols = []string{}
		out  = csv.NewWriter(r)
	)
	r.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Header().Set("Content-Disposition", `attachment; filename="attendance.csv"`)
	out.UseCRLF = true
	for _, row := range report.Header {
		for _, cell := range row {
			cols = append(cols, cell.Text)
			for i := 1; i < cell.Span; i++ {
				cols = append(cols, "")
			}
		}
		out.Write(cols)
	}
	for _, row := range report.Body {
		out.Write(row)
	}
	for _, row := range report.Footer {
		out.Write(row)
	}
	out.Flush()
}

func renderAttendanceValue(
	data map[model.PersonID]map[columnKey]map[eventTypeAbbr]int, pid model.PersonID, key columnKey, etype eventTypeAbbr,
	stats string,
) string {
	var value int
	if data[pid] != nil && data[pid][key] != nil {
		value = data[pid][key][etype]
	}
	if value == 0 {
		return ""
	}
	if stats == "hours" {
		return fmt.Sprintf("%.1f", float64(value)/2.0)
	}
	return strconv.Itoa(value)
}

func renderAttendanceTotal(
	data map[model.PersonID]map[columnKey]map[eventTypeAbbr]int, pid model.PersonID, key columnKey, etype eventTypeAbbr,
	stats string,
) string {
	var value int
	if data[pid] != nil && data[pid][key] != nil {
		value = data[pid][key][etype]
	}
	if stats == "hours" {
		return fmt.Sprintf("%.1f", float64(value)/2.0)
	}
	return strconv.Itoa(value)
}
