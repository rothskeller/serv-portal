package main

import (
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/session"
)

func convertSessions(tx *ostore.Tx, st *nstore.Store) {
	for _, osession := range tx.FetchSessions() {
		np := person.WithID(st, person.ID(osession.Person.ID), person.FID|person.FInformalName)
		session.Create(st, np, string(osession.Token), string(osession.CSRF), osession.Expires)
	}
}
