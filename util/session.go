package util

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util/config"
)

// Lifetime of a login session.
const sessionExpiration = time.Hour

// Lifetime of a remember-me login session.  Note that some browsers won't allow
// any cookie to last this long — Chrome appears to have a six month limit — so
// the remember-me may not be forever.
const rememberExpiration = 10 * 365 * 24 * time.Hour

// Forbidden is the error returned when the calling session lacks the privileges
// needed for the call it issued.
var Forbidden = HTTPError(http.StatusForbidden, "403 Forbidden")

// ValidateSession decodes the auth and/or remember tokens in the request, if
// any, and sets the Person field of the request appropriately.  If the session
// authorization is not valid, the Person field is left unchanged (i.e., nil).
func ValidateSession(r *Request) {
	var (
		c      *http.Cookie
		newexp time.Time
		err    error
	)
	if c, err = r.Cookie("auth"); err != nil {
		return
	}
	r.Tx.DeleteExpiredSessions()
	if r.Session = r.Tx.FetchSession(model.SessionToken(c.Value)); r.Session == nil {
		return
	}
	if r.Method != "GET" && r.Session.CSRF != model.CSRFToken(r.Request.Header.Get("X-CSRF-Token")) {
		r.Session = nil
		return
	}
	r.Person = r.Session.Person
	newexp = time.Now().Add(sessionExpiration)
	if newexp.After(r.Session.Expires) {
		r.Session.Expires = newexp
		r.Tx.UpdateSession(r.Session)
		http.SetCookie(r, &http.Cookie{
			Name:     "auth",
			Value:    string(r.Session.Token),
			Path:     "/",
			Expires:  r.Session.Expires,
			HttpOnly: true,
			Secure:   strings.HasPrefix(config.Get("siteURL"), "https://"),
			SameSite: http.SameSiteLaxMode,
		})
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
func CreateSession(r *Request, remember bool) {
	r.Session = &model.Session{
		Token:  model.SessionToken(RandomToken()),
		Person: r.Person,
		CSRF:   model.CSRFToken(RandomToken()),
	}
	if remember {
		r.Session.Expires = time.Now().Add(rememberExpiration)
	} else {
		r.Session.Expires = time.Now().Add(sessionExpiration)
	}
	r.Tx.CreateSession(r.Session)
	http.SetCookie(r, &http.Cookie{
		Name:     "auth",
		Value:    string(r.Session.Token),
		Path:     "/",
		Expires:  r.Session.Expires,
		HttpOnly: true,
		Secure:   strings.HasPrefix(config.Get("siteURL"), "https://"),
		SameSite: http.SameSiteLaxMode,
	})
}
