package report

import (
	"encoding/csv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

type clearanceOrg struct {
	org       model.Org
	privLevel model.PrivLevel
	title     string
}
type clearanceRow struct {
	id              model.PersonID
	sortName        string
	orgs            []clearanceOrg
	volgistics      bool
	dswCERT         bool
	dswComm         bool
	bgCheckRecorded model.BGCheckType
	bgCheckAssumed  model.BGCheckType
	bgCheck         bool
	identification  model.IdentType
}

// GetClearance handles GET /api/reports/clearance requests.
func GetClearance(r *util.Request) error {
	var (
		params clearanceParameters
		data   []*clearanceRow
	)
	params = readClearanceParameters(r)
	if len(params.allowedRoles) == 0 {
		return util.Forbidden
	}
	data = getClearanceData(r, params)
	r.Tx.Commit()
	if params.renderCSV {
		renderClearanceCSV(r, data, params)
	} else {
		renderClearanceJSON(r, data, params)
	}
	return nil
}

func getClearanceData(r *util.Request, params clearanceParameters) (data []*clearanceRow) {
	var now = time.Now()

	for _, p := range r.Tx.FetchPeople() {
		if !clearancePersonMatch(p, params) {
			continue
		}
		if params.with != "" && !matchClearanceRestriction(p, params.with) {
			continue
		}
		if params.without != "" && matchClearanceRestriction(p, params.without) {
			continue
		}
		var row = clearanceRow{
			id:         p.ID,
			sortName:   p.SortName,
			volgistics: p.VolgisticsID != 0,
			dswCERT: p.DSWRegistrations != nil &&
				!p.DSWRegistrations[model.DSWCERT].IsZero() &&
				!p.DSWRegistrations[model.DSWCERT].After(now) &&
				(p.DSWUntil == nil || p.DSWUntil[model.DSWCERT].Before(now)),
			dswComm: p.DSWRegistrations != nil &&
				!p.DSWRegistrations[model.DSWComm].IsZero() &&
				!p.DSWRegistrations[model.DSWComm].After(now) &&
				(p.DSWUntil == nil || p.DSWUntil[model.DSWComm].Before(now)),
			identification: p.Identification,
		}
		// In calculating the orgs and privLevels to display, we don't
		// use p.Orgs because that includes implied roles.  For the
		// purpose of this report, we only want explicit roles.  The one
		// exception, which is coded explicitly below, is that we want
		// all leads to be shown as members of OrgAdmin.
		var orgs = map[model.Org]clearanceOrg{}
		for _, role := range r.Tx.FetchRoles() {
			if role.Title == "" || !p.Roles[role.ID] {
				continue
			}
			if orgs[role.Org].privLevel < role.PrivLevel {
				orgs[role.Org] = clearanceOrg{role.Org, role.PrivLevel, role.Title}
			}
			if role.PrivLevel >= model.PrivLeader && orgs[model.OrgAdmin].privLevel < model.PrivMember {
				orgs[model.OrgAdmin] = clearanceOrg{model.OrgAdmin, model.PrivMember, role.Title}
			}
		}
		for _, o := range model.AllOrgs {
			if orgs[o].privLevel != model.PrivNone {
				row.orgs = append(row.orgs, orgs[o])
			}
		}
		for _, bc := range p.BGChecks {
			if bc.Assumed {
				row.bgCheckAssumed |= bc.Type
			} else {
				row.bgCheckRecorded |= bc.Type
			}
		}
		row.bgCheckAssumed &^= row.bgCheckRecorded
		if p.Identification&model.IDCardKey != 0 {
			const needed = model.BGCheckDOJ | model.BGCheckFBI | model.BGCheckPHS
			row.bgCheck = (row.bgCheckAssumed|row.bgCheckRecorded)&needed == needed
		} else {
			const needed = model.BGCheckDOJ | model.BGCheckFBI
			row.bgCheck = (row.bgCheckAssumed|row.bgCheckRecorded)&needed == needed
		}
		data = append(data, &row)
	}
	return data
}

func renderClearanceCSV(r *util.Request, data []*clearanceRow, params clearanceParameters) {
	var (
		cols = []string{}
		out  = csv.NewWriter(r)
	)
	r.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Header().Set("Content-Disposition", `attachment; filename="clearance.csv"`)
	out.UseCRLF = true
	cols = append(cols, "Name")
	if r.Person.IsAdminLeader() {
		cols = append(cols, "Volgistics")
	} else {
		cols = append(cols, "Volunteer")
	}
	cols = append(cols, "DSW CERT", "DSW Comm", "BG Check", "Photo ID", "Card Key", "Green CERT LS Shirt", "Green CERT SS Shirt", "Tan SERV Shirt")
	out.Write(cols)
	for _, row := range data {
		cols = cols[:0]
		cols = append(cols, row.sortName, bool2CSV(row.volgistics), bool2CSV(row.dswCERT), bool2CSV(row.dswComm))
		if r.Person.IsAdminLeader() {
			var s []string
			for _, t := range model.AllBGCheckTypes {
				if row.bgCheckRecorded&t != 0 {
					s = append(s, t.String())
				} else if row.bgCheckAssumed&t != 0 {
					s = append(s, strings.ToLower(t.String()))
				}
			}
			cols = append(cols, strings.Join(s, " "))
		} else {
			cols = append(cols, bool2CSV(row.bgCheck))
		}
		cols = append(cols,
			bool2CSV(row.identification&model.IDPhoto != 0),
			bool2CSV(row.identification&model.IDCardKey != 0),
			bool2CSV(row.identification&model.IDCERTShirtLS != 0),
			bool2CSV(row.identification&model.IDCERTShirtSS != 0),
			bool2CSV(row.identification&model.IDSERVShirt != 0),
		)
		out.Write(cols)
	}
	out.Flush()
}
func bool2CSV(b bool) string {
	if b {
		return "X"
	}
	return ""
}

func renderClearanceJSON(r *util.Request, data []*clearanceRow, params clearanceParameters) {
	var out jwriter.Writer

	out.RawByte('{')
	clearanceRenderParams(r, &out, params)
	out.RawString(`,"rows":[`)
	for i, row := range data {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(row.id))
		out.RawString(`,"sortName":`)
		out.String(row.sortName)
		out.RawString(`,"orgs":{`)
		for i, o := range row.orgs {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(o.org.String())
			out.RawString(`:{"privLevel":`)
			out.String(o.privLevel.String())
			out.RawString(`,"title":`)
			out.String(o.title)
			out.RawByte('}')
		}
		out.RawString(`},"dswCERT":`)
		out.Bool(row.dswCERT)
		out.RawString(`,"dswComm":`)
		out.Bool(row.dswComm)
		out.RawString(`,"cardKey":`)
		out.Bool(row.identification&model.IDCardKey != 0)
		out.RawString(`,"certShirtLS":`)
		out.Bool(row.identification&model.IDCERTShirtLS != 0)
		out.RawString(`,"certShirtSS":`)
		out.Bool(row.identification&model.IDCERTShirtSS != 0)
		out.RawString(`,"idPhoto":`)
		out.Bool(row.identification&model.IDPhoto != 0)
		out.RawString(`,"servShirt":`)
		out.Bool(row.identification&model.IDSERVShirt != 0)
		out.RawString(`,"volgistics":`)
		out.Bool(row.volgistics)
		if r.Person.IsAdminLeader() {
			out.RawString(`,"bgCheckDOJ":`)
			out.String(renderBGCheck(row, model.BGCheckDOJ))
			out.RawString(`,"bgCheckFBI":`)
			out.String(renderBGCheck(row, model.BGCheckFBI))
			out.RawString(`,"bgCheckPHS":`)
			out.String(renderBGCheck(row, model.BGCheckPHS))
		} else {
			out.RawString(`,"bgCheck":`)
			out.Bool(row.bgCheck)
		}
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
}
func renderBGCheck(row *clearanceRow, bg model.BGCheckType) string {
	if row.bgCheckRecorded&bg != 0 {
		return "recorded"
	} else if row.bgCheckAssumed&bg != 0 {
		return "assumed"
	}
	return ""
}
