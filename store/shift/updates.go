package shift

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FTask | FStart | FEnd | FVenue | FMin | FMax

// Updater is a structure that can be filled with data for a new or changed
// Shift, and then later applied.  For creating new Shifts, it can simply be
// instantiated with new().  For updating existing Shifts, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the Shift being changed.
type Updater struct {
	ID    ID
	Event *event.Event
	Task  *task.Task
	Start string
	End   string
	Venue *venue.Venue
	Min   uint
	Max   uint
}

// Updater returns a new Updater for the receiver Shift, with its data matching
// the current data for the Shift.  The Shift must have fetched UpdaterFields.
// The parent Event and Task and any associated Venue *may* be given as
// arguments to save looking them up.
func (s *Shift) Updater(storer phys.Storer, e *event.Event, t *task.Task, v *venue.Venue) *Updater {
	const eventFields = event.FID | event.FStart | event.FName
	const taskFields = task.FID | task.FEvent | task.FName
	const venueFields = venue.FID | venue.FName

	if s.fields&UpdaterFields != UpdaterFields {
		panic("Shift.Updater called without fetching UpdaterFields")
	}
	if t == nil {
		t = task.WithID(storer, s.task, taskFields)
	}
	if e == nil {
		e = event.WithID(storer, t.Event(), eventFields)
	}
	if s.venue == 0 {
		v = nil
	} else if v == nil {
		v = venue.WithID(storer, s.venue, venueFields)
	}
	return &Updater{
		ID:    s.id,
		Event: e,
		Task:  t,
		Start: s.start,
		End:   s.end,
		Venue: v,
		Min:   s.min,
		Max:   s.max,
	}
}

const createSQL = `INSERT INTO shift (id, task, start, end, venue, min, max) VALUES (?,?,?,?,?,?,?)`

// Create creates a new Shift, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (s *Shift) {
	s = new(Shift)
	s.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			s.id = u.ID
		} else {
			s.id = ID(phys.LastInsertRowID(storer))
		}
	})
	s.auditAndUpdate(storer, u, true)
	return s
}

const updateSQL = `UPDATE shift SET task=?, start=?, end=?, venue=?, min=?, max=? WHERE id=?`

// Update updates the existing Shift, with the data in the Updater.
func (s *Shift) Update(storer phys.Storer, u *Updater) {
	if s.fields&UpdaterFields != UpdaterFields {
		panic("Shift.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(s.id))
		stmt.Step()
	})
	s.auditAndUpdate(storer, u, false)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Task.ID()))
	stmt.BindText(u.Start)
	stmt.BindText(u.End)
	if u.Venue != nil {
		stmt.BindInt(int(u.Venue.ID()))
	} else {
		stmt.BindNull()
	}
	stmt.BindInt(int(u.Min))
	stmt.BindNullInt(int(u.Max))
}

func (s *Shift) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Event %s %q [%d]:: Task %q [%d]:: Shift %d", u.Event.Start()[:10], u.Event.Name(), u.Event.ID(), u.Task.Name(), u.Task.ID(), s.id)
	if create {
		context = "ADD " + context
	}
	if u.Task.ID() != s.task {
		phys.Audit(storer, "%s:: task = Event %s %q [%d] Task %q [%d]", context, u.Event.Start()[:10], u.Event.Name(), u.Event.ID(), u.Task.Name(), u.Task.ID())
		s.task = u.Task.ID()
	}
	if u.Start != s.start {
		phys.Audit(storer, "%s:: start = %s", context, u.Start)
		s.start = u.Start
	}
	if u.End != s.end {
		phys.Audit(storer, "%s:: end = %s", context, u.End)
		s.end = u.End
	}
	if u.Venue != nil {
		if u.Venue.ID() != s.venue {
			phys.Audit(storer, "%s:: venue = %q [%d]", context, u.Venue.Name(), u.Venue.ID())
			s.venue = u.Venue.ID()
		}
	} else if s.venue != 0 {
		phys.Audit(storer, "%s:: venue = nil", context)
		s.venue = 0
	}
	if u.Min != s.min {
		phys.Audit(storer, "%s:: min = %d", context, u.Min)
		s.min = u.Min
	}
	if u.Max != s.max {
		phys.Audit(storer, "%s:: max = %d", context, u.Max)
		s.max = u.Max
	}
}

const overlappingShiftSQL = `SELECT 1 FROM shift WHERE id!=? AND task=? AND venue IS ? AND ?<end AND ?>start`

// OverlappingShift returns whether the data specified in the Updater would
// overlap another Shift if applied.
func (u *Updater) OverlappingShift(storer phys.Storer) (found bool) {
	phys.SQL(storer, overlappingShiftSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.Task.ID()))
		if u.Venue == nil {
			stmt.BindNull()
		} else {
			stmt.BindInt(int(u.Venue.ID()))
		}
		stmt.BindText(u.Start)
		stmt.BindText(u.End)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver Shift.  The parent Event and Task *may* be
// specified to avoid a lookup.
func (s *Shift) Delete(storer phys.Storer, e *event.Event, t *task.Task) {
	const eventFields = event.FID | event.FStart | event.FName
	const taskFields = task.FID | task.FName
	if t == nil || t.Fields()&taskFields != taskFields || t.ID() != s.task {
		t = task.WithID(storer, s.task, taskFields)
	}
	if e == nil || e.Fields()&eventFields != eventFields || e.ID() != t.Event() {
		e = event.WithID(storer, t.Event(), eventFields)
	}
	phys.SQL(storer, `DELETE FROM shift WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "Event %s %q [%d]:: Task %q [%d]:: DELETE Shift %d", e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), s.ID())
}
