package venue

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList generates a comma-separated list of column names for the specified
// venue fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("v.id")
	}
	if fields&FName != 0 {
		sb.WriteString(sep())
		sb.WriteString("v.name")
	}
	if fields&FURL != 0 {
		sb.WriteString(sep())
		sb.WriteString("v.url")
	}
	if fields&FFlags != 0 {
		sb.WriteString(sep())
		sb.WriteString("v.flags")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (v *Venue) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		v.id = ID(stmt.ColumnInt())
	}
	if fields&FName != 0 {
		v.name = stmt.ColumnText()
	}
	if fields&FURL != 0 {
		v.url = stmt.ColumnText()
	}
	if fields&FFlags != 0 {
		v.flags = Flag(stmt.ColumnHexInt())
	}
	v.fields |= fields
}
