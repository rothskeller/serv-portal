// This program gets called for incoming mail to the SERV portal, with the mail
// being delivered on stdin.
//
// Mail is handled as follows:
//   - If its MessageID header is one we've seen before, ignore it.
//   - If addressed to bounce@, it gets stored in the database as a bounced
//     message, and a notification gets sent to the webmasters.  Exception: if
//     the content of the bounced message is itself such a notification, no
//     second notification is sent.
//   - If addressed to a known, unmoderated email list, and the sender is on
//     the list, it gets resent to the list.
//   - If addressed to a known, moderated list, and the sender is a moderator of
//     the list, and the sender is authenticated, it gets resent to the list.
//   - If addressed to a known email list, and not handled above, it gets stored
//     in the database as an email awaiting moderation, and a notification of
//     that gets sent to the moderators of the list.
//   - Any other message get stored in the database as an unknown message, and
//     a notification of that is sent daily to the webmasters.
package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"log"
	"net/textproto"
	"os"
	"regexp"
	"time"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
)

var removeAddressRE = regexp.MustCompile(`\s*<[^>]*>`)

func main() {
	var (
		logf *os.File
		tx   *db.Tx
		raw  []byte
		rdr  textproto.Reader
		hdr  textproto.MIMEHeader
		em   model.EmailMessage
		err  error
	)
	if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
		panic(err)
	}
	if logf, err = os.OpenFile("mail-hook.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		panic(err)
	} else {
		os.Stderr = logf
		log.SetOutput(logf)
		log.SetFlags(log.Ldate | log.Ltime)
	}
	db.Open("serv.db")
	tx = db.Begin()
	if raw, err = ioutil.ReadAll(os.Stdin); err != nil {
		log.Fatalf("can't read message from stdin: %s", err)
		panic(err)
	}
	rdr.R = bufio.NewReader(bytes.NewReader(raw))
	if hdr, err = rdr.ReadMIMEHeader(); err != nil {
		log.Printf("can't parse headers of message: %s", err)
		// We'll use an empty set of headers and go on, because we'll
		// want to record the message in the database.  We have to give
		// it a message ID, though, because the database requires those
		// to be unique.
		hdr = make(textproto.MIMEHeader)
		hdr.Set("Message-Id", time.Now().Format(time.RFC3339Nano))
	}
	em.MessageID = hdr.Get("Message-Id")
	em.Timestamp = time.Now()
	em.Type = model.EmailUnrecognized // for now
	em.Attention = true
	em.From = removeAddressRE.ReplaceAllLiteralString(hdr.Get("From"), "")
	em.Subject = hdr.Get("Subject")
	tx.CreateEmailMessage(&em, raw)
	tx.Commit()
	log.Printf("received and recorded %s", em.MessageID)
}
