package event

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/venue"
)

// ColumnList generates a comma-separated list of column names for the specified
// Event fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.id")
	}
	if fields&FName != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.name")
	}
	if fields&FStart != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.start")
	}
	if fields&FEnd != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.end")
	}
	if fields&FVenue != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.venue")
	}
	if fields&FVenueURL != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.venue_url")
	}
	if fields&FActivation != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.activation")
	}
	if fields&FDetails != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.details")
	}
	if fields&FFlags != 0 {
		sb.WriteString(sep())
		sb.WriteString("e.flags")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (e *Event) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		e.id = ID(stmt.ColumnInt())
	}
	if fields&FName != 0 {
		e.name = stmt.ColumnText()
	}
	if fields&FStart != 0 {
		e.start = stmt.ColumnText()
	}
	if fields&FEnd != 0 {
		e.end = stmt.ColumnText()
	}
	if fields&FVenue != 0 {
		e.venue = venue.ID(stmt.ColumnInt())
	}
	if fields&FVenueURL != 0 {
		e.venueURL = stmt.ColumnText()
	}
	if fields&FActivation != 0 {
		e.activation = stmt.ColumnText()
	}
	if fields&FDetails != 0 {
		e.details = stmt.ColumnText()
	}
	if fields&FFlags != 0 {
		e.flags = Flag(stmt.ColumnHexInt())
	}
	e.fields |= fields
}
