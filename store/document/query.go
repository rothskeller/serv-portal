package document

import (
	"fmt"
	"os"

	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

const withIDSQL = `SELECT ` + ColumnList + ` FROM document d WHERE d.id=?`

// WithID returns the document with the specified ID, or nil if there is none.
// Note that the resulting document may be archived.
func WithID(storer phys.Storer, id ID) (d *Document) {
	phys.SQL(storer, withIDSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			d = new(Document)
			d.Scan(stmt)
		}
	})
	return d
}

const withNameSQL = `SELECT ` + ColumnList + ` FROM document d WHERE folder=? AND name=? AND NOT archived`

// WithName returns the (non-archived) document with the specified name in the
// specified folder, or nil if there is none.
func WithName(storer phys.Storer, fid folder.ID, name string) (d *Document) {
	phys.SQL(storer, withNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(fid))
		stmt.BindText(name)
		if stmt.Step() {
			d = new(Document)
			d.Scan(stmt)
		}
	})
	return d
}

// Open returns an open file handle to the specified document.  It must be
// closed by the caller.
func Open(did ID) (fh *os.File) {
	var fname = fmt.Sprintf("documents/%02d/%02d", did/100, did%100)
	var err error
	if fh, err = os.OpenFile(fname, os.O_RDONLY, 0); err != nil {
		panic(fmt.Sprintf("document file %d not found in file system: %s", did, err))
	}
	return fh
}

const allInFolderSQL = `SELECT ` + ColumnList + ` FROM document d WHERE folder=? AND NOT archived ORDER BY name`

// AllInFolder fetches each of the documents in the specified folder, in order
// by name.  Archived documents are not returned.
func AllInFolder(storer phys.Storer, fid folder.ID, fn func(*Document)) {
	var d Document

	phys.SQL(storer, allInFolderSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(fid))
		for stmt.Step() {
			d.Scan(stmt)
			fn(&d)
		}
	})
}

const existInFolderSQL = `SELECT 1 FROM document WHERE folder=? LIMIT 1`

// ExistInFolder returns whether any documents, including archived documents,
// exist in the specified folder.
func ExistInFolder(storer phys.Storer, fid folder.ID) (found bool) {
	phys.SQL(storer, existInFolderSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(fid))
		found = stmt.Step()
	})
	return found
}
