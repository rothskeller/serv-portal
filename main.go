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
	"sunnyvaleserv.org/portal/api/folder"
	"sunnyvaleserv.org/portal/api/list"
	"sunnyvaleserv.org/portal/api/person"
	"sunnyvaleserv.org/portal/api/report"
	"sunnyvaleserv.org/portal/api/role"
	"sunnyvaleserv.org/portal/api/search"
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
	case r.Method == "GET" && c[0] == "dl":
		return folder.GetDocument(r)
	case r.Method == "GET" && c[1] == "folders":
		return folder.GetFolder(r)
	case r.Method == "POST" && c[1] == "login" && c[2] == "":
		return authn.PostLogin(r)
	case r.Method == "POST" && c[1] == "password-reset" && c[2] == "":
		return authn.PostPasswordReset(r)
	case r.Method == "GET" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return authn.GetPasswordResetToken(r, c[2])
	case r.Method == "POST" && c[1] == "password-reset" && c[2] != "" && c[3] == "":
		return authn.PostPasswordResetToken(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "hours" && c[4] != "" && c[5] == "":
		return person.GPPersonHoursMonth(r, c[2], c[4])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "hours" && c[4] != "" && c[5] == "":
		return person.GPPersonHoursMonth(r, c[2], c[4])
	case r.Method == "POST" && c[1] == "unsubscribe" && c[2] != "" && c[3] != "" && c[4] == "":
		return email.PostUnsubscribeList(r, c[2], c[3])
	case r.Person == nil:
		return util.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
	case r.Method == "GET" && c[1] == "login" && c[2] == "":
		return authn.GetLogin(r)
	case r.Method == "POST" && c[1] == "logout" && c[2] == "":
		return authn.PostLogout(r)
	case r.Method == "GET" && c[1] == "document":
		return folder.GetDocument(r)
	case r.Method == "POST" && c[1] == "document":
		return folder.PostDocument(r)
	case r.Method == "DELETE" && c[1] == "document":
		return folder.DeleteDocument(r)
	case r.Method == "GET" && c[1] == "events" && c[2] == "":
		return event.GetEvents(r)
	case r.Method == "GET" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.GetEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "":
		return event.PostEvent(r, c[2])
	case r.Method == "POST" && c[1] == "events" && c[2] != "" && c[3] == "attendance" && c[4] == "":
		return event.PostEventAttendance(r, c[2])
	case r.Method == "POST" && c[1] == "folders":
		return folder.PostFolder(r)
	case r.Method == "DELETE" && c[1] == "folders":
		return folder.DeleteFolder(r)
	case r.Method == "GET" && c[1] == "lists" && c[2] == "":
		return list.GetLists(r)
	case r.Method == "GET" && c[1] == "lists" && c[2] != "" && c[3] == "":
		return list.GetList(r, c[2])
	case r.Method == "POST" && c[1] == "lists" && c[2] != "" && c[3] == "":
		return list.PostList(r, c[2])
	case r.Method == "DELETE" && c[1] == "lists" && c[2] != "" && c[3] == "":
		return list.DeleteList(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] == "":
		return person.GetPeople(r)
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "":
		return person.GetPerson(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "contact" && c[4] == "":
		return person.GetPersonContact(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "contact" && c[4] == "":
		return person.PostPersonContact(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "lists" && c[4] == "":
		return person.GetPersonLists(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "lists" && c[4] == "":
		return person.PostPersonLists(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "names" && c[4] == "":
		return person.GetPersonNames(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "names" && c[4] == "":
		return person.PostPersonNames(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "notes" && c[4] == "":
		return person.GetPersonNotes(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "notes" && c[4] == "":
		return person.PostPersonNotes(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "password" && c[4] == "":
		return person.GetPersonPassword(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "password" && c[4] == "":
		return person.PostPersonPassword(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "roles" && c[4] == "":
		return person.GetPersonRoles(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "roles" && c[4] == "":
		return person.PostPersonRoles(r, c[2])
	case r.Method == "GET" && c[1] == "people" && c[2] != "" && c[3] == "status" && c[4] == "":
		return person.GetPersonStatus(r, c[2])
	case r.Method == "POST" && c[1] == "people" && c[2] != "" && c[3] == "status" && c[4] == "":
		return person.PostPersonStatus(r, c[2])
	case r.Method == "GET" && c[1] == "reports" && c[2] == "attendance" && c[3] == "":
		return report.GetAttendance(r)
	case r.Method == "GET" && c[1] == "reports" && c[2] == "clearance" && c[3] == "":
		return report.GetClearance(r)
	case r.Method == "GET" && c[1] == "roles" && c[2] == "":
		return role.GetRoles(r)
	case r.Method == "POST" && c[1] == "roles" && c[2] == "":
		return role.PostRoles(r)
	case r.Method == "GET" && c[1] == "roles" && c[2] != "" && c[3] == "":
		return role.GetRole(r, c[2])
	case r.Method == "POST" && c[1] == "roles" && c[2] != "" && c[3] == "":
		return role.PostRole(r, c[2])
	case r.Method == "DELETE" && c[1] == "roles" && c[2] != "" && c[3] == "":
		return role.DeleteRole(r, c[2])
	case r.Method == "GET" && c[1] == "search" && c[2] == "":
		return search.GetSearch(r)
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
