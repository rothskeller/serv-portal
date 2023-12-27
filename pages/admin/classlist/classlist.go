package classlist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/classes requests.
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
		Title:    "Classes",
		MenuItem: "admin",
		Tabs: []ui.PageTab{
			{Name: "Roles", URL: "/admin/roles", Target: "main"},
			{Name: "Lists", URL: "/admin/lists", Target: "main"},
			{Name: "Classes", URL: "/admin/classes", Target: "main", Active: true},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		main.E("div class=classlistButtons").
			E("a href=/admin/classes/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add Class")
		grid := main.E("div class=classlistGrid")
		row := grid.E("div class=classlistHeading")
		row.E("div>Type")
		row.E("div>Date")
		row.E("div>Description")
		class.All(r, class.FID|class.FType|class.FStart|class.FEnDesc, func(c *class.Class) {
			row = grid.E("div class=classlistRow")
			row.E("div>%s", c.Type().String())
			row.E("div").E("a href=/admin/classes/%d up-layer=new up-size=grow up-dismissable=key up-history=false>%s", c.ID(), c.Start())
			row.E("div class=classlistDesc").R(c.EnDesc())
		})
	})
}
