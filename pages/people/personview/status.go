package personview

import (
	"time"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const statusPersonFields = person.FPrivLevels | person.FVolgisticsID | person.FDSWRegistrations | person.FBGChecks | person.FIdentification | person.FFlags

func showStatus(r *request.Request, main *htmlb.Element, user, p *person.Person) {
	if p.ID() != user.ID() && !user.HasPrivLevel(0, enum.PrivLeader) {
		return
	}
	section := main.E("div class=personviewSection")
	sheader := section.E("div class=personviewSectionHeader")
	sheader.E("div class=personviewSectionHeaderText").R(r.LangString("Volunteer Status", "Estado de voluntario"))
	if user.IsAdminLeader() {
		sheader.E("div class=personviewSectionHeaderEdit").
			E("a href=/people/%d/edstatus up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Edit", p.ID())
	}
	section = section.E("div class=personviewStatus")
	showVolgistics(r, section, user, p)
	showDSWCERT(r, section, p)
	showDSWCommunications(r, section, p)
	if user.IsAdminLeader() {
		showBGChecksAL(section, p)
	} else {
		showBGChecksNotAL(r, section, user, p)
	}
	if user.HasPrivLevel(0, enum.PrivLeader) {
		showIdentifications(section, p)
	}
}

func showVolgistics(r *request.Request, section *htmlb.Element, user, p *person.Person) {
	if user.IsAdminLeader() {
		section.E("div>Volgistics")
		if p.VolgisticsID() != 0 {
			section.E("div>#%d", p.VolgisticsID())
		} else if p.Flags()&person.VolgisticsPending != 0 {
			section.E("div", !p.HasPrivLevel(0, enum.PrivMember), "class=personviewStatus-needed").R("Registration pending")
		} else {
			section.E("div", !p.HasPrivLevel(0, enum.PrivMember), "class=personviewStatus-needed").R("Not registered")
		}
	} else if p.VolgisticsID() == 0 {
		section.E("div").R(r.LangString("City volunteer", "Voluntario de la ciudad"))
		if p.Flags()&person.VolgisticsPending != 0 {
			section.E("div").R(r.LangString("Registration pending", "Inscripción pendiente"))
		} else if user.ID() == p.ID() {
			section.E("div").E("a href=/people/%d/vregister up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Inscribirse", p.ID())
		} else {
			section.E("div", !p.HasPrivLevel(0, enum.PrivMember), "class=personviewStatus-needed").R("No se inscribido")
		}
	}
}

func showDSWCERT(r *request.Request, section *htmlb.Element, p *person.Person) {
	needed := p.HasPrivLevel(enum.OrgCERTD, enum.PrivMember) || p.HasPrivLevel(enum.OrgCERTT, enum.PrivMember)
	if cert := p.DSWRegistrations().CERT; cert != nil {
		if cert.Expiration.IsZero() {
			section.E("div>DSW CERT")
			section.E("div>%s %s", r.LangString("Registered", "Registrado"), formatDate(cert.Registered))
		} else if cert.Expiration.After(time.Now()) {
			section.E("div>DSW CERT")
			section.E("div>%s %s, %s %s", r.LangString("Registered", "Registrado"), formatDate(cert.Registered), r.LangString("expires", "expirá"), formatDate(cert.Expiration))
		} else {
			section.E("div>DSW CERT")
			section.E("div", needed, "class=personviewStatus-needed").TF("Expired on %s", formatDate(cert.Expiration))
		}
	} else if needed {
		section.E("div>DSW CERT")
		section.E("div class=personviewStatus-needed").R(r.LangString("Not registered", "No registrado"))
	}
}

func showDSWCommunications(r *request.Request, section *htmlb.Element, p *person.Person) {
	needed := p.HasPrivLevel(enum.OrgSARES, enum.PrivMember)
	if comm := p.DSWRegistrations().Communications; comm != nil {
		if comm.Expiration.IsZero() {
			section.E("div>DSW SARES")
			section.E("div>%s %s", r.LangString("Registered", "Registrado"), formatDate(comm.Registered))
		} else if comm.Expiration.After(time.Now()) {
			section.E("div>DSW SARES")
			section.E("div>%s %s, %s %s", r.LangString("Registered", "Registrado"), formatDate(comm.Registered), r.LangString("expires", "expirá"), formatDate(comm.Expiration))
		} else {
			section.E("div>DSW SARES")
			section.E("div", needed, "class=personviewStatus-needed").R(r.LangString("Expired on ", "Expiró ")).T(formatDate(comm.Expiration))
		}
	} else if needed {
		section.E("div>DSW SARES")
		section.E("div class=personviewStatus-needed").R(r.LangString("Not registered", "No registrado"))
	}
}

func showBGChecksAL(section *htmlb.Element, p *person.Person) {
	bg := p.BGChecks()
	if bg.DOJ == nil && bg.FBI == nil && bg.PHS == nil && !p.HasPrivLevel(0, enum.PrivMember) {
		return
	}
	section.E("div>Background checks")
	div := section.E("div class=personviewStatusBGChecks")
	showBGCheck(div, bg.DOJ, p.HasPrivLevel(0, enum.PrivMember), "DOJ", "NLI")
	showBGCheck(div, bg.FBI, p.HasPrivLevel(0, enum.PrivMember), "FBI", "NLI")
	showBGCheck(div, bg.PHS, p.Identification()&person.IDCardKey != 0, "PHS", "rescinded")
}
func showBGCheck(div *htmlb.Element, check *person.BGCheck, needed bool, label, nliLabel string) {
	if !needed && check == nil {
		return
	}
	div.E("div>%s:", label)
	div = div.E("div", needed && (check == nil || !check.NLI.IsZero()), "class=personviewStatus-needed")
	if check != nil {
		if check.Assumed {
			div.R("assumed ")
		}
		div.R("cleared")
		if !check.Cleared.IsZero() {
			div.R(" ")
			div.R(formatDate(check.Cleared))
		}
		if !check.NLI.IsZero() {
			div.R(", ")
			div.R(nliLabel)
			div.R(" ")
			div.R(formatDate(check.NLI))
		}
	} else if needed {
		div.R("needed")
	}
}

func showBGChecksNotAL(r *request.Request, section *htmlb.Element, user, p *person.Person) {
	bg := p.BGChecks()
	cleared := bg.DOJ != nil && bg.DOJ.NLI.IsZero() && bg.FBI != nil && bg.FBI.NLI.IsZero()
	if p.Identification()&person.IDCardKey != 0 && (bg.PHS == nil || !bg.PHS.NLI.IsZero()) {
		cleared = false
	}
	if cleared {
		section.E("div").R(r.LangString("Background check", "Verificación de antecedentes"))
		section.E("div").R(r.LangString("Cleared", "Aprobada"))
		return
	}
	needed := p.HasPrivLevel(0, enum.PrivMember) || p.Identification()&person.IDCardKey != 0
	if needed {
		section.E("div").R(r.LangString("Background check", "Verificación de antecedentes"))
		section.E("div class=personviewStatus-needed").R(r.LangString("Needed", "Necesitada"))
	}
}

func showIdentifications(section *htmlb.Element, p *person.Person) {
	ids := p.Identification()
	if ids == 0 {
		return
	}
	section.E("div>IDs issued")
	div := section.E("div class=personviewStatusIdents")
	for _, id := range person.AllIdentTypes {
		if ids&id != 0 {
			div.E("div>%s", id.String())
		}
	}
}

func formatDate(t time.Time) string {
	return t.Format("2006\u201101\u201102") // non-breaking hyphens
}
