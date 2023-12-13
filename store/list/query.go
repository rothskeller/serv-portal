package list

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
)

const withIDSQL = `SELECT type, name FROM list WHERE id=?`

// WithID returns the list with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID) (l *List) {
	phys.SQL(storer, withIDSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			l = new(List)
			l.ID = id
			l.Type = Type(stmt.ColumnInt())
			l.Name = stmt.ColumnText()
		}
	})
	return l
}

const allSQL = `SELECT id, type, name FROM list ORDER BY name`

// All reads each list from the database, in order by name.
func All(storer phys.Storer, fn func(*List)) {
	phys.SQL(storer, allSQL, func(stmt *phys.Stmt) {
		var l List
		for stmt.Step() {
			l.ID = ID(stmt.ColumnInt())
			l.Type = Type(stmt.ColumnInt())
			l.Name = stmt.ColumnText()
			fn(&l)
		}
	})
}
