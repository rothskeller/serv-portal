// Main program for the serv.rothskeller.net/api server.
//
// This program handles requests to https://sunnyvaleserv.org/portal/*, for
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

	"sunnyvaleserv.org/portal/api/authn"
	"sunnyvaleserv.org/portal/api/email"
	"sunnyvaleserv.org/portal/api/event"
	"sunnyvaleserv.org/portal/api/group"
	"sunnyvaleserv.org/portal/api/person"
	"sunnyvaleserv.org/portal/api/report"
	"sunnyvaleserv.org/portal/api/role"
	"sunnyvaleserv.org/portal/api/text"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
)

var (
	txh store.Tx
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

// txWrapper opens the database and wraps the request in a transaction.  It also
// checks authorization, and sets the user's identity in r.Person and r.Auth if
// properly validated.
func txWrapper(r *util.Request) error {
	// Open the database and start a transaction.
	store.Open("serv.db")
	r.Tx = store.Begin(r.LogEntry)
	defer func() {
		r.Tx.Rollback()
	}()
	r.Auth = r.Tx.Authorizer()
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
		return authn.PostLogin(r)
	case r.Method == "POST" && c[1] == "password-reset" && c[2] == "":
		return authn.PostPasswordReset(r)
	case r.Method == "GET" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return authn.GetPasswordResetToken(r, c[2])
	case r.Method == "POST" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return authn.PostPasswordResetToken(r, c[2])
	case r.Method == "GET" && c[1] == "unsubscribe" && c[2] != "":
		return email.GetUnsubscribe(r, c[2])
	case r.Method == "POST" && c[1] == "unsubscribe" && c[2] != "" && c[3] == "":
		return email.PostUnsubscribe(r, c[2])
	case r.Method == "POST" && c[1] == "unsubscribe" && c[2] != "" && c[3] != "" && c[4] == "":
		return email.PostUnsubscribeList(r, c[2], c[3])
	case r.Person == nil:
		return util.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	case r.Method == "GET" && c[1] == "login" && c[2] == "":
		return authn.GetLogin(r)
	case r.Method == "POST" && c[1] == "logout" && c[2] == "":
		return authn.PostLogout(r)
	case r.Method == "GET" && c[1] == "emails" && c[2] == "":
		return email.GetEmails(r)
	case r.Method == "POST" && c[1] == "emails" && c[2] != "" && c[3] == "":
		return email.PostEmail(r, c[2])
	case r.Method == "GET" && c[1] == "events" && c[2] == "":
		return event.GetEvents(r)
	case r.Method == "GET" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.GetEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.PostEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "attendance" && c[4] == "":
		return event.PostEventAttendance(r, c[2])
	case r.Method == "GET" && c[1] == "groups" && c[2] == "":
		return group.GetGroups(r)
	case r.Method == "GET" && c[1] == "groups" && c[2] != "" && c[3] == "":
		return group.GetGroup(r, c[2])
	case r.Method == "POST" && c[1] == "groups" && c[2] != "" && c[3] == "":
		return group.PostGroup(r, c[2])
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
	case r.Method == "GET" && c[1] == "roles" && c[2] == "":
		return role.GetRoles(r)
	case r.Method == "GET" && c[1] == "roles" && c[2] != "" && c[3] == "":
		return role.GetRole(r, c[2])
	case r.Method == "POST" && c[1] == "roles" && c[2] != "" && c[3] == "":
		return role.PostRole(r, c[2])
	case r.Method == "GET" && c[1] == "sms" && c[2] == "":
		return text.GetSMS(r)
	case r.Method == "POST" && c[1] == "sms" && c[2] == "":
		return text.PostSMS(r)
	case r.Method == "GET" && c[1] == "sms" && c[2] == "NEW" && c[3] == "":
		return text.GetSMSNew(r)
	case r.Method == "GET" && c[1] == "sms" && c[2] != "" && c[3] == "":
		return text.GetSMS1(r, c[2])
	}
	return util.NotFound
}
