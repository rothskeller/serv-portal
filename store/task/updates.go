package task

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FEvent | FName | FOrg | FFlags | FDetails

// Updater is a structure that can be filled with data for a new or changed
// Task, and then later applied.  For creating new Tasks, it can simply be
// instantiated with new().  For updating existing Tasks, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the Task being changed.
type Updater struct {
	ID      ID
	Event   *event.Event
	Name    string
	Org     enum.Org
	Flags   Flag
	Details string
}

// Updater returns a new Updater for the receiver Task, with its data matching
// the current data for the Task.  The Task must have fetched UpdaterFields.
// The parent Event *may* be given as an argument to save looking it up.
func (t *Task) Updater(storer phys.Storer, e *event.Event) *Updater {
	const eventFields = event.FID | event.FStart | event.FName

	if t.fields&UpdaterFields != UpdaterFields {
		panic("Task.Updater called without fetching UpdaterFields")
	}
	if e == nil {
		e = event.WithID(storer, t.event, eventFields)
	}
	return &Updater{
		ID:      t.id,
		Event:   e,
		Name:    t.name,
		Org:     t.org,
		Flags:   t.flags,
		Details: t.details,
	}
}

const nextSortSQL = `SELECT COALESCE(MAX(sort), 0) FROM task WHERE event=?`
const createSQL = `INSERT INTO task (id, sort, event, name, org, flags, details) VALUES (?,?,?,?,?,?,?)`

// Create creates a new Task, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (t *Task) {
	var sort int

	t = new(Task)
	t.fields = UpdaterFields
	phys.SQL(storer, nextSortSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.Event.ID()))
		stmt.Step()
		sort = stmt.ColumnInt() + 1
	})
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		stmt.BindInt(sort)
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			t.id = u.ID
		} else {
			t.id = ID(phys.LastInsertRowID(storer))
		}
	})
	t.auditAndUpdate(storer, u, true)
	return t
}

const updateSQL = `UPDATE task SET event=?, name=?, org=?, flags=?, details=? WHERE id=?`

// Update updates the existing Task, with the data in the Updater.
func (t *Task) Update(storer phys.Storer, u *Updater) {
	if t.fields&UpdaterFields != UpdaterFields {
		panic("Task.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(t.id))
		stmt.Step()
	})
	t.auditAndUpdate(storer, u, false)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Event.ID()))
	stmt.BindNullText(u.Name)
	stmt.BindInt(int(u.Org))
	stmt.BindHexInt(int(u.Flags))
	stmt.BindNullText(u.Details)
}

func (t *Task) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Event %s %q [%d]:: Task %q [%d]", u.Event.Start()[:10], u.Event.Name(), u.Event.ID(), u.Name, t.id)
	if create {
		context = "ADD " + context
	}
	if u.Event.ID() != t.event {
		phys.Audit(storer, "%s:: event = %s %q [%d]", context, u.Event.Start()[:10], u.Event.Name(), u.Event.ID())
		t.event = u.Event.ID()
	}
	if u.Name != t.name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		t.name = u.Name
	}
	if u.Org != t.org {
		phys.Audit(storer, "%s:: org = %s [%d]", context, u.Org, u.Org)
		t.org = u.Org
	}
	if u.Flags != t.flags {
		phys.Audit(storer, "%s:: flags = 0x%x", context, u.Flags)
		t.flags = u.Flags
	}
	if u.Details != t.details {
		phys.Audit(storer, "%s:: details = %q", context, u.Details)
		t.details = u.Details
	}
}

const duplicateNameSQL = `SELECT 1 FROM task WHERE id!=? AND event=? AND name=?`

// DuplicateName returns whether the name specified in the Updater would be a
// duplicate if applied.
func (u *Updater) DuplicateName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.Event.ID()))
		stmt.BindText(u.Name)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver Task.  The parent Event *may* be specified to
// avoid a lookup.
func (t *Task) Delete(storer phys.Storer, e *event.Event) {
	const eventFields = event.FID | event.FStart | event.FName
	if e == nil || e.Fields()&eventFields != eventFields || e.ID() != t.event {
		e = event.WithID(storer, t.event, eventFields)
	}
	phys.SQL(storer, `DELETE FROM task WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "Event %s %q [%d]:: DELETE Task %q [%d]", e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID())
}
