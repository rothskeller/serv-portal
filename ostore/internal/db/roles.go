package db

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchRoles retrieves all of the roles from the database.
func (tx *Tx) FetchRoles() []*model.Role {
	var (
		data  []byte
		roles model.Roles
	)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM roles`).Scan(&data))
	panicOnError(roles.Unmarshal(data))
	return roles.Roles
}

// SaveRoles saves the list of roles in the database.
func (tx *Tx) SaveRoles(list []*model.Role) {
	var (
		roles model.Roles
		data  []byte
		err   error
	)
	roles.Roles = list
	data, err = roles.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE roles SET data=?`, data))
}

// indexRole updates the search index with information about a roster role.
func (tx *Tx) IndexRole(r *model.Role, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='role' AND id=?`, r.ID))
	}
	if !r.ShowRoster {
		return
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, roleName, roleTitle) VALUES ('role',?,?,?)`, r.ID, r.Name, r.Title))
}
