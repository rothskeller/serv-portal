package taskperson

import (
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
)

const SetEventFields = event.FID | event.FStart | event.FName
const SetTaskFields = task.FID | task.FEvent | task.FName
const SetPersonFields = person.FID | person.FInformalName

// Set sets the relationship between the specified Person and Task.  The parent
// Event and the previous relationship *may* be provided to avoid lookups.
func Set(storer phys.Storer, e *event.Event, t *task.Task, p *person.Person, minutes uint, flags Flag) {
	if e == nil {
		e = event.WithID(storer, t.Event(), SetEventFields)
	}
	pminutes, pflags := Get(storer, t.ID(), p.ID())
	if minutes == pminutes && flags == pflags {
		return
	}
	if minutes == 0 && flags == 0 {
		phys.SQL(storer, "DELETE FROM task_person WHERE task=? AND person=?", func(stmt *phys.Stmt) {
			stmt.BindInt(int(t.ID()))
			stmt.BindInt(int(p.ID()))
			stmt.Step()
		})
		goto AUDIT
	}
	phys.SQL(storer, "UPDATE task_person SET flags=?, minutes=? WHERE task=? AND person=?", func(stmt *phys.Stmt) {
		stmt.BindHexInt(int(flags))
		stmt.BindNullInt(int(minutes))
		stmt.BindInt(int(t.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		goto AUDIT
	}
	phys.SQL(storer, "INSERT INTO task_person (task, person, flags, minutes) VALUES (?,?,?,?)", func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.BindHexInt(int(flags))
		stmt.BindNullInt(int(minutes))
		stmt.Step()
	})
AUDIT:
	if flags != pflags {
		phys.Audit(storer, "Event %s %q [%d]:: Task %q [%d]:: Person %q [%d]:: flags = 0x%x",
			e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), p.InformalName(), p.ID(), flags)
	}
	if minutes != pminutes {
		phys.Audit(storer, "Event %s %q [%d]:: Task %q [%d]:: Person %q [%d]:: minutes = %d",
			e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), p.InformalName(), p.ID(), minutes)
	}
}
