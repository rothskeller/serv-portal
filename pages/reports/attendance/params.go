package attrep

import (
	"time"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type parameters struct {
	dateRange     string
	dateFrom      string
	dateTo        string
	orgs          map[enum.Org]bool
	cells         string
	collapseX     bool
	collapseY     bool
	includeZerosY bool
	groupByOrg    bool
	renderCSV     bool
	// Not really parameters, but cached here for convenience:
	allowedOrgs map[enum.Org]bool
	dateRanges  []dateRange
}
type dateRange struct {
	tag      string
	label    string
	dateFrom string
	dateTo   string
}

func readParameters(r *request.Request, user *person.Person) (params parameters) {
	params.allowedOrgs = allowedOrgs(user)
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
	params.cells = r.FormValue("cells")
	if params.cells != "h" && params.cells != "a" && params.cells != "c" {
		params.cells = "h"
	}
	params.includeZerosY = r.FormValue("includeZerosY") != ""
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

	// Organizations
	params.orgs = make(map[enum.Org]bool)
	for _, oidstr := range r.Form["orgs"] {
		if o := enum.Org(util.ParseID(oidstr)); o.Valid() {
			params.orgs[o] = params.allowedOrgs[o]
		}
	}
	if len(params.orgs) == 0 { // No organizations selected, so include all allowed.
		params.orgs = params.allowedOrgs
	}
	return params
}

func renderParams(main *htmlb.Element, params parameters) {
	form := main.E("form class=attrepForm")
	grid := form.E("div class=attrepParams")
	// First, the date range box.
	box := grid.E("div class=attrepParamsBox")
	box.E("div class=attrepParamsBoxTitle>Date Range")
	sel := box.E("select id=attrepParamsDaterange name=dateRange")
	for _, dr := range params.dateRanges {
		sel.E("option value=%s data-from=%s data-to=%s", dr.tag, dr.dateFrom, dr.dateTo,
			params.dateRange == dr.tag, "selected").T(dr.label)
	}
	box.E("div id=attrepParamsDates>%s to\n%s", params.dateFrom, params.dateTo)
	// Next, the rows choice.
	box = grid.E("div class=attrepParamsBox")
	box.E("div class=attrepParamsBoxTitle>Rows")
	box.E("s-radio name=rows value=p label=Person", params.collapseY && !params.groupByOrg, "checked")
	box.E("s-radio name=rows value=o label=Org", params.collapseY && params.groupByOrg, "checked")
	box.E("s-radio name=rows value=po label='Person, Org'", !params.collapseY && !params.groupByOrg, "checked")
	box.E("s-radio name=rows value=op label='Org, Person'", !params.collapseY && params.groupByOrg, "checked")
	box.E("input type=checkbox name=includeZerosY class='s-check attrepParamsZeros' label='Include Zeros'", params.includeZerosY, "checked")
	// Next the columns choice.
	box = grid.E("div class=attrepParamsBox")
	box.E("div class=attrepParamsBoxTitle>Columns")
	box.E("s-radio name=columns value=e label=Dates", !params.collapseX, "checked")
	box.E("s-radio name=columns value=m label=Months", params.collapseX, "checked")
	// Next, the cells choice.
	box = grid.E("div class=attrepParamsBox")
	box.E("div class=attrepParamsBoxTitle>Cells")
	box.E("s-radio name=cells value=h label=Hours", params.cells != "a" && params.cells != "c", "checked")
	box.E("s-radio name=cells value=a label=Sign-Ins", params.cells == "a", "checked")
	box.E("s-radio name=cells value=c label=Credits", params.cells == "c", "checked")
	// Next, the organizations choice.
	box = grid.E("div class=attrepParamsBox")
	box.E("div class=attrepParamsBoxTitle>Orgs")
	for _, org := range enum.AllOrgs {
		if params.allowedOrgs[org] {
			box.E("div").E("input type=checkbox class=s-check name=orgs value=%d label=%s", org, org.Label(),
				params.orgs[org] || len(params.orgs) == 0, "checked")
		}
	}
}

// allowedOrgs returns a map of which organizations the current caller has
// reporting capabilities on.
func allowedOrgs(user *person.Person) (orgs map[enum.Org]bool) {
	orgs = make(map[enum.Org]bool)
	for _, o := range enum.AllOrgs {
		if user.HasPrivLevel(o, enum.PrivLeader) {
			orgs[o] = true
		}
	}
	return orgs
}

// dateRanges generates the set of date ranges supported by the reports.
func dateRanges() (ranges []dateRange) {
	now := time.Now()
	ranges = []dateRange{
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
			tag:      "l3",
			label:    "last three months",
			dateFrom: time.Date(now.Year(), now.Month()-3, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
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
		{
			tag:      "l12",
			label:    "last 12 months",
			dateFrom: time.Date(now.Year()-1, now.Month(), 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			dateTo:   time.Date(now.Year(), now.Month(), 0, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
		},
	}
	if now.Month() >= time.July {
		ranges = append(ranges, []dateRange{
			{
				tag:      "tf",
				label:    "this fiscal year",
				dateFrom: time.Date(now.Year(), 7, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				dateTo:   time.Date(now.Year()+1, 6, 30, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			},
			{
				tag:      "lf",
				label:    "last fiscal year",
				dateFrom: time.Date(now.Year()-1, 7, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				dateTo:   time.Date(now.Year(), 6, 30, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			},
		}...)
	} else {
		ranges = append(ranges, []dateRange{
			{
				tag:      "tf",
				label:    "this fiscal year",
				dateFrom: time.Date(now.Year()-1, 7, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				dateTo:   time.Date(now.Year(), 6, 30, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			},
			{
				tag:      "lf",
				label:    "last fiscal year",
				dateFrom: time.Date(now.Year()-2, 7, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
				dateTo:   time.Date(now.Year()-1, 6, 30, 0, 0, 0, 0, time.Local).Format("2006-01-02"),
			},
		}...)
	}
	return ranges
}
