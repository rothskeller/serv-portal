// Package textmsg defines the TextMessage type, which describes a text (SMS)
// message sent by this system.
package textmsg

import (
	"time"

	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
)

// ID uniquely identifies a text message.
type ID int

// Fields is a bitmask of flags identifying specified fields of the TextMessage
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FSender
	FTimestamp
	FMessage
	FLists
)

// TextMessage describes a text (SMS) message sent by this system.
type TextMessage struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields    Fields // which fields of the structure are populated
	id        ID
	sender    person.ID
	timestamp time.Time
	message   string
	lists     []TextToList
}

// TextToList describes a list to which a text message was sent.
type TextToList struct {
	// ID is the identifier of the list.  It may be zero if the list has
	// been subsequently deleted.
	ID list.ID
	// Name of the list at the time the message was sent.  It will not
	// reflect any subsequent changes to the list name.
	Name string
}
