package main

import (
	"sunnyvaleserv.org/portal/model"
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
)

func convertLists(tx *ostore.Tx, st *nstore.Store) {
	for _, olist := range tx.FetchLists() {
		var nlist = list.Updater{
			ID:   list.ID(olist.ID),
			Type: list.Type(olist.Type),
			Name: olist.Name,
		}
		list.Create(st, &nlist)
	}
}

func convertSubscriptions(tx *ostore.Tx, st *nstore.Store) {
	for _, olist := range tx.FetchLists() {
		nl := list.WithID(st, list.ID(olist.ID))
		for pid, lps := range olist.People {
			np := person.WithID(st, person.ID(pid), person.FID|person.FInformalName)
			if lps&model.ListSubscribed != 0 {
				listperson.Subscribe(st, nl, np)
			}
			if lps&model.ListUnsubscribed != 0 {
				listperson.Unsubscribe(st, nl, np)
			}
		}
	}
}
