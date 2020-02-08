package email

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/authz"
	"sunnyvaleserv.org/portal/config"
	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// PostEmail handles POST /api/emails/$id requests.
func PostEmail(r *util.Request, idstr string) error {
	var (
		msg *model.EmailMessage
	)
	if msg = r.Tx.FetchEmailMessage(model.EmailMessageID(util.ParseID(idstr))); msg == nil {
		return util.NotFound
	}
	switch r.FormValue("action") {
	case "accept":
		if msg.Type != model.EmailModerated {
			return errors.New("can't accept message that isn't waiting for moderation")
		}
		SendMessage(r.Tx, r.Auth, msg)
		msg.Type = model.EmailSent
		msg.Attention = false
		r.Tx.UpdateEmailMessage(msg)
		r.Tx.Commit()
		return nil
	case "forward":
		ForwardMessage(r, msg)
		r.Tx.Commit()
		return nil
	default:
		return errors.New("invalid action")
	}
}

// SendMessage sends an email message to the groups to which it's addressed.
func SendMessage(tx *db.Tx, auth *authz.Authorizer, email *model.EmailMessage) {
	var (
		raw     []byte
		msg     *mail.Message
		body    []byte
		root    *messagePart
		client  *smtp.Client
		tlsconf tls.Config
		login   loginAuth
		err     error
	)
	raw = tx.FetchEmailMessageBody(email.ID)
	if msg, err = mail.ReadMessage(bytes.NewReader(raw)); err != nil {
		panic(err)
	}
	body, _ = ioutil.ReadAll(msg.Body)
	root, _ = makeMessagePart(textproto.MIMEHeader(msg.Header), body)
	for hdr := range msg.Header {
		switch hdr {
		case "Cc", "Content-Transfer-Encoding", "Content-Type", "Date", "In-Reply-To", "Message-Id", "Mime-Version",
			"Organization", "Reply-To", "Subject", "To":
			break
		case "From":
			if msg.Header["Reply-To"] == nil {
				msg.Header["Reply-To"] = msg.Header["From"]
			}
			fallthrough
		default:
			delete(msg.Header, hdr)
		}
	}
	if client, err = smtp.Dial(config.Get("sendGridServerPort")); err != nil {
		goto SEND_ERROR
	}
	defer client.Close()
	tlsconf.ServerName = config.Get("sendGridServer")
	if err = client.StartTLS(&tlsconf); err != nil {
		goto SEND_ERROR
	}
	login.username = config.Get("sendGridUsername")
	login.password = config.Get("sendGridPassword")
	if err = client.Auth(&login); err != nil {
		goto SEND_ERROR
	}
	for _, group := range email.Groups {
		if err = sendMessageToGroup(tx, auth, client, email, msg, root, group); err != nil {
			goto SEND_ERROR
		}
	}
	email.Type = model.EmailSent
	email.Attention = false
	tx.UpdateEmailMessage(email)
	return

SEND_ERROR:
	email.Type = model.EmailSendFailed
	email.Error = "Send Failed: " + err.Error()
	email.Attention = true
}

func sendMessageToGroup(
	tx *db.Tx, auth *authz.Authorizer, client *smtp.Client, email *model.EmailMessage, msg *mail.Message, root *messagePart,
	gid model.GroupID,
) error {
	group := auth.FetchGroup(gid)
	pids := make(map[model.PersonID]bool)
	for _, pid := range auth.PeopleG(gid) {
		pids[pid] = true
	}
	for _, rid := range auth.RolesAG(model.PrivBCC, gid) {
		for _, pid := range auth.PeopleR(rid) {
			pids[pid] = true
		}
	}
PEOPLE:
	for pid := range pids {
		for _, ne := range group.NoEmail {
			if pid == ne {
				continue PEOPLE
			}
		}
		person := tx.FetchPerson(pid)
		if person.NoEmail {
			continue
		}
		if err := sendMessageToPerson(client, email, msg, root, group, person); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToPerson(
	client *smtp.Client, email *model.EmailMessage, msg *mail.Message, root *messagePart, group *model.Group,
	person *model.Person,
) error {
	if person.Email != "" {
		if err := sendMessageToEmail(client, email, msg, root, group, person, person.Email); err != nil {
			return err
		}
	}
	if person.Email2 != "" {
		if err := sendMessageToEmail(client, email, msg, root, group, person, person.Email2); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToEmail(
	client *smtp.Client, email *model.EmailMessage, msg *mail.Message, root *messagePart, group *model.Group,
	person *model.Person, address string,
) error {
	var (
		buf bytes.Buffer
		wr  io.WriteCloser
		err error
	)
	emitFrom(&buf, email, group)
	fmt.Fprintf(&buf, "Sender: %s <%s@sunnyvaleserv.org>\r\n", quoteIfNeeded(group.Name), group.Email)
	fmt.Fprintf(&buf, "Errors-To: admin@sunnyvaleserv.org\r\n")
	emitHeaders(&buf, msg.Header)
	rewrite(&buf, root, group.Email, person.InformalName, address)
	if err = client.Mail(group.Email + "@sunnyvaleserv.org"); err != nil {
		return err
	}
	if err = client.Rcpt(address); err != nil {
		return err
	}
	if wr, err = client.Data(); err != nil {
		return err
	}
	if _, err = wr.Write(buf.Bytes()); err != nil {
		return err
	}
	if err = wr.Close(); err != nil {
		return err
	}
	return nil
}

func emitFrom(buf io.Writer, email *model.EmailMessage, group *model.Group) {
	var from = email.From
	if idx := strings.IndexByte(from, '@'); idx >= 0 {
		from = from[:idx]
	}
	fmt.Fprintf(buf, "From: %s via %s <%s@sunnyvaleserv.org>\r\n", quoteIfNeeded(from), quoteIfNeeded(group.Name), group.Email)
}

func emitHeaders(buf io.Writer, headers mail.Header) {
	for h, va := range headers {
		for _, v := range va {
			emitHeader(buf, h, v)
		}
	}
	fmt.Fprint(buf, "\r\n")
}

func emitHeader(buf io.Writer, name string, value string) {
	str := name + ": " + value
	for len(str) > 78 {
		idx := strings.LastIndex(str[:78], ", ")
		if idx >= 0 {
			fmt.Fprint(buf, str[:idx], ",\r\n ")
			str = str[idx+2:]
			continue
		}
		idx = strings.LastIndex(str[:78], " ")
		if idx >= 0 {
			fmt.Fprint(buf, str[:idx], "\r\n ")
			str = str[idx+1:]
			continue
		}
		idx = strings.IndexByte(str, ' ')
		if idx >= 0 {
			fmt.Fprint(buf, str[:idx], "\r\n ")
			str = str[idx+1:]
			continue
		}
		fmt.Fprint(buf, str, "\r\n")
		return
	}
	if len(str) != 0 {
		fmt.Fprint(buf, str, "\r\n")
	}
}

var unquotedRE = regexp.MustCompile("^[-a-zA-Z0-9!#$%&'*+/=?^_`{}|~.]+$")

func quoteIfNeeded(s string) string {
	if unquotedRE.MatchString(s) {
		return s
	}
	return `"` + strings.Replace(s, `"`, `\"`, -1) + `"`
}

// ForwardMessage forwards an email message to the caller's primary email
// address.
func ForwardMessage(r *util.Request, email *model.EmailMessage) {
	var (
		raw      []byte
		msg      bytes.Buffer
		login    loginAuth
		boundary string
	)
	login.username = config.Get("sendGridUsername")
	login.password = config.Get("sendGridPassword")
	raw = r.Tx.FetchEmailMessageBody(email.ID)
	boundary = util.RandomToken()
	fmt.Fprintf(&msg, "From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>\r\nTo: %s <%s>\r\nSubject: Forward of moderated email\r\nContent-Type: multipart/mixed; boundary=%s\r\n\r\n\r\n--%[3]s\r\nContent-Type: text/plain\r\n\r\nThe attached email was forwarded at your request from the SunnyvaleSERV.org portal.\r\n\r\n--%[3]s\r\nContent-Type: message/rfc822\r\n\r\n", quoteIfNeeded(r.Person.InformalName), r.Person.Email, boundary)
	msg.Write(raw)
	fmt.Fprintf(&msg, "\r\n--%s--\r\n", boundary)
	smtp.SendMail(config.Get("sendGridServerPort"), &login, "admin@sunnyvaleserv.org", []string{r.Person.Email}, msg.Bytes())
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
