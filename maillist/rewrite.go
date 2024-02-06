package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net/textproto"
	"strings"
)

type messageRewriter struct {
	header     textproto.MIMEHeader // header of the message or part
	rewritable bool                 // true if this part is rewritable
	// for multipart parts:
	boundary    string             // boundary string
	alternative bool               // whether the parts are alternatives
	parts       []*messageRewriter // list of parts
	primaryPart *messageRewriter   // which part to modify for footer
	// for leaf parts:
	prefix          []byte // decoded content before footer insertion
	suffix          []byte // decoded content after footer insertion
	quotedPrintable bool   // true if CTE is quoted-printable
	base64          bool   // true if CTE is base64
	plain           bool   // true if content is plain text
	html            bool   // true if content is HTML
}

// newMessageRewriter creates a new messageRewriter structure reflecting the
// provided message or message part.
func newMessageRewriter(header textproto.MIMEHeader, body []byte) (mr *messageRewriter) {
	var (
		mediaType string
		params    map[string]string
		err       error
	)
	mr = &messageRewriter{header: header, prefix: body}
	if ct := header.Get("Content-Type"); ct == "" {
		mediaType = "text/plain"
	} else if mediaType, params, err = mime.ParseMediaType(ct); err != nil {
		mediaType = ""
	}
	switch {
	case mediaType == "text/plain":
		mr.parsePlainBody()
	case mediaType == "text/html":
		mr.parseHTMLBody()
	case strings.HasPrefix(mediaType, "multipart/"):
		mr.parseMultipartBody(mediaType[10:], params["boundary"])
	}
	return mr
}

// parsePlainBody parses the body of a text/plain message or message part.
func (mr *messageRewriter) parsePlainBody() {
	var (
		qpr     *quotedprintable.Reader
		br      io.Reader
		decoded []byte
		err     error
	)
	switch strings.ToLower(mr.header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		decoded = mr.prefix
	case "quoted-printable":
		qpr = quotedprintable.NewReader(bytes.NewReader(mr.prefix))
		decoded, err = io.ReadAll(qpr)
		if err != nil {
			return
		}
		mr.quotedPrintable = true
	case "base64":
		br = base64.NewDecoder(base64.StdEncoding, bytes.NewReader(mr.prefix))
		decoded, err = io.ReadAll(br)
		if err != nil {
			return
		}
		mr.base64 = true
	default:
		return
	}
	mr.plain, mr.rewritable = true, true
	mr.prefix = decoded
}

// parseHTMLBody parses the body of a text/html message or message part.
func (mr *messageRewriter) parseHTMLBody() {
	var (
		qpr     *quotedprintable.Reader
		br      io.Reader
		decoded []byte
		idx     int
		err     error
	)
	switch strings.ToLower(mr.header.Get("Content-Transfer-Encoding")) {
	case "", "7bit", "8bit", "binary":
		decoded = mr.prefix
	case "quoted-printable":
		qpr = quotedprintable.NewReader(bytes.NewReader(mr.prefix))
		decoded, err = io.ReadAll(qpr)
		if err != nil {
			return
		}
		mr.quotedPrintable = true
	case "base64":
		br = base64.NewDecoder(base64.StdEncoding, bytes.NewReader(mr.prefix))
		decoded, err = io.ReadAll(br)
		if err != nil {
			return
		}
		mr.base64 = true
	default:
		return
	}
	mr.html, mr.rewritable = true, true
	if idx = bytes.LastIndex(decoded, []byte("</body>")); idx >= 0 {
		mr.prefix = decoded[:idx]
		mr.suffix = decoded[idx:]
	} else {
		mr.prefix = decoded
	}
}

// parseMultipartBody parses the body of a multipart/* message or message part.
func (mr *messageRewriter) parseMultipartBody(subtype, boundary string) {
	var (
		mpr     *multipart.Reader
		part    *multipart.Part
		childmr *messageRewriter
		pbody   []byte
		err     error
	)
	mr.boundary = boundary
	mr.alternative = subtype == "alternative"
	mpr = multipart.NewReader(bytes.NewReader(mr.prefix), boundary)
	for part, err = mpr.NextPart(); err == nil; part, err = mpr.NextPart() {
		if pbody, err = io.ReadAll(part); err != nil {
			return
		}
		childmr = newMessageRewriter(part.Header, pbody)
		if childmr.rewritable {
			mr.rewritable = true
		}
		switch subtype {
		case "mixed":
			if childmr.rewritable {
				mr.primaryPart = childmr
			}
		case "related":
			if childmr.rewritable && mr.primaryPart == nil {
				mr.primaryPart = childmr
			}
		}
		mr.parts = append(mr.parts, childmr)
	}
	if err != io.EOF {
		mr.primaryPart, mr.parts, mr.rewritable = nil, nil, false
	}
}

// rewrite rewrites a message part (or a whole message) for distribution to an
// individual list recipient.  It writes the result to the provided Writer.
func (mr *messageRewriter) rewrite(w io.Writer, list string, recip *Receiver) (err error) {
	if len(mr.parts) != 0 {
		var (
			mpw    *multipart.Writer
			childp io.Writer
		)
		mpw = multipart.NewWriter(w)
		mpw.SetBoundary(mr.boundary)
		for _, childmr := range mr.parts {
			if childp, err = mpw.CreatePart(childmr.header); err != nil {
				return err
			}
			if mr.alternative || mr.primaryPart == childmr {
				err = childmr.rewrite(childp, list, recip)
			} else {
				err = childmr.copy(childp)
			}
			if err != nil {
				return err
			}
		}
		return mpw.Close()
	}
	if mr.quotedPrintable {
		w = quotedprintable.NewWriter(w)
	} else if mr.base64 {
		sw := NewSplit76Writer(w)
		defer sw.Close()
		w = base64.NewEncoder(base64.StdEncoding, sw)
	} else {
		w = NewCRLFWriter(w)
	}
	if _, err = w.Write(mr.prefix); err != nil {
		return err
	}
	if mr.plain {
		fmt.Fprintf(w, "\n________\nThis message was sent to %s <%s> via the %s@SunnyvaleSERV.org mailing list.\nTo unsubscribe, visit https://SunnyvaleSERV.org/unsubscribe/%s.\n",
			recip.Name, recip.Addr, list, recip.Token)
	}
	if mr.html {
		fmt.Fprintf(w, `<div style="height:1em;width:5em;border-bottom:1px solid #888"></div><div style="color:#888">This message was sent to %s &lt;<a style="color:#888">%s</a>&gt; via the <a style="color:#888">%s@SunnyvaleSERV.org</a> mailing list.<br>To unsubscribe, visit our <a style="color:#888" href="https://SunnyvaleSERV.org/unsubscribe/%s">unsubscribe page</a>.</div>`,
			html.EscapeString(recip.Name), html.EscapeString(recip.Addr), list, recip.Token)
	}
	if _, err = w.Write(mr.suffix); err != nil {
		return err
	}
	if mr.quotedPrintable {
		if err = w.(*quotedprintable.Writer).Close(); err != nil {
			return err
		}
	}
	if mr.base64 {
		if err = w.(io.Closer).Close(); err != nil {
			return err
		}
	}
	return nil
}

// copy copies a message or message part to the provided Writer, without change.
func (mr *messageRewriter) copy(w io.Writer) (err error) {
	if len(mr.parts) != 0 {
		var (
			mpw *multipart.Writer
			cp  io.Writer
		)
		mpw = multipart.NewWriter(w)
		mpw.SetBoundary(mr.boundary)
		for _, childmr := range mr.parts {
			if cp, err = mpw.CreatePart(childmr.header); err != nil {
				return err
			}
			if err = childmr.copy(cp); err != nil {
				return err
			}
		}
		return mpw.Close()
	}
	if mr.quotedPrintable {
		w = quotedprintable.NewWriter(w)
	} else if mr.base64 {
		sw := NewSplit76Writer(w)
		defer sw.Close()
		w = base64.NewEncoder(base64.StdEncoding, sw)
	} else {
		w = NewCRLFWriter(w)
	}
	if _, err = w.Write(mr.prefix); err != nil {
		return err
	}
	if _, err = w.Write(mr.suffix); err != nil {
		return err
	}
	if mr.quotedPrintable {
		if err = w.(*quotedprintable.Writer).Close(); err != nil {
			return err
		}
	}
	if mr.base64 {
		if err = w.(io.Closer).Close(); err != nil {
			return err
		}
	}
	return nil
}
