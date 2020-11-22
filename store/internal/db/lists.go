package db

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchLists retrieves all of the lists from the database.
func (tx *Tx) FetchLists() []*model.List {
	var (
		data  []byte
		lists model.Lists
	)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM lists`).Scan(&data))
	panicOnError(lists.Unmarshal(data))
	return lists.Lists
}

// SaveLists saves the list of lists in the database.
func (tx *Tx) SaveLists(list []*model.List) {
	var (
		lists model.Lists
		data  []byte
		err   error
	)
	lists.Lists = list
	data, err = lists.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE lists SET data=?`, data))
}
