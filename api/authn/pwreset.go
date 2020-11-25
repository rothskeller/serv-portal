package authn

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/sendmail"
)

// Time during which the password reset sequence must be completed.
const pwresetThreshold = time.Hour

// PostPasswordReset handles POST /api/password-reset requests.
func PostPasswordReset(r *util.Request) error {
	var (
		person *model.Person
		body   bytes.Buffer
		emails []string
		email  = r.FormValue("email")
	)
	if person = r.Tx.FetchPersonByUsername(email); person == nil {
		return nil
	}
	if person.Roles[model.DisabledUser] {
		return nil // person is disabled
	}
	if !person.HasPrivLevel(model.PrivStudent) {
		return nil // person belongs to no orgs
	}
	r.Tx.WillUpdatePerson(person)
	emails = append(emails, person.Email)
	if person.Email2 != "" {
		emails = append(emails, person.Email2)
	}
	r.Tx.DeleteSessionsForPerson(person, "")
	person.PWResetToken = util.RandomToken()
	person.BadLoginCount = 0
	person.PWResetTime = time.Now()
	r.Tx.UpdatePerson(person)
	r.Tx.Commit()
	fmt.Fprintf(&body, "From: %s\r\nDate: %s\r\nTo: ", config.Get("fromEmail"), time.Now().Format(time.RFC1123Z))
	for i, e := range emails {
		if i != 0 {
			body.WriteString(", ")
		}
		fmt.Fprint(&body, &mail.Address{Name: person.InformalName, Address: e})
	}
	fmt.Fprintf(&body, "\r\nSubject: SunnyvaleSERV.org Password Reset\r\n\r\nGreetings, %s,\r\n\r\nTo reset your password on SunnyvaleSERV.org, click this link:\r\n    %s/password-reset/%s\r\n\r\nIf you have any problems, reply to this email. If you did not request a password reset, you can safely ignore this email.\r\n",
		person.InformalName, config.Get("siteURL"), person.PWResetToken)
	if err := sendmail.SendMessage(config.Get("fromAddr"), append(emails, config.Get("adminEmail")), body.Bytes()); err != nil {
		panic(err)
	}
	return nil
}

// GetPasswordResetToken handles GET /api/password-reset/$token requests.
func GetPasswordResetToken(r *util.Request, token string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPersonByPWResetToken(token); person == nil || time.Since(person.PWResetTime) > pwresetThreshold {
		time.Sleep(5 * time.Second)
		return util.HTTPError(http.StatusConflict, "The password reset token is invalid or expired.")
	}
	r.Tx.Commit()
	out.RawByte('[')
	for i, h := range SERVPasswordHints {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(h)
	}
	out.RawByte(',')
	out.String(person.InformalName)
	out.RawByte(',')
	out.String(person.FormalName)
	if person.CallSign != "" {
		out.RawByte(',')
		out.String(person.CallSign)
	}
	if person.HomeAddress.Address != "" {
		out.RawByte(',')
		out.String(person.HomeAddress.Address)
	}
	if person.MailAddress.Address != "" {
		out.RawByte(',')
		out.String(person.MailAddress.Address)
	}
	if person.WorkAddress.Address != "" {
		out.RawByte(',')
		out.String(person.WorkAddress.Address)
	}
	if person.CellPhone != "" {
		out.RawByte(',')
		out.String(person.CellPhone)
	}
	if person.HomePhone != "" {
		out.RawByte(',')
		out.String(person.HomePhone)
	}
	if person.WorkPhone != "" {
		out.RawByte(',')
		out.String(person.WorkPhone)
	}
	out.RawByte(']')
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostPasswordResetToken handles POST /api/password-reset/$token requests.
func PostPasswordResetToken(r *util.Request, token string) error {
	var (
		person   *model.Person
		password = r.FormValue("password")
	)
	if person = r.Tx.FetchPersonByPWResetToken(token); person == nil || time.Since(person.PWResetTime) > pwresetThreshold {
		return util.HTTPError(http.StatusConflict, "The password reset token is invalid or expired.")
	}
	if !StrongPassword(person, password) {
		return errors.New("bad password")
	}
	r.Tx.WillUpdatePerson(person)
	SetPassword(r, person, password)
	person.PWResetToken = ""
	r.Tx.UpdatePerson(person)
	r.Person = person
	r.Auth.SetMe(person)
	util.CreateSession(r, false)
	r.Tx.Commit()
	return GetLogin(r)
}
