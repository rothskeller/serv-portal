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
		var roles = model.Roles{Roles: make([]*model.Role2, 0, len(p.Roles))}
		for rid, direct := range p.Roles {
			r := tx.FetchRole(rid)
			r.People = append(r.People, p.ID)
			if direct {
				roles.Roles = append(roles.Roles, r)
			}
			if p.Orgs[r.Org].PrivLevel < r.PrivLevel && !p.Roles[model.DisabledUser] {
				p.Orgs[r.Org].PrivLevel = r.PrivLevel
			}
		}
		sort.Sort(roles)
		if p.Roles[model.DisabledUser] {
			continue
		}
		for _, r := range roles.Roles {
			if p.Orgs[r.Org].Title == "" && r.Title != "" {
				p.Orgs[r.Org].Title = r.Title
			}
		}
	}
	// Populate the senders and subscribers of lists.  Note that disabled
	// users can still be subscribed to lists, but sending to those lists
	// will omit them.  That way their manual subscriptions are retained if
	// their account is re-enabled.
	for _, p := range tx.FetchPeople() {
		for rid := range p.Roles {
			r := tx.FetchRole(rid)
			for lid, rtl := range r.Lists {
				l := tx.FetchList(lid)
				if rtl.Sender() {
					l.People[p.ID] |= model.ListSender
				}
				if l.People[p.ID]&model.ListUnsubscribed != 0 {
					continue
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

// CanSubscribe returns a map of which lists the specified person is allowed to
// subscribe to.
func CanSubscribe(tx *store.Tx, person *model.Person) (can map[model.ListID]bool) {
	can = make(map[model.ListID]bool)
	for rid := range person.Roles {
		role := tx.FetchRole(rid)
		for lid, rtl := range role.Lists {
			if rtl.SubModel() != model.ListNoSub {
				can[lid] = true
			}
		}
	}
	return can
}
