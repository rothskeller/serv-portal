package shift

import (
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
)

// Fields returns the set of fields that have been retrieved for this Shift.
func (s *Shift) Fields() Fields {
	return s.fields
}

// ID is the unique identifier of the Shift.
func (s *Shift) ID() ID {
	if s == nil {
		return 0
	}
	if s.fields&FID == 0 {
		panic("Shift.ID called without having fetched FID")
	}
	return s.id
}

// Task is the ID of the Task to which the Shift belongs.
func (s *Shift) Task() task.ID {
	if s.fields&FTask == 0 {
		panic("Shift.Task called without having fetched FTask")
	}
	return s.task
}

// Start is the Shift's start time.
func (s *Shift) Start() string {
	if s.fields&FStart == 0 {
		panic("Shift.Start called without having fetched FStart")
	}
	return s.start
}

// End is the Shift's end time.
func (s *Shift) End() string {
	if s.fields&FEnd == 0 {
		panic("Shift.End called without having fetched FEnd")
	}
	return s.end
}

// Venue is the ID of the Shift's venue, or 0 if the shift has no venue.
func (s *Shift) Venue() venue.ID {
	if s.fields&FVenue == 0 {
		panic("Shift.Venue called without having fetched FVenue")
	}
	return s.venue
}

// Min is the minimum number of people that we need for the Shift.
func (s *Shift) Min() uint {
	if s.fields&FMin == 0 {
		panic("Shift.Min called without having fetched FMin")
	}
	return s.min
}

// Max is the maximum number of people who can be signed up for the Shift.  If
// it is zero, the number of signups is unlimited.
func (s *Shift) Max() uint {
	if s.fields&FMax == 0 {
		panic("Shift.Max called without having fetched FMax")
	}
	return s.max
}
