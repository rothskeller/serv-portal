package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/mail"
	"net/textproto"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"sunnyvaleserv.org/portal/maillist"
)

func sendListEmail(tf *os.File, messageID string, list *maillist.List, md *mailMetadata) (err error) {
	// Start by reading the message.
	var (
		fname string
		raw   []byte
		msg   *mail.Message
		body  []byte
	)
	log.Printf("  Sending to %s:", list.Name)
	fname = filepath.Join("maillist/QUEUE", messageID)
	if raw, err = os.ReadFile(fname); err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	if msg, err = mail.ReadMessage(bytes.NewReader(raw)); err != nil {
		return fmt.Errorf("parse message: %w", err)
	}
	if body, err = io.ReadAll(msg.Body); err != nil {
		return fmt.Errorf("read message body: %w", err)
	}
	// Create a rewriter for it.
	rewr := newMessageRewriter(textproto.MIMEHeader(msg.Header), body)
	var last time.Time
	for email, rdata := range list.Recipients {
		if md.sent.Has(email) {
			log.Printf("    Already sent to %s.", email)
			continue
		}
		time.Sleep(time.Until(last.Add(100 * time.Millisecond))) // TODO shouldn't be hard-coded
		last = time.Now()
		if err = resendMessageToOne(msg.Header, raw, rewr, list, email, rdata); err != nil {
			return err
		}
		fmt.Fprintf(tf, "S %s %s\n", time.Now().Format(time.RFC3339), email)
	}
	list.Reason = "as a bcc: for the " + list.Name + " list"
	list.NoUnsubscribe = true
	for email, rdata := range list.Bcc {
		if md.sent.Has(email) {
			log.Printf("    Already sent to %s.", email)
			continue
		}
		time.Sleep(time.Until(last.Add(100 * time.Millisecond))) // TODO shouldn't be hard-coded
		last = time.Now()
		if err = resendMessageToOne(msg.Header, raw, rewr, list, email, rdata); err != nil {
			return err
		}
		fmt.Fprintf(tf, "S %s %s\n", time.Now().Format(time.RFC3339), email)
	}
	tstamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(tf, "L %s %s\n", tstamp, list.Name)
	tstamp = tstamp[:len(tstamp)-6] // remove time zone
	dirname := filepath.Join("maillist", list.Name)
	if err := os.MkdirAll(dirname, 0777); err != nil {
		log.Fatalf("ERROR: creating %s: %s", dirname, err)
	}
	baseFName := filepath.Join(dirname, tstamp)
	fname = baseFName
	seq := 1
	for {
		if fh, err := os.OpenFile(fname, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666); err != nil {
			seq++
			fname = fmt.Sprintf("%s-%d", baseFName, seq)
		} else {
			fh.Close()
			break
		}
	}
	if fh, err := os.Create(fname); err != nil {
		log.Fatalf("ERROR: creating %s: %s", fname, err)
	} else if _, err = fh.Write(raw); err != nil {
		log.Fatalf("ERROR: archiving %s: %s", fname, err)
	} else if err = fh.Close(); err != nil {
		log.Fatalf("ERROR: archiving %s: %s", fname, err)
	}
	fname += ".data"
	if fh, err := os.Create(fname); err != nil {
		log.Fatalf("ERROR: creating %s: %s", fname, err)
	} else if _, err = tf.Seek(0, 0); err != nil {
		log.Fatalf("ERROR: seek in tracking file: %s", err)
	} else if _, err = io.Copy(fh, tf); err != nil {
		log.Fatalf("ERROR: archiving %s: %s", fname, err)
	} else if err = fh.Close(); err != nil {
		log.Fatalf("ERROR: archiving %s: %s", fname, err)
	}
	return nil
}

// resendMessageToOne sends the incoming message to a single recipient on the
// mailing list it was addressed to.
func resendMessageToOne(
	hdr mail.Header, raw []byte,
	rewr *messageRewriter, list *maillist.List, email string, rdata *maillist.RecipientData,
) (err error) {
	var (
		from *mail.Address
		buf  bytes.Buffer
		crlf CRLFWriter
	)
	crlf = NewCRLFWriter(&buf)
	from, _ = mail.ParseAddress(hdr.Get("From")) // previously parsed, so we know there is no error
	if from.Name == "" {
		from.Name = from.Address
	}
	fmt.Fprintf(crlf, "From: %s via %s <%s@sunnyvaleserv.org>\n", quoteIfNeeded(from.Name), list.DisplayName, list.Name)
	if !emitSelectedHeaders(&buf, raw) {
		// No Reply-To header, so generate one.
		fmt.Fprintf(crlf, "Reply-To: %s\n", hdr.Get("From"))
	}
	fmt.Fprintf(crlf, "List-Unsubscribe: <mailto:admin@SunnyvaleSERV.org?subject=%s>, <https://SunnyvaleSERV.org/unsubscribe/%s/%s>\n",
		url.QueryEscape(fmt.Sprintf("Unsubscribe %s from %s", email, list.Name)), rdata.UnsubscribeToken, list.Name)
	fmt.Fprintf(crlf, "List-Unsubscribe-Post: List-Unsubscribe=One-Click\n")
	fmt.Fprintln(crlf)
	if err = rewr.rewrite(&buf, list, email, rdata); err != nil {
		return err
	}
	cset := "serv-outgoing"
	_, err = sesClient.SendRawEmail(context.Background(), &ses.SendRawEmailInput{
		RawMessage:           &types.RawMessage{Data: buf.Bytes()},
		Destinations:         []string{email},
		ConfigurationSetName: &cset,
	})
	if err != nil {
		log.Printf("    Sending to %s:", email)
	} else {
		log.Printf("    Sent to %s.", email)
	}
	return err
}

// emitSelectedHeaders emits desired headers from the incoming email to the
// outgoing one, keeping their order and formatting from the original.  It
// returns whether a Reply-To header was seen and emitted.
func emitSelectedHeaders(out io.Writer, raw []byte) (sawReplyTo bool) {
	var lastEmitted bool
	for {
		var (
			line   []byte
			idx    int
			header string
		)
		idx = bytes.IndexByte(raw, '\n')
		if idx == 0 || (idx == 1 && raw[0] == '\r') {
			return
		}
		idx++
		line, raw = raw[:idx], raw[idx:]
		if line[0] == ' ' || line[0] == '\t' {
			if lastEmitted {
				out.Write(line)
			}
			continue
		}
		if idx = bytes.IndexByte(line, ':'); idx < 0 {
			lastEmitted = false
			continue
		}
		header = strings.ToLower(strings.TrimSpace(string(line[:idx])))
		switch header {
		case "reply-to":
			sawReplyTo = true
			fallthrough
		case "cc", "content-transfer-encoding", "content-type", "date", "in-reply-to",
			"mime-version", "organization", "subject", "to":
			out.Write(line)
			lastEmitted = true
		default:
			lastEmitted = false
		}
	}
}

var unquotedRE = regexp.MustCompile("^[-a-zA-Z0-9!#$%&'*+/=?^_`{}|~ ]+$")

// quoteIfNeeded returns the string passed to it, quoted appropriately for
// inclusion in an email header if quoting is needed.
func quoteIfNeeded(s string) string {
	if unquotedRE.MatchString(s) {
		return s
	}
	return `"` + strings.ReplaceAll(s, `"`, `\"`) + `"`
}
