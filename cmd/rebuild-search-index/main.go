// This program rebuilds the search index.
package main

import (
	"context"
	"os"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/search"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		os.Chdir("/home/snyserv/sunnyvaleserv.org/data")
	case "/Users/stever":
		os.Chdir("/Users/stever/src/serv-portal/data")
	default:
	}
	entry := log.New("", "rebuild-search-index")
	store.Connect(context.Background(), entry, func(st *store.Store) {
		st.Transaction(func() {
			search.RebuildSearchIndex(st)
		})
	})
	entry.Log()
}
