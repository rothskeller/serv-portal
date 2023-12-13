package homepage

import (
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/util/request"
)

// Serve handles "/" requests.
func Serve(r *request.Request) {
	if user := auth.SessionUser(r, 0, false); user == nil {
		servePublic(r)
	} else {
		serveUser(r, user)
	}
}
