// Package model contains the data model types and constants for the SERV
// portal.
package model

//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. model.proto
//go:generate easyjson -all model.pb.go

import (
	"errors"
	"time"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// An ApprovalID identifies an item waiting for approval.
type ApprovalID int

// Equal compares two background check records for equality.
func (bc *BackgroundCheck) Equal(o *BackgroundCheck) bool {
	return bc.Date == o.Date && bc.Type == o.Type && bc.Assumed == o.Assumed
}

// A CSRFToken is a random string used to verify that submitted forms came from
// our own site and not from a forgery.
type CSRFToken string

// A Document represents a document (either a file or a link) in a folder.
type Document struct {
	// Name is the name of the document.  For files, it is the filename.
	// For links, it is the link title.
	Name string
	// URL is the URL for link documents.  For files, it is empty.
	URL string
}

// A DSWClass identifies a DSW classification as defined in state regulations
// (19 CCR § 2572.1).
type DSWClass int

const (
	// DSWComm is the "Communications" classification.
	DSWComm DSWClass = 2
	// DSWCERT is the "Community Emergency Response Team Member"
	// classification.
	DSWCERT = 3
	// Other state-defined classifications are not used by SERV.
)

// AllDSWClasses gives the list of all DSW classes.
var AllDSWClasses = []DSWClass{DSWComm, DSWCERT}

// DSWClassNames gives the names of the various DSW classes.
var DSWClassNames = map[DSWClass]string{
	DSWComm: "Communications",
	DSWCERT: "CERT",
}

// MarshalEasyJSON encodes the DSWClass into JSON.
func (c DSWClass) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(DSWClassNames[c])
}

// UnmarshalEasyJSON decodes the DSWClass from JSON.
func (c *DSWClass) UnmarshalEasyJSON(l *jlexer.Lexer) {
	s := l.UnsafeString()
	for org, name := range DSWClassNames {
		if s == name {
			*c = org
			return
		}
	}
	l.AddError(errors.New("unrecognized value for DSWClass"))
}

// EmContactRelationships is a list of allowed values for the Relationship
// field of an EmContact.
var EmContactRelationships = []string{
	"Co-worker", "Daughter", "Father", "Friend", "Mother", "Neighbor",
	"Other", "Relative", "Son", "Spouse", "Supervisor",
}

// An EventID is a positive integer uniquely identifying an Event.
type EventID int

// An EventType identifies the type of an Event.  (This used to be a bitmask of
// multiple types, but now each event is restricted to a single type.  I kept
// the values unchanged for convenience.)
type EventType uint32

// The known event types.
const (
	// EventPublicService represents a public service event, such as the
	// State of the City, Pancake Breakfast, etc.
	EventPublicService EventType = 1 << iota
	// EventClass represents a class that is taught by SERV or DPS
	// instructors.
	EventClass
	eventUnused1 // no longer used
	eventUnused2 // no longer used
	EventEmergency
	EventMeeting
	EventSocial
	EventTraining
	eventUnused3 // no longer used
	// EventHours is a placeholder event to record "other" volunteer hours
	// not associated with a true event.  Events of this type are never
	// visible.
	EventHours
)

// AllEventTypes is the list of all known event types.
var AllEventTypes = []EventType{
	EventClass,
	EventEmergency,
	EventMeeting,
	EventPublicService,
	EventSocial,
	EventTraining,
	EventHours,
}

// EventTypeNames maps event types to their names.
var EventTypeNames = map[EventType]string{
	EventClass:         "Class",
	EventEmergency:     "Emergency",
	EventMeeting:       "Meeting",
	EventPublicService: "Public Service",
	EventSocial:        "Social",
	EventTraining:      "Training",
	EventHours:         "Other Hours",
}

// Hours returns the number of hours the event lasted.
func (e *Event) Hours() float64 {
	start, _ := time.Parse("15:04", e.Start)
	end, _ := time.Parse("15:04", e.End)
	return float64(end.Sub(start)) / float64(time.Hour)
}

// A Folder represents a folder of documents in the site's repository.
type Folder struct {
	// Name is the display name of the folder.
	Name string
	// URL is the URL of the folder, relative to the root of the folder
	// tree.  It should start with a slash unless it is empty.
	URL string
	// Visibility determines who can see the folder.
	Visibility FolderVisibility
	// Org determines which organization can see the folder when Visibility
	// is set to FolderVisibleToOrg.
	Org Org
}

// An IdentType is a type of identification (or a bitmask of multiple types)
// that has been issued to a volunteer (and should be retrieved if they leave).
type IdentType uint8

// Values for IdentType.
const (
	// IDPhoto is a regular DPS Volunteer photo ID badge.
	IDPhoto IdentType = 1 << iota
	// IDCardKey is a photo badge with card key access to DPS buildings.
	IDCardKey
	// IDSERVShirt is a tan long-sleeved button-down shirt identifying the
	// person as a SERV leader and/or class instructor.
	IDSERVShirt
	// IDCERTShirt is a green long-sleeved tee shirt identifying the person
	// as a member of a CERT deployment team.
	IDCERTShirtLS
	// IDCERTShirt is a green short-sleeved tee shirt identifying the person
	// as a member of a CERT deployment team.
	IDCERTShirtSS
)

// AllIdentTypes is the list of all identification types.
var AllIdentTypes = []IdentType{IDPhoto, IDCardKey, IDSERVShirt, IDCERTShirtLS, IDCERTShirtSS}

// IdentTypeNames gives the names for each identification type.
var IdentTypeNames = map[IdentType]string{
	IDPhoto:       "photo ID",
	IDCardKey:     "access card",
	IDSERVShirt:   "tan SERV shirt",
	IDCERTShirtLS: "green CERT shirt (LS)",
	IDCERTShirtSS: "green CERT shirt (SS)",
}

// A ListID identifies a List.
type ListID int

// ListPersonStatus is a bitmask of flags describing a person's status on a
// list.
type ListPersonStatus uint8

// Values for ListPersonStatus.
const (
	// ListSubscribed indicates that the person is subscribed to the list.
	ListSubscribed ListPersonStatus = 1 << iota
	// ListUnsubscribed indicates that the person has unsubscribed from the
	// list.
	ListUnsubscribed
	// ListSender indicates that the person is allowed to send to the list.
	ListSender
)

// ListSubModel describes the subscription model a given role grants to a given
// list.
type ListSubModel uint8

// Values for ListSubModel.
const (
	// ListNoSub indicates that holders of the role are not granted any
	// subscription privileges on the list.
	ListNoSub ListSubModel = iota
	// ListAllowSub indicates that holders of the role are allowed to
	// subscribe to the list.
	ListAllowSub
	// ListAutoSub indicates that holders of the role are automatically
	// subscribed to the list.
	ListAutoSub
	// ListWarnUnsub is like ListAutoSub, but people trying to unsubscribe
	// from the list are warned that they may lose the role if they do.
	ListWarnUnsub
)

// AllListSubModels is a list of all list subscription models values.
var AllListSubModels = []ListSubModel{ListAllowSub, ListAutoSub, ListWarnUnsub}

// ListSubModelNames gives the display names of the list subscription models.
var ListSubModelNames = map[ListSubModel]string{
	ListAllowSub:  "allow",
	ListAutoSub:   "auto",
	ListWarnUnsub: "warn",
}

// ListType is the type of a list.
type ListType uint8

// Values for ListType.
const (
	// ListNone is an unspecified type.
	ListNone ListType = iota
	// ListEmail is an email distribution list.
	ListEmail
	// ListSMS is a text messaging distribution list.
	ListSMS
)

// ListTypeNames gives the display names of all list types.
var ListTypeNames = map[ListType]string{
	ListEmail: "email",
	ListSMS:   "sms",
}

// A PersonID is a positive integer uniquely identifying a Person.
type PersonID int

// AdminPersonID is the PersonID for the dedicated admin user.
const AdminPersonID PersonID = 1

// HasPrivLevel returns whether the receiver person has at least the specified
// privilege level on any org.
func (p *Person) HasPrivLevel(privLevel PrivLevel) bool {
	for _, org := range AllOrgs {
		if p.Orgs[org].PrivLevel >= privLevel {
			return true
		}
	}
	return false
}

// IsAdminLeader returns whether the receiver person is a Leader in the Admin
// org.
func (p *Person) IsAdminLeader() bool {
	return p.Orgs[OrgAdmin].PrivLevel >= PrivLeader
}

// A RoleID identifies a Role.
type RoleID int

// Webmaster is a role with special treatment.  It always has ID 1.
const Webmaster RoleID = 1

// DisabledUser is a role with special treatment.  It always has ID 2.
const DisabledUser RoleID = 2

// A RoleToList is a bitmask indicating the relationship between a list and the
// people holding a role.  The lower nibble contains the ListSubModel value and
// the upper nibble contains the Sender flag.
type RoleToList uint8

const rtlSender RoleToList = 0x10
const rtlSubModelMask RoleToList = 0x0F

// Sender returns whether holders of the role are allowed to send messages to
// the list.
func (rtl RoleToList) Sender() bool { return rtl&rtlSender != 0 }

// SubModel returns the list subscription model for holders of the role.
func (rtl RoleToList) SubModel() ListSubModel { return ListSubModel(rtl & rtlSubModelMask) }

// SetSender sets whether holders of the role are allowed to send messages to
// the list.
func (rtl *RoleToList) SetSender(sender bool) {
	if sender {
		*rtl |= rtlSender
	} else {
		*rtl &^= rtlSender
	}
}

// SetSubModel sets the list subscription model for holders of the role.
func (rtl *RoleToList) SetSubModel(subModel ListSubModel) {
	*rtl = (*rtl &^ rtlSubModelMask) | RoleToList(subModel)
}

// A SessionToken is a string that uniquely identifies a login session.
type SessionToken string

// A Session describes a login session.
type Session struct {
	Token   SessionToken
	Person  *Person
	Expires time.Time
	CSRF    CSRFToken
}

// A TextMessageID uniquely identifies an outgoing text message.
type TextMessageID int

// A VenueID uniquely identifies an event Venue.
type VenueID int
