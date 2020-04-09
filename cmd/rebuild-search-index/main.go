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
		store.Open("/home/snyserv/sunnyvaleserv.org/data/serv.db")
	case "/Users/stever":
		store.Open("/Users/stever/src/serv-portal/data/serv.db")
	default:
		store.Open("./serv.db")
	}
	tx = store.Begin(nil)
	tx.RebuildSearchIndex()
	tx.Commit()
}
