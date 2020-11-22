package db

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchRoles retrieves all of the roles from the database.
func (tx *Tx) FetchRoles() []*model.Role2 {
	var (
		data  []byte
		roles model.Roles
	)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM roles`).Scan(&data))
	panicOnError(roles.Unmarshal(data))
	return roles.Roles
}

// SaveRoles saves the list of roles in the database.
func (tx *Tx) SaveRoles(list []*model.Role2) {
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
