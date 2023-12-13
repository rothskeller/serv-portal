package session

import (
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
)

const timestampFormat = "2006-01-02T15:04:05"
const createSQL = `INSERT INTO session (token, person, expires, csrf) VALUES (?,?,?,?)`

// Create creates a new session for a person.
func Create(storer phys.Storer, p *person.Person, token, csrf string, expires time.Time) {
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindText(token)
		stmt.BindInt(int(p.ID()))
		stmt.BindText(expires.In(time.Local).Format(timestampFormat))
		stmt.BindText(csrf)
		stmt.Step()
	})
	phys.Audit(storer, "AuthN:: ADD Session %s for person %q [%d] expires %s", token, p.InformalName(), p.ID(), expires.In(time.Local).Format(timestampFormat))
}

const deleteSQL = `DELETE FROM session WHERE token=?`

// Delete deletes the session with the specified token.
func Delete(storer phys.Storer, token string, p *person.Person) {
	phys.SQL(storer, deleteSQL, func(stmt *phys.Stmt) {
		stmt.BindText(token)
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "AuthN:: DELETE Session %s for person %q [%d]", token, p.InformalName(), p.ID())
	}
}

const deleteExpiredSQL = `DELETE FROM session WHERE expires<?`

// deleteExpired deletes all expired sessions.
func deleteExpired(storer phys.Storer) {
	phys.SQL(storer, deleteExpiredSQL, func(stmt *phys.Stmt) {
		stmt.BindText(time.Now().Format(timestampFormat))
		stmt.Step()
	})
}

const deleteForPersonSQL = `DELETE FROM session WHERE person=? AND token!=?`

// DeleteForPerson deletes all sessions for the specified person, except the one
// with the specified token (which can be "" if that exception is not needed).
func DeleteForPerson(storer phys.Storer, p *person.Person, token string) {
	phys.SQL(storer, deleteForPersonSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(p.ID()))
		stmt.BindText(token)
		stmt.Step()
	})
	if ra := phys.RowsAffected(storer); ra > 0 {
		if token != "" {
			phys.Audit(storer, "AuthN:: DELETE %d Sessions for person %q [%d] (all except %q)", ra, p.InformalName(), p.ID(), token)
		} else {
			phys.Audit(storer, "AuthN:: DELETE all %d Sessions for person %q [%d]", ra, p.InformalName(), p.ID())
		}
	}
}

const extendSQL = `UPDATE session SET expires=? WHERE token=?`

// Extend changes the expiration time of the specified session to the specified
// time.
func Extend(storer phys.Storer, token string, expires time.Time) {
	phys.SQL(storer, extendSQL, func(stmt *phys.Stmt) {
		stmt.BindText(expires.In(time.Local).Format(timestampFormat))
		stmt.BindText(token)
		stmt.Step()
	})
	// Intentionally not audited due to noise.
}
