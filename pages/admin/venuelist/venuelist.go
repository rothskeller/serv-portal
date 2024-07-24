package venuelist

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/venues requests.
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
		Title:    "Venues",
		MenuItem: "admin",
		Tabs: []ui.PageTab{
			{Name: "Roles", URL: "/admin/roles", Target: "main"},
			{Name: "Lists", URL: "/admin/lists", Target: "main"},
			{Name: "Venues", URL: "/admin/venues", Target: "main", Active: true},
			{Name: "Classes", URL: "/admin/classes", Target: "main"},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		grid := main.E("div class=venuelistGrid")
		row := grid.E("div class=venuelistHeading")
		row.E("div>Name")
		row.E("div>Map")
		row.E("div>Flags")
		venue.All(r, venue.FID|venue.FName|venue.FURL|venue.FFlags, func(v *venue.Venue) {
			row = grid.E("div class=venuelistRow")
			row.E("div").E("a href=/admin/venues/%d up-layer=new up-size=grow up-dismissable=key up-history=false", v.ID()).T(v.Name())
			if v.URL() != "" {
				row.E("div").E("a href=%s target=_blank", v.URL()).R("map")
			} else {
				row.E("div")
			}
			if v.Flags()&venue.CanOverlap != 0 {
				row.E("div>can overlap")
			} else {
				row.E("div")
			}
		})
		main.E("div class=venuelistButtons").
			E("a href=/admin/venues/NEW up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary'>Add Venue")
	})
}
