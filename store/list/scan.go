package list

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList is a comma-separated list of column names for the specified
// list fields.  It is used in constructing SQL SELECT statements.
const ColumnList = `l.id, l.type, l.name, l.moderators`

// Scan reads columns from the specified statement into the receiver.
func (l *List) Scan(stmt *phys.Stmt) {
	l.ID = ID(stmt.ColumnInt())
	l.Type = Type(stmt.ColumnInt())
	l.Name = stmt.ColumnText()
	l.Moderators = unpackModerators(stmt.ColumnText())
}
