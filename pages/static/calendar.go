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
	ui.Page(r, user, ui.PageOpts{Title: "SERV Calendar Subscription"}, func(main *htmlb.Element) {
		main.A("class=static").R(`<p>You can subscribe to the SERV calendar so that SERV events will automatically appear in the calendar app on your
  phone, or in your desktop calendar software. Please see the instructions for your phone or software below.
<h1>iPhone or iPad Calendar App</h1>
<ol>
  <li>Open the Settings app.
  <li>Go to “Passwords &amp; Accounts”.
  <li>Tap on “Other”.
  <li>Tap on “Add Subscribed Calendar”.
  <li>In the “Server” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Tap “Next”.
  <li>Optional: change the “Description” field to a name that’s meaningful to you, such as “SERV Calendar”.
  <li>Tap “Save”.
</ol>
<h1>Google Calendar (including Android Phones)</h1>
<ol>
  <li>In a web browser, go to Google Calendar (<code>https://calendar.google.com</code>). Log in if necessary.
  <li>In the left sidebar, click the large “+” sign next to “Other Calendars”.
  <li>Click “From URL”.
  <li>In the “URL of calendar” field, enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Add calendar”.
</ol>
<h1>Microsoft Outlook</h1>
<ol>
  <li>Open Microsoft Outlook.
  <li>Go to the calendar page.
  <li>In the Home ribbon, click on “Open Calendar”, then “From Internet”.
  <li>Enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Yes”.
  <li>In the left sidebar, under “Other Calendars”, right-click on “Untitled” and choose “Rename Calendar”.
  <li>Give the calendar a name meaningful to you, such as “SERV Calendar”.
</ol>
<h1>Mac Calendar App</h1>
<ol>
  <li>Open the Calendar app.
  <li>From the menu, choose File → New Calendar Subscription.
  <li>Enter <code>https://sunnyvaleserv.org/calendar.ics</code>.
  <li>Click “Subscribe”.
  <li>Set the options to suit your preferences and click “OK”.
</ol>
<h1>Other Software</h1>
<p>Most calendar software has the ability to subscribe to Internet calendars. Consult the documentation for your
  software to find out how. The address of the SERV calendar is <code>https://sunnyvaleserv.org/calendar.ics</code>.
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>`)
	})
}
