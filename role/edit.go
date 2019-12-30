package role

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// EditRole handles GET and POST /teams/$tid/roles/$rid requests (where $id may be "NEW").
func EditRole(r *util.Request, tidstr, ridstr string) error {
	var (
		team  *model.Team
		erd   editRoleData
		title string
	)
	if team = r.Tx.FetchTeam(model.TeamID(util.ParseID(tidstr))); team == nil {
		return util.NotFound
	}
	if ridstr == "NEW" {
		erd.Role = &model.Role{
			Team:    team,
			PrivMap: make(model.PrivilegeMap),
		}
		erd.Role.PrivMap.Merge(team.PrivMap)
		title = "New Role"
		erd.NoDelete = true
	} else {
		if erd.Role = r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); erd.Role == nil {
			return util.NotFound
		}
		if erd.Role.Team != team {
			return errors.New("role does not belong to team")
		}
		title = "Role: " + erd.Role.Name
		erd.NoDelete = len(team.Roles) == 1
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	erd.AllTeams = r.Tx.FetchTeams()
	if r.Method == http.MethodPost {
		if r.FormValue("delete") != "" && erd.Role.ID != 0 && len(erd.Role.Team.Roles) > 1 {
			r.Tx.DeleteRole(erd.Role)
			http.Redirect(r, r.Request, "/teams", http.StatusSeeOther)
			r.Tx.Commit()
			return nil
		}
		erd.Role.Name = strings.TrimSpace(r.FormValue("name"))
		for _, t := range team.Roles {
			if t != erd.Role && t.Name == erd.Role.Name {
				erd.NameError = "Another role already has this name."
			}
		}
		erd.Role.PrivMap = make(model.PrivilegeMap)
		erd.Role.PrivMap.Set(team, model.PrivMember)
		for _, t := range erd.AllTeams {
			tid := strconv.Itoa(int(t.ID))
			if r.FormValue("member-"+tid) != "" {
				erd.Role.PrivMap.Add(t, model.PrivMember)
			}
			if r.FormValue("manager-"+tid) != "" {
				erd.Role.PrivMap.Add(t, model.PrivView|model.PrivAdmin|model.PrivManage)
			} else if r.FormValue("admin-"+tid) != "" {
				erd.Role.PrivMap.Add(t, model.PrivView|model.PrivAdmin)
			} else if r.FormValue("viewer-"+tid) != "" {
				erd.Role.PrivMap.Add(t, model.PrivView)
			}
		}
		if erd.NameError != "" {
			goto SHOWFORM
		}
		r.Tx.SaveRole(erd.Role)
		r.Tx.Commit()
		http.Redirect(r, r.Request, "/teams", http.StatusSeeOther)
		return nil
	}
SHOWFORM:
	erd.AllTeams = util.SortTeamsHierarchically(erd.AllTeams)
	r.Tx.Commit()
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "roles",
		BodyData: &erd,
	}, template.Must(template.New("editRole").Parse(editRoleTemplate)))
	return nil
}

type editRoleData struct {
	NameError string
	NoDelete  bool
	Role      *model.Role
	AllTeams  []*model.Team
}

const editRoleTemplate = `{{ define "body" -}}
<div id="editRole">
  <div class="pageTitle">{{ if .Role.ID }}Edit Role{{ else }}Create Role{{ end }}</div>
  <form method="POST">
    <input type="hidden" name="team" value="{{ .Role.Team.ID }}">
    <div id="editRole-name-row">
      <label for="editRole-name" id="editRole-name-label">Role name</label>
      <input id="editRole-name" name="name" class="form-control" autofocus value="{{ .Role.Name }}">
    </div>
    {{- if .NameError }}<div id="editRole-name-error">{{ .NameError }}</div>{{ end }}
    <div id="editRole-team-row">
      <div id="editRole-team-label">Team</div>
      <div>{{ .Role.Team.Name }}</div>
    </div>
    <div id="editRole-privileges-label">This role has the following privileges on teams:</div>
    <table id="editRole-privileges-table">
      <tr>
        <th></th>
	<th>Member</th>
	<th>Viewer</th>
	<th>Admin</th>
	<th>Manager</th>
      </tr>
      {{ $pm := .Role.PrivMap }}
      {{- range .AllTeams }}
        <tr>
	  <td class="indent-{{ .Depth }}">{{ .Name }}</td>
	  <td><input name="member-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 1 3 5 7 9 11 13 15 }} checked{{ end }}></td>
	  <td><input name="viewer-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 2 3 6 7 10 11 14 15 }} checked{{ end }}></td>
	  <td><input name="admin-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 4 5 6 7 12 13 14 15 }} checked{{ end }}></td>
	  <td><input name="manager-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 8 9 10 11 12 13 14 15 }} checked{{ end }}></td>
	</tr>
      {{ end }}
    </table>
    <div id="editRole-submit-row">
      <button type="submit" class="btn btn-primary">{{ if .Role.ID }}Save Role{{ else }}Create Role{{ end }}</button>
      <a class="btn btn-secondary" href="/teams">Cancel</a>
      {{- if not .NoDelete }}
        <button id="editRole-delete" name="delete" type="submit" class="btn btn-danger">Delete Role</button>
      {{ end }}
    </div>
  </form>
</div>
{{- end }}`
