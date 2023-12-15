package request

import (
	"net/http"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
)

// A Request represents an in-progress web request.  It contains the request
// data, the response data, the caller's session data if any, and similar
// tracking data.
type Request struct {
	*http.Request
	*store.Store
	ResponseWriter
	SessionToken string
	CSRF         string
	Path         string
	Language     string
	LogEntry     *log.Entry
}

// Header and Write on Request resolve the ambiguities between http.Request
// and http.ResponseWriter, in favor of the latter.  This allows Request to be
// used in the context of an http.ResponseWriter, such as in calls to
// http.Error.
func (r *Request) Header() http.Header           { return r.ResponseWriter.Header() }
func (r *Request) Write(buf []byte) (int, error) { return r.ResponseWriter.Write(buf) }

// DisableCompression disables compression of the response.
func (r *Request) DisableCompression() {
	if r.StatusCode() != 0 {
		panic("DisableCompression called after response header written")
	}
	switch rw := r.ResponseWriter.(type) {
	case *UncompressedResponse:
		return
	case *CompressedResponse:
		r.ResponseWriter = rw.ResponseWriter
	default:
		panic("DisableCompression called on unknown response writer type")
	}
}
