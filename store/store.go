// Package store contains the data store code for the SERV portal.  This handles
// caching, auditing of changes, and data storage.
package store

import (
	"sunnyvaleserv.org/portal/log"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store/authz"
	"sunnyvaleserv.org/portal/store/internal/db"
)

// Open opens the database.
func Open(path string) {
	db.Open(path)
}

// Tx is a handle to a transaction on the data store.
type Tx struct {
	tx         *db.Tx
	entry      *log.Entry
	auth       *authz.Authorizer
	people     map[model.PersonID]*model.Person
	personList []*model.Person
	venues     map[model.VenueID]*model.Venue
	venueList  []*model.Venue
}

// Begin starts a transaction, returning our Tx wrapper.
func Begin(entry *log.Entry) (tx *Tx) {
	return &Tx{tx: db.Begin(), entry: entry, people: make(map[model.PersonID]*model.Person)}
}

// Authorizer returns an authorizer for the transaction.
func (tx *Tx) Authorizer() *authz.Authorizer {
	if tx.auth == nil {
		tx.auth = authz.NewAuthorizer(tx.tx, tx.entry)
	}
	return tx.auth
}

// Commit commits a transaction.
func (tx *Tx) Commit() {
	tx.tx.Commit()
}

// Rollback rolls back a transaction.
func (tx *Tx) Rollback() {
	if tx.tx != nil {
		tx.tx.Rollback()
	}
}
