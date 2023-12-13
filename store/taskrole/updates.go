package taskrole

import (
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/task"
)

const addRoleSQL = `INSERT INTO task_role (task, role) VALUES (?,?)`
const removeRoleSQL = `DELETE FROM task_role WHERE task=? AND role=?`

// Set sets the Roles for the specified Task.  The Event containing the Task and
// the previous set of Roles *may* be specified to lookups and allocations.
func Set(storer phys.Storer, e *event.Event, t *task.Task, roles, prev []*role.Role) {
	const eventFields = event.FID | event.FStart | event.FName
	const roleFields = role.FID | role.FName
	var (
		rmap = make(map[role.ID]*role.Role)
		pmap = make(map[role.ID]*role.Role)
	)
	if e == nil || e.Fields()&eventFields != eventFields || e.ID() != t.Event() {
		e = event.WithID(storer, t.Event(), eventFields)
	}
	for _, rl := range roles {
		if rl.Fields()&roleFields != roleFields {
			rl = role.WithID(storer, rl.ID(), roleFields)
		}
		rmap[rl.ID()] = rl
	}
	if prev == nil {
		Get(storer, t.ID(), roleFields, func(rl *role.Role) {
			if rl2 := rmap[rl.ID()]; rl2 != nil {
				pmap[rl.ID()] = rl2
			} else {
				pmap[rl.ID()] = rl.Clone()
			}
		})
	} else {
		for _, rl := range prev {
			if rl2 := rmap[rl.ID()]; rl2 != nil {
				pmap[rl.ID()] = rl2
			} else if rl.Fields()&roleFields == roleFields {
				pmap[rl.ID()] = rl.Clone()
			} else {
				pmap[rl.ID()] = role.WithID(storer, rl.ID(), roleFields)
			}
		}
	}
	for rid, rl := range rmap {
		if pmap[rid] == nil {
			phys.SQL(storer, addRoleSQL, func(stmt *phys.Stmt) {
				stmt.BindInt(int(t.ID()))
				stmt.BindInt(int(rid))
				stmt.Step()
			})
			phys.Audit(storer, "Event %s %q [%d]:: Task %q [%d]:: ADD Role %q [%d]",
				e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), rl.Name(), rid)
		}
	}
	for rid, rl := range pmap {
		if rmap[rid] == nil {
			phys.SQL(storer, removeRoleSQL, func(stmt *phys.Stmt) {
				stmt.BindInt(int(t.ID()))
				stmt.BindInt(int(rid))
				stmt.Step()
			})
			phys.Audit(storer, "Event %s %q [%d]:: Task %q [%d]:: REMOVE Role %q [%d]",
				e.Start()[:10], e.Name(), e.ID(), t.Name(), t.ID(), rl.Name(), rid)
		}
	}
}
