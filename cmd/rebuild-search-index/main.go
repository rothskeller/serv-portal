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
	os.Chdir(`C:\SERV`)
	entry := log.New("", "rebuild-search-index")
	store.Connect(context.Background(), entry, func(st *store.Store) {
		st.Transaction(func() {
			search.RebuildSearchIndex(st)
		})
	})
	entry.Log()
}
