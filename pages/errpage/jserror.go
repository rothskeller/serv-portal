package errpage

import (
	"io"
	"net/http"

	"sunnyvaleserv.org/portal/util/request"
)

func PostJSError(r *request.Request) {
	details, _ := io.ReadAll(r.Body)
	r.LogEntry.Problems.Add(string(details))
	r.WriteHeader(http.StatusNoContent)
}
