package personrole

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
)

const personHasRoleSQL = `SELECT explicit FROM person_role WHERE person=? AND role=? LIMIT 1`

// PersonHasRole returns whether the specified person has the specified role.
func PersonHasRole(storer phys.Storer, pid person.ID, rid role.ID) (held, explicit bool) {
	phys.SQL(storer, personHasRoleSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		stmt.BindInt(int(rid))
		held = stmt.Step()
		if held {
			explicit = stmt.ColumnBool()
		}
	})
	return held, explicit
}

const rolesForPersonSQL1 = `SELECT `
const rolesForPersonSQL2 = `, pr.explicit FROM role r, person_role pr WHERE r.id=pr.role AND pr.person=? ORDER BY r.priority`

var rolesForPersonSQLCache map[role.Fields]string

// RolesForPerson fetches the roles assigned to the specified person, in
// priority order.
func RolesForPerson(storer phys.Storer, pid person.ID, fields role.Fields, fn func(r *role.Role, explicit bool)) {
	if rolesForPersonSQLCache == nil {
		rolesForPersonSQLCache = make(map[role.Fields]string)
	}
	if _, ok := rolesForPersonSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(rolesForPersonSQL1)
		role.ColumnList(&sb, fields)
		sb.WriteString(rolesForPersonSQL2)
		rolesForPersonSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, rolesForPersonSQLCache[fields], func(stmt *phys.Stmt) {
		var r role.Role
		var explicit bool

		stmt.BindInt(int(pid))
		for stmt.Step() {
			r.Scan(stmt, fields)
			explicit = stmt.ColumnBool()
			fn(&r, explicit)
		}
	})
}

const peopleForRoleSQL1 = `SELECT `
const peopleForRoleSQL2 = `, pr.explicit FROM person p, person_role pr WHERE p.id=pr.person AND pr.role=? ORDER BY p.sort_name`

var peopleForRoleSQLCache map[person.Fields]string

// PeopleForRole fetches the people who hold the specified role, in order by
// person sort name.
func PeopleForRole(storer phys.Storer, rid role.ID, fields person.Fields, fn func(p *person.Person, explicit bool)) {
	if peopleForRoleSQLCache == nil {
		peopleForRoleSQLCache = make(map[person.Fields]string)
	}
	if _, ok := peopleForRoleSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(peopleForRoleSQL1)
		person.ColumnList(&sb, fields)
		sb.WriteString(peopleForRoleSQL2)
		peopleForRoleSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, peopleForRoleSQLCache[fields], func(stmt *phys.Stmt) {
		var p person.Person
		var explicit bool

		stmt.BindInt(int(rid))
		for stmt.Step() {
			p.Scan(stmt, fields)
			explicit = stmt.ColumnBool()
			fn(&p, explicit)
		}
	})
}

// PeopleCountForRole returns the number of people who hold the specified role.
func PeopleCountForRole(storer phys.Storer, rid role.ID) (count int) {
	phys.SQL(storer, "SELECT COUNT(*) FROM person_role WHERE role=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(rid))
		stmt.Step()
		count = stmt.ColumnInt()
	})
	return count
}
