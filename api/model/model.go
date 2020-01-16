// Package model contains the data model types and constants for the SERV
// portal.
package model

import (
	"time"
)

// An EventID is a positive integer uniquely identifying an Event.
type EventID int

// An EventType is a bitmask identifying the type(s) of an Event (usually only
// one, but sometimes more).
type EventType uint16

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

// An Event structure contains the data describing an event at which SERV
// volunteers may participate.
type Event struct {
	ID        EventID
	Name      string
	Date      string // 2006-01-02
	Start     string // 13:45
	End       string // 14:45
	Venue     *Venue
	Details   string
	Type      EventType
	Roles     []*Role
	SccAresID string
}

// Hours returns the number of hours the event lasted.
func (e *Event) Hours() float64 {
	start, _ := time.Parse("15:04", e.Start)
	end, _ := time.Parse("15:04", e.End)
	return float64(end.Sub(start)) / float64(time.Hour)
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

	// RoleSccAres identifies the role that is invited to events imported
	// from the scc-ares-races.org database.
	RoleSccAres = "scc-ares"
)

// A Role describes a role that a Person can hold.
type Role struct {
	ID          RoleID
	Tag         RoleTag
	Name        string
	MemberLabel string
	ImplyOnly   bool
	Individual  bool
	PrivMap     PrivilegeMap // persistent, non-transitive
	TransPrivs  PrivilegeMap // transient, transitive
}

// A SessionToken is a string that uniquely identifies a login session.
type SessionToken string

// A Session describes a login session.
type Session struct {
	Token   SessionToken
	Person  *Person
	Expires time.Time
}

// A VenueID uniquely identifies an event Venue.
type VenueID int

// A Venue is a place where Events occur.
type Venue struct {
	ID      VenueID
	Name    string
	Address string
	City    string
	URL     string
}
