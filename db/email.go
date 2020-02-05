package db

import (
	"database/sql"

	"sunnyvaleserv.org/portal/model"
)

// FetchEmailMessage returns the email message with the specified ID, or nil if
// there is none.
func (tx *Tx) FetchEmailMessage(id model.EmailMessageID) (e *model.EmailMessage) {
	var data []byte
	e = new(model.EmailMessage)
	switch err := tx.tx.QueryRow(`SELECT data FROM email_message WHERE id=?`, id).Scan(&data); err {
	case nil:
		panicOnError(e.Unmarshal(data))
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}

}

// FetchEmailMessageByMessageID returns the email message with the specified
// MessageID header, or nil if there is none.
func (tx *Tx) FetchEmailMessageByMessageID(messageID string) (e *model.EmailMessage) {
	var data []byte
	e = new(model.EmailMessage)
	switch err := tx.tx.QueryRow(`SELECT data FROM email_message WHERE message_id=?`, messageID).Scan(&data); err {
	case nil:
		panicOnError(e.Unmarshal(data))
		return e
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}

}

// FetchEmailMessageBody returns the body of the email message with the
// specified ID.
func (tx *Tx) FetchEmailMessageBody(id model.EmailMessageID) (body []byte) {
	panicOnError(tx.tx.QueryRow(`SELECT body FROM email_message_body WHERE id=?`, id).Scan(&body))
	return body
}

// FetchEmailMessages calls the supplied function for all email messages, in
// reverse chronological order.  It stops when the function return false.
func (tx *Tx) FetchEmailMessages(handler func(*model.EmailMessage) bool) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM email_message ORDER BY timestamp DESC`)
	panicOnError(err)
	for rows.Next() {
		var data []byte
		var e model.EmailMessage
		panicOnError(rows.Scan(&data))
		panicOnError(e.Unmarshal(data))
		if !handler(&e) {
			rows.Close()
			break
		}
	}
	panicOnError(rows.Err())
}

// CreateEmailMessage creates an email message in the database, including
// saving its body.
func (tx *Tx) CreateEmailMessage(em *model.EmailMessage, body []byte) {
	var (
		data []byte
		err  error
	)
	panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id), 0) FROM email_message`).Scan(&em.ID))
	em.ID++
	data, err = em.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`INSERT INTO email_message (id, message_id, timestamp, data) VALUES (?,?,?,?)`, em.ID, em.MessageID, Time(em.Timestamp), data))
	panicOnExecError(tx.tx.Exec(`INSERT INTO email_message_body (id, body) VALUES (?,?)`, em.ID, body))
	tx.audit("email_message", em.ID, data)
}

// UpdateEmailMessage saves changes to an existing email message to the
// database.
func (tx *Tx) UpdateEmailMessage(em *model.EmailMessage) {
	var (
		data []byte
		err  error
	)
	data, err = em.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE email_message SET (message_id, timestamp, data) = (?,?,?) WHERE id=?`, em.MessageID, Time(em.Timestamp), data, em.ID))
	tx.audit("email_message", em.ID, data)
}
