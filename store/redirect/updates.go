package redirect

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// Updater is a structure that can be filled with data for a new or changed
// redirect, and then later applied.  For creating new redirects, it can simply
// be instantiated with new().  For updating existing redirects, either *every*
// field in it must be set, or it should be instantiated with the Updater method
// of the redirect being changed.
type Updater Redirect

// Updater returns a new Updater for the specified redirect, with its data
// matching the current data for the redirect.  The redirect must have fetched
// UpdaterFields.
func (r *Redirect) Updater() (u *Updater) {
	u = &Updater{
		ID:     r.ID,
		Entry:  r.Entry,
		Target: r.Target,
	}
	return u
}

const createSQL = `INSERT INTO redirect (id, entry, target) VALUES (?,?,?)`

// Create creates a new redirect, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (r *Redirect) {
	r = new(Redirect)
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		stmt.BindText(u.Entry)
		stmt.BindText(u.Target)
		stmt.Step()
		if u.ID != 0 {
			r.ID = u.ID
		} else {
			r.ID = ID(phys.LastInsertRowID(storer))
		}
	})
	r.auditAndUpdate(storer, u, true)
	return r
}

const updateSQL = `UPDATE redirect SET entry=?, target=? WHERE id=?`

// Update updates the existing redirect, with the data in the Updater.
func (r *Redirect) Update(storer phys.Storer, u *Updater) {
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		stmt.BindText(u.Entry)
		stmt.BindText(u.Target)
		stmt.BindInt(int(u.ID))
		stmt.Step()
	})
	r.auditAndUpdate(storer, u, false)
}

func (r *Redirect) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Redirect %q [%d]", u.Entry, r.ID)
	if create {
		context = "ADD " + context
	}
	if u.Entry != r.Entry {
		phys.Audit(storer, "%s:: entry = %q", context, u.Entry)
		r.Entry = u.Entry
	}
	if u.Target != r.Target {
		phys.Audit(storer, "%s:: target = %q", context, u.Target)
		r.Target = u.Target
	}
}

const duplicateEntrySQL = `SELECT 1 FROM redirect WHERE id!=? AND entry=?`

// DuplicateEntry returns whether the entry URL specified in the Updater would
// be a duplicate if applied.
func (u *Updater) DuplicateEntry(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateEntrySQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Entry)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver redirect.
func (r *Redirect) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM redirect WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(r.ID))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE Redirect %q [%d]", r.Entry, r.ID)
}
