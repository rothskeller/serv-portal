package role

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList generates a comma-separated list of column names for the specified
// role fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.id")
	}
	if fields&FName != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.name")
	}
	if fields&FTitle != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.title")
	}
	if fields&FPriority != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.priority")
	}
	if fields&FOrg != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.org")
	}
	if fields&FPrivLevel != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.privlevel")
	}
	if fields&FFlags != 0 {
		sb.WriteString(sep())
		sb.WriteString("r.flags")
	}
	if fields&FImplies != 0 {
		panic("cannot fetch FImplies using ColumnList/Scan")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (r *Role) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		r.id = ID(stmt.ColumnInt())
	}
	if fields&FName != 0 {
		r.name = stmt.ColumnText()
	}
	if fields&FTitle != 0 {
		r.title = stmt.ColumnText()
	}
	if fields&FPriority != 0 {
		r.priority = uint(stmt.ColumnInt())
	}
	if fields&FOrg != 0 {
		r.org = enum.Org(stmt.ColumnInt())
	}
	if fields&FPrivLevel != 0 {
		r.privLevel = enum.PrivLevel(stmt.ColumnInt())
	}
	if fields&FFlags != 0 {
		r.flags = Flags(stmt.ColumnInt())
	}
	r.fields |= fields &^ FImplies
}

const readImpliesSQL = `SELECT implied FROM role_implies WHERE implier=? ORDER BY implied`

func (r *Role) readImplies(storer phys.Storer) {
	r.implies = r.implies[:0]
	phys.SQL(storer, readImpliesSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(r.ID()))
		for stmt.Step() {
			r.implies = append(r.implies, ID(stmt.ColumnInt()))
		}
	})
	r.fields |= FImplies
}
