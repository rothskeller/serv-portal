package text

import (
	"errors"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

// ValidateTextMessage ensures the correctness of a text message.
func ValidateTextMessage(tx *store.Tx, tm *model.TextMessage) (err error) {
	if tm.ID < 0 {
		return errors.New("invalid ID")
	}
	if tx.FetchPerson(tm.Sender) == nil {
		return errors.New("nonexistent sender")
	}
	var seenLists = make(map[model.ListID]bool)
	for _, lid := range tm.Lists {
		if seenLists[lid] {
			return errors.New("duplicate list")
		}
		if list := tx.FetchList(lid); list == nil {
			return errors.New("nonexistent list")
		} else if list.Type != model.ListSMS {
			return errors.New("invalid list type")
		}
	}
	if tm.Timestamp.IsZero() {
		return errors.New("missing timestamp")
	}
	if tm.Message == "" {
		return errors.New("missing message")
	}
	var seenRecipients = make(map[model.PersonID]bool)
	for _, tr := range tm.Recipients {
		if seenRecipients[tr.Recipient] {
			return errors.New("duplicate recipient")
		}
		if tx.FetchPerson(tr.Recipient) == nil {
			return errors.New("nonexistent recipient")
		}
	}
	return nil
}
