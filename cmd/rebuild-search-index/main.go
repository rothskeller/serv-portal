// This program rebuilds the search index.
package main

import (
	"os"

	"sunnyvaleserv.org/portal/store"
)

func main() {
	var tx *store.Tx

	switch os.Getenv("HOME") {
	case "/home/snyserv":
		os.Chdir("/home/snyserv/sunnyvaleserv.org/data")
	case "/Users/stever":
		os.Chdir("/Users/stever/src/serv-portal/data")
	default:
	}
	store.Open("./serv.db")
	tx = store.Begin(nil)
	tx.RebuildSearchIndex()
	tx.Commit()
}
