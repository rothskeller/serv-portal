package textmsg

import (
	"time"

	"sunnyvaleserv.org/portal/store/person"
)

// Fields returns the set of fields that have been retrieved for this text
// message.
func (t *TextMessage) Fields() Fields {
	return t.fields
}

// ID is the unique identifier of the TextMessage.
func (t *TextMessage) ID() ID {
	if t == nil {
		return 0
	}
	if t.fields&FID == 0 {
		panic("TextMessage.ID called without having fetched FID")
	}
	return t.id
}

// Sender is the ID of the Person who sent the TextMessage.
func (t *TextMessage) Sender() person.ID {
	if t.fields&FSender == 0 {
		panic("TextMessage.Sender called without having fetched FSender")
	}
	return t.sender
}

// Timestamp is the time at which the TextMessage was sent.
func (t *TextMessage) Timestamp() time.Time {
	if t.fields&FTimestamp == 0 {
		panic("TextMessage.Timestamp called without having fetched FTimestamp")
	}
	return t.timestamp
}

// Message is the text of the TextMessage that was sent.
func (t *TextMessage) Message() string {
	if t.fields&FMessage == 0 {
		panic("TextMessage.Message called without having fetched FMessage")
	}
	return t.message
}

// Lists is set of lists to which the TextMessage was sent.
func (t *TextMessage) Lists() []TextToList {
	if t.fields&FLists == 0 {
		panic("TextMessage.Lists called without having fetched FLists")
	}
	return t.lists
}
