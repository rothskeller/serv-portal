package main

import (
	"fmt"
	"net/http"
	"os"
	"syscall"

	"sunnyvaleserv.org/portal/server"
)

func main() {
	var err error

	ensureSingleton()
	if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if err = http.ListenAndServe("localhost:3000", server.Server); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}

// lockFH is the singleton lock file used in ensureSingleton.  It is declared at
// global scope so that it never gets garbage collected.
var lockFH *os.File

// ensureSingleton makes sure there is only one instance of server running at a
// time.  Redundant instances exit immediately and silently.
func ensureSingleton() {
	var err error

	// Open (or create) the run.lock file.
	if lockFH, err = os.OpenFile("run.lock", os.O_CREATE|os.O_WRONLY, 0666); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: open run.lock: %s", err)
		os.Exit(1)
	}
	// Acquire an exclusive lock on the run.lock file.
	switch err = syscall.Flock(int(lockFH.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err {
	case nil:
		// Lock successfully acquired, so we are the only running
		// instance.  We will hold the lock until our process exits.
		return
	case syscall.EWOULDBLOCK:
		// Another process has the lock, so there is already another
		// running instance.  Exit immediately and silently.
		os.Exit(0)
	default:
		// Unable to acquire the lock, for some reason other than
		// another process holding it.  Report the error and exit.
		fmt.Fprintf(os.Stderr, "ERROR: lock run.lock: %s", err)
		os.Exit(1)
	}
}
