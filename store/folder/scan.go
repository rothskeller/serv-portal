package folder

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList generates a comma-separated list of column names for the specified
// folder fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.id")
	}
	if fields&FParent != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.parent")
	}
	if fields&FName != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.name")
	}
	if fields&FURLName != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.url_name")
	}
	if fields&FViewer != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.view_org, f.view_priv")
	}
	if fields&FEditor != 0 {
		sb.WriteString(sep())
		sb.WriteString("f.edit_org, f.edit_priv")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (f *Folder) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		f.id = ID(stmt.ColumnInt())
	}
	if fields&FParent != 0 {
		f.parent = ID(stmt.ColumnInt())
	}
	if fields&FName != 0 {
		f.name = stmt.ColumnText()
	}
	if fields&FURLName != 0 {
		f.urlName = stmt.ColumnText()
	}
	if fields&FViewer != 0 {
		f.viewOrg = enum.Org(stmt.ColumnInt())
		f.viewPriv = enum.PrivLevel(stmt.ColumnInt())
	}
	if fields&FEditor != 0 {
		f.editOrg = enum.Org(stmt.ColumnInt())
		f.editPriv = enum.PrivLevel(stmt.ColumnInt())
	}
	f.fields |= fields
}
