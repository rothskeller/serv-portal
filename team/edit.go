package team

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// EditTeam handles GET and POST /teams/$id requests (where $id may be "NEW").
func EditTeam(r *util.Request, idstr string) error {
	var (
		etd   editTeamData
		title string
	)
	if idstr == "NEW" {
		etd.Team = &model.Team{PrivMap: make(model.PrivilegeMap)}
		etd.Team.PrivMap.Set(etd.Team, model.PrivMember)
		if etd.Team.Parent = r.Tx.FetchTeam(model.TeamID(util.ParseID(r.FormValue("parent")))); etd.Team.Parent == nil {
			return errors.New("bad parent team")
		}
		etd.Team.Parent.Children = append(etd.Team.Parent.Children, etd.Team)
		etd.Parent = etd.Team.Parent.ID
		etd.Team.PrivMap.Merge(etd.Team.Parent.PrivMap)
		for _, t := range r.Tx.FetchTeams() {
			t.PrivMap.Set(etd.Team, t.PrivMap.Get(etd.Team.Parent)&^model.PrivMember)
		}
		title = "New Team"
		etd.NoDelete = true
	} else {
		if etd.Team = r.Tx.FetchTeam(model.TeamID(util.ParseID(idstr))); etd.Team == nil {
			return util.NotFound
		}
		title = "Team: " + etd.Team.Name
		etd.NoDelete = len(etd.Team.Children) != 0
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	etd.AllTeams = r.Tx.FetchTeams()
	if etd.Team.ID == 0 {
		etd.AllTeams = append(etd.AllTeams, etd.Team)
	}
	if r.Method == http.MethodPost {
		if r.FormValue("delete") != "" && etd.Team.ID != 0 {
			r.Tx.DeleteTeam(etd.Team)
			http.Redirect(r, r.Request, "/teams", http.StatusSeeOther)
			r.Tx.Commit()
			return nil
		}
		etd.Team.Name = strings.TrimSpace(r.FormValue("name"))
		if etd.Team.Name == "" {
			etd.NameError = "The team name is required."
		} else {
			for _, t := range r.Tx.FetchTeams() {
				if t != etd.Team && t.Name == etd.Team.Name {
					etd.NameError = "Another team already has this name."
				}
			}
		}
		etd.Team.Email = strings.TrimSpace(r.FormValue("email"))
		etd.Team.PrivMap = make(model.PrivilegeMap)
		etd.Team.PrivMap.Set(etd.Team, model.PrivMember)
		for _, t := range etd.AllTeams {
			tid := strconv.Itoa(int(t.ID))
			if r.FormValue("member-"+tid) != "" {
				etd.Team.PrivMap.Add(t, model.PrivMember)
			}
			if r.FormValue("manager-"+tid) != "" {
				etd.Team.PrivMap.Add(t, model.PrivView|model.PrivAdmin|model.PrivManage)
			} else if r.FormValue("admin-"+tid) != "" {
				etd.Team.PrivMap.Add(t, model.PrivView|model.PrivAdmin)
			} else if r.FormValue("viewer-"+tid) != "" {
				etd.Team.PrivMap.Add(t, model.PrivView)
			}
			if t != etd.Team {
				tid := strconv.Itoa(int(t.ID))
				if r.FormValue("inmember-"+tid) != "" {
					t.PrivMap.Set(etd.Team, model.PrivMember)
				} else {
					t.PrivMap.Set(etd.Team, 0)
				}
				if r.FormValue("inmanager-"+tid) != "" {
					t.PrivMap.Add(etd.Team, model.PrivView|model.PrivAdmin|model.PrivManage)
				} else if r.FormValue("inadmin-"+tid) != "" {
					t.PrivMap.Add(etd.Team, model.PrivView|model.PrivAdmin)
				} else if r.FormValue("inviewer-"+tid) != "" {
					t.PrivMap.Add(etd.Team, model.PrivView)
				}
			}
		}
		r.Tx.FetchTeamByTag(model.TeamWebmasters).PrivMap.Add(etd.Team, model.PrivView|model.PrivAdmin|model.PrivManage)
		if etd.NameError != "" {
			goto SHOWFORM
		}
		r.Tx.SaveTeam(etd.Team)
		r.Tx.Commit()
		http.Redirect(r, r.Request, "/teams", http.StatusSeeOther)
		return nil
	}
SHOWFORM:
	etd.AllTeams = util.SortTeamsHierarchically(etd.AllTeams)
	r.Tx.Commit()
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "teams",
		BodyData: &etd,
	}, template.Must(template.New("editTeam").Parse(editTeamTemplate)))
	return nil
}

type editTeamData struct {
	NameError string
	NoDelete  bool
	Parent    model.TeamID
	Team      *model.Team
	AllTeams  []*model.Team
}

const editTeamTemplate = `{{ define "body" -}}
<div id="editTeam">
  <div class="pageTitle">{{ if .Team.ID }}Edit Team{{ else }}Create Team{{ end }}</div>
  <form method="POST">
    {{- if .Parent }}
      <input type="hidden" name="parent" value="{{ .Parent }}">
    {{ end }}
    <div id="editTeam-name-row">
      <label for="editTeam-name" id="editTeam-name-label">Team name</label>
      <input id="editTeam-name" name="name" class="form-control" autofocus value="{{ .Team.Name }}">
    </div>
    {{- if .NameError }}<div id="editTeam-name-error">{{ .NameError }}</div>{{ end }}
    <div id="editTeam-parent-row">
      <div id="editTeam-parent-label">Parent team</div>
      <div>{{ with .Team.Parent }}{{ .Name }}{{ else }}(none){{ end }}</div>
    </div>
    <div id="editTeam-email-row">
      <label for="editTeam-email" id="editTeam-email-label">Team email</label>
      <input id="editTeam-email" name="email" type="email" class="form-control" value="{{ .Team.Email }}">
    </div>
    <div id="editTeam-privileges-label">This team has the following privileges on other teams:</div>
    <table id="editTeam-privileges-table">
      <tr>
        <th></th>
	<th>Member</th>
	<th>Viewer</th>
	<th>Admin</th>
	<th>Manager</th>
      </tr>
      {{ $pm := .Team.PrivMap }}
      {{- range .AllTeams }}
        <tr>
	  <td class="indent-{{ .Depth }}">{{ if .ID }}{{ .Name }}{{ else }}(new team){{ end }}</td>
	  <td><input name="member-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 1 3 5 7 9 11 13 15 }} checked{{ end }}></td>
	  <td><input name="viewer-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 2 3 6 7 10 11 14 15 }} checked{{ end }}></td>
	  <td><input name="admin-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 4 5 6 7 12 13 14 15 }} checked{{ end }}></td>
	  <td><input name="manager-{{ .ID }}" type="checkbox"{{ if eq (index $pm .) 8 9 10 11 12 13 14 15 }} checked{{ end }}></td>
	</tr>
      {{ end }}
    </table>
    <div id="editTeam-inPrivs-label">These other teams have the following privileges on this team:</div>
    <table id="editTeam-inPrivs-table">
      <tr>
        <th></th>
	<th>Member</th>
	<th>Viewer</th>
	<th>Admin</th>
	<th>Manager</th>
      </tr>
      {{ $team := .Team }}
      {{- range .AllTeams }}{{ if ne .ID $team.ID }}
        <tr>
	  <td class="editTeam-indent-{{ .Depth }}">{{ .Name }}</td>
	  <td><input name="inmember-{{ .ID }}" type="checkbox"{{ if eq (index .PrivMap $team) 1 3 5 7 9 11 13 15 }} checked{{ end }}></td>
	  <td><input name="inviewer-{{ .ID }}" type="checkbox"{{ if eq (index .PrivMap $team) 2 3 6 7 10 11 14 15 }} checked{{ end }}></td>
	  <td><input name="inadmin-{{ .ID }}" type="checkbox"{{ if eq (index .PrivMap $team) 4 5 6 7 12 13 14 15 }} checked{{ end }}></td>
	  <td><input name="inmanager-{{ .ID }}" type="checkbox"{{ if eq (index .PrivMap $team) 8 9 10 11 12 13 14 15 }} checked{{ end }}></td>
	</tr>
      {{ end }}{{ end }}
    </table>
    <div id="editTeam-submit-row">
      <button type="submit" class="btn btn-primary">{{ if .Team.ID }}Save Team{{ else }}Create Team{{ end }}</button>
      <a class="btn btn-secondary" href="/teams">Cancel</a>
      {{- if not .NoDelete }}
        <button id="editTeam-delete" name="delete" type="submit" class="btn btn-danger">Delete Team</button>
      {{ end }}
    </div>
  </form>
</div>
{{- end }}`
