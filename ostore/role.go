package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateRole creates a new role in the database, with the next available ID.
func (tx *Tx) CreateRole(role *model.Role) {
	tx.Tx.CreateRole(role)
	tx.entry.Change("create role [%d]", role.ID)
	tx.entry.Change("set role [%d] name to %q", role.ID, role.Name)
	if role.Title != "" {
		tx.entry.Change("set role [%d] title to %q", role.ID, role.Title)
	}
	tx.entry.Change("set role [%d] org to %s", role.ID, role.Org.String())
	if role.PrivLevel != model.PrivNone {
		tx.entry.Change("set role [%d] privLevel to %s", role.ID, role.PrivLevel)
	}
	if role.ShowRoster {
		tx.entry.Change("set role [%d] showRoster", role.ID)
	}
	if role.ImplicitOnly {
		tx.entry.Change("set role [%d] implicitOnly", role.ID)
	}
	tx.entry.Change("set role [%d] priority to %d", role.ID, role.Priority)
	for irid, direct := range role.Implies {
		if direct {
			tx.entry.Change("set role [%d] implies role %q [%d]", role.ID, tx.FetchRole(irid).Name, irid)
		}
	}
	for lid, rtl := range role.Lists {
		if rtl.SubModel() != model.ListNoSub {
			tx.entry.Change("set role [%d] list %q [%d] subModel to %s", role.ID, tx.FetchList(lid).Name, lid, model.ListSubModelNames[rtl.SubModel()])
		}
		if rtl.Sender() {
			tx.entry.Change("set role [%d] list %q [%d] sender", role.ID, tx.FetchList(lid).Name, lid)
		}
	}
}

// WillUpdateRole saves a copy of a role before it's updated, so that we can
// compare against it to generate audit log entries.
func (tx *Tx) WillUpdateRole(r *model.Role) {
	if tx.originalRoles[r.ID] != nil {
		return
	}
	var or = *r
	or.Implies = make(map[model.RoleID]bool, len(r.Implies))
	for irid, direct := range r.Implies {
		or.Implies[irid] = direct
	}
	or.Lists = make(map[model.ListID]model.RoleToList)
	for lid, rtl := range r.Lists {
		or.Lists[lid] = rtl
	}
	tx.originalRoles[r.ID] = &or
}

// UpdateRole updates a role in the database.
func (tx *Tx) UpdateRole(r *model.Role) {
	var or = tx.originalRoles[r.ID]

	if or == nil {
		panic("must call WillUpdateRole before UpdateRole")
	}
	tx.Tx.UpdateRole(r)
	if r.Name != or.Name {
		tx.entry.Change("set role [%d] name to %q", r.ID, r.Name)
	}
	if r.Title != or.Title {
		tx.entry.Change("set role [%d] title to %q", r.ID, r.Title)
	}
	if r.Org != or.Org {
		tx.entry.Change("set role [%d] org to %s", r.ID, r.Org.String())
	}
	if r.PrivLevel != or.PrivLevel {
		if r.PrivLevel != model.PrivNone {
			tx.entry.Change("set role [%d] privLevel to %s", r.ID, r.PrivLevel)
		} else {
			tx.entry.Change("clear role [%d] privLevel", r.ID)
		}
	}
	if r.ShowRoster != or.ShowRoster {
		if r.ShowRoster {
			tx.entry.Change("set role [%d] showRoster", r.ID)
		} else {
			tx.entry.Change("clear role [%d] showRoster", r.ID)
		}
	}
	if r.ImplicitOnly != or.ImplicitOnly {
		if r.ImplicitOnly {
			tx.entry.Change("set role [%d] implicitOnly", r.ID)
		} else {
			tx.entry.Change("clear role [%d] implicitOnly", r.ID)
		}
	}
	if r.Priority != or.Priority {
		tx.entry.Change("set role [%d] priority to %d", r.ID, r.Priority)
	}
	for irid, direct := range r.Implies {
		if direct {
			if odirect, ok := or.Implies[irid]; !ok || !odirect {
				tx.entry.Change("set role [%d] implies role %q [%d]", r.ID, tx.FetchRole(irid).Name, irid)
			}
		} else if odirect, ok := or.Implies[irid]; ok && odirect {
			tx.entry.Change("clear role [%d] implies role %q [%d]", r.ID, tx.FetchRole(irid).Name, irid)
		}
	}
	for irid, odirect := range or.Implies {
		if odirect {
			if _, ok := r.Implies[irid]; !ok {
				tx.entry.Change("clear role [%d] implies role %q [%d]", r.ID, tx.FetchRole(irid).Name, irid)
			}
		}
	}
	for lid, rtl := range r.Lists {
		ortl := or.Lists[lid]
		if rtl.SubModel() != ortl.SubModel() {
			if rtl.SubModel() != model.ListNoSub {
				tx.entry.Change("set role [%d] list %q [%d] subModel to %s", r.ID, tx.FetchList(lid).Name, lid, model.ListSubModelNames[rtl.SubModel()])
			} else {
				tx.entry.Change("clear role [%d] list %q [%d] subModel", r.ID, tx.FetchList(lid).Name, lid)
			}
		}
		if rtl.Sender() != ortl.Sender() {
			if rtl.Sender() {
				tx.entry.Change("set role [%d] list %q [%d] sender", r.ID, tx.FetchList(lid).Name, lid)
			} else {
				tx.entry.Change("clear role [%d] list %q [%d] sender", r.ID, tx.FetchList(lid).Name, lid)
			}
		}
	}
	for lid, ortl := range or.Lists {
		if _, ok := r.Lists[lid]; ok {
			continue
		}
		if ortl.SubModel() != model.ListNoSub {
			tx.entry.Change("clear role [%d] list %q [%d] subModel", r.ID, tx.FetchList(lid).Name, lid)
		}
		if ortl.Sender() {
			tx.entry.Change("clear role [%d] list %q [%d] sender", r.ID, tx.FetchList(lid).Name, lid)
		}
	}
	delete(tx.originalRoles, r.ID)
}

// DeleteRole deletes a role from the database.
func (tx *Tx) DeleteRole(role *model.Role) {
	tx.Tx.DeleteRole(role)
	tx.entry.Change("delete role %q [%d]", role.Name, role.ID)
}
