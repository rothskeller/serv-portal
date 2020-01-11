// Package model contains the data model types and constants for the SERV
// portal.
package model

import (
	"time"
)

// An EventID is a positive integer uniquely identifying an Event.
type EventID int

// An EventType is a string identifying a type of Event.
type EventType string

// The known event types.
const (
	EventTraining EventType = "Train"
	EventDrill              = "Drill"
	EventCivic              = "Civic"
	EventIncident           = "Incid"
	EventContEd             = "CE"
	EventMeeting            = "Meeting"
	EventClass              = "Class"
)

// AllEventTypes is the list of all known event types.
var AllEventTypes = []EventType{EventTraining, EventDrill, EventCivic, EventIncident, EventContEd, EventMeeting, EventClass}

// An Event structure contains the data describing an event at which SERV
// volunteers may participate.
type Event struct {
	ID    EventID
	Date  string // 2006-01-02
	Name  string
	Hours float64
	Type  EventType
	Roles []*Role
}

// A PersonID is a positive integer uniquely identifying a Person.
type PersonID int

// A Person structure contains the data describing a person involved (or
// formerly involved) with SERV, and therefore a person with a (potentially
// disabled) login to the SERV portal.
type Person struct {
	ID            PersonID
	FirstName     string
	LastName      string
	Email         string
	Phone         string
	Password      string
	BadLoginCount int
	BadLoginTime  time.Time
	PWResetToken  string
	PWResetTime   time.Time
	Roles         []*Role
	PrivMap       PrivilegeMap // transient, transitive
}

// A Privilege is something members of an actor team get to do to a target team.
// The type can be used as a single privilege or a bitmask of multiple
// privileges.
type Privilege uint8

// Known privilege values.
const (
	// PrivHoldsRole isn't a privilege per se; it denotes holding the target
	// role.
	PrivHoldsRole Privilege = 1 << iota

	// PrivViewHolders denotes the ability to view the list of people who
	// hold the target role.
	PrivViewHolders

	// PrivAssignRole denotes the ability to assign people to the target
	// role or remove them from it.
	PrivAssignRole

	// PrivManageEvents denotes the ability to manage events to which the
	// role is invited.
	PrivManageEvents
)

// A RoleID is a positive integer uniquely identifying a Role.
type RoleID int

// A RoleTag is a string that uniquely identifies some teams that are treated as
// special cases by the web site code.
type RoleTag string

// Values for RoleTag.
const (
	// RoleDisabled identifies the role to which disabled users belong.
	// Holding this role blocks logging into the web site.
	RoleDisabled RoleTag = "disabled"

	// RoleWebmaster identifies the role to which all webmasters belong.
	// Webmasters have all privileges on all roles.
	RoleWebmaster = "webmaster"
)

// A Role describes a role that a Person can hold.
type Role struct {
	ID          RoleID
	Tag         RoleTag
	Name        string
	MemberLabel string
	SERVGroup   SERVGroup
	ImplyOnly   bool
	Individual  bool
	PrivMap     PrivilegeMap // persistent, non-transitive
	TransPrivs  PrivilegeMap // transient, transitive
}

// A SERVGroup identifies one of the main SERV groups.
type SERVGroup string

// Values for SERVGroup.
const (
	GroupSERV           SERVGroup = "SERV"
	GroupCERTDeployment           = "CERT-D"
	GroupCERTTraining             = "CERT-T"
	GroupListos                   = "Listos"
	GroupOutreach                 = "Outreach"
	GroupPEP                      = "PEP"
	GroupSARES                    = "SARES"
	GroupCountyARES               = "SCC ARES"
	GroupSNAP                     = "SNAP"
)

// AllSERVGroups is the list of all SERV groups.
var AllSERVGroups = []SERVGroup{
	GroupSERV, GroupCERTDeployment, GroupCERTTraining, GroupListos,
	GroupOutreach, GroupPEP, GroupSARES, GroupCountyARES, GroupSNAP,
}

// A SessionToken is a string that uniquely identifies a login session.
type SessionToken string

// A Session describes a login session.
type Session struct {
	Token   SessionToken
	Person  *Person
	Expires time.Time
}
