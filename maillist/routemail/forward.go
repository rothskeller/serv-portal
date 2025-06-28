package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/mail"
	"net/textproto"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
)

// forwardMessage takes an incoming message and forwards it as an attachment to
// the specified addresses.  The body of the forward contains the specified
// comment (HTML format) and a summary of the structure of the incoming message.
func forwardMessage(
	hdr mail.Header, body, raw []byte, from, replyTo string, to []string, subject, comment string,
) (err error) {
	var (
		buf      bytes.Buffer
		crlf     CRLFWriter
		qp       *quotedprintable.Writer
		boundary string
	)
	boundary = randomToken()
	crlf = NewCRLFWriter(&buf)
	fmt.Fprintf(crlf, "From: %s\n", from)
	if replyTo != "" {
		fmt.Fprintf(crlf, "Reply-To: %s\n", replyTo)
	}
	fmt.Fprintf(crlf, `To: %s
Subject: %s
MIME-Version: 1.0
Content-Type: multipart/mixed; boundary="%s"
Date: %s

--%s
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`, strings.Join(to, ", "), subject, boundary, time.Now().Format(time.RFC1123Z), boundary)
	qp = quotedprintable.NewWriter(&buf)
	fmt.Fprintf(qp, `<div style="margin-bottom:1em">%s</div>`, comment)
	addMessageHeaders(qp, hdr)
	addMessageStructure(qp, hdr, body)
	qp.Close()
	fmt.Fprintf(crlf, `

--%s
Content-Type: message/rfc822

`, boundary)
	buf.Write(raw)
	fmt.Fprintf(crlf, `

--%s--
`, boundary)
	cset := "serv-outgoing"
	_, err = sesClient.SendRawEmail(context.Background(), &ses.SendRawEmailInput{RawMessage: &types.RawMessage{Data: buf.Bytes()}, ConfigurationSetName: &cset})
	return err
}

func addMessageHeaders(qp *quotedprintable.Writer, hdrs mail.Header) {
	fmt.Fprint(qp, `<table style="margin-bottom:1em">`)
	for _, hdr := range []string{"From", "To", "Cc", "Subject", "Date"} {
		fmt.Fprintf(qp, `<tr><td style="padding-right:2em">%s:</td><td>%s</td></tr>`,
			hdr, html.EscapeString(hdrs.Get(hdr)))
	}
	fmt.Fprint(qp, `</table>`)
}

func addMessageStructure(qp *quotedprintable.Writer, hdr mail.Header, body []byte) {
	header := textproto.MIMEHeader(hdr)
	addMessageStructureInner(qp, header, body)
}

func addMessageStructureInner(qp *quotedprintable.Writer, header textproto.MIMEHeader, body []byte) {
	var (
		ct        string
		cd        string
		cte       string
		mediaType string
		params    map[string]string
		err       error
	)
	ct = header.Get("Content-Type")
	cd = header.Get("Content-Disposition")
	cte = header.Get("Content-Transfer-Encoding")
	if ct == "" {
		ct = "text/plain"
	}
	if mediaType, params, err = mime.ParseMediaType(ct); err != nil {
		mediaType = "INVALID"
	}
	if cd != "" {
		fmt.Fprintf(qp, `<div>%s %s</div>`, html.EscapeString(mediaType), html.EscapeString(cd))
	} else {
		fmt.Fprintf(qp, `<div>%s</div>`, html.EscapeString(mediaType))
	}
	if mediaType == "text/plain" && (cte == "" || cte == "quoted-printable" || cte == "base64") {
		var show []byte
		switch cte {
		case "quoted-printable":
			show, _ = io.ReadAll(quotedprintable.NewReader(bytes.NewReader(body)))
		case "base64":
			show = make([]byte, base64.StdEncoding.DecodedLen(len(body)))
			base64.StdEncoding.Decode(show, body)
		default:
			show = body
		}
		if len(show) > 500 {
			show = show[:500]
		}
		fmt.Fprintf(qp, `<pre style="margin-left:2em">%s</pre>`, html.EscapeString(string(show)))
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		fmt.Fprint(qp, `<div style="margin-left:2em">`)
		mpr := multipart.NewReader(bytes.NewReader(body), params["boundary"])
		for part, err := mpr.NextPart(); err == nil; part, err = mpr.NextPart() {
			pbody, _ := io.ReadAll(part)
			addMessageStructureInner(qp, part.Header, pbody)
		}
		fmt.Fprint(qp, `</div>`)
	}
}
