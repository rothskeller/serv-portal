package redirect

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
)

const withIDSQL = `SELECT entry, target FROM redirect WHERE id=?`

// WithID returns the redirect with the specified ID, or nil if it does not
// exist.
func WithID(storer phys.Storer, id ID) (r *Redirect) {
	phys.SQL(storer, withIDSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			r = new(Redirect)
			r.ID = id
			r.Entry = stmt.ColumnText()
			r.Target = stmt.ColumnText()
		}
	})
	return r
}

const withEntrySQL = `SELECT id, target FROM redirect WHERE entry=?`

// WithEntry returns the redirect with the specified entry URL, or nil if it
// does not exist.
func WithEntry(storer phys.Storer, entry string) (r *Redirect) {
	phys.SQL(storer, withEntrySQL, func(stmt *phys.Stmt) {
		stmt.BindText(entry)
		if stmt.Step() {
			r = new(Redirect)
			r.Entry = entry
			r.ID = ID(stmt.ColumnInt())
			r.Target = stmt.ColumnText()
		}
	})
	return r
}

const allSQL = `SELECT id, entry, target FROM redirect ORDER BY entry`

// All reads each redirect from the database, in order by entry URL.
func All(storer phys.Storer, fn func(*Redirect)) {
	phys.SQL(storer, allSQL, func(stmt *phys.Stmt) {
		var r Redirect
		for stmt.Step() {
			r.ID = ID(stmt.ColumnInt())
			r.Entry = stmt.ColumnText()
			r.Target = stmt.ColumnText()
			fn(&r)
		}
	})
}
