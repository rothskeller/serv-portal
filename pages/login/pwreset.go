package login

import (
	"bytes"
	"fmt"
	"net/http"
	"net/mail"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/session"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/sendmail"
)

// pwresetThreshold is how long a password reset token is good for.
const pwresetThreshold = time.Hour

// HandlePWReset handles /password-reset requests.
func HandlePWReset(r *request.Request) {
	if r.Method == http.MethodGet {
		ui.Page(r, nil, ui.PageOpts{}, func(main *htmlb.Element) {
			main.A("class=login")
			main.E("div class=loginBanner").T(r.Loc("Password Reset"))
			main.E("div class=loginExplain").T(r.Loc("To reset your password, please enter your email address.  If it’s one we have on file, we’ll send a password reset link to it."))
			form := main.E("form class='form form-centered form-2col loginForm' method=POST up-target=body")
			row := form.E("div class=formRow")
			row.E("label for=pwresetEmail class=formLabel").T(r.Loc("Email address"))
			row.E("input type=text id=pwresetEmail name=email autocomplete=email autocapitalize=none inputmode=email autofocus")
			row = form.E("div class='formRow-3col loginSubmit'")
			row.E("input type=submit class='sbtn sbtn-primary' value=%s", r.Loc("Reset Password"))
		})
		return
	}
	const personFields = person.FID | person.FInformalName | person.FPrivLevels | person.FEmail | person.FEmail2 | person.FPWResetToken | person.FPWResetTime | person.FCallSign
	var (
		p      *person.Person
		body   bytes.Buffer
		emails []string
		email  = r.FormValue("email")
	)
	if p = person.WithEmail(r, email, personFields); p == nil {
		goto RESPOND // email not recognized
	}
	if held, _ := personrole.PersonHasRole(r, p.ID(), role.Disabled); held {
		goto RESPOND // person is disabled
	}
	if !p.HasPrivLevel(0, enum.PrivStudent) {
		goto RESPOND // person belongs to no orgs
	}
	emails = append(emails, p.Email())
	if p.Email2() != "" {
		emails = append(emails, p.Email2())
	}
	r.Transaction(func() {
		session.DeleteForPerson(r, p, "")
		up := p.Updater()
		if up.PWResetToken == "" || time.Now().After(up.PWResetTime.Add(pwresetThreshold)) {
			up.PWResetToken = util.RandomToken()
		}
		up.PWResetTime = time.Now()
		p.Update(r, up, person.FPWResetToken|person.FPWResetTime)
	})
	fmt.Fprintf(&body, "From: %s\r\nTo: ", config.Get("fromEmail"))
	for i, e := range emails {
		if i != 0 {
			body.WriteString(", ")
		}
		fmt.Fprint(&body, &mail.Address{Name: p.InformalName(), Address: e})
	}
	fmt.Fprintf(&body, r.Loc("\r\nSubject: SunnyvaleSERV.org Password Reset\r\n\r\nGreetings, %s,\r\n\r\nTo reset your password on SunnyvaleSERV.org, click this link:\r\n    %s/password-reset/%s\r\n\r\nIf you have any problems, reply to this email. If you did not request a password reset, you can safely ignore this email.\r\n"),
		p.InformalName(), config.Get("siteURL"), p.PWResetToken())
	if err := sendmail.SendMessage(config.Get("fromAddr"), append(emails, config.Get("adminEmail")), body.Bytes()); err != nil {
		panic(err)
	}
RESPOND:
	ui.Page(r, nil, ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=login")
		main.E("div class=loginBanner").T(r.Loc("Password Reset"))
		exp := main.E("div class=loginExplain")
		exp.E("p").T(r.Loc("We have sent a password reset link to the email address you provided. It is valid for one hour. Please check your email and follow the link we sent to reset your password."))
		exp.E("p").R(r.Loc("If you do not receive an email with a password reset link, it may be that the email address you provided is not the one we have on file for you. Contact <a href=\"mailto:admin@sunnyvaleserv.org\">admin@SunnyvaleSERV.org</a> for assistance."))
	})
}

// HandlePWResetToken handles /password-reset/${token} requests.
func HandlePWResetToken(r *request.Request, token string) {
	const personFields = person.FID | person.FPrivLevels | person.FPWResetToken | person.FPWResetTime | auth.StrongPasswordPersonFields
	var p *person.Person

	if p = person.WithPWResetToken(r, token, personFields); p == nil {
		goto INVALID // unknown token
	}
	if held, _ := personrole.PersonHasRole(r, p.ID(), role.Disabled); held {
		goto INVALID // person is disabled
	}
	if time.Now().After(p.PWResetTime().Add(pwresetThreshold)) {
		goto INVALID // token has expired
	}
	if r.Method == http.MethodPost {
		postPWResetToken(r, p)
	} else {
		getPWResetToken(r, p)
	}
	return
INVALID:
	ui.Page(r, nil, ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=login")
		main.E("div class=loginBanner").T(r.Loc("Password Reset"))
		main.E("div class=loginExplain").T(r.Loc("This password reset link is invalid or has expired."))
		main.E("div class=loginSubmit").E("a class='sbtn sbtn-primary' href=/password-reset up-target=body").T(r.Loc("Try Again"))
	})
}

func getPWResetToken(r *request.Request, p *person.Person) {
	ui.Page(r, nil, ui.PageOpts{}, func(main *htmlb.Element) {
		main.A("class=login")
		main.E("div class=loginBanner").T(r.Loc("Password Reset"))
		form := main.E("form class='form form-centered form-2col loginForm pwResetForm' method=POST up-target=body")
		row := form.E("div class=formRow")
		row.E("label for=pwresetNewPassword").T(r.Loc("New Password"))
		row.E("s-password id=pwresetNewPassword name=newpwd hints=%s", strings.Join(auth.StrongPasswordHints(p), ","))
		row = form.E("div class='formRow-3col loginSubmit'")
		row.E("input type=submit class='sbtn sbtn-primary' value=%s", r.Loc("Reset Password"))
	})
}

func postPWResetToken(r *request.Request, p *person.Person) {
	var newpwd = r.FormValue("newpwd")
	if newpwd == "" || !auth.StrongPassword(p, newpwd) {
		getPWResetToken(r, p) // somehow a weak password got through
		return
	}
	r.Transaction(func() {
		auth.SetPassword(r, p, newpwd)
		auth.CreateSession(r, p, false)
	})
	http.Redirect(r, r.Request, "/", http.StatusSeeOther)
}
