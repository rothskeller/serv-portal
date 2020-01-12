package db

import (
	"database/sql"

	"rothskeller.net/serv/model"
)

// FetchSccAresEventNames retrieves the name rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventNames() (rewrites map[string]string) {
	var (
		scc  string
		serv string
		rows *sql.Rows
		err  error
	)
	rewrites = make(map[string]string)
	rows, err = tx.tx.Query(`SELECT scc, serv FROM scc_ares_event_name`)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&scc, &serv))
		rewrites[scc] = serv
	}
	panicOnError(rows.Err())
	return rewrites
}

// FetchSccAresEventVenues retrieves the location rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventVenues() (venues map[string]*model.Venue) {
	var (
		scc  string
		serv model.VenueID
		rows *sql.Rows
		err  error
	)
	venues = make(map[string]*model.Venue)
	tx.cacheVenues()
	rows, err = tx.tx.Query(`SELECT scc, serv FROM scc_ares_event_location`)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&scc, (*ID)(&serv)))
		venues[scc] = tx.venues[serv]
	}
	panicOnError(rows.Err())
	return venues
}

// FetchSccAresEventTypes retrieves the event type rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventTypes() (types map[string]model.EventType) {
	var (
		scc  string
		serv model.EventType
		rows *sql.Rows
		err  error
	)
	types = make(map[string]model.EventType)
	rows, err = tx.tx.Query(`SELECT scc, serv FROM scc_ares_event_type`)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&scc, &serv))
		types[scc] = serv
	}
	panicOnError(rows.Err())
	return types
}
