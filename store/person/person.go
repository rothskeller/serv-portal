// Package person defines the Person type, which describes a SERV-related person
// (volunteer, class student, etc.).
package person

import (
	"time"

	"sunnyvaleserv.org/portal/store/enum"
)

// ID uniquely identifies a person.
type ID int

// AdminID is the ID of the admin user, which has special protections.
const AdminID ID = 1

// Flags is a bitmask of flags describing the Person.
type Flags uint

// Values for Flags:
const (
	// NoEmail indicates that the Person should never receive any emails
	// from this system: they have hit the "unsubscribe from all" button.
	NoEmail Flags = 1 << iota
	// NoText indicates that the Person should never receive any SMS
	// messages from this system: they have hit the "unsubscribe from all"
	// button.
	NoText
	// HoursReminder indicates that the Person should receive a reminder (on
	// the 8th of the month) to report their volunteer hours in the previous
	// month.  This flag is set for all eligible volunteers on the 1st of
	// the month, and cleared when they actually record hours.
	HoursReminder
	// VolgisticsPending indicates that the Person has submitted a volunteer
	// application that is outstanding.
	VolgisticsPending
	// VisibleToAll indicates that the Person, and their work contact info,
	// is visible to anyone with a login, bypassing the normal visibility
	// rules.  (This is primarily used for the SERV coordinator.)
	VisibleToAll
)

// Fields is a bitmask of flags identifying specified fields of the Person
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FVolgisticsID
	FInformalName
	FFormalName
	FSortName
	FCallSign
	FPronouns
	FEmail
	FEmail2
	FCellPhone
	FHomePhone
	FWorkPhone
	FPassword
	FBadLoginCount
	FBadLoginTime
	FPWResetToken
	FPWResetTime
	FUnsubscribeToken
	FHoursToken
	FIdentification
	FBirthdate
	FLanguage
	FFlags
	FAddresses
	FBGChecks
	FDSWRegistrations
	FNotes
	FEmContacts
	FPrivLevels
)

// Person describes a SERV-related person (volunteer, class student, etc.).
type Person struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields           Fields // which fields of the structure are populated
	id               ID
	volgisticsID     uint
	informalName     string
	formalName       string
	sortName         string
	callSign         string
	pronouns         string
	email            string
	email2           string
	cellPhone        string
	homePhone        string
	workPhone        string
	password         string
	badLoginCount    uint
	badLoginTime     time.Time
	pwresetToken     string
	pwresetTime      time.Time
	unsubscribeToken string
	hoursToken       string
	birthdate        string
	identification   IdentType
	language         string
	flags            Flags
	addresses        Addresses
	bgChecks         BGChecks
	dswRegistrations DSWRegistrations
	notes            Notes
	emContacts       EmContacts
	privLevels       []enum.PrivLevel
}

func (p *Person) Clone() (c *Person) {
	c = new(Person)
	*c = *p
	c.addresses.Home = p.addresses.Home.clone()
	c.addresses.Work = p.addresses.Work.clone()
	c.addresses.Mail = p.addresses.Mail.clone()
	c.bgChecks.DOJ = p.bgChecks.DOJ.clone()
	c.bgChecks.FBI = p.bgChecks.FBI.clone()
	c.bgChecks.PHS = p.bgChecks.PHS.clone()
	c.dswRegistrations.CERT = p.dswRegistrations.CERT.clone()
	c.dswRegistrations.Communications = p.dswRegistrations.Communications.clone()
	c.notes = p.notes.clone()
	c.emContacts = p.emContacts.clone()
	return c
}

// Address describes an address associated with a Person.
type Address struct {
	// SameAsHome indicates that this (work or mailing) address is the same
	// as the person's home address.  If this flag is set, none of the other
	// fields are filled in.
	SameAsHome bool
	// Address is the full address being specified, with commas between the
	// lines of the address.
	Address string
	// Latitude is the latitude of the address, if it is a physical address.
	// If it is a P.O. box or similar non-physical address, Latitude is
	// zero.
	Latitude float64
	// Longitude is the longitude of the address, if it is a physical
	// address.  If it is a P.O. box or similar non-physical address,
	// Longitude is zero.
	Longitude float64
	// FireDistrict is the Sunnyvale fire district (1 through 6) containing
	// the address, if it is a physical address within Sunnyvale's fire
	// coverage.  Otherwise, it is zero.
	FireDistrict uint
}

func (a *Address) clone() (c *Address) {
	if a == nil {
		return nil
	}
	c = new(Address)
	*c = *a
	return c
}

// Addresses contains the set of addresses for a Person.
type Addresses struct {
	// Home is the person's home (residence) address.
	Home *Address
	// Work is the person's work address.
	Work *Address
	// Mail is the person's postal mailing address.
	Mail *Address
}

func (a Addresses) clone() Addresses {
	return Addresses{a.Home.clone(), a.Work.clone(), a.Mail.clone()}
}

// addressType is the type of an Address when stored in the person_address
// database table.
type addressType uint

// Values for addressType:
const (
	addressHome addressType = iota
	addressWork
	addressMail
)

// BGCheck describes a background check performed on a Person.
type BGCheck struct {
	// Cleared is the date on which the background check was passed.
	// (Background checks that are not passed are not recorded.)  It may be
	// zero if the date is not known.
	Cleared time.Time
	// NLI is the date on which the relevant agency is told that we are "no
	// longer interested" in updates about this person, and therefore the
	// background check should no longer be considered effective.  This only
	// applies to DOJ and FBI checks.  For PHS checks, if this field is
	// nonzero, it means DPS has decided for whatever reason that the person
	// no longer qualifies for unattended access.
	NLI time.Time
	// Assumed is true if the background check is assumed to have occurred
	// but we have no record to prove it.
	Assumed bool
}

// Valid returns whether the background check is valid and current.
func (b *BGCheck) Valid() bool {
	if b == nil {
		return false
	}
	return b.NLI.IsZero() // a future NLI doesn't make sense
}

func (b *BGCheck) clone() (c *BGCheck) {
	if b == nil {
		return nil
	}
	c = new(BGCheck)
	*c = *b
	return c
}

// BGChecks is the set of background checks for a Person.
type BGChecks struct {
	// DOJ is a fingerprint check with the California Department of Justice.
	// It is required of all *city* volunteers, including DPS/OES
	// volunteers.
	DOJ *BGCheck
	// FBI is a fingerprint check with the Federal Bureau of Investigations.
	// It is required of all DPS/OES volunteers.
	FBI *BGCheck
	// PHS is a thorough background check with references based on a
	// Personal History Statement.  It is required for unattended access to
	// DPS facilities (i.e., card key access).
	PHS *BGCheck
}

func (b BGChecks) clone() BGChecks {
	return BGChecks{b.DOJ.clone(), b.FBI.clone(), b.PHS.clone()}
}

// bgCheckType is the type of a BGCheck when stored in the person_bgcheck
// database table.
type bgCheckType uint

// Values for bgCheckType:
const (
	bgCheckDOJ bgCheckType = iota
	bgCheckFBI
	bgCheckPHS
)

// DSWRegistration describes a Person's registration as a Disaster Service
// Worker in a particular classification.
type DSWRegistration struct {
	// Registered is the date the Person registered as a DSW in this
	// classification.
	Registered time.Time
	// Expiration is the date the Person's DSW registration in this
	// classification expired or will expire.  It may be zero if the
	// registration is current with no planned expiration.
	Expiration time.Time
}

// Valid returns whether the DSW registration is valid and current.
func (d *DSWRegistration) Valid() bool {
	if d == nil || d.Registered.IsZero() {
		return false
	}
	if d.Expiration.IsZero() {
		return true
	}
	return d.Expiration.After(time.Now())
}

func (d *DSWRegistration) clone() (c *DSWRegistration) {
	if d == nil {
		return nil
	}
	c = new(DSWRegistration)
	*c = *d
	return c
}

// DSWRegistrations is the set of DSW registrations for a Person.
type DSWRegistrations struct {
	// CERT is the registration in the "Community Emergency Response Team
	// Member" classification.
	CERT *DSWRegistration
	// Communications is the registration in the "Communications"
	// classification.
	Communications *DSWRegistration
}

func (d DSWRegistrations) clone() DSWRegistrations {
	return DSWRegistrations{d.CERT.clone(), d.Communications.clone()}
}

// dswClass is the class of a DSWRegistration when stored in the person_dswreg
// database table.
type dswClass uint

// Values for dswClass:
const (
	dswCommunications dswClass = 2
	dswCERT           dswClass = 3
	// There are actually 14 DSW classifications defined by the state, but
	// these are the two we use.
)

// Note is a dated note associated with a Person.
type Note struct {
	// Note is the text of the note.
	Note string
	// Date is the date of the note.
	Date time.Time
	// Visibility indicates who is allowed to see the note.
	Visibility NoteVisibility
}

func (n *Note) clone() (c *Note) {
	c = new(Note)
	*c = *n
	return c
}

// Notes is the set of notes associated with a Person.
type Notes []*Note

func (n Notes) clone() (c Notes) {
	if len(n) == 0 {
		return nil
	}
	c = make([]*Note, len(n))
	for i := range n {
		c[i] = n[i].clone()
	}
	return c
}

// EmContact is an emergency contact for a Person.
type EmContact struct {
	// Name is the name of the emergency contact.
	Name string
	// HomePhone is the home phone number of the emergency contact.
	HomePhone string
	// CellPhone is the cell phone number of the emergency contact.
	CellPhone string
	// Relationship is the relationship of the emergency contact to the
	// Person.
	Relationship string
}

func (ec *EmContact) clone() (c *EmContact) {
	c = new(EmContact)
	*c = *ec
	return c
}

// EmContacts is the set of emergency contacts associated with a Person.
type EmContacts []*EmContact

func (ec EmContacts) clone() (c EmContacts) {
	if len(ec) == 0 {
		return nil
	}
	c = make([]*EmContact, len(ec))
	for i := range ec {
		c[i] = ec[i].clone()
	}
	return c
}
