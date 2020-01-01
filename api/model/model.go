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
	Teams []*Team
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
	PrivMap       PrivilegeMap // transient
}

// A Privilege is something members of an actor team get to do to a target team.
// The type can be used as a single privilege or a bitmask of multiple
// privileges.
type Privilege uint8

// Known privilege values.
const (
	// PrivMember isn't a privilege per se; it denotes membership in the
	// target team.
	PrivMember Privilege = 1 << iota

	// PrivView denotes the ability to view the roster of the target team.
	PrivView

	// PrivAdmin denotes the ability to administer the target team, i.e.,
	// change the roles its members have (but not add or remove members),
	// and schedule events.
	PrivAdmin

	// PrivManage denotes the ability to manage the target team, i.e., add
	// and remove members.
	PrivManage
)

// A PrivilegeMap expresses the privileges held by an acting Person, Role, or
// Team.  It is a map from the target team to the privileges the actor has on
// that target team.
type PrivilegeMap map[*Team]Privilege

// A RoleID is a positive integer uniquely identifying a Role.
type RoleID int

// A Role describes a role that a Person can hold in a team.  A Role hAS no
// meaning outside the context of the team that defines it.
type Role struct {
	ID      RoleID
	Team    *Team
	Name    string
	PrivMap PrivilegeMap
}

// A RememberMeToken is a string that uniquely identifies a remember-me request.
type RememberMeToken string

// A RememberMe describes a remember-me request.
type RememberMe struct {
	Token   RememberMeToken
	Person  *Person
	Expires time.Time
}

// A SessionToken is a string that uniquely identifies a login session.
type SessionToken string

// A Session describes a login session.
type Session struct {
	Token   SessionToken
	Person  *Person
	Expires time.Time
}

// A TeamID is an integer that uniquely identifies a team.
type TeamID int

// A TeamTag is a string that uniquely identifies some teams that are treated as
// special cases by the web site code.
type TeamTag string

// Values for TeamTag.
const (
	// TeamLogin identifies the team to which all web site users belong.
	// Membership in this team implies the right to log in to the web site.
	TeamLogin TeamTag = "login"

	// TeamWebmasters identifies the team to which all webmasters belong.
	// Webmasters have all privileges on all teams.
	TeamWebmasters = "webmaster"
)

// TeamType specifies how the team membership is constructed.
type TeamType uint8

// Values for TeamType.
const (
	// TeamExplicit states that team membership comprises normal (non-tied)
	// roles defined on that team, and the people who hold them.
	TeamNormal TeamType = iota

	// TeamTiedRoles states that team membership comprises tied roles on
	// that team, and the people who hold them.
	TeamTiedRoles

	// TeamAncestor states that team membership comprises the membership of
	// child teams; the team has no explicit roles of its own.
	TeamAncestor
)

// A Team is a group of people who are visible and/or mailable.
type Team struct {
	ID       TeamID
	Parent   *Team
	Tag      TeamTag
	Type     TeamType
	Name     string
	Email    string
	PrivMap  PrivilegeMap
	Roles    []*Role // transient
	Children []*Team // transient
}
