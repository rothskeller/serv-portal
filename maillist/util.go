package main

import (
	"bytes"
	"crypto/rand"
	"encoding/base64"
	"io"
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

// A Split76Writer is a writer that writes data transparently, inserting a CRLF
// after every 76 characters and just before closing.  It is used to introduce
// line breaks in a base64 encoded block.
type Split76Writer struct {
	w  io.Writer
	ll int
}

// NewSplit76Writer returns a new Split76Writer wrapping the supplied writer
// (which is probably a base64 encoder).
func NewSplit76Writer(w io.Writer) *Split76Writer {
	return &Split76Writer{w: w}
}

// Write writes to the Split76Writer.
func (w *Split76Writer) Write(b []byte) (n int, err error) {
	for len(b) > 0 {
		if w.ll == 76 {
			if _, err = w.w.Write([]byte{'\r', '\n'}); err != nil {
				return n, err
			}
			w.ll = 0
		}
		c := 76 - w.ll
		if c > len(b) {
			c = len(b)
		}
		if wn, err := w.w.Write(b[:c]); err != nil {
			n += wn
			return n, err
		}
		n += c
		w.ll += c
		b = b[c:]
	}
	return n, nil
}

// Close closes the Split76Writer and, if it supports closing, the writer it
// wraps.
func (w *Split76Writer) Close() error {
	if w.ll != 0 {
		if _, err := w.w.Write([]byte{'\r', '\n'}); err != nil {
			return err
		}
		w.ll = 0
	}
	if cw, ok := w.w.(io.Closer); ok {
		return cw.Close()
	}
	return nil
}

// randomToken returns a random token string, used for various purposes.
func randomToken() string {
	var (
		tokenb [12]byte
		err    error
	)
	if _, err = rand.Read(tokenb[:]); err != nil {
		panic(err)
	}
	return base64.URLEncoding.EncodeToString(tokenb[:])
}
