package static

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// AboutPage handles GET /about requests.
func AboutPage(r *request.Request) {
	ui.Page(r, auth.SessionUser(r, 0, false), ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=static").R(`<h1>Privacy Policy</h1>
<p>
  This web site collects information about people who work for, volunteer for,
  or take classes organized through the Office of Emergency Services (OES) in
  the Sunnyvale Department of Public Safety (DPS).  The information we collect
  includes:
</p>
<ul>
  <li>Basic Information
  <ul>
    <li>name
    <li>amateur radio call sign
    <li>contact information (email addresses, phone numbers, and physical and postal addresses)
    <li>memberships in, and roles held in, SERV volunteer groups
    <li>emergency response classes taken and certificates issued
    <li>credentials that are relevant to SERV operations
    <li>other information voluntarily provided such as skills, languages spoken, available equipment, etc.
  </ul>
  <li>Restricted Information
  <ul>
    <li>attendance at SERV events, and hours spent at them
    <li>Disaster Service Worker registration status
    <li>photo IDs and card access keys issued
    <li>Live Scan fingerprinting success, with date (see note below)
    <li>background check success, with date (see note below)
  </ul>
  <li>Targeted Information
  <ul>
    <li>email messages sent to any SunnyvaleSERV.org address
    <li>text messages sent through this web site
  </ul>
  <li>Private Information
  <ul>
    <li>logs of web site visits and actions taken
  </ul>
</ul>
<p>
  All of the above information is available to the paid and volunteer staff of
  OES and their delegates, including the web site maintainers.  Private
  information is not available to anyone else.
<p>
  If you are a student in an OES-organized class, such as CERT, LISTOS, or
  PEP, your basic and restricted information may be shared with the class
  instructors as long as the class is in progress.
<p>
  If you are a volunteer in a SERV volunteer group, your basic information may
  be shared with other volunteers in that group, and your restricted
  information may be shared with the leaders of that group.
<p>
  If you are a volunteer in a SERV volunteer group, and you have successfully
  completed Live Scan fingerprinting and/or background checks, that fact (with
  no detail other than the date) may be shared with the leaders of your
  volunteer group.  A negative result will not be shared with them.
<p>
  If you have sent any email or text messages (targeted information) through
  the site, they may be shared with any member of the group(s) to which you
  sent them, including members who join those groups after you send the
  messages.
<p>
  If you volunteer for mutual aid or training with another emergency response
  organization or jurisdiction, we may share your basic and/or restricted
  information with them.
<p>
  The OES staff may share anonymized, aggregate data derived from the above
  information with anyone at their discretion.
</p>
<h1>Cookies</h1>
<p>
  This site uses browser cookies.  While you are logged in, a browser cookie
  contains your session identification; this cookie goes away when you log out
  or your login session expires.  More permanent cookies are used to store
  some of your user interface preferences, such as whether you prefer to see
  the events page in calendar or list form.  No personally identifiable
  information is ever stored in browser cookies.
</p>
<h1>Credits and Copyrights</h1>
<p>
  This site was developed by Steven Roth, as a volunteer for the Sunnyvale
  Department of Public Safety.  The site software is copyrighted © 2020–2021
  by Steven Roth.  Steven Roth has granted the Sunnyvale Department of Public
  Safety a non-exclusive, perpetual, royalty-free, worldwide license to use
  this software.  The Sunnyvale Department of Public Safety owns the
  SunnyvaleSERV.org domain and funds the ongoing usage and maintenance of the
  site.
</p>
<h1>Technologies and Services</h1>
<p>
  The software for this web site is written in
  <a href="https://golang.org" target="_blank">Go</a>, with data storage in a
  <a href="https://sqlite.org" target="_blank">SQLite</a> database.  This web
  site is hosted by
  <a href="https://www.dreamhost.com/" target="_blank">Dreamhost</a>.  It uses
  <a href="https://smartystreets.com/" target="_blank">SmartyStreets</a> for
  geolocation,
  <a href="https://www.google.com/maps" target="_blank">Google Maps</a> for
  mapping, and <a href="https://www.twilio.com/" target="_blank">Twilio</a>
  for text messaging.
</p>
<div style="margin:1.5rem 0"><button class="sbtn sbtn-primary" onclick="history.back()">Back</button></div>
`)
	})
}
