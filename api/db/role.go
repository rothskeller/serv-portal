package db

import (
	"database/sql"
	"sort"

	"rothskeller.net/serv/model"
)

func (tx *Tx) cacheRoles() {
	var (
		rows *sql.Rows
		err  error
	)
	tx.roles = make(map[model.RoleID]*model.Role)
	rows, err = tx.tx.Query(`SELECT id, tag, name, member_label, imply_only, individual, privileges FROM role`)
	panicOnError(err)
	for rows.Next() {
		var (
			role model.Role
			tag  sql.NullString
		)
		panicOnError(rows.Scan(&role.ID, &tag, &role.Name, &role.MemberLabel, &role.ImplyOnly, &role.Individual, &role.PrivMap))
		role.Tag = model.RoleTag(tag.String)
		tx.roles[role.ID] = &role
	}
	panicOnError(rows.Err())
	tx.roleList = make([]*model.Role, 0, len(tx.roles))
	tx.recalcRoleCache()
}
func (tx *Tx) recalcRoleCache() {
	tx.maxRoleID = 0
	tx.roleList = tx.roleList[:0]
	tx.roleTags = make(map[model.RoleTag]*model.Role)
	for _, role := range tx.roles {
		if role.Tag != "" {
			tx.roleTags[role.Tag] = role
		}
		tx.roleList = append(tx.roleList, role)
		if tx.maxRoleID < role.ID {
			tx.maxRoleID = role.ID
		}
	}
	sort.Sort(model.RoleSort(tx.roleList))
	if !model.RecalculateTransitivePrivilegeMaps(tx.roleList) {
		panic("cycle in role graph")
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

// FetchRoles retrieves all of the roles from the database.
func (tx *Tx) FetchRoles() []*model.Role {
	return tx.roleList
}

// UnusedRoleID returns a role ID for a new role.
func (tx *Tx) UnusedRoleID() (roleID model.RoleID) {
	for roleID = 1; tx.roles[roleID] != nil; roleID++ {
	}
	return roleID
}

// SaveRole saves a role definition to the database.  If its supplied ID is
// zero, it creates a new role in the database, and puts its ID in the supplied
// role structure.  Note that this does not save the role privileges; do that
// with a SavePrivileges call.
func (tx *Tx) SaveRole(role *model.Role) {
	var tag = sql.NullString{Valid: role.Tag != "", String: string(role.Tag)}

	if tx.roles[role.ID] == nil {
		panicOnExecError(tx.tx.Exec(`INSERT INTO role (id, tag, name, member_label, imply_only, individual, privileges) VALUES (?,?,?,?,?,?,x'')`, role.ID, tag, role.Name, role.MemberLabel, role.ImplyOnly, role.Individual))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE role SET tag=?, name=?, member_label=?, imply_only=?, individual=? WHERE id=?`, tag, role.Name, role.MemberLabel, role.ImplyOnly, role.Individual, role.ID))
	}
	tx.roles[role.ID] = role
	tx.recalcRoleCache()
	tx.audit(model.AuditRecord{Role: role})
}

// DeleteRole deletes a role definition from the database.
func (tx *Tx) DeleteRole(role *model.Role) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM role WHERE id=?`, role.ID))
	delete(tx.roles, role.ID)
	tx.recalcRoleCache()
	tx.audit(model.AuditRecord{Role: &model.Role{ID: role.ID}})
}

// SavePrivileges saves all of the role privileges of all roles in the system.
func (tx *Tx) SavePrivileges() {
	var (
		stmt1 *sql.Stmt
		stmt2 *sql.Stmt
		err   error
	)
	panicOnExecError(tx.tx.Exec(`DELETE FROM role_privilege`))
	stmt1, err = tx.tx.Prepare(`INSERT INTO role_privilege (actor, target, privileges) VALUES (?,?,?)`)
	panicOnError(err)
	stmt2, err = tx.tx.Prepare(`UPDATE role SET privileges=? WHERE id=?`)
	panicOnError(err)
	for _, actor := range tx.roles {
		for _, target := range tx.roles {
			if priv := actor.PrivMap.Get(target.ID); priv != 0 {
				panicOnExecError(stmt1.Exec(actor.ID, target.ID, priv))
			}
		}
		panicOnNoRows(stmt2.Exec(actor.PrivMap, actor.ID))
	}
	panicOnError(stmt1.Close())
	panicOnError(stmt2.Close())
}
