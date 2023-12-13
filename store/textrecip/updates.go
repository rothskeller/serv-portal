package textrecip

import (
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
)

// Milliseconds are used because some systems send multiple adjancent
// autoreplies that can arrive during the same second, and we need milliseconds
// to know the proper order.
const timestampFormat = "2006-01-02T15:04:05.000"
const addRecipientSQL = `INSERT INTO textmsg_recipient (textmsg, recipient, number, status, timestamp) VALUES (?,?,?,?,?)`
const addNumberSQL = `
INSERT INTO textmsg_number (number, textmsg) VALUES (?,?)
ON CONFLICT DO UPDATE SET textmsg=?2`

// AddRecipient adds a single recipient to a text message.
func AddRecipient(storer phys.Storer, t *textmsg.TextMessage, p *person.Person, number, status string, timestamp time.Time) {
	phys.SQL(storer, addRecipientSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.BindNullText(number)
		stmt.BindNullText(status)
		stmt.BindText(timestamp.In(time.Local).Format(timestampFormat))
		stmt.Step()
	})
	if number != "" {
		phys.SQL(storer, addNumberSQL, func(stmt *phys.Stmt) {
			stmt.BindText(number)
			stmt.BindInt(int(t.ID()))
			stmt.Step()
		})
	}
	phys.Audit(storer, "ADD TextMessage %d:: Recipient %q [%d] = number %q status %q", t.ID(), p.InformalName(), p.ID(), number, status)
}

const updateStatusSQL = `UPDATE textmsg_recipient SET status=?, timestamp=? WHERE textmsg=? AND recipient=?`

// UpdateStatus updates the status of a single recipient to a text message.
func UpdateStatus(storer phys.Storer, t *textmsg.TextMessage, p *person.Person, status string, timestamp time.Time) {
	phys.SQL(storer, updateStatusSQL, func(stmt *phys.Stmt) {
		stmt.BindNullText(status)
		stmt.BindText(timestamp.In(time.Local).Format(timestampFormat))
		stmt.BindInt(int(t.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "TextMessage %d:: Recipient %q [%d] = status %q", t.ID(), p.InformalName(), p.ID(), status)
}

const addReplySQL = `INSERT INTO textmsg_reply (textmsg, recipient, reply, timestamp) VALUES (?,?,?,?)`

// AddReply adds a reply from a recipient of a text message.
func AddReply(storer phys.Storer, t *textmsg.TextMessage, p *person.Person, reply string, timestamp time.Time) {
	phys.SQL(storer, addReplySQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		stmt.BindInt(int(p.ID()))
		stmt.BindText(reply)
		stmt.BindText(timestamp.In(time.Local).Format(timestampFormat))
		stmt.Step()
	})
	phys.Audit(storer, "TextMessage %d:: Recipient %q [%d]:: Reply %q", t.ID(), p.InformalName(), p.ID(), reply)
}
