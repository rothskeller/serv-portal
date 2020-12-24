package model

import (
	"errors"
)

// A PrivLevel is a privilege level for membership in an Org.
type PrivLevel uint8

// Values for PrivLevel.
const (
	// PrivNone indicates no membership in the Org.
	PrivNone PrivLevel = iota
	// PrivStudent indicates student-level membership in the Org, with
	// essentially no privileges other than being on lists.
	PrivStudent
	// PrivMember indicates full membership in the Org.
	PrivMember
	// PrivLeader indicates a leader of the Org.
	PrivLeader
)

// String returns the string form of the PrivLevel, as used in APIs.
func (pl PrivLevel) String() string {
	switch pl {
	case PrivStudent:
		return "student"
	case PrivMember:
		return "member"
	case PrivLeader:
		return "leader"
	default:
		return ""
	}
}

// ParsePrivLevel translates the string form of a PrivLevel (as returned by
// String(), above) into a PrivLevel value.  It returns an error if the string
// is not recognized.
func ParsePrivLevel(s string) (PrivLevel, error) {
	switch s {
	case "student":
		return PrivStudent, nil
	case "member":
		return PrivMember, nil
	case "leader":
		return PrivLeader, nil
	default:
		return PrivNone, errors.New("invalid privLevel")
	}
}

// Valid returns whether an PrivLevel value is valid.
func (pl PrivLevel) Valid() bool {
	return pl == PrivStudent || pl == PrivMember || pl == PrivLeader
}

// AllPrivLevels is a list of all privilege levels.
var AllPrivLevels = []PrivLevel{PrivStudent, PrivMember, PrivLeader}
