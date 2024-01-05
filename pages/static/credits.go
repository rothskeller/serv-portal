package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// CreditsPage handles GET /site-credits requests.
func CreditsPage(r *request.Request) {
	ui.Page(r, auth.SessionUser(r, 0, false), ui.PageOpts{}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("h1").R(r.Loc("Credits and Copyrights"))
		main.E("p").R(r.Loc("This site was developed by Steven Roth, as a volunteer for the Sunnyvale Department of Public Safety.  The site software is copyrighted © 2020–2021 by Steven Roth.  Steven Roth has granted the Sunnyvale Department of Public Safety a non-exclusive, perpetual, royalty-free, worldwide license to use this software.  The Sunnyvale Department of Public Safety owns the SunnyvaleSERV.org domain and funds the ongoing usage and maintenance of the site."))
		main.E("h1").R(r.Loc("Technologies and Services"))
		main.E("p").R(r.Loc("The software for this web site is written in <a href=https://golang.org target=_blank>Go</a>, with data storage in a <a href=https://sqlite.org target=_blank>SQLite</a> database.  This web site is hosted by <a href=https://www.dreamhost.com/ target=_blank>Dreamhost</a>.  It uses <a href=https://www.google.com/maps target=_blank>Google Maps</a> for geolocation and mapping, <a href=https://www.twilio.com/ target=_blank>Twilio</a> for text messaging, and <a href=https://www.algolia.com/ target=_blank>Algolia</a> for searching."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
