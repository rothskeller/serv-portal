package person

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// EditPerson handles GET and POST /people/$id requests (where $id may be "NEW").
func EditPerson(r *util.Request, idstr string) error {
	var (
		epd   editPersonData
		title string
	)
	epd.Teams, epd.AdminTeams, epd.ManageTeams = editPersonTeams(r)
	if idstr == "NEW" {
		if len(epd.ManageTeams) == 0 {
			return util.Forbidden
		}
		epd.Person = new(model.Person)
		title = "New Person"
		epd.CanEditInfo = true
	} else {
		if epd.Person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); epd.Person == nil {
			return util.NotFound
		}
		if !r.Person.CanViewPerson(epd.Person) {
			return util.Forbidden
		}
		title = epd.Person.FirstName + " " + epd.Person.LastName
		epd.CanEditInfo = r.Person == epd.Person || r.Person.IsWebmaster()
	}
	if r.Method == http.MethodPost {
		if len(epd.AdminTeams) == 0 && !epd.CanEditInfo {
			return util.Forbidden
		}
		if epd.CanEditInfo {
			if epd.Person.FirstName = strings.TrimSpace(r.FormValue("firstName")); epd.Person.FirstName == "" {
				epd.FirstNameError = "The person's first name is required."
			}
			if epd.Person.LastName = strings.TrimSpace(r.FormValue("lastName")); epd.Person.LastName == "" {
				epd.LastNameError = "The person's last name is required."
			}
			if epd.Person.Email = strings.TrimSpace(r.FormValue("email")); epd.Person.Email == "" {
				epd.EmailError = "The person's email address is required."
			} else if !emailRE.MatchString(epd.Person.Email) {
				epd.EmailError = "This is not a valid email address."
			} else if emailInUse(r, epd.Person) {
				epd.EmailError = "This email address is in use by a different person."
			}
			if epd.Person.Phone = strings.TrimSpace(r.FormValue("phone")); epd.Person.Phone != "" {
				ph := strings.Map(keepDigits, epd.Person.Phone)
				if len(ph) != 10 {
					epd.PhoneError = "A valid phone number has exactly 10 digits."
				} else {
					epd.Person.Phone = ph[0:3] + "-" + ph[3:6] + "-" + ph[6:10]
				}
			}
		}
		var rmap = make(map[*model.Team]*model.Role, len(epd.Person.Roles))
		for _, r := range epd.Person.Roles {
			rmap[r.Team] = r
		}
		var tmap = make(map[string]bool, len(r.Form["team"]))
		for _, tidstr := range r.Form["team"] {
			tmap[tidstr] = true
		}
		for _, t := range epd.Teams {
			if !epd.AdminTeams[t] {
				continue
			}
			tidstr := strconv.Itoa(int(t.ID))
			if !tmap[tidstr] && epd.ManageTeams[t] {
				delete(rmap, t)
			} else {
				if role := r.Tx.FetchRole(model.RoleID(util.ParseID(r.FormValue("role-" + tidstr)))); role != nil && role.Team == t {
					rmap[t] = role
				} else {
					return errors.New("invalid role")
				}
			}
		}
		epd.Person.Roles = epd.Person.Roles[:0]
		for _, t := range r.Tx.FetchTeams() {
			if r := rmap[t]; r != nil {
				epd.Person.Roles = append(epd.Person.Roles, r)
			}
		}
		if epd.FirstNameError != "" || epd.LastNameError != "" || epd.EmailError != "" || epd.PhoneError != "" {
			goto SHOWFORM
		}
		r.Tx.SavePerson(epd.Person)
		r.Tx.Commit()
		http.Redirect(r, r.Request, "/people", http.StatusSeeOther)
		return nil
	}
SHOWFORM:
	r.Tx.Commit()
	epd.TeamMap = make(map[*model.Team]bool)
	epd.RoleMap = make(map[*model.Role]bool)
	for _, role := range epd.Person.Roles {
		epd.RoleMap[role] = true
		epd.TeamMap[role.Team] = true
	}
	util.RenderPage(r, &util.Page{
		Title:    title,
		MenuItem: "people",
		BodyData: &epd,
	}, template.Must(template.New("editPerson").Parse(editPersonTemplate)))
	return nil
}

func editPersonTeams(r *util.Request) (teams []*model.Team, admin, manage map[*model.Team]bool) {
	admin = make(map[*model.Team]bool)
	manage = make(map[*model.Team]bool)
	for _, t := range r.Tx.FetchTeams() {
		if t.Type != model.TeamNormal {
			continue
		}
		teams = append(teams, t)
		if r.Person.PrivMap.Has(t, model.PrivManage) {
			manage[t] = true
		}
		if r.Person.PrivMap.Has(t, model.PrivAdmin) {
			admin[t] = true
		}
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams, admin, manage
}

func emailInUse(r *util.Request, person *model.Person) bool {
	for _, p := range r.Tx.FetchPeople() {
		if p.ID == person.ID {
			continue
		}
		if strings.EqualFold(p.Email, person.Email) {
			return true
		}
	}
	return false
}

func keepDigits(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}

type editPersonData struct {
	FirstNameError string
	LastNameError  string
	EmailError     string
	PhoneError     string
	Person         *model.Person
	CanEditInfo    bool
	RoleMap        map[*model.Role]bool
	TeamMap        map[*model.Team]bool
	Teams          []*model.Team
	AdminTeams     map[*model.Team]bool
	ManageTeams    map[*model.Team]bool
}

const editPersonTemplate = `{{ define "body" -}}
<div id="editPerson">
  <div class="pageTitle">{{ if .Person.ID }}Edit Person{{ else }}Create Person{{ end }}</div>
  <form method="POST">
    <div id="editPerson-firstName-row">
      <label for="editPerson-firstName" id="editPerson-firstName-label">First name</label>
      {{- if .CanEditInfo }}
        <input id="editPerson-firstName" name="firstName" class="form-control" autofocus value="{{ .Person.FirstName }}">
      {{ else }}
        <div>{{ .Person.FirstName }}</div>
      {{ end }}
    </div>
    {{- if .FirstNameError }}<div id="editPerson-firstName-error">{{ .FirstNameError }}</div>{{ end }}
    <div id="editPerson-lastName-row">
      <label for="editPerson-lastName" id="editPerson-lastName-label">Last name</label>
      {{- if .CanEditInfo }}
        <input id="editPerson-lastName" name="lastName" class="form-control" autofocus value="{{ .Person.LastName }}">
      {{ else }}
        <div>{{ .Person.LastName }}</div>
      {{ end }}
    </div>
    {{- if .LastNameError }}<div id="editPerson-lastName-error">{{ .LastNameError }}</div>{{ end }}
    <div id="editPerson-email-row">
      <label for="editPerson-email" id="editPerson-email-label">Email address</label>
      {{- if .CanEditInfo }}
        <input id="editPerson-email" name="email" class="form-control" autofocus value="{{ .Person.Email }}">
      {{ else }}
        <div><a href="mailto:{{ .Person.Email }}">{{ .Person.Email }}</a></div>
      {{ end }}
    </div>
    {{- if .EmailError }}<div id="editPerson-email-error">{{ .EmailError }}</div>{{ end }}
    <div id="editPerson-phone-row">
      <label for="editPerson-phone" id="editPerson-phone-label">Phone number</label>
      {{- if .CanEditInfo }}
        <input id="editPerson-phone" name="phone" class="form-control" autofocus value="{{ .Person.Phone }}">
      {{ else }}
        <div><a href="tel:{{ .Person.Phone }}">{{ .Person.Phone }}</a></div>
      {{ end }}
    </div>
    {{- if .PhoneError }}<div id="editPerson-phone-error">{{ .PhoneError }}</div>{{ end }}
    {{- $epd := . }}
    <div class="editPerson-group-label">This person belongs to these teams:</div>
    <table id="editPerson-team-table">
    {{- range .Teams }}
      <tr>
        <td>
	  {{ if $epd.ManageTeams }}
            <input id="editPerson-team-{{ .ID }}" name="team" type="checkbox" value="{{ .ID }}"{{ if index $epd.TeamMap . }} checked{{ end }}{{ if not (index $epd.ManageTeams .) }} disabled{{ end }}>
	  {{ end }}
	  <label for="editPerson-team-{{ .ID }}">{{ .Name }}</label>
	</td>
	<td>
	  {{- if and (index $epd.AdminTeams .) (gt (len .Roles) 1) }}
	    <select id="editPerson-role-{{ .ID }}" name="role-{{ .ID }}"{{ if not (index $epd.TeamMap .) }} style="display:none"{{ end }}>
	      {{- range .Roles }}
	        <option value="{{ .ID }}"{{ if index $epd.RoleMap . }} selected{{ end }}>{{ if .Name }}{{ .Name }}{{ else }}(member){{ end }}</option>
	      {{ end }}
	    </select>
	  {{ else if gt (len .Roles) 1 }}
	    {{ .Name }}
	  {{ else if index $epd.AdminTeams . }}
	    <input type="hidden" name="role-{{ .ID }}" value="{{ (index .Roles 0).ID }}">
	  {{ end }}
	</td>
      </tr>
    {{ end }}
    </table>
    {{ if $epd.AdminTeams }}
      <div id="editPerson-submit-row">
        <button type="submit" class="btn btn-primary">{{ if .Person.ID }}Save Person{{ else }}Create Person{{ end }}</button>
        <a class="btn btn-secondary" href="/people">Cancel</a>
      </div>
    {{ end }}
  </form>
</div>
{{- end }}`
