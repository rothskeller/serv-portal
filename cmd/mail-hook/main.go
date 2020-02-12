// This program gets called for incoming mail to the SERV portal, with the mail
// being delivered on stdin.
//
// Mail is handled as follows:
//   - If its MessageID header is one we've seen before, ignore it.
//   - If addressed to one or more known lists, and the sender is allowed to
//     send to all of those lists, resend the message to those lists and store
//     it in the database as EmailSent.
//   - If addressed to one or more known lists, and the sender is not allowed to
//     send to all of those lists, store it in the database as EmailModerated,
//     and send a notification to the email moderators.
//   - If addressed to none of the known lists, store it in the database as
//     EmailUnrecognized.
package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/toorop/go-dkim"
	"sunnyvaleserv.org/portal/api/email"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/authz"
	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/sendmail"
)

var removeAddressRE = regexp.MustCompile(`\s*<[^>]*>`)

func main() {
	var (
		entry      *log.Entry
		tx         *store.Tx
		auth       *authz.Authorizer
		raw        []byte
		msg        *mail.Message
		em         model.EmailMessage
		sender     *model.Person
		unauth     string
		notify     bytes.Buffer
		err        error
		recipients = map[*model.Group]bool{}
	)
	// Set up execution environment.
	if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
		panic(err)
	}
	entry = log.New("", "mail-hook")
	defer entry.Log()
	store.Open("serv.db")
	tx = store.Begin(entry)
	auth = tx.Authorizer()
	// Read and parse the email message.
	em.Timestamp = time.Now()
	if raw, err = ioutil.ReadAll(os.Stdin); err != nil {
		entry.Error = fmt.Sprintf("can't read message from stdin: %s", err)
		tx.Rollback()
		entry.Log()
		return
	}
	if msg, err = mail.ReadMessage(bytes.NewReader(raw)); err != nil {
		// We'll want to record the bogus message in the database, so
		// we'll need to make up a unique message ID for it.
		em.MessageID = time.Now().Format(time.RFC3339Nano)
		em.Type = model.EmailUnrecognized
		em.Error = "ReadMessage: " + err.Error()
		em.Attention = true
		goto DONE
	}
	em.Subject = msg.Header.Get("Subject")
	if em.MessageID = msg.Header.Get("Message-Id"); em.MessageID == "" {
		// A message without an ID is bogus; we'll treat it as
		// unrecognized.  We need to make up a unique message ID to
		// store it, though.
		em.MessageID = time.Now().Format(time.RFC3339Nano)
		em.Type = model.EmailUnrecognized
		em.Error = "No Message-Id header"
		em.Attention = true
		goto DONE
	}
	if f, err := mail.ParseAddress(msg.Header.Get("From")); err == nil {
		if f.Name != "" {
			em.From = f.Name
		} else {
			em.From = f.Address
		}
		if sender = tx.FetchPersonByEmail(strings.ToLower(f.Address)); sender == nil {
			unauth = "unknown sender"
		}
	} else {
		em.Type = model.EmailUnrecognized
		em.Error = "ParseAddress(From): " + err.Error()
		em.Attention = true
		goto DONE
	}
	// If we have seen this MessageID before, ignore the message altogether.
	if tx.FetchEmailMessageByMessageID(em.MessageID) != nil {
		entry.Change("received and ignored duplicate of %s", em.MessageID)
		tx.Rollback()
		entry.Log()
		return
	}
	// Who is it addressed to?
	for _, hdr := range []string{"To", "Cc"} {
		if addrs, err := msg.Header.AddressList(hdr); err != nil && err != mail.ErrHeaderNotPresent {
			em.Type = model.EmailUnrecognized
			em.Error = "ParseAddressList(" + hdr + "): " + err.Error()
			em.Attention = true
			goto DONE
		} else {
			for _, addr := range addrs {
				addr.Address = strings.ToLower(addr.Address)
				if !strings.HasSuffix(addr.Address, "@sunnyvaleserv.org") {
					continue
				}
				if addr.Address == "admin@sunnyvaleserv.org" {
					continue
				}
				if group := auth.FetchGroupByEmail(addr.Address[:len(addr.Address)-len("@sunnyvaleserv.org")]); group == nil {
					em.Type = model.EmailUnrecognized
					em.Error = "Unknown recipient " + addr.Address
					em.Attention = true
					goto DONE
				} else {
					if !recipients[group] {
						em.Groups = append(em.Groups, group.ID)
						recipients[group] = true
						if unauth == "" && !auth.CanPAG(sender.ID, model.PrivSendEmailMessages, group.ID) {
							unauth = "sender not authorized to send to " + group.Email
						}
					}
				}
			}
		}
	}
	if len(recipients) == 0 {
		// No recognized group.
		em.Type = model.EmailUnrecognized
		em.Error = "No groups on To or Cc list"
		em.Attention = true
		goto DONE
	}
	// We checked above whether the sender has privilege to send to all of
	// the recipient groups.  If so, let's take the time to verify that the
	// sender isn't spoofed.
	if unauth == "" {
		if valid, err := dkim.Verify(&raw); err != nil || valid != dkim.SUCCESS {
			unauth = "sender not DKIM-verified"
		}
	}
	// If not authorized, mark it as moderated and set the attention flag.
	if unauth != "" {
		em.Type = model.EmailModerated
		em.Error = "Message requires moderation: " + unauth
		em.Attention = true
		goto DONE
	}
	// Resend the message to every member of the recipient groups.
	tx.CreateEmailMessage(&em, raw)
	email.SendMessage(tx, &em)
	tx.Commit()
	return
DONE:
	tx.CreateEmailMessage(&em, raw)
	tx.Commit()
	if em.Attention {
		var toLists []string
		for _, g := range em.Groups {
			toLists = append(toLists, tx.Authorizer().FetchGroup(g).Name)
		}
		fmt.Fprintf(&notify, "From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>\r\nTo: admin@sunnyvaleserv.org\r\nSubject: Email Needs Attention\r\n\r\nSunnyvaleSERV.org has received an email that needs attention:\n\nFrom: %s\nTo: %s\nSubject: %s\nType: %s\n",
			em.From, strings.Join(toLists, ", "), em.Subject, model.EmailMessageTypeNames[em.Type])
		if em.Error != "" {
			fmt.Fprintf(&notify, "Error: %s\n", em.Error)
		}
		fmt.Fprintf(&notify, "\nPlease visit https://SunnyvaleSERV.org/admin/emails to address it.\n")
		sendmail.SendMessage("admin@sunnyvaleserv.org", []string{"admin@sunnyvaleserv.org"}, notify.Bytes())
	}
}
