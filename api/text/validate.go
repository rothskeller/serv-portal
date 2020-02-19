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
	var seenGroups = make(map[model.GroupID]bool)
	for _, gid := range tm.Groups {
		if seenGroups[gid] {
			return errors.New("duplicate group")
		}
		seenGroups[gid] = true
		if tx.Authorizer().FetchGroup(gid) == nil {
			return errors.New("nonexistent group")
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
