package login

import (
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

// HandleLogin handles GET and POST /login and /login/* requests.
func HandleLogin(r *request.Request) {
	const personFields = person.FID | person.FInformalName | person.FBadLoginCount | person.FBadLoginTime | person.FPrivLevels | person.FPassword | person.FCallSign
	var (
		email      string
		remember   bool
		statusCode = http.StatusOK
	)
	if auth.SessionUser(r, 0, false) != nil { // Already logged in.
		if len(r.Path) > 6 { // Redirect path in URL.
			http.Redirect(r, r.Request, r.Path[6:], http.StatusSeeOther)
		} else {
			http.Redirect(r, r.Request, "/", http.StatusSeeOther)
		}
		return
	}
	if r.Method == http.MethodPost {
		var (
			p        *person.Person
			password = r.FormValue("password")
		)
		email = r.FormValue("email")
		remember = r.FormValue("remember") != ""
		// Check that the login is valid.
		if p = person.WithEmail(r, email, personFields); p == nil {
			goto FAIL // no person with that username
		}
		if p.ID() != person.AdminID { // admin cannot be disabled or locked out
			if p.BadLoginCount() >= maxBadLogins && time.Now().Before(p.BadLoginTime().Add(badLoginThreshold)) {
				goto FAIL // locked out
			}
			if held, _ := personrole.PersonHasRole(r, p.ID(), role.Disabled); held {
				goto FAIL // person is disabled
			}
		} else { // admin can not be remembered
			remember = false
		}
		if !auth.CheckPassword(r, p, password) {
			goto FAIL // password mismatch
		}
		// The login is valid.  Record it and create a session.
		r.Transaction(func() {
			if p.BadLoginCount() > 0 {
				up := p.Updater()
				up.BadLoginCount = 0
				up.BadLoginTime = time.Time{}
				p.Update(r, up, person.FBadLoginCount|person.FBadLoginTime)
			}
			auth.CreateSession(r, p, remember)
		})
		if len(r.Path) > 6 { // Redirect path in URL.
			http.Redirect(r, r.Request, r.Path[6:], http.StatusSeeOther)
		} else {
			http.Redirect(r, r.Request, "/", http.StatusSeeOther)
		}
		return
	FAIL:
		statusCode = http.StatusUnprocessableEntity
		if p != nil {
			// Record the bad login attempt.
			r.Transaction(func() {
				up := p.Updater()
				if time.Now().Before(up.BadLoginTime.Add(badLoginThreshold)) {
					up.BadLoginCount++
				} else {
					up.BadLoginCount = 1
				}
				up.BadLoginTime = time.Now()
				p.Update(r, up, person.FBadLoginCount|person.FBadLoginTime)
			})
		}
	}
	ui.Page(r, nil, ui.PageOpts{
		Title:      "Login",
		Banner:     "Sunnyvale SERV",
		StatusCode: statusCode,
	}, func(main *htmlb.Element) {
		main.A("class=login")
		main.E("div class=loginBanner>Please log in.")
		main.E("div class=loginExplain>This web site is for SERV volunteers only. If you are interested in joining one of the SERV volunteer organizations, send us email at <a href=mailto:serv@sunnyvaleserv.org>SERV@SunnyvaleSERV.org</a>.")
		main.E("div class=loginBrowserwarn>Your browser is out of date and lacks features needed by this web site. The site may not look or behave correctly.")
		form := main.E("form class='form form-centered form-2col loginForm' method=POST up-target=body up-fail-target=form")

		// Email row.
		row := form.E("div class=formRow")
		row.E("label for=loginEmail class=formLabel>Email address")
		row.E("input name=email type=text id=loginEmail autocomplete=email autocapitalize=none inputmode=email value=%s",
			email, statusCode == http.StatusOK, "autofocus")

		// Password row.
		row = form.E("div class=formRow")
		row.E("label for=loginPassword class=formLabel>Password")
		row.E("input name=password type=password id=loginPassword autocomplete=password autocapitalize=none",
			statusCode != http.StatusOK, "autofocus")

		// Remember row.
		row = form.E("div class=formRow")
		row.Element("input type=checkbox class='s-check formInput' name=remember label='Remember me'", remember, "checked")

		// Submit button row.
		row = form.E("div class='formRow-3col loginSubmit'")
		row.E("input type=submit class='sbtn sbtn-primary' value=%s", "Log in")

		// Failure notice.
		if statusCode != http.StatusOK {
			form.E("div class='formRow-3col loginFailed'>Login incorrect. Please try again.")
		}

		// Reset password link.
		main.E("div class=loginReset").E("a href=/password-reset up-follow>Reset my password")

		// Website information link.
		main.E("div class=loginAbout").E("a href=/about up-follow>Web Site Information")
	})
}

// HandleLogout handles GET /logout requests.
func HandleLogout(r *request.Request) {
	if user := auth.SessionUser(r, person.FID|person.FInformalName, false); user != nil {
		r.Transaction(func() {
			auth.DeleteSession(r, user)
		})
	}
	http.Redirect(r, r.Request, "/", http.StatusSeeOther)
}
