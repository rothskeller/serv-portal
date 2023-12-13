package rolelist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/roles requests.
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
		Title:    "Roles",
		MenuItem: "admin",
		Tabs: []ui.PageTab{
			{Name: "Roles", URL: "/admin/roles", Target: "main", Active: true},
			{Name: "Lists", URL: "/admin/lists", Target: "main"},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		grid := main.E("div class=rolelistGrid")
		row := grid.E("div class=rolelistHeading")
		row.E("div>Org")
		row.E("div>Priv")
		row.E("div>Role")
		role.All(r, role.FID|role.FOrg|role.FPrivLevel|role.FName|role.FFlags, func(rl *role.Role) {
			row = grid.E("div class=rolelistRow",
				rl.Flags()&role.Archived != 0, "class=archived",
				rl.Flags()&role.Filter != 0, "class=filter")
			row.E("div>%s", rl.Org().String())
			row.E("div>%s", privNames[rl.PrivLevel()])
			ndiv := row.E("div")
			ndiv.E("a href=/admin/roles/%d up-layer=new up-size=grow up-dismissable=key up-history=false>%s", rl.ID(), rl.Name())
			ndiv.E("a href=/people?role=%d class=rolelistCount up-target=.pageCanvas>%d", rl.ID(), personrole.PeopleCountForRole(r, rl.ID()))
		})
		main.E("div class=rolelistButtons").
			E("a href=/admin/roles/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add Role")
	})
}

var privNames = map[enum.PrivLevel]string{
	0:                "—",
	enum.PrivStudent: "Student",
	enum.PrivMember:  "Member",
	enum.PrivLeader:  "Leader",
	enum.PrivMaster:  "Master",
}
