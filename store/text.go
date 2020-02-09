package store

import (
	"fmt"
	"strings"

	"sunnyvaleserv.org/portal/model"
)

// FetchTextMessages returns a list of outgoing text messages, in reverse
// chronological order.
func (tx *Tx) FetchTextMessages() (messages []*model.TextMessage) {
	return tx.tx.FetchTextMessages()
}

// FetchTextMessage returns the text message with the specified ID, or nil if
// there is none.
func (tx *Tx) FetchTextMessage(id model.TextMessageID) (message *model.TextMessage) {
	return tx.tx.FetchTextMessage(id)
}

// FetchTextMessageByNumber returns the text message most recently sent to the
// specified phone number, or nil if there is none.
func (tx *Tx) FetchTextMessageByNumber(number string) (message *model.TextMessage) {
	return tx.tx.FetchTextMessageByNumber(number)
}

// CreateTextMessage creates a new text message in the database, with the next
// available ID.
func (tx *Tx) CreateTextMessage(message *model.TextMessage) {
	var gstr []string

	tx.tx.CreateTextMessage(message)
	tx.entry.Change("create text [%d]", message.ID)
	tx.entry.Change("set text [%d] sender to person %q [%d]", message.ID, tx.FetchPerson(message.Sender).InformalName, message.Sender)
	if len(message.Groups) != 0 {
		for _, g := range message.Groups {
			gstr = append(gstr, fmt.Sprintf("%q [%d]", tx.Authorizer().FetchGroup(g).Name, g))
		}
		tx.entry.Change("set text [%d] groups to %s", message.ID, strings.Join(gstr, ", "))
	}
	tx.entry.Change("set text [%d] timestamp to %s", message.ID, message.Timestamp.Format("2006-01-02 15:04:05"))
	tx.entry.Change("set text [%d] message to %q", message.ID, message.Message)
	for _, r := range message.Recipients {
		if r.Status != "" {
			tx.entry.Change("add text [%d] recipient %q [%d] status %q", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Status)
		} else {
			tx.entry.Change("add text [%d] recipient %q [%d] number %s", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Number)
		}
	}
}

// UpdateTextMessage updates an existing text message in the database.
func (tx *Tx) UpdateTextMessage(message *model.TextMessage) {
	var om = tx.tx.FetchTextMessage(message.ID)

	tx.tx.UpdateTextMessage(message)
	// Nothing can change in a text except the status, timestamp, and
	// responses of each recipient.
	if len(message.Recipients) != len(om.Recipients) {
		panic("text recipient list should not change")
	}
	for i, r := range message.Recipients {
		or := om.Recipients[i]
		if r.Recipient != or.Recipient {
			panic("text recipient list should not change")
		}
		if r.Status != or.Status {
			tx.entry.Change("set text [%d] recipient %s [%d] status to %q", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Status)
		}
		if r.Timestamp != or.Timestamp {
			tx.entry.Change("set text [%d] recipient %s [%d] timestamp to %s", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Timestamp.Format("2006-01-02 15:04:05"))
		}
		if len(or.Responses) > len(r.Responses) {
			panic("text responses should not be removed")
		}
		for j := len(or.Responses); j < len(r.Responses); j++ {
			rr := r.Responses[j]
			tx.entry.Change("add text [%d] recipient %q [%d] response %s %q", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, rr.Timestamp.Format("2006-01-02 15:04:05"), rr.Response)
		}
	}
}
