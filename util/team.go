package util

import (
	"sort"

	"serv.rothskeller.net/portal/model"
)

// SortTeamsHierarchically takes a list of teams and returns a new slice with
// the same teams in the appropriate order for display in a tree list.  Teams
// under the same parent are sorted by name.  The argument slice not changed.
func SortTeamsHierarchically(teams []*model.Team) (nt []*model.Team) {
	nt = make([]*model.Team, 0, len(teams))
	for _, t := range teams {
		if t.Parent == nil {
			nt = appendChildHierarchically(nt, t)
		}
	}
	return nt
}
func appendChildHierarchically(teams []*model.Team, team *model.Team) []*model.Team {
	teams = append(teams, team)
	sort.Slice(team.Children, func(i, j int) bool {
		return team.Children[i].Name < team.Children[j].Name
	})
	for _, c := range team.Children {
		teams = appendChildHierarchically(teams, c)
	}
	return teams
}
