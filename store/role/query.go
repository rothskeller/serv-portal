package role

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

var withIDSQLCache map[Fields]string

// WithID returns the role with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (r *Role) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields&^FImplies]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FImplies)
		sb.WriteString(" FROM role r WHERE r.id=?")
		withIDSQLCache[fields&^FImplies] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields&^FImplies], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			r = new(Role)
			r.Scan(stmt, fields&^FImplies)
			r.id = id
			r.fields |= FID
			if fields&FImplies != 0 {
				r.readImplies(storer)
			}
		}
	})
	return r
}

var allSQLCache map[Fields]string

// All reads each role from the database, in order by priority.
func All(storer phys.Storer, fields Fields, fn func(*Role)) {
	if allSQLCache == nil {
		allSQLCache = make(map[Fields]string)
	}
	if _, ok := allSQLCache[fields&^FImplies]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FImplies)
		sb.WriteString(" FROM role r ORDER BY r.priority")
		allSQLCache[fields&^FImplies] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields&^FImplies], func(stmt *phys.Stmt) {
		var r Role
		for stmt.Step() {
			r.Scan(stmt, fields&^FImplies)
			if fields&FImplies != 0 {
				r.readImplies(storer)
			}
			fn(&r)
		}
	})
}

var allWithOrgSQLCache map[Fields]string

// AllWithOrg reads each role in the specified organization from the database,
// in order by priority.
func AllWithOrg(storer phys.Storer, fields Fields, org enum.Org, fn func(*Role)) {
	if allWithOrgSQLCache == nil {
		allWithOrgSQLCache = make(map[Fields]string)
	}
	if _, ok := allWithOrgSQLCache[fields&^FImplies]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FImplies)
		sb.WriteString(" FROM role r WHERE r.org=? ORDER BY r.priority")
		allWithOrgSQLCache[fields&^FImplies] = sb.String()
	}
	phys.SQL(storer, allWithOrgSQLCache[fields&^FImplies], func(stmt *phys.Stmt) {
		var r Role
		stmt.BindInt(int(org))
		for stmt.Step() {
			r.Scan(stmt, fields&^FImplies)
			if fields&FImplies != 0 {
				r.readImplies(storer)
			}
			fn(&r)
		}
	})
}

const allThatImplySQL = `
WITH RECURSIVE implies AS (
	SELECT ri.implier, ri.implied FROM role_implies ri
	UNION
	SELECT implies.implier, ri.implied FROM role_implies ri, implies WHERE implies.implied=ri.implier
)
SELECT implier FROM implies WHERE implied=?`

// AllThatImply reads from the database the set of IDs of roles that imply the
// specified one, directly or indirectly.
func AllThatImply(storer phys.Storer, implied ID) (impliers map[ID]struct{}) {
	impliers = make(map[ID]struct{})
	phys.SQL(storer, allThatImplySQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(implied))
		for stmt.Step() {
			impliers[ID(stmt.ColumnInt())] = struct{}{}
		}
	})
	return impliers
}
