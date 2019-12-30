package util

import (
	"html/template"
	"serv.rothskeller.net/portal/model"
)

// A Page structure contains the data needed to render the standard page layout
// for a page, plus the data to be passed into the page body template.
type Page struct {
	Title     string
	ShowMenu  bool
	ShowTeams bool
	MenuItem  string
	MyID      model.PersonID
	MyName    string
	BodyData  interface{}
}

// RenderPage renders a page in the standard page layout.  The supplied "page"
// contains the layout variables and the context to be supplied to the page
// body template.  The supplied "tmpl" is the page body template, which must
// define a template named "body" and any other named templates used by it.
func RenderPage(r *Request, page *Page, tmpl *template.Template) {
	tmpl = template.Must(tmpl.Parse(layoutTemplate))
	page.ShowMenu = r.Person != nil
	page.ShowTeams = r.Person.IsWebmaster()
	if r.Person != nil {
		page.MyID = r.Person.ID
		page.MyName = r.Person.FirstName + " " + r.Person.LastName
	}
	if err := tmpl.ExecuteTemplate(r, "layout", page); err != nil {
		panic(err)
	}
}

const layoutTemplate = `{{ define "layout" -}}
<!DOCTYPE html><html lang="en"><head>
  <meta charset="utf-8">
  <meta http-equiv="X-UA-Compatible" content="IE=edge">
  <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">
  <title>{{ .Title }}</title>
  <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
  <link rel="stylesheet" href="/portal.css">
  <script type="text/javascript" src="/portal.js"></script>
</head>
<body class="mouse">
  <div id="layout-heading">
    <div id="layout-menu-trigger-box">
      <svg id="layout-menu-trigger" class="bi bi-list" viewBox="0 0 20 20" fill="currentColor" xmlns="http://www.w3.org/2000/svg">
        <path fill-rule="evenodd" d="M4.5 13.5A.5.5 0 015 13h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 9h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5zm0-4A.5.5 0 015 5h10a.5.5 0 010 1H5a.5.5 0 01-.5-.5z" clip-rule="evenodd"/>
      </svg>
    </div>
    <div id="layout-titlebox">
      <div id="layout-title">{{ .Title }}</div>
    </div>
    <div id="layout-menu-spacer"></div>
  </div>
  <div id="layout-main">
    {{- if .ShowMenu }}
      <div id="layout-menu">
        <div id="layout-menu-welcome">Welcome<br><b>{{ .MyName }}</b></div>
        <a class="layout-menu-item{{ if eq .MenuItem "events" }} layout-menu-item-active{{ end }}" href="/events">Events</a>
        <a class="layout-menu-item{{ if eq .MenuItem "people" }} layout-menu-item-active{{ end }}" href="/people">People</a>
	{{- if .ShowTeams }}
          <a class="layout-menu-item{{ if eq .MenuItem "teams" }} layout-menu-item-active{{ end }}" href="/teams">Teams</a>
	{{ end }}
        <a class="layout-menu-item{{ if eq .MenuItem "reports" }} layout-menu-item-active{{ end }}" href="/reports">Reports</a>
        <a class="layout-menu-item" href="/people/{{ .MyID }}">Profile</a>
        <a class="layout-menu-item" href="/logout">Logout</a>
      </div>
    {{ end }}
    <div id="layout-content">
      {{ block "body" .BodyData }}{{ end }}
    </div>
  </div>
  <script src="https://code.jquery.com/jquery-3.4.1.slim.min.js" integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
  <script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js" integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo" crossorigin="anonymous"></script>
  <script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js" integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6" crossorigin="anonymous"></script>
</body></html>{{ end }}`
