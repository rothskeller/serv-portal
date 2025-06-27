package classes

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// GetMoulage handles GET /moulage requests.
func GetMoulage(r *request.Request) {
	user := auth.SessionUser(r, 0, false)
	ui.Page(r, user, ui.PageOpts{
		Title:    r.Loc("Moulage Training"),
		MenuItem: "classes",
	}, func(main *htmlb.Element) {
		main.E("div class=moulageHeading>Moulage Training")
		text := main.E("div class=moulageIntro")
		text.E("p>Help us take CERT exercises to a higher level by learning how to apply fake wounds to live volunteer “victims.”  Live victims amplify the realism of CERT exercises, such as the disaster scenario at the end of each CERT Basic Training class, or the annual county-wide CERT exercises.  Making those live victims look injured is a valued skill.")
		text.E("p>In this class, you’ll learn how to apply different types of fake wounds, ranging from scratches to amputated hands.  You’ll also learn how to coach volunteer victims on how to act out their injuries for greater realism.")
		text.E("p>This class size is limited.  If it fills, preference will be given to Sunnyvale volunteers and/or past moulage helpers.")
		getClassesCommon(r, user, main, class.Moulage)
		main.E("div class=classesSERV").R(r.Loc("This class is presented by Sunnyvale Emergency Response Volunteers (SERV), the volunteer arm of the Sunnyvale Office of Emergency Services."))
	})
}
