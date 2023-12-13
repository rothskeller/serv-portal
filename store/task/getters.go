package task

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
)

// Fields returns the set of fields that have been retrieved for this Task.
func (t *Task) Fields() Fields {
	return t.fields
}

// ID is the unique identifier of the Task.
func (t *Task) ID() ID {
	if t == nil {
		return 0
	}
	if t.fields&FID == 0 {
		panic("Task.ID called without having fetched FID")
	}
	return t.id
}

// Event is the ID of the Event to which the Task belongs.
func (t *Task) Event() event.ID {
	if t.fields&FEvent == 0 {
		panic("Task.Event called without having fetched FEvent")
	}
	return t.event
}

// Name is the Task's name.  It will normally be empty for a Task that is the
// only Task in its Event; it must be non-empty when the parent Event has
// multiple Tasks.
func (t *Task) Name() string {
	if t.fields&FName == 0 {
		panic("Task.Name called without having fetched FName")
	}
	return t.name
}

// Org is the SERV organization with which this Task is associated.  This
// determines which bucket volunteer hours are recorded in, which dot color(s)
// are associated with the parent event in the calendar, and which people can
// edit the task definition.
func (t *Task) Org() enum.Org {
	if t.fields&FOrg == 0 {
		panic("Task.Org called without having fetched FOrg")
	}
	return t.org
}

// Flags is a bitmask of flags describing the Task.
func (t *Task) Flags() Flag {
	if t.fields&FFlags == 0 {
		panic("Task.Flags called without having fetched FFlags")
	}
	return t.flags
}

// Details is free-form HTML text describing the Task.  It will generally be
// empty when the Task is the only Task under its parent Event.
func (t *Task) Details() string {
	if t.fields&FDetails == 0 {
		panic("Task.Details called without having fetched FDetails")
	}
	return t.details
}
