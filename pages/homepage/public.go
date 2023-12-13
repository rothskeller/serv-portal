package homepage

import (
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func servePublic(r *request.Request) {
	ui.Page(r, nil, ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=pubhome")
		pubhomeHead(main)
		orgs := main.E("div class=pubhomeOrgs")
		pubhomeListos(orgs)
		pubhomeSpanish(orgs)
		pubhomeSNAP(orgs)
		pubhomeCERT(orgs)
		pubhomeSARES(orgs)
		pubhomeFooter(main)
	})
}

func pubhomeHead(main *htmlb.Element) {
	head := main.E("div class=pubhomeHead")
	head.E("div class=pubhomeHeadName>Sunnyvale Emergency Response&nbsp;Volunteers")
	head.E("img class=pubhomeHeadLogo src=%s", ui.AssetURL("serv-logo.png"))
	head.E(`div class=pubhomeHeadText>SERV is the volunteer arm of the
Sunnyvale Office of Emergency Services. SERV volunteers teach disaster
preparedness classes, assist uniformed Public Safety officers in emergencies,
and respond in disasters when professional responders are&nbsp;overloaded.`)
	head.E("div class=pubhomeHeadLogin").
		E("a href=/login class='sbtn sbtn-primary' up-target=main>Volunteer&nbsp;Login")
}

func pubhomeListos(orgs *htmlb.Element) {
	org := orgs.E("div class=pubhomeOrg id=folder-listos")
	org.E("div class=pubhomeOrgName>Disaster Preparedness")
	org.E("img class=pubhomeOrgLogo src=%s", ui.AssetURL("listos-logo.png"))
	org.E(`div class=pubhomeOrgText><b>Listos California</b> disaster
preparedness classes, taught by SERV volunteers, teach how to prepare your home
and family for disasters and&nbsp;emergencies.`)
	links := org.E("div class=pubhomeOrgLinks")
	links.E("div").E("a href=/classes target=_blank>Class schedules and registration")
	links.E("div").E("a href=/files/disaster-preparedness up-target=main>Disaster preparedness materials")
	links.E("div>Email ").E("a href=mailto:listos@sunnyvaleserv.org>Listos@SunnyvaleSERV.org")
}

func pubhomeSpanish(orgs *htmlb.Element) {
	org := orgs.E("div class=pubhomeOrg id=folder-spanish")
	org.E("div class=pubhomeOrgName>Preparación para desastres")
	org.E("img class=pubhomeOrgLogo src=%s", ui.AssetURL("listos-logo.png"))
	org.E(`div class=pubhomeOrgText>Las clases de preparación para desastres
de <b>Listos California,</b> impartidas por voluntarios de SERV, enseñan cómo
preparar su hogar y su familia para desastres y emergencias. Las clases están
disponibles en&nbsp;español.`)
	links := org.E("div class=pubhomeOrgLinks")
	links.E("div").E("a href=/classes target=_blank>Horarios de clases y registro")
	links.E("div").E("a href=/files/preparacion-para-desastres up-target=main>Materiales de preparación para desastres")
	links.E("div>Email ").E("a href=mailto:listos@sunnyvaleserv.org>Listos@SunnyvaleSERV.org")
}

func pubhomeSNAP(orgs *htmlb.Element) {
	org := orgs.E("div class=pubhomeOrg id=folder-snap")
	org.E("div class=pubhomeOrgName>Neighborhood Organization")
	org.E("img class=pubhomeOrgLogo src=%s", ui.AssetURL("snap-logo.png"))
	org.E(`div class=pubhomeOrgText><b>Sunnyvale Neighborhoods Actively
Prepare (SNAP)</b> provides support for neighborhood organization of disaster
preparedness and response&nbsp;volunteers.`)
	links := org.E("div class=pubhomeOrgLinks")
	links.E("div").E("a href=/files/neighborhood-organization up-target=main>Neighborhood organization materials")
	links.E("div>Email ").E("a href=mailto:snap@sunnyvaleserv.org>SNAP@SunnyvaleSERV.org")
}

func pubhomeCERT(orgs *htmlb.Element) {
	org := orgs.E("div class=pubhomeOrg id=folder-cert")
	org.E("div class=pubhomeOrgName>Disaster Response")
	org.E("img class=pubhomeOrgLogo src=%s", ui.AssetURL("cert-logo.png"))
	org.E(`div class=pubhomeOrgText>Sunnyvale’s <b>Community Emergency
Response Team (CERT)</b> supplements uniformed officers in responding to
emergencies, disasters, and public service&nbsp;events.`)
	links := org.E("div class=pubhomeOrgLinks")
	links.E("div").E("a href=/classes target=_blank>Class schedules and registration")
	links.E("div").E("a href=/files/disaster-response up-target=main>Disaster response and CERT materials")
	links.E("div>Email ").E("a href=mailto:cert@sunnyvaleserv.org>CERT@SunnyvaleSERV.org")
}

func pubhomeSARES(orgs *htmlb.Element) {
	org := orgs.E("div class=pubhomeOrg id=folder-sares")
	org.E("div class=pubhomeOrgName>Emergency Communications")
	org.E("img class=pubhomeOrgLogo src=%s", ui.AssetURL("sares-logo.png"))
	org.E(`div class=pubhomeOrgText>The <b>Sunnyvale Amateur Radio Service
(SARES)</b> provides radio communication services during a disaster when regular
methods are overloaded or unavailable. SARES provides a wealth of training for
radio&nbsp;operators.`)
	links := org.E("div class=pubhomeOrgLinks")
	links.E("div").E("a href=https://saresrg.org/ target=_blank>SARES Repeater Group website")
	links.E("div").E("a href=https://www.scc-ares-races.org/activities/events.php target=_blank>Class schedules and registration")
	links.E("div>Email ").E("a href=mailto:sares@sunnyvaleserv.org>SARES@SunnyvaleSERV.org")
}

func pubhomeFooter(main *htmlb.Element) {
	main.E("div class=pubhomeFooter").
		E("a href=/about up-target=main>Web Site Information")
}
