package classes

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// GetCERT handles GET /cert-basic requests.
func GetCERT(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{
		Title: r.Loc("CERT Basic Training"),
	}, func(main *htmlb.Element) {
		main.E("div class=certHeading").
			E("img class=certLogo src=%s", ui.AssetURL(r.Loc("cert-logo.png"))).P().
			E("div class=certSlogan").R(r.Loc("How to help your community after a disaster"))
		main.E("div class=pepIntro").R(r.Loc("<p>In a disaster, professional emergency responders will be overwhelmed, and people will have to rely on their neighbors for help.  If you want to be one of the helpers, this class is for you.  It teaches basic emergency response skills, and how to use them safely.  Topics include:</p><ul><li>Disaster Preparedness<li>The CERT Organization<li>Usage of Personal Protective Equipment (PPE)<li>Disaster Medical Operations<li>Triaging, Assessing, and Treating Patients<li>Disaster Psychology<li>Fire Safety and Utility Control<li>Extinguishing Small Fires<li>Light Search and Rescue<li>Terrorism and CERT<li>Disaster Simulation Exercise</ul><p>This class meets for seven weekday evenings and one full Saturday (see dates below).  On successful completion of the class, you will be invited to join the Sunnyvale CERT Deployment Team, which supports the professional responders in Sunnyvale's Department of Public Safety.</p><p>IMPORTANT:  Space in this class is limited.  Please do not sign up unless you fully expect to attend all of the sessions.  This class is open to anyone aged 18 or over, but preference will be given to Sunnyvale residents.  High school students under age 18 are welcome if their parent or other responsible adult is also in the class.</p>"))
		getClassesCommon(r, main, class.CERTBasic)
	})
}
