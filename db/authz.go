package db

func (tx *Tx) CreateAuthorizer(data []byte) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM authorizer`))
	panicOnExecError(tx.tx.Exec(`INSERT INTO authorizer (data) VALUES (?)`, data))
}

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
