package email

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/mail"
	"net/textproto"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/sendmail"
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
		SendMessage(r.Tx, msg)
		msg.Type = model.EmailSent
		msg.Attention = false
		r.Tx.UpdateEmailMessage(msg)
		r.Tx.Commit()
		return nil
	case "sendToMe":
		SendMessageToMe(r, msg)
		return nil
	case "discard":
		msg.Attention = false
		r.Tx.UpdateEmailMessage(msg)
		r.Tx.Commit()
		return nil
	default:
		return errors.New("invalid action")
	}
}

// SendMessage sends an email message to the groups to which it's addressed.
func SendMessage(tx *store.Tx, email *model.EmailMessage) {
	var (
		raw    []byte
		msg    *mail.Message
		body   []byte
		root   *messagePart
		mailer *sendmail.Mailer
		err    error
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
	if mailer, err = sendmail.OpenMailer(); err != nil {
		goto SEND_ERROR
	}
	defer mailer.Close()
	for _, group := range email.Groups {
		if err = sendMessageToGroup(tx, mailer, email, msg, root, group); err != nil {
			goto SEND_ERROR
		}
	}
	email.Type = model.EmailSent
	email.Attention = false
	return

SEND_ERROR:
	email.Type = model.EmailSendFailed
	email.Error = "Send Failed: " + err.Error()
	email.Attention = true
}

func sendMessageToGroup(
	tx *store.Tx, mailer *sendmail.Mailer, email *model.EmailMessage, msg *mail.Message, root *messagePart, gid model.GroupID,
) error {
	group := tx.Authorizer().FetchGroup(gid)
	pids := make(map[model.PersonID]bool)
	for _, pid := range tx.Authorizer().PeopleG(gid) {
		pids[pid] = true
	}
	for _, rid := range tx.Authorizer().RolesAG(model.PrivBCC, gid) {
		for _, pid := range tx.Authorizer().PeopleR(rid) {
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
		if err := sendMessageToPerson(mailer, email, msg, root, group, person); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToPerson(
	mailer *sendmail.Mailer, email *model.EmailMessage, msg *mail.Message, root *messagePart, group *model.Group,
	person *model.Person,
) error {
	if person.Email != "" {
		if err := sendMessageToEmail(mailer, email, msg, root, group, person, person.Email); err != nil {
			return err
		}
	}
	if person.Email2 != "" {
		if err := sendMessageToEmail(mailer, email, msg, root, group, person, person.Email2); err != nil {
			return err
		}
	}
	return nil
}

func sendMessageToEmail(
	mailer *sendmail.Mailer, email *model.EmailMessage, msg *mail.Message, root *messagePart, group *model.Group,
	person *model.Person, address string,
) error {
	var (
		buf bytes.Buffer
	)
	emitFrom(&buf, email, group)
	fmt.Fprintf(&buf, "Sender: %s <%s@sunnyvaleserv.org>\r\n", quoteIfNeeded(group.Name), group.Email)
	fmt.Fprintf(&buf, "Errors-To: admin@sunnyvaleserv.org\r\n")
	fmt.Fprintf(&buf, "List-Unsubscribe: <https://sunnyvaleserv.org/unsubscribe/%s/%s>\r\n", person.UnsubscribeToken, group.Email)
	fmt.Fprintf(&buf, "List-Unsubscribe-Post: List-Unsubscribe=One-Click\r\n")
	emitHeaders(&buf, msg.Header)
	rewrite(&buf, root, group.Email, person.InformalName, address, person.UnsubscribeToken)
	return mailer.SendMessage(group.Email+"@sunnyvaleserv.org", []string{address}, buf.Bytes())
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

// SendMessageToMe sends an email message to the caller's primary email address
// (rather than to the list(s) it was addressed to).
func SendMessageToMe(r *util.Request, email *model.EmailMessage) {
	raw := r.Tx.FetchEmailMessageBody(email.ID)
	r.Tx.Commit()
	sendmail.SendMessage("admin@sunnyvaleserv.org", []string{r.Person.Email}, raw)
}
