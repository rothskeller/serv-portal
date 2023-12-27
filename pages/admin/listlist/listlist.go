package listlist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/lists requests.
func Get(r *request.Request) {
	var (
		user *person.Person
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !user.IsWebmaster() {
		errpage.Forbidden(r, user)
		return
	}
	Render(r, user)
}

func Render(r *request.Request, user *person.Person) {
	var opts = ui.PageOpts{
		Title:    "Lists",
		MenuItem: "admin",
		Tabs: []ui.PageTab{
			{Name: "Roles", URL: "/admin/roles", Target: "main"},
			{Name: "Lists", URL: "/admin/lists", Target: "main", Active: true},
			{Name: "Classes", URL: "/admin/classes", Target: "main"},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		grid := main.E("div class=listlistGrid")
		row := grid.E("div class=listlistHeading")
		row.E("div>List")
		row.E("div>Sub")
		row.E("div>Unsub")
		row.E("div>Send")
		list.All(r, func(l *list.List) {
			var name string
			var senderCount, subCount, unsubCount int

			if l.Type == list.SMS {
				name = "SMS: " + l.Name
			} else {
				name = l.Name + "@sunnyvaleserv.org"
			}
			listperson.All(r, l.ID, 0, func(_ *person.Person, sender, sub, unsub bool) {
				if sender {
					senderCount++
				}
				if sub && !unsub {
					subCount++
				}
				if unsub {
					unsubCount++
				}
			})
			row = grid.E("div class=listlistRow")
			row.E("a href=/admin/lists/%d up-layer=new up-size=grow up-dismissable=key up-history=false>%s", l.ID, name)
			if subCount != 0 {
				row.E("a href=/admin/lists/%d/sub up-layer=new up-size=grow up-history=false>%d", l.ID, subCount)
			} else {
				row.E("div>0")
			}
			if unsubCount != 0 {
				row.E("a href=/admin/lists/%d/unsub up-layer=new up-size=grow up-history=false>%d", l.ID, unsubCount)
			} else {
				row.E("div>0")
			}
			if senderCount != 0 {
				row.E("a href=/admin/lists/%d/sender up-layer=new up-size=grow up-history=false>%d", l.ID, senderCount)
			} else {
				row.E("div>0")
			}
		})
		main.E("div class=listlistButtons").
			E("a href=/admin/lists/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add List")
	})
}
