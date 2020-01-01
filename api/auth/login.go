package auth

import (
	"net/http"
	"time"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

// Lifetime of a remember-me request.  (A year, more or less.)
const rememberMeExpiration = 365 * 24 * time.Hour

// PostLogin handles POST /api/login requests.
func PostLogin(r *util.Request) error {
	var (
		person   *model.Person
		email    = r.FormValue("email")
		password = r.FormValue("password")
		remember = r.FormValue("remember") != ""
	)
	// Check that the login is valid.
	if person = r.Tx.FetchPersonByEmail(email); person == nil {
		goto FAIL // no person with that email
	}
	if !person.CanLogIn() {
		goto FAIL // person not a member of any team
	}
	if person.PWResetToken != "" {
		goto FAIL // password reset in progress
	}
	if person.BadLoginCount >= maxBadLogins && time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
		goto FAIL // locked out
	}
	if !checkPassword(person, password) {
		goto FAIL // password mismatch
	}
	// The login is valid.  Record it and create a session.
	r.Person = person
	if person.BadLoginCount > 0 {
		person.BadLoginCount = 0
		r.Tx.SavePerson(person)
	}
	util.CreateSession(r)
	if remember {
		var rm = &model.RememberMe{
			Token:   model.RememberMeToken(util.RandomToken()),
			Person:  person,
			Expires: time.Now().Add(rememberMeExpiration),
		}
		r.Tx.CreateRememberMe(rm)
		http.SetCookie(r, &http.Cookie{
			Name:    "remember",
			Value:   string(rm.Token),
			Path:    "/",
			Expires: rm.Expires,
		})
	}
	r.Tx.Commit()
	return nil

FAIL:
	if person != nil {
		// Record the bad login attempt.
		if time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
			person.BadLoginCount++
		} else {
			person.BadLoginCount = 1
		}
		person.BadLoginTime = time.Now()
		r.Tx.SavePerson(person)
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
	if c, err := r.Cookie("remember"); err == nil {
		r.Tx.DeleteRememberMe(&model.RememberMe{Token: model.RememberMeToken(c.Value), Person: r.Person})
		http.SetCookie(r, &http.Cookie{
			Name:   "remember",
			Value:  c.Value,
			Path:   "/",
			MaxAge: -1,
		})
	}
	r.Tx.Commit()
	return nil
}
