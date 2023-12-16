package search

import (
	"net/url"
	"path"
	"path/filepath"
	"strings"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/search"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /search requests.
func Handle(r *request.Request) {
	var (
		user      *person.Person
		query     string
		results   []any
		searchErr error
		seen      bool
		seenAny   bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if query = r.FormValue("q"); query != "" {
		results, searchErr = search.Search(r, query)
	}
	r.HTMLNoCache()
	ui.Page(r, user, ui.PageOpts{Title: r.LangString("Search", "Buscar")}, func(main *htmlb.Element) {
		form := main.E("form class=searchForm method=GET")
		form.E("input type=search name=q class=formInput value=%s autofocus", query)
		form.E("input type=submit class='sbtn sbtn-primary' value=%s", r.LangString("Search", "Buscar"))
		if searchErr != nil {
			main.E("div class=searchErr>%s", searchErr.Error())
		}
		rdiv := main.E("div class=searchResults")
		for _, result := range results {
			if e, ok := result.(*event.Event); ok {
				if !seen {
					rdiv.E("div class=searchHeading").R(r.LangString("Events", "Eventos"))
					seen, seenAny = true, true
				}
				rdiv.E("div class=searchResult").
					E("a href=/events/%d up-target=.pageCanvas>%s %s", e.ID(), e.Start()[:10], e.Name())
			}
		}
		seen = false
		for _, result := range results {
			if p, ok := result.(*person.Person); ok {
				if !personVisibleToUser(user, p) {
					continue
				}
				if !seen {
					rdiv.E("div class=searchHeading").R(r.LangString("People", "Personas"))
					seen, seenAny = true, true
				}
				rdiv.E("div class=searchResult").
					E("a href=/people/%d up-target=.pageCanvas>%s", p.ID(), p.SortName())
			}
		}
		seen = false
		for _, result := range results {
			if f, ok := result.(*folder.Folder); ok {
				if !user.HasPrivLevel(f.Viewer()) {
					continue
				}
				if !seen {
					rdiv.E("div class=searchHeading").R(r.LangString("Folders", "Carpetas"))
					seen, seenAny = true, true
				}
				rdiv.E("div class=searchResult").
					E("a href=%s up-target=.pageCanvas>%s", f.Path(r), f.Name())
			}
		}
		seen = false
		for _, result := range results {
			if d, ok := result.(*document.Document); ok {
				f := folder.WithID(r, d.Folder, folder.FID|folder.FViewer|folder.FURLName|folder.FName|folder.FParent)
				if !user.HasPrivLevel(f.Viewer()) {
					continue
				}
				if !seen {
					rdiv.E("div class=searchHeading").R(r.LangString("Documents", "Archivos"))
					seen, seenAny = true, true
				}
				sr := rdiv.E("div class=searchResult")
				if d.URL != "" {
					sr.E("a href=%s target=_blank>%s", d.URL, d.Name)
				} else {
					var newtab bool

					switch strings.ToLower(filepath.Ext(d.Name)) {
					case ".pdf", ".png", ".jpeg", ".jpg":
						newtab = true
					}
					sr.E("a href=%s", path.Join(f.Path(r), url.PathEscape(d.Name)), newtab, "target=_blank").
						T(d.Name)
				}
				sr.E("span class=searchContext> %s ", r.LangString("in folder", "en carpeta")).E("a href=%s>%s", f.Path(r), f.Name())
			}
		}
		seen = false
		for _, result := range results {
			if v, ok := result.(*venue.Venue); ok {
				if !seen {
					rdiv.E("div class=searchHeading").R(r.LangString("Venues", "Sitios"))
					seen, seenAny = true, true
				}
				if v.URL() != "" {
					rdiv.E("div class=searchResult").
						E("a href=%s target=_blank>%s", v.URL(), v.Name())
				} else {
					rdiv.E("div class=searchResult>%s", v.Name())
				}
			}
		}
		seen = false
		for _, result := range results {
			if rl, ok := result.(*role.Role); ok {
				if !user.HasPrivLevel(rl.Org(), enum.PrivStudent) && !user.HasPrivLevel(enum.OrgAdmin, enum.PrivMember) {
					continue
				}
				if !seen {
					rdiv.E("div class=searchHeading>Roles")
					seen, seenAny = true, true
				}
				rdiv.E("div class=searchResult").
					E("a href=/people?role=%d up-target=.pageCanvas>%s", rl.ID(), rl.Name())
			}
		}
		if !seenAny {
			rdiv.E("div class=searchHeading").R(r.LangString("Nothing matched your search.", "No se encontró nada en su búsqueda."))
		}
		main.E("div class=searchLogo>Search provided by").E("img src=%s", ui.AssetURL("algolia-logo.png"))
	})
}

func personVisibleToUser(user, p *person.Person) bool {
	if user.HasPrivLevel(enum.OrgAdmin, enum.PrivMember) {
		return true
	}
	for _, org := range enum.AllOrgs {
		if user.HasPrivLevel(org, enum.PrivStudent) && p.HasPrivLevel(org, enum.PrivMember) {
			return true
		}
	}
	return false
}
