package classreg

import (
	"strings"

	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ClassIsFull returns whether the number of registrations for the class is
// >= its limit.
func ClassIsFull(storer phys.Storer, cid class.ID) bool {
	var count, limit int
	phys.SQL(storer, "SELECT COUNT(*) FROM classreg WHERE class=?", func(stmt *phys.Stmt) {
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

// CountForClass returns the number of people registered for the specified class.
func CountForClass(storer phys.Storer, cid class.ID) (count uint) {
	phys.SQL(storer, "SELECT COUNT(*) FROM classreg WHERE class=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(cid))
		stmt.Step()
		count = uint(stmt.ColumnInt())
	})
	return count
}

var allForClassSQLCache map[Fields]string

// AllForClass reads the list of people registered for the specified class, in
// the order that they were registered.
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
		for stmt.Step() {
			cr.Scan(stmt, fields)
			fn(&cr)
		}
	})
}
