package session

import (
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
)

const withTokenSQL = `SELECT person, expires, csrf FROM session WHERE token=?`

// WithToken returns the session data associated with the specified session
// token.  It returns zero values if the session does not exist or has expired.
func WithToken(storer phys.Storer, token string) (pid person.ID, expires time.Time, csrf string) {
	// We do not attempt to join with the person table and return a person
	// object, because the caller almost certainly wants joins against the
	// person sub-tables.  We'll just return the person ID and let them call
	// person.WithID to get what they want.
	deleteExpired(storer)
	phys.SQL(storer, withTokenSQL, func(stmt *phys.Stmt) {
		stmt.BindText(token)
		if stmt.Step() {
			pid = person.ID(stmt.ColumnInt())
			expires, _ = time.ParseInLocation(timestampFormat, stmt.ColumnText(), time.Local)
			csrf = stmt.ColumnText()
		}
	})
	return
}
