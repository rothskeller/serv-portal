package db

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchVenues retrieves all of the venues from the database.
func (tx *Tx) FetchVenues() []*model.Venue {
	var (
		data   []byte
		venues model.Venues
	)
	panicOnError(tx.tx.QueryRow(`SELECT data FROM venue`).Scan(&data))
	panicOnError(venues.Unmarshal(data))
	return venues.Venues
}

// SaveVenues saves the list of venues in the database.
func (tx *Tx) SaveVenues(list []*model.Venue) {
	var (
		venues model.Venues
		data   []byte
		err    error
	)
	venues.Venues = list
	data, err = venues.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE venue SET data=?`, data))
}
