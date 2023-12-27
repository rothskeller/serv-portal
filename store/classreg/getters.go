package classreg

import (
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/person"
)

// Fields returns the set of fields that have been retrieved for this venue.
func (c *ClassReg) Fields() Fields {
	return c.fields
}

// ID is the unique identifier of the ClassReg.
func (c *ClassReg) ID() ID {
	if c == nil {
		return 0
	}
	if c.fields&FID == 0 {
		panic("ClassReg.ID called without having fetched FID")
	}
	return c.id
}

// Class is the class instance for which the person was registered.
func (c *ClassReg) Class() class.ID {
	if c.fields&FClass == 0 {
		panic("ClassReg.Class called without having fetched FClass")
	}
	return c.class
}

// Person is the ID of the person who was registered, if they have an entry in
// the person table.  (Many registrants don't, in which case this is zero.)
func (c *ClassReg) Person() person.ID {
	if c.fields&FPerson == 0 {
		panic("ClassReg.Person called without having fetched FPerson")
	}
	return c.person
}

// RegisteredBy is the ID of the person who made the registration, which may not
// be the same as the person who was registered.  Registrations must be made by
// someone logged in, so this is never zero.
func (c *ClassReg) RegisteredBy() person.ID {
	if c.fields&FRegisteredBy == 0 {
		panic("ClassReg.RegisteredBy called without having fetched FRegisteredBy")
	}
	return c.registeredBy
}

// FirstName is the first name of the person being registered.
func (c *ClassReg) FirstName() string {
	if c.fields&FFirstName == 0 {
		panic("ClassReg.FirstName called without having fetched FFirstName")
	}
	return c.firstName
}

// LastName is the last name of the person being registered.
func (c *ClassReg) LastName() string {
	if c.fields&FLastName == 0 {
		panic("ClassReg.LastName called without having fetched FLastName")
	}
	return c.lastName
}

// Email is the email address of the person being registered.  It is optional.
func (c *ClassReg) Email() string {
	if c.fields&FEmail == 0 {
		panic("ClassReg.Email called without having fetched FEmail")
	}
	return c.email
}

// CellPhone is the cell phone number of the person being registered.  It is
// optional.
func (c *ClassReg) CellPhone() string {
	if c.fields&FCellPhone == 0 {
		panic("ClassReg.CellPhone called without having fetched FCellPhone")
	}
	return c.cellPhone
}
