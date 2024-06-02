package list

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// Updater is a structure that can be filled with data for a new or changed
// list, and then later applied.  For creating new lists, it can simply be
// instantiated with new().  For updating existing lists, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the list being changed.
type Updater List

// Updater returns a new Updater for the specified list, with its data matching
// the current data for the list.  The list must have fetched UpdaterFields.
func (l *List) Updater() (u *Updater) {
	u = &Updater{
		ID:   l.ID,
		Type: l.Type,
		Name: l.Name,
	}
	if l.Moderators != nil {
		u.Moderators = l.Moderators.Clone()
	}
	return u
}

const createSQL = `INSERT INTO list (id, type, name, moderators) VALUES (?,?,?,?)`

// Create creates a new list, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (l *List) {
	l = new(List)
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		stmt.BindInt(int(u.Type))
		stmt.BindText(u.Name)
		stmt.BindNullText(packModerators(u.Moderators))
		stmt.Step()
		if u.ID != 0 {
			l.ID = u.ID
		} else {
			l.ID = ID(phys.LastInsertRowID(storer))
		}
	})
	l.auditAndUpdate(storer, u, true)
	return l
}

const updateSQL = `UPDATE list SET type=?, name=?, moderators=? WHERE id=?`

// Update updates the existing list, with the data in the Updater.
func (l *List) Update(storer phys.Storer, u *Updater) {
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.Type))
		stmt.BindText(u.Name)
		stmt.BindNullText(packModerators(u.Moderators))
		stmt.BindInt(int(u.ID))
		stmt.Step()
	})
	l.auditAndUpdate(storer, u, false)
}

func packModerators(mods sets.Set[string]) string {
	if !mods.HasAny() {
		return ""
	}
	return strings.Join(mods.UnsortedList(), ",")
}

func (l *List) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("List %q [%d]", u.Name, l.ID)
	if create {
		context = "ADD " + context
	}
	if u.Type != l.Type {
		phys.Audit(storer, "%s:: type = %s [%d]", context, u.Type, u.Type)
		l.Type = u.Type
	}
	if u.Name != l.Name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		l.Name = u.Name
	}
	if u.Moderators != nil && (l.Moderators == nil || !u.Moderators.Equal(l.Moderators)) {
		phys.Audit(storer, "%s:: moderators = %q", context, packModerators(u.Moderators))
	} else if u.Moderators == nil && l.Moderators != nil {
		phys.Audit(storer, "%s:: moderators = nil", context)
	}
}

const duplicateNameSQL = `SELECT 1 FROM list WHERE id!=? AND name=?`

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

// Delete deletes the receiver list.
func (l *List) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM list WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(l.ID))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE List %q [%d]", l.Name, l.ID)
}
