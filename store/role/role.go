// Package role defines the Role type, which describes roles held by people.
package role

import "sunnyvaleserv.org/portal/store/enum"

// ID uniquely identifies a role.
type ID int

// Well known role IDs:
const (
	Webmaster ID = 1
	Disabled  ID = 2
)

// Flags is a bitmask of flags describing a Role.
type Flags uint

// Values for Flags:
const (
	// Filter indicates that the role is useful as a filter for a list of
	// people, i.e., it is held by a group of people and is relevant to
	// current activities.
	Filter Flags = 1 << iota
	// ImplicitOnly indicates that a person cannot be directly assigned the
	// Role; they can only get it indirectly through implication by another
	// Role.
	ImplicitOnly
	// Archived indicates that the Role is no longer relevant to current
	// activities.  Assignments to it are retained and displayed, but it is
	// omitted from other contexts.  This is mostly used for student roles
	// for past classes.
	Archived
)

// Fields is a bitmask of flags identifying specified fields of the Role
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FName
	FTitle
	FPriority
	FOrg
	FPrivLevel
	FFlags
	FImplies
)

// Role describes a role that can be held by people.
type Role struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields    Fields // which fields of the structure are populated
	id        ID
	name      string
	title     string
	priority  uint
	org       enum.Org
	privLevel enum.PrivLevel
	flags     Flags
	implies   []ID
}

// Clone creates a copy of a Role.
func (r *Role) Clone() (r2 *Role) {
	r2 = new(Role)
	*r2 = *r
	if r.implies != nil {
		r2.implies = make([]ID, len(r.implies))
		copy(r2.implies, r.implies)
	}
	return r2
}
