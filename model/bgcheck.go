package model

import (
	"errors"
	strings "strings"
)

// A BGCheckStatus indicates the status of a person's background check.
type BGCheckStatus uint8

// Values for BGCheckStatus.
const (
	// BGCheckNone means the person has not passed a background check.
	BGCheckNone BGCheckStatus = iota
	// BGCheckAssumed means that we assume the person had a background
	// check, but we have no record of it.
	BGCheckAssumed
	// BGCheckRecorded means that the person's background check is recorded.
	BGCheckRecorded
)

// String returns the string form of a BGCheckStatus value, as used in APIs.
func (v BGCheckStatus) String() string {
	switch v {
	case BGCheckAssumed:
		return "assumed"
	case BGCheckRecorded:
		return "recorded"
	default:
		return ""
	}
}

// ParseBGCheckStatus translates the string form of a BGCheckStatus (as returned
// by String(), above) into a BGCheckStatus value.  It returns an error if the
// string is not recognized.
func ParseBGCheckStatus(s string) (BGCheckStatus, error) {
	switch s {
	case "":
		return BGCheckNone, nil
	case "assumed":
		return BGCheckAssumed, nil
	case "recorded":
		return BGCheckRecorded, nil
	default:
		return BGCheckStatus(0), errors.New("invalid bgCheckStatus")
	}
}

// Valid returns whether a BGCheckStatus value is valid.
func (v BGCheckStatus) Valid() bool {
	switch v {
	case BGCheckNone, BGCheckAssumed, BGCheckRecorded:
		return true
	default:
		return false
	}
}

// AllBGCheckStatuses is the ordered list of BGCheckStatus values, which is
// used for iteration.
var AllBGCheckStatuses = []BGCheckStatus{BGCheckNone, BGCheckAssumed, BGCheckRecorded}

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
