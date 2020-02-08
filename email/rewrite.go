package email

import (
	"bytes"
	"fmt"
	"html"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

type messagePart struct {
	header          textproto.MIMEHeader
	prefix          []byte
	suffix          []byte
	boundary        string
	alternative     bool
	parts           []*messagePart
	lastPart        *messagePart
	quotedPrintable bool
	plain           bool
	html            bool
}

func makeMessagePart(header textproto.MIMEHeader, body []byte) (mp *messagePart, rewritable bool) {
	var (
		mediaType string
		params    map[string]string
		err       error
	)
	if mediaType, params, err = mime.ParseMediaType(header.Get("Content-Type")); err != nil {
		return &messagePart{header: header, prefix: body}, false
	}
	switch mediaType {
	case "text/plain":
		return makePlainMessagePart(header, body)
	case "text/html":
		return makeHTMLMessagePart(header, body)
	case "multipart/alternative":
		return makeMultiMessagePart(header, body, params["boundary"], true)
	case "multipart/mixed":
		return makeMultiMessagePart(header, body, params["boundary"], false)
	default:
		return &messagePart{header: header, prefix: body}, false
	}
}

func makePlainMessagePart(header textproto.MIMEHeader, body []byte) (mp *messagePart, rewritable bool) {
	var (
		qpr *quotedprintable.Reader
		err error
	)
	mp = &messagePart{header: header}
	switch strings.ToLower(header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		mp.prefix = body
		mp.plain = true
	case "quoted-printable":
		qpr = quotedprintable.NewReader(bytes.NewReader(body))
		mp.prefix, err = ioutil.ReadAll(qpr)
		if err != nil {
			mp.prefix = body
			break
		}
		mp.quotedPrintable = true
		mp.plain = true
	default:
		mp.prefix = body
	}
	return mp, mp.plain
}

func makeHTMLMessagePart(header textproto.MIMEHeader, body []byte) (mp *messagePart, rewritable bool) {
	var (
		qpr     *quotedprintable.Reader
		decoded []byte
		idx     int
		err     error
	)
	mp = &messagePart{header: header}
	switch strings.ToLower(header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		decoded = body
	case "quoted-printable":
		qpr = quotedprintable.NewReader(bytes.NewReader(body))
		decoded, err = ioutil.ReadAll(qpr)
		if err != nil {
			mp.prefix = body
			return mp, false
		}
		mp.quotedPrintable = true
	default:
		mp.prefix = body
		return mp, false
	}
	mp.html = true
	if idx = bytes.LastIndex(decoded, []byte("</body>")); idx >= 0 {
		mp.prefix = decoded[:idx]
		mp.suffix = decoded[idx:]
	} else {
		mp.prefix = decoded
	}
	return mp, true
}

func makeMultiMessagePart(header textproto.MIMEHeader, body []byte, boundary string, alternative bool) (mp *messagePart, rewritable bool) {
	var (
		mpr   *multipart.Reader
		part  *multipart.Part
		cmp   *messagePart
		rewr  bool
		pbody []byte
		err   error
	)
	mp = &messagePart{header: header, boundary: boundary, alternative: alternative}
	mpr = multipart.NewReader(bytes.NewReader(body), boundary)
	for part, err = mpr.NextPart(); err == nil; part, err = mpr.NextPart() {
		if pbody, err = ioutil.ReadAll(part); err != nil {
			mp.prefix = body
			return mp, false
		}
		cmp, rewr = makeMessagePart(part.Header, pbody)
		if rewr {
			rewritable = true
		}
		if !alternative && rewr {
			mp.lastPart = cmp
		}
		mp.parts = append(mp.parts, cmp)
	}
	if err != io.EOF {
		mp.prefix = body
		mp.lastPart = nil
		mp.parts = nil
		return mp, false
	}
	return mp, rewritable
}

func rewrite(w io.Writer, mp *messagePart, groupAddress, personName, personAddress string) (err error) {
	if len(mp.parts) != 0 {
		var (
			mpw *multipart.Writer
			cp  io.Writer
		)
		mpw = multipart.NewWriter(w)
		mpw.SetBoundary(mp.boundary)
		for _, c := range mp.parts {
			if cp, err = mpw.CreatePart(c.header); err != nil {
				return err
			}
			if mp.alternative || mp.lastPart == c {
				err = rewrite(cp, c, groupAddress, personName, personAddress)
			} else {
				err = copyPart(cp, c)
			}
			if err != nil {
				return err
			}
		}
		return mpw.Close()
	}
	if mp.quotedPrintable {
		w = quotedprintable.NewWriter(w)
	}
	if _, err = w.Write(mp.prefix); err != nil {
		return err
	}
	if mp.plain {
		if _, err = fmt.Fprintf(w, "________\r\nThis message was sent to %s <%s> via the %s@SunnyvaleSERV.org mailing list.\r\nTo unsubscribe, visit https://SunnyvaleSERV.org/unsubscribe/%s.\r\n", personName, personAddress, groupAddress, personAddress); err != nil {
			return err
		}
	}
	if mp.html {
		if _, err = fmt.Fprintf(w, "<div style=\"height:1em;width:5em;border-bottom:1px solid #888\"></div><div style=\"color:#888\">This message was sent to %s &lt;<a style=\"color:#888\">%s</a>&gt; via the <a style=\"color:#888\">%s@SunnyvaleSERV.org</a> mailing list.<br>To unsubscribe, visit <a style=\"color:#888\" href=\"https://SunnyvaleSERV.org/unsubscribe/%[2]s\">https://SunnyvaleSERV.org/unsubscribe/%[2]s</a>.</div>", html.EscapeString(personName), html.EscapeString(personAddress), groupAddress); err != nil {
			return err
		}
	}
	if _, err = w.Write(mp.suffix); err != nil {
		return err
	}
	if mp.quotedPrintable {
		if err = w.(*quotedprintable.Writer).Close(); err != nil {
			return err
		}
	}
	return nil
}

func copyPart(w io.Writer, mp *messagePart) (err error) {
	if len(mp.parts) != 0 {
		var (
			mpw *multipart.Writer
			cp  io.Writer
		)
		mpw = multipart.NewWriter(w)
		mpw.SetBoundary(mp.boundary)
		for _, c := range mp.parts {
			if cp, err = mpw.CreatePart(c.header); err != nil {
				return err
			}
			if err = copyPart(cp, c); err != nil {
				return err
			}
		}
		return mpw.Close()
	}
	if mp.quotedPrintable {
		w = quotedprintable.NewWriter(w)
	}
	if _, err = w.Write(mp.prefix); err != nil {
		return err
	}
	if _, err = w.Write(mp.suffix); err != nil {
		return err
	}
	if mp.quotedPrintable {
		if err = w.(*quotedprintable.Writer).Close(); err != nil {
			return err
		}
	}
	return nil
}
