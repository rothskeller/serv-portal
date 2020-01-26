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
	switch err := dbh.QueryRow(`SELECT data FROM person WHERE id=?`, id).Scan(&data); err {
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
	switch err := dbh.QueryRow(`SELECT data FROM person WHERE username=?`, username).Scan(&data); err {
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
	switch err := dbh.QueryRow(`SELECT data FROM person WHERE pwreset_token=?`, token).Scan(&data); err {
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
	rows, err = dbh.Query(`SELECT data FROM person`)
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

// SavePerson saves a person to the database.  If the supplied person ID is
// zero, a new person is added to the database; otherwise, the identified person
// is updated.
func (tx *Tx) SavePerson(p *model.Person) {
	var (
		data []byte
		err  error
	)
	tx.recalcPersonPrivileges(p)
	if p.ID == 0 {
		panicOnError(tx.tx.QueryRow(`SELECT max(id) FROM person`).Scan(&p.ID))
		p.ID++
		data, err = p.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`INSERT INTO person (id, username, pwreset_token, data) VALUES (?,?,?,?)`, p.ID, IDStr(p.Username), IDStr(p.PWResetToken), data))
	} else {
		data, err = p.Marshal()
		panicOnError(err)
		panicOnExecError(tx.tx.Exec(`UPDATE person SET (username, pwreset_token, data) = (?,?,?) WHERE id=?`, IDStr(p.Username), IDStr(p.PWResetToken), data, p.ID))
	}
	tx.audit("person", p.ID, data)
}

// recalcAllPersonPrivileges recalculates the privileges map for every person in
// the database, from the privilege masks for the role(s) held by that person.
// This is done whenever the role privilege masks may have changed, i.e., when
// roles or groups are edited.
func (tx *Tx) recalcAllPersonPrivileges() {
	people := tx.FetchPeople()
	for _, p := range people {
		tx.SavePerson(p)
	}
}

// recalcPersonPrivileges recalculates the privileges map for a person from the
// privilege masks for the role(s) held by that person.
func (tx *Tx) recalcPersonPrivileges(p *model.Person) {
	p.Privileges.Clear()
	for _, r := range p.Roles {
		p.Privileges.Merge(&tx.roles[r].Privileges)
	}
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
	tx.audit("session", s.Token, buf.Bytes())
}

// UpdateSession updates a session in the database.
func (tx *Tx) UpdateSession(s *model.Session) {
	panicOnNoRows(tx.tx.Exec(`UPDATE session SET expires=? WHERE token=?`, Time(s.Expires), s.Token))
	// deliberately not auditing
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE token=?`, s.Token))
	tx.audit("session", s.Token, nil)
}

// DeleteSessionsForPerson deletes all sessions for the specified person, except
// the supplied one if any.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person, except model.SessionToken) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session where person=? AND token != ?`, p.ID, except))
	tx.audit("session", p.Username, nil)
}

// DeleteExpiredSessions deletes all expired sessions.
func (tx *Tx) DeleteExpiredSessions() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE expires<?`, Time(time.Now())))
	// deliberately not auditing
}
