package cache

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) cacheLists() {
	if tx.lists != nil {
		return
	}
	tx.listList = tx.Tx.FetchLists()
	tx.lists = make(map[model.ListID]*model.List, len(tx.listList))
	for _, v := range tx.listList {
		tx.lists[v.ID] = v
		if v.People == nil {
			v.People = make(map[model.PersonID]model.ListPersonStatus)
		}
	}
}

// FetchList retrieves a single list from the database.  It returns nil if the
// specified list doesn't exist.
func (tx *Tx) FetchList(id model.ListID) *model.List {
	tx.cacheLists()
	return tx.lists[id]
}

// FetchLists retrieves all of the lists from the database.
func (tx *Tx) FetchLists() []*model.List {
	tx.cacheLists()
	return tx.listList
}

// CreateList creates a new list in the database, with the next available ID.
func (tx *Tx) CreateList(list *model.List) {
	tx.cacheLists()
	for list.ID = 1; tx.lists[list.ID] != nil; list.ID++ {
	}
	tx.listList = append(tx.listList, list)
	sort.Sort(model.Lists{Lists: tx.listList})
	tx.lists[list.ID] = list
	tx.listsDirty = true
}

// UpdateList updates an existing list in the database.
func (tx *Tx) UpdateList(list *model.List) {
	tx.cacheLists()
	if list != tx.lists[list.ID] {
		panic("list must be updated in place")
	}
	sort.Sort(model.Lists{Lists: tx.listList})
	tx.listsDirty = true
}

// DeleteList deletes a list from the database.
func (tx *Tx) DeleteList(list *model.List) {
	tx.cacheLists()
	if list != tx.lists[list.ID] {
		panic("deleting list that is not in cache")
	}
	delete(tx.lists, list.ID)
	j := 0
	for _, l := range tx.listList {
		if l != list {
			tx.listList[j] = l
			j++
		}
	}
	tx.listList = tx.listList[:j]
	tx.listsDirty = true
}
