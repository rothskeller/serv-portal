package clearrep

import (
	"time"

	"sunnyvaleserv.org/portal/store/person"
)

func validRestrictions(user *person.Person) []string {
	if user.IsAdminLeader() {
		return []string{"bgDOJ-assumed", "bgDOJ", "bgFBI-assumed", "bgFBI", "bgPHS-assumed", "bgPHS", "cardKey", "idPhoto", "dswCERT", "dswComm", "certShirtLS", "certShirtSS", "volgistics", "servShirt"}
	}
	return []string{"bgCheck", "volgistics", "cardKey", "idPhoto", "dswCERT", "dswComm", "certShirtLS", "certShirtSS", "servShirt"}
}

var restrictionLabels = map[string]string{
	"bgCheck":       "BG Check",
	"bgDOJ-assumed": "BG Check: DOJ (assumed)",
	"bgDOJ":         "BG Check: DOJ (recorded)",
	"bgFBI-assumed": "BG Check: FBI (assumed)",
	"bgFBI":         "BG Check: FBI (recorded)",
	"bgPHS-assumed": "BG Check: PHS (assumed)",
	"bgPHS":         "BG Check: PHS (recorded)",
	"cardKey":       "DPS Card Key",
	"certShirtLS":   "Green CERT Shirt (LS)",
	"certShirtSS":   "Green CERT Shirt (SS)",
	"dswCERT":       "DSW for CERT",
	"dswComm":       "DSW for Communications",
	"idPhoto":       "DPS Photo ID",
	"servShirt":     "Tan SERV Shirt",
	"volgistics":    "In Volgistics",
}

func validRestriction(user *person.Person, s string) bool {
	switch s {
	case "bgCheck":
		return !user.IsAdminLeader()
	case "bgDOJ-assumed", "bgDOJ", "bgFBI-assumed", "bgFBI", "bgPHS-assumed", "bgPHS":
		return user.IsAdminLeader()
	case "cardKey", "certShirtLS", "certShirtSS", "dswCERT", "dswComm", "idPhoto", "servShirt", "volgistics":
		return true
	default:
		return false
	}
}

func matchRestriction(p *person.Person, s string) bool {
	switch s {
	case "bgCheck":
		if p.Identification()&person.IDCardKey != 0 {
			return hasBGCheck(p.BGChecks().PHS, true)
		}
		return hasBGCheck(p.BGChecks().FBI, true)
	case "bgDOJ-assumed":
		return hasBGCheck(p.BGChecks().DOJ, true)
	case "bgDOJ":
		return hasBGCheck(p.BGChecks().DOJ, false)
	case "bgFBI-assumed":
		return hasBGCheck(p.BGChecks().FBI, true)
	case "bgFBI":
		return hasBGCheck(p.BGChecks().FBI, false)
	case "bgPHS-assumed":
		return hasBGCheck(p.BGChecks().PHS, true)
	case "bgPHS":
		return hasBGCheck(p.BGChecks().PHS, false)
	case "cardKey":
		return p.Identification()&person.IDCardKey != 0
	case "certShirtLS":
		return p.Identification()&person.IDCERTShirtLS != 0
	case "certShirtSS":
		return p.Identification()&person.IDCERTShirtSS != 0
	case "dswCERT":
		return hasDSW(p.DSWRegistrations().CERT)
	case "dswComm":
		return hasDSW(p.DSWRegistrations().Communications)
	case "idPhoto":
		return p.Identification()&person.IDPhoto != 0
	case "servShirt":
		return p.Identification()&person.IDSERVShirt != 0
	case "volgistics":
		return p.VolgisticsID() != 0
	}
	return false
}

func hasDSW(dsw *person.DSWRegistration) bool {
	if dsw == nil {
		return false
	}
	return dsw.Expiration.IsZero() || dsw.Expiration.After(time.Now())
}

func hasBGCheck(bc *person.BGCheck, assumedOK bool) bool {
	if bc == nil || (bc.Assumed && !assumedOK) {
		return false
	}
	return bc.NLI.IsZero()
}
