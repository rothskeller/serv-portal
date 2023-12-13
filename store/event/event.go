// Package evnt defines the data model for SERV events.
package event

import (
	"sunnyvaleserv.org/portal/store/venue"
)

// ID uniquely identifies an event.
type ID int

// Flag is a flag, or a bitmask of flags, describing an event.
type Flag uint

// Values for Flags:
const (
	// OtherHours is a flag indicating that this "event" is a pseudo-event
	// used for recording volunteer hours not associated with a real event.
	OtherHours Flag = 1 << iota
)

// Fields is a bitmask of flags identifying specified fields of the Event
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FName
	FStart
	FEnd
	FVenue
	FVenueURL
	FActivation
	FDetails
	FFlags
)

// Event describes a single event on the SERV calendar.
type Event struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields     Fields // which fields of the structure are populated
	id         ID
	name       string
	start      string
	end        string
	venue      venue.ID
	venueURL   string
	activation string
	details    string
	flags      Flag
}

// Clone returns a clone of the receiver Event.
func (e *Event) Clone() (c *Event) {
	c = new(Event)
	*c = *e
	return c
}
