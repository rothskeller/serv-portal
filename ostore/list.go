package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateList creates a new list in the database, with the next available ID.
func (tx *Tx) CreateList(list *model.List) {
	tx.Tx.CreateList(list)
	tx.entry.Change("create list [%d]", list.ID)
	tx.entry.Change("set list [%d] type to %s", list.ID, model.ListTypeNames[list.Type])
	tx.entry.Change("set list [%d] name to %q", list.ID, list.Name)
	for pid, lps := range list.People {
		if lps&model.ListSubscribed != 0 {
			tx.entry.Change("add list [%d] subscriber %q [%d]", list.ID, tx.FetchPerson(pid).InformalName, pid)
		}
		if lps&model.ListUnsubscribed != 0 {
			tx.entry.Change("add list [%d] unsubscribe %q [%d]", list.ID, tx.FetchPerson(pid).InformalName, pid)
		}
		if lps&model.ListSender != 0 {
			tx.entry.Change("add list [%d] sender %q [%d]", list.ID, tx.FetchPerson(pid).InformalName, pid)
		}
	}
}

// WillUpdateList saves a copy of a list before it's updated, so that we can
// compare against it to generate audit log entries.
func (tx *Tx) WillUpdateList(l *model.List) {
	if tx.originalLists[l.ID] != nil {
		return
	}
	var ol = *l
	ol.People = make(map[model.PersonID]model.ListPersonStatus, len(l.People))
	for pid, lps := range l.People {
		ol.People[pid] = lps
	}
	tx.originalLists[l.ID] = &ol
}

// UpdateList updates a list in the database.
func (tx *Tx) UpdateList(l *model.List) {
	var ol = tx.originalLists[l.ID]

	if ol == nil {
		panic("must call WillUpdateList before UpdateList")
	}
	tx.Tx.UpdateList(l)
	if l.Type != ol.Type {
		tx.entry.Change("set list [%d] type to %s", l.ID, model.ListTypeNames[l.Type])
	}
	if l.Name != ol.Name {
		tx.entry.Change("set list [%d] name to %q", l.ID, l.Name)
	}
	for pid, lps := range l.People {
		olps := ol.People[pid]
		if lps&model.ListSubscribed != olps&model.ListSubscribed {
			if lps&model.ListSubscribed != 0 {
				tx.entry.Change("add list [%d] subscriber %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			} else {
				tx.entry.Change("remove list [%d] subscriber %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			}
		}
		if lps&model.ListSubscribed != olps&model.ListSubscribed {
			if lps&model.ListUnsubscribed != 0 {
				tx.entry.Change("add list [%d] unsubscribe %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			} else {
				tx.entry.Change("remove list [%d] unsubscribe %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			}
		}
		if lps&model.ListSender != olps&model.ListSender {
			if lps&model.ListSender != 0 {
				tx.entry.Change("add list [%d] sender %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			} else {
				tx.entry.Change("remove list [%d] sender %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
			}
		}
	}
	for pid, olps := range ol.People {
		if _, ok := l.People[pid]; ok {
			continue
		}
		if olps&model.ListSubscribed != 0 {
			tx.entry.Change("remove list [%d] subscriber %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
		}
		if olps&model.ListUnsubscribed != 0 {
			tx.entry.Change("remove list [%d] unsubscribe %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
		}
		if olps&model.ListSender != 0 {
			tx.entry.Change("remove list [%d] sender %q [%d]", l.ID, tx.FetchPerson(pid).InformalName, pid)
		}
	}
	delete(tx.originalLists, l.ID)
}

// DeleteList deletes a list from the database.
func (tx *Tx) DeleteList(list *model.List) {
	tx.Tx.DeleteList(list)
	tx.entry.Change("delete list %q [%d]", list.Name, list.ID)
}
