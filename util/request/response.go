package request

import (
	"compress/gzip"
	"net/http"
)

// ResponseWriter is the interface that all response writers used in this
// application need to satisfy.  It adds the StatusCode method to the standard
// library's ResponseWriter interface.
type ResponseWriter interface {
	http.ResponseWriter

	// StatusCode returns the status code that has been sent to the browser,
	// or zero if headers and status code have not yet been sent.
	StatusCode() int
}

// UncompressedResponse is a minimal response writer.
type UncompressedResponse struct {
	http.ResponseWriter
	statusCode int
}

var _ ResponseWriter = (*UncompressedResponse)(nil) // verify interface compliance

// NewUncompressedResponse returns a new minimal response writer.
func NewUncompressedResponse(w http.ResponseWriter) (ur *UncompressedResponse) {
	return &UncompressedResponse{ResponseWriter: w}
}

// Write implements the http.ResponseWriter interface.
func (ur *UncompressedResponse) Write(buf []byte) (int, error) {
	if ur.statusCode == 0 {
		ur.statusCode = http.StatusOK
	}
	return ur.ResponseWriter.Write(buf)
}

// WriteHeader implements the http.ResponseWriter interface.
func (ur *UncompressedResponse) WriteHeader(statusCode int) {
	ur.statusCode = statusCode
	ur.ResponseWriter.WriteHeader(statusCode)
}

// StatusCode implements the server.ResponseWriter interface.
func (ur *UncompressedResponse) StatusCode() int { return ur.statusCode }

// CompressedResponse is a response writer that sends the response in compressed
// format.
type CompressedResponse struct {
	ResponseWriter
	gz *gzip.Writer
}

// NewCompressedResponse returns a new compressed response writer.
func NewCompressedResponse(rw ResponseWriter) (cr *CompressedResponse) {
	return &CompressedResponse{ResponseWriter: rw}
}

// Write implements the http.ResponseWriter interface.
func (cr *CompressedResponse) Write(buf []byte) (int, error) {
	if cr.gz == nil {
		cr.gz = gzip.NewWriter(cr.ResponseWriter)
		cr.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	return cr.gz.Write(buf)
}

// WriteHeader implements the http.ResponseWriter interface.
func (cr *CompressedResponse) WriteHeader(statusCode int) {
	if cr.gz == nil {
		cr.gz = gzip.NewWriter(cr.ResponseWriter)
		cr.ResponseWriter.Header().Set("Content-Encoding", "gzip")
	}
	cr.ResponseWriter.WriteHeader(statusCode)
}

// Flush implements the http.Flusher interface.
func (cr *CompressedResponse) Flush() {
	if cr.gz != nil && cr.StatusCode() != http.StatusNoContent && cr.StatusCode() != http.StatusNotModified {
		if err := cr.gz.Close(); err != nil {
			panic(err)
		}
	}
	if rw, ok := cr.ResponseWriter.(http.Flusher); ok {
		rw.Flush()
	}
}

// HTMLNoCache sends the appropriate Content-Type and Cache-Control headers to
// set an HTML content type and disable caching.
func (r *Request) HTMLNoCache() {
	r.Header().Set("Content-Type", "text/html; charset=utf-8")
	r.Header().Set("Cache-Control", "no-store")
}
