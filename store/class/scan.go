package class

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ColumnList generates a comma-separated list of column names for the specified
// class fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.id")
	}
	if fields&FType != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.type")
	}
	if fields&FStart != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.start")
	}
	if fields&FEnDesc != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.en_desc")
	}
	if fields&FEsDesc != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.es_desc")
	}
	if fields&FLimit != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.elimit")
	}
	if fields&FReferrals != 0 {
		sb.WriteString(sep())
		sb.WriteString("c.referrals")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (c *Class) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		c.id = ID(stmt.ColumnInt())
	}
	if fields&FType != 0 {
		c.ctype = Type(stmt.ColumnInt())
	}
	if fields&FStart != 0 {
		c.start = stmt.ColumnText()
	}
	if fields&FEnDesc != 0 {
		c.enDesc = stmt.ColumnText()
	}
	if fields&FEsDesc != 0 {
		c.esDesc = stmt.ColumnText()
	}
	if fields&FLimit != 0 {
		c.limit = uint(stmt.ColumnInt())
	}
	if fields&FReferrals != 0 {
		refmask := uint64(stmt.ColumnInt())
		c.referrals = make([]uint, len(AllReferrals)+1)
		for _, ref := range AllReferrals {
			c.referrals[ref] = uint((refmask >> (ref * 8)) & 0xFF)
		}
	}
	c.fields |= fields
}
