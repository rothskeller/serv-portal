package authz

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

// UpdateAuthz recalculates all of the interactions between people, roles, and
// lists, and enforces consistency between them.  It is an expensive call,
// because it updates every person, role, and list in the database.
func UpdateAuthz(tx *store.Tx) {
	var manualSubs = map[model.ListID]map[model.PersonID]bool{}

	// We're going to update every role, every list, and every person.
	for _, r := range tx.FetchRoles() {
		tx.WillUpdateRole(r)
	}
	for _, l := range tx.FetchLists() {
		tx.WillUpdateList(l)
	}
	for _, p := range tx.FetchPeople() {
		tx.WillUpdatePerson(p)
	}
	// Clean out the computed data so that it can be recomputed.
	for _, r := range tx.FetchRoles() {
		for irid, direct := range r.Implies {
			if !direct || tx.FetchRole(irid) == nil {
				delete(r.Implies, irid)
			}
		}
		for lid, rtl := range r.Lists {
			if rtl == 0 || tx.FetchList(lid) == nil {
				delete(r.Lists, lid)
			}
		}
		r.People = r.People[:0]
	}
	for _, l := range tx.FetchLists() {
		manualSubs[l.ID] = make(map[model.PersonID]bool)
		for pid := range l.People {
			if l.People[pid]&model.ListSubscribed != 0 {
				manualSubs[l.ID][pid] = true
			}
			l.People[pid] &^= model.ListSender | model.ListSubscribed
		}
	}
	for _, p := range tx.FetchPeople() {
		for r, direct := range p.Roles {
			if role := tx.FetchRole(r); role == nil || role.ImplicitOnly || !direct {
				delete(p.Roles, r)
			}
		}
		p.Orgs = make([]model.OrgMembership, model.NumOrgs)
	}
	// Add recursive indirect role implications based on current direct
	// implications.
	var seen = map[model.Role2ID]bool{}
	for _, r := range tx.FetchRoles() {
		addIndirectImplications(tx, r, seen)
	}
	// Add indirect role implications to people.
	for _, p := range tx.FetchPeople() {
		for rid := range p.Roles {
			r := tx.FetchRole(rid)
			for irid := range r.Implies {
				if !p.Roles[irid] {
					p.Roles[irid] = false
				}
			}
		}
	}
	// Fill in org memberships of people, and people list of each role.
	for _, p := range tx.FetchPeople() {
		var roles = make([]*model.Role2, 0, len(p.Roles))
		for rid := range p.Roles {
			r := tx.FetchRole(rid)
			roles = append(roles, r)
			if r.PrivLevel > model.PrivNone && p.Orgs[r.Org].PrivLevel == model.PrivNone {
				r.People = append(r.People, p.ID)
			}
			if p.Orgs[r.Org].PrivLevel < r.PrivLevel {
				p.Orgs[r.Org].PrivLevel = r.PrivLevel
			}
		}
		sort.Slice(roles, func(i, j int) bool {
			if roles[i].Org != roles[j].Org {
				return roles[i].Org < roles[j].Org
			}
			return roles[i].Priority < roles[j].Priority
		})
		for _, r := range roles {
			if p.Orgs[r.Org].Title == "" && r.Title != "" {
				p.Orgs[r.Org].Title = r.Title
			}
		}
	}
	// Populate the senders and subscribers of lists.
	for _, p := range tx.FetchPeople() {
		for rid := range p.Roles {
			r := tx.FetchRole(rid)
			for lid, rtl := range r.Lists {
				l := tx.FetchList(lid)
				if rtl.Sender() {
					l.People[p.ID] |= model.ListSender
				}
				switch rtl.SubModel() {
				case model.ListAllowSub:
					if manualSubs[lid][p.ID] {
						l.People[p.ID] |= model.ListSubscribed
					}
				case model.ListAutoSub, model.ListWarnUnsub:
					l.People[p.ID] |= model.ListSubscribed
				}
			}
		}
	}
	// Save everything.
	for _, r := range tx.FetchRoles() {
		tx.UpdateRole(r)
	}
	for _, l := range tx.FetchLists() {
		tx.UpdateList(l)
	}
	for _, p := range tx.FetchPeople() {
		tx.UpdatePerson(p)
	}
}

// addIndirectImplications takes a role with only direct implications and
// recursively fills in its indirect implications.
func addIndirectImplications(tx *store.Tx, r *model.Role2, seen map[model.Role2ID]bool) {
	if seen[r.ID] {
		return
	}
	for irid, direct := range r.Implies {
		if !direct {
			continue
		}
		ir := tx.FetchRole(irid)
		addIndirectImplications(tx, ir, seen)
		for indirect := range ir.Implies {
			r.Implies[indirect] = false
		}
	}
	seen[r.ID] = true
}
