package team

import (
	"html/template"

	"serv.rothskeller.net/portal/util"
)

// ListTeams handles GET /teams requests.
func ListTeams(r *util.Request) error {
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	teams := r.Tx.FetchTeams()
	teams = util.SortTeamsHierarchically(teams)
	r.Tx.Commit()
	util.RenderPage(r, &util.Page{
		Title:    "Teams and Roles",
		MenuItem: "teams",
		BodyData: teams,
	}, template.Must(template.New("listTeams").Parse(listTeamsTemplate)))
	return nil
}

const listTeamsTemplate = `{{ define "body" -}}
<div id="listTeams">
  <div class="pageTitle">Teams and Roles</div>
  <table id="listTeams-table">
    <thead>
      <tr>
        <th>Team</th>
	<th>Email</th>
	<th>Roles</th>
	<th></th>
      </tr>
    </thead>
    <tbody>
      {{- range . }}
        <tr>
	  <td class="indent-{{ .Depth }}"><a href="/teams/{{ .ID }}">{{ .Name }}</a></td>
	  <td>{{ if .Email }}{{ .Email }}{{ else }}&mdash;{{ end }}</td>
	  <td>
	    {{- range .Roles }}
	      <div><a href="/teams/{{ .Team.ID }}/roles/{{ .ID }}">{{ with .Name }}{{.}}{{ else }}(member){{ end }}</a></div>
	    {{ end }}
	  </td>
	  <td>
	    {{- if eq .Type 0 }}
	      <a href="/teams/{{ .ID }}/roles/NEW">Add Role</a>
	    {{ else if eq .Type 2 }}
	      <a href="/teams/NEW?parent={{ .ID }}">Add Child Team</a>
	    {{ end }}
	  </td>
	</tr>
      {{ end }}
    </tbody>
  </table>
</div>
{{- end }}`
