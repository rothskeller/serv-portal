package model

import (
	"sort"
)

// Add adds a privilege to a map.
func (pm PrivilegeMap) Add(team *Team, priv Privilege) {
	pm[team] |= priv
}

// Set sets the privilege on a team in a map.
func (pm PrivilegeMap) Set(team *Team, priv Privilege) {
	pm[team] = priv
}

// Get returns the privileges on a team in a map.
func (pm PrivilegeMap) Get(team *Team) Privilege {
	return pm[team]
}

// Has returns whether the specified privilege(s) exist in the receiver map for
// the specified team.
func (pm PrivilegeMap) Has(team *Team, priv Privilege) bool {
	return pm[team]&priv == priv
}

// HasAny returns whether the receiver map has the specified privilege(s) on any
// team.
func (pm PrivilegeMap) HasAny(priv Privilege) bool {
	for _, p := range pm {
		if p&priv == priv {
			return true
		}
	}
	return false
}

// HasTag returns whether the specified privilege(s) exist in the receiver map
// for the team with the specified tag.
func (pm PrivilegeMap) HasTag(tag TeamTag, priv Privilege) bool {
	for t, p := range pm {
		if t.Tag == tag {
			return p&priv == priv
		}
	}
	return false
}

// Merge merges all of the privileges in the argument map into the receiver map.
func (pm PrivilegeMap) Merge(other PrivilegeMap) {
	for t, p := range other {
		pm[t] |= p
	}
}

// TeamsWith returns a sorted list of all teams for which the map contains the
// specified privilege(s).
func (pm PrivilegeMap) TeamsWith(privs Privilege) (teams []*Team) {
	for t, p := range pm {
		if p&privs == privs {
			teams = append(teams, t)
		}
	}
	sort.Slice(teams, func(i, j int) bool {
		return teams[i].Name < teams[j].Name
	})
	return teams
}
