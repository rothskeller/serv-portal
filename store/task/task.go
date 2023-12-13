// Package task defines the data model for SERV tasks.
package task

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
)

// ID uniquely identifies a task.
type ID int

// Flag is a flag (or bitmask of flags) describing aspects of the Task.
type Flag uint

// Value for Flag:
const (
	// RecordHours is a flag indicating that hours spent on this task count
	// as volunteer hours that should be recorded.
	RecordHours Flag = 1 << iota
	// CoveredByDSW is a flag indicating that this activity is covered by
	// Sunnyvale DSW, for participants who are properly registered, etc.  It
	// also implies that current Sunnyvale DSW registration is required in
	// order to sign up for the task.
	CoveredByDSW
	// RequiresBGCheck is a flag indicating that this task is for people
	// who have passed a background check.  Other people cannot sign up.
	RequiresBGCheck
	// SignupsOpen is a flag indicating that eligible people can sign up for
	// this task (if it is in the future) or can record hours spent on it
	// even if they hadn't signed up for it (if it is in the recent past).
	SignupsOpen
	// HasAttended is a flag indicating that at least one person has been
	// marked as having attended the Task.
	HasAttended
	// HasCredited is a flag indicating that at least one person has been
	// credited with participation in the Task.
	HasCredited
)

// Fields is a bitmask of flags identifying specified fields of the Task
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FEvent
	FName
	FOrg
	FFlags
	FDetails
)

// Task describes a single task in an event on the SERV calendar.
type Task struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields  Fields // which fields of the structure are populated
	id      ID
	event   event.ID
	name    string
	org     enum.Org
	flags   Flag
	details string
}

func (t *Task) Clone() (c *Task) {
	c = new(Task)
	*c = *t
	return c
}
