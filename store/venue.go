package store

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) cacheVenues() {
	if tx.venues != nil {
		return
	}
	tx.venueList = tx.tx.FetchVenues()
	tx.venues = make(map[model.VenueID]*model.Venue, len(tx.venueList))
	for _, v := range tx.venueList {
		tx.venues[v.ID] = v
	}
}

// FetchVenue retrieves a single venue from the database.  It returns nil if the
// specified venue doesn't exist.
func (tx *Tx) FetchVenue(id model.VenueID) *model.Venue {
	tx.cacheVenues()
	return tx.venues[id]
}

// FetchVenues retrieves all of the venues from the database.
func (tx *Tx) FetchVenues() []*model.Venue {
	tx.cacheVenues()
	return tx.venueList
}

// CreateVenue creates a new venue in the database, with the next available ID.
func (tx *Tx) CreateVenue(venue *model.Venue) {
	tx.cacheVenues()
	for venue.ID = 1; tx.venues[venue.ID] != nil; venue.ID++ {
	}
	tx.venueList = append(tx.venueList, venue)
	sort.Sort(model.Venues{Venues: tx.venueList})
	tx.venues[venue.ID] = venue
	tx.tx.SaveVenues(tx.venueList)
	tx.entry.Change("create venue [%d]", venue.ID)
	tx.entry.Change("set venue [%d] name to %q", venue.Name)
	if venue.Address != "" {
		tx.entry.Change("set venue [%d] address to %q", venue.Address)
	}
	if venue.City != "" {
		tx.entry.Change("set venue [%d] city to %q", venue.City)
	}
	if venue.URL != "" {
		tx.entry.Change("set venue [%d] URL to %q", venue.URL)
	}
}
