package report

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
)

func emitValidClearanceRestrictions(caller *model.Person, out *jwriter.Writer) {
	if caller.IsAdminLeader() {
		out.RawString(`[{"value":"bgDOJ-assumed","label":"BG Check: DOJ (assumed)"},{"value":"bgDOJ","label":"BG Check: DOJ (recorded)"},{"value":"bgFBI-assumed","label":"BG Check: FBI (assumed)"},{"value":"bgFBI","label":"BG Check: FBI (recorded)"},{"value":"bgPHS-assumed","label":"BG Check: PHS (assumed)"},{"value":"bgPHS","label":"BG Check: PHS (recorded)"},{"value":"cardKey","label":"DPS Card Key"},{"value":"idPhoto","label":"DPS Photo ID"},{"value":"dswCERT","label":"DSW for CERT"},{"value":"dswComm","label":"DSW for Communications"},{"value":"certShirtLS","label":"Green CERT Shirt (LS)"},{"value":"certShirtSS","label":"Green CERT Shirt (SS)"},{"value":"volgistics","label":"In Volgistics"},{"value":"servShirt","label":"Tan SERV Shirt"}]`)
	} else {
		out.RawString(`[{"value":"bgCheck","label":"BG Check"},{"value":"volgistics","label":"City Volunteer"},{"value":"cardKey","label":"DPS Card Key"},{"value":"idPhoto","label":"DPS Photo ID"},{"value":"dswCERT","label":"DSW for CERT"},{"value":"dswComm","label":"DSW for Communications"},{"value":"certShirtLS","label":"Green CERT Shirt (LS)"},{"value":"certShirtSS","label":"Green CERT Shirt (SS)"},{"value":"servShirt","label":"Tan SERV Shirt"}]`)
	}
}

func validClearanceRestriction(caller *model.Person, s string) bool {
	switch s {
	case "bgCheck":
		return !caller.IsAdminLeader()
	case "bgDOJ-assumed", "bgDOJ", "bgFBI-assumed", "bgFBI", "bgPHS-assumed", "bgPHS":
		return caller.IsAdminLeader()
	case "cardKey", "certShirtLS", "certShirtSS", "dswCERT", "dswComm", "idPhoto", "servShirt", "volgistics":
		return true
	default:
		return false
	}
}

func matchClearanceRestriction(p *model.Person, s string) bool {
	switch s {
	case "bgCheck":
		if p.Identification&model.IDCardKey != 0 {
			return hasBGCheck(p, model.BGCheckPHS, true)
		}
		return hasBGCheck(p, model.BGCheckFBI, true)
	case "bgDOJ-assumed":
		return hasBGCheck(p, model.BGCheckDOJ, true)
	case "bgDOJ":
		return hasBGCheck(p, model.BGCheckDOJ, false)
	case "bgFBI-assumed":
		return hasBGCheck(p, model.BGCheckFBI, true)
	case "bgFBI":
		return hasBGCheck(p, model.BGCheckFBI, false)
	case "bgPHS-assumed":
		return hasBGCheck(p, model.BGCheckPHS, true)
	case "bgPHS":
		return hasBGCheck(p, model.BGCheckPHS, false)
	case "cardKey":
		return p.Identification&model.IDCardKey != 0
	case "certShirtLS":
		return p.Identification&model.IDCERTShirtLS != 0
	case "certShirtSS":
		return p.Identification&model.IDCERTShirtSS != 0
	case "dswCERT":
		return p.DSWUntil != nil && p.DSWUntil[model.DSWCERT].After(time.Now())
	case "dswComm":
		return p.DSWUntil != nil && p.DSWUntil[model.DSWComm].After(time.Now())
	case "idPhoto":
		return p.Identification&model.IDPhoto != 0
	case "servShirt":
		return p.Identification&model.IDSERVShirt != 0
	case "volgistics":
		return p.VolgisticsID != 0
	}
	return false
}

func hasBGCheck(p *model.Person, check model.BGCheckType, assumedOK bool) bool {
	for _, bc := range p.BGChecks {
		if bc.Type&check != 0 && (assumedOK || !bc.Assumed) {
			return true
		}
	}
	return false
}
