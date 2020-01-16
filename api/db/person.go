package db

import (
	"database/sql"
	"sort"
	"time"

	"rothskeller.net/serv/model"
)

// FetchPerson retrieves a single person from the database by ID.  It returns
// nil if no such person exists.
func (tx *Tx) FetchPerson(id model.PersonID) *model.Person {
	return tx.fetchPerson(`WHERE p.id=?`, id)
}

// FetchPersonByEmail retrieves a single person from the database, given their
// email address.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByEmail(email string) *model.Person {
	return tx.fetchPerson(`WHERE p.email=?`, email)
}

// FetchPersonByPWResetToken retrieves a single person from the database, given
// a password reset token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByPWResetToken(token string) *model.Person {
	return tx.fetchPerson(`WHERE p.pwreset_token=?`, token)
}

func (tx *Tx) fetchPerson(where string, args ...interface{}) (p *model.Person) {
	var (
		rows *sql.Rows
		err  error
	)
	p = new(model.Person)
	err = tx.tx.QueryRow(`SELECT p.id, p.first_name, p.last_name, p.email, p.phone, p.password, p.bad_login_count, p.bad_login_time, p.pwreset_token, p.pwreset_time FROM person p `+where, args...).Scan(&p.ID, &p.FirstName, &p.LastName, (*IDStr)(&p.Email), &p.Phone, &p.Password, &p.BadLoginCount, (*Time)(&p.BadLoginTime), (*IDStr)(&p.PWResetToken), (*Time)(&p.PWResetTime))
	if err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	rows, err = tx.tx.Query(`SELECT role FROM person_role WHERE person=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		var rid model.RoleID
		panicOnError(rows.Scan(&rid))
		role := tx.FetchRole(rid)
		p.Roles = append(p.Roles, role)
	}
	panicOnError(rows.Err())
	sort.Sort(model.RoleSort(p.Roles))
	tx.setPersonPrivileges(p)
	return p
}

// FetchPeople returns all of the people in the database, in alphabetical order
// by last name.
func (tx *Tx) FetchPeople() (people []*model.Person) {
	var (
		rows *sql.Rows
		err  error
		pmap = make(map[model.PersonID]*model.Person)
	)
	rows, err = tx.tx.Query(`SELECT id, first_name, last_name, email, phone, password, bad_login_count, bad_login_time, pwreset_token, pwreset_time FROM person ORDER BY last_name, first_name`)
	panicOnError(err)
	for rows.Next() {
		var p model.Person
		panicOnError(rows.Scan(&p.ID, &p.FirstName, &p.LastName, (*IDStr)(&p.Email), &p.Phone, &p.Password, &p.BadLoginCount, (*Time)(&p.BadLoginTime), (*IDStr)(&p.PWResetToken), (*Time)(&p.PWResetTime)))
		pmap[p.ID] = &p
		people = append(people, &p)
	}
	panicOnError(rows.Err())
	rows, err = tx.tx.Query(`SELECT person, role FROM person_role`)
	panicOnError(err)
	for rows.Next() {
		var (
			pid  model.PersonID
			rid  model.RoleID
			p    *model.Person
			role *model.Role
		)
		panicOnError(rows.Scan(&pid, &rid))
		p = pmap[pid]
		role = tx.FetchRole(rid)
		p.Roles = append(p.Roles, role)
	}
	panicOnError(rows.Err())
	for _, p := range people {
		sort.Sort(model.RoleSort(p.Roles))
		tx.setPersonPrivileges(p)
	}
	return people
}

func (tx *Tx) setPersonPrivileges(p *model.Person) {
	p.PrivMap = make(model.PrivilegeMap, tx.maxRoleID+1)
	for _, r := range p.Roles {
		p.PrivMap = p.PrivMap.Merge(r.TransPrivs)
	}
}

// SavePerson saves a person to the database.  If the supplied person ID is
// zero, a new person is added to the database; otherwise, the identified person
// is updated.
func (tx *Tx) SavePerson(p *model.Person) {
	var err error

	if p.ID == 0 {
		var result sql.Result
		result, err = tx.tx.Exec(`INSERT INTO person (first_name, last_name, email, phone, password, bad_login_count, bad_login_time, pwreset_token, pwreset_time) VALUES (?,?,?,?,?,?,?,?,?)`, p.FirstName, p.LastName, IDStr(p.Email), p.Phone, p.Password, p.BadLoginCount, Time(p.BadLoginTime), IDStr(p.PWResetToken), Time(p.PWResetTime))
		panicOnError(err)
		p.ID = model.PersonID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE person SET first_name=?, last_name=?, email=?, phone=?, password=?, bad_login_count=?, bad_login_time=?, pwreset_token=?, pwreset_time=? WHERE id=?`, p.FirstName, p.LastName, IDStr(p.Email), p.Phone, p.Password, p.BadLoginCount, Time(p.BadLoginTime), IDStr(p.PWResetToken), Time(p.PWResetTime), p.ID))
		panicOnExecError(tx.tx.Exec(`DELETE FROM person_role WHERE person=?`, p.ID))
	}
	for _, role := range p.Roles {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_role (person, role) VALUES (?,?)`, p.ID, role.ID))
	}
	tx.audit(model.AuditRecord{Person: p})
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
	panicOnExecError(tx.tx.Exec(`INSERT INTO session (token, person, expires) VALUES (?,?,?)`, s.Token, s.Person.ID, Time(s.Expires)))
	tx.audit(model.AuditRecord{Session: s})
}

// UpdateSession updates a session in the database.
func (tx *Tx) UpdateSession(s *model.Session) {
	panicOnNoRows(tx.tx.Exec(`UPDATE session SET expires=? WHERE token=?`, Time(s.Expires), s.Token))
	// deliberately not auditing
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE token=?`, s.Token))
	tx.audit(model.AuditRecord{Session: &model.Session{Token: s.Token, Person: s.Person}})
}

// DeleteSessionsForPerson deletes all sessions for the specified person, except
// the supplied one if any.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person, except model.SessionToken) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session where person=? AND token != ?`, p.ID, except))
	tx.audit(model.AuditRecord{Session: &model.Session{Person: p}})
}

// DeleteExpiredSessions deletes all expired sessions.
func (tx *Tx) DeleteExpiredSessions() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE expires<?`, Time(time.Now())))
	// deliberately not auditing
}
