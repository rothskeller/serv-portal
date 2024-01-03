package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SARESPage handles GET /sares requests.
func SARESPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: "Sunnyvale ARES"}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").E("b").R(r.Loc("Sunnyvale Amateur Radio Emergency Service"))
		main.E("p").R(r.Loc("The Sunnyvale Amateur Radio Emergency Service (SARES) is the local chapter of the nationwide Amateur Radio Emergency Service operated by the Amateur Radio Relay League (ARRL).  During times of emergency, it also operates as a local branch of the federal Radio Amateur Civil Emergency Service (RACES).  SARES provides emergency communications services, usually but not always using amateur radio, when regular communications methods are unavailable or saturated."))
		main.E("p").R(r.Loc("In a disaster, telephones and the Internet will likely be down.  Or, if they are working, they will be unable to keep up with demand.  Radio communications serve as an effective backup because they do not rely on massive, fragile infrastructure.  SARES operators can provide essential emergency communications when no other methods are working.  Outside of emergencies, SARES operators provide ongoing community service by supplying communications assistance at public events on request."))
		main.E("p").R(r.Loc("Membership in SARES requires a current FCC amateur radio license.  If you are interested in emergency communications but do not have a license, SARES members will connect you with resources to help you get one."))
		main.E("p").R(r.Loc("For more information about SARES or amateur radio, write to <a href=mailto:sares@sunnyvale.ca.gov target=_blank>sares@sunnyvale.ca.gov</a>."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
