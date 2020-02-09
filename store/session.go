package store

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchSessions fetches all sessions in the database.
func (tx *Tx) FetchSessions() (sessions []*model.Session) {
	return tx.tx.FetchSessions()
}

// FetchSession fetches the session with the specified token.  It does not check
// for session expiration.  It returns nil if no such session exists.
func (tx *Tx) FetchSession(token model.SessionToken) (s *model.Session) {
	return tx.tx.FetchSession(token)
}

// CreateSession creates a session in the database.
func (tx *Tx) CreateSession(s *model.Session) {
	tx.tx.CreateSession(s)
	tx.entry.Change("created session %s for person %q [%d]", s.Token, s.Person.InformalName, s.Person.ID)
}

// UpdateSession updates a session in the database.
func (tx *Tx) UpdateSession(s *model.Session) {
	tx.tx.UpdateSession(s)
	// deliberately not auditing
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	tx.tx.DeleteSession(s)
	tx.entry.Change("deleted session %s", s.Token)
}

// DeleteSessionsForPerson deletes all sessions for the specified person, except
// the supplied one if any.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person, except model.SessionToken) {
	tx.tx.DeleteSessionsForPerson(p, except)
	if except != "" {
		tx.entry.Change("deleted all sessions for person %q [%d] except %s", p.InformalName, p.ID, except)
	} else {
		tx.entry.Change("deleted all sessions for person %q [%d]", p.InformalName, p.ID)
	}
}

// DeleteExpiredSessions deletes all expired sessions.
func (tx *Tx) DeleteExpiredSessions() {
	tx.tx.DeleteExpiredSessions()
	// deliberately not auditing
}
