package db

import (
	"database/sql"
	"fmt"

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

// FetchTextDeliveries returns the set of delivery records for a text message.
// FetchTextMessages returns a list of outgoing text messages, in reverse
// chronological order.
func (tx *Tx) FetchTextDeliveries(id model.TextMessageID) (deliveries []*model.TextDelivery) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM text_delivery WHERE message=?`, id)
	panicOnError(err)
	for rows.Next() {
		var (
			data     []byte
			delivery model.TextDelivery
		)
		panicOnError(rows.Scan(&data))
		panicOnError(delivery.Unmarshal(data))
		deliveries = append(deliveries, &delivery)
	}
	panicOnError(rows.Err())
	return deliveries
}

// FetchTextDelivery returns the text message delivery record for the specified
// message and recipient, or nil if there is none.
func (tx *Tx) FetchTextDelivery(message model.TextMessageID, recipient model.PersonID) (delivery *model.TextDelivery) {
	var (
		data []byte
		err  error
	)
	if err = tx.tx.QueryRow(`SELECT data FROM text_delivery WHERE message=? AND recipient=?`, message, recipient).Scan(&data); err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	delivery = new(model.TextDelivery)
	panicOnError(delivery.Unmarshal(data))
	return delivery
}

// FetchNewestTextDelivery returns the most recent text message delivery record
// for the specified recipient, or nil if there is none.
func (tx *Tx) FetchNewestTextDelivery(recipient model.PersonID) (delivery *model.TextDelivery) {
	var (
		data []byte
		err  error
	)
	if err = tx.tx.QueryRow(`SELECT data FROM text_delivery WHERE recipient=? ORDER BY message DESC LIMIT 1`, recipient).Scan(&data); err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	delivery = new(model.TextDelivery)
	panicOnError(delivery.Unmarshal(data))
	return delivery
}

// SaveTextMessage saves the supplied text message in the database.  If it does
// not already have an ID, it assigns one.  If deliveries are specified, it
// creates them.
func (tx *Tx) SaveTextMessage(message *model.TextMessage, deliveries []*model.TextDelivery) {
	var (
		data []byte
		err  error
	)
	if message.ID == 0 {
		panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id),0)+1 FROM text_message`).Scan(&message.ID))
		data, err = message.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`INSERT INTO text_message (id, data) VALUES (?,?)`, message.ID, data))
	} else {
		data, err = message.Marshal()
		panicOnError(err)
		panicOnNoRows(tx.tx.Exec(`UPDATE text_message SET data=? WHERE id=?`, data, message.ID))
	}
	tx.audit("text_message", message.ID, data)
	if deliveries == nil {
		return
	}
	for _, d := range deliveries {
		d.Message = message.ID
		data, err = d.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`INSERT INTO text_delivery (message, recipient, data) VALUES (?,?,?)`, d.Message, d.Recipient, data))
		tx.audit("text_delivery", fmt.Sprintf("%d-%d", d.Message, d.Recipient), data)
	}
}

// SaveTextDelivery saves the supplied text delivery in the database.
func (tx *Tx) SaveTextDelivery(delivery *model.TextDelivery) {
	var (
		data []byte
		err  error
	)
	data, err = delivery.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE text_delivery SET data=? WHERE message=? AND number=?`, data, delivery.Message, delivery.Recipient))
	tx.audit("text_delivery", fmt.Sprintf("%d-%d", delivery.Message, delivery.Recipient), data)
}
