package auth

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"serv.rothskeller.net/portal/model"
	"serv.rothskeller.net/portal/util"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

// Lifetime of a remember-me request.  (A year, more or less.)
const rememberMeExpiration = 365 * 24 * time.Hour

// GetLogin displays the login page.
func GetLogin(r *util.Request) error {
	util.RenderPage(r, &util.Page{
		Title: "SERV Portal",
		BodyData: &loginData{
			Target: r.FormValue("target"),
		},
	}, template.Must(template.New("login").Parse(loginTemplate)))
	return nil
}

// PostLogin handles a login request.
func PostLogin(r *util.Request) error {
	var (
		person   *model.Person
		email    = strings.ToLower(strings.TrimSpace(r.FormValue("email")))
		password = r.FormValue("password")
		remember = r.FormValue("remember") != ""
		target   = r.FormValue("target")
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
	if target == "" {
		target = "/"
	}
	http.Redirect(r, r.Request, target, http.StatusSeeOther)
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
		r.Tx.Commit()
	}
	// Show the login page again, with a login failure message.
	util.RenderPage(r, &util.Page{
		Title: "SERV Portal",
		BodyData: &loginData{
			Email:  r.FormValue("email"),
			Target: r.FormValue("target"),
			Failed: true,
		},
	}, template.Must(template.New("login").Parse(loginTemplate)))
	return nil
}

type loginData struct {
	Email  string
	Target string
	Failed bool
}

const loginTemplate = `{{ define "body" -}}
<div id="login-top">
  <div id="login-banner">
    Please log in.
  </div>
  <div id="login-forserv">
    This web site is for SERV volunteers only. If you are interested in joining
    one of the SERV volunteer organizations, please visit Sunnyvale’s <a href="https://sunnyvale.ca.gov/government/safety/emergency.htm">emergency response&nbsp;page</a>.
  </div>
  <form id="login-form" method="POST">
    {{- if .Target }}<input type="hidden" name="target" value="{{ .Target }}">{{ end }}
    <div id="login-email-row">
      <label for="login-email" class="login-form-label">Email address</label>
      <input id="login-email" name="email" class="form-control" autocorrect="off" autocapitalize="none" {{ if .Email }}value="{{ .Email }}"{{ else }}autofocus{{ end }}>
    </div>
    <div id="login-password-row">
      <label for="login-password" class="login-form-label">Password</label>
      <input id="login-password" name="password" class="form-control" type="password"{{ if .Email }} autofocus{{ end }}>
    </div>
    <div id="login-remember-row" class="custom-control custom-switch">
      <input id="login-remember" name="remember" type="checkbox" class="custom-control-input">
      <label for="login-remember" class="custom-control-label">Remember me</label>
    </div>
    <div id="login-submit-row">
      <button type="submit" class="btn btn-primary">Log in</button>
    </div>
  </form>
  {{- if .Failed }}
    <div id="login-failed">
      Login incorrect. Please try again.
    </div>
  {{ end }}
  <div id="login-reset">
    <a class="btn btn-secondary" href="/login/reset">Reset my password</a>
  </div>
</div>
{{- end }}`
