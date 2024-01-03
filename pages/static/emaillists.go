package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// EmailListsPage handles GET /email-lists requests.
func EmailListsPage(r *request.Request) {
	var user *person.Person

	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("SERV Email Lists")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("h1").R(r.Loc("SERV Email Lists"))
		main.E("p").R(r.Loc("The SunnyvaleSERV.org site offers a number of email distribution lists. We have one for each volunteer program, that we give out to the general public who might want more information about the program.  Email sent to these lists is delivered to designated public contact people for each program:"))
		ul := main.E("ul class=emaillist")
		ul.E("li>cert@sunnyvaleserv.org")
		ul.E("li>listos@sunnyvaleserv.org")
		ul.E("li>sares@sunnyvaleserv.org")
		ul.E("li>snap@sunnyvaleserv.org")
		main.E("p").R(r.Loc("There are also lists for the volunteers on each of our teams:"))
		ul = main.E("ul class=emaillist")
		ul.E("li>cert-alpha@sunnyvaleserv.org")
		ul.E("li>cert-committee@sunnyvaleserv.org")
		ul.E("li>listos-team@sunnyvaleserv.org")
		ul.E("li>outreach-team@sunnyvaleserv.org")
		ul.E("li>sares-active@sunnyvaleserv.org")
		ul.E("li>sares-leads@sunnyvaleserv.org")
		ul.E("li>snap-team@sunnyvaleserv.org")
		main.E("p").R(r.Loc("and for the students in each CERT class:"))
		ul = main.E("ul class=emaillist")
		ul.E("li>cert-60@sunnyvaleserv.org")
		ul.E("li>cert-61@sunnyvaleserv.org")
		ul.E("li>cert-62@sunnyvaleserv.org")
		ul.E("li style=font:inherit>etc.")
		main.E("p").R(r.Loc("Finally, there are some broader lists for special purposes:"))
		ul = main.E("ul class=emaillist")
		ul.E("li>serv-all@sunnyvaleserv.org")
		ul.E("li>volunteer-hours@sunnyvaleserv.org")
		main.E("p").R(r.Loc("All of these email lists have restricted access.  For the team lists, only members of the team can send mail to them; for the class lists, only the instructors can send mail to them; and for the broader lists, only DPS staff can send mail to them.  Any mail sent to any of our lists from someone else is held for approval before being routed to the list.  Messages on topics unrelated to SERV will generally be rejected."))
		main.E("p").R(r.Loc("If you are receiving email from one of these lists that you do not want, there is an “unsubscribe” link at the bottom of every email.  If you are receiving email at the wrong address, you can change your email address in the “Profile” section of this web site."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
