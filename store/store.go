// Package store contains the data store code for the SERV portal.  This handles
// caching, auditing of changes, and data storage.
package store

import (
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store/authz"
	"sunnyvaleserv.org/portal/store/internal/cache"
	"sunnyvaleserv.org/portal/util/log"
)

// Open opens the database.
func Open(path string) {
	cache.Open(path)
}

// Tx is a handle to a transaction on the data store.
type Tx struct {
	*cache.Tx
	entry          *log.Entry
	auth           *authz.Authorizer
	originalLists  map[model.ListID]*model.List
	originalPeople map[model.PersonID]*model.Person
	originalRoles  map[model.Role2ID]*model.Role2
}

// Begin starts a transaction, returning our Tx wrapper.
func Begin(entry *log.Entry) (tx *Tx) {
	return &Tx{
		Tx:             cache.Begin(),
		entry:          entry,
		originalLists:  make(map[model.ListID]*model.List),
		originalPeople: make(map[model.PersonID]*model.Person),
		originalRoles:  make(map[model.Role2ID]*model.Role2),
	}
}

// Authorizer returns an authorizer for the transaction.
func (tx *Tx) Authorizer() *authz.Authorizer {
	if tx.auth == nil {
		tx.auth = authz.NewAuthorizer(tx.Tx, tx.entry)
	}
	return tx.auth
}
