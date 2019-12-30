package person

import (
	"html/template"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// ListPeople handles GET /people requests.
func ListPeople(r *util.Request) error {
	var (
		lpd = listPeopleData{
			TeamID: model.TeamID(util.ParseID(r.FormValue("team"))),
		}
		focus   *model.Team
		canView bool
	)
	switch {
	case lpd.TeamID < 0:
		lpd.TeamID = 0
	case lpd.TeamID > 0:
		focus = r.Tx.FetchTeam(lpd.TeamID)
	}
	lpd.People = r.Tx.FetchPeople()
	j := 0
	for _, p := range lpd.People {
		if !r.Person.CanViewPerson(p) {
			continue
		}
		canView = true
		if focus == nil || p.IsMember(focus) {
			lpd.People[j] = p
			j++
		}
	}
	if !canView {
		return util.Forbidden
	}
	lpd.People = lpd.People[:j]
	lpd.Teams = r.Person.ViewableTeams()
	if len(lpd.Teams) == 0 {
		lpd.Teams = nil
	} else {
		lpd.Teams = append(lpd.Teams, nil)
		copy(lpd.Teams[1:], lpd.Teams)
		lpd.Teams[0] = &model.Team{ID: 0, Name: "(all)"}
	}
	r.Tx.Commit()
	util.RenderPage(r, &util.Page{
		Title:    "People",
		MenuItem: "people",
		BodyData: lpd,
	}, template.Must(template.New("listPeople").Parse(listPeopleTemplate)))
	return nil
}

type listPeopleData struct {
	People []*model.Person
	Teams  []*model.Team
	TeamID model.TeamID
}

const listPeopleTemplate = `{{ define "body" -}}
<div id="listPeople">
  <form id="listPeople-title" method="GET" class="pageTitle">
    People
    {{- if .Teams }}
      <select id="listPeople-team" name="team">
        {{- $tid := .TeamID }}
	{{- range .Teams }}
	  <option value="{{ .ID }}"{{ if eq .ID $tid }} selected{{ end }}>{{ .Name }}</option>
	{{ end }}
      </select>
    {{ end }}
  </form>
  <table id="listPeople-table">
    <thead>
      <tr>
        <th>Person</th>
	<th>Contact Info</th>
	<th>Roles</th>
      </tr>
    </thead>
    <tbody>
      {{- range .People }}
        <tr>
	  <td><a href="/people/{{ .ID }}">{{ .LastName }}, {{ .FirstName }}</a></td>
	  <td>
	    <div><a href="mailto:{{ .Email }}">{{ .Email }}</div>
	    {{- if .Phone }}
	      <div><a href="tel:{{ .Phone }}">{{ .Phone }}</div>
	    {{ end }}
	  </td>
	  <td>
	    {{- range .Roles }}
	      <div>{{ .Team.Name }}{{ if .Name }}: {{ .Name }}{{ end }}</div>
	    {{ else }}
	      <div>&mdash;</div>
	    {{ end }}
	  </td>
	</tr>
      {{ end }}
    </tbody>
  </table>
  <div id="listPeople-buttons">
    <a class="btn btn-secondary" href="/people/NEW">Add Person</a>
  </div>
</div>
{{- end }}`
