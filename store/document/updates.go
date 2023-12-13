package document

import (
	"fmt"
	"os"

	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// Updater is a structure that can be filled with data for a new or changed
// document, and then later applied.  For creating new documents, it can simply
// be instantiated with new().  For updating existing documents, either *every*
// field in it must be set, or it should be instantiated with the Updater method
// of the document being changed.
type Updater struct {
	ID       ID
	Folder   *folder.Folder
	Name     string
	URL      string
	Archived bool
	Contents []byte
	LinkTo   string
}

// Updater returns a new Updater for the specified document, with its data
// matching the current data for the document.  If folder is non-nil, it is used
// as the folder object, avoiding a fetch.
func (d *Document) Updater(storer phys.Storer, f *folder.Folder) *Updater {
	if f == nil || f.ID() != d.Folder {
		f = folder.WithID(storer, d.Folder, folder.FID|folder.FName)
	}
	return &Updater{
		ID:       d.ID,
		Folder:   f,
		Name:     d.Name,
		URL:      d.URL,
		Archived: d.Archived,
	}
}

const createSQL = `INSERT INTO document (id, folder, name, url, archived) VALUES (?,?,?,?,?)`

// Create creates a new document, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (d *Document) {
	d = new(Document)
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			d.ID = u.ID
		} else {
			d.ID = ID(phys.LastInsertRowID(storer))
			u.ID = d.ID
		}
	})
	switch {
	case u.URL != "":
		if u.Contents != nil || u.LinkTo != "" {
			panic("specify exactly one of URL, Contents, or LinkTo")
		}
	case u.Contents != nil:
		if u.LinkTo != "" {
			panic("specify exactly one of URL, Contents, or LinkTo")
		}
		if err := os.MkdirAll(fmt.Sprintf("documents/%02d", d.ID/100), 0777); err != nil {
			panic(err)
		}
		if err := os.WriteFile(fmt.Sprintf("documents/%02d/%02d", d.ID/100, d.ID%100), u.Contents, 0666); err != nil {
			panic(err)
		}
	case u.LinkTo != "":
		if err := os.MkdirAll(fmt.Sprintf("documents/%02d", d.ID/100), 0777); err != nil {
			panic(err)
		}
		if err := os.Link(u.LinkTo, fmt.Sprintf("documents/%02d/%02d", d.ID/100, d.ID%100)); err != nil {
			panic(err)
		}
	default:
		panic("specify exactly one of URL, Contents, or LinkTo")
	}
	d.auditAndUpdate(storer, u, true)
	phys.Index(storer, d)
	return d
}

const updateSQL = `UPDATE document SET folder=?, name=?, url=?, archived=? WHERE id=?`

// Update updates the existing document, with the data in the Updater.
func (d *Document) Update(storer phys.Storer, u *Updater) {
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(d.ID))
		stmt.Step()
	})
	d.auditAndUpdate(storer, u, false)
	phys.Index(storer, d)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Folder.ID()))
	stmt.BindText(u.Name)
	stmt.BindNullText(u.URL)
	stmt.BindBool(u.Archived)
}

func (d *Document) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Document %q [%d]", u.Name, d.ID)
	if create {
		context = "ADD " + context
	}
	if u.Folder.ID() != d.Folder {
		phys.Audit(storer, "%s:: folder = %q [%d]", context, u.Folder.Name(), u.Folder.ID())
		d.Folder = u.Folder.ID()
	}
	if u.Name != d.Name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		d.Name = u.Name
	}
	if u.URL != d.URL {
		phys.Audit(storer, "%s:: url = %q", context, u.URL)
		d.URL = u.URL
	}
	if u.Archived != d.Archived {
		phys.Audit(storer, "%s:: archived = %v", context, u.Archived)
		d.Archived = u.Archived
	}
}

const duplicateNameSQL = `SELECT 1 FROM document WHERE id!=? AND folder=? AND name=? AND NOT archived`

// DuplicateName returns whether the name specified in the Updater would be a
// duplicate if applied.
func (u *Updater) DuplicateName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.Folder.ID()))
		stmt.BindText(u.Name)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver document from the database.  (It does not delete
// the corresponding file, if any, from the file system.)
func (d *Document) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM document WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(d.ID))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE Document %q [%d]", d.Name, d.ID)
	phys.Unindex(storer, d)
}
