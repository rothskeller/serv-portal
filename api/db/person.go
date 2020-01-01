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

// FetchPersonByRememberMeToken retrieves a single person from the database,
// given a remember-me token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByRememberMeToken(token string) *model.Person {
	return tx.fetchPerson(`, remember_me r WHERE r.person=p.ID AND r.token=?`, token)
}

// FetchPersonBySessionToken retrieves a single person from the database, given
// a session token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonBySessionToken(token string) *model.Person {
	return tx.fetchPerson(`, session s WHERE s.person=p.ID AND s.token=?`, token)
}

func (tx *Tx) fetchPerson(where string, args ...interface{}) (p *model.Person) {
	var (
		rows *sql.Rows
		err  error
	)
	p = new(model.Person)
	err = tx.tx.QueryRow(`SELECT p.id, p.first_name, p.last_name, p.email, p.phone, p.password, p.bad_login_count, p.bad_login_time, p.pwreset_token, p.pwreset_time FROM person p `+where, args...).Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Password, &p.BadLoginCount, (*Time)(&p.BadLoginTime), (*IDStr)(&p.PWResetToken), (*Time)(&p.PWResetTime))
	if err == sql.ErrNoRows {
		return nil
	}
	panicOnError(err)
	p.PrivMap = make(model.PrivilegeMap)
	rows, err = tx.tx.Query(`SELECT role FROM person_role WHERE person=?`, p.ID)
	panicOnError(err)
	for rows.Next() {
		var rid model.RoleID
		panicOnError(rows.Scan(&rid))
		role := tx.FetchRole(rid)
		p.Roles = append(p.Roles, role)
		p.PrivMap.Merge(role.PrivMap)
	}
	panicOnError(rows.Err())
	sort.Slice(p.Roles, func(i, j int) bool {
		return roleLess(p.Roles[i], p.Roles[j])
	})
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
		panicOnError(rows.Scan(&p.ID, &p.FirstName, &p.LastName, &p.Email, &p.Phone, &p.Password, &p.BadLoginCount, (*Time)(&p.BadLoginTime), (*IDStr)(&p.PWResetToken), (*Time)(&p.PWResetTime)))
		p.PrivMap = make(model.PrivilegeMap)
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
		p.PrivMap.Merge(role.PrivMap)
	}
	panicOnError(rows.Err())
	for _, p := range people {
		sort.Slice(p.Roles, func(i, j int) bool {
			return roleLess(p.Roles[i], p.Roles[j])
		})
	}
	return people
}

func roleLess(i, j *model.Role) bool {
	switch {
	case i.Team.Name < j.Team.Name:
		return true
	case i.Team.Name > j.Team.Name:
		return false
	default:
		return i.Name < j.Name
	}
}

// SavePerson saves a person to the database.  If the supplied person ID is
// zero, a new person is added to the database; otherwise, the identified person
// is updated.
func (tx *Tx) SavePerson(p *model.Person) {
	var err error

	if p.ID == 0 {
		var result sql.Result
		result, err = tx.tx.Exec(`INSERT INTO person (first_name, last_name, email, phone, password, bad_login_count, bad_login_time, pwreset_token, pwreset_time) VALUES (?,?,?,?,?,?,?,?,?)`, p.FirstName, p.LastName, p.Email, p.Phone, p.Password, p.BadLoginCount, Time(p.BadLoginTime), IDStr(p.PWResetToken), Time(p.PWResetTime))
		panicOnError(err)
		p.ID = model.PersonID(lastInsertID(result))
	} else {
		panicOnNoRows(tx.tx.Exec(`UPDATE person SET first_name=?, last_name=?, email=?, phone=?, password=?, bad_login_count=?, bad_login_time=?, pwreset_token=?, pwreset_time=? WHERE id=?`, p.FirstName, p.LastName, p.Email, p.Phone, p.Password, p.BadLoginCount, Time(p.BadLoginTime), IDStr(p.PWResetToken), Time(p.PWResetTime), p.ID))
		panicOnExecError(tx.tx.Exec(`DELETE FROM person_role WHERE person=?`, p.ID))
	}
	for _, role := range p.Roles {
		panicOnExecError(tx.tx.Exec(`INSERT INTO person_role (person, role) VALUES (?,?)`, p.ID, role.ID))
	}
	tx.audit(model.AuditRecord{Person: p})
}

// CreateSession creates a session in the database.
func (tx *Tx) CreateSession(s *model.Session) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO session (token, person, expires) VALUES (?,?,?)`, s.Token, s.Person.ID, Time(s.Expires)))
	tx.audit(model.AuditRecord{Session: s})
}

// DeleteSession deletes a session from the database.
func (tx *Tx) DeleteSession(s *model.Session) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE token=?`, s.Token))
	tx.audit(model.AuditRecord{Session: &model.Session{Token: s.Token, Person: s.Person}})
}

// DeleteSessionsForPerson deletes all sessions for the specified person.
func (tx *Tx) DeleteSessionsForPerson(p *model.Person) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session where person=?`, p.ID))
	tx.audit(model.AuditRecord{Session: &model.Session{Person: p}})
}

// DeleteExpiredSessions deletes all expired sessions.
func (tx *Tx) DeleteExpiredSessions() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM session WHERE expires<?`, Time(time.Now())))
	// deliberately not auditing
}

// CreateRememberMe creates a remember-me request in the database.
func (tx *Tx) CreateRememberMe(rm *model.RememberMe) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO remember_me (token, person, expires) VALUES (?,?,?)`, rm.Token, rm.Person.ID, Time(rm.Expires)))
	tx.audit(model.AuditRecord{RememberMe: rm})
}

// DeleteRememberMe deletes a remember-me request in the database.
func (tx *Tx) DeleteRememberMe(rm *model.RememberMe) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM remember_me WHERE token=?`, rm.Token))
	tx.audit(model.AuditRecord{RememberMe: &model.RememberMe{Token: rm.Token, Person: rm.Person}})
}

// DeleteRememberMeForPerson deletes all remember-me requests for the specified
// person.
func (tx *Tx) DeleteRememberMeForPerson(p *model.Person) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM remember_me where person=?`, p.ID))
	tx.audit(model.AuditRecord{RememberMe: &model.RememberMe{Person: p}})
}

// DeleteExpiredRememberMeTokens deletes all expired remember-me tokens.
func (tx *Tx) DeleteExpiredRememberMeTokens() {
	panicOnExecError(tx.tx.Exec(`DELETE FROM remember_me WHERE expires<?`, Time(time.Now())))
	// deliberately not auditing
}
