package folder

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FName | FParent | FURLName | FViewer | FEditor

// Updater is a structure that can be filled with data for a new or changed
// folder, and then later applied.  For creating new folders, it can simply be
// instantiated with new().  For updating existing folders, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the folder being changed.
type Updater struct {
	ID       ID
	Parent   *Folder
	Name     string
	URLName  string
	ViewOrg  enum.Org
	ViewPriv enum.PrivLevel
	EditOrg  enum.Org
	EditPriv enum.PrivLevel
}

// Updater returns a new Updater for the specified folder, with its data
// matching the current data for the folder.  The folder must have fetched
// UpdaterFields.  If parent is non-nil, it is used as the parent folder object,
// avoiding a fetch.
func (f *Folder) Updater(storer phys.Storer, parent *Folder) *Updater {
	if f.fields&UpdaterFields != UpdaterFields {
		panic("Folder.Updater called without fetching UpdaterFields")
	}
	if parent == nil || parent.ID() != f.parent {
		parent = WithID(storer, f.parent, FID|FName)
	}
	return &Updater{
		ID:       f.id,
		Parent:   parent,
		Name:     f.name,
		URLName:  f.urlName,
		ViewOrg:  f.viewOrg,
		ViewPriv: f.viewPriv,
		EditOrg:  f.editOrg,
		EditPriv: f.editPriv,
	}
}

const createSQL = `INSERT INTO folder (id, parent, name, url_name, view_org, view_priv, edit_org, edit_priv) VALUES (?,?,?,?,?,?,?,?)`

// Create creates a new folder, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (f *Folder) {
	f = new(Folder)
	f.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			f.id = u.ID
		} else {
			f.id = ID(phys.LastInsertRowID(storer))
		}
	})
	f.auditAndUpdate(storer, u, true)
	phys.Index(storer, f)
	return f
}

const updateSQL = `UPDATE folder SET parent=?, name=?, url_name=?, view_org=?, view_priv=?, edit_org=?, edit_priv=? WHERE id=?`

// Update updates the existing folder, with the data in the Updater.
func (f *Folder) Update(storer phys.Storer, u *Updater) {
	if f.fields&UpdaterFields != UpdaterFields {
		panic("Folder.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(f.id))
		stmt.Step()
	})
	f.auditAndUpdate(storer, u, false)
	phys.Index(storer, f)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Parent.ID()))
	stmt.BindText(u.Name)
	stmt.BindText(u.URLName)
	stmt.BindInt(int(u.ViewOrg))
	stmt.BindInt(int(u.ViewPriv))
	stmt.BindInt(int(u.EditOrg))
	stmt.BindInt(int(u.EditPriv))
}

func (f *Folder) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Folder %q [%d]", u.Name, f.id)
	if create {
		context = "ADD " + context
	}
	if u.Parent.ID() != f.parent {
		phys.Audit(storer, "%s:: parent = %q [%d]", context, u.Parent.Name(), u.Parent.ID())
		f.parent = u.Parent.ID()
	}
	if u.Name != f.name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		f.name = u.Name
	}
	if u.URLName != f.urlName {
		phys.Audit(storer, "%s:: urlName = %q", context, u.URLName)
		f.urlName = u.URLName
	}
	if u.ViewOrg != f.viewOrg || u.ViewPriv != f.viewPriv {
		phys.Audit(storer, "%s:: viewer = %s [%d] / %s [%d]", context, u.ViewOrg, u.ViewOrg, u.ViewPriv, u.ViewPriv)
		f.viewOrg, f.viewPriv = u.ViewOrg, u.ViewPriv
	}
	if u.EditOrg != f.editOrg || u.EditPriv != f.editPriv {
		phys.Audit(storer, "%s:: editor = %s [%d] / %s [%d]", context, u.EditOrg, u.EditOrg, u.EditPriv, u.EditPriv)
		f.editOrg, f.editPriv = u.EditOrg, u.EditPriv
	}
}

const duplicateURLNameSQL1 = `SELECT 1 FROM folder WHERE id!=? AND parent=? AND url_name=?`
const duplicateURLNameSQL2 = `SELECT 1 FROM document WHERE folder=? AND name=? AND NOT archived`

// DuplicateURLName returns whether the URL name specified in the Updater would
// be a duplicate if applied.
func (u *Updater) DuplicateURLName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateURLNameSQL1, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.Parent.ID()))
		stmt.BindText(u.URLName)
		found = stmt.Step()
	})
	if found {
		return true
	}
	phys.SQL(storer, duplicateURLNameSQL2, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.URLName)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver folder.
func (f *Folder) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM folder WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(f.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE Folder %q [%d]", f.Name(), f.ID())
	phys.Unindex(storer, f)
}
