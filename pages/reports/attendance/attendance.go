package attrep

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type dataKey struct {
	pid person.ID
	org enum.Org
}
type column struct {
	ctype string
	label string
	etype string
}
type row struct {
	rtype  string
	label1 string
	label2 string
	data   []int
}

// Get handles GET /reports/attendance requests.
func Get(r *request.Request) {
	var (
		user    *person.Person
		params  parameters
		data    map[dataKey]map[string]int
		columns []column
		rows    []row
		pcount  int
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if params = readParameters(r, user); len(params.allowedOrgs) == 0 {
		errpage.Forbidden(r, user)
		return
	}
	data = getData(r, params)
	columns = makeColumns(data, params)
	rows, pcount = makeRows(r, columns, data, params)
	if params.renderCSV {
		renderCSV(r, rows, columns, params)
		return
	}
	ui.Page(r, user, ui.PageOpts{
		Title:    "Attendance",
		Banner:   "Attendance Report",
		MenuItem: "reports",
		Tabs: []ui.PageTab{
			{Name: "Attendance", URL: "/reports/attendance", Alias: "/reports/attendance?*", Target: "main", Active: true},
			{Name: "Clearance", URL: "/reports/clearance", Alias: "/reports/clearance?*", Target: "main"},
		},
	}, func(e *htmlb.Element) {
		e.Attr("class=attrep")
		renderParams(e, params)
		renderReport(e, rows, columns, params, pcount)
	})
}

func getData(r *request.Request, params parameters) (data map[dataKey]map[string]int) {
	data = make(map[dataKey]map[string]int)
	// Create rows for every person/org if includeZeros is set.
	if params.includeZerosY {
		if params.collapseY && params.groupByOrg {
			for org := range params.orgs {
				data[dataKey{org: org}] = nil
			}
		} else {
			priv := enum.PrivMember
			if params.cells != "h" {
				priv = enum.PrivStudent
			}
			person.All(r, person.FID|person.FPrivLevels, func(p *person.Person) {
				for o := range params.orgs {
					if p.HasPrivLevel(o, priv) {
						if params.collapseY {
							data[dataKey{pid: p.ID()}] = nil
						} else {
							data[dataKey{pid: p.ID(), org: o}] = nil
						}
					}
				}
			})
		}
	}
	// Load the requested data.
	event.AllBetween(r, params.dateFrom, params.dateTo+"Z", event.FID|event.FStart|event.FFlags, 0, func(e *event.Event, _ *venue.Venue) {
		// Adding a "Z" after dateTo allows it to include events on the
		// final date.  (Alphabetic "<" comparison.)
		var date = e.Start()[:10]
		if params.collapseX {
			date = date[:7]
		} else if e.Flags()&event.OtherHours != 0 {
			date = date[:8] + "??"
		}
		task.AllForEvent(r, e.ID(), task.FID|task.FOrg, func(t *task.Task) {
			if !params.orgs[t.Org()] {
				return
			}
			taskperson.PeopleForTask(r, t.ID(), person.FID, func(p *person.Person, minutes uint, flags taskperson.Flag) {
				switch params.cells {
				case "h":
					if minutes == 0 {
						return
					}
				case "c":
					if flags&taskperson.Credited == 0 {
						return
					}
				case "a":
					if flags&taskperson.Attended == 0 {
						return
					}
				}
				var key dataKey
				if params.groupByOrg || !params.collapseY {
					key.org = t.Org()
				}
				if !params.groupByOrg || !params.collapseY {
					key.pid = p.ID()
				}
				if data[key] == nil {
					data[key] = make(map[string]int)
				}
				if params.cells == "h" {
					data[key][date] += int(minutes)
				} else {
					data[key][date]++
				}
			})
		})

	})
	return data
}

func makeColumns(data map[dataKey]map[string]int, params parameters) (columns []column) {
	var dates = make(map[string]bool)
	for _, d := range data {
		for date := range d {
			dates[date] = true
		}
	}
	if len(dates) == 0 {
		return nil
	}
	var datelist = make([]string, 0, len(dates))
	for date := range dates {
		datelist = append(datelist, date)
	}
	sort.Strings(datelist)
	for _, date := range datelist {
		columns = append(columns, column{ctype: "C", label: date})
	}
	for i := range columns {
		if i == 0 || columns[i-1].label[:7] != columns[i].label[:7] {
			columns[i].ctype = "S"
		}
	}
	if len(columns) == 1 {
		columns[0].ctype = "1"
		return columns
	}
	for key := range data {
		var total = 0
		for _, v := range data[key] {
			total += v
		}
		if data[key] == nil {
			data[key] = make(map[string]int)
		}
		data[key]["TOTAL"] = total
	}
	return append(columns, column{ctype: "T", label: "TOTAL"})
}

func makeRows(r *request.Request, columns []column, data map[dataKey]map[string]int, params parameters) (rows []row, pcount int) {
	var pnames = make(map[person.ID]string)

	// Generate the rows, in random order.
	for key := range data {
		if key.pid != 0 && pnames[key.pid] == "" {
			pnames[key.pid] = person.WithID(r, key.pid, person.FSortName).SortName()
		}
		var r row
		switch {
		case params.collapseY && params.groupByOrg:
			r.label1 = orgNames[key.org]
		case !params.collapseY && params.groupByOrg:
			r.label1 = orgNames[key.org]
			r.label2 = pnames[key.pid]
		case params.collapseY && !params.groupByOrg:
			r.label1 = pnames[key.pid]
		case !params.collapseY && !params.groupByOrg:
			r.label1 = pnames[key.pid]
			r.label2 = orgNames[key.org]
		}
		for _, col := range columns {
			r.data = append(r.data, data[key][col.label])
		}
		rows = append(rows, r)
	}
	// Now, sort the rows.
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].label1 < rows[j].label1 {
			return true
		}
		if rows[i].label1 > rows[j].label1 {
			return false
		}
		return rows[i].label2 < rows[j].label2
	})
	// Set the row types.
	switch {
	case params.collapseY:
		for i := range rows {
			if i == 0 {
				rows[i].rtype = "S"
			} else if i%2 == 1 {
				rows[i].rtype = "C2"
			} else {
				rows[i].rtype = "C"
			}
		}
	case !params.collapseY && params.groupByOrg:
		var colored bool
		var last string
		for i, r := range rows {
			if r.label1 != last {
				last = r.label1
				rows[i].rtype = "S"
				colored = false
			} else if colored {
				rows[i].rtype = "C"
				colored = false
			} else {
				rows[i].rtype = "C2"
				colored = true
			}
		}
	case !params.collapseY && !params.groupByOrg:
		var colored bool
		var last string
		for i, r := range rows {
			if i == 0 {
				rows[i].rtype = "S"
				last = r.label1
				colored = false
			} else {
				if r.label1 != last {
					colored = !colored
					last = r.label1
				}
				if colored {
					rows[i].rtype = "C2"
				} else {
					rows[i].rtype = "C"
				}
			}
		}
	}
	pcount = len(rows)
	// Add a total row.
	if len(rows) > 1 {
		var totals = make([]int, len(columns))
		for _, row := range rows {
			addValues(totals, row.data)
		}
		rows = append(rows, row{rtype: "T", label1: "TOTAL", data: totals})
	}
	// Add total rows for each grouping, and remove redundant labels.
	if !params.collapseY {
		var nr = make([]row, 0, len(rows))
		var totals []int
		var count = 0
		var last = ""
		var lastRType = ""
		for _, rw := range rows {
			if rw.label1 != last {
				if count > 1 {
					var rtype string
					if params.groupByOrg && lastRType == "C2" {
						rtype = "T2"
					} else if params.groupByOrg || lastRType == "C2" {
						rtype = "Tc2"
					} else {
						rtype = "T2"
					}
					nr = append(nr, row{
						rtype:  rtype,
						label1: "",
						label2: "TOTAL",
						data:   totals,
					})
				}
				count = 0
				last = rw.label1
				lastRType = rw.rtype
				totals = make([]int, len(columns))
			} else {
				rw.label1 = ""
				lastRType = rw.rtype
				pcount--
			}
			addValues(totals, rw.data)
			count++
			nr = append(nr, rw)
		}
		if count > 1 {
			nr = append(nr, row{
				rtype:  "T2",
				label1: "",
				label2: "TOTAL",
				data:   totals,
			})
		}
		rows = nr
	}
	if params.groupByOrg {
		pcount = 0
	}
	return rows, pcount
}

func addValues(to, from []int) {
	for idx := range from {
		to[idx] += from[idx]
	}
}

func renderCSV(r *request.Request, rows []row, columns []column, params parameters) {
	var (
		cols = []string{}
		out  = csv.NewWriter(r)
	)
	r.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Header().Set("Content-Disposition", `attachment; filename="attendance.csv"`)
	out.UseCRLF = true
	cols = append(cols, "")
	if !params.collapseY {
		cols = append(cols, "")
	}
	for _, col := range columns {
		cols = append(cols, col.label)
	}
	out.Write(cols)
	if !params.collapseX {
		cols = cols[:1]
		if !params.collapseY {
			cols = append(cols, "")
		}
		for _, col := range columns {
			cols = append(cols, col.etype)
		}
		out.Write(cols)
	}
	for _, row := range rows {
		cols = cols[:0]
		cols = append(cols, row.label1)
		if !params.collapseY {
			cols = append(cols, row.label2)
		}
		for _, v := range row.data {
			if params.cells == "h" {
				cols = append(cols, fmt.Sprintf("%.1f", float64(v)/60.0))
			} else {
				cols = append(cols, strconv.Itoa(v))
			}
		}
		out.Write(cols)
	}
	out.Flush()
}

func renderReport(main *htmlb.Element, rows []row, columns []column, params parameters, pcount int) {
	table := main.E("table class=attrepTable")
	// First, write the top row with the column labels.
	tr := table.E("tr class=attrepRowH")
	tr.E("td class=attrepColH") // empty
	if !params.collapseY {
		tr.E("td class=attrepColH2") // empty
	}
	for _, col := range columns {
		td := tr.E("td class=attrepCol%s", col.ctype)
		if col.label != "" {
			td.E("div class=attrepCell-vertical").T(col.label)
		}
	}
	// Now, the data rows.
	for _, row := range rows {
		tr = table.E("tr class=attrepRow%s", row.rtype)
		tr.E("td class=attrepColH").T(row.label1)
		if !params.collapseY {
			tr.E("td class=attrepColH2").T(row.label2)
		}
		for ci, val := range row.data {
			td := tr.E("td class=attrepCol%s", columns[ci].ctype)
			if val == 0 && columns[ci].ctype[0] != 'T' {
				// nothing
			} else if params.cells == "h" {
				td.R(strconv.FormatFloat(float64(val)/60.0, 'f', 1, 64))
			} else {
				td.R(strconv.Itoa(val))
			}
		}
	}
	// Add a count below the table.
	if pcount == 1 {
		main.E("div class=attrepPcount>1 person listed")
	} else if pcount > 1 {
		main.E("div class=attrepPcount>%d people listed", pcount)
	}
	// Add an export button.
	if len(rows) != 0 {
		main.E("div class=attrepButtons").E("button type=button id=attrepExport class='sbtn sbtn-primary'>Export")
	}
}

var orgNames = map[enum.Org]string{
	enum.OrgAdmin:  "Admin",
	enum.OrgCERTD:  "CERT-D",
	enum.OrgCERTT:  "CERT-T",
	enum.OrgListos: "Listos",
	enum.OrgSARES:  "SARES",
	enum.OrgSNAP:   "SNAP",
}
