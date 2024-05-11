package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SNAPPage handles GET /sares requests.
func SNAPPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: "SNAP"}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").E("b").R(r.Loc("Sunnyvale Neighborhoods Actively Prepare (SNAP)"))
		main.E("p").R(r.Loc("SNAP is our neighborhood disaster preparedness program.  Following a disaster, Sunnyvale residents will need to rely on each other for several days if city and county services are overwhelmed.  While our <a href=/listos up-target=main>Listos</a> program teaches preparedness for individuals and families, SNAP prepares neighbors to organize a timely response and to support each other in a disaster."))
		main.E("p").R(r.Loc("Using the “Map Your Neighborhood” (MYN) program provided by the Washington State Emergency Management Division, we lead a two-hour meeting of around 15–25 households.  Neighbors learn the 9 Steps to take following a disaster, identify resources and skills available in their neighborhood that will be useful in a disaster response, and “map” any special challenges or people with particular needs.  As part of this model, neighbors get to know each other and are better prepared to work together responding to a disaster."))
		main.E("p").R(r.Loc("For more information about SNAP, or to arrange a MYN meeting for your neighborhood, write to <a href=mailto:snap@sunnyvale.ca.gov target=_blank>snap@sunnyvale.ca.gov</a>."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
