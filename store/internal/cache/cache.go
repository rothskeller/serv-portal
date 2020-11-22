// Package cache handles caching for the data store code for the SERV portal.
package cache

import (
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store/internal/db"
)

// Open opens the database.
func Open(path string) {
	db.Open(path)
}

// Tx is a handle to a transaction on the data store.
type Tx struct {
	*db.Tx
	folders    map[model.FolderID]*model.FolderNode
	rootFolder *model.FolderNode
	lists      map[model.ListID]*model.List
	listList   []*model.List
	listsDirty bool
	people     map[model.PersonID]*model.Person
	personList []*model.Person
	roles      map[model.Role2ID]*model.Role2
	roleList   []*model.Role2
	rolesDirty bool
	venues     map[model.VenueID]*model.Venue
	venueList  []*model.Venue
}

// Begin starts a transaction, returning our Tx wrapper.
func Begin() (tx *Tx) {
	return &Tx{Tx: db.Begin(), people: make(map[model.PersonID]*model.Person)}
}

// Commit flushes the lists and roles, if changed, before committing.
func (tx *Tx) Commit() {
	if tx.listsDirty {
		tx.Tx.SaveLists(tx.listList)
	}
	if tx.rolesDirty {
		tx.Tx.SaveRoles(tx.roleList)
	}
	tx.Tx.Commit()
}
