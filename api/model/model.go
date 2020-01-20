// Package model contains the data model types and constants for the SERV
// portal.
package model

//go:generate protoc -I=. -I=$GOPATH/src -I=$GOPATH/src/github.com/gogo/protobuf/protobuf --gogo_out=Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types:. model.proto
//go:generate easyjson -all model.pb.go

import (
	"time"
)

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

/*
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
*/

// Hours returns the number of hours the event lasted.
func (e *Event) Hours() float64 {
	start, _ := time.Parse("15:04", e.Start)
	end, _ := time.Parse("15:04", e.End)
	return float64(end.Sub(start)) / float64(time.Hour)
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

	// GroupSccAres identifies the group that is invited to events imported
	// from the scc-ares-races.org database.
	GroupSccAres = "scc-ares"
)

// A PersonID is a positive integer uniquely identifying a Person.
type PersonID int

/*
// A Person structure contains the data describing a person involved (or
// formerly involved) with SERV, and therefore a person with a (potentially
// disabled) login to the SERV portal.
type Person struct {
	ID            PersonID
	FirstName     string
	LastName      string
	Nickname      string
	Suffix        string
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
*/

// A Privilege is something holders of a role get to do to a target group.  The
// type can be used as a single privilege or a bitmask of multiple privileges.
type Privilege uint8

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
)

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
)

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

/*
// A Venue is a place where Events occur.
type Venue struct {
	ID      VenueID
	Name    string
	Address string
	City    string
	URL     string
}
*/
