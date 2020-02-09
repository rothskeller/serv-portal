package store

import (
	"fmt"
	"strings"

	"sunnyvaleserv.org/portal/model"
)

// FetchEmailMessage returns the email message with the specified ID, or nil if
// there is none.
func (tx *Tx) FetchEmailMessage(id model.EmailMessageID) (e *model.EmailMessage) {
	return tx.tx.FetchEmailMessage(id)
}

// FetchEmailMessageByMessageID returns the email message with the specified
// MessageID header, or nil if there is none.
func (tx *Tx) FetchEmailMessageByMessageID(messageID string) (e *model.EmailMessage) {
	return tx.tx.FetchEmailMessageByMessageID(messageID)
}

// FetchEmailMessageBody returns the body of the email message with the
// specified ID.
func (tx *Tx) FetchEmailMessageBody(id model.EmailMessageID) (body []byte) {
	return tx.tx.FetchEmailMessageBody(id)
}

// FetchEmailMessages calls the supplied function for all email messages, in
// reverse chronological order.  It stops when the function return false.
func (tx *Tx) FetchEmailMessages(handler func(*model.EmailMessage) bool) {
	tx.tx.FetchEmailMessages(handler)
}

// CreateEmailMessage creates an email message in the database, including
// saving its body.
func (tx *Tx) CreateEmailMessage(em *model.EmailMessage, body []byte) {
	var gstr []string

	tx.tx.CreateEmailMessage(em, body)
	tx.entry.Change("create email [%d]", em.ID)
	tx.entry.Change("set email [%d] messageID to %q", em.ID, em.MessageID)
	tx.entry.Change("set email [%d] timestamp to %s", em.ID, em.Timestamp.Format("2006-01-02 15:04:05"))
	tx.entry.Change("set email [%d] type to %s", em.ID, model.EmailMessageTypeNames[em.Type])
	if em.Attention {
		tx.entry.Change("set email [%d] attention flag", em.ID)
	}
	if em.From != "" {
		tx.entry.Change("set email [%d] from to %q", em.ID, em.From)
	}
	if len(em.Groups) != 0 {
		for _, g := range em.Groups {
			gstr = append(gstr, fmt.Sprintf("%q [%d]", tx.Authorizer().FetchGroup(g).Name, g))
		}
		tx.entry.Change("set email [%d] groups to %s", em.ID, strings.Join(gstr, ", "))
	}
	if em.Subject != "" {
		tx.entry.Change("set email [%d] subject to %q", em.ID, em.Subject)
	}
	if em.Error != "" {
		tx.entry.Change("set email [%d] error to %q", em.ID, em.Error)
	}
}

// UpdateEmailMessage saves changes to an existing email message to the
// database.
func (tx *Tx) UpdateEmailMessage(em *model.EmailMessage) {
	var oem *model.EmailMessage

	oem = tx.tx.FetchEmailMessage(em.ID)
	tx.tx.UpdateEmailMessage(em)
	if em.MessageID != oem.MessageID {
		tx.entry.Change("set email [%d] messageID to %q", em.ID, em.MessageID)
	}
	if em.Timestamp != oem.Timestamp {
		tx.entry.Change("set email [%d] timestamp to %s", em.ID, em.Timestamp.Format("2006-01-02 15:04:05"))
	}
	if em.Type != oem.Type {
		tx.entry.Change("set email [%d] type to %s", em.ID, model.EmailMessageTypeNames[em.Type])
	}
	if em.Attention != oem.Attention {
		if em.Attention {
			tx.entry.Change("set email [%d] attention flag", em.ID)
		} else {
			tx.entry.Change("clear email [%d] attention flag", em.ID)
		}
	}
	if em.From != oem.From {
		tx.entry.Change("set email [%d] from to %q", em.ID, em.From)
	}
GROUPS1:
	for _, og := range oem.Groups {
		for _, g := range em.Groups {
			if og == g {
				continue GROUPS1
			}
		}
		tx.entry.Change("remove email [%d] group %q [%d]", em.ID, tx.Authorizer().FetchGroup(og).Name, og)
	}
GROUPS2:
	for _, g := range em.Groups {
		for _, og := range oem.Groups {
			if og == g {
				continue GROUPS2
			}
		}
		tx.entry.Change("add email [%d] group %q [%d]", em.ID, tx.Authorizer().FetchGroup(g).Name, g)
	}
	if em.Subject != oem.Subject {
		tx.entry.Change("set email [%d] subject to %q", em.ID, em.Subject)
	}
	if em.Error != oem.Error {
		tx.entry.Change("set email [%d] error to %q", em.ID, em.Error)
	}
}
