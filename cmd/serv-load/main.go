// serv-load reads JSON-formatted objects and installs them in the SERV Portal
// database, overwriting any previous versions of those objects.  Use with care.
//
// usage: serv-load object-type < json-file
//     or, more commonly,
// usage: serv-dump object-type | jq 'some-filter' | serv-load object-type
package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mailru/easyjson/jlexer"
	"sunnyvaleserv.org/portal/db"
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: serv-load object-type < json-file
    where object-type is one of:
        audit
	event
	group
	person
	role
	session
	text_message
	venue
    or an abbreviation of one of those.
`)
	os.Exit(2)
}

func main() {
	var (
		tx  *db.Tx
		buf []byte
		in  *jlexer.Lexer
		err error
	)
	if len(os.Args) != 2 || len(os.Args[1]) == 0 {
		usage()
	}
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err = os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	if buf, err = ioutil.ReadAll(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: stdin: %s\n", err)
		os.Exit(1)
	}
	maybeMakeBackup()
	db.Open("./serv.db")
	tx = db.Begin()
	in = &jlexer.Lexer{Data: buf}
	switch {
	// case strings.HasPrefix("audit", os.Args[1]):
	// 	dumpAudit(tx)
	// case strings.HasPrefix("events", os.Args[1]):
	// 	dumpEvents(tx)
	// case strings.HasPrefix("groups", os.Args[1]):
	// 	dumpGroups(tx)
	case strings.HasPrefix("person", os.Args[1]) || strings.HasPrefix("people", os.Args[1]):
		loadPeople(tx, in)
	// case strings.HasPrefix("roles", os.Args[1]):
	// 	dumpRoles(tx)
	// case strings.HasPrefix("sessions", os.Args[1]):
	// 	dumpSessions(tx)
	// case strings.HasPrefix("text_messages", os.Args[1]):
	// 	dumpTextMessages(tx)
	// case strings.HasPrefix("venues", os.Args[1]):
	// 	dumpVenues(tx)
	default:
		usage()
	}
	tx.Commit()
}