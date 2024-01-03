package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// CERTPage handles GET /cert requests.
func CERTPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Sunnyvale CERT")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").E("b").R(r.Loc("Community Emergency Response Team (CERT)"))
		main.E("p").R(r.Loc("CERT is a nationwide program, managed by the Federal Emergency Management Agency (FEMA), that prepares residents to care for themselves and their communities during and after major disasters.  Its emphasis is on training residents to be able to respond safely and effectively during an emergency."))
		main.E("p").R(r.Loc("The CERT program was created by the Los Angeles Fire Department after seeing the significant loss of life of volunteer rescuers in the 1985 Mexico City earthquake.  Volunteers are credited with having saved many lives in the aftermath of that earthquake, but many of the volunteers were killed because they did not know how to keep themselves safe while doing such work.  LAFD created the CERT program to ensure that the same thing didn't happen on their watch.  The 1987 Whittier earthquake near Los Angeles underscored the value of this program.  It the early 1990s, FEMA expanded the program to cover other disasters besides earthquakes, and spread it nationwide."))
		main.E("p").R(r.Loc("In Sunnyvale, we teach the FEMA-standard <a href=/cert-basic up-target=main>CERT Basic Training<a> class, with some local enhancements, to anyone who wants it. This is a 30-hour class, taught over seven weeks, covering all aspects of volunteer disaster response.  For the graduates of that class, we also teach occasional refresher classes on specific CERT topics to help our volunteers keep their skills and knowledge fresh."))
		main.E("p").R(r.Loc("Sunnyvale also has a “CERT Deployment Team.”  This is a group of CERT-trained volunteers who have agreed to be on call to assist the professional responders in the Department of Public Safety when needed.  Our CERT Deployment Team receives additional, monthly training covering both the CERT topics and more advanced public safety skills."))
		main.E("p").R(r.Loc("For more information about our CERT program, write to <a href=mailto:cert@sunnyvale.ca.gov target=_blank>cert@sunnyvale.ca.gov</a>."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
