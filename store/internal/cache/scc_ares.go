package cache

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchSccAresEventVenues retrieves the location rewriting map for the
// scc-ares-races.org event imports.
func (tx *Tx) FetchSccAresEventVenues() (venues map[string]*model.Venue) {
	var vmap = tx.Tx.FetchSccAresEventVenues()
	tx.cacheVenues()
	venues = make(map[string]*model.Venue)
	for scc, serv := range vmap {
		venues[scc] = tx.venues[serv]
	}
	return venues
}
