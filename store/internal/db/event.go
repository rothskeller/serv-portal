package db

import (
	"database/sql"
	"sort"

	"sunnyvaleserv.org/portal/model"
)

// FetchEvent retrieves a single event from the database by ID.  It returns nil
// if no such event exists.
func (tx *Tx) FetchEvent(id model.EventID) (e *model.Event) {
	var data []byte
	e = new(model.Event)
	switch err := tx.tx.QueryRow(`SELECT data FROM event WHERE id=?`, id).Scan(&data); err {
	case nil:
		panicOnError(e.Unmarshal(data))
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchEvents returns all of the events within the specified time range, in
// chronological order.  The time range is inclusive; each time must be in
// 2006-01-02 format.
func (tx *Tx) FetchEvents(from, to string) (events []*model.Event) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM event WHERE date>=? AND date<=?`, from, to)
	panicOnError(err)
	for rows.Next() {
		var data []byte
		var e model.Event
		panicOnError(rows.Scan(&data))
		panicOnError(e.Unmarshal(data))
		events = append(events, &e)
	}
	panicOnError(rows.Err())
	sort.Sort(model.EventSort(events))
	return events
}

// CreateEvent creates a new event in the database, assigning it the next
// available ID.
func (tx *Tx) CreateEvent(e *model.Event) {
	var (
		data []byte
		err  error
	)
	panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id), 0) FROM event`).Scan(&e.ID))
	e.ID++
	data, err = e.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`INSERT INTO event (id, date, data) VALUES (?,?,?)`, e.ID, e.Date, data))
	tx.indexEvent(e, false)
}

// UpdateEvent updates an existing event in the database.
func (tx *Tx) UpdateEvent(e *model.Event) {
	var (
		data []byte
		err  error
	)
	data, err = e.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE event SET (date, data) = (?,?) WHERE id=?`, e.Date, data, e.ID))
	tx.indexEvent(e, true)
}

func (tx *Tx) indexEvent(e *model.Event, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='event' and id=?`, e.ID))
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, eventName, eventDetails, eventDate) VALUES ('event',?,?,?,?)`, e.ID, e.Name, IDStr(e.Details), e.Date))
}

// DeleteEvent deletes an event from the database.
func (tx *Tx) DeleteEvent(e *model.Event) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
}

// FetchAttendanceByEvent retrieves the attendance at a specific event.
func (tx *Tx) FetchAttendanceByEvent(e *model.Event) (attend map[model.PersonID]model.AttendanceInfo) {
	var (
		rows *sql.Rows
		err  error
	)
	attend = make(map[model.PersonID]model.AttendanceInfo)
	rows, err = tx.tx.Query(`SELECT person, type, minutes FROM attendance WHERE event=?`, e.ID)
	panicOnError(err)
	for rows.Next() {
		var ai model.AttendanceInfo
		var pid model.PersonID
		panicOnError(rows.Scan(&pid, &ai.Type, &ai.Minutes))
		attend[pid] = ai
	}
	panicOnError(rows.Err())
	return attend
}

// SaveEventAttendance saves the attendance for a specific event.
func (tx *Tx) SaveEventAttendance(e *model.Event, attend map[model.PersonID]model.AttendanceInfo) {
	var (
		stmt *sql.Stmt
		err  error
	)
	panicOnExecError(tx.tx.Exec(`DELETE FROM attendance WHERE event=?`, e.ID))
	stmt, err = tx.tx.Prepare(`INSERT INTO attendance (event, person, type, minutes) VALUES (?,?,?,?)`)
	panicOnError(err)
	for pid, att := range attend {
		panicOnExecError(stmt.Exec(e.ID, pid, att.Type, att.Minutes))

	}
	panicOnError(stmt.Close())
}

// FetchAttendanceByPerson retrieves the attendance for a specific person.
func (tx *Tx) FetchAttendanceByPerson(p *model.Person) (attend map[model.EventID]model.AttendanceInfo) {
	var (
		rows *sql.Rows
		err  error
	)
	attend = make(map[model.EventID]model.AttendanceInfo)
	rows, err = tx.tx.Query(`SELECT event, type, minutes FROM attendance WHERE person=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		var ai model.AttendanceInfo
		var eid model.EventID
		panicOnError(rows.Scan(&eid, &ai.Type, &ai.Minutes))
		attend[eid] = ai
	}
	panicOnError(rows.Err())
	return attend
}
