package db

import (
	"database/sql"

	"rothskeller.net/serv/model"
)

// FetchEvent retrieves a single event from the database by ID.  It returns nil
// if no such event exists.
func (tx *Tx) FetchEvent(id model.EventID) (e *model.Event) {
	var (
		rows    *sql.Rows
		venueID model.VenueID
		err     error
	)
	e = &model.Event{ID: id}
	err = tx.tx.QueryRow(`SELECT name, date, start, end, venue, type FROM event WHERE id=?`, id).Scan(&e.Name, &e.Date, &e.Start, &e.End, (*ID)(&venueID), &e.Type)
	if err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	if venueID != 0 {
		e.Venue = tx.FetchVenue(venueID)
	}
	rows, err = tx.tx.Query(`SELECT role FROM event_role WHERE event=?`, id)
	panicOnError(err)
	for rows.Next() {
		var role model.RoleID
		panicOnError(rows.Scan(&role))
		e.Roles = append(e.Roles, tx.FetchRole(role))
	}
	panicOnError(rows.Err())
	return e
}

// FetchEvents returns all of the events within the specified time range, in
// chronological order.  The time range is inclusive; each time must be in
// 2006-01-02 format.
func (tx *Tx) FetchEvents(from, to string) (events []*model.Event) {
	var (
		rows    *sql.Rows
		stmt    *sql.Stmt
		venueID model.VenueID
		err     error
	)
	rows, err = tx.tx.Query(`SELECT id, name, date, start, end, venue, type FROM event WHERE date>=? AND date<=? ORDER BY date, start, name`, from, to)
	panicOnError(err)
	for rows.Next() {
		var e model.Event
		panicOnError(rows.Scan(&e.ID, &e.Name, &e.Date, &e.Start, &e.End, (*ID)(&venueID), &e.Type))
		if venueID != 0 {
			e.Venue = tx.FetchVenue(venueID)
		}
		events = append(events, &e)
	}
	panicOnError(rows.Err())
	stmt, err = tx.tx.Prepare(`SELECT role FROM event_role WHERE event=?`)
	panicOnError(err)
	for _, e := range events {
		rows, err = stmt.Query(e.ID)
		panicOnError(err)
		for rows.Next() {
			var role model.RoleID
			panicOnError(rows.Scan(&role))
			e.Roles = append(e.Roles, tx.FetchRole(role))
		}
		panicOnError(rows.Err())
	}
	panicOnError(stmt.Close())
	return events
}

// SaveEvent saves an event to the database.  If the supplied event ID is
// zero, a new event is added to the database; otherwise, the identified event
// is updated.
func (tx *Tx) SaveEvent(e *model.Event) {
	var (
		venueID model.VenueID
		err     error
	)
	if e.Venue != nil {
		venueID = e.Venue.ID
	}
	if e.ID == 0 {
		var result sql.Result
		result, err = tx.tx.Exec(`INSERT INTO event (name, date, start, end, venue, type) VALUES (?,?,?,?,?,?)`, e.Name, e.Date, e.Start, e.End, ID(venueID), e.Type)
		panicOnError(err)
		e.ID = model.EventID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE event SET name=?, date=?, start=?, end=?, venue=?, type=? WHERE id=?`, e.Name, e.Date, e.Start, e.End, ID(venueID), e.Type, e.ID))
		panicOnExecError(tx.tx.Exec(`DELETE FROM event_role WHERE event=?`, e.ID))
	}
	for _, r := range e.Roles {
		panicOnExecError(tx.tx.Exec(`INSERT INTO event_role (event, role) VALUES (?,?)`, e.ID, r.ID))
	}
	tx.audit(model.AuditRecord{Event: e})
}

// DeleteEvent deletes an event from the database.
func (tx *Tx) DeleteEvent(e *model.Event) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
	tx.audit(model.AuditRecord{Event: &model.Event{ID: e.ID}})
}

// FetchAttendanceByEvent retrieves the attendance at a specific event.
func (tx *Tx) FetchAttendanceByEvent(e *model.Event) (attend map[model.PersonID]bool) {
	var (
		rows *sql.Rows
		pid  model.PersonID
		err  error
	)
	attend = make(map[model.PersonID]bool)
	rows, err = tx.tx.Query(`SELECT person FROM attendance WHERE event=?`, e.ID)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&pid))
		attend[pid] = true
	}
	panicOnError(rows.Err())
	return attend
}

// SaveEventAttendance saves the attendance for a specific event.
func (tx *Tx) SaveEventAttendance(e *model.Event, people []*model.Person) {
	var (
		stmt *sql.Stmt
		err  error
	)
	panicOnExecError(tx.tx.Exec(`DELETE FROM attendance WHERE event=?`, e.ID))
	stmt, err = tx.tx.Prepare(`INSERT INTO attendance (event, person) VALUES (?,?)`)
	panicOnError(err)
	for _, p := range people {
		panicOnExecError(stmt.Exec(e.ID, p.ID))
	}
	panicOnError(stmt.Close())
}

// FetchAttendanceByPerson retrieves the attendance for a specific person.
func (tx *Tx) FetchAttendanceByPerson(p *model.Person) (attend map[model.EventID]bool) {
	var (
		rows *sql.Rows
		eid  model.EventID
		err  error
	)
	attend = make(map[model.EventID]bool)
	rows, err = tx.tx.Query(`SELECT event FROM attendance WHERE person=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		panicOnError(rows.Scan(&eid))
		attend[eid] = true
	}
	panicOnError(rows.Err())
	return attend
}
