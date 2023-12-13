package taskrole

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/task"
)

const rolesForTaskSQL1 = `SELECT `
const rolesForTaskSQL2 = ` FROM role r, task_role tr WHERE r.id=tr.role AND tr.task=?`

var rolesForTaskSQLCache map[role.Fields]string

// Get fetches the set of Roles for a Task, in unspecified order.
func Get(storer phys.Storer, tid task.ID, fields role.Fields, fn func(*role.Role)) {
	if rolesForTaskSQLCache == nil {
		rolesForTaskSQLCache = make(map[role.Fields]string)
	}
	if _, ok := rolesForTaskSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(rolesForTaskSQL1)
		role.ColumnList(&sb, fields)
		sb.WriteString(rolesForTaskSQL2)
		rolesForTaskSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, rolesForTaskSQLCache[fields], func(stmt *phys.Stmt) {
		var r role.Role

		stmt.BindInt(int(tid))
		for stmt.Step() {
			r.Scan(stmt, fields)
			fn(&r)
		}
	})
}
