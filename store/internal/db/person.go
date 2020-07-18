package db

import (
	"database/sql"
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

// FetchPersonByHoursToken retrieves a single person from the database, given an
// hours token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByHoursToken(token string) (p *model.Person) {
	var data []byte
	p = new(model.Person)
	switch err := tx.tx.QueryRow(`SELECT data FROM person WHERE hours_token=?`, token).Scan(&data); err {
	case nil:
		panicOnError(p.Unmarshal(data))
		return p
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchPersonByEmail retrieves a single person from the database, given an
// email address.  It returns nil if no such person exists or if more than one
// person has that email address.
func (tx *Tx) FetchPersonByEmail(email string) (p *model.Person) {
	var (
		rows *sql.Rows
		pid  model.PersonID
		err  error
	)
	p = new(model.Person)
	rows, err = tx.tx.Query(`SELECT person FROM person_email WHERE email=?`, email)
	panicOnError(err)
	for rows.Next() {
		if pid != 0 {
			rows.Close()
			return nil
		}
		panicOnError(rows.Scan(&pid))
	}
	panicOnError(rows.Err())
	return tx.FetchPerson(pid)
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
	panicOnExecError(tx.tx.Exec(`INSERT INTO person (id, username, pwreset_token, cell_phone, unsubscribe, hours_token, data) VALUES (?,?,?,?,?,?,?)`, p.ID, IDStr(p.Username), IDStr(p.PWResetToken), IDStr(p.CellPhone), p.UnsubscribeToken, IDStr(p.HoursToken), data))
	if p.Email != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email))
	}
	if p.Email2 != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email2))
	}
	tx.indexPerson(p, false)
}

// UpdatePerson updates a person in the database.
func (tx *Tx) UpdatePerson(p *model.Person) {
	var (
		data []byte
		err  error
	)
	data, err = p.Marshal()
	panicOnError(err)
	panicOnNoRows(tx.tx.Exec(`UPDATE person SET (username, pwreset_token, cell_phone, unsubscribe, hours_token, data) = (?,?,?,?,?,?) WHERE id=?`, IDStr(p.Username), IDStr(p.PWResetToken), IDStr(p.CellPhone), p.UnsubscribeToken, IDStr(p.HoursToken), data, p.ID))
	panicOnExecError(tx.tx.Exec(`DELETE FROM person_email WHERE person=?`, p.ID))
	if p.Email != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email))
	}
	if p.Email2 != "" {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_email (person, email) VALUES (?,?)`, p.ID, p.Email2))
	}
	tx.indexPerson(p, true)
}

// indexPerson updates the search index with information about a person.
func (tx *Tx) indexPerson(p *model.Person, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='person' AND id=?`, p.ID))
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, personInformalName, personFormalName, personCallSign, personEmail, personEmail2, personHomeAddress, personWorkAddress,personMailAddress) VALUES ('person',?,?,?,?,?,?,?,?,?)`, p.ID, p.InformalName, p.FormalName, IDStr(p.CallSign), IDStr(p.Email), IDStr(p.Email2), IDStr(p.HomeAddress.Address), IDStr(p.WorkAddress.Address), IDStr(p.MailAddress.Address)))
}

// FetchSessions fetches all sessions in the database.
func (tx *Tx) FetchSessions() (sessions []*model.Session) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT token, person, expires, csrf FROM session`)
	panicOnError(err)
	for rows.Next() {
		var session model.Session
		var pid model.PersonID
		panicOnError(rows.Scan(&session.Token, &pid, (*Time)(&session.Expires), &session.CSRF))
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
	switch err := tx.tx.QueryRow(`SELECT person, expires, csrf FROM session WHERE token=?`, token).Scan(&pid, (*Time)(&s.Expires), &s.CSRF); err {
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
	panicOnExecError(tx.tx.Exec(`INSERT INTO session (token, person, expires, csrf) VALUES (?,?,?,?)`, s.Token, s.Person.ID, Time(s.Expires), s.CSRF))
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
