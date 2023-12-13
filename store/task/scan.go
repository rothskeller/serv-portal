package task

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList generates a comma-separated list of column names for the specified
// Task fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.id")
	}
	if fields&FEvent != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.event")
	}
	if fields&FName != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.name")
	}
	if fields&FOrg != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.org")
	}
	if fields&FFlags != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.flags")
	}
	if fields&FDetails != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.details")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (t *Task) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		t.id = ID(stmt.ColumnInt())
	}
	if fields&FEvent != 0 {
		t.event = event.ID(stmt.ColumnInt())
	}
	if fields&FName != 0 {
		t.name = stmt.ColumnText()
	}
	if fields&FOrg != 0 {
		t.org = enum.Org(stmt.ColumnInt())
	}
	if fields&FFlags != 0 {
		t.flags = Flag(stmt.ColumnHexInt())
	}
	if fields&FDetails != 0 {
		t.details = stmt.ColumnText()
	}
	t.fields |= fields
}
