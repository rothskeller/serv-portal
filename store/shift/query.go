package shift

import (
	"strings"

	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
)

// Exists returns whether a Shift with the specified ID exists.
func Exists(storer phys.Storer, id ID) (found bool) {
	phys.SQL(storer, `SELECT 1 FROM task WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		found = stmt.Step()
	})
	return found
}

var withIDSQLCache map[Fields]string

// WithID returns the Shift with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (s *Shift) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM shift s WHERE s.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			s = new(Shift)
			s.Scan(stmt, fields)
			s.id = id
			s.fields |= FID
		}
	})
	return s
}

// ExistsForTask returns whether the specified Task has any Shifts.
func ExistsForTask(storer phys.Storer, tid task.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM shift WHERE task=? LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(tid))
		found = stmt.Step()
	})
	return found
}

var allForTaskSQLCache map[Fields]map[venue.Fields]string

// AllForTask fetches all Shifts for the specified Task, in order by venue name
// and then by time range.
func AllForTask(storer phys.Storer, tid task.ID, fields Fields, venueFields venue.Fields, fn func(*Shift, *venue.Venue)) {
	if venueFields != 0 {
		fields |= FVenue
	}
	if allForTaskSQLCache == nil {
		allForTaskSQLCache = make(map[Fields]map[venue.Fields]string)
	}
	if allForTaskSQLCache[fields] == nil {
		allForTaskSQLCache[fields] = make(map[venue.Fields]string)
	}
	if _, ok := allForTaskSQLCache[fields][venueFields]; !ok {
		var sb strings.Builder
		sb.WriteString(`SELECT `)
		ColumnList(&sb, fields)
		if venueFields != 0 {
			sb.WriteString(`, `)
		}
		venue.ColumnList(&sb, venueFields)
		sb.WriteString(` FROM shift s LEFT JOIN venue v ON s.venue=v.id WHERE s.task=? ORDER BY v.name, s.start, s.end`)
		allForTaskSQLCache[fields][venueFields] = sb.String()
	}
	phys.SQL(storer, allForTaskSQLCache[fields][venueFields], func(stmt *phys.Stmt) {
		var s Shift
		var v *venue.Venue
		if venueFields != 0 {
			v = new(venue.Venue)
		}
		stmt.BindInt(int(tid))
		for stmt.Step() {
			s.Scan(stmt, fields)
			if venueFields != 0 && s.venue != 0 {
				v.Scan(stmt, venueFields)
				fn(&s, v)
			} else {
				fn(&s, nil)
			}
		}
	})
}

// AllAfter fetches all shifts that start on or after the specified time, along
// with their corresponding events, tasks, and venues if the requisite fields
// are requested.  The shifts are fetched in event, task, and shift order.
func AllAfter(storer phys.Storer, datetime string, eventFields event.Fields, taskFields task.Fields, shiftFields Fields, venueFields venue.Fields, fn func(*event.Event, *task.Task, *Shift, *venue.Venue)) {
	var sb strings.Builder

	sb.WriteString(`SELECT `)
	ColumnList(&sb, shiftFields)
	if eventFields != 0 {
		sb.WriteString(", ")
		event.ColumnList(&sb, eventFields)
	}
	if taskFields != 0 {
		sb.WriteString(", ")
		task.ColumnList(&sb, taskFields)
	}
	if venueFields != 0 {
		venueFields |= venue.FID
		sb.WriteString(", ")
		venue.ColumnList(&sb, venueFields)
	}
	sb.WriteString(" FROM event e, task t, shift s LEFT JOIN venue v ON s.venue=v.id WHERE s.start>=? AND s.task=t.id AND t.event=e.id ORDER BY e.start, e.end, e.id, t.sort, v.name, s.start, s.end")
	phys.SQL(storer, sb.String(), func(stmt *phys.Stmt) {
		var (
			e *event.Event
			t *task.Task
			s Shift
			v *venue.Venue
		)
		if eventFields != 0 {
			e = new(event.Event)
		}
		if taskFields != 0 {
			t = new(task.Task)
		}
		if venueFields != 0 {
			v = new(venue.Venue)
		}
		stmt.BindText(datetime)
		for stmt.Step() {
			s.Scan(stmt, shiftFields)
			if eventFields != 0 {
				e.Scan(stmt, eventFields)
			}
			if taskFields != 0 {
				t.Scan(stmt, taskFields)
			}
			if venueFields != 0 {
				v.Scan(stmt, venueFields)
			}
			if venueFields != 0 && v.ID() != 0 {
				fn(e, t, &s, v)
			} else {
				fn(e, t, &s, nil)
			}
		}
	})
}
