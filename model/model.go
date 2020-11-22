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

// AllAttendanceTypes lists all defined attendance types.
var AllAttendanceTypes = []AttendanceType{AttendAsVolunteer, AttendAsStudent, AttendAsAuditor, AttendAsAbsent}

// AttendanceTypeNames gives the names for the attendance types.
var AttendanceTypeNames = map[AttendanceType]string{
	AttendAsVolunteer: "Volunteer",
	AttendAsStudent:   "Student",
	AttendAsAuditor:   "Audit",
	AttendAsAbsent:    "Absent",
}

// A CSRFToken is a random string used to verify that submitted forms came from
// our own site and not from a forgery.
type CSRFToken string

// A DocumentID is a positive integer uniquely identifying a document within its
// folder.  For cache-busting purposes, each new revision of a document gets a
// new DocumentID.
type DocumentID int

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
	// GroupStudents identifies the group to which all class students
	// belong.  Volunteer hours are not recorded for students (or, more
	// precisely, for people whose only organization-carrying role also puts
	// them in the students group).
	GroupStudents = "students"
)

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
	IDCERTShirt
)

// AllIdentTypes is the list of all identification types.
var AllIdentTypes = []IdentType{IDPhoto, IDCardKey, IDSERVShirt, IDCERTShirt}

// IdentTypeNames gives the names for each identification type.
var IdentTypeNames = map[IdentType]string{
	IDPhoto:     "photo ID",
	IDCardKey:   "access card",
	IDSERVShirt: "tan SERV shirt",
	IDCERTShirt: "green CERT shirt",
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

// An Org identifies one of the SERV volunteer organizations.
type Org uint8

// Values for Org.
const (
	OrgNone2 Org = iota
	OrgAdmin2
	OrgCERTD2
	OrgCERTT2
	OrgListos
	OrgSARES2
	OrgSNAP2
	NumOrgs
)

// AllOrgs gives the list of all Orgs.
var AllOrgs = []Org{OrgAdmin2, OrgCERTD2, OrgCERTT2, OrgListos, OrgSARES2, OrgSNAP2}

// OrgNames gives the display names of the Orgs.
var OrgNames = map[Org]string{
	OrgAdmin2: "admin",
	OrgCERTD2: "cert-d",
	OrgCERTT2: "cert-t",
	OrgListos: "listos",
	OrgSARES2: "sares",
	OrgSNAP2:  "snap",
}

// MembersCanViewContactInfo returns whether Members of the receiver Org can
// view contact info for members of the Org.
func (o Org) MembersCanViewContactInfo() bool { return o != OrgSARES2 }

// DSWClass returns the DSW class for the receiver Org.
func (o Org) DSWClass() DSWClass {
	switch o {
	case OrgSARES2:
		return DSWComm
	case OrgCERTD2, OrgCERTT2:
		return DSWCERT
	default:
		return 0
	}
}

// An Organization identifies one of the SERV volunteer organizations.
type Organization uint8

// Values for Organization.
const (
	OrgNone Organization = iota
	OrgAdmin
	OrgCERTT
	OrgLISTOS
	OrgOutreach
	OrgPEP
	OrgSARES
	OrgSNAP
	OrgCERTD
)

// AllOrganizations gives the list of all Organizations.
var AllOrganizations = []Organization{OrgAdmin, OrgCERTD, OrgCERTT, OrgLISTOS, OrgOutreach, OrgPEP, OrgSARES, OrgSNAP}

// CurrentOrganizations gives the list of Organizations that are currently in
// use (as opposed to historical).
var CurrentOrganizations = []Organization{OrgAdmin, OrgCERTD, OrgCERTT, OrgLISTOS, OrgSARES, OrgSNAP}

// OrganizationNames gives the names of the various Organizations.
var OrganizationNames = map[Organization]string{
	OrgNone:     "",
	OrgAdmin:    "Admin",
	OrgCERTD:    "CERT Deployment",
	OrgCERTT:    "CERT Training",
	OrgLISTOS:   "LISTOS",
	OrgOutreach: "Outreach",
	OrgPEP:      "PEP",
	OrgSARES:    "SARES",
	OrgSNAP:     "SNAP",
}

// OrganizationToDSWClass gives the DSW classes associated with each
// organization.
var OrganizationToDSWClass = map[Organization]DSWClass{
	OrgCERTD:    DSWCERT,
	OrgCERTT:    DSWCERT,
	OrgLISTOS:   DSWCERT,
	OrgOutreach: DSWCERT,
	OrgPEP:      DSWCERT,
	OrgSARES:    DSWComm,
	OrgSNAP:     DSWCERT,
}

// MarshalEasyJSON encodes the organization into JSON.
func (o Organization) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(OrganizationNames[o])
}

// UnmarshalEasyJSON decodes the organization from JSON.
func (o *Organization) UnmarshalEasyJSON(l *jlexer.Lexer) {
	s := l.UnsafeString()
	for org, name := range OrganizationNames {
		if s == name {
			*o = org
			return
		}
	}
	l.AddError(errors.New("unrecognized value for Organization"))
}

// A Permission is something holders of a role get to do.  Unlike a Privilege,
// it is not specific to a target group.  The type can be used as a single
// permission or a bitmask of multiple permissions.
type Permission uint16

// Known permission values.
const (
	// PermViewClearances allows its holders to view clearances (i.e.,
	// DSW registrations, background checks, etc.) for people whom they can
	// otherwise see.
	PermViewClearances Permission = 1 << iota

	// PermEditClearances allows its holders to edit clearances.
	PermEditClearances
)

// AllPermissions lists all possible permissions.
var AllPermissions = []Permission{
	PermViewClearances, PermEditClearances,
}

// PermissionNames gives the names of all of the permissions.
var PermissionNames = map[Permission]string{
	PermViewClearances: "permViewClearances",
	PermEditClearances: "permEditClearances",
}

// MarshalEasyJSON encodes the permission into JSON.
func (p Permission) MarshalEasyJSON(w *jwriter.Writer) {
	w.String(PermissionNames[p])
}

// UnmarshalEasyJSON decodes the permission from JSON.
func (p *Permission) UnmarshalEasyJSON(l *jlexer.Lexer) {
	s := l.UnsafeString()
	if s == "" {
		*p = 0
		return
	}
	for perm, name := range PermissionNames {
		if s == name {
			*p = perm
			return
		}
	}
	l.AddError(errors.New("unrecognized value for Permission"))
}

// A PersonID is a positive integer uniquely identifying a Person.
type PersonID int

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

// A PrivLevel is a privilege level for membership in an Org.
type PrivLevel uint8

// Values for PrivLevel.
const (
	// PrivNone indicates no membership in the Org.
	PrivNone PrivLevel = iota
	// PrivStudent indicates student-level membership in the Org, with
	// essentially no privileges other than being on lists.
	PrivStudent
	// PrivMember2 indicates full membership in the Org.
	PrivMember2
	// PrivLeader indicates a leader of the Org.
	PrivLeader
)

// AllPrivLevels is a list of all privilege levels.
var AllPrivLevels = []PrivLevel{PrivStudent, PrivMember2, PrivLeader}

// PrivLevelNames gives the display names of the privilege levels.
var PrivLevelNames = map[PrivLevel]string{
	PrivStudent: "student",
	PrivMember2: "member",
	PrivLeader:  "leader",
}

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

// A Role2ID identifies a Role.
type Role2ID int

// Webmaster is a role with special treatment.  It always has ID 1.
const Webmaster Role2ID = 1

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
