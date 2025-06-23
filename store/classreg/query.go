package classreg

import (
	"strings"

	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ClassHasSignups returns whether there are any registrations for the class
// (even if on the waitlist).
func ClassHasSignups(storer phys.Storer, cid class.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM classreg WHERE class=? LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(cid))
		found = stmt.Step()
	})
	return found
}

// ClassHasWaitlist returns whether there is anyone on the waitlist for the
// class.
func ClassHasWaitlist(storer phys.Storer, cid class.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM classreg WHERE class=? AND waitlist LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(cid))
		found = stmt.Step()
	})
	return found
}

// ClassIsFull returns whether the number of non-waitlist registrations for the
// class is >= its limit.
func ClassIsFull(storer phys.Storer, cid class.ID) bool {
	var count, limit int
	phys.SQL(storer, "SELECT COUNT(*) FROM classreg WHERE class=? AND NOT waitlist", func(stmt *phys.Stmt) {
		stmt.BindInt(int(cid))
		stmt.Step()
		count = stmt.ColumnInt()
	})
	if count == 0 {
		return false
	}
	phys.SQL(storer, "SELECT elimit FROM class WHERE id=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(cid))
		stmt.Step()
		limit = stmt.ColumnInt()
	})
	return limit != 0 && count >= limit
}

var withIDSQLCache map[Fields]string

// WithID returns the class registration with the specified ID, or nil if it
// does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (cr *ClassReg) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM classreg cr WHERE cr.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			cr = new(ClassReg)
			cr.Scan(stmt, fields)
			cr.id = id
			cr.fields |= FID
		}
	})
	return cr
}

var allForClassSQLCache map[Fields]string

// AllForClass reads the list of people registered for the specified class, in
// the order that they were registered.  It includes the waitlisted students.
func AllForClass(storer phys.Storer, cid class.ID, fields Fields, fn func(*ClassReg)) {
	if allForClassSQLCache == nil {
		allForClassSQLCache = make(map[Fields]string)
	}
	if _, ok := allForClassSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM classreg cr WHERE cr.class=? ORDER BY cr.id")
		allForClassSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allForClassSQLCache[fields], func(stmt *phys.Stmt) {
		var cr ClassReg

		stmt.BindInt(int(cid))
		for stmt.Step() {
			cr.Scan(stmt, fields)
			fn(&cr)
		}
	})
}
