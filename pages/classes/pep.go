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
	user := auth.SessionUser(r, 0, false)
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
		classes := main.E("div class=classesRegisterGrid")
		if r.Language == "es" {
			classes.E("div").R(`Lunes, el 6 de enero, 6:00–8:00pm
Impartido en inglés
Biblioteca Pública de Sunnyvale
665 W. Olive Avenue, Sunnyvale`)
		} else {
			classes.E("div").R(`Monday, January 5, 6:00–8:00pm
Taught in English
Sunnyvale Public Library
665 W. Olive Avenue, Sunnyvale`)
		}
		classes.E("div").E("a href=https://sunnyvale.libcal.com/event/15686104 target=_blank class='sbtn sbtn-primary sbtn-small'").R(r.Loc("Sign Up"))
		getClassesCommon(r, user, main, class.PEP)
		classes = main.E("div class=classesRegisterGrid")
		classes.E("div").R(r.Loc("Subscribe to our email list to be notified when additional classes are scheduled (English or Spanish)."))
		classes.E("div").E("a href=/pep/notify up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary sbtn-small'").R(r.Loc("Subscribe"))
		text = main.E("div class=pepIntro")
		text.E("p").R(r.Loc("We also teach tailored versions of the class for private groups such as apartment complexes, churches, and businesses.  To arrange a class for your group, please contact us at pep@sunnyvaleserv.org."))
		main.E("div class=classesSERV").R(r.Loc("This class is presented by Sunnyvale Emergency Response Volunteers (SERV), the volunteer arm of the Sunnyvale Office of Emergency Services."))
	})
}
