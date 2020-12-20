package report

import (
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

type attrepParameters struct {
	dateRange       string
	dateFrom        string
	dateTo          string
	eventTypes      map[model.EventType]bool
	orgs            map[model.Org]bool
	attendanceTypes map[model.AttendanceType]bool
	collapseX       bool
	collapseY       bool
	includeZerosX   bool
	includeZerosY   bool
	sumHours        bool
	groupByOrg      bool
	renderCSV       bool
	// Not really parameters, but cached here for convenience:
	allowedOrgs map[model.Org]bool
	dateRanges  []attrepDateRange
}
type attrepDateRange struct {
	tag      string
	label    string
	dateFrom string
	dateTo   string
}

func readAttrepParameters(r *util.Request) (params attrepParameters) {
	params.allowedOrgs = allowedOrgs(r)
	params.dateRanges = dateRanges()

	// Basic report options.
	switch r.FormValue("rows") {
	default:
		params.collapseY, params.groupByOrg = true, false
	case "o":
		params.collapseY, params.groupByOrg = true, true
	case "op":
		params.collapseY, params.groupByOrg = false, true
	case "po":
		params.collapseY, params.groupByOrg = false, false
	}
	params.collapseX = r.FormValue("columns") == "m"
	params.sumHours = r.FormValue("cells") != "c"
	params.includeZerosX, _ = strconv.ParseBool(r.FormValue("includeZerosX"))
	params.includeZerosY, _ = strconv.ParseBool(r.FormValue("includeZerosY"))
	params.renderCSV = r.FormValue("format") == "csv"

	// Date range.
	params.dateRange = r.FormValue("dateRange")
	for _, dr := range params.dateRanges {
		if params.dateRange == dr.tag {
			params.dateFrom, params.dateTo = dr.dateFrom, dr.dateTo
			break
		}
	}
	if params.dateFrom == "" {
		params.dateRange = params.dateRanges[0].tag
		params.dateFrom = params.dateRanges[0].dateFrom
		params.dateTo = params.dateRanges[0].dateTo
	}

	// Event types.
	params.eventTypes = make(map[model.EventType]bool)
	for _, etidstr := range strings.Split(r.FormValue("eventTypes"), ",") {
		if et := model.EventType(util.ParseID(etidstr)); model.EventTypeNames[et] != "" {
			params.eventTypes[et] = true
		}
	}
	if len(params.eventTypes) == 0 { // No event types selected, so include them all.
		for _, et := range model.AllEventTypes {
			if params.sumHours || et != model.EventHours {
				params.eventTypes[et] = true
			}
		}
	}

	// Organizations
	params.orgs = make(map[model.Org]bool)
	for _, oidstr := range strings.Split(r.FormValue("orgs"), ",") {
		if o := model.Org(util.ParseID(oidstr)); orgNames[o] != "" {
			params.orgs[o] = params.allowedOrgs[o]
		}
	}
	if len(params.orgs) == 0 { // No organizations selected, so include all allowed.
		params.orgs = params.allowedOrgs
	}

	// Attendance types.
	params.attendanceTypes = make(map[model.AttendanceType]bool)
	for _, atidstr := range strings.Split(r.FormValue("attendanceTypes"), ",") {
		if at := model.AttendanceType(util.ParseID(atidstr)); at.Valid() {
			params.attendanceTypes[at] = true
		}
	}
	if len(params.attendanceTypes) == 0 { // No types selected.
		params.attendanceTypes[model.AttendAsVolunteer] = true
		if params.sumHours {
			params.attendanceTypes[model.AttendAsAbsent] = true
		}
	}

	// If we're doing attendance counts, forcibly leave out non-attendance
	// events and hours.
	if !params.sumHours {
		params.attendanceTypes[model.AttendAsAbsent] = false
		params.eventTypes[model.EventHours] = false
	}
	return params
}

func attrepRenderParams(out *jwriter.Writer, params attrepParameters) {
	out.RawString(`"parameters":{"dateRange":`)
	out.String(params.dateRange)
	out.RawString(`,"dateFrom":`)
	out.String(params.dateFrom)
	out.RawString(`,"dateTo":`)
	out.String(params.dateTo)
	out.RawString(`,"rows":`)
	switch {
	case params.collapseY && params.groupByOrg:
		out.String("o")
	case params.collapseY && !params.groupByOrg:
		out.String("p")
	case !params.collapseY && params.groupByOrg:
		out.String("op")
	case !params.collapseY && !params.groupByOrg:
		out.String("po")
	}
	if params.collapseX {
		out.RawString(`,"columns":"m"`)
	} else {
		out.RawString(`,"columns":"e"`)
	}
	if params.sumHours {
		out.RawString(`,"cells":"h"`)
	} else {
		out.RawString(`,"cells":"c"`)
	}
	out.RawString(`,"includeZerosX":`)
	out.Bool(params.includeZerosX)
	out.RawString(`,"includeZerosY":`)
	out.Bool(params.includeZerosY)
	out.RawString(`,"orgs":[`)
	var first = true
	for o, v := range params.orgs {
		if !v {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.Int(int(o))
	}
	out.RawString(`],"eventTypes":[`)
	first = true
	for et, v := range params.eventTypes {
		if !v {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.Int(int(et))
	}
	out.RawString(`],"attendanceTypes":[`)
	first = true
	for at, v := range params.attendanceTypes {
		if !v {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.Int(int(at))
	}
	out.RawString(`]},"options":{"dateRanges":[`)
	for i, dr := range params.dateRanges {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"tag":`)
		out.String(dr.tag)
		out.RawString(`,"label":`)
		out.String(dr.label)
		out.RawString(`,"dateFrom":`)
		out.String(dr.dateFrom)
		out.RawString(`,"dateTo":`)
		out.String(dr.dateTo)
		out.RawByte('}')
	}
	out.RawString(`],"orgs":[`)
	first = true
	for _, o := range model.AllOrgs {
		if !params.allowedOrgs[o] {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(o))
		out.RawString(`,"label":`)
		out.String(orgNames[o])
		out.RawByte('}')
	}
	out.RawString(`],"eventTypes":[`)
	for i, et := range model.AllEventTypes {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(et))
		out.RawString(`,"label":`)
		out.String(model.EventTypeNames[et])
		out.RawByte('}')
	}
	out.RawString(`],"attendanceTypes":[`)
	for i, at := range model.AllAttendanceTypes {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(at))
		out.RawString(`,"label":`)
		out.String(at.String())
		out.RawByte('}')
	}
	out.RawString(`]}`)
}

// allowedOrgs returns a map of which organizations the current caller has
// reporting capabilities on.
func allowedOrgs(r *util.Request) (orgs map[model.Org]bool) {
	orgs = make(map[model.Org]bool)
	for _, o := range model.AllOrgs {
		if r.Person.Orgs[o].PrivLevel >= model.PrivLeader {
			orgs[o] = true
		}
	}
	return orgs
}

// dateRanges generates the set of date ranges supported by the reports.
func dateRanges() []attrepDateRange {
	var now = time.Now()
	return []attrepDateRange{
		{
			tag:      "tm",
			label:    "this month",
			dateFrom: time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
		{
			tag:      "lm",
			label:    "last month",
			dateFrom: time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
		{
			tag:      "tq",
			label:    "this quarter",
			dateFrom: time.Date(now.Year(), (now.Month()-1)/3*3+1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), (now.Month()-1)/3*3+4, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
		{
			tag:      "lq",
			label:    "last quarter",
			dateFrom: time.Date(now.Year(), (now.Month()-1)/3*3-2, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), (now.Month()-1)/3*3+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
		{
			tag:      "ty",
			label:    "this year",
			dateFrom: time.Date(now.Year(), 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
		{
			tag:      "ly",
			label:    "last year",
			dateFrom: time.Date(now.Year()-1, 1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year()-1, 12, 31, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
	}
}
