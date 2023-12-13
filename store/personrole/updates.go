package personrole

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
)

const addRoleSQL = `INSERT OR IGNORE INTO person_role (person, role, explicit) VALUES (?,?,1)`

// AddRole adds the specified role to the specified person.  It is a no-op if
// the person already holds the role.
func AddRole(storer phys.Storer, p *person.Person, r *role.Role) {
	phys.SQL(storer, addRoleSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		stmt.BindInt(int(r.ID()))
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "Person %q [%d]:: ADD Role %q [%d]", p.InformalName(), p.ID(), r.Name(), r.ID())
	}
}

const removeRoleSQL = `DELETE FROM person_role WHERE person=? AND role=?`

// RemoveRole removes the specified role from the specified person.  It is a
// no-op if the person lacks the role.
func RemoveRole(storer phys.Storer, p *person.Person, r *role.Role) {
	phys.SQL(storer, removeRoleSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		stmt.BindInt(int(r.ID()))
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "Person %q [%d]:: REMOVE Role %q [%d]", p.InformalName(), p.ID(), r.Name(), r.ID())
	}
}
