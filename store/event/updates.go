package event

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/venue"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FName | FStart | FEnd | FVenue | FVenueURL | FActivation | FDetails | FFlags

// Updater is a structure that can be filled with data for a new or changed
// Event, and then later applied.  For creating new events, it can simply be
// instantiated with new().  For updating existing events, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the Event being changed.
type Updater struct {
	ID         ID
	Name       string
	Start      string
	End        string
	Venue      *venue.Venue
	VenueURL   string
	Activation string
	Details    string
	Flags      Flag
}

// Updater returns a new Updater for the receiver Event, with its data matching
// the current data for the event.  The Event must have fetched UpdaterFields.
// The associated Venue *may* be given as an argument to save looking it up.
func (e *Event) Updater(storer phys.Storer, v *venue.Venue) *Updater {
	const venueFields = venue.FID | venue.FName

	if e.fields&UpdaterFields != UpdaterFields {
		panic("Event.Updater called without fetching UpdaterFields")
	}
	if e.venue == 0 {
		v = nil
	} else if v == nil || v.Fields()&venueFields != venueFields || v.ID() != e.venue {
		v = venue.WithID(storer, e.venue, venueFields)
	}
	return &Updater{
		ID:         e.id,
		Name:       e.name,
		Start:      e.start,
		End:        e.end,
		Venue:      v,
		VenueURL:   e.venueURL,
		Activation: e.activation,
		Details:    e.details,
		Flags:      e.flags,
	}
}

const createSQL = `INSERT INTO event (id, name, start, end, venue, venue_url, activation, details, flags) VALUES (?,?,?,?,?,?,?,?,?)`

// Create creates a new Event, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (e *Event) {
	e = new(Event)
	e.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			e.id = u.ID
		} else {
			e.id = ID(phys.LastInsertRowID(storer))
		}
	})
	e.auditAndUpdate(storer, u, true)
	phys.Index(storer, e)
	return e
}

const updateSQL = `UPDATE event SET name=?, start=?, end=?, venue=?, venue_url=?, activation=?, details=?, flags=? WHERE id=?`

// Update updates the existing event, with the data in the Updater.
func (e *Event) Update(storer phys.Storer, u *Updater) {
	if e.fields&UpdaterFields != UpdaterFields {
		panic("Event.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(e.id))
		stmt.Step()
	})
	e.auditAndUpdate(storer, u, false)
	phys.Index(storer, e)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindText(u.Name)
	stmt.BindText(u.Start)
	stmt.BindText(u.End)
	if u.Venue != nil {
		stmt.BindInt(int(u.Venue.ID()))
	} else {
		stmt.BindNull()
	}
	stmt.BindNullText(u.VenueURL)
	stmt.BindNullText(u.Activation)
	stmt.BindNullText(u.Details)
	stmt.BindHexInt(int(u.Flags))
}

func (e *Event) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Event %s %q [%d]", u.Start[:10], u.Name, e.id)
	if create {
		context = "ADD " + context
	}
	if u.Name != e.name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		e.name = u.Name
	}
	if u.Start != e.start {
		phys.Audit(storer, "%s:: start = %s", context, u.Start)
		e.start = u.Start
	}
	if u.End != e.end {
		phys.Audit(storer, "%s:: end = %s", context, u.End)
		e.end = u.End
	}
	if vid := u.Venue.ID(); vid != e.venue {
		if vid == 0 {
			phys.Audit(storer, "%s:: venue = nil", context)
		} else {
			phys.Audit(storer, "%s:: venue = %q [%d]", context, u.Venue.Name(), vid)
		}
		e.venue = vid
	}
	if u.VenueURL != e.venueURL {
		phys.Audit(storer, "%s:: venueURL = %q", context, u.VenueURL)
		e.venueURL = u.VenueURL
	}
	if u.Activation != e.activation {
		phys.Audit(storer, "%s:: activation = %q", context, u.Activation)
		e.activation = u.Activation
	}
	if u.Details != e.details {
		phys.Audit(storer, "%s:: details = %q", context, u.Details)
		e.details = u.Details
	}
	if u.Flags != e.flags {
		phys.Audit(storer, "%s:: flags = 0x%x", context, u.Flags)
		e.flags = u.Flags
	}
}

const duplicateNameSQL = `SELECT 1 FROM event WHERE id!=? AND name=? AND start LIKE ?`

// DuplicateName returns whether the name specified in the Updater would be a
// duplicate if applied.
func (u *Updater) DuplicateName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Name)
		stmt.BindText(u.Start[:10] + "%")
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver event.
func (e *Event) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM event WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(e.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE Event %s %q [%d]", e.Start()[:10], e.Name(), e.ID())
	phys.Unindex(storer, e)
}
