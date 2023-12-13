package venue

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FName | FURL | FFlags

// Updater is a structure that can be filled with data for a new or changed
// venue, and then later applied.  For creating new venues, it can simply be
// instantiated with new().  For updating existing venues, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the venue being changed.
type Updater struct {
	ID    ID
	Name  string
	URL   string
	Flags Flag
}

// Updater returns a new Updater for the specified venue, with its data matching
// the current data for the venue.  The venue must have fetched UpdaterFields.
func (v *Venue) Updater() *Updater {
	if v.fields&UpdaterFields != UpdaterFields {
		panic("Venue.Updater called without fetching UpdaterFields")
	}
	return &Updater{
		ID:    v.id,
		Name:  v.name,
		URL:   v.url,
		Flags: v.flags,
	}
}

const createSQL = `INSERT INTO venue (id, name, url, flags) VALUES (?,?,?,?)`

// Create creates a new venue, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (v *Venue) {
	v = new(Venue)
	v.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			v.id = u.ID
		} else {
			v.id = ID(phys.LastInsertRowID(storer))
		}
	})
	v.auditAndUpdate(storer, u, true)
	phys.Index(storer, v)
	return v
}

const updateSQL = `UPDATE venue SET name=?, url=?, flags=? WHERE id=?`

// Update updates the existing venue, with the data in the Updater.
func (v *Venue) Update(storer phys.Storer, u *Updater) {
	if v.fields&UpdaterFields != UpdaterFields {
		panic("Venue.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(v.id))
		stmt.Step()
	})
	v.auditAndUpdate(storer, u, false)
	phys.Index(storer, v)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindText(u.Name)
	stmt.BindNullText(u.URL)
	stmt.BindHexInt(int(u.Flags))
}

func (v *Venue) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Venue %q [%d]", u.Name, v.id)
	if create {
		context = "ADD " + context
	}
	if u.Name != v.name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		v.name = u.Name
	}
	if u.URL != v.url {
		phys.Audit(storer, "%s:: url = %q", context, u.URL)
		v.url = u.URL
	}
	if u.Flags != v.flags {
		phys.Audit(storer, "%s:: flags = 0x%x", context, u.Flags)
		v.flags = u.Flags
	}
}

const duplicateNameSQL = `SELECT 1 FROM venue WHERE id!=? AND name=?`

// DuplicateName returns whether the name specified in the Updater
// would be a duplicate if applied.
func (u *Updater) DuplicateName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Name)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver venue.
func (v *Venue) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM venue WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(v.ID()))
		stmt.Step()
	})
	phys.Unindex(storer, v)
	phys.Audit(storer, "DELETE Venue %q [%d]", v.Name(), v.ID())
}
