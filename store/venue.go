package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateVenue creates a new venue in the database, with the next available ID.
func (tx *Tx) CreateVenue(venue *model.Venue) {
	tx.Tx.CreateVenue(venue)
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
