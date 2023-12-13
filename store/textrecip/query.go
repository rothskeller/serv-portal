package textrecip

import (
	"strings"
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
)

// FormatNumberForTwilio formats a phone number to be accepted by Twilio: all
// digits, no punctuation, with a +1 prefix.
func FormatNumberForTwilio(n string) string {
	if n == "" {
		return ""
	}
	n = strings.TrimPrefix(n, "+1")
	return "+1" + strings.Map(func(r rune) rune {
		if r < '0' || r > '9' {
			return -1
		}
		return r
	}, n)
}

// WithNumber returns the recipient of the specified text message who has the
// specified number, or nil if there is none.
func WithNumber(storer phys.Storer, tmid textmsg.ID, number string, fields person.Fields) (p *person.Person) {
	var sb strings.Builder
	sb.WriteString("SELECT ")
	person.ColumnList(&sb, fields)
	sb.WriteString(" FROM textmsg_recipient tr, person p WHER tr.textmsg=? AND tr.number=? AND tr.recipient=p.id")
	phys.SQL(storer, sb.String(), func(stmt *phys.Stmt) {
		stmt.BindInt(int(tmid))
		stmt.BindText(number)
		if stmt.Step() {
			p = new(person.Person)
			p.Scan(stmt, fields)
		}
	})
	return p
}

const allForTextSQL1 = `SELECT tr.number, tr.status, tr.timestamp, `
const allForTextSQL2 = ` FROM textmsg_recipient tr, person p WHERE tr.recipient=p.id AND tr.textmsg=? ORDER BY p.sort_name`

var allForTextSQLCache map[person.Fields]string

// AllRecipientsOfText fetches all of the recipients of the specified text message, in
// order by the recipient's sort name.
func AllRecipientsOfText(
	storer phys.Storer, tmid textmsg.ID, fields person.Fields,
	fn func(p *person.Person, number, status string, timestamp time.Time),
) {
	if allForTextSQLCache == nil {
		allForTextSQLCache = make(map[person.Fields]string)
	}
	if _, ok := allForTextSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(allForTextSQL1)
		person.ColumnList(&sb, fields)
		sb.WriteString(allForTextSQL2)
		allForTextSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allForTextSQLCache[fields], func(stmt *phys.Stmt) {
		var (
			p         person.Person
			number    string
			status    string
			timestamp time.Time
		)
		stmt.BindInt(int(tmid))
		for stmt.Step() {
			number = stmt.ColumnText()
			status = stmt.ColumnText()
			timestamp, _ = time.ParseInLocation(timestampFormat, stmt.ColumnText(), time.Local)
			p.Scan(stmt, fields)
			fn(&p, number, status, timestamp)
		}
	})
}

const allRepliesFromRecipientSQL = `SELECT reply, timestamp FROM textmsg_reply WHERE textmsg=? AND recipient=? ORDER BY timestamp DESC`

// AllRepliesFromRecipient fetches all of the replies to the specified text
// message by the specified recipient, in reverse chronological order.  The
// called function may return false to abort the fetch.
func AllRepliesFromRecipient(storer phys.Storer, tmid textmsg.ID, pid person.ID, fn func(string, time.Time) bool) {
	phys.SQL(storer, allRepliesFromRecipientSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(tmid))
		stmt.BindInt(int(pid))
		for stmt.Step() {
			var reply = stmt.ColumnText()
			var timestamp, _ = time.ParseInLocation(timestampFormat, stmt.ColumnText(), time.Local)
			if !fn(reply, timestamp) {
				return
			}
		}
	})
}
