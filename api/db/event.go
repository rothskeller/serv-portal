package db

import (
	"database/sql"
	"sort"

	"rothskeller.net/serv/model"
)

// FetchEvent retrieves a single event from the database by ID.  It returns nil
// if no such event exists.
func (tx *Tx) FetchEvent(id model.EventID) (e *model.Event) {
	var data []byte
	e = new(model.Event)
	switch err := dbh.QueryRow(`SELECT data FROM event WHERE id=?`, id).Scan(&data); err {
	case nil:
		panicOnError(e.Unmarshal(data))
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchEventBySccAresID retrieves a single event from the database by its
// scc-ares-races.org ID.  It returns nil if no such event exists.
func (tx *Tx) FetchEventBySccAresID(id string) (e *model.Event) {
	var data []byte
	e = new(model.Event)
	switch err := dbh.QueryRow(`SELECT data FROM event WHERE scc_ares_id=?`, id).Scan(&data); err {
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
	rows, err = dbh.Query(`SELECT data FROM event WHERE date>=? AND date<=?`, from, to)
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

// SaveEvent saves an event to the database.  If the supplied event ID is
// zero, a new event is added to the database; otherwise, the identified event
// is updated.
func (tx *Tx) SaveEvent(e *model.Event) {
	var (
		data []byte
		err  error
	)
	if e.ID == 0 {
		panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id), 0) FROM event`).Scan(&e.ID))
		e.ID++
		data, err = e.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`INSERT INTO event (id, date, scc_ares_id, data) VALUES (?,?,?,?)`, e.ID, e.Date, IDStr(e.SccAresID), data))
	} else {
		data, err = e.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`UPDATE event SET (date, scc_ares_id, data) = (?,?,?) WHERE id=?`, e.Date, IDStr(e.SccAresID), data, e.ID))
	}
	tx.audit("event", e.ID, data)
}

// DeleteEvent deletes an event from the database.
func (tx *Tx) DeleteEvent(e *model.Event) {
	panicOnNoRows(tx.tx.Exec(`DELETE FROM event WHERE id=?`, e.ID))
	tx.audit("event", e.ID, nil)
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
