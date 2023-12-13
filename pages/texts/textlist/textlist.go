package textlist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /texts requests.
func Get(r *request.Request) {
	const textmsgFields = textmsg.FID | textmsg.FTimestamp | textmsg.FLists | textmsg.FMessage
	var (
		user *person.Person
		opts ui.PageOpts
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !listperson.CanSendText(r, user.ID()) {
		errpage.Forbidden(r, user)
		return
	}
	opts = ui.PageOpts{
		Title:    "Text Messages",
		MenuItem: "texts",
		Tabs: []ui.PageTab{
			{Name: "Messages", URL: "/texts", Target: "main", Active: true},
		},
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		main.E("div class=textlistButtons").
			E("a href=/texts/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>New Message")
		var table *htmlb.Element
		textmsg.All(r, textmsgFields, func(t *textmsg.TextMessage) {
			var visible = user.IsAdminLeader()
			for _, tl := range t.Lists() {
				if tl.ID != 0 && listperson.CanSend(r, user.ID(), tl.ID) {
					visible = true
				}
			}
			if !visible {
				return
			}
			if table == nil {
				table = main.E("div class=textlistTable")
				row := table.E("div class=textlistHeading")
				row.E("div class=textlistTime>Time Sent")
				row.E("div class=textlistLists>Recipients")
				row.E("div class=textlistText>Message")
			}
			row := table.E("div class=textlistRow")
			row.E("a href=/texts/%d class=textlistTime up-target=.pageCanvas>%s", t.ID(), t.Timestamp().Format("2006-01-02 15:04"))
			lists := row.E("div class=textlistLists")
			for _, tl := range t.Lists() {
				lists.E("div>%s", tl.Name)
			}
			row.E("div class=textlistText>%s", t.Message())
		})
		if table == nil {
			main.E("div>No messages have been sent which are visible to you.")
		}
	})
}
