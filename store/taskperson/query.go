package taskperson

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
)

// Get returns the relationship between the specified Task and Person.
func Get(storer phys.Storer, tid task.ID, pid person.ID) (minutes uint, flags Flag) {
	phys.SQL(storer, "SELECT flags, minutes FROM task_person WHERE task=? AND person=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(tid))
		stmt.BindInt(int(pid))
		if stmt.Step() {
			flags = Flag(stmt.ColumnHexInt())
			minutes = uint(stmt.ColumnInt())
		}
	})
	return
}

// ExistsForTask returns whether there are any taskperson records (i.e., anyone
// recorded as attended or credited, or anyone who has recorded volunteer hours)
// for the specified Task.
func ExistsForTask(storer phys.Storer, tid task.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM task_person WHERE task=? LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(tid))
		found = stmt.Step()
	})
	return found
}

// ExistsForEvent returns whether there are any taskperson records (i.e., anyone
// recorded as attended or credited, or anyone who has recorded volunteer hours)
// for any Task of the specified Event.
func ExistsForEvent(storer phys.Storer, eid event.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM task_person tp, task t WHERE tp.task=t.id AND t.event=? LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(eid))
		found = stmt.Step()
	})
	return found
}

var peopleForTaskSQLCache map[person.Fields]string

// PeopleForTask fetches all relationships between any Person and the specified
// Task.  They are returned in random order.
func PeopleForTask(storer phys.Storer, tid task.ID, personFields person.Fields, fn func(*person.Person, uint, Flag)) {
	if peopleForTaskSQLCache == nil {
		peopleForTaskSQLCache = make(map[person.Fields]string)
	}
	if _, ok := peopleForTaskSQLCache[personFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT tp.flags, tp.minutes, ")
		person.ColumnList(&sb, personFields)
		sb.WriteString(" FROM person p, task_person tp WHERE tp.task=? AND p.id=tp.person")
		peopleForTaskSQLCache[personFields] = sb.String()
	}
	phys.SQL(storer, peopleForTaskSQLCache[personFields], func(stmt *phys.Stmt) {
		var p person.Person
		var minutes uint
		var flags Flag

		stmt.BindInt(int(tid))
		for stmt.Step() {
			flags = Flag(stmt.ColumnHexInt())
			minutes = uint(stmt.ColumnInt())
			p.Scan(stmt, personFields)
			fn(&p, minutes, flags)
		}
	})
}

// HasRoleForTask returns whether the specified person has any of the roles for
// the specified task.
func HasRoleForTask(storer phys.Storer, tid task.ID, pid person.ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM task_role tr, person_role pr WHERE tr.role=pr.role AND tr.task=? AND pr.person=? LIMIT 1", func(stmt *phys.Stmt) {
		stmt.BindInt(int(tid))
		stmt.BindInt(int(pid))
		found = stmt.Step()
	})
	return found
}

// AllBetween returns all task/person relationships for the specified person for
// events between the specified dates.
func AllBetween(storer phys.Storer, start, end string, pid person.ID, eventFields event.Fields, taskFields task.Fields, fn func(e *event.Event, t *task.Task, minutes uint, flags Flag)) {
	var sb strings.Builder

	sb.WriteString("SELECT tp.minutes, tp.flags")
	if eventFields != 0 {
		sb.WriteString(", ")
		event.ColumnList(&sb, eventFields)
	}
	if taskFields != 0 {
		sb.WriteString(", ")
		task.ColumnList(&sb, taskFields)
	}
	sb.WriteString(" FROM task_person tp, task t, event e WHERE tp.task=t.id AND t.event=e.id AND tp.person=? AND e.start>=? AND e.end<? ORDER BY e.start, e.end, e.id, t.sort")
	phys.SQL(storer, sb.String(), func(stmt *phys.Stmt) {
		var e event.Event
		var t task.Task
		var minutes uint
		var flags Flag

		stmt.BindInt(int(pid))
		stmt.BindText(start)
		stmt.BindText(end)
		for stmt.Step() {
			minutes = uint(stmt.ColumnInt())
			flags = Flag(stmt.ColumnHexInt())
			e.Scan(stmt, eventFields)
			t.Scan(stmt, taskFields)
			fn(&e, &t, minutes, flags)
		}
	})
}

const minutesBetweenSQL = `
SELECT e.id, tp.person, t.org, tp.minutes
FROM task_person tp, task t, event e, person p
WHERE tp.task=t.id AND t.event=e.id
  AND tp.minutes!=0
  AND e.start>=?
  AND e.end<?`

// MinutesBetween fetches all minutes entries for events between the specified
// dates.
func MinutesBetween(storer phys.Storer, start, end string, fn func(event.ID, person.ID, enum.Org, uint)) {
	phys.SQL(storer, minutesBetweenSQL, func(stmt *phys.Stmt) {
		stmt.BindText(start)
		stmt.BindText(end)
		for stmt.Step() {
			var eid = event.ID(stmt.ColumnInt())
			var pid = person.ID(stmt.ColumnInt())
			var org = enum.Org(stmt.ColumnInt())
			var minutes = uint(stmt.ColumnInt())
			fn(eid, pid, org, minutes)
		}
	})
}
