package classes

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// GetPEP handles GET /pep and GET /ppde requests.
func GetPEP(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{
		Title:    r.Loc("Personal Emergency Preparedness"),
		MenuItem: "classes",
	}, func(main *htmlb.Element) {
		main.E("div class=pepHeading").
			E("img class=pepLogo src=%s", ui.AssetURL(r.Loc("pep-logo.png"))).P().
			E("div class=pepSlogan").R(r.Loc("Are you prepared\nfor a disaster?"))
		text := main.E("div class=pepIntro")
		text.E("p").R(r.Loc("Earthquakes, fires, floods, pandemics, power outages, chemical spills ... these are just some of the disasters than can strike our area without warning.  After a disaster strikes, professional emergency services may not be available to help you for several days.  Are you fully prepared to take care of yourself and your family if the need arises?"))
		text.E("p").R(r.Loc("Our <b>Personal Emergency Preparedness</b> class can help you prepare for disasters.  It will teach you about the various disasters you might face, what preparations you can make for them, and how to prioritize."))
		getClassesCommon(r, user, main, class.PEP)
	})
}
