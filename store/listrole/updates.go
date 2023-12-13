package listrole

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/role"
)

const deleteListRoleSQL = `DELETE FROM list_role WHERE list=? AND role=?`
const setListRoleSQL = `
INSERT INTO list_role (list, role, sender, submodel) VALUES (?,?,?,?)
ON CONFLICT DO UPDATE SET sender=?3, submodel=?4 WHERE sender!=?3 OR submodel!=?4`

// SetListRole sets the sender privilege and subscription model on the specified
// list for the specified role.
func SetListRole(storer phys.Storer, l *list.List, r *role.Role, sender bool, submodel SubscriptionModel) {
	if !sender && submodel == 0 {
		phys.SQL(storer, deleteListRoleSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(l.ID))
			stmt.BindInt(int(r.ID()))
			stmt.Step()
		})
		if phys.RowsAffected(storer) != 0 {
			phys.Audit(storer, "List %q [%d]:: DELETE Role %q [%d]", l.Name, l.ID, r.Name(), r.ID())
		}
	} else {
		phys.SQL(storer, setListRoleSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(l.ID))
			stmt.BindInt(int(r.ID()))
			stmt.BindBool(sender)
			stmt.BindInt(int(submodel))
			stmt.Step()
		})
		if phys.RowsAffected(storer) != 0 {
			phys.Audit(storer, "List %q [%d]:: Role %q [%d]:: sender = %v", l.Name, l.ID, r.Name(), r.ID(), sender)
			phys.Audit(storer, "List %q [%d]:: Role %q [%d]:: submodel = %s [%d]", l.Name, l.ID, r.Name(), r.ID(), submodel, submodel)
		}
	}
}
