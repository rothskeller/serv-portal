package role

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/util"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FName | FTitle | FPriority | FOrg | FPrivLevel | FFlags | FImplies

// Updater is a structure that can be filled with data for a new or changed
// role, and then later applied.  For creating new roles, it can simply be
// instantiated with new().  For updating existing roles, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the role being changed.
type Updater struct {
	ID        ID
	Name      string
	Title     string
	Priority  uint
	Org       enum.Org
	PrivLevel enum.PrivLevel
	Flags     Flags
	Implies   []ID
}

// Updater returns a new Updater for the specified role, with its data matching
// the current data for the role.  The role must have fetched UpdaterFields.
func (r *Role) Updater() *Updater {
	var implies []ID

	if r.fields&UpdaterFields != UpdaterFields {
		panic("Role.Updater called without fetching UpdaterFields")
	}
	if len(r.implies) != 0 {
		implies = make([]ID, len(r.implies))
		copy(implies, r.implies)
	}
	return &Updater{
		ID:        r.id,
		Name:      r.name,
		Title:     r.title,
		Priority:  r.priority,
		Org:       r.org,
		PrivLevel: r.privLevel,
		Flags:     r.flags,
		Implies:   implies,
	}
}

const createSQL = `INSERT INTO role (id, name, title, priority, org, privlevel, flags) VALUES (?,?,?,?,?,?,?)`

// Create creates a new role, with the data in the Updater.  If priority is 0,
// it is set to the end of the list (i.e., max(existing)+1).
func Create(storer phys.Storer, u *Updater) (r *Role) {
	if u.Priority == 0 {
		phys.SQL(storer, "SELECT MAX(priority) FROM role", func(stmt *phys.Stmt) {
			stmt.Step()
			u.Priority = uint(stmt.ColumnInt()) + 1
		})
	}
	r = new(Role)
	r.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			r.id = u.ID
		} else {
			r.id = ID(phys.LastInsertRowID(storer))
		}
	})
	storeImplies(storer, r.id, u.Implies, false)
	r.auditAndUpdate(storer, u, true)
	phys.Index(storer, r)
	return r
}

const updateSQL = `UPDATE role SET name=?, title=?, priority=?, org=?, privlevel=?, flags=? WHERE id=?`

// Update updates the existing venue, with the data in the Updater.
func (r *Role) Update(storer phys.Storer, u *Updater) {
	if r.fields&UpdaterFields != UpdaterFields {
		panic("Role.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(r.id))
		stmt.Step()
	})
	storeImplies(storer, r.id, u.Implies, true)
	r.auditAndUpdate(storer, u, false)
	phys.Index(storer, r)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindText(u.Name)
	stmt.BindNullText(u.Title)
	stmt.BindNullInt(int(u.Priority))
	stmt.BindInt(int(u.Org))
	stmt.BindInt(int(u.PrivLevel))
	stmt.BindInt(int(u.Flags))
}

const deleteImpliesSQL = `DELETE FROM role_implies WHERE implier=?`
const insertImpliesSQL = `INSERT INTO role_implies (implier, implied) VALUES (?,?)`

func storeImplies(storer phys.Storer, implier ID, implies []ID, delfirst bool) {
	if delfirst {
		phys.SQL(storer, deleteImpliesSQL, func(stmt *phys.Stmt) {
			stmt.BindInt(int(implier))
			stmt.Step()
		})
	}
	if len(implies) != 0 {
		util.SortIDList(implies)
		phys.SQL(storer, insertImpliesSQL, func(stmt *phys.Stmt) {
			for _, implied := range implies {
				stmt.BindInt(int(implier))
				stmt.BindInt(int(implied))
				stmt.Step()
				stmt.Reset()
			}
		})
	}
}

const roleNameSQL = `SELECT name FROM role WHERE id=?`

func (r *Role) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Role %q [%d]", u.Name, r.id)
	if create {
		context = "ADD " + context
	}
	if u.Name != r.name {
		phys.Audit(storer, "%s:: name = %q", context, u.Name)
		r.name = u.Name
	}
	if u.Title != r.title {
		phys.Audit(storer, "%s:: title = %q", context, u.Title)
		r.title = u.Title
	}
	if u.Priority != r.priority {
		phys.Audit(storer, "%s:: priority = %d", context, u.Priority)
		r.priority = u.Priority
	}
	if u.Org != r.org {
		phys.Audit(storer, "%s:: org = %s [%d]", context, u.Org, u.Org)
		r.org = u.Org
	}
	if u.PrivLevel != r.privLevel {
		phys.Audit(storer, "%s:: privLevel = %s [%d]", context, u.PrivLevel, u.PrivLevel)
		r.privLevel = u.PrivLevel
	}
	if u.Flags != r.flags {
		phys.Audit(storer, "%s:: flags = 0x%x", context, u.Flags)
		r.flags = u.Flags
	}
	if !util.EqualIDList(r.implies, u.Implies) {
		if len(u.Implies) == 0 {
			phys.Audit(storer, "%s:: implies = []", context)
			r.implies = nil
		} else {
			phys.SQL(storer, roleNameSQL, func(stmt *phys.Stmt) {
				for i, implied := range u.Implies {
					stmt.BindInt(int(implied))
					if stmt.Step() {
						phys.Audit(storer, "%s:: implies:: [%d] = %q [%d]", context, i+1, stmt.ColumnText(), implied)
					} else {
						panic("implies list includes nonexistent role")
					}
					stmt.Reset()
				}
			})
		}
	}
}

const duplicateNameSQL = `SELECT 1 FROM role WHERE id!=? AND name=?`

// DuplicateName returns whether the name specified in the Updater would be a
// duplicate if applied.
func (u *Updater) DuplicateName(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateNameSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Name)
		found = stmt.Step()
	})
	return found
}

const duplicateTitleSQL = `SELECT 1 FROM role WHERE id!=? AND title=?`

// DuplicateTitle returns whether the title specified in the Updater would be a
// duplicate if applied.
func (u *Updater) DuplicateTitle(storer phys.Storer) (found bool) {
	if u.Title == "" {
		return false
	}
	phys.SQL(storer, duplicateTitleSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindText(u.Title)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver role.
func (r *Role) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM role WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(r.ID()))
		stmt.Step()
	})
	phys.SQL(storer, "UPDATE role SET priority=priority-1 WHERE priority>?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(r.Priority()))
		stmt.Step()
	})
	phys.Unindex(storer, r)
	phys.Audit(storer, "DELETE Role %q [%d]", r.Name(), r.ID())
}

// Reorder changes the priority of a role, shifting the other role priorities
// to keep the list unique and compact.  If "to" is zero, the role is moved to
// the end of the list.
func (r *Role) Reorder(storer phys.Storer, to uint) {
	if to == 0 {
		phys.SQL(storer, "SELECT MAX(priority) FROM role", func(stmt *phys.Stmt) {
			stmt.Step()
			to = uint(stmt.ColumnInt())
		})
	}
	if r.Priority() == to {
		return
	}
	if r.Priority() < to {
		phys.SQL(storer, "UPDATE role SET priority=priority-1 WHERE priority BETWEEN ? AND ?", func(stmt *phys.Stmt) {
			stmt.BindInt(int(r.Priority()))
			stmt.BindInt(int(to))
			stmt.Step()
		})
	} else {
		phys.SQL(storer, "UPDATE role SET priority=priority+1 WHERE priority BETWEEN ? AND ?", func(stmt *phys.Stmt) {
			stmt.BindInt(int(to))
			stmt.BindInt(int(r.Priority()))
			stmt.Step()
		})
	}
	phys.SQL(storer, "UPDATE role SET priority=? WHERE id=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(to))
		stmt.BindInt(int(r.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "Role %q [%d]:: priority = %d (shift from %d)", r.Name(), r.ID(), to, r.Priority())
	r.priority = to
}
