package store

// FetchAuthorizer fetches the authorizer data from the database.  The authz
// package handles unmarshaling it.
func (tx *Tx) FetchAuthorizer() (data []byte) {
	return tx.tx.FetchAuthorizer()
}

// SaveAuthorizer saves all of the authorizer data to the database.  The authz
// package handles marshaling it.
func (tx *Tx) SaveAuthorizer(data []byte) {
	tx.tx.SaveAuthorizer(data)
	// Authorizer creates its own log changes.
}
