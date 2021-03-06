package sendmail

import (
	"bytes"
	"io"
	"regexp"
	"strings"
)

// A CRLFWriter is a writer that converts \n to \r\n.
type CRLFWriter struct {
	w io.Writer
}

// NewCRLFWriter returns a new CRLFWriter wrapping the supplied writer.
func NewCRLFWriter(w io.Writer) CRLFWriter {
	return CRLFWriter{w}
}

// Write writes to the CRLFWriter.
func (w CRLFWriter) Write(b []byte) (n int, err error) {
	var idx int
	var wn int
	for idx = bytes.IndexByte(b, '\n'); idx >= 0; idx = bytes.IndexByte(b, '\n') {
		wn, err = w.w.Write(b[:idx])
		n += wn
		if err != nil {
			return n, err
		}
		if _, err = w.w.Write([]byte{'\r', '\n'}); err != nil {
			return n, err
		}
		n++
		b = b[idx+1:]
	}
	if len(b) > 0 {
		wn, err = w.w.Write(b)
		n += wn
	}
	return
}

var unquotedRE = regexp.MustCompile("^[-a-zA-Z0-9!#$%&'*+/=?^_`{}|~.]+$")

// QuoteIfNeeded returns the string passed to it, quoted appropriately for
// inclusion in an email header if quoting is needed.
func QuoteIfNeeded(s string) string {
	if unquotedRE.MatchString(s) {
		return s
	}
	return `"` + strings.Replace(s, `"`, `\"`, -1) + `"`
}
