package classes

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// GetMYN handles GET /myn requests.
func GetMYN(r *request.Request) {
	user := auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{
		Title:    r.Loc("Map Your Neighborhood"),
		MenuItem: "classes",
	}, func(main *htmlb.Element) {
		main.E("div class=mynHeading").
			E("img class=mynLogo src=%s", ui.AssetURL(r.Loc("myn-logo.png"))).P().
			E("div class=mynSlogan").R(r.Loc("Planning for disasters\nwith your neighbors"))
		text := main.E("div class=mynIntro")
		text.E("p").R(r.Loc("Following a disaster, Sunnyvale residents will need to rely on each other for several days if city and county services are overwhelmed.  The “Map Your Neighborhood” (MYN) program prepares neighbors to organize a timely response and to support each other in a disaster."))
		text.E("p").R(r.Loc("In this program, we lead a two-hour meeting of around 15–25 households.  Neighbors learn the 9 Steps to take following a disaster, identify resources and skills available in their neighborhood that will be useful in a disaster response, and “map” any special challenges or people with particular needs.  As part of this model, neighbors get to know each other and are better prepared to work together responding to a disaster."))
		text.E("p").R(r.Loc("For more information about this program, or to arrange a MYN meeting for your neighborhood, click the button below and fill out the contact form.  Alternatively, you can write to <a href=mailto:myn@sunnyvale.ca.gov target=_blank>myn@sunnyvale.ca.gov</a>."))
		text.E("div class=staticBack").E("a href=https://forms.gle/bXfpRsGohY9biBi87 target=_blank class='sbtn sbtn-primary'").R(r.Loc("Request Information"))
	})
}
