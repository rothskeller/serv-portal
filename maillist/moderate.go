package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"html"
	"log"
	"net/mail"
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
	ctx context.Context, s3Client *s3.Client, sesClient *ses.Client, listname string, ald AllListData, body []byte,
) (err error) {
	var (
		ld    *ListData
		msgid string
		hdr   mail.Header
		raw   []byte
	)
	if ld = ald[listname]; ld == nil {
		return fmt.Errorf("moderation response for unknown list %q", listname)
	}
	if idx := bytes.Index(body, []byte("[MOD]MSGID: ")); idx < 0 {
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
