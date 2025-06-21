package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"log"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/util/config"
)

func messageNeedsModeration(messageID string, list *maillist.List, md *mailMetadata) (problems []string) {
	if len(md.verdicts) != 5 {
		problems = append(problems, "Missing DKIM/DMARC/SPF/Spam/Virus verdicts")
	} else {
		if md.verdicts[0] != "PASS" && md.verdicts[0] != "GRAY" {
			problems = append(problems, "DKIM verdict is "+md.verdicts[0])
		}
		if md.verdicts[1] != "PASS" && md.verdicts[1] != "GRAY" {
			problems = append(problems, "DMARC verdict is "+md.verdicts[1])
		}
		if md.verdicts[2] != "PASS" && md.verdicts[2] != "GRAY" {
			problems = append(problems, "SPF verdict is "+md.verdicts[2])
		}
		if md.verdicts[3] != "PASS" {
			problems = append(problems, "Spam verdict is "+md.verdicts[3])
		}
		if md.verdicts[4] != "PASS" {
			problems = append(problems, "Virus verdict is "+md.verdicts[4])
		}
	}
	// As a special case, "admin" is not moderated for sender.
	if list.Name == "admin" {
		return problems
	}
	// To check the sender, we need to read the message and get the From
	// line.
	var (
		fname string
		fh    *os.File
		msg   *mail.Message
		from  string
		err   error
	)
	fname = filepath.Join("maillist/QUEUE", messageID)
	if fh, err = os.Open(fname); err != nil {
		log.Fatalf("ERROR: open message: %s", err)
	}
	defer fh.Close()
	if msg, err = mail.ReadMessage(fh); err != nil {
		log.Fatalf("ERROR: read message %s: %s", messageID, err)
	}
	from = msg.Header.Get("From")
	if addr, err := mail.ParseAddress(from); err != nil {
		problems = append(problems, "Invalid From: header")
	} else if !list.Senders.Has(strings.ToLower(addr.Address)) {
		problems = append(problems, fmt.Sprintf("%s is not an authorized sender to %s", addr.Address, list.Name))
	}
	return problems
}

func requestModeration(tf *os.File, messageID string, list *maillist.List, problems []string) (err error) {
	var (
		fname   string
		raw     []byte
		msg     *mail.Message
		body    []byte
		from    string
		subject string
		comment string
	)
	log.Printf("  Requesting moderation for %s", list.Name)
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
	from = fmt.Sprintf("%s Moderator <%s.mod@mx.sunnyvaleserv.org>", list.Name, list.Name)
	subject = fmt.Sprintf("[MOD] %s", msg.Header.Get("Subject"))
	for i := range problems {
		problems[i] = html.EscapeString(problems[i])
	}
	comment = "<p>This message needs moderation because:<br>• " + strings.Join(problems, "<br>• ") + "</p><p>To approve, reply to this message.  To reject, ignore this message.</p><p>[MOD]MSGID: " + messageID + "</p>"
	if err = forwardMessage(msg.Header, body, raw, from, "", strings.Split(config.Get("listModerators"), ","), subject, comment); err != nil {
		return err
	}
	tstamp := time.Now().Format(time.RFC3339)
	fmt.Fprintf(tf, "M %s %s\n", tstamp, list.Name)
	return nil
}

func handleModerationResponse(messageID, listname string) (err error) {
	var (
		fname    string
		raw      []byte
		msg      *mail.Message
		body     []byte
		modMsgID string
		modTF    *os.File
		modMeta  *mailMetadata
	)
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
	if body, err = extractPlainText(textproto.MIMEHeader(msg.Header), body); err != nil || body == nil {
		log.Printf("  WARNING: ignoring moderation response with no plain text content")
		return nil
	} else if idx := bytes.Index(body, []byte("[MOD]MSGID: ")); idx < 0 {
		log.Printf("  WARNING: ignoring moderation response with no message ID")
		return nil
	} else {
		if idx2 := bytes.IndexFunc(body[idx+12:], func(r rune) bool {
			return r == '\r' || r == '\n'
		}); idx2 < 0 {
			log.Printf("  WARNING: ignoring moderation response with no newline after message ID")
			return nil
		} else {
			modMsgID = string(body[idx+12 : idx+12+idx2])
		}
	}
	fname = filepath.Join("maillist/QUEUE", modMsgID+".data")
	if modTF, err = os.OpenFile(fname, os.O_RDWR, 0666); os.IsNotExist(err) {
		log.Printf("  WARNING: ignoring moderation response for nonexistent message %s", modMsgID)
		return nil
	} else if err != nil {
		return fmt.Errorf("read message: %w", err)
	}
	defer modTF.Close()
	if modMeta, err = readTracking(modTF); err != nil {
		return fmt.Errorf("QUEUE/%s.data: %w", modMsgID, err)
	}
	if !modMeta.moderating.Has(listname) {
		log.Printf("  WARNING: ignoring moderation response for untargeted list %s", listname)
		return nil
	}
	if modMeta.approved.Has(listname) {
		log.Printf("  WARNING: ignoring duplicate approval for list %s", listname)
		return nil
	}
	from := msg.Header.Get("From")
	if addr, err := mail.ParseAddress(from); err == nil {
		from = addr.Address
	}
	if !slices.Contains(strings.Split(config.Get("listModerators"), ","), from) {
		log.Printf("  WARNING: ignoring moderation response from non-moderator %s", from)
		return nil
	}
	if fields := strings.Fields(from); len(fields) > 1 {
		from = fields[0]
	} else if len(fields) == 0 {
		from = "-"
	}
	fmt.Fprintf(modTF, "A %s %s %s\n", time.Now().Format(time.RFC3339), listname, from)
	if err = modTF.Chmod(0644); err != nil {
		return fmt.Errorf("chmod QUEUE/%s.data: %w", modMsgID, err)
	}
	toHandle.Insert(modMsgID)
	os.Remove("maillist/QUEUE/" + messageID)
	os.Remove("maillist/QUEUE/" + messageID + ".data")
	log.Printf("  Moderation approval for %s.", listname)
	return nil
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
