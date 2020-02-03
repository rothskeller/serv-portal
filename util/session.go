package util

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/model"
)

// Lifetime of a login session.
const sessionExpiration = time.Hour

// Forbidden is the error returned when the calling session lacks the privileges
// needed for the call it issued.
var Forbidden = HTTPError(http.StatusForbidden, "403 Forbidden")

// ValidateSession decodes the auth and/or remember tokens in the request, if
// any, and sets the Person field of the request appropriately.  If the session
// authorization is not valid, the Person field is left unchanged (i.e., nil).
func ValidateSession(r *Request) {
	var (
		c   *http.Cookie
		err error
	)
	if c, err = r.Cookie("auth"); err != nil {
		return
	}
	r.Tx.DeleteExpiredSessions()
	if r.Session = r.Tx.FetchSession(model.SessionToken(c.Value)); r.Session == nil {
		return
	}
	if r.Method != "GET" && r.Session.Token != model.SessionToken(r.Request.Header.Get("X-XSRF-TOKEN")) {
		r.Session = nil
		return
	}
	r.Person = r.Session.Person
	r.Auth.SetMe(r.Session.Person)
	r.Session.Expires = time.Now().Add(sessionExpiration)
	r.Tx.UpdateSession(r.Session)
	r.Tx.SetUsername(r.Person.Username)
	http.SetCookie(r, &http.Cookie{
		Name:    "auth",
		Value:   string(r.Session.Token),
		Path:    "/",
		Expires: r.Session.Expires,
	})
}

// RandomToken returns a random token string, used for various purposes.
func RandomToken() string {
	var (
		tokenb [24]byte
		err    error
	)
	if _, err = rand.Read(tokenb[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(tokenb[:])
}

// CreateSession creates a session for the person in the request, and sets a
// response cookie with the session token.
func CreateSession(r *Request) {
	r.Session = &model.Session{
		Token:   model.SessionToken(RandomToken()),
		Person:  r.Person,
		Expires: time.Now().Add(sessionExpiration),
	}
	r.Tx.CreateSession(r.Session)
	http.SetCookie(r, &http.Cookie{
		Name:    "auth",
		Value:   string(r.Session.Token),
		Path:    "/",
		Expires: r.Session.Expires,
	})
}
