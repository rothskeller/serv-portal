package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateTextMessage creates a new text message in the database, with the next
// available ID.
func (tx *Tx) CreateTextMessage(message *model.TextMessage) {
	tx.Tx.CreateTextMessage(message)
	tx.entry.Change("create text [%d]", message.ID)
	tx.entry.Change("set text [%d] sender to person %q [%d]", message.ID, tx.FetchPerson(message.Sender).InformalName, message.Sender)
	tx.entry.Change("set text [%d] timestamp to %s", message.ID, message.Timestamp.Format("2006-01-02 15:04:05"))
	tx.entry.Change("set text [%d] message to %q", message.ID, message.Message)
	for _, r := range message.Recipients {
		if r.Status != "" {
			tx.entry.Change("add text [%d] recipient %q [%d] status %q", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Status)
		} else {
			tx.entry.Change("add text [%d] recipient %q [%d] number %s", message.ID, tx.FetchPerson(r.Recipient).InformalName, r.Recipient, r.Number)
		}
	}
	for _, l := range message.Lists {
		tx.entry.Change("add text [%d] list %q [%d]", message.ID, l, tx.FetchList(l).Name)
	}
}

// UpdateTextMessage updates an existing text message in the database.
func (tx *Tx) UpdateTextMessage(message *model.TextMessage) {
	var om = tx.Tx.FetchTextMessage(message.ID)

	tx.Tx.UpdateTextMessage(message)
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
