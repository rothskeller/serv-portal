// Package shift defines the data model for SERV shifts.
package shift

import (
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
)

// ID uniquely identifies a shift.
type ID int

// Fields is a bitmask of flags identifying specified fields of the Shift
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FTask
	FStart
	FEnd
	FVenue
	FMin
	FMax
)

// Shift describes a single task in an event on the SERV calendar.
type Shift struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields Fields // which fields of the structure are populated
	id     ID
	task   task.ID
	start  string
	end    string
	venue  venue.ID
	min    uint
	max    uint
}

// Clone creates a copy of a Shift.
func (s *Shift) Clone() (c *Shift) {
	c = new(Shift)
	*c = *s
	return c
}
