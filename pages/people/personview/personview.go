package personview

import (
	"fmt"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// PersonFields are the fields that must be retrieved in order to display the
// entire PersonView page.
const PersonFields = person.FID | person.FInformalName | person.FPrivLevels | namesPersonFields | contactPersonFields | statusPersonFields | notesPersonFields | subscriptionsPersonFields

// Get handles GET /people/${id} requests.
func Get(r *request.Request, idstr string) {
	var (
		user               *person.Person
		p                  *person.Person
		canView            bool
		canViewContactInfo bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), PersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if canView, canViewContactInfo = canViewPerson(user, p); !canView {
		errpage.Forbidden(r, user)
		return
	}
	Render(r, user, p, canViewContactInfo, "")
}

// Render renders the person view page, or a particular section of it.  It is
// called by Get, above, and also by the edit dialogs after accepting a change
// to a person.
func Render(r *request.Request, user, p *person.Person, canViewContactInfo bool, section string) {
	// Show the page.
	opts := ui.PageOpts{
		Title:    p.InformalName(),
		MenuItem: "people",
		Tabs: []ui.PageTab{
			{Name: "List", URL: "/people", Target: ".pageCanvas"},
			{Name: "Map", URL: "/people/map", Target: ".pageCanvas"},
			{Name: "Details", URL: fmt.Sprintf("/people/%d", p.ID()), Target: "main", Active: true},
		},
	}
	if user.ID() == p.ID() || user.HasPrivLevel(0, enum.PrivLeader) {
		opts.Tabs = append(opts.Tabs, ui.PageTab{Name: "Activity", URL: fmt.Sprintf("/people/%d/activity/current", p.ID()), Target: "main"})
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		main.A("class=personview")
		if section == "" || section == "names" {
			showNames(main, user, p)
			main.E("div class=personviewSpacer")
		}
		if section == "" {
			showRoles(r, main, user, p)
		}
		if section == "" || section == "contact" {
			showContact(main, user, p, canViewContactInfo)
		}
		if section == "" || section == "status" {
			showStatus(r, main, user, p)
		}
		if section == "" || section == "notes" {
			showNotes(r, main, user, p, canViewContactInfo)
		}
		if section == "" || section == "subscriptions" {
			showSubscriptions(r, main, user, p)
		}
		if section == "" || section == "password" {
			showPassword(r, main, user, p)
		}
	})
}

// canViewPerson returns whether the specified viewer is allowed to see the
// specified viewee.  It returns two flags: one indicating whether the viewee
// can be seen in the roster at all; the other indicating whether the viewee's
// contact information can be seen.
func canViewPerson(viewer, viewee *person.Person) (roster, contact bool) {
	if viewer.ID() == viewee.ID() || viewer.HasPrivLevel(0, enum.PrivLeader) {
		return true, true
	}
	for o, op := range viewer.PrivLevels() {
		if op < enum.PrivMember {
			continue
		}
		if viewee.PrivLevels()[o] == 0 {
			continue
		}
		roster = true
		if op < enum.PrivMember {
			continue
		}
		if op == enum.PrivMember && viewee.PrivLevels()[o] < enum.PrivLeader && !enum.Org(o).MembersCanViewContactInfo() {
			continue
		}
		contact = true
		return
	}
	return
}
