// Package class defines the Class type, which describes an instance of a class
// that we offer.
package class

import "slices"

// ID uniquely identifies a class.
type ID int

// Fields is a bitmask of flags identifying specified fields of the Class
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FType
	FStart
	FEnDesc
	FEsDesc
	FLimit
	FReferrals
)

// Class describes an instance of a class that we offer.
type Class struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields    Fields // which fields of the structure are populated
	id        ID
	ctype     Type
	start     string
	enDesc    string
	esDesc    string
	limit     uint
	referrals []uint
}

// Clone creates a clone of the class.
func (c *Class) Clone() (clone *Class) {
	clone = new(Class)
	*clone = *c
	clone.referrals = slices.Clone(c.referrals)
	return clone
}
