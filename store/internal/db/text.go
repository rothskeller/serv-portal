package db

import (
	"database/sql"

	"sunnyvaleserv.org/portal/model"
)

// FetchTextMessages returns a list of outgoing text messages, in reverse
// chronological order.
func (tx *Tx) FetchTextMessages() (messages []*model.TextMessage) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM text_message ORDER BY id DESC`)
	panicOnError(err)
	for rows.Next() {
		var (
			data    []byte
			message model.TextMessage
		)
		panicOnError(rows.Scan(&data))
		panicOnError(message.Unmarshal(data))
		messages = append(messages, &message)
	}
	panicOnError(rows.Err())
	return messages
}

// FetchTextMessage returns the text message with the specified ID, or nil if
// there is none.
func (tx *Tx) FetchTextMessage(id model.TextMessageID) (message *model.TextMessage) {
	var (
		data []byte
		err  error
	)
	if err = tx.tx.QueryRow(`SELECT data FROM text_message WHERE id=?`, id).Scan(&data); err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	message = new(model.TextMessage)
	panicOnError(message.Unmarshal(data))
	return message
}

// FetchTextMessageByNumber returns the text message most recently sent to the
// specified phone number, or nil if there is none.
func (tx *Tx) FetchTextMessageByNumber(number string) (message *model.TextMessage) {
	var (
		data []byte
		err  error
	)
	if err = tx.tx.QueryRow(`SELECT m.data FROM text_number n, text_message m WHERE n.number=? AND n.mid=m.id`, number).Scan(&data); err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	message = new(model.TextMessage)
	panicOnError(message.Unmarshal(data))
	return message
}

// CreateTextMessage creates a new text message in the database, with the next
// available ID.
func (tx *Tx) CreateTextMessage(message *model.TextMessage) {
	var (
		data []byte
		err  error
	)
	panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id),0)+1 FROM text_message`).Scan(&message.ID))
	data, err = message.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`INSERT INTO text_message (id, data) VALUES (?,?)`, message.ID, data))
	for _, r := range message.Recipients {
		if r.Number != "" {
			panicOnExecError(tx.tx.Exec(`INSERT OR REPLACE INTO text_number (number, mid) VALUES (?,?)`, r.Number, message.ID))
		}
	}
}

// UpdateTextMessage updates an existing text message in the database.
func (tx *Tx) UpdateTextMessage(message *model.TextMessage) {
	var (
		data []byte
		err  error
	)
	data, err = message.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE text_message SET data=? WHERE id=?`, data, message.ID))
}
