package personedit

import (
	"bytes"
	"fmt"
	"net/http"
	"net/mail"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/sendmail"
)

const pwResetPersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | person.FPassword | person.FBadLoginCount | person.FBadLoginTime | person.FPWResetToken | person.FPWResetTime | person.FEmail | person.FEmail2

// HandlePWReset handles requests for /people/$id/pwreset.
func HandlePWReset(r *request.Request, idstr string) {
	var (
		user *person.Person
		p    *person.Person
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), pwResetPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() == p.ID() || p.ID() == person.AdminID || !user.HasPrivLevel(0, enum.PrivLeader) || p.Email() == "" {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		postPWReset(r, user, p)
	} else {
		getPWReset(r, p)
	}
}

func getPWReset(r *request.Request, p *person.Person) {
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewPassword")
	form.E("div class='formTitle formTitle-warning'>Reset Password")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("div class=formRow-3col").
		R("This will replace ").
		T(p.InformalName()).
		R("’s current password with a new, randomly generated one, and send it to them by email. You should do this only if they have explicitly asked you to. Are you sure you want to reset ").
		T(p.InformalName()).
		R("’s password?")
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-warning' value='Reset Password'")
}

func postPWReset(r *request.Request, user, p *person.Person) {
	var (
		password string
		emails   []string
		body     bytes.Buffer
		crlf     = sendmail.NewCRLFWriter(&body)
	)
	password = auth.RandomPassword()
	auth.SetPassword(r, p, password)
	emails = append(emails, p.Email())
	fmt.Fprintf(crlf, "To: %s\n", (&mail.Address{Name: p.InformalName(), Address: p.Email()}).String())
	if p.Email2() != "" {
		emails = append(emails, p.Email2())
		fmt.Fprintf(crlf, "To: %s\n", (&mail.Address{Name: p.InformalName(), Address: p.Email2()}).String())
	}
	fmt.Fprintf(crlf, r.Loc(`From: %s
Subject: SunnyvaleSERV.org Password Reset
Content-Type: text/plain; charset=utf8

Hello, %s,

%s has reset the password for your account on SunnyvaleSERV.org.  Your new login information is:

    Email:    %s
    Password: %s

This password is three words chosen randomly from a dictionary — a method that generally produces a very secure and often memorable password.  If the resulting phrase has any meaning, it’s unintentional coincidence.

You can change this password by logging into SunnyvaleSERV.org and clicking the “Change Password” button on your Profile page.  If you have any questions, just reply to this email.

Regards,
SunnyvaleSERV.org
`), config.Get("fromEmail"), p.InformalName(), user.InformalName(), p.Email(), password)
	if err := sendmail.SendMessage(config.Get("fromAddr"), emails, body.Bytes()); err != nil {
		panic(err)
	}
	personview.Render(r, user, p, person.ViewFull, "password")
}
