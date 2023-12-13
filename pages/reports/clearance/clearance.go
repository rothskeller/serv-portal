package clearrep

import (
	"encoding/csv"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type orgdata struct {
	org       enum.Org
	privLevel enum.PrivLevel
	title     string
}
type rowdata struct {
	id             person.ID
	sortName       string
	orgs           []orgdata
	volgistics     bool
	dswCERT        bool
	dswComm        bool
	bgDOJRecorded  bool
	bgDOJAssumed   bool
	bgFBIRecorded  bool
	bgFBIAssumed   bool
	bgPHSRecorded  bool
	bgPHSAssumed   bool
	bgCheck        bool
	identification person.IdentType
}

// Get handles GET /reports/clearance requests.
func Get(r *request.Request) {
	var (
		user   *person.Person
		params parameters
		data   []*rowdata
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if params = readParameters(r, user); len(params.allowedRoles) == 0 {
		errpage.Forbidden(r, user)
		return
	}
	data = getData(r, params)
	if params.renderCSV {
		renderCSV(r, user, data, params)
		return
	}
	ui.Page(r, user, ui.PageOpts{
		Title:    "Clearance",
		Banner:   "Clearance Report",
		MenuItem: "reports",
		Tabs: []ui.PageTab{
			{Name: "Attendance", URL: "/reports/attendance", Alias: "/reports/attendance?*", Target: "main"},
			{Name: "Clearance", URL: "/reports/clearance", Alias: "/reports/clearance?*", Target: "main", Active: true},
		},
	}, func(e *htmlb.Element) {
		renderReport(e, user, data, params)
	})
}

func getData(r *request.Request, params parameters) (data []*rowdata) {
	const personFields = person.FID | person.FBGChecks | person.FDSWRegistrations | person.FIdentification | person.FSortName | person.FVolgisticsID

	person.All(r, personFields, func(p *person.Person) {
		if !personMatch(r, p, params) {
			return
		}
		if params.with != "" && !matchRestriction(p, params.with) {
			return
		}
		if params.without != "" && matchRestriction(p, params.without) {
			return
		}
		var row = rowdata{
			id:             p.ID(),
			sortName:       p.SortName(),
			volgistics:     p.VolgisticsID() != 0,
			dswCERT:        hasDSW(p.DSWRegistrations().CERT),
			dswComm:        hasDSW(p.DSWRegistrations().Communications),
			identification: p.Identification(),
		}
		// In calculating the orgs and privLevels to display, we don't
		// use privlevels because that includes implied roles.  For the
		// purpose of this report, we only want explicit roles.  The one
		// exception, which is coded explicitly below, is that we want
		// all leads to be shown as members of OrgAdmin.
		var orgs = map[enum.Org]orgdata{}
		role.All(r, role.FID|role.FOrg|role.FPrivLevel|role.FTitle, func(rl *role.Role) {
			if rl.Title() == "" {
				return
			}
			if held, _ := personrole.PersonHasRole(r, p.ID(), rl.ID()); !held {
				return
			}
			if orgs[rl.Org()].privLevel < rl.PrivLevel() {
				orgs[rl.Org()] = orgdata{rl.Org(), rl.PrivLevel(), rl.Title()}
			}
			if rl.PrivLevel() >= enum.PrivLeader && orgs[enum.OrgAdmin].privLevel < enum.PrivMember {
				orgs[enum.OrgAdmin] = orgdata{enum.OrgAdmin, enum.PrivMember, rl.Title()}
			}
		})
		for _, o := range enum.AllOrgs {
			if orgs[o].privLevel != 0 {
				row.orgs = append(row.orgs, orgs[o])
			}
		}
		if hasBGCheck(p.BGChecks().DOJ, false) {
			row.bgDOJRecorded = true
		} else if hasBGCheck(p.BGChecks().DOJ, true) {
			row.bgDOJAssumed = true
		}
		if hasBGCheck(p.BGChecks().FBI, false) {
			row.bgFBIRecorded = true
		} else if hasBGCheck(p.BGChecks().FBI, true) {
			row.bgFBIAssumed = true
		}
		if hasBGCheck(p.BGChecks().PHS, false) {
			row.bgPHSRecorded = true
		} else if hasBGCheck(p.BGChecks().PHS, true) {
			row.bgPHSAssumed = true
		}
		row.bgCheck = (row.bgDOJRecorded || row.bgDOJAssumed) && (row.bgFBIRecorded || row.bgFBIAssumed)
		if p.Identification()&person.IDCardKey != 0 && !row.bgPHSRecorded && !row.bgPHSAssumed {
			row.bgCheck = false
		}
		data = append(data, &row)
	})
	return data
}

func renderCSV(r *request.Request, user *person.Person, data []*rowdata, params parameters) {
	var (
		cols = []string{}
		out  = csv.NewWriter(r)
	)
	r.Header().Set("Content-Type", "text/csv; charset=utf-8")
	r.Header().Set("Content-Disposition", `attachment; filename="clearance.csv"`)
	out.UseCRLF = true
	cols = append(cols, "Name")
	if user.IsAdminLeader() {
		cols = append(cols, "Volgistics")
	} else {
		cols = append(cols, "Volunteer")
	}
	cols = append(cols, "DSW CERT", "DSW Comm", "BG Check", "Photo ID", "Card Key", "Green CERT LS Shirt", "Green CERT SS Shirt", "Tan SERV Shirt")
	out.Write(cols)
	for _, row := range data {
		cols = cols[:0]
		cols = append(cols, row.sortName, bool2CSV(row.volgistics), bool2CSV(row.dswCERT), bool2CSV(row.dswComm))
		if user.IsAdminLeader() {
			var s string
			switch {
			case row.bgDOJRecorded:
				s = "D"
			case row.bgDOJAssumed:
				s = "d"
			default:
				s = " "
			}
			switch {
			case row.bgFBIRecorded:
				s += " F"
			case row.bgFBIAssumed:
				s += " f"
			default:
				s += "  "
			}
			switch {
			case row.bgPHSRecorded:
				s += " P"
			case row.bgPHSAssumed:
				s += " p"
			default:
				s += "  "
			}
			cols = append(cols, s)
		} else {
			cols = append(cols, bool2CSV(row.bgCheck))
		}
		cols = append(cols,
			bool2CSV(row.identification&person.IDPhoto != 0),
			bool2CSV(row.identification&person.IDCardKey != 0),
			bool2CSV(row.identification&person.IDCERTShirtLS != 0),
			bool2CSV(row.identification&person.IDCERTShirtSS != 0),
			bool2CSV(row.identification&person.IDSERVShirt != 0),
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

func renderReport(main *htmlb.Element, user *person.Person, data []*rowdata, params parameters) {
	hasBGCheckDetail := user.IsAdminLeader()
	renderParams(main, params)
	renderTable(main, data, hasBGCheckDetail)
	renderCount(main, data)
	renderCSVButton(main, data)
}
func renderTable(main *htmlb.Element, data []*rowdata, hasBGCheckDetail bool) {
	if len(data) == 0 {
		main.E("div class=clearrep-noData>No one matches these report criteria.")
		return
	}
	var omap = map[enum.Org]struct{}{}
	for _, row := range data {
		for _, org := range row.orgs {
			omap[org.org] = struct{}{}
		}
	}
	var orgs []enum.Org
	for _, org := range enum.AllOrgs {
		if _, ok := omap[org]; ok {
			orgs = append(orgs, org)
		}
	}
	table := main.E("div class=clearrepTable")
	renderTableHeading(table, hasBGCheckDetail)
	for _, p := range data {
		renderTableRow(table, p, orgs, hasBGCheckDetail)
	}
}
func renderTableHeading(table *htmlb.Element, hasBGCheckDetail bool) {
	table.E("div class=clearrepHeading>Orgs")
	table.E("div class=clearrepHeading>Name")
	table.E("div class=clearrepHeading>V")
	table.E("div class=clearrepHeading>DSW")
	if hasBGCheckDetail {
		table.E("div class=clearrepHeading>BG")
	} else {
		table.E("div class=clearrepHeading>B")
	}
	table.E("div class=clearrepHeading>Identification")
}
func renderTableRow(table *htmlb.Element, p *rowdata, orgs []enum.Org, hasBGCheckDetail bool) {
	renderOrgBadgeCells(table, p, orgs)
	renderNameLinkCell(table, p)
	renderVolgisticsCell(table, p)
	renderDSWCells(table, p)
	if hasBGCheckDetail {
		renderBGCheckCells(table, p)
	} else {
		renderBGCheckCell(table, p)
	}
	renderIdentCells(table, p)
}
func renderOrgBadgeCells(table *htmlb.Element, p *rowdata, orgs []enum.Org) {
	div := table.E("div class=clearrepBoxes")
	for _, org := range orgs {
		var found bool
		for _, po := range p.orgs {
			if po.org == org {
				div.E("div class=clearrepOrg%d%d title=%s>%s", org, po.privLevel, po.title, orgBadgeLabels[org])
				found = true
				break
			}
		}
		if !found {
			div.E("div class=clearrepOrgPH")
		}
	}
}
func renderNameLinkCell(table *htmlb.Element, p *rowdata) {
	table.E("a href=/people/%d up-target=.pageCanvas", p.id).T(p.sortName)
}
func renderVolgisticsCell(table *htmlb.Element, p *rowdata) {
	if p.volgistics {
		table.E("div class=clearrepVolgistics title='City Volunteer'>V")
	} else {
		table.E("div class=clearrepVolgistics")
	}
}
func renderDSWCells(table *htmlb.Element, p *rowdata) {
	div := table.E("div class=clearrepBoxes")
	if p.dswCERT {
		div.E("div class=clearrepDSWCERT title='DSW for CERT'>C")
	} else {
		div.E("div class=clearrepDSWCERT")
	}
	if p.dswComm {
		div.E("div class=clearrepDSWComm title='DSW for Communications'>S")
	} else {
		div.E("div class=clearrepDSWComm")
	}
}
func renderBGCheckCell(table *htmlb.Element, p *rowdata) {
	if p.bgCheck {
		table.E("div class=clearrepBGCheck title='Background Check'>B")
	} else {
		table.E("div class=clearrepBGCheck")
	}
}
func renderBGCheckCells(table *htmlb.Element, p *rowdata) {
	div := table.E("div class=clearrepBoxes")
	if p.bgDOJRecorded {
		div.E("div class=clearrepBGDOJ-recorded title='LiveScan/DOJ'>D")
	} else if p.bgDOJAssumed {
		div.E("div class=clearrepBGDOJ-assumed title='LiveScan/DOJ (assumed)'>D")
	} else {
		div.E("div class=clearrepBGDOJ")
	}
	if p.bgFBIRecorded {
		div.E("div class=clearrepBGFBI-recorded title='LiveScan/FBI'>F")
	} else if p.bgFBIAssumed {
		div.E("div class=clearrepBGFBI-assumed title='LiveScan/FBI (assumed)'>F")
	} else {
		div.E("div class=clearrepBGFBI")
	}
	if p.bgPHSRecorded {
		div.E("div class=clearrepBGPHS-recorded title='Personal History'>P")
	} else if p.bgPHSAssumed {
		div.E("div class=clearrepBGPHS-assumed title='Personal History (assumed)'>P")
	} else {
		div.E("div class=clearrepBGPHS")
	}
}
func renderIdentCells(table *htmlb.Element, p *rowdata) {
	div := table.E("div class=clearrepBoxes")
	if p.identification&person.IDPhoto != 0 {
		div.E("div class=clearrepIDPhoto title='Photo ID'>P")
	} else {
		div.E("div class=clearrepIDPhoto")
	}
	if p.identification&person.IDCardKey != 0 {
		div.E("div class=clearrepIDCardkey title='Card Key'>C")
	} else {
		div.E("div class=clearrepIDCardkey")
	}
	if p.identification&person.IDCERTShirtLS != 0 {
		div.E("div class=clearrepIDCERTShirtls title='Green CERT Shirt (LS)'>S")
	} else {
		div.E("div class=clearrepIDCERTShirtls")
	}
	if p.identification&person.IDCERTShirtSS != 0 {
		div.E("div class=clearrepIDCERTShirtss title='Green CERT Shirt (SS)'>S")
	} else {
		div.E("div class=clearrepIDCERTShirtss")
	}
	if p.identification&person.IDSERVShirt != 0 {
		div.E("div class=clearrepIDSERVShirt title='Tan SERV Shirt'>S")
	} else {
		div.E("div class=clearrepIDSERVShirt")
	}
}
func renderCount(main *htmlb.Element, data []*rowdata) {
	switch len(data) {
	case 0:
		break
	case 1:
		main.E("div class=clearrepCount>1 person listed")
	default:
		main.E("div class=clearrepCount>%d people listed", len(data))
	}
}
func renderCSVButton(main *htmlb.Element, data []*rowdata) {
	if len(data) != 0 {
		main.E("div class=clearrepButtons").
			E("button type=button id=clearrepExport class='sbtn sbtn-primary'>Export")
	}
}

var orgBadgeLabels = map[enum.Org]string{
	enum.OrgAdmin:  "A",
	enum.OrgCERTD:  "D",
	enum.OrgCERTT:  "T",
	enum.OrgListos: "L",
	enum.OrgSARES:  "S",
	enum.OrgSNAP:   "S",
}
