package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/ses"
)

// maybeSendMessageToList sends the message to everyone on the list, but only if
// the sender is allowed to do so.  Otherwise it sends it to the moderators of
// the list for moderation.
func maybeSendMessageToList(
	ctx context.Context, client *ses.Client, rcpt *SESReceipt, msgid string, hdr mail.Header, body, raw []byte,
	listname string, list *ListData,
) (err error) {
	var (
		froma    *mail.Address
		problems []string
	)
	if len(list.Senders) == 1 && list.Senders[0] == "*" {
		// nothing
	} else if froma, err = mail.ParseAddress(hdr.Get("From")); err != nil {
		problems = append(problems, "The message From address could not be parsed.")
	} else if !slices.Contains(list.Senders, strings.ToLower(froma.Address)) {
		problems = append(problems, fmt.Sprintf("The sender %s is not authorized to send to %s.", html.EscapeString(froma.String()), listname))
	}
	if rcpt.DMARCVerdict.Status != "PASS" && rcpt.DMARCVerdict.Status != "GRAY" {
		problems = append(problems, fmt.Sprintf("The DMARC verdict is %q.", rcpt.DMARCVerdict))
	}
	if rcpt.DKIMVerdict.Status != "PASS" && rcpt.DKIMVerdict.Status != "GRAY" {
		problems = append(problems, fmt.Sprintf("The DKIM verdict is %q.", rcpt.DKIMVerdict))
	}
	if rcpt.SPFVerdict.Status != "PASS" && rcpt.SPFVerdict.Status != "GRAY" {
		problems = append(problems, fmt.Sprintf("The SPF verdict is %q.", rcpt.SPFVerdict))
	}
	if rcpt.SpamVerdict.Status != "PASS" {
		problems = append(problems, fmt.Sprintf("The spam verdict is %q.", rcpt.SpamVerdict))
	}
	if rcpt.VirusVerdict.Status != "PASS" {
		problems = append(problems, fmt.Sprintf("The virus verdict is %q.", rcpt.VirusVerdict))
	}
	if len(problems) != 0 {
		return sendForModeration(ctx, client, msgid, hdr, body, raw, listname, list, problems)
	}
	return resendMessageToList(ctx, client, hdr, body, raw, listname, list.Receivers)
}

func sendForModeration(
	ctx context.Context, client *ses.Client, msgid string, hdr mail.Header, body, raw []byte,
	listname string, list *ListData, reasons []string,
) (err error) {
	log.Printf("- sending to moderation for %s", listname)
	var from = fmt.Sprintf("%s Moderator <%s.mod@mx.sunnyvaleserv.org>", listname, listname)
	var subject = fmt.Sprintf("[MOD] %s", hdr.Get("Subject"))
	var comment = "<p>This message needs moderation because:<br>• " + strings.Join(reasons, "<br>• ") + "</p><p>To approve, reply to this message.  To reject, ignore this message.</p><p>[MOD]MSGID: " + msgid + "</p>"
	return forwardMessage(ctx, client, hdr, body, raw, from, "", list.Moderators, subject, comment)
}

func handleModerationResponse(
	ctx context.Context, s3Client *s3.Client, sesClient *ses.Client, listname string, ald AllListData, hdr mail.Header, body []byte,
) (err error) {
	var (
		ld    *ListData
		msgid string
		raw   []byte
	)
	if ld = ald[listname]; ld == nil {
		return fmt.Errorf("moderation response for unknown list %q", listname)
	}
	if body, err = extractPlainText(textproto.MIMEHeader(hdr), body); err != nil || body == nil {
		return errors.New("no plain text content in moderation response")
	} else if idx := bytes.Index(body, []byte("[MOD]MSGID: ")); idx < 0 {
		return errors.New("[MOD]MSGID not found in moderation response")
	} else {
		if idx2 := bytes.IndexFunc(body[idx+12:], func(r rune) bool {
			return (r < 'a' || r > 'z') && (r < '0' || r > '9')
		}); idx2 < 0 {
			return errors.New("[MOD]MSGID end not found in moderation response")
		} else {
			msgid = string(body[idx+12 : idx+12+idx2])
		}
	}
	log.Printf("- moderation approval for %s on %s", msgid, listname)
	if hdr, body, raw, err = readMail(ctx, s3Client, msgid); err != nil {
		return fmt.Errorf("reading moderated message %s: %s", msgid, err)
	}
	return resendMessageToList(ctx, sesClient, hdr, body, raw, listname, ld.Receivers)
}

// extractPlainText extracts the plain text portion of a message from its body.
// It returns a nil body if there is none.  This is a recursive function, to
// handled nested multipart bodies.
func extractPlainText(header textproto.MIMEHeader, body []byte) (nbody []byte, err error) {
	var (
		mediatype string
		params    map[string]string
	)
	// Decode any content transfer encoding.  If we come across an encoding
	// we can't handle, or we have an error decoding, return an empty body
	// with a notplain indicator.
	switch strings.ToLower(header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		break // no decoding needed
	case "quoted-printable":
		if body, err = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(body))); err != nil {
			return nil, err
		}
	case "base64":
		if body, err = io.ReadAll(base64.NewDecoder(base64.StdEncoding, bytes.NewReader(body))); err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}
	// Decode the content type.
	if ct := header.Get("Content-Type"); ct != "" {
		if mediatype, params, err = mime.ParseMediaType(header.Get("Content-Type")); err != nil {
			return nil, err // can't decode Content-Type
		}
	} else {
		mediatype, params = "text/plain", map[string]string{}
	}
	// If the content type is multipart, look for the last plain text part
	// in it.  This is a recursive call.
	if strings.HasPrefix(mediatype, "multipart/") {
		var (
			mr       *multipart.Reader
			part     *multipart.Part
			partbody []byte
			found    []byte
		)
		mr = multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for {
			part, err = mr.NextRawPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				return nil, err // Can't decode multipart body
			}
			partbody, _ = io.ReadAll(part)
			plain, err := extractPlainText(part.Header, partbody)
			if err != nil {
				return nil, err
			}
			if plain != nil {
				found = plain
			}
		}
		return found, nil
	}
	// If the content type is anything other than text/plain, we're out of
	// luck.
	if mediatype != "text/plain" {
		return nil, nil
	}
	// In theory we also ought to check the charset, but we'll elide that
	// until experience proves a need.
	return body, nil
}
