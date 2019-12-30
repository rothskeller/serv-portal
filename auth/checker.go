package auth

import (
	"sort"

	"serv.rothskeller.net/portal/model"
)

// Checker encapsulates all of the logic of checking authorization.
type Checker struct {
	teams []*model.Team
	tmap  map[model.TeamID]*model.Team
}

// NewChecker creates a new auth checker.
func NewChecker() *Checker {
	return &Checker{tmap: make(map[model.TeamID]*model.Team)}
}

// AddTeams adds a set of teams to the auth checker.  They have no privileges
// initially.
func (c *Checker) AddTeams(teams []*model.Team) {
	c.teams = append(c.teams, teams...)
	for _, t := range teams {
		c.tmap[t.ID] = t
	}
	sort.Sort(teamSorter(c.teams))
}

type teamSorter []*model.Team

func (ts teamSorter) Len() int           { return len(ts) }
func (ts teamSorter) Less(i, j int) bool { return ts[i].Name < ts[j].Name }
func (ts teamSorter) Swap(i, j int)      { ts[i], ts[j] = ts[j], ts[i] }
