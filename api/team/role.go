package team

import (
	"errors"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetRole handles GET /api/teams/$tid/roles/$rid requests (where $rid may be "NEW").
func GetRole(r *util.Request, tidstr, ridstr string) error {
	var (
		team     *model.Team
		allTeams []*model.Team
		noDelete bool
		role     *model.Role
		out      jwriter.Writer
	)
	if team = r.Tx.FetchTeam(model.TeamID(util.ParseID(tidstr))); team == nil {
		return util.NotFound
	}
	if ridstr == "NEW" {
		role = &model.Role{
			Team:    team,
			PrivMap: make(model.PrivilegeMap),
		}
		role.PrivMap.Merge(team.PrivMap)
		noDelete = true
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); role == nil {
			return util.NotFound
		}
		if role.Team != team {
			return errors.New("role does not belong to team")
		}
		noDelete = len(team.Roles) == 1
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	allTeams = util.SortTeamsHierarchically(r.Tx.FetchTeams())
	r.Tx.Commit()
	out.RawString(`{"team":{"id":`)
	out.Int(int(team.ID))
	out.RawString(`,"name":`)
	out.String(team.Name)
	out.RawString(`},"role":{"id":`)
	out.Int(int(role.ID))
	out.RawString(`,"name":`)
	out.String(role.Name)
	out.RawString(`},"canDelete":`)
	out.Bool(!noDelete)
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
		out.RawString(`,"member":`)
		out.Bool(role.PrivMap.Has(t, model.PrivMember))
		out.RawString(`,"access":`)
		switch {
		case role.PrivMap.Has(t, model.PrivManage):
			out.String("manage")
		case role.PrivMap.Has(t, model.PrivAdmin):
			out.String("admin")
		case role.PrivMap.Has(t, model.PrivView):
			out.String("view")
		default:
			out.String("none")
		}
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRole handles POST /api/teams/$tid/roles/$rid requests (where $rid may be "NEW").
func PostRole(r *util.Request, tidstr, ridstr string) error {
	var (
		team *model.Team
		role *model.Role
	)
	if team = r.Tx.FetchTeam(model.TeamID(util.ParseID(tidstr))); team == nil {
		return util.NotFound
	}
	if ridstr == "NEW" {
		role = &model.Role{
			Team:    team,
			PrivMap: make(model.PrivilegeMap),
		}
		role.PrivMap.Merge(team.PrivMap)
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(ridstr))); role == nil {
			return util.NotFound
		}
		if role.Team != team {
			return errors.New("role does not belong to team")
		}
	}
	if !r.Person.IsWebmaster() {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && role.ID != 0 && len(role.Team.Roles) > 1 {
		r.Tx.DeleteRole(role)
		r.Tx.Commit()
		return nil
	}
	role.Name = strings.TrimSpace(r.FormValue("name"))
	for _, t := range team.Roles {
		if t != role && t.Name == role.Name {
			r.Header().Set("Content-Type", "application/json")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		}
	}
	role.PrivMap = make(model.PrivilegeMap)
	role.PrivMap.Set(team, model.PrivMember)
	for _, t := range r.Tx.FetchTeams() {
		tid := strconv.Itoa(int(t.ID))
		if r.FormValue("member-"+tid) == "true" {
			role.PrivMap.Add(t, model.PrivMember)
		}
		switch r.FormValue("access-" + tid) {
		case "manage":
			role.PrivMap.Add(t, model.PrivView|model.PrivAdmin|model.PrivManage)
		case "admin":
			role.PrivMap.Add(t, model.PrivView|model.PrivAdmin)
		case "view":
			role.PrivMap.Add(t, model.PrivView)
		}
	}
	r.Tx.SaveRole(role)
	r.Tx.Commit()
	return nil
}
