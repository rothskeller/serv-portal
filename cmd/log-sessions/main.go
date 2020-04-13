// This program adds a log entry for every existing session.  It is run at the
// start of each month so that log reports during that month can map session
// tokens to user names for the sessions started in prior months.
package main

import (
	"os"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	var tx *store.Tx

	switch os.Getenv("HOME") {
	case "/home/snyserv":
		os.Chdir("/home/snyserv/sunnyvaleserv.org/data")
	case "/Users/stever":
		os.Chdir("/Users/stever/src/serv-portal/data")
	}
	store.Open("serv.db")
	tx = store.Begin(nil)
	for _, session := range tx.FetchSessions() {
		var entry = log.New("", "continued-session")
		entry.Session = string(session.Token)
		entry.Change("continued session %s for person %q", session.Token, session.Person.InformalName)
		entry.Log()
	}
	tx.Commit()
}
