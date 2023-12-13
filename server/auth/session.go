package auth

import (
	"net/http"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/session"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/request"
)

// Lifetime of a login session.
const sessionExpiration = time.Hour

// Lifetime of a remember-me login session.  Note that some browsers won't allow
// any cookie to last this long — Chrome appears to have a six month limit — so
// the remember-me may not be forever.
const rememberExpiration = 10 * 365 * 24 * time.Hour // ten years, ish

// SessionUser validates the session and returns the session user, with all of
// the specified fields fetched.  (As a convenience, ID, name, and privilege
// levels are always fetched even if not specified.)  If the session is not
// valid, it returns nil.  If respond is true and the session is not valid, it
// issues the appropriate web response before returning nil.
func SessionUser(r *request.Request, fields person.Fields, respond bool) (p *person.Person) {
	var (
		pid     person.ID
		expires time.Time
		extend  time.Time
		csrf    string
	)
	// If there's no session token, there is no user logged in.
	if r.SessionToken == "" {
		goto UNAUTHORIZED
	}
	// Get the session data, and extend the session expiration if found.
	r.Transaction(func() {
		if pid, expires, csrf = session.WithToken(r, r.SessionToken); pid != 0 {
			extend = time.Now().Add(time.Hour)
			if extend.After(expires) {
				session.Extend(r, r.SessionToken, extend)
			}
		}
	})
	if pid == 0 {
		goto UNAUTHORIZED
	}
	// Get the session user's data.  Note that we always retrieve the user's
	// ID, name, and privilege levels, even if not requested.  The name is
	// needed for session logging, and all three are needed to display the
	// page menu.
	p = person.WithID(r, pid, fields|person.FInformalName|person.FPrivLevels)
	r.LogEntry.User = p.InformalName()
	r.CSRF = csrf
	return p

UNAUTHORIZED:
	if respond {
		if strings.Contains(r.Request.Header.Get("Accept"), "text/html") {
			http.Redirect(r, r.Request, "/login"+r.Path, http.StatusSeeOther)
		} else {
			http.Error(r, "401 Unauthorized", http.StatusUnauthorized)
		}
	}
	return nil
}

// CheckCSRF returns whether the CSRF token received in the form data matches the
// CSRF token for the session.  If it does not, it issues the appropriate web
// response before returning false.
func CheckCSRF(r *request.Request, user *person.Person) bool {
	if r.Method != http.MethodPost || r.FormValue("csrf") == r.CSRF {
		return true
	}
	r.Problems().Add("invalid CSRF token")
	errpage.Forbidden(r, user)
	return false
}

// CreateSession creates a session for the person in the request, and sets a
// response cookie with the session token.
func CreateSession(r *request.Request, user *person.Person, remember bool) {
	var expires time.Time

	r.SessionToken = util.RandomToken()
	r.CSRF = util.RandomToken()
	r.LogEntry.User = user.InformalName()
	if remember {
		expires = time.Now().Add(rememberExpiration)
	} else {
		expires = time.Now().Add(sessionExpiration)
	}
	r.Transaction(func() {
		session.Create(r, user, r.SessionToken, r.CSRF, expires)
	})
	http.SetCookie(r, &http.Cookie{
		Name:     "auth",
		Value:    string(r.SessionToken),
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		Secure:   strings.HasPrefix(config.Get("siteURL"), "https://"),
		SameSite: http.SameSiteLaxMode,
	})
}

// DeleteSession ends and deletes the current session.
func DeleteSession(r *request.Request, user *person.Person) {
	session.Delete(r, r.SessionToken, user)
	http.SetCookie(r, &http.Cookie{
		Name:     "auth",
		Path:     "/",
		Expires:  time.Date(2000, 0, 0, 0, 0, 0, 0, time.UTC),
		HttpOnly: true,
		Secure:   strings.HasPrefix(config.Get("siteURL"), "https://"),
		SameSite: http.SameSiteLaxMode,
	})
}
