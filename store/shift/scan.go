package shift

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
)

// ColumnList generates a comma-separated list of column names for the specified
// Shift fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.id")
	}
	if fields&FTask != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.task")
	}
	if fields&FStart != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.start")
	}
	if fields&FEnd != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.end")
	}
	if fields&FVenue != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.venue")
	}
	if fields&FMin != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.min")
	}
	if fields&FMax != 0 {
		sb.WriteString(sep())
		sb.WriteString("s.max")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (s *Shift) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		s.id = ID(stmt.ColumnInt())
	}
	if fields&FTask != 0 {
		s.task = task.ID(stmt.ColumnInt())
	}
	if fields&FStart != 0 {
		s.start = stmt.ColumnText()
	}
	if fields&FEnd != 0 {
		s.end = stmt.ColumnText()
	}
	if fields&FVenue != 0 {
		s.venue = venue.ID(stmt.ColumnInt())
	}
	if fields&FMin != 0 {
		s.min = uint(stmt.ColumnInt())
	}
	if fields&FMax != 0 {
		s.max = uint(stmt.ColumnInt())
	}
	s.fields |= fields
}
