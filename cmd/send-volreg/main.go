// send-volreg sends someone's volunteer registration to Volgistics.  This
// should happen automatically, but can be called manually if a resend is
// needed.
package main

import (
	"fmt"
	"os"

	"sunnyvaleserv.org/portal/api/person"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	var (
		entry *log.Entry
		tx    *store.Tx
		p     *model.Person
		err   error
	)
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "usage: send-volreg person-id [interest...]\n")
		os.Exit(2)
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
	store.Open("serv.db")
	entry = log.New("", "send-volreg")
	defer entry.Log()
	tx = store.Begin(entry)
	if p = tx.FetchPerson(model.PersonID(util.ParseID(os.Args[1]))); p == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no such person %s\n", os.Args[1])
		os.Exit(1)
	}
	if !p.VolgisticsPending {
		fmt.Fprintf(os.Stderr, "ERROR: %s does not have a pending volgistics registration\n", p.InformalName)
		os.Exit(1)
	}
	if err = person.SendVolunteerRegistration(p, os.Args[2:]); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
	}
	tx.Commit()
}
