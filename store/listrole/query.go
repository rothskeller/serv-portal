package listrole

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/role"
)

const getSQL = `SELECT sender, submodel FROM list_role WHERE list=? AND role=?`

// Get returns the privileges assigned to a list for a role.
func Get(storer phys.Storer, lid list.ID, rid role.ID) (sender bool, submodel SubscriptionModel) {
	phys.SQL(storer, getSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(lid))
		stmt.BindInt(int(rid))
		if stmt.Step() {
			sender = stmt.ColumnBool()
			submodel = SubscriptionModel(stmt.ColumnInt())
		}
	})
	return sender, submodel
}

const allListsForRoleSQL1 = `SELECT `
const allListsForRoleSQL2 = `, lr.sender, lr.submodel FROM list l, list_role lr WHERE lr.list=l.id AND lr.role=? ORDER BY l.name`

var allListsForRoleSQLCache string

// AllListsForRole fetches all of the lists for which the specified role grants
// privileges, in order by name.
func AllListsForRole(storer phys.Storer, rid role.ID, fn func(list *list.List, sender bool, submodel SubscriptionModel)) {
	if allListsForRoleSQLCache == "" {
		var sb strings.Builder
		sb.WriteString(allListsForRoleSQL1)
		sb.WriteString(list.ColumnList)
		sb.WriteString(allListsForRoleSQL2)
		allListsForRoleSQLCache = sb.String()
	}
	phys.SQL(storer, allListsForRoleSQLCache, func(stmt *phys.Stmt) {
		var l list.List
		var sender bool
		var submodel SubscriptionModel

		stmt.BindInt(int(rid))
		for stmt.Step() {
			l.Scan(stmt)
			sender = stmt.ColumnBool()
			submodel = SubscriptionModel(stmt.ColumnInt())
			fn(&l, sender, submodel)
		}
	})
}

const allRolesForListSQL1 = `SELECT `
const allRolesForListSQL2 = `, lr.sender, lr.submodel FROM role r, list_role lr WHERE lr.role=r.id AND lr.list=? ORDER BY r.name`

var allRolesForListSQLCache map[role.Fields]string

// AllRolesForList fetches all of the roles for which privileges are granted on
// the specified list, in order by name.
func AllRolesForList(storer phys.Storer, lid list.ID, fields role.Fields, fn func(rl *role.Role, sender bool, submodel SubscriptionModel)) {
	if allRolesForListSQLCache == nil {
		allRolesForListSQLCache = make(map[role.Fields]string)
	}
	if allRolesForListSQLCache[fields] == "" {
		var sb strings.Builder
		sb.WriteString(allRolesForListSQL1)
		role.ColumnList(&sb, fields)
		sb.WriteString(allRolesForListSQL2)
		allRolesForListSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allRolesForListSQLCache[fields], func(stmt *phys.Stmt) {
		var rl role.Role
		var sender bool
		var submodel SubscriptionModel

		stmt.BindInt(int(lid))
		for stmt.Step() {
			rl.Scan(stmt, fields)
			sender = stmt.ColumnBool()
			submodel = SubscriptionModel(stmt.ColumnInt())
			fn(&rl, sender, submodel)
		}
	})
}
