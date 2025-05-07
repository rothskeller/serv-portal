package classes

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// GetClasses handles GET /classes or GET /clases requests.
func GetClasses(r *request.Request) {
	user := auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{
		Title:    r.Loc("Classes and Training"),
		MenuItem: "classes",
	}, func(main *htmlb.Element) {
		classes := main.E("div class=classes")
		pep := classes.E("div class=classesBlock")
		pep.E("div class=classesBlockHeading").R(r.Loc("Personal Emergency Preparedness"))
		pep.E("div class=pepHeading").
			E("img class=pepLogo src=%s", ui.AssetURL(r.Loc("pep-logo.png"))).P().
			E("div class=pepSlogan").R(r.Loc("Are you prepared\nfor a disaster?"))
		text := pep.E("div class=pepIntro")
		text.E("p").R(r.Loc("Earthquakes, fires, floods, pandemics, power outages, chemical spills ... these are just some of the disasters than can strike our area without warning.  After a disaster strikes, professional emergency services may not be available to help you for several days.  Are you fully prepared to take care of yourself and your family if the need arises?"))
		pep.E("button type=button class='classesViewMore viewmore sbtn sbtn-small sbtn-primary' data-target=classesPEPMore").R(r.Loc("View More"))
		more := pep.E("div id=classesPEPMore class=classesMore")
		text = more.E("div class=pepIntro")
		text.E("p").R(r.Loc("Our <b>Personal Emergency Preparedness</b> class can help you prepare for disasters.  It will teach you about the various disasters you might face, what preparations you can make for them, and how to prioritize."))
		cgrid := main.E("div class=classesRegisterGrid")
		if r.Language == "es" {
			cgrid.E("div").R(`Miércoles, el 25 de junio, 6:30–8:30pm
Impartido en inglés
Biblioteca Pública de Sunnyvale
665 W. Olive Avenue, Sunnyvale`)
		} else {
			cgrid.E("div").R(`Wednesday, June 25, 6:30–8:30pm
Taught in English
Sunnyvale Public Library
665 W. Olive Avenue, Sunnyvale`)
		}
		cgrid.E("div").E("a href=https://sunnyvale.libcal.com/event/14553820 target=_blank class='sbtn sbtn-primary sbtn-small'").R(r.Loc("Sign Up"))
		getClassesCommon(r, user, more, class.PEP)
		cert := classes.E("div class=classesBlock")
		cert.E("div class=classesBlockHeading").R(r.Loc("CERT Basic Training"))
		cert.E("div class=certHeading").
			E("img class=certLogo src=%s", ui.AssetURL("cert-logo.png")).P().
			E("div class=certSlogan").R(r.Loc("How to help your community after a disaster"))
		text = cert.E("div class=certIntro")
		text.E("p").R(r.Loc("In a disaster, professional emergency responders will be overwhelmed, and people will have to rely on their neighbors for help.  If you want to be one of the helpers, the <b>Community Emergency Response Team (CERT) Basic Training</b> class is for you.  It teaches basic emergency response skills, and how to use them safely."))
		cert.E("button type=button class='classesViewMore viewmore sbtn sbtn-small sbtn-primary' data-target=classesCERTMore").R(r.Loc("View More"))
		more = cert.E("div id=classesCERTMore class=classesMore")
		text = more.E("div class=certIntro")
		text.E("p").R(r.Loc("Topics include:<ul><li>Disaster Preparedness<li>The CERT Organization<li>Usage of Personal Protective Equipment (PPE)<li>Disaster Medical Operations<li>Triaging, Assessing, and Treating Patients<li>Disaster Psychology<li>Fire Safety and Utility Control<li>Extinguishing Small Fires<li>Light Search and Rescue<li>Terrorism and CERT<li>Disaster Simulation Exercise</ul>"))
		text.E("p").R(r.Loc("This class meets for seven weekday evenings and one full Saturday (see dates below).  On successful completion of the class, you will be invited to join the Sunnyvale CERT Deployment Team, which supports the professional responders in Sunnyvale's Department of Public Safety."))
		text.E("p").R(r.Loc("<b>IMPORTANT:</b>  Space in this class is limited.  Please do not sign up unless you fully expect to attend all of the sessions.  This class is open to anyone aged 18 or over, but preference will be given to Sunnyvale residents.  High school students under age 18 are welcome if their parent or other responsible adult is also in the class."))
		if r.Language != "en" {
			text.E("p").R(r.Loc("<b>IMPORTANT:</b>  This class is taught only in English.  However, the printed materials are available in Spanish."))
		}
		getClassesCommon(r, user, more, class.CERTBasic)
	})
}
