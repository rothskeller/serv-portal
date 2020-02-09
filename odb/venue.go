package db

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

func (tx *Tx) cacheVenues() {
	var (
		data   []byte
		venues model.Venues
	)
	tx.venues = make(map[model.VenueID]*model.Venue)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM venue`).Scan(&data))
	panicOnError(venues.Unmarshal(data))
	tx.venueList = venues.Venues
	for _, venue := range tx.venueList {
		tx.venues[venue.ID] = venue
	}
}

// FetchVenue retrieves a single venue from the database.  It returns nil if the
// specified venue doesn't exist.
func (tx *Tx) FetchVenue(id model.VenueID) *model.Venue {
	if tx.venues == nil {
		tx.cacheVenues()
	}
	return tx.venues[id]
}

// FetchVenues retrieves all of the venues from the database.
func (tx *Tx) FetchVenues() []*model.Venue {
	if tx.venues == nil {
		tx.cacheVenues()
	}
	return tx.venueList
}

// SaveVenue saves a venue definition to the database.  If its supplied ID is
// zero, it creates a new venue in the database, and puts its ID in the supplied
// venue structure.
func (tx *Tx) SaveVenue(venue *model.Venue) {
	var (
		venues model.Venues
		data   []byte
		err    error
	)
	if tx.venues == nil {
		tx.cacheVenues()
	}
	if venue.ID == 0 {
		for venue.ID = 1; tx.venues[venue.ID] != nil; venue.ID++ {
		}
		tx.venueList = append(tx.venueList, venue)
	}
	tx.venues[venue.ID] = venue
	venues.Venues = tx.venueList
	sort.Sort(venues)
	data, err = venues.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE venue SET data=?`, data))
	tx.audit("venues", 0, data)
}

// DeleteVenue deletes a venue definition from the database.
func (tx *Tx) DeleteVenue(venue *model.Venue) {
	var (
		venues model.Venues
		data   []byte
		err    error
	)
	if tx.venues == nil {
		tx.cacheVenues()
	}
	delete(tx.venues, venue.ID)
	j := 0
	for _, v := range tx.venueList {
		if v != venue {
			tx.venueList[j] = v
			j++
		}
	}
	tx.venueList = tx.venueList[:j]
	venues.Venues = tx.venueList
	data, err = venues.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE venue SET data=?`, data))
	tx.audit("venues", 0, data)
}
