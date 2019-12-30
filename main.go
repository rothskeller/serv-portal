// Main program for the serv.rothskeller.net/portal server.
//
// This program handles requests to https://serv.rothskeller.net/portal/*, for
// management of the Sunnyvale Emergency Response Volunteers (SERV) database.
// It is invoked as a CGI "script" by the Dreamhost web server.
//
// This program expects to be run in the web root directory, which must contain
// a mode-700 "data" subdirectory.  The data subdirectory must contain the
// serv.db database and the config.json configuration file.  The audit.log log
// file will be created there.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"net/url"
	"os"
	"strings"

	"serv.rothskeller.net/portal/auth"
	"serv.rothskeller.net/portal/db"
	"serv.rothskeller.net/portal/event"
	"serv.rothskeller.net/portal/person"
	"serv.rothskeller.net/portal/report"
	"serv.rothskeller.net/portal/role"
	"serv.rothskeller.net/portal/team"
	"serv.rothskeller.net/portal/util"
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
	case r.Method == "GET" && r.Path == "/login":
		return auth.GetLogin(r)
	case r.Method == "POST" && r.Path == "/login":
		return auth.PostLogin(r)
	case r.Person == nil:
		http.Redirect(r, r.Request, "/login?target="+url.QueryEscape(r.Path), http.StatusSeeOther)
		return nil
	case r.Method == "GET" && (r.Path == "/" || r.Path == "/index.html"):
		http.Redirect(r, r.Request, "/events", http.StatusSeeOther)
		return nil
	case r.Method == "GET" && r.Path == "/events":
		return event.ListEvents(r)
	case r.Method == "GET" && c[0] == "events" && c[1] != "" && c[2] == "":
		return event.EditEvent(r, c[1])
	case r.Method == "POST" && c[0] == "events" && c[1] != "" && c[2] == "":
		return event.EditEvent(r, c[1])
	case r.Method == "GET" && c[0] == "events" && c[1] != "" && c[2] == "attendance" && c[3] == "":
		return event.GetEventAttendance(r, c[1])
	case r.Method == "POST" && c[0] == "events" && c[1] != "" && c[2] == "attendance" && c[3] == "":
		return event.PostEventAttendance(r, c[1])
	case r.Method == "GET" && r.Path == "/people":
		return person.ListPeople(r)
	case r.Method == "GET" && c[0] == "people" && c[1] != "" && c[2] == "":
		return person.EditPerson(r, c[1])
	case r.Method == "POST" && c[0] == "people" && c[1] != "" && c[2] == "":
		return person.EditPerson(r, c[1])
	case r.Method == "GET" && r.Path == "/reports":
		return report.GetIndex(r)
	case r.Method == "GET" && r.Path == "/reports/cert-attendance":
		return report.CERTAttendanceReport(r)
	case r.Method == "GET" && r.Path == "/teams":
		return team.ListTeams(r)
	case r.Method == "GET" && c[0] == "teams" && c[1] != "" && c[2] == "":
		return team.EditTeam(r, c[1])
	case r.Method == "POST" && c[0] == "teams" && c[1] != "" && c[2] == "":
		return team.EditTeam(r, c[1])
	case r.Method == "GET" && c[0] == "teams" && c[1] != "" && c[2] == "roles" && c[3] != "" && c[4] == "":
		return role.EditRole(r, c[1], c[3])
	case r.Method == "POST" && c[0] == "teams" && c[1] != "" && c[2] == "roles" && c[3] != "" && c[4] == "":
		return role.EditRole(r, c[1], c[3])
	}
	return util.NotFound
}
