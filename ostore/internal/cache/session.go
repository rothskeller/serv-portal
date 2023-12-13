package cache

import (
	"sunnyvaleserv.org/portal/model"
)

// FetchSessions fetches all sessions in the database.
func (tx *Tx) FetchSessions() (sessions []*model.Session) {
	sessions = tx.Tx.FetchSessions()
	for _, s := range sessions {
		if p := tx.people[s.Person.ID]; p != nil {
			s.Person = p
		} else {
			tx.people[s.Person.ID] = s.Person
		}
	}
	return sessions
}

// FetchSession fetches the session with the specified token.  It does not check
// for session expiration.  It returns nil if no such session exists.
func (tx *Tx) FetchSession(token model.SessionToken) (s *model.Session) {
	s = tx.Tx.FetchSession(token)
	if s != nil {
		if p := tx.people[s.Person.ID]; p != nil {
			s.Person = p
		} else {
			tx.people[s.Person.ID] = s.Person
		}
	}
	return s
}
