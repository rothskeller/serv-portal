package db

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) cacheAuth() {
	var (
		data  []byte
		authz model.AuthzData
	)
	tx.groups = make(map[model.GroupID]*model.Group)
	tx.groupTags = make(map[model.GroupTag]*model.Group)
	tx.roles = make(map[model.RoleID]*model.Role)
	tx.roleTags = make(map[model.RoleTag]*model.Role)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM authz`).Scan(&data))
	panicOnError(authz.Unmarshal(data))
	tx.groupList = authz.Groups
	tx.roleList = authz.Roles
	for _, group := range tx.groupList {
		tx.groups[group.ID] = group
		if group.Tag != "" {
			tx.groupTags[group.Tag] = group
		}
	}
	for _, role := range tx.roleList {
		tx.roles[role.ID] = role
		if role.Tag != "" {
			tx.roleTags[role.Tag] = role
		}
	}
}

// FetchRole retrieves a single role from the database.  It returns nil if the
// specified role doesn't exist.
func (tx *Tx) FetchRole(id model.RoleID) *model.Role {
	return tx.roles[id]
}

// FetchRoleByTag retrieves the role with the specified tag from the database.
// It returns nil if no such team exists.
func (tx *Tx) FetchRoleByTag(tag model.RoleTag) *model.Role {
	return tx.roleTags[tag]
}

// FetchRoles retrieves all of the roles from the database, in order by name.
func (tx *Tx) FetchRoles() []*model.Role {
	return tx.roleList
}

// FetchGroup retrieves a single group from the database.  It returns nil if the
// specified group doesn't exist.
func (tx *Tx) FetchGroup(id model.GroupID) *model.Group {
	return tx.groups[id]
}

// FetchGroupByTag retrieves the group with the specified tag from the database.
// It returns nil if no such team exists.
func (tx *Tx) FetchGroupByTag(tag model.GroupTag) *model.Group {
	return tx.groupTags[tag]
}

// FetchGroups retrieves all of the groups from the database, in order by name.
func (tx *Tx) FetchGroups() []*model.Group {
	return tx.groupList
}

// CreateGroup adds the supplied group to the list of groups, giving it an
// available ID.  This call does not persist the change; call SaveAuthz later
// for that.
func (tx *Tx) CreateGroup(group *model.Group) {
	for group.ID = 1; tx.groups[group.ID] != nil; group.ID++ {
	}
	tx.groupList = append(tx.groupList, group)
	tx.groups[group.ID] = group
	// Since this may be a re-used ID, we need to make sure none of the
	// roles have this group in their map.  This will also enlarge their
	// maps for the new group if needed.
	for _, r := range tx.roles {
		r.Privileges.Set(group, 0)
	}
}

// DeleteGroup deletes the supplied group from the list of groups.  It does not
// persist the change; call SaveAuthz later for that.
func (tx *Tx) DeleteGroup(group *model.Group) {
	j := 0
	for _, g := range tx.groupList {
		if g != group {
			tx.groupList[j] = g
			j++
		}
	}
	tx.groupList = tx.groupList[:j]
	delete(tx.groups, group.ID)
	delete(tx.groupTags, group.Tag)
}

// CreateRole adds the supplied role to the list of roles, giving it an
// available ID.  This call does not persist the change; call SaveAuthz later
// for that.
func (tx *Tx) CreateRole(role *model.Role) {
	for role.ID = 1; tx.roles[role.ID] != nil; role.ID++ {
	}
	tx.roleList = append(tx.roleList, role)
	tx.roles[role.ID] = role
}

// DeleteRole deletes the supplied role from the list of roles.  It does not
// persist the change; call SaveAuthz later for that.
func (tx *Tx) DeleteRole(role *model.Role) {
	j := 0
	for _, g := range tx.roleList {
		if g != role {
			tx.roleList[j] = g
			j++
		}
	}
	tx.roleList = tx.roleList[:j]
	delete(tx.roles, role.ID)
	delete(tx.roleTags, role.Tag)
}

// SaveAuthz saves the entire set of groups, roles, and privileges to the
// database.  It also recalculates the privilege maps for every person in the
// database.  It's a very expensive operation, but it's only called when making
// changes to roles and groups, so that should be OK.
func (tx *Tx) SaveAuthz() {
	var (
		authz model.AuthzData
		data  []byte
		err   error
	)
	sort.Sort(model.GroupSort(tx.groupList))
	authz.Groups = tx.groupList
	sort.Sort(model.RoleSort(tx.roleList))
	authz.Roles = tx.roleList
	data, err = authz.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE authz SET data=?`, data))
	tx.audit("authz", 0, data)
	tx.recalcAllPersonPrivileges()
}
