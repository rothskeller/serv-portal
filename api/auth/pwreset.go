package auth

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/smtp"
	"time"

	"github.com/mailru/easyjson/jwriter"
	"rothskeller.net/serv/config"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// Time during which the password reset sequence must be completed.
const pwresetThreshold = time.Hour

// PostPasswordReset handles POST /api/password-reset requests.
func PostPasswordReset(r *util.Request) error {
	var (
		person   *model.Person
		body     bytes.Buffer
		emails   []string
		username = r.FormValue("username")
	)
	if person = r.Tx.FetchPersonByUsername(username); person == nil {
		return nil
	}
	if !IsEnabled(r, person) {
		return nil
	}
	for _, e := range person.Emails {
		if !e.Bad {
			emails = append(emails, e.Email)
		}
	}
	if len(emails) == 0 {
		return nil
	}
	r.Tx.DeleteSessionsForPerson(person, "")
	person.PWResetToken = util.RandomToken()
	person.BadLoginCount = 0
	person.PWResetTime = time.Now()
	r.Tx.SavePerson(person)
	r.Tx.Commit()
	fmt.Fprintf(&body, "From: %s\r\nTo: ", config.Get("fromEmail"))
	for i, e := range emails {
		if i != 0 {
			body.WriteString(", ")
		}
		fmt.Fprintf(&body, "%s <%s>", person.FormalName, e)
	}
	fmt.Fprintf(&body, "\r\nSubject: SERV Portal Password Reset\r\n\r\nGreetings, %s,\r\n\r\nTo reset your password on the SERV Portal, click this link:\r\n    %s/password-reset/%s\r\n\r\nIf you have any problems, reply to this email. If you did not request a password reset, you can safely ignore this email.\r\n",
		person.InformalName, config.Get("siteURL"), person.PWResetToken)
	if err := smtp.SendMail(
		config.Get("smtpServer"),
		&loginAuth{config.Get("smtpUsername"), config.Get("smtpPassword")},
		config.Get("fromAddr"),
		append(emails, config.Get("adminEmail")),
		body.Bytes(),
	); err != nil {
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
	out.RawByte(',')
	out.String(person.CallSign)
	out.RawByte(',')
	out.String(person.Username)
	for _, a := range person.Addresses {
		out.RawByte(',')
		out.String(a.Address)
		out.RawByte(',')
		out.String(a.City)
		out.RawByte(',')
		out.String(a.State)
		out.RawByte(',')
		out.String(a.Zip)
	}
	for _, e := range person.Emails {
		out.RawByte(',')
		out.String(e.Email)
	}
	for _, p := range person.Phones {
		out.RawByte(',')
		out.String(p.Phone)
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
	if !StrongPassword(r, person, password) {
		return errors.New("bad password")
	}
	SetPassword(r, person, password)
	person.PWResetToken = ""
	r.Tx.SavePerson(person)
	r.Person = person
	util.CreateSession(r)
	r.Tx.Commit()
	return GetLogin(r)
}

type loginAuth struct{ username, password string }

func (a *loginAuth) Start(server *smtp.ServerInfo) (string, []byte, error) {
	return "LOGIN", []byte(a.username), nil
}

func (a *loginAuth) Next(fromServer []byte, more bool) ([]byte, error) {
	if more {
		switch string(fromServer) {
		case "Username:":
			return []byte(a.username), nil
		case "Password:":
			return []byte(a.password), nil
		default:
			return nil, errors.New("Unknown fromServer")
		}
	}
	return nil, nil
}
