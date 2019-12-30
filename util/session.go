package util

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"time"

	"serv.rothskeller.net/portal/model"
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
	if c, err = r.Cookie("auth"); err == nil {
		r.Tx.DeleteExpiredSessions()
		r.Person = r.Tx.FetchPersonBySessionToken(c.Value)
	}
	if r.Person == nil {
		if c, err = r.Cookie("remember"); err == nil {
			r.Tx.DeleteExpiredRememberMeTokens()
			r.Person = r.Tx.FetchPersonByRememberMeToken(c.Value)
			if r.Person != nil {
				CreateSession(r)
			}
		}
	}
	if r.Person != nil {
		r.Tx.SetUsername(r.Person.Email)
	}
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
	session := &model.Session{
		Token:   model.SessionToken(RandomToken()),
		Person:  r.Person,
		Expires: time.Now().Add(sessionExpiration),
	}
	r.Tx.CreateSession(session)
	http.SetCookie(r, &http.Cookie{
		Name:    "auth",
		Value:   string(session.Token),
		Path:    "/",
		Expires: session.Expires,
	})
}
