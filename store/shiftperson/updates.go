package shiftperson

import (
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
)

const nextSignupSQL = `SELECT COALESCE(MAX(signed_up), 0) FROM shift_person WHERE shift=?`
const signUpSQL = `INSERT INTO shift_person (shift, person, signed_up) VALUES (?,?,?) ON CONFLICT DO UPDATE SET signed_up=?3 WHERE shift_person.signed_up<0`
const SignUpEventFields = event.FID | event.FStart | event.FName
const SignUpTaskFields = task.FID | task.FEvent | task.FName
const SignUpShiftFields = shift.FID | shift.FTask
const SignUpPersonFields = person.FID | person.FInformalName

// SignUp signs the specified Person up for the specified Shift.  The parent
// Event and Task *may* be provided to avoid lookups.  This function is a no-op
// if the Person is already signed up.
func SignUp(storer phys.Storer, e *event.Event, t *task.Task, s *shift.Shift, p *person.Person) {
	var signedUp int

	if t == nil {
		t = task.WithID(storer, s.Task(), SignUpTaskFields)
	}
	if e == nil {
		e = event.WithID(storer, t.Event(), SignUpEventFields)
	}
	phys.SQL(storer, nextSignupSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.Step()
		signedUp = stmt.ColumnInt() + 1
		if signedUp < 1 {
			signedUp = 1
		}
	})
	phys.SQL(storer, signUpSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.BindInt(signedUp)
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "Event %s %q [%d]:: Task %s [%d]:: Shift %d:: sign up %q [%d]",
			e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), s.ID(), p.InformalName(), p.ID())
	}
}

const nextDeclineSQL = `SELECT COALESCE(MIN(signed_up), 0) FROM shift_person WHERE shift=?`
const declineSQL = `INSERT INTO shift_person (shift, person, signed_up) VALUES (?,?,?) ON CONFLICT DO UPDATE SET signed_up=?3 WHERE shift_person.signed_up>0`

// Decline marks the specified Person as having declined the specified Shift
// (and in the process removes any existing signup by that Person for that
// Shift).  The parent Event and Task *may* be provided to avoid lookups.  This
// function is a no-op if the Person has already declined the Shift.
func Decline(storer phys.Storer, e *event.Event, t *task.Task, s *shift.Shift, p *person.Person) {
	var signedUp int

	if t == nil {
		t = task.WithID(storer, s.Task(), SignUpTaskFields)
	}
	if e == nil {
		e = event.WithID(storer, t.Event(), SignUpEventFields)
	}
	phys.SQL(storer, nextDeclineSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.Step()
		signedUp = stmt.ColumnInt() - 1
		if signedUp > -1 {
			signedUp = -1
		}
	})
	phys.SQL(storer, declineSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.BindInt(signedUp)
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "Event %s %q [%d]:: Task %s [%d]:: Shift %d:: decline %q [%d]",
			e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), s.ID(), p.InformalName(), p.ID())
	}
}
