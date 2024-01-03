package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// SubscribeCalendarPage handles GET /subscribe-calendar requests.
func SubscribeCalendarPage(r *request.Request) {
	var user *person.Person

	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("SERV Calendar Subscription")}, func(main *htmlb.Element) {
		main = main.A("class=static")
		main.E("p").R(r.Loc("You can subscribe to the SERV calendar so that SERV events will automatically appear in the calendar app on your phone, or in your desktop calendar software. Please see the instructions for your phone or software below."))
		main.E("h1").R(r.Loc("iPhone or iPad Calendar App"))
		ol := main.E("ol")
		ol.E("li").R(r.Loc("Open the Settings app."))
		ol.E("li").R(r.Loc("Go to “Calendar”."))
		ol.E("li").R(r.Loc("Go to “Accounts”."))
		ol.E("li").R(r.Loc("Go to “Add Account”."))
		ol.E("li").R(r.Loc("Tap on “Other”."))
		ol.E("li").R(r.Loc("Tap on “Add Subscribed Calendar”."))
		ol.E("li").R(r.Loc("In the “Server” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>."))
		ol.E("li").R(r.Loc("Tap “Next”."))
		ol.E("li").R(r.Loc("Optional: change the “Description” field to a name that’s meaningful to you, such as “SERV Calendar”."))
		ol.E("li").R(r.Loc("Tap “Save”."))
		main.E("h1").R(r.Loc("Google Calendar (including Android Phones)"))
		ol = main.E("ol")
		ol.E("li").R(r.Loc("In a web browser, go to Google Calendar (<code>https://calendar.google.com</code>). Log in if necessary."))
		ol.E("li").R(r.Loc("In the left sidebar, click the large “+” sign next to “Other Calendars”."))
		ol.E("li").R(r.Loc("Click “From URL”."))
		ol.E("li").R(r.Loc("In the “URL of calendar” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>."))
		ol.E("li").R(r.Loc("Click “Add calendar”."))
		main.E("h1").R(r.Loc("Microsoft Outlook"))
		ol = main.E("ol")
		ol.E("li").R(r.Loc("Open Microsoft Outlook."))
		ol.E("li").R(r.Loc("Go to the calendar page."))
		ol.E("li").R(r.Loc("In the Home ribbon, click on “Open Calendar”, then “From Internet”."))
		ol.E("li").R(r.Loc("Enter <code>https://sunnyvaleserv.org/calendar.ics</code>."))
		ol.E("li").R(r.Loc("Click “Yes”."))
		ol.E("li").R(r.Loc("In the left sidebar, under “Other Calendars”, right-click on “Untitled” and choose “Rename Calendar”."))
		ol.E("li").R(r.Loc("Give the calendar a name meaningful to you, such as “SERV Calendar”."))
		main.E("h1").R(r.Loc("Mac Calendar App"))
		ol = main.E("ol")
		ol.E("li").R(r.Loc("Open the Calendar app."))
		ol.E("li").R(r.Loc("From the menu, choose File → New Calendar Subscription."))
		ol.E("li").R(r.Loc("Enter <code>https://sunnyvaleserv.org/calendar.ics</code>."))
		ol.E("li").R(r.Loc("Click “Subscribe”."))
		ol.E("li").R(r.Loc("Set the options to suit your preferences and click “OK”."))
		main.E("h1").R(r.Loc("Other Software"))
		main.E("p").R(r.Loc("Most calendar software has the ability to subscribe to Internet calendars. Consult the documentation for your software to find out how. The address of the SERV calendar is <code>https://sunnyvaleserv.org/calendar.ics</code>."))
		main.E("div class=staticBack").E("button class='sbtn sbtn-primary' onclick='history.back()'").R(r.Loc("Back"))
	})
}
