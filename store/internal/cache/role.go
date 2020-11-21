package cache

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) cacheRoles() {
	if tx.roles != nil {
		return
	}
	tx.roleList = tx.Tx.FetchRoles()
	tx.roles = make(map[model.Role2ID]*model.Role2, len(tx.roleList))
	for _, v := range tx.roleList {
		tx.roles[v.ID] = v
		if v.Implies == nil {
			v.Implies = make(map[model.Role2ID]bool)
		}
		if v.Lists == nil {
			v.Lists = make(map[model.ListID]model.RoleToList)
		}
	}
}

// FetchRole retrieves a single role from the database.  It returns nil if the
// specified role doesn't exist.
func (tx *Tx) FetchRole(id model.Role2ID) *model.Role2 {
	tx.cacheRoles()
	return tx.roles[id]
}

// FetchRoles retrieves all of the roles from the database.
func (tx *Tx) FetchRoles() []*model.Role2 {
	tx.cacheRoles()
	return tx.roleList
}

// CreateRole creates a new role in the database, with the next available ID.
func (tx *Tx) CreateRole(role *model.Role2) {
	tx.cacheRoles()
	for role.ID = 1; tx.roles[role.ID] != nil; role.ID++ {
	}
	tx.roleList = append(tx.roleList, role)
	sort.Sort(model.Roles{Roles: tx.roleList})
	tx.roles[role.ID] = role
	tx.rolesDirty = true
}

// UpdateRole updates an existing role in the database.
func (tx *Tx) UpdateRole(role *model.Role2) {
	tx.cacheRoles()
	if role != tx.roles[role.ID] {
		panic("role must be updated in place")
	}
	sort.Sort(model.Roles{Roles: tx.roleList})
	tx.rolesDirty = true
}

// DeleteRole deletes a role from the database.
func (tx *Tx) DeleteRole(role *model.Role2) {
	tx.cacheRoles()
	if role != tx.roles[role.ID] {
		panic("deleting role that is not in cache")
	}
	delete(tx.roles, role.ID)
	j := 0
	for _, l := range tx.roleList {
		if l != role {
			tx.roleList[j] = l
			j++
		}
	}
	tx.roleList = tx.roleList[:j]
	tx.rolesDirty = true
}
