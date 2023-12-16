package peoplelist

import (
	"sort"
	"strings"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

type personData struct {
	*person.Person
	viewLevel  person.ViewLevel
	callPrefix string
	callSuffix string
	role1      string
	roles      []string
	rolePrio   uint
}

// Handle handles GET and POST /people requests.
func Handle(r *request.Request) {
	var (
		user        *person.Person
		opts        ui.PageOpts
		focus       *role.Role
		roleOptions []*role.Role
		title       string
		currsort    string
	)
	if user = auth.SessionUser(r, person.CanViewViewerFields, true); user == nil {
		return
	}
	// Figure out what role to focus on, if any.
	if _, ok := r.Form["role"]; ok {
		focus = role.WithID(r, role.ID(util.ParseID(r.FormValue("role"))), role.FOrg|role.FName)
	} else if rid := state.GetFocusRole(r); rid != 0 {
		focus = role.WithID(r, rid, role.FOrg|role.FName)
	}
	if focus != nil && !user.HasPrivLevel(focus.Org(), enum.PrivMember) {
		focus = nil
	}
	// Get the list of roles the caller is allowed to focus on.
	role.All(r, role.FID|role.FName|role.FOrg|role.FFlags, func(rl *role.Role) {
		if user.HasPrivLevel(rl.Org(), enum.PrivMember) && (rl.Flags()&role.Filter != 0 || rl.ID() == focus.ID()) {
			clone := *rl
			roleOptions = append(roleOptions, &clone)
		}
	})
	// If there's only one such role, focus on it.
	if len(roleOptions) == 1 {
		focus = roleOptions[0]
	}
	if focus != nil {
		title = focus.Name()
		state.SetFocusRole(r, focus.ID())
	} else {
		title = r.LangString("People", "Personas")
		state.SetFocusRole(r, 0)
	}
	// Figure out the sort order.
	currsort = r.FormValue("sort")
	if currsort != "name" && currsort != "call" && currsort != "suffix" {
		currsort = "priority"
	}
	if (focus == nil || focus.Org() != enum.OrgSARES) && currsort != "name" {
		currsort = "priority"
	}
	// Show the page.
	opts = ui.PageOpts{
		Title:    title,
		MenuItem: "people",
		Tabs: []ui.PageTab{
			{Name: r.LangString("List", "Lista"), URL: "/people", Target: "main", Active: true},
			{Name: r.LangString("Map", "Mapa"), URL: "/people/map", Target: "main"},
		},
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		main.A("class=peoplelist")
		listControls(r, user, main, focus, roleOptions, currsort)
		people := getPeople(r, user, focus, currsort == "suffix")
		sort.Slice(people, func(i, j int) bool {
			return personDataLess(people[i], people[j], currsort)
		})
		grid := main.E("div class=peoplelistGrid",
			currsort == "suffix", "class=peoplelistGrid-callsuffix",
			currsort != "suffix" && focus != nil && focus.Org() == enum.OrgSARES, "class=peoplelistGrid-callsign",
		)
		for _, p := range people {
			showPerson(r, grid, p, focus != nil && focus.Org() == enum.OrgSARES, currsort == "suffix")
		}
		if len(people) == 1 {
			main.E("div class=peoplelistCount").R(r.LangString("1 person listed.", "1 persona en la lista."))
		} else {
			main.E("div class=peoplelistCount>%d ", len(people)).R(r.LangString("people listed.", "personas en la lista."))
		}
		if user.HasPrivLevel(0, enum.PrivLeader) {
			// TODO add user
		}
	})
}

// listControls displays the controls bar (focus choice and sort order) at the
// top of the list.
func listControls(r *request.Request, user *person.Person, main *htmlb.Element, focus *role.Role, roleOptions []*role.Role, currsort string) {
	var nextsort string

	form := main.E("form class=peoplelistForm method=POST")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	// If they have more than one choice, give them a select box; otherwise
	// just add the choice as a hidden element.
	if len(roleOptions) > 1 {
		sort.Slice(roleOptions, func(i, j int) bool { return roleOptions[i].Name() < roleOptions[j].Name() })
		sel := form.E("select id=peoplelistRole name=role")
		sel.E("option value=0", focus == nil, "selected").R(r.LangString("(all)", "(todos)"))
		for _, role := range roleOptions {
			sel.E("option value=%d", role.ID(), focus != nil && focus.ID() == role.ID(), "selected").T(role.Name())
		}
	} else if focus != nil {
		form.E("input type=hidden name=role value=%d", focus.ID())
	}
	// Given the current sort order, determine the next sort order in
	// sequence, and emit a button that will change to that sort order.
	switch currsort {
	case "priority":
		nextsort = "name"
	case "name":
		if focus != nil && focus.Org() == enum.OrgSARES {
			nextsort = "suffix"
		} else {
			nextsort = "priority"
		}
	case "suffix":
		nextsort = "priority"
	}
	form.E("button type=submit name=sort class='sbtn sbtn-small sbtn-secondary' value=%s", nextsort).R(r.LangString("Sort", "Ordenar"))
	// Emit the current sort order.  If they click the button, this will be
	// ignored because it will be the second value of "sort" in the form.
	// Otherwise, this will retain the current sort order when changing
	// focus.
	form.E("input type=hidden name=sort value=%s", currsort)
}

// getPeople gets the set of people to be shown (i.e., those who are visible to
// the calling user and who have the focus role).  It also computes information
// about the person that is needed for both sorting and display.  (Information
// needed only for display is computed in showPerson.)
func getPeople(r *request.Request, user *person.Person, focus *role.Role, splitCallSign bool) (people []*personData) {
	const personFields = person.FID | person.FSortName | person.FInformalName | person.FCallSign | person.FHomePhone | person.FCellPhone | person.FWorkPhone | person.FEmail | person.FEmail2 | person.CanViewTargetFields
	person.All(r, personFields, func(p *person.Person) {
		viewLevel := user.CanView(p)
		if viewLevel == person.ViewNone {
			return
		}
		if focus != nil {
			if held, _ := personrole.PersonHasRole(r, p.ID(), focus.ID()); !held {
				return
			}
		}
		clone := *p
		pd := &personData{Person: &clone, viewLevel: viewLevel}
		// Compute which role title to show for the person.
		pd.roles, pd.rolePrio = rolesToShow(r, p, focus)
		switch len(pd.roles) {
		case 0:
			break
		case 1:
			pd.role1 = pd.roles[0]
		default:
			pd.role1 = pd.roles[0] + ", ..."
		}
		// Split the call sign into prefix and suffix.
		idx := strings.IndexAny(p.CallSign(), "0123456789")
		if idx >= 0 {
			pd.callPrefix, pd.callSuffix = p.CallSign()[:idx+1], p.CallSign()[idx+1:]
		}
		people = append(people, pd)
	})
	return people
}

// rolesToShow returns a list of role titles to show for a particular person,
// when viewing a particular focus group.  It also returns the highest role
// priority (numerically lowest number) of any of those roles.
func rolesToShow(r *request.Request, p *person.Person, focus *role.Role) (list []string, minPrio uint) {
	// We want to show all roles that:
	//   - the person holds directly
	//   - have titles
	//   - imply, directly or indirectly, the focus role if any
	//   - are not the focus role itself
	// If there is more than one such role, they are sorted by priority.
	var (
		roles   []*role.Role
		include map[role.ID]struct{}
	)
	minPrio = 99999
	if focus != nil {
		include = role.AllThatImply(r, focus.ID())
	}
	personrole.RolesForPerson(r, p.ID(), role.FID|role.FTitle|role.FPriority, func(rl *role.Role, explicit bool) {
		if !explicit || rl.ID() == focus.ID() || rl.Title() == "" {
			return
		}
		if include != nil {
			if _, ok := include[rl.ID()]; !ok {
				return
			}
		}
		clone := *rl
		roles = append(roles, &clone)
		if rl.Priority() < minPrio {
			minPrio = rl.Priority()
		}
	})
	if len(roles) == 0 {
		return nil, minPrio
	}
	sort.Slice(roles, func(i, j int) bool {
		return roles[i].Priority() < roles[j].Priority()
	})
	list = make([]string, len(roles))
	for i, rl := range roles {
		list[i] = rl.Title()
	}
	return list, minPrio
}

// personDataLess is used in sorting the list of people to display.
func personDataLess(a, b *personData, currsort string) bool {
	switch currsort {
	case "priority":
		if a.rolePrio != b.rolePrio {
			return a.rolePrio < b.rolePrio
		}
	case "suffix":
		if a.callSuffix != b.callSuffix {
			return a.callSuffix < b.callSuffix
		}
	}
	return a.SortName() < b.SortName()
}

// showPerson displays a single person.
func showPerson(r *request.Request, grid *htmlb.Element, p *personData, showCall, splitCall bool) {
	// Show the call sign column(s) if needed.
	if splitCall {
		grid.E("div class=peoplelistPersonCall1", p.role1 != "", "peoplelistPerson-withrole").T(p.callPrefix)
		grid.E("div class=peoplelistPersonCall2", p.role1 != "", "peoplelistPerson-withrole").T(p.callSuffix)
	} else if showCall {
		grid.E("div class=peoplelistPersonCall1", p.role1 != "", "peoplelistPerson-withrole").T(p.CallSign())
	}
	// Show the name and principal role.
	nr := grid.E("div class=peoplelistPersonNameroles")
	nr.E("div class=peoplelistPersonName").
		E("a class=peoplelistPersonName href=/people/%d up-target=.pageCanvas>%s", p.ID(), p.SortName())
	nr.E("div class=peoplelistPersonRoles>%s", p.role1)
	// Show the email and phone number if allowed.
	ep := grid.E("div class=peoplelistPersonEmailphone")
	e := ep.E("div class=peoplelistPersonEmail")
	if p.viewLevel >= person.ViewWorkContact {
		if p.Email() != "" {
			e.E("a href=mailto:%s target=_blank>%s", p.Email(), p.Email())
		} else if p.Email2() != "" {
			e.E("a href=mailto:%s target=_blank>%s", p.Email2(), p.Email2())
		}
	}
	ph := ep.E("div class=peoplelistPersonPhone")
	if p.viewLevel == person.ViewFull {
		if p.CellPhone() != "" {
			ph.T(p.CellPhone())
		} else if p.HomePhone() != "" {
			ph.T(p.HomePhone())
		} else if p.WorkPhone() != "" {
			ph.T(p.WorkPhone())
		}
	} else if p.viewLevel == person.ViewWorkContact && p.WorkPhone() != "" {
		ph.T(p.WorkPhone())
	}
	// Show the details popup.
	det := grid.E("div class=peoplelistPersonDetails")
	det.E("s-icon icon=info")
	showPersonDetails(r, det.E("div class=peoplelistDetails style=display:none"), p)
}

// showPersonDetails renders the details popup for a person.
func showPersonDetails(r *request.Request, box *htmlb.Element, p *personData) {
	// Show the name and call sign.
	n := box.E("div class=peoplelistDetailsName>%s", p.InformalName())
	if p.CallSign() != "" {
		n.E("span class=peoplelistDetailsCall>%s", p.CallSign())
	}
	// Show all of the person's roles.
	rb := box.E("div class=peoplelistDetailsRoles")
	for _, role := range p.roles {
		rb.E("div>%s", role)
	}
	// If the caller can't see the person's contact info, show no more.
	if p.viewLevel <= person.ViewNoContact {
		return
	}
	// Show the email addresses.
	if p.Email() != "" || p.Email2() != "" {
		e := box.E("div class=peoplelistDetailsEmails")
		if p.Email() != "" {
			e.E("div class=peoplelistDetailsIconline").
				E("div class=peoplelistDetailsEmail>%s", p.Email()).P().
				E("div class=peoplelistDetailsIcons").
				E("div class=peoplelistDetailsIcon").
				E("a href=mailto:%s target=_blank>", p.Email()).
				E("s-icon icon=email")
		}
		if p.Email2() != "" {
			e.E("div class=peoplelistDetailsIconline").
				E("div class=peoplelistDetailsEmail>%s", p.Email2()).P().
				E("div class=peoplelistDetailsIcons").
				E("div class=peoplelistDetailsIcon").
				E("a href=mailto:%s target=_blank>", p.Email2()).
				E("s-icon icon=email")
		}
	}
	// Show the phone numbers.
	if (p.viewLevel == person.ViewFull && (p.CellPhone() != "" || p.HomePhone() != "")) || p.WorkPhone() != "" {
		ph := box.E("div class=peoplelistDetailsPhones")
		if p.viewLevel == person.ViewFull && p.CellPhone() != "" {
			ph.E("div class=peoplelistDetailsIconline").
				E("div class=peoplelistDetailsPhone>%s (%s)", p.CellPhone(), r.LangString("cell", "móvil")).P().
				E("div class=peoplelistDetailsIcons").
				E("div class=peoplelistDetailsIcon").
				E("a href=sms:%s target=_blank>", p.CellPhone()).
				E("s-icon icon=message").P().P().P().
				E("div class=peoplelistDetailsIcon").
				E("a href=tel:%s target=_blank>", p.CellPhone()).
				E("s-icon icon=phone")
		}
		if p.viewLevel == person.ViewFull && p.HomePhone() != "" {
			ph.E("div class=peoplelistDetailsIconline").
				E("div class=peoplelistDetailsPhone>%s (%s)", p.HomePhone(), r.LangString("home", "casa")).P().
				E("div class=peoplelistDetailsIcons").
				E("div class=peoplelistDetailsIcon").
				E("a href=tel:%s target=_blank>", p.HomePhone()).
				E("s-icon icon=phone")
		}
		if p.WorkPhone() != "" {
			ph.E("div class=peoplelistDetailsIconline").
				E("div class=peoplelistDetailsPhone>%s (%s)", p.WorkPhone(), r.LangString("work", "trabajo")).P().
				E("div class=peoplelistDetailsIcons").
				E("div class=peoplelistDetailsIcon").
				E("a href=tel:%s target=_blank>", p.WorkPhone()).
				E("s-icon icon=phone")
		}
	}
}
