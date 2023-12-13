package venue

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// Exists returns whether a venue with the specified ID exists.
func Exists(storer phys.Storer, id ID) (found bool) {
	phys.SQL(storer, `SELECT 1 FROM venue WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		found = stmt.Step()
	})
	return found
}

var withIDSQLCache map[Fields]string

// WithID returns the venue with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (v *Venue) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM venue v WHERE v.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			v = new(Venue)
			v.Scan(stmt, fields)
			v.id = id
			v.fields |= FID
		}
	})
	return v
}

var allSQLCache map[Fields]string

// All reads each venue from the database, in order by name.
func All(storer phys.Storer, fields Fields, fn func(*Venue)) {
	if allSQLCache == nil {
		allSQLCache = make(map[Fields]string)
	}
	if _, ok := allSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM venue v ORDER BY v.name")
		allSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields], func(stmt *phys.Stmt) {
		var v Venue
		for stmt.Step() {
			v.Scan(stmt, fields)
			fn(&v)
		}
	})
}
