// Package classreg defines the ClassReg type, which describes a single
// registration for a class.
package classreg

import (
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/person"
)

// ID uniquely identifies a class registration.
type ID int

// Fields is a bitmask of flags identifying specified fields of the ClassReg
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FClass
	FPerson
	FRegisteredBy
	FFirstName
	FLastName
	FEmail
	FCellPhone
	FWaitlist
)

// ClassReg describes a registration for a class.
type ClassReg struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields       Fields // which fields of the structure are populated
	id           ID
	class        class.ID
	person       person.ID
	registeredBy person.ID
	firstName    string
	lastName     string
	email        string
	cellPhone    string
	waitlist     bool
}

// Clone creates a clone of the class registration.
func (cr *ClassReg) Clone() (clone *ClassReg) {
	clone = new(ClassReg)
	*clone = *cr
	return clone
}
