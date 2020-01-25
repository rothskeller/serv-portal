// Main program for the serv.rothskeller.net/api server.
//
// This program handles requests to https://rothskeller.net/serv/*, for
// management of the Sunnyvale Emergency Response Volunteers (SERV) database.
// It is invoked as a CGI "script" by the Dreamhost web server.
//
// This program expects to be run in the web root directory, which must contain
// a mode-700 "data" subdirectory.  The data subdirectory must contain the
// serv.db database and the config.json configuration file.  The audit.log and
// request.log log files will be created there.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"strings"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/db"
	"rothskeller.net/serv/event"
	"rothskeller.net/serv/person"
	"rothskeller.net/serv/report"
	"rothskeller.net/serv/text"
	"rothskeller.net/serv/util"
)

var (
	txh db.Tx
)

func main() {
	// Change working directory to the data subdirectory of the CGI script
	// location.  This directory should be mode 700 so that it not directly
	// readable by the web server.
	if err := os.Chdir("data"); err != nil {
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	// Run the request.
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		util.RunRequest(w, r, txWrapper)
	}))
}

// txWrapper opens the database and wraps the request in a transaction.
func txWrapper(r *util.Request) error {
	// Open the database and start a transaction.
	db.Open("serv.db")
	r.Tx = db.Begin()
	defer func() {
		r.Tx.Rollback()
	}()
	r.Tx.SetRequest(r.Method + " " + r.Path)
	return authWrapper(r)
}

// authWrapper looks for authorization cookies in the request and, if present,
// validates the session.  It never fails; it just doesn't set r.Person if the
// session is invalid.
func authWrapper(r *util.Request) error {
	util.ValidateSession(r)
	return router(r)
}

// router sends the request to the appropriate handler given its method and
// path.
func router(r *util.Request) error {
	c := strings.Split(r.Path[1:], "/")
	for len(c) < 6 {
		c = append(c, "")
	}
	switch {
	case r.Method == "POST" && c[1] == "login" && c[2] == "":
		return auth.PostLogin(r)
	case r.Method == "POST" && c[1] == "password-reset" && c[2] == "":
		return auth.PostPasswordReset(r)
	case r.Method == "GET" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return auth.GetPasswordResetToken(r, c[2])
	case r.Method == "POST" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return auth.PostPasswordResetToken(r, c[2])
	case r.Person == nil:
		return util.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	case r.Method == "GET" && c[1] == "login" && c[2] == "":
		return auth.GetLogin(r)
	case r.Method == "POST" && c[1] == "logout" && c[2] == "":
		return auth.PostLogout(r)
	case r.Method == "GET" && c[1] == "events" && c[2] == "":
		return event.GetEvents(r)
	case r.Method == "GET" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.GetEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.PostEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "attendance" && c[4] == "":
		return event.PostEventAttendance(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] == "":
		return person.GetPeople(r)
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "":
		return person.GetPerson(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "":
		return person.PostPerson(r, c[2])
	case r.Method == "GET" && c[1] == "reports" && c[2] == "":
		return report.GetIndex(r)
	case r.Method == "GET" && c[1] == "reports" && c[2] == "cert-attendance" && c[3] == "":
		return report.CERTAttendanceReport(r)
		/*
			case r.Method == "GET" && c[1] == "roles" && c[2] == "":
				return role.GetRoles(r)
			case r.Method == "GET" && c[1] == "roles" && c[2] != "" && c[3] == "":
				return role.GetRole(r, c[2])
			case r.Method == "POST" && c[1] == "roles" && c[2] != "" && c[3] == "":
				return role.PostRole(r, c[2])
			case r.Method == "POST" && c[1] == "roles" && c[2] != "" && c[3] == "reloadPrivs" && c[4] == "":
				return role.PostRoleReloadPrivs(r, c[2])
		*/
	case r.Method == "POST" && c[1] == "textMessage" && c[2] == "":
		return text.PostTextMessage(r)
	case r.Method == "GET" && c[1] == "textMessage" && c[2] != "" && c[3] == "":
		return text.GetTextMessage(r, c[2])
	}
	return util.NotFound
}
