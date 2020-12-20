package model

import (
	"errors"
	strings "strings"
)

// A BGCheckType is a bitmask of types of background check that a person has
// had.
type BGCheckType uint8

// Values for BGCheckType.
const (
	// BGCheckDOJ means the person has been fingerprinted and has had a
	// criminal records check through the California Department of Justice.
	BGCheckDOJ BGCheckType = 1 << iota
	// BGCheckFBI means the person has been fingerprinted and has had a
	// criminal records check through the Federal Bureau of Investigations.
	BGCheckFBI
	// BGCheckPHS means the person has submitted a personal history
	// statement with references, and those references have been
	// interviewed.
	BGCheckPHS
)

// String returns the string form of a BGCheckType value, as used in APIs.  Note
// that this can only be successfully called on a single bit, not a bitmask.
func (v BGCheckType) String() string {
	switch v {
	case BGCheckDOJ:
		return "DOJ"
	case BGCheckFBI:
		return "FBI"
	case BGCheckPHS:
		return "PHS"
	default:
		return ""
	}
}

// MaskString returns the string form of a BGCheckType bitmask.
func (v BGCheckType) MaskString() string {
	var s []string
	for _, t := range AllBGCheckTypes {
		if v&t != 0 {
			s = append(s, t.String())
		}
	}
	return strings.Join(s, "+")
}

// ParseBGCheckType translates the string form of a BGCheckType (as returned by
// String(), above) into a BGCheckType value.  It returns an error if the string
// is not recognized.
func ParseBGCheckType(s string) (BGCheckType, error) {
	switch s {
	case "DOJ":
		return BGCheckDOJ, nil
	case "FBI":
		return BGCheckFBI, nil
	case "PHS":
		return BGCheckPHS, nil
	default:
		return BGCheckType(0), errors.New("invalid bgCheckType")
	}
}

// Valid returns whether a BGCheckType value is valid as a single bit
// (mask=false) or as a bitmask (mask=true).
func (v BGCheckType) Valid(mask bool) bool {
	if mask {
		return v&^(BGCheckDOJ|BGCheckFBI|BGCheckPHS) == 0
	}
	switch v {
	case BGCheckDOJ, BGCheckFBI, BGCheckPHS:
		return true
	default:
		return false
	}
}

// AllBGCheckTypes is the ordered list of BGCheckType values, which is
// used for iteration.
var AllBGCheckTypes = []BGCheckType{BGCheckDOJ, BGCheckFBI, BGCheckPHS}
