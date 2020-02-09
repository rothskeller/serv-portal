package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateSession creates a session in the database.
func (tx *Tx) CreateSession(s *model.Session) {
	tx.Tx.CreateSession(s)
	tx.entry.Change("created session %s for person %q [%d]", s.Token, s.Person.InformalName, s.Person.ID)
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	tx.Tx.DeleteSession(s)
	tx.entry.Change("deleted session %s", s.Token)
}

// DeleteSessionsForPerson deletes all sessions for the specified person, except
// the supplied one if any.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person, except model.SessionToken) {
	tx.Tx.DeleteSessionsForPerson(p, except)
	if except != "" {
		tx.entry.Change("deleted all sessions for person %q [%d] except %s", p.InformalName, p.ID, except)
	} else {
		tx.entry.Change("deleted all sessions for person %q [%d]", p.InformalName, p.ID)
	}
}
