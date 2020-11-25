package report

import (
	"encoding/csv"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

type attrepDataKey struct {
	pid model.PersonID
	org model.Org
}
type attrepColumn struct {
	ctype string
	label string
	etype string
}
type attrepRow struct {
	rtype  string
	label1 string
	label2 string
	data   []int
}

// GetAttendance handles GET /api/reports/attendance requests.
func GetAttendance(r *util.Request) error {
	var (
		params  attrepParameters
		events  []*model.Event
		data    map[attrepDataKey][]int
		columns []attrepColumn
		rows    []attrepRow
	)
	params = readAttrepParameters(r)
	if len(params.allowedOrgs) == 0 {
		return util.Forbidden
	}
	events = getAttrepEvents(r, params)
	data = getAttrepData(r, events, params)
	if params.collapseX {
		columns = makeAttrepMonthColumns(events, data, params)
	} else {
		columns = makeAttrepEventColumns(events, data, params)
	}
	rows = makeAttrepRows(r, columns, data, params)
	if !params.collapseX && len(params.orgs) > 1 {
		columns = mergeAttrepHours(rows, columns)
	}
	r.Tx.Commit()
	if params.renderCSV {
		renderAttrepCSV(r, rows, columns, params)
	} else {
		renderAttrepJSON(r, rows, columns, params)
	}
	return nil
}

func getAttrepEvents(r *util.Request, params attrepParameters) (events []*model.Event) {
	for _, e := range r.Tx.FetchEvents(params.dateFrom, params.dateTo) {
		if !params.eventTypes[e.Type] {
			continue
		}
		if !params.orgs[e.Org] {
			continue
		}
		events = append(events, e)
	}
	return events
}

func getAttrepData(r *util.Request, events []*model.Event, params attrepParameters) (data map[attrepDataKey][]int) {
	data = make(map[attrepDataKey][]int)
	for idx, e := range events {
		for pid, att := range r.Tx.FetchAttendanceByEvent(e) {
			var value = 1
			var key = attrepDataKey{pid: pid, org: e.Org}

			if !params.attendanceTypes[att.Type] && e.Type != model.EventHours {
				continue
			}
			if params.sumHours {
				if att.Minutes == 0 {
					continue
				}
				value = int(att.Minutes)
			}
			if data[key] == nil {
				data[key] = make([]int, len(events))
			}
			data[key][idx] = value
		}
	}
	return data
}

func makeAttrepMonthColumns(events []*model.Event, data map[attrepDataKey][]int, params attrepParameters) (columns []attrepColumn) {
	var monthIndex = make(map[string]int)
	if params.includeZerosX {
		timeFrom, _ := time.ParseInLocation("2006-01-02", params.dateFrom, time.Local)
		timeTo, _ := time.ParseInLocation("2006-01-02", params.dateTo, time.Local)
		timeFrom = time.Date(timeFrom.Year(), timeFrom.Month(), 1, 0, 0, 0, 0, time.Local)
		for timeFrom.Before(timeTo) {
			var month = timeFrom.Format("2006-01")
			monthIndex[month] = len(columns)
			columns = append(columns, attrepColumn{ctype: "c", label: month})
			timeFrom = timeFrom.AddDate(0, 1, 0)
		}
	} else {
		events = removeZeroEvents(events, data)
		for _, e := range events {
			var month = e.Date[0:7]
			if len(columns) == 0 || columns[len(columns)-1].label != month {
				monthIndex[month] = len(columns)
				columns = append(columns, attrepColumn{ctype: "c", label: month})
			}
		}
	}
	if len(columns) > 0 {
		columns[0].ctype = "s"
	}
	var newdata = make(map[attrepDataKey][]int, len(data))
	for key := range data {
		newdata[key] = make([]int, len(columns))
	}
	for oi, e := range events {
		var month = e.Date[0:7]
		var ni = monthIndex[month]
		for key := range data {
			newdata[key][ni] += data[key][oi]
		}
	}
	for key := range data {
		data[key] = newdata[key]
	}
	return addTotalColumn(columns, data)
}

func makeAttrepEventColumns(events []*model.Event, data map[attrepDataKey][]int, params attrepParameters) (columns []attrepColumn) {
	var lastmonth = ""
	if !params.includeZerosX {
		events = removeZeroEvents(events, data)
	}
	for _, e := range events {
		var ctype = "c"
		var label = e.Date
		var etype = model.EventTypeNames[e.Type][:1]
		if e.Date[:7] != lastmonth {
			ctype = "s"
			lastmonth = e.Date[:7]
		}
		if e.Type == model.EventHours {
			label = e.Date[:7] + "-??"
			etype = "?"
		}
		columns = append(columns, attrepColumn{
			ctype: ctype,
			label: label,
			etype: etype,
		})
	}
	return addTotalColumn(columns, data)
}

func removeZeroEvents(events []*model.Event, data map[attrepDataKey][]int) []*model.Event {
	var j = 0
	for i := range events {
		var nonzero = false
		for _, v := range data {
			if v[i] != 0 {
				nonzero = true
				break
			}
		}
		if nonzero {
			events[j] = events[i]
			for _, v := range data {
				v[j] = v[i]
			}
			j++
		}
	}
	for key := range data {
		data[key] = data[key][:j]
	}
	return events[:j]
}

func addTotalColumn(columns []attrepColumn, data map[attrepDataKey][]int) []attrepColumn {
	if len(columns) == 0 {
		return columns
	}
	if len(columns) == 1 {
		columns[0].ctype = "1"
		return columns
	}
	if len(columns) < 2 {
		return columns
	}
	for key := range data {
		var total = 0
		for _, v := range data[key] {
			total += v
		}
		data[key] = append(data[key], total)
	}
	return append(columns, attrepColumn{ctype: "t", label: "TOTAL"})
}

func makeAttrepRows(r *util.Request, columns []attrepColumn, data map[attrepDataKey][]int, params attrepParameters) (rows []attrepRow) {
	var pnames = make(map[model.PersonID]string)

	// First, handle includeZeros.  If it's true, we need to add all of the
	// people who didn't have attendance.  If it's false, we need to remove
	// rows that have only zero values.
	if params.includeZerosY {
		for _, p := range r.Tx.FetchPeople() {
			for o := range params.orgs {
				if p.Orgs[o].PrivLevel >= model.PrivMember2 {
					key := attrepDataKey{pid: p.ID, org: o}
					if data[key] == nil {
						data[key] = make([]int, len(columns))
					}
				}
			}
		}
	} else {
		for key := range data {
			var nonzero = false
			for _, v := range data[key] {
				if v != 0 {
					nonzero = true
					break
				}
			}
			if !nonzero {
				delete(data, key)
			}
		}
	}
	// Now, generate the rows, in random order.
	switch {
	case params.collapseY && params.groupByOrg:
		var index = make(map[model.Org]int)
		for key, vals := range data {
			idx, ok := index[key.org]
			if !ok {
				idx = len(rows)
				index[key.org] = idx
				rows = append(rows, attrepRow{
					label1: orgNames[key.org],
					data:   make([]int, len(columns)),
				})
			}
			addValues(rows[idx].data, vals)
		}
	case !params.collapseY && params.groupByOrg:
		for key, vals := range data {
			if pnames[key.pid] == "" {
				pnames[key.pid] = r.Tx.FetchPerson(key.pid).SortName
			}
			rows = append(rows, attrepRow{
				label1: orgNames[key.org],
				label2: pnames[key.pid],
				data:   vals,
			})
		}
	case params.collapseY && !params.groupByOrg:
		var index = make(map[model.PersonID]int)
		for key, vals := range data {
			idx, ok := index[key.pid]
			if !ok {
				idx = len(rows)
				index[key.pid] = idx
				rows = append(rows, attrepRow{
					label1: r.Tx.FetchPerson(key.pid).SortName,
					data:   make([]int, len(columns)),
				})
			}
			addValues(rows[idx].data, vals)
		}
	case !params.collapseY && !params.groupByOrg:
		for key, vals := range data {
			if pnames[key.pid] == "" {
				pnames[key.pid] = r.Tx.FetchPerson(key.pid).SortName
			}
			rows = append(rows, attrepRow{
				label1: pnames[key.pid],
				label2: orgNames[key.org],
				data:   vals,
			})
		}
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
				rows[i].rtype = "s"
			} else if i%2 == 1 {
				rows[i].rtype = "c2"
			} else {
				rows[i].rtype = "c"
			}
		}
	case !params.collapseY && params.groupByOrg:
		var colored bool
		var last string
		for i, r := range rows {
			if r.label1 != last {
				last = r.label1
				rows[i].rtype = "s"
				colored = false
			} else if colored {
				rows[i].rtype = "c"
				colored = false
			} else {
				rows[i].rtype = "c2"
				colored = true
			}
		}
	case !params.collapseY && !params.groupByOrg:
		var colored bool
		var last string
		for i, r := range rows {
			if i == 0 {
				rows[i].rtype = "s"
				last = r.label1
				colored = false
			} else {
				if r.label1 != last {
					colored = !colored
					last = r.label1
				}
				if colored {
					rows[i].rtype = "c2"
				} else {
					rows[i].rtype = "c"
				}
			}
		}
	}
	// Add a total row.
	if len(rows) > 1 {
		var totals = make([]int, len(columns))
		for _, row := range rows {
			addValues(totals, row.data)
		}
		rows = append(rows, attrepRow{rtype: "t", label1: "TOTAL", data: totals})
	}
	// Add total rows for each grouping, and remove redundant labels.
	if !params.collapseY {
		var nr = make([]attrepRow, 0, len(rows))
		var totals []int
		var count = 0
		var last = ""
		var lastRType = ""
		for _, row := range rows {
			if row.label1 != last {
				if count > 1 {
					var rtype string
					if params.groupByOrg && lastRType == "c2" {
						rtype = "t2"
					} else if params.groupByOrg || lastRType == "c2" {
						rtype = "tc2"
					} else {
						rtype = "t2"
					}
					nr = append(nr, attrepRow{
						rtype:  rtype,
						label1: "",
						label2: "TOTAL",
						data:   totals,
					})
				}
				count = 0
				last = row.label1
				lastRType = row.rtype
				totals = make([]int, len(columns))
			} else {
				row.label1 = ""
				lastRType = row.rtype
			}
			addValues(totals, row.data)
			count++
			nr = append(nr, row)
		}
		if count > 1 {
			nr = append(nr, attrepRow{
				rtype:  "t2",
				label1: "",
				label2: "TOTAL",
				data:   totals,
			})
		}
		rows = nr
	}
	return rows
}

func addValues(to, from []int) {
	for idx := range from {
		to[idx] += from[idx]
	}
}

func mergeAttrepHours(rows []attrepRow, columns []attrepColumn) []attrepColumn {
	// If we're displaying multiple orgs, we may have multiple column
	// headings for "other hours".  That's ugly, so we'll merge them.
	j := 0
	last := ""
	for i := range columns {
		if columns[i].etype == "?" && columns[i].label == last {
			for r := range rows {
				rows[r].data[j-1] += rows[r].data[i]
			}
		} else {
			last = columns[i].label
			columns[j] = columns[i]
			for r := range rows {
				rows[r].data[j] = rows[r].data[i]
			}
			j++
		}
	}
	for r := range rows {
		rows[r].data = rows[r].data[:j]
	}
	return columns[:j]
}

func renderAttrepCSV(r *util.Request, rows []attrepRow, columns []attrepColumn, params attrepParameters) {
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
			if params.sumHours {
				cols = append(cols, fmt.Sprintf("%.1f", float64(v)/60.0))
			} else {
				cols = append(cols, strconv.Itoa(v))
			}
		}
		out.Write(cols)
	}
	out.Flush()
}

func renderAttrepJSON(r *util.Request, rows []attrepRow, columns []attrepColumn, params attrepParameters) {
	var out jwriter.Writer

	out.RawByte('{')
	attrepRenderParams(&out, params)
	if params.collapseY {
		out.RawString(`,"columns":["h"`)
	} else {
		out.RawString(`,"columns":["h","h2"`)
	}
	for _, col := range columns {
		out.RawByte(',')
		out.String(col.ctype)
	}
	if params.collapseX {
		out.RawString(`],"rows":["h"`)
	} else {
		out.RawString(`],"rows":["h","h2"`)
	}
	for _, row := range rows {
		out.RawByte(',')
		out.String(row.rtype)
	}
	out.RawString(`],"cells":[[""`)
	if !params.collapseY {
		out.RawString(`,""`)
	}
	for _, col := range columns {
		out.RawByte(',')
		out.String(col.label)
	}
	if !params.collapseX {
		out.RawString(`],[""`)
		if !params.collapseY {
			out.RawString(`,""`)
		}
		for _, col := range columns {
			out.RawByte(',')
			out.String(col.etype)
		}
	}
	out.RawByte(']')
	for _, row := range rows {
		out.RawString(`,[`)
		out.String(row.label1)
		if !params.collapseY {
			out.RawByte(',')
			out.String(row.label2)
		}
		for _, val := range row.data {
			out.RawByte(',')
			if val == 0 {
				out.RawString(`""`)
			} else if params.sumHours {
				out.String(fmt.Sprintf("%.1f", float64(val)/60.0))
			} else {
				out.IntStr(val)
			}
		}
		out.RawByte(']')
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
}

var orgNames = map[model.Org]string{
	model.OrgAdmin2: "Admin",
	model.OrgCERTD2: "CERT-D",
	model.OrgCERTT2: "CERT-T",
	model.OrgListos: "Listos",
	model.OrgSARES2: "SARES",
	model.OrgSNAP2:  "SNAP",
}
