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
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// resendMessageToList sends the incoming message to the specified list.
func resendMessageToList(
	ctx context.Context, client *ses.Client, hdr mail.Header, body, raw []byte, list string, recipients []Receiver,
) (err error) {
	log.Printf("- sending to %s", list)
	var rewr = newMessageRewriter(textproto.MIMEHeader(hdr), body)
	var last time.Time
	for i := range recipients {
		time.Sleep(time.Until(last.Add(100 * time.Millisecond)))
		last = time.Now()
		if err = resendMessageToOne(ctx, client, hdr, raw, rewr, list, &recipients[i]); err != nil {
			return err
		}
	}
	return nil
}

// resendMessageToOne sends the incoming message to a single recipient on the
// mailing list it was addressed to.
func resendMessageToOne(
	ctx context.Context, client *ses.Client, hdr mail.Header, raw []byte,
	rewr *messageRewriter, list string, recip *Receiver,
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
	fmt.Fprintf(crlf, "From: %s via %s <%s@sunnyvaleserv.org>\n", quoteIfNeeded(from.Name), list, list)
	if !emitSelectedHeaders(&buf, raw) {
		// No Reply-To header, so generate one.
		fmt.Fprintf(crlf, "Reply-To: %s\n", hdr.Get("From"))
	}
	fmt.Fprintf(crlf, "List-Unsubscribe: <mailto:admin@SunnyvaleSERV.org?subject=%s>, <https://SunnyvaleSERV.org/unsubscribe/%s/%s>\n",
		url.QueryEscape(fmt.Sprintf("Unsubscribe %s from %s", recip.Addr, list)), recip.Token, list)
	fmt.Fprintf(crlf, "List-Unsubscribe-Post: List-Unsubscribe=One-Click\n")
	fmt.Fprintln(crlf)
	if err = rewr.rewrite(&buf, list, recip); err != nil {
		return err
	}
	_, err = client.SendRawEmail(ctx, &ses.SendRawEmailInput{
		RawMessage:   &types.RawMessage{Data: buf.Bytes()},
		Destinations: []string{recip.Addr},
	})
	return err
}

// emitSelectedHeaders emits desired headers from the incoming email to the
// outgoing one, keeping their order and formatting from the original.  It
// returns whether a Reply-To header was seen and emitted.
func emitSelectedHeaders(out io.Writer, raw []byte) (sawReplyTo bool) {
	var (
		lastEmitted bool
	)
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

var unquotedRE = regexp.MustCompile("^[-a-zA-Z0-9!#$%&'*+/=?^_`{}|~.]+$")

// quoteIfNeeded returns the string passed to it, quoted appropriately for
// inclusion in an email header if quoting is needed.
func quoteIfNeeded(s string) string {
	if unquotedRE.MatchString(s) {
		return s
	}
	return `"` + strings.Replace(s, `"`, `\"`, -1) + `"`
}
