package redirlist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/redirect"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/redirects requests.
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
		Title:    "RedirectS",
		MenuItem: "admin",
		Tabs: []ui.PageTab{
			{Name: "Roles", URL: "/admin/roles", Target: "main"},
			{Name: "Lists", URL: "/admin/lists", Target: "main"},
			{Name: "Venues", URL: "/admin/venues", Target: "main"},
			{Name: "Classes", URL: "/admin/classes", Target: "main"},
			{Name: "Redirects", URL: "/admin/redirects", Target: "main", Active: true},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		grid := main.E("div class=redirlistGrid")
		row := grid.E("div class=redirlistHeading")
		row.E("div>Entry")
		row.E("div>Target")
		redirect.All(r, func(rd *redirect.Redirect) {
			row = grid.E("div class=redirlistRow")
			row.E("div").E("a href=/admin/redirects/%d up-layer=new up-size=grow up-dismissable=key up-history=false", rd.ID).T(rd.Entry)
			row.E("div").T(rd.Target)
		})
		main.E("div class=redirlistButtons").
			E("a href=/admin/redirects/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add Redirect")
	})
}
