package redirect

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList is a comma-separated list of column names for the specified
// redirect fields.  It is used in constructing SQL SELECT statements.
const ColumnList = `l.id, l.entry, l.target`

// Scan reads columns from the specified statement into the receiver.
func (r *Redirect) Scan(stmt *phys.Stmt) {
	r.ID = ID(stmt.ColumnInt())
	r.Entry = stmt.ColumnText()
	r.Target = stmt.ColumnText()
}
