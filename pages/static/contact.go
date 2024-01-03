package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// ContactUsPage handles GET /contact requests.
func ContactUsPage(r *request.Request) {
	var user = auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Contact Us")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").R(r.Loc("Sunnyvale Emergency Response Volunteers (SERV) is the volunteer arm of the Sunnyvale Office of Emergency Services, which is part of the city’s Department of Public Safety."))
		bq := main.E("blockquote")
		bq.E("a href=mailto:serv@sunnyvale.ca.gov target=_blank>serv@sunnyvale.ca.gov")
		bq.E("br")
		bq.R(r.Loc("(408) 730–7190 English (messages only)"))
		bq.E("br")
		bq.R(r.Loc("(408) 730-7294 Spanish (messages only)"))
		main.E("p").R(r.Loc("Our offices are at"))
		bq = main.E("blockquote")
		bq.R(r.Loc("Sunnyvale Public Safety Headquarters"))
		bq.E("br")
		bq.R("700 All America Way<br>Sunnyvale, CA  94086")
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
