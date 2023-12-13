package enum

import "errors"

// PrivLevel is a privilege level that a person may hold in an organization, by
// virtue of the roles held by that person.
type PrivLevel uint

// Values for PrivLevel.
const (
	_ PrivLevel = iota
	// PrivStudent indicates student-level membership in an organization,
	// with essentially no privileges other than being on lists.
	PrivStudent
	// PrivMember indicates full membership in an organization, but not
	// leadership privileges.
	PrivMember
	// PrivLeader indicates leadership of an organization.
	PrivLeader
	// PrivMaster indicates a webmaster.
	PrivMaster
)

// String returns the string form of the PrivLevel, as used in APIs.
func (pl PrivLevel) String() string {
	switch pl {
	case PrivStudent:
		return "Student"
	case PrivMember:
		return "Member"
	case PrivLeader:
		return "Leader"
	case PrivMaster:
		return "Master"
	default:
		return ""
	}
}

// ParsePrivLevel translates the string form of a PrivLevel (as returned by
// String) into a PrivLevel value.  It returns an error if the string is not
// recognized.
func ParsePrivLevel(s string) (PrivLevel, error) {
	switch s {
	case "Student":
		return PrivStudent, nil
	case "Member":
		return PrivMember, nil
	case "Leader":
		return PrivLeader, nil
	case "Master":
		return PrivMaster, nil
	default:
		return 0, errors.New("invalid privLevel")
	}
}

// Valid returns whether an PrivLevel value is valid.
func (pl PrivLevel) Valid() bool {
	return pl == PrivStudent || pl == PrivMember || pl == PrivLeader || pl == PrivMaster
}

// AllPrivLevels is a list of all privilege levels.
var AllPrivLevels = []PrivLevel{PrivStudent, PrivMember, PrivLeader, PrivMaster}
