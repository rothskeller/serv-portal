package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// ListosPage handles GET /listos requests.
func ListosPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: "Listos California"}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").E("b>Listos California")
		main.E("p").R(r.Loc("Listos California is a state program, managed by the California Office of Emergency Services (CalOES), focusing on disaster preparedness education for California residents.  Under their umbrella, the Sunnyvale Listos program provides disaster preparedness education in Sunnyvale."))
		main.E("p").R(r.Loc("Our flagship offering is our <a href=/pep up-target=main>Personal Emergency Preparedness</a> class.  This is a two-hour class that teaches home and family preparedness.  We offer this class to the general public every 2–3 months, in both English and Spanish.  We also offer it to neighborhood associations, businesses, etc. when requested."))
		main.E("p").R(r.Loc("Our disaster preparedness education efforts also include Outreach booths and tables at public events (the Arts and Wine Festival, the Diwali Festival, the Firefighters Pancake Breakfast, neighborhood block parties, etc.).  At these events, we set up tables and distribute disaster preparedness information to participants."))
		main.E("p").R(r.Loc("For more information about Listos California or our disaster preparedness education programs, write us at <a href=mailto:listos@sunnyvale.ca.gov target=_blank>listos@sunnyvale.ca.gov</a>. Also write to us if you want to arrange a private preparedness class for your neighborhood or group, or have a preparedness table at your event."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
