package authn

import (
	"net/http"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

// Lifetime of a remember-me request.  (A year, more or less.)
const rememberMeExpiration = 365 * 24 * time.Hour

// GetLogin handles GET /api/login requests.
func GetLogin(r *util.Request) error {
	var out jwriter.Writer
	out.RawString(`{"id":`)
	out.Int(int(r.Person.ID))
	out.RawString(`,"informalName":`)
	out.String(r.Person.InformalName)
	out.RawString(`,"webmaster":`)
	out.Bool(r.Auth.IsWebmaster())
	out.RawString(`,"canAddEvents":`)
	out.Bool(r.Auth.CanA(model.PrivManageEvents))
	out.RawString(`,"canAddPeople":`)
	out.Bool(r.Auth.CanA(model.PrivManageMembers))
	out.RawString(`,"canSendTextMessages":`)
	out.Bool(r.Auth.CanA(model.PrivSendTextMessages))
	out.RawString(`,"canViewReports":`)
	out.Bool(r.Auth.CanAG(model.PrivManageEvents, r.Auth.FetchGroupByTag("cert-teams").ID))
	out.RawByte('}')
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostLogin handles POST /api/login requests.
func PostLogin(r *util.Request) error {
	var (
		person   *model.Person
		username = r.FormValue("username")
		password = r.FormValue("password")
	)
	// Check that the login is valid.
	if person = r.Tx.FetchPersonByUsername(username); person == nil {
		goto FAIL // no person with that username
	}
	if username != "admin" { // admin cannot be disabled or locked out
		if !IsEnabled(r, person) {
			goto FAIL // person is disabled
		}
		if person.BadLoginCount >= maxBadLogins && time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
			goto FAIL // locked out
		}
	}
	if !CheckPassword(person, password) {
		goto FAIL // password mismatch
	}
	// The login is valid.  Record it and create a session.
	r.Person = person
	r.Auth.SetMe(person)
	if person.BadLoginCount > 0 {
		r.Tx.WillUpdatePerson(person)
		person.BadLoginCount = 0
		r.Tx.UpdatePerson(person)
	}
	util.CreateSession(r)
	r.Tx.Commit()
	return GetLogin(r)

FAIL:
	if person != nil {
		// Record the bad login attempt.
		r.Tx.WillUpdatePerson(person)
		if time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
			person.BadLoginCount++
		} else {
			person.BadLoginCount = 1
		}
		person.BadLoginTime = time.Now()
		r.Tx.UpdatePerson(person)
	}
	r.Tx.Commit()
	return util.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
}

// PostLogout handles POST /api/logout requests.
func PostLogout(r *util.Request) error {
	if c, err := r.Cookie("auth"); err == nil {
		r.Tx.DeleteSession(&model.Session{Token: model.SessionToken(c.Value), Person: r.Person})
		http.SetCookie(r, &http.Cookie{
			Name:   "auth",
			Value:  c.Value,
			Path:   "/",
			MaxAge: -1,
		})
	}
	r.Tx.Commit()
	return nil
}
