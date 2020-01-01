package team

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetTeam handles GET /api/teams/$id requests (where $id may be "NEW").
func GetTeam(r *util.Request, idstr string) error {
	var (
		noDelete bool
		team     *model.Team
		allTeams []*model.Team
		out      jwriter.Writer
	)
	if idstr == "NEW" {
		team = &model.Team{PrivMap: make(model.PrivilegeMap)}
		team.PrivMap.Set(team, model.PrivMember)
		team.Parent = r.Tx.FetchTeam(model.TeamID(util.ParseID(r.FormValue("parent"))))
		if team.Parent == nil || team.Parent.Type != model.TeamAncestor {
			return errors.New("invalid parent team")
		}
		team.Parent.Children = append(team.Parent.Children, team)
		team.PrivMap.Merge(team.Parent.PrivMap)
		for _, t := range r.Tx.FetchTeams() {
			t.PrivMap.Set(team, t.PrivMap.Get(team.Parent)&^model.PrivMember)
		}
		noDelete = true
	} else {
		if team = r.Tx.FetchTeam(model.TeamID(util.ParseID(idstr))); team == nil {
			return util.NotFound
		}
		noDelete = len(team.Children) != 0
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	allTeams = r.Tx.FetchTeams()
	r.Tx.Commit()
	if team.ID == 0 {
		allTeams = append(allTeams, team)
	}
	allTeams = util.SortTeamsHierarchically(allTeams)
	out.RawString(`{"id":`)
	out.Int(int(team.ID))
	out.RawString(`,"name":`)
	out.String(team.Name)
	out.RawString(`,"email":`)
	out.String(team.Email)
	out.RawString(`,"canDelete":`)
	out.Bool(!noDelete)
	if team.Parent == nil {
		out.RawString(`,"parent":null`)
	} else {
		out.RawString(`,"parent":`)
		out.String(team.Parent.Name)
	}
	out.RawString(`,"privs":[`)
	for i, t := range allTeams {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"indent":`)
		out.Int(t.Depth())
		out.RawString(`,"id":`)
		out.Int(int(t.ID))
		out.RawString(`,"name":`)
		out.String(t.Name)
		out.RawString(`,"actor":{"member":`)
		out.Bool(team.PrivMap.Has(t, model.PrivMember))
		out.RawString(`,"access":`)
		switch {
		case team.PrivMap.Has(t, model.PrivManage):
			out.String("manage")
		case team.PrivMap.Has(t, model.PrivAdmin):
			out.String("admin")
		case team.PrivMap.Has(t, model.PrivView):
			out.String("view")
		default:
			out.String("none")
		}
		out.RawString(`},"target":{"member":`)
		out.Bool(t.PrivMap.Has(team, model.PrivMember))
		out.RawString(`,"access":`)
		switch {
		case t.PrivMap.Has(team, model.PrivManage):
			out.String("manage")
		case t.PrivMap.Has(team, model.PrivAdmin):
			out.String("admin")
		case t.PrivMap.Has(team, model.PrivView):
			out.String("view")
		default:
			out.String("none")
		}
		out.RawString(`}}`)
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostTeam handles POST /api/teams/$id requests (where $id may be "NEW").
func PostTeam(r *util.Request, idstr string) error {
	var team *model.Team

	if idstr == "NEW" {
		team = &model.Team{PrivMap: make(model.PrivilegeMap)}
		if team.Parent = r.Tx.FetchTeam(model.TeamID(util.ParseID(r.FormValue("parent")))); team.Parent == nil || team.Parent.Type != model.TeamAncestor {
			return errors.New("invalid parent team")
		}
	} else {
		if team = r.Tx.FetchTeam(model.TeamID(util.ParseID(idstr))); team == nil {
			return util.NotFound
		}
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && team.ID != 0 {
		r.Tx.DeleteTeam(team)
		r.Tx.Commit()
		return nil
	}
	team.Name = strings.TrimSpace(r.FormValue("name"))
	if team.Name == "" {
		return errors.New("missing name")
	}
	for _, t := range r.Tx.FetchTeams() {
		if t != team && t.Name == team.Name {
			r.Header().Set("Content-Type", "application/json")
			r.Write([]byte(`{"nameError":"Another team already has this name."}`))
			return nil
		}
	}
	team.Email = strings.TrimSpace(r.FormValue("email"))
	team.PrivMap = make(model.PrivilegeMap)
	team.PrivMap.Set(team, model.PrivMember)
	for _, t := range r.Tx.FetchTeams() {
		tid := strconv.Itoa(int(t.ID))
		if r.FormValue("a:member-"+tid) != "" {
			team.PrivMap.Add(t, model.PrivMember)
		}
		switch r.FormValue("a:access-" + tid) {
		case "manage":
			team.PrivMap.Add(t, model.PrivView|model.PrivAdmin|model.PrivManage)
		case "admin":
			team.PrivMap.Add(t, model.PrivView|model.PrivAdmin)
		case "view":
			team.PrivMap.Add(t, model.PrivView)
		}
		if t != team {
			if r.FormValue("t:member-"+tid) != "" {
				t.PrivMap.Set(team, model.PrivMember)
			} else {
				t.PrivMap.Set(team, 0)
			}
			switch r.FormValue("t:access-" + tid) {
			case "manage":
				t.PrivMap.Add(team, model.PrivView|model.PrivAdmin|model.PrivManage)
			case "admin":
				t.PrivMap.Add(team, model.PrivView|model.PrivAdmin)
			case "view":
				t.PrivMap.Add(team, model.PrivView)
			}
		}
	}
	r.Tx.FetchTeamByTag(model.TeamWebmasters).PrivMap.Add(team, model.PrivView|model.PrivAdmin|model.PrivManage)
	r.Tx.SaveTeam(team)
	r.Tx.Commit()
	return nil
}
