package db

import (
	"bytes"
	"database/sql"
	"fmt"
	"sort"
	"time"

	"sunnyvaleserv.org/portal/model"
)

// FetchPerson retrieves a single person from the database by ID.  It returns
// nil if no such person exists.
func (tx *Tx) FetchPerson(id model.PersonID) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE id=?`, id).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPersonByUsername retrieves a single person from the database, given
// their username.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByUsername(username string) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE username=?`, username).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPersonByPWResetToken retrieves a single person from the database, given
// a password reset token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByPWResetToken(token string) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE pwreset_token=?`, token).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPersonByCellPhone retrieves a single person from the database, given a
// cell phone number.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByCellPhone(token string) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE cell_phone=?`, token).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPersonByUnsubscribe retrieves a single person from the database, given
// an unsubscribe token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByUnsubscribe(token string) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE unsubscribe=?`, token).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPeople returns all of the people in the database, in order by sortname.
func (tx *Tx) FetchPeople() (people []*model.Person) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM person`)
	panicOnError(err)
	for rows.Next() {
		var data []byte
		var p model.Person
		panicOnError(rows.Scan(&data))
		p.Unmarshal(data)
		people = append(people, &p)
	}
	panicOnError(rows.Err())
	sort.Sort(model.PersonSort(people))
	return people
}

// CreatePerson creates a new person in the database, assigning the next
// available person ID.
func (tx *Tx) CreatePerson(p *model.Person) {
	var (
		data []byte
		err  error
	)
	panicOnError(tx.tx.QueryRow(`SELECT max(id) FROM person`).Scan(&p.ID))
	p.ID++
	data, err = p.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`INSERT INTO person (id, username, pwreset_token, cell_phone, unsubscribe, data) VALUES (?,?,?,?,?,?)`, p.ID, IDStr(p.Username), IDStr(p.PWResetToken), IDStr(p.CellPhone), p.UnsubscribeToken, data))
	if p.Email != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email))
	}
	if p.Email2 != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email2))
	}
}

// UpdatePerson updates a person in the database.
func (tx *Tx) UpdatePerson(p *model.Person) {
	var (
		data []byte
		err  error
	)
	data, err = p.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE person SET (username, pwreset_token, cell_phone, unsubscribe, data) = (?,?,?,?,?) WHERE id=?`, IDStr(p.Username), IDStr(p.PWResetToken), IDStr(p.CellPhone), p.UnsubscribeToken, data, p.ID))
	panicOnExecError(tx.tx.Exec(`DELETE FROM person_email WHERE person=?`, p.ID))
	if p.Email != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email))
	}
	if p.Email2 != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email2))
	}
}

// FetchSessions fetches all sessions in the database.
func (tx *Tx) FetchSessions() (sessions []*model.Session) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT token, person, expires FROM session`)
	panicOnError(err)
	for rows.Next() {
		var session model.Session
		var pid model.PersonID
		panicOnError(rows.Scan(&session.Token, &pid, (*Time)(&session.Expires)))
		session.Person = tx.FetchPerson(pid)
		sessions = append(sessions, &session)
	}
	panicOnError(rows.Err())
	return sessions
}

// FetchSession fetches the session with the specified token.  It does not check
// for session expiration.  It returns nil if no such session exists.
func (tx *Tx) FetchSession(token model.SessionToken) (s *model.Session) {
	var pid model.PersonID

	s = &model.Session{Token: token}
	switch err := tx.tx.QueryRow(`SELECT person, expires FROM session WHERE token=?`, token).Scan(&pid, (*Time)(&s.Expires)); err {
	case nil:
		s.Person = tx.FetchPerson(pid)
		return s
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// CreateSession creates a session in the database.
func (tx *Tx) CreateSession(s *model.Session) {
	var buf bytes.Buffer
	panicOnExecError(tx.tx.Exec(`INSERT INTO session (token, person, expires) VALUES (?,?,?)`, s.Token, s.Person.ID, Time(s.Expires)))
	fmt.Fprintf(&buf, "person:%s expires:%s", s.Person.Username, s.Expires.Format("2006-01-02 15:04:05"))
}

// UpdateSession updates a session in the database.
func (tx *Tx) UpdateSession(s *model.Session) {
	panicOnNoRows(tx.tx.Exec(`UPDATE session SET expires=? WHERE token=?`, Time(s.Expires), s.Token))
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE token=?`, s.Token))
}

// DeleteSessionsForPerson deletes all sessions for the specified person, except
// the supplied one if any.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person, except model.SessionToken) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session where person=? AND token != ?`, p.ID, except))
}

// DeleteExpiredSessions deletes all expired sessions.
func (tx *Tx) DeleteExpiredSessions() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE expires<?`, Time(time.Now())))
}
