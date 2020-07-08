package sendmail

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/smtp"

	"sunnyvaleserv.org/portal/util/config"
)

// A Mailer is a connection to the SMTP server used for sending email.
type Mailer struct {
	client *smtp.Client
}

// OpenMailer creates a connection to the SMTP server for sending email.
func OpenMailer() (m *Mailer, err error) {
	var (
		tlsconf tls.Config
		conn    net.Conn
		login   loginAuth
	)
	m = new(Mailer)
	tlsconf.ServerName = config.Get("smtpServer")
	if conn, err = tls.Dial("tcp", config.Get("smtpServerPort"), &tlsconf); err != nil {
		return nil, fmt.Errorf("tls.Dial: %s", err)
	}
	if m.client, err = smtp.NewClient(conn, config.Get("smtpServerPort")); err != nil {
		conn.Close()
		return nil, fmt.Errorf("smtp.NewClient: %s", err)
	}
	login.username = config.Get("smtpUsername")
	login.password = config.Get("smtpPassword")
	if err = m.client.Auth(&login); err != nil {
		m.client.Close()
		return nil, fmt.Errorf("smtp.Auth: %s", err)
	}
	return m, nil
}

// SendMessage sends a single message through the Mailer.  If it returns an
// error, the Mailer is no longer usable.
func (m *Mailer) SendMessage(from string, to []string, body []byte) (err error) {
	var wr io.WriteCloser

	if err = m.client.Mail(from); err != nil {
		m.client.Close()
		return err
	}
	for _, t := range to {
		if err = m.client.Rcpt(t); err != nil {
			m.client.Close()
			return err
		}
	}
	if wr, err = m.client.Data(); err != nil {
		m.client.Close()
		return err
	}
	if _, err = wr.Write(body); err != nil {
		m.client.Close()
		return err
	}
	if err = wr.Close(); err != nil {
		m.client.Close()
		return err
	}
	return nil
}

// Close closes the connection to SendGrid.  The Mailer may not be used after
// this is called.
func (m *Mailer) Close() {
	m.client.Close()
}

// SendMessage sends a single message through SendGrid.  It is a shortcut for
// creating a Mailer, calling SendMessage on it, and then closing it.
func SendMessage(from string, to []string, body []byte) (err error) {
	var m *Mailer

	if m, err = OpenMailer(); err != nil {
		return err
	}
	if err = m.SendMessage(from, to, body); err != nil {
		return err
	}
	m.Close()
	return nil
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
