package shiftperson

import (
	"strings"

	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
)

// Get returns whether the specified Person is signed up for the specified
// Shift.  Zero means no; a negative value means an explicit decline; a positive
// value means yes.
func Get(storer phys.Storer, sid shift.ID, pid person.ID) (signedup int) {
	phys.SQL(storer, "SELECT signed_up FROM shift_person WHERE shift=? AND person=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(sid))
		stmt.BindInt(int(pid))
		if stmt.Step() {
			signedup = stmt.ColumnInt()
		}
	})
	return
}

// HasSignups returns whether any people are signed up for the specified Shift.
func HasSignups(storer phys.Storer, sid shift.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM shift_person WHERE shift=? AND signed_up>0 LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(sid))
		found = stmt.Step()
	})
	return found
}

// TaskHasSignups returns whether any people are signed up for any Shift on the
// specified Task.
func TaskHasSignups(storer phys.Storer, tid task.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM shift_person sp, shift s WHERE sp.shift=s.id AND s.task=? AND sp.signed_up>0 LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(tid))
		found = stmt.Step()
	})
	return found
}

// EventHasSignups returns whether any people are signed up for any Shift of any
// Task of the specified Event.
func EventHasSignups(storer phys.Storer, eid event.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM shift_person sp, shift s, task t WHERE sp.shift=s.id AND s.task=t.id AND t.event=? AND sp.signed_up>0 LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(eid))
		found = stmt.Step()
	})
	return found
}

var peopleForShiftSQLCache map[person.Fields]string

// PeopleForShift fetches all people signed up for the specified Shift, in the
// order they signed up.
func PeopleForShift(storer phys.Storer, sid shift.ID, personFields person.Fields, fn func(*person.Person)) {
	if peopleForShiftSQLCache == nil {
		peopleForShiftSQLCache = make(map[person.Fields]string)
	}
	if _, ok := peopleForShiftSQLCache[personFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		person.ColumnList(&sb, personFields)
		sb.WriteString(" FROM person p, shift_person sp WHERE sp.shift=? AND p.id=sp.person AND sp.signed_up>0 ORDER BY sp.signed_up")
		peopleForShiftSQLCache[personFields] = sb.String()
	}
	phys.SQL(storer, peopleForShiftSQLCache[personFields], func(stmt *phys.Stmt) {
		var p person.Person

		stmt.BindInt(int(sid))
		for stmt.Step() {
			p.Scan(stmt, personFields)
			fn(&p)
		}
	})
}

// OverlappingSignup returns whether the specified person is signed up for any
// shift of any task overlaps the specified time range.
func OverlappingSignup(storer phys.Storer, pid person.ID, start, end string) (overlap bool) {
	phys.SQL(storer, "SELECT 1 FROM shift_person sp, shift s WHERE sp.shift=s.id AND sp.person=? AND sp.signed_up>0 AND ?<s.end AND ?>s.start", func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		stmt.BindText(start)
		stmt.BindText(end)
		overlap = stmt.Step()
	})
	return overlap
}
