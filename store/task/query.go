package task

import (
	"strings"

	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// Exists returns whether a Task with the specified ID exists.
func Exists(storer phys.Storer, id ID) (found bool) {
	phys.SQL(storer, `SELECT 1 FROM task WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		found = stmt.Step()
	})
	return found
}

var withIDSQLCache map[Fields]string

// WithID returns the Task with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (t *Task) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM task t WHERE t.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			t = new(Task)
			t.Scan(stmt, fields)
			t.id = id
			t.fields |= FID
		}
	})
	return t
}

// CountForEvent returns the number of Tasks for the specified Event.
func CountForEvent(storer phys.Storer, eid event.ID) (count int) {
	phys.SQL(storer, "SELECT COUNT(*) FROM task WHERE event=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(eid))
		stmt.Step()
		count = stmt.ColumnInt()
	})
	return count
}

var allForEventSQLCache map[Fields]string

// AllForEvent fetches all Tasks for the specified Event, in sorted order.
func AllForEvent(storer phys.Storer, eid event.ID, fields Fields, fn func(*Task)) {
	if allForEventSQLCache == nil {
		allForEventSQLCache = make(map[Fields]string)
	}
	if _, ok := allForEventSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(`SELECT `)
		ColumnList(&sb, fields)
		sb.WriteString(` FROM task t`)
		sb.WriteString(` WHERE t.event=? ORDER BY t.sort`)
		allForEventSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allForEventSQLCache[fields], func(stmt *phys.Stmt) {
		var t Task
		stmt.BindInt(int(eid))
		for stmt.Step() {
			t.Scan(stmt, fields)
			fn(&t)
		}
	})
}
