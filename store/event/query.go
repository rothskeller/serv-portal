package event

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/venue"
)

// Exists returns whether an event with the specified ID exists.
func Exists(storer phys.Storer, id ID) (found bool) {
	phys.SQL(storer, `SELECT 1 FROM event WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		found = stmt.Step()
	})
	return found
}

var withIDSQLCache map[Fields]string

// WithID returns the event with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (e *Event) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM event e WHERE e.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			e = new(Event)
			e.Scan(stmt, fields)
			e.id = id
			e.fields |= FID
		}
	})
	return e
}

var allBetweenSQLCache map[Fields]string

// AllBetween reads each event between the specified dates from the database, in
// chronological order.  The date range is inclusive start, exclusive end.
func AllBetween(storer phys.Storer, start, end string, fields Fields, venueFields venue.Fields, fn func(*Event, *venue.Venue)) {
	var vcache = make(map[venue.ID]*venue.Venue)

	if venueFields != 0 {
		fields |= FVenue
	}
	if allBetweenSQLCache == nil {
		allBetweenSQLCache = make(map[Fields]string)
	}
	if _, ok := allBetweenSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM event e WHERE e.start>=? AND e.start<? ORDER BY e.start, e.end, e.id")
		allBetweenSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allBetweenSQLCache[fields], func(stmt *phys.Stmt) {
		var e Event
		stmt.BindText(start)
		stmt.BindText(end)
		for stmt.Step() {
			e.Scan(stmt, fields)
			if venueFields != 0 && e.venue != 0 {
				if vcache[e.venue] == nil {
					vcache[e.venue] = venue.WithID(storer, e.venue, venueFields)
				}
			}
			fn(&e, vcache[e.venue])
		}
	})
}
