package db

import (
	"database/sql"
	"sort"

	"rothskeller.net/serv/model"
)

func (tx *Tx) cacheVenues() {
	var (
		rows *sql.Rows
		err  error
	)
	tx.venues = make(map[model.VenueID]*model.Venue)
	rows, err = tx.tx.Query(`SELECT id, name, address, city, url FROM venue`)
	panicOnError(err)
	for rows.Next() {
		var venue model.Venue
		panicOnError(rows.Scan(&venue.ID, &venue.Name, &venue.Address, &venue.City, &venue.URL))
		tx.venues[venue.ID] = &venue
		tx.venueList = append(tx.venueList, &venue)
	}
	panicOnError(rows.Err())
	sort.Slice(tx.venueList, func(i, j int) bool { return tx.venueList[i].Name < tx.venueList[j].Name })
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
	if tx.venues == nil {
		tx.cacheVenues()
	}
	if venue.ID == 0 {
		result, err := tx.tx.Exec(`INSERT INTO venue (name, address, city, url) VALUES (?,?,?,?)`, venue.Name, venue.Address, venue.City, venue.URL)
		panicOnError(err)
		venue.ID = model.VenueID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE venue SET name=?, address=?, city=?, url=? WHERE id=?`, venue.Name, venue.Address, venue.City, venue.URL, venue.ID))
	}
	tx.venues[venue.ID] = venue
	tx.audit(model.AuditRecord{Venue: venue})
}

// DeleteVenue deletes a venue definition from the database.
func (tx *Tx) DeleteVenue(venue *model.Venue) {
	if tx.venues == nil {
		tx.cacheVenues()
	}
	panicOnNoRows(tx.tx.Exec(`DELETE FROM venue WHERE id=?`, venue.ID))
	delete(tx.venues, venue.ID)
	tx.audit(model.AuditRecord{Venue: &model.Venue{ID: venue.ID}})
}
