package store

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchSccAresEventNames retrieves the name rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventNames() (rewrites map[string]string) {
	return tx.tx.FetchSccAresEventNames()
}

// FetchSccAresEventVenues retrieves the location rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventVenues() (venues map[string]*model.Venue) {
	var vmap = tx.tx.FetchSccAresEventVenues()
	tx.cacheVenues()
	venues = make(map[string]*model.Venue)
	for scc, serv := range vmap {
		venues[scc] = tx.venues[serv]
	}
	return venues
}

// FetchSccAresEventTypes retrieves the event type rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventTypes() (types map[string]model.EventType) {
	return tx.tx.FetchSccAresEventTypes()
}
