package main

import (
	"fmt"
	"os"
	"sort"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/util/config"

	"zombiezen.com/go/sqlite"
)

func main() {
	var (
		dbconn *sqlite.Conn
		list   *maillist.List
		err    error
	)
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: serv-list listname")
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
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if dbconn, err = sqlite.OpenConn(config.Get("databaseFilename"), sqlite.OpenReadOnly|sqlite.OpenNoMutex); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: open DB: %s\n", err)
		os.Exit(1)
	}
	if list = maillist.GetList(dbconn, os.Args[1]); list == nil {
		fmt.Fprintf(os.Stderr, "ERROR: no such list %q\n", os.Args[1])
		os.Exit(1)
	}
	if list.DisplayName != list.Name {
		fmt.Printf("==== LIST %s (from sender \"via %s\")\n", list.Name, list.DisplayName)
	} else {
		fmt.Printf("==== LIST %s\n", list.Name)
	}
	fmt.Printf("== Sent to «addr» %s.\n", list.Reason)
	if list.NoUnsubscribe {
		fmt.Println("== No unsubscribe link.")
	}
	if list.Senders.Len() != 0 {
		fmt.Println("== Unmoderated Senders:")
		ss := list.Senders.UnsortedList()
		sort.Strings(ss)
		for _, s := range ss {
			fmt.Println(s)
		}
	} else {
		fmt.Println("== No unmoderated senders")
	}
	if len(list.Recipients) != 0 {
		fmt.Println("== Recipients")
		var ss []string
		for e, r := range list.Recipients {
			ss = append(ss, fmt.Sprintf("%s <%s>", r.Name, e))
		}
		sort.Strings(ss)
		for _, s := range ss {
			fmt.Println(s)
		}
	} else {
		fmt.Println("== No recipients")
	}
}
