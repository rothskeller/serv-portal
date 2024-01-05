package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// PrivacyPage handles GET /privacy-policy requests.
func PrivacyPage(r *request.Request) {
	ui.Page(r, auth.SessionUser(r, 0, false), ui.PageOpts{}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("h1").R(r.Loc("Privacy Policy"))
		main.E("p").R(r.Loc("This web site collects information about people who work for, volunteer for, or take classes organized through the Office of Emergency Services (OES) in the Sunnyvale Department of Public Safety (DPS).  The information we collect includes:"))
		l1 := main.E("ul")
		l1.E("li").R(r.Loc("Basic Information"))
		l2 := l1.E("ul")
		l2.E("li").R(r.Loc("name"))
		l2.E("li").R(r.Loc("amateur radio call sign"))
		l2.E("li").R(r.Loc("contact information (email addresses, phone numbers, and physical and postal addresses)"))
		l2.E("li").R(r.Loc("memberships in, and roles held in, SERV volunteer groups"))
		l2.E("li").R(r.Loc("emergency response classes taken and certificates issued"))
		l2.E("li").R(r.Loc("credentials that are relevant to SERV operations"))
		l2.E("li").R(r.Loc("other information voluntarily provided such as skills, languages spoken, available equipment, etc."))
		l1.E("li").R(r.Loc("Restricted Information"))
		l2 = l1.E("ul")
		l2.E("li").R(r.Loc("attendance at SERV events, and hours spent at them"))
		l2.E("li").R(r.Loc("Disaster Service Worker registration status"))
		l2.E("li").R(r.Loc("photo IDs and card access keys issued"))
		l2.E("li").R(r.Loc("Live Scan fingerprinting success, with date (see note below)"))
		l2.E("li").R(r.Loc("background check success, with date (see note below)"))
		l1.E("li").R(r.Loc("Targeted Information"))
		l2 = l1.E("ul")
		l2.E("li").R(r.Loc("email messages sent to any SunnyvaleSERV.org address"))
		l2.E("li").R(r.Loc("text messages sent through this web site"))
		l1.E("li").R(r.Loc("Private Information"))
		l2 = l1.E("ul")
		l2.E("li").R(r.Loc("logs of web site visits and actions taken"))
		main.E("p").R(r.Loc("All of the above information is available to the paid and volunteer staff of OES and their delegates, including the web site maintainers.  Private information is not available to anyone else."))
		main.E("p").R(r.Loc("If you are a student in an OES-organized class, such as CERT, Listos, or PEP, your basic and restricted information may be shared with the class instructors as long as the class is in progress."))
		main.E("p").R(r.Loc("If you are a volunteer in a SERV volunteer group, your basic information may be shared with other volunteers in that group, and your restricted information may be shared with the leaders of that group."))
		main.E("p").R(r.Loc("If you are a volunteer in a SERV volunteer group, and you have successfully completed Live Scan fingerprinting and/or background checks, that fact (with no detail other than the date) may be shared with the leaders of your volunteer group.  A negative result will not be shared with them."))
		main.E("p").R(r.Loc("If you have sent any email or text messages (targeted information) through the site, they may be shared with any member of the group(s) to which you sent them, including members who join those groups after you send the messages."))
		main.E("p").R(r.Loc("If you volunteer for mutual aid or training with another emergency response organization or jurisdiction, we may share your basic and/or restricted information with them."))
		main.E("p").R(r.Loc("The OES staff may share anonymized, aggregate data derived from the above information with anyone at their discretion."))
		main.E("h1").R(r.Loc("Cookies"))
		main.E("p").R(r.Loc("This site uses browser cookies.  While you are logged in, a browser cookie contains your session identification; this cookie goes away when you log out or your login session expires.  More permanent cookies are used to store some of your user interface preferences, such as your preferred language and whether you prefer to see the events page in calendar or list form.  No personally identifiable information is ever stored in browser cookies."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
