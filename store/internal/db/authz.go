package db

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchAuthorizer fetches the authorizer data from the database.  The authz
// package handles unmarshaling it.
func (tx *Tx) FetchAuthorizer() (data []byte) {
	panicOnError(tx.tx.QueryRow(`SELECT data FROM authorizer`).Scan(&data))
	return data
}

// SaveAuthorizer saves all of the authorizer data to the database.  The authz
// package handles marshaling it.
func (tx *Tx) SaveAuthorizer(data []byte) {
	panicOnNoRows(tx.tx.Exec(`UPDATE authorizer SET data=?`, data))
}

// IndexGroups updates the search index with the current list of groups.
func (tx *Tx) IndexGroups(groups []*model.Group, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='group'`))
	}
	for _, g := range groups {
		panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, groupName) VALUES ('group',?,?)`, g.ID, g.Name))
	}
}
