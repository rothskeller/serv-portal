package document

import (
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList is a comma-separated list of column names for the
// document fields.  It is used in constructing SQL SELECT statements.
const ColumnList = `d.id, d.folder, d.name, d.url, d.archived`

// Scan reads columns corresponding to ColumnList from the specified statement
// into the receiver.
func (d *Document) Scan(stmt *phys.Stmt) {
	d.ID = ID(stmt.ColumnInt())
	d.Folder = folder.ID(stmt.ColumnInt())
	d.Name = stmt.ColumnText()
	d.URL = stmt.ColumnText()
	d.Archived = stmt.ColumnBool()
}
