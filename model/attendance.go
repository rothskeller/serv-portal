package model

import (
	"errors"
)

// An AttendanceInfo structure gives information about a person's attendance at
// an event, and volunteer hours for the event.
type AttendanceInfo struct {
	Type    AttendanceType
	Minutes uint16
}

// An AttendanceType indicates the role that a person played in attending an
// event (though it's carefully not called "role" to avoid confusion with
// authorization roles).
type AttendanceType uint8

// Values for AttendanceType.
const (
	AttendAsVolunteer AttendanceType = iota
	AttendAsStudent
	AttendAsAuditor
	// AttendAsAbsent the attendance type used in an AttendanceInfo record
	// when a person has reported hours for the event, but wasn't actually
	// recorded as being in attendance.
	AttendAsAbsent
)

// String returns the string form of the AttendanceType, as used in APIs.  It
// is not the display label.
func (v AttendanceType) String() string {
	switch v {
	case AttendAsVolunteer:
		return "Volunteer"
	case AttendAsStudent:
		return "Student"
	case AttendAsAuditor:
		return "Audit"
	case AttendAsAbsent:
		return "Absent"
	default:
		return ""
	}
}

// ParseAttendanceType translates the string form of a AttendanceType (as
// returned by String(), above) into a AttendanceType value.  It returns an
// error if the string is not recognized.
func ParseAttendanceType(s string) (AttendanceType, error) {
	switch s {
	case "Volunteer":
		return AttendAsVolunteer, nil
	case "Student":
		return AttendAsStudent, nil
	case "Audit":
		return AttendAsAuditor, nil
	case "Absent":
		return AttendAsAbsent, nil
	default:
		return AttendanceType(0), errors.New("invalid visibility")
	}
}

// Valid returns whether an AttendanceType value is valid.
func (v AttendanceType) Valid() bool {
	switch v {
	case AttendAsVolunteer, AttendAsStudent, AttendAsAuditor, AttendAsAbsent:
		return true
	default:
		return false
	}
}

// AllAttendanceTypes is the ordered list of AttendanceType values, which is
// used for iteration.
var AllAttendanceTypes = []AttendanceType{AttendAsVolunteer, AttendAsStudent, AttendAsAuditor, AttendAsAbsent}
