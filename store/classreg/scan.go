package classreg

import (
	"strings"

	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
)

// ColumnList generates a comma-separated list of column names for the specified
// class registration fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.id")
	}
	if fields&FClass != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.class")
	}
	if fields&FPerson != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.person")
	}
	if fields&FRegisteredBy != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.registered_by")
	}
	if fields&FFirstName != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.first_name")
	}
	if fields&FLastName != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.last_name")
	}
	if fields&FEmail != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.email")
	}
	if fields&FCellPhone != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.cell_phone")
	}
	if fields&FWaitlist != 0 {
		sb.WriteString(sep())
		sb.WriteString("cr.waitlist")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (cr *ClassReg) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		cr.id = ID(stmt.ColumnInt())
	}
	if fields&FClass != 0 {
		cr.class = class.ID(stmt.ColumnInt())
	}
	if fields&FPerson != 0 {
		cr.person = person.ID(stmt.ColumnInt())
	}
	if fields&FRegisteredBy != 0 {
		cr.registeredBy = person.ID(stmt.ColumnInt())
	}
	if fields&FFirstName != 0 {
		cr.firstName = stmt.ColumnText()
	}
	if fields&FLastName != 0 {
		cr.lastName = stmt.ColumnText()
	}
	if fields&FEmail != 0 {
		cr.email = stmt.ColumnText()
	}
	if fields&FCellPhone != 0 {
		cr.cellPhone = stmt.ColumnText()
	}
	if fields&FWaitlist != 0 {
		cr.waitlist = stmt.ColumnBool()
	}
	cr.fields |= fields
}
