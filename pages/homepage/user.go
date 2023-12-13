package homepage

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func serveUser(r *request.Request, user *person.Person) {
	ui.Page(r, user, ui.PageOpts{}, func(main *htmlb.Element) {
		// Normally, on a narrow display, if the user opens the menu,
		// it overlays the page canvas.  In this one case, we want the
		// page canvas shifted to the right instead.
		main.A("class='userhome pageCanvas-menuShift'")
		main.E("div class=userhomeWelcome>Welcome to SunnyvaleSERV.org!")
		// When the menu is closed, we'll draw an arrow pointing to it
		// and tell people they can open it.
		main.E("svg class=userhomeArrow viewBox=%s", "0 0 448 2048").
			E("path fill=currentColor d=%s", "M34.9 289.5l-22.2-22.2c-9.4-9.4-9.4-24.6 0-33.9L207 39c9.4-9.4 24.6-9.4 33.9 0l194.3 194.3c9.4 9.4 9.4 24.6 0 33.9L413 289.4c-9.5 9.5-25 9.3-34.3-.4L264 168.6V1992c0 13.3-10.7 24-24 24h-32c-13.3 0-24-10.7-24-24V168.6L69.2 289.1c-9.3 9.8-24.8 10-34.3.4z")
		main.E("div class=userhomeClosed>Click here to open the menu.")
		// When the menu is open, we'll display helper text next to each
		// menu item.
		helpers := main.E("div class=userhomeHelpers")
		// First, the Events item.
		helpers.E("div class=userhomeHelper").
			E("s-icon icon=left class=userhomeHelper-left").P().
			E("div>Click here to get details of upcoming classes and events.")
		// Next, the People item.
		if user.HasPrivLevel(0, enum.PrivStudent) {
			helpers.E("div class=userhomeHelper").
				E("s-icon icon=left class=userhomeHelper-left").P().
				E("div>Click here to view team rosters and maps.")
		}
		// Next, the Files item.
		helpers.E("div class=userhomeHelper").
			E("s-icon icon=left class=userhomeHelper-left").P().
			E("div>Click here for class materials and other documents.")
		// Next, the Reports item.
		if user.HasPrivLevel(0, enum.PrivLeader) {
			helpers.E("div class=userhomeHelper").
				E("s-icon icon=left class=userhomeHelper-left").P().
				E("div>Click here to generate reports.")
		}
		// Next, the Texts item.
		if listperson.CanSendText(r, user.ID()) {
			helpers.E("div class=userhomeHelper").
				E("s-icon icon=left class=userhomeHelper-left").P().
				E("div>Click here to send group text messages.")
		}
		// Next, the Admin item.
		if user.IsWebmaster() {
			helpers.E("div class=userhomeHelper").
				E("s-icon icon=left class=userhomeHelper-left").P().
				E("div>If you can see this, you shouldn't need handholding.")
		}
		// Next, the Profile item.
		helpers.E("div class=userhomeHelper").
			E("s-icon icon=left class=userhomeHelper-left").P().
			E("div>Click here to change your password or contact info.")
		// Lastly, the Logout item.
		helpers.E("div class=userhomeHelper").
			E("s-icon icon=left class=userhomeHelper-left").P().
			E("div>Click here to log out of the web site.")
		// Below the menu item helpers, we add some other help
		// information.
		helpers.E("div class=userhomeSeealso").
			E("div>Also see:").P().
			E("div").E("a href=/subscribe-calendar up-follow>Subscribe to the SERV calendar on your phone").P().
			E("div").E("a href=/email-lists up-follow>Information about SERV email lists")
	})
}
