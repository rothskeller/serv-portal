package ui

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/state"
)

// A PageOpts structure gives the details needed to render a page.
type PageOpts struct {
	// Request is the request being processed.  It's required.
	Request *request.Request
	// User is the person making the request, or nil for an unauthenticated
	// user.
	User *person.Person
	// Title is the string to be placed in the browser title bar.  If it is
	// empty, "Sunnyvale SERV" is used.
	Title string
	// Banner is the string to be placed in the content banner.  If it is
	// empty, Title is used.
	Banner string
	// StatusCode is the status code to send to the client.  If it is zero,
	// 200 OK is sent.
	StatusCode int
	// MenuItem is the ID of the active menu item.
	MenuItem string
	// Tabs is the list of tabs to be shown in the tab bar on the page.  If
	// the list is empty, no tab bar is shown.
	Tabs []PageTab
	// NoHome is a flag indicating that the home page icon should not be
	// displayed.  (It is true when showing the home page itself.)
	NoHome bool
}

// A PageTab describes one tab on a page with a tab bar.
type PageTab struct {
	Name   string
	URL    string
	Target string
	Alias  string
	Active bool
	Hide   bool
}

// Page displays a page.  It calls the supplied function with the <main> element
// into which the page contents should be rendered.
func Page(r *request.Request, user *person.Person, opts PageOpts, fn func(*htmlb.Element)) {
	r.HTMLNoCache()
	if opts.StatusCode != 0 {
		r.WriteHeader(opts.StatusCode)
	}

	html := htmlb.HTML(r).Attr("lang=%s translate=no", r.Language)
	defer html.Close()
	pageHead(html, opts.Title)

	body := html.E("body class=page", user == nil, "class=page-noMenu")
	pageTitle(r, body, user, opts.Banner, opts.Title, opts.NoHome)
	if user != nil {
		pageMenu(body, r, user, opts.MenuItem)
	}

	if len(opts.Tabs) != 0 {
		page := body.E("div class=pageCanvas").E("div class=pageTabbed")
		tabs := page.E("nav class=pageTabBar").E("ul class=pageTabs up-nav")
		for _, tab := range opts.Tabs {
			if !tab.Hide {
				tabs.E("li class=pageTab").
					E("a href=%s up-target=%s class=pageTabLink", tab.URL, tab.Target,
						tab.Alias != "", "up-alias=%s", tab.Alias,
						tab.Active, "class=up-current").T(tab.Name)
			}
		}
		fn(page.E("main class=pageTabContent"))
	} else {
		fn(body.E("main class=pageCanvas"))
	}
}

func pageHead(h *htmlb.Element, title string) {
	h.E("meta charset=utf-8")
	h.E("meta name=viewport content='width=device-width, initial-scale=1.0'")
	if title != "" {
		h.E("title").T(title).R(" - Sunnyvale SERV")
	} else {
		h.E("title>Sunnyvale SERV")
	}
	h.E("link rel=stylesheet href=%s", AssetURL("styles.css"))
	h.E("script src=%s", AssetURL("script.js"))
	h.E("script").
		R("window.algoliaApplicationID='").R(config.Get("algoliaApplicationID")).R("';\n").
		R("window.algoliaSearchKey='").R(config.Get("algoliaSearchKey")).R("';\n").
		R("window.algoliaIndex='").R(config.Get("algoliaIndex")).R("';\n")
}

func pageTitle(r *request.Request, h *htmlb.Element, user *person.Person, banner, title string, noHome bool) {
	h = h.E("div class=pageTitle")
	if user != nil {
		h.E("div id=pageMenuTrigger class=pageTitleMenu").E("s-icon icon=menu")
	} else if !noHome {
		h.E("div class=pageTitleMenu").E("a href=/ up-target=.pageCanvas").E("s-icon icon=home")
	}
	switch {
	case banner != "":
		h.E("div class=pageTitleText up-hungry").T(banner)
	case title != "":
		h.E("div class=pageTitleText up-hungry").T(title)
	default:
		h.E("div class=pageTitleText up-hungry>Sunnyvale SERV")
	}
	h.E("div class=pageTitleSearch").E("a class=nolink href=/search").E("s-icon icon=search")
}

func pageMenu(h *htmlb.Element, r *request.Request, user *person.Person, menuItem string) {
	h = h.E("nav class=pageMenu")
	welcome := h.E("div class=pageMenuWelcome").R(r.Loc("Welcome") + "<br>").E("b").T(user.InformalName())
	if r.Language == "es" {
		welcome.E("a href=/en%s class='sbtn sbtn-xsmall sbtn-primary pageMenuLangSel' up-target=body", r.Path).R("View in English")
	} else {
		welcome.E("a href=/es%s class='sbtn sbtn-xsmall sbtn-primary pageMenuLangSel' up-target=body", r.Path).R("Ver en español")
	}
	ul := h.E("ul class=pageMenuList up-nav")
	ul.E("li").E("a href=/ up-target=.pageCanvas class=pageMenuItem", menuItem == "home", "class=up-current").R(r.Loc("Home[PAGE]"))
	if user.HasPrivLevel(0, enum.PrivStudent) {
		ul.E("li").E("a href=%s up-target=.pageCanvas up-alias=/events/* class=pageMenuItem", state.GetEventsURL(r),
			menuItem == "events", "class=up-current").R(r.Loc("Events"))
	}
	ul.E("li").E("a href=/classes up-target=.pageCanvas up-alias='/clases /pep /cert-basic' class=pageMenuItem",
		menuItem == "classes", "class=up-current").R(r.Loc("Classes"))
	if user.HasPrivLevel(0, enum.PrivStudent) {
		ul.E("li").E("a href=/people up-target=.pageCanvas up-alias='/people/* -/people/%d -/people/%d/*' class=pageMenuItem", user.ID(), user.ID(),
			menuItem == "people", "class=up-current").R(r.Loc("People"))
	}
	ul.E("li").E("a href=/files up-target=.pageCanvas up-alias=/files/* class=pageMenuItem",
		menuItem == "files", "class=up-current").R(r.Loc("Files"))
	if user.HasPrivLevel(0, enum.PrivLeader) {
		ul.E("li").E("a href=/reports/attendance up-target=.pageCanvas up-alias=/reports/* class=pageMenuItem",
			menuItem == "reports", "class=up-current").R("Reports")
	}
	if listperson.CanSendText(r, user.ID()) {
		ul.E("li").E("a href=/texts up-target=.pageCanvas up-alias=/texts/* class=pageMenuItem",
			menuItem == "texts", "class=up-current").R("Texts")
	}
	if user.IsWebmaster() {
		ul.E("li").E("a href=/admin/roles up-target=.pageCanvas up-alias=/admin/* class=pageMenuItem",
			menuItem == "admin", "class=up-current").R("Admin")
	}
	if user.ID() != person.AdminID {
		ul.E("li").E("a href=/people/%d up-target=.pageCanvas up-alias=/people/%d/* class=pageMenuItem", user.ID(), user.ID(),
			menuItem == "profile", "class=up-current").R(r.Loc("Profile"))
	}
	ul.E("li").E("a href=/logout up-target=body class=pageMenuItem").R(r.Loc("Logout"))
	h.E("a href=/about up-target=.pageCanvas class=pageMenuAbout").R(r.Loc("Web Site Info"))
}
