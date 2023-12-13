// Package venue defines the Venue type, which describes an event venue.
package venue

// ID uniquely identifies a venue.
type ID int

// Flag is a flag, or bitmask of flags, for a venue.
type Flag uint

// Values for Flag
const (
	// CanOverlap indicates that multiple events can occur simultaneously
	// at this venue.  Typically used for online meetings.
	CanOverlap Flag = 1 << iota
)

// Fields is a bitmask of flags identifying specified fields of the Venue
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FName
	FURL
	FFlags
)

// Venue describes a venue at which events can be held.
type Venue struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields Fields // which fields of the structure are populated
	id     ID
	name   string
	url    string
	flags  Flag
}

// Clone creates a clone of the venue.
func (v *Venue) Clone() (c *Venue) {
	if v == nil {
		return nil
	}
	c = new(Venue)
	*c = *v
	return c
}
