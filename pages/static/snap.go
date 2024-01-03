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
		main.E("p").R(r.Loc("The SNAP program (“Sunnyvale Neighborhoods Actively Prepare”) is our neighborhood disaster preparedness program.  While our <a href=/listos up-target=main>Listos</a> program teaches preparedness for individuals and families, SNAP prepares neighbors to support each other in a disaster.  When someone is willing to host a preparedness event for their neighborhood (ideally 15–20 homes), we can help facilitate the event and raise the preparedness level of the whole group."))
		main.E("p").R(r.Loc("To do this, we make use of the “Map Your Neighborhood” program provided by the Washington State Emergency Management Division.  We help the neighbors build a map and a common understanding of the resources and skills available in their neighborhood in a disaster, and any special challenges or people with particular needs.  In the process, the neighbors get to know each other better and are more prepared to face a disaster together."))
		main.E("p").R(r.Loc("For more information about SNAP, or to arrange an event for your neighborhood, write to <a href=mailto:snap@sunnyvale.ca.gov target=_blank>snap@sunnyvale.ca.gov</a>."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
