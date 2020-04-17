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

// An AttendanceInfo structure gives information about a person's attendance at
// an event.
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
)

// AttendanceTypeNames gives the names for the attendance types.
var AttendanceTypeNames = map[AttendanceType]string{
	AttendAsVolunteer: "Volunteer",
	AttendAsStudent:   "Student",
	AttendAsAuditor:   "Audit",
}

// A CSRFToken is a random string used to verify that submitted forms came from
// our own site and not from a forgery.
type CSRFToken string

// A DocumentID is a positive integer uniquely identifying a document within its
// folder.  For cache-busting purposes, each new revision of a document gets a
// new DocumentID.
type DocumentID int

// An EmailMessageID is a positive integer uniquely identifying an email message
// handled by the portal.
type EmailMessageID int

// An EmailMessageType describes the type of an email message handled by the
// portal.
type EmailMessageType byte

// Values for EmailMessageType
const (
	EmailBounce EmailMessageType = iota
	EmailSent
	EmailSendFailed
	EmailModerated
	EmailUnrecognized
)

// AllEmailMessageTypes lists all of the known email message types.
var AllEmailMessageTypes = []EmailMessageType{EmailBounce, EmailSent, EmailSendFailed, EmailModerated, EmailUnrecognized}

// EmailMessageTypeNames gives names to all of the EmailMessageType values.
var EmailMessageTypeNames = map[EmailMessageType]string{
	EmailBounce:       "bounce",
	EmailSent:         "sent",
	EmailSendFailed:   "send_failed",
	EmailModerated:    "moderated",
	EmailUnrecognized: "unrecognized",
}

// An EventID is a positive integer uniquely identifying an Event.
type EventID int

// An EventType is a bitmask identifying the type(s) of an Event (usually only
// one, but sometimes more).
type EventType uint32

// The known event types.
const (
	EventCivic EventType = 1 << iota
	EventClass
	EventContEd
	EventDrill
	EventIncident
	EventMeeting
	EventParty
	EventTraining
	EventWork
)

// AllEventTypes is the list of all known event types.
var AllEventTypes = []EventType{
	EventCivic,
	EventClass,
	EventContEd,
	EventDrill,
	EventIncident,
	EventMeeting,
	EventParty,
	EventTraining,
	EventWork,
}

// EventTypeNames maps event types to their names.
var EventTypeNames = map[EventType]string{
	EventCivic:    "Civic Event",
	EventClass:    "Class",
	EventContEd:   "Continuing Education",
	EventDrill:    "Drill",
	EventIncident: "Incident",
	EventMeeting:  "Meeting",
	EventParty:    "Party",
	EventTraining: "Training",
	EventWork:     "Work Event",
}

// Hours returns the number of hours the event lasted.
func (e *Event) Hours() float64 {
	start, _ := time.Parse("15:04", e.Start)
	end, _ := time.Parse("15:04", e.End)
	return float64(end.Sub(start)) / float64(time.Hour)
}

// A FolderID is a positive integer uniquely identifying a Folder.
type FolderID int

// A FolderNode is a Folder, plus the link fields necessary to construct the
// tree of folders.
type FolderNode struct {
	*Folder
	ParentNode *FolderNode
	ChildNodes []*FolderNode
}

// A GroupID is a positive integer uniquely identifying a Group.
type GroupID int

// A GroupTag is a string that uniquely identifies a group for programmatic
// access.  This is used for some groups that are treated as special cases by
// the web site code.
type GroupTag string

// Values for GroupTag.
const (
	// GroupDisabled identifies the group to which disabled users belong.
	// Members of this group are blocked from logging into the web site.
	GroupDisabled GroupTag = "disabled"
)

// An Organization identifies one of the SERV volunteer organizations.
type Organization uint8

// Values for Organization.
const (
	OrgNone Organization = iota
	OrgAdmin
	OrgCERT
	OrgLISTOS
	OrgOutreach
	OrgPEP
	OrgSARES
	OrgSNAP
)

// AllOrganizations gives the list of all Organizations.
var AllOrganizations = []Organization{OrgAdmin, OrgCERT, OrgLISTOS, OrgOutreach, OrgPEP, OrgSARES, OrgSNAP}

// OrganizationNames gives the names of the various Organizations.
var OrganizationNames = map[Organization]string{
	OrgNone:     "",
	OrgAdmin:    "Admin",
	OrgCERT:     "CERT",
	OrgLISTOS:   "LISTOS",
	OrgOutreach: "Outreach",
	OrgPEP:      "PEP",
	OrgSARES:    "SARES",
	OrgSNAP:     "SNAP",
}

// A PersonID is a positive integer uniquely identifying a Person.
type PersonID int

// A Privilege is something holders of a role get to do to a target group.  The
// type can be used as a single privilege or a bitmask of multiple privileges.
type Privilege uint16

// Known privilege values.
const (
	// PrivMember isn't a privilege per se; it denotes that holders of the
	// actor role are members of the target group.
	PrivMember Privilege = 1 << iota

	// PrivViewMembers denotes the ability to view the list of people who
	// are members of the target group.
	PrivViewMembers

	// PrivViewContactInfo denotes the ability to view contact information
	// for the members of the target group.
	PrivViewContactInfo

	// PrivManageMembers denotes the ability to manage the membership of the
	// target group, i.e., to add or remove members from it, and assign its
	// roles to its members.  Holding this privilege against any target
	// group implicitly denotes the ability to create new users.
	PrivManageMembers

	// PrivManageEvents denotes the ability to manage events to which the
	// target group is invited.
	PrivManageEvents

	// PrivSendTextMessages denotes the ability to send text messages to the
	// members of the target group.
	PrivSendTextMessages

	// PrivSendEmailMessages denotes the ability to send unmoderated email
	// messages to the members of the target group.
	PrivSendEmailMessages

	// PrivBCC denotes that holders of the actor role receive BCC copies of
	// emails sent to the target group.
	PrivBCC

	// PrivManageFolders denotes the ability to manage folders (and the
	// documents within them) belonging to the target group.
	PrivManageFolders
)

// AllPrivileges lists all possible privileges.
var AllPrivileges = []Privilege{
	PrivMember, PrivViewMembers, PrivViewContactInfo, PrivManageMembers, PrivManageEvents, PrivSendTextMessages,
	PrivSendEmailMessages, PrivBCC, PrivManageFolders,
}

// PrivilegeNames gives the names of all of the privileges.
var PrivilegeNames = map[Privilege]string{
	PrivMember:            "member",
	PrivViewMembers:       "roster",
	PrivViewContactInfo:   "contact",
	PrivManageMembers:     "admin",
	PrivManageEvents:      "events",
	PrivSendTextMessages:  "texts",
	PrivSendEmailMessages: "emails",
	PrivBCC:               "bcc",
	PrivManageFolders:     "folders",
}

// MarshalEasyJSON encodes the privilege into JSON.
func (p Privilege) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(PrivilegeNames[p])
}

// UnmarshalEasyJSON decodes the privilege from JSON.
func (p *Privilege) UnmarshalEasyJSON(l *jlexer.Lexer) {
	s := l.UnsafeString()
	if s == "" {
		*p = 0
		return
	}
	for priv, name := range PrivilegeNames {
		if s == name {
			*p = priv
			return
		}
	}
	l.AddError(errors.New("unrecognized value for Privilege"))
}

// A RoleID is a positive integer uniquely identifying a Role.
type RoleID int

// A RoleTag is a string that uniquely identifies a role for programmatic
// access.  This is used for some roles that are treated as special cases by the
// web site code.
type RoleTag string

// Values for RoleTag.
const (
	// RoleWebmaster identifies the webmaster role.  People holding this
	// role have all privileges on all groups.
	RoleWebmaster = "webmaster"

	// RoleDisabled identifies the role assigned to disabled users.  People
	// holding this role are in the Disabled Users group and therefore are
	// blocked from logging in.
	RoleDisabled = "disabled"
)

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
