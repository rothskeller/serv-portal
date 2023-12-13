package event

import (
	"sunnyvaleserv.org/portal/store/venue"
)

// Fields returns the set of fields that have been retrieved for this event.
func (e *Event) Fields() Fields {
	return e.fields
}

// ID is the unique identifier of the Event.
func (e *Event) ID() ID {
	if e == nil {
		return 0
	}
	if e.fields&FID == 0 {
		panic("Event.ID called without having fetched FID")
	}
	return e.id
}

// Name is the Event's name.
func (e *Event) Name() string {
	if e.fields&FName == 0 {
		panic("Event.Name called without having fetched FName")
	}
	return e.name
}

// Start is the Event's starting time, in YYYY-MM-DDTHH:MM format (local time).
func (e *Event) Start() string {
	if e.fields&FStart == 0 {
		panic("Event.Start called without having fetched FStart")
	}
	return e.start
}

// End is the Event's ending time, in YYYY-MM-DDTHH:MM format (local time).
func (e *Event) End() string {
	if e.fields&FEnd == 0 {
		panic("Event.End called without having fetched FEnd")
	}
	return e.end
}

// Venue is the Event's venue ID (which may be zero for no venue).
func (e *Event) Venue() venue.ID {
	if e.fields&FVenue == 0 {
		panic("Event.Venue called without having fetched FVenue")
	}
	return e.venue
}

// VenueURL is the URL for the Event's venue, overriding any URL specified in
// the venue object.
func (e *Event) VenueURL() string {
	if e.fields&FVenueURL == 0 {
		panic("Event.VenueURL called without having fetched FVenueURL")
	}
	return e.venueURL
}

// Activation is the Event's activation number (which may be empty).
func (e *Event) Activation() string {
	if e.fields&FActivation == 0 {
		panic("Event.Activation called without having fetched FActivation")
	}
	return e.activation
}

// Details is free-form HTML text describing the Event.
func (e *Event) Details() string {
	if e.fields&FDetails == 0 {
		panic("Event.Details called without having fetched FDetails")
	}
	return e.details
}

// Flags are flags describing the Event.
func (e *Event) Flags() Flag {
	if e.fields&FFlags == 0 {
		panic("Event.Flags called without having fetched FFlags")
	}
	return e.flags
}
