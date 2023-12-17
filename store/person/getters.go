package person

import (
	"time"

	"sunnyvaleserv.org/portal/store/enum"
)

// Fields returns the set of fields that have been retrieved for this person.
func (p *Person) Fields() Fields {
	return p.fields
}

// ID is the unique identifier of the Person.
func (p *Person) ID() ID {
	if p == nil {
		return 0
	}
	if p.fields&FID == 0 {
		panic("Person.ID called without having fetched FID")
	}
	return p.id
}

// VolgisticsID is ID the of the Person in HR's Volgistics volunteer management
// software, or 0 if they are not registered as a city volunteer.
func (p *Person) VolgisticsID() uint {
	if p.fields&FVolgisticsID == 0 {
		panic("Person.VolgisticsID called without having fetched FVolgisticsID")
	}
	return p.volgisticsID
}

// InformalName is the informal name of the Person, in "First Last" order.  It
// is how we normally address them.
func (p *Person) InformalName() string {
	if p.fields&FInformalName == 0 {
		panic("Person.InformalName called without having fetched FInformalName")
	}
	return p.informalName
}

// FormalName is the formal name of the Person, in "First Last" order.  It is
// how their name would appear on certificates and other formal documents.
func (p *Person) FormalName() string {
	if p.fields&FFormalName == 0 {
		panic("Person.FormalName called without having fetched FFormalName")
	}
	return p.formalName
}

// SortName is the name of the Person as it would appear in a sortable list,
// i.e., the informal name in "Last, First" order.
func (p *Person) SortName() string {
	if p.fields&FSortName == 0 {
		panic("Person.SortName called without having fetched FSortName")
	}
	return p.sortName
}

// CallSign is the amateur radio FCC call sign of the Person, if any.
func (p *Person) CallSign() string {
	if p.fields&FCallSign == 0 {
		panic("Person.CallSign called without having fetched FCallSign")
	}
	return p.callSign
}

// Pronouns are the preferred pronouns of the Person, e.g., "he/him/his".
func (p *Person) Pronouns() string {
	if p.fields&FPronouns == 0 {
		panic("Person.Pronouns called without having fetched FPronouns")
	}
	return p.pronouns
}

// Email is the primary email address of the Person.
func (p *Person) Email() string {
	if p.fields&FEmail == 0 {
		panic("Person.Email called without having fetched FEmail")
	}
	return p.email
}

// Email2 is the secondary email address of the Person, if any.
func (p *Person) Email2() string {
	if p.fields&FEmail2 == 0 {
		panic("Person.Email2 called without having fetched FEmail2")
	}
	return p.email2
}

// CellPhone is the cell phone number of the Person (10 digits, no punctuation).
func (p *Person) CellPhone() string {
	if p.fields&FCellPhone == 0 {
		panic("Person.CellPhone called without having fetched FCellPhone")
	}
	return p.cellPhone
}

// HomePhone is the home phone number of the Person (10 digits, no punctuation).
func (p *Person) HomePhone() string {
	if p.fields&FHomePhone == 0 {
		panic("Person.HomePhone called without having fetched FHomePhone")
	}
	return p.homePhone
}

// WorkPhone is the work phone number of the Person (10 or more digits, no
// punctuation; any digits beyond 10 are assumed to be an extension number).
func (p *Person) WorkPhone() string {
	if p.fields&FWorkPhone == 0 {
		panic("Person.WorkPhone called without having fetched FWorkPhone")
	}
	return p.workPhone
}

// Password is the hashed password of the Person.
func (p *Person) Password() string {
	if p.fields&FPassword == 0 {
		panic("Person.Password called without having fetched FPassword")
	}
	return p.password
}

// BadLoginCount is number of failed login attempts the Person has had since
// their last successful login.
func (p *Person) BadLoginCount() uint {
	if p.fields&FBadLoginCount == 0 {
		panic("Person.BadLoginCount called without having fetched FBadLoginCount")
	}
	return p.badLoginCount
}

// BadLoginTime is the time at which the Person's last failed login attempt
// occurred.  It may be the zero time indicating they have never had one.
func (p *Person) BadLoginTime() time.Time {
	if p.fields&FBadLoginTime == 0 {
		panic("Person.BadLoginTime called without having fetched FBadLoginTime")
	}
	return p.badLoginTime
}

// PWResetToken is the token for the Person's in-progress password reset
// attempt, if any.
func (p *Person) PWResetToken() string {
	if p.fields&FPWResetToken == 0 {
		panic("Person.PWResetToken called without having fetched FPWResetToken")
	}
	return p.pwresetToken
}

// PWResetTime is the time at which the Person's in-progress password reset
// attempt was started.  It may be the zero time indicating that they have none.
func (p *Person) PWResetTime() time.Time {
	if p.fields&FPWResetTime == 0 {
		panic("Person.PWResetTime called without having fetched FPWResetTime")
	}
	return p.pwresetTime
}

// UnsubscribeToken is the authentication token placed in the footer of list
// emails sent to the Person, allowing them to unsubscribe from the list without
// logging in.
func (p *Person) UnsubscribeToken() string {
	if p.fields&FUnsubscribeToken == 0 {
		panic("Person.UnsubscribeToken called without having fetched FUnsubscribeToken")
	}
	return p.unsubscribeToken
}

// HoursToken is the authentication token placed in emails sent to the Person
// requesting volunteer hours, allowing them to log hours without logging in.
func (p *Person) HoursToken() string {
	if p.fields&FHoursToken == 0 {
		panic("Person.HoursToken called without having fetched FHoursToken")
	}
	return p.hoursToken
}

// Identification is a bitmask of identifications that have been issued to the
// Person, such as logo shirts, photo IDs, etc.
func (p *Person) Identification() IdentType {
	if p.fields&FIdentification == 0 {
		panic("Person.Identification called without having fetched FIdentification")
	}
	return p.identification
}

// Birthdate is the person's birthdate, in YYYY-MM-DD format.
func (p *Person) Birthdate() string {
	if p.fields&FBirthdate == 0 {
		panic("Person.Birthdate called without having fetched FBirthdate")
	}
	return p.birthdate
}

// Flags is a set of flags describing the Person.
func (p *Person) Flags() Flags {
	if p.fields&FFlags == 0 {
		panic("Person.Flags called without having fetched FFlags")
	}
	return p.flags
}

// Addresses is the set of addresses for the Person.  It is a slice indexed by
// AddressType.  It always has length numAddressTypes, but any element of it
// could be nil, indicating no address of that type.
func (p *Person) Addresses() Addresses {
	if p.fields&FAddresses == 0 {
		panic("Person.Addresses called without having fetched FAddresses")
	}
	return p.addresses
}

// BGChecks is the set of background checks of the Person.  It is a slice
// indexed by BGCheckType.  It always has length numBGCheckTypes, but any
// element of it could be nil, indicating no background check of that type.
func (p *Person) BGChecks() BGChecks {
	if p.fields&FBGChecks == 0 {
		panic("Person.BGChecks called without having fetched FBGChecks")
	}
	return p.bgChecks
}

// DSWRegistrations is the set of DSW registrations of the Person.
func (p *Person) DSWRegistrations() DSWRegistrations {
	if p.fields&FDSWRegistrations == 0 {
		panic("Person.DSWRegistrations called without having fetched FDSWRegistrations")
	}
	return p.dswRegistrations
}

// DSWRegistrationForOrg returns the DSW registration appropriate for the
// specified Org.  It returns nil, false if the org doesn't have an associated
// DSW classification.
func (p *Person) DSWRegistrationForOrg(org enum.Org) (*DSWRegistration, bool) {
	switch org {
	case enum.OrgSARES:
		return p.DSWRegistrations().Communications, true
	case enum.OrgCERTD, enum.OrgCERTT:
		return p.DSWRegistrations().CERT, true
	default:
		return nil, false
	}
}

// Notes is the list of dated notes about the Person, in chronological order.
func (p *Person) Notes() []*Note {
	if p.fields&FNotes == 0 {
		panic("Person.Notes called without having fetched FNotes")
	}
	return p.notes
}

// EmContacts is the list of emergency contacts for the Person.
func (p *Person) EmContacts() []*EmContact {
	if p.fields&FEmContacts == 0 {
		panic("Person.EmContacts called without having fetched FEmContacts")
	}
	return p.emContacts
}

// PrivLevels is the set of privilege levels for the Person.  It is a slice
// indexed by enum.Org.  It always has length enum.NumOrgs.
func (p *Person) PrivLevels() []enum.PrivLevel {
	if p.fields&FPrivLevels == 0 {
		panic("Person.PrivLevels called without having fetched FPrivLevels")
	}
	return p.privLevels
}

// HasPrivLevel returns whether the Person has the specified privilege level, or
// any higher privilege level, on the specified organization (or, if org is
// zero, any organization).
func (p *Person) HasPrivLevel(org enum.Org, level enum.PrivLevel) bool {
	if level == 0 {
		return true
	}
	if p == nil {
		return false
	}
	if p.fields&FPrivLevels == 0 {
		panic("Person.HasPrivLevel called without having fetched FPrivLevels")
	}
	if org != 0 {
		return p.privLevels[org] >= level
	}
	for _, org := range enum.AllOrgs {
		if p.privLevels[org] >= level {
			return true
		}
	}
	return false
}

// IsAdminLeader returns whether the Person has PrivLeader on OrgAdmin.
func (p *Person) IsAdminLeader() bool { return p.HasPrivLevel(enum.OrgAdmin, enum.PrivLeader) }

// IsWebmaster returns whether the Person has PrivMaster on OrgAdmin.
func (p *Person) IsWebmaster() bool { return p.HasPrivLevel(enum.OrgAdmin, enum.PrivMaster) }
