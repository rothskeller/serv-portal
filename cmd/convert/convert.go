package main

import (
	"context"
	_ "embed"
	"log"
	"os"

	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/search"
	"sunnyvaleserv.org/portal/util/config"
	nlog "sunnyvaleserv.org/portal/util/log"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

//go:embed schema.sql
var schema string

// //go:embed fixups.sql
// var fixups string

func main() {
	var (
		tx    *ostore.Tx
		conn  *sqlite.Conn
		entry *nlog.Entry
		err   error
	)

	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			log.Fatal(err)
		}
	case "/Users/stever":
		if err = os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("where am I?")
	}
	os.Remove(config.Get("databaseFilename"))
	os.RemoveAll("documents")
	search.EmptyEntireIndex()
	if conn, err = sqlite.OpenConn(config.Get("databaseFilename"), sqlite.OpenReadWrite|sqlite.OpenCreate); err != nil {
		log.Fatal(err)
	}
	if err = sqlitex.ExecuteScript(conn, schema, nil); err != nil {
		log.Fatal(err)
	}
	if err = conn.Close(); err != nil {
		log.Fatal(err)
	}
	ostore.Open("serv.db")
	tx = ostore.Begin(nil)
	entry = nlog.New("", "convert")
	err = nstore.Connect(context.Background(), entry, func(st *nstore.Store) {
		st.Transaction(func() {
			convertLists(tx, st)
			convertRoles(tx, st)
			convertPeople(tx, st)
			convertSubscriptions(tx, st)
			role.Recalculate(st)
			convertVenues(tx, st)
			convertEvents(tx, st)
			convertTextMessages(tx, st)
			convertSessions(tx, st)
			convertFolders(tx, st)
		})
	})
	entry.Log()
	if err != nil {
		log.Fatal(err)
	}
	tx.Rollback()
	// if conn, err = sqlite.OpenConn(config.Get("databaseFilename"), sqlite.OpenReadWrite); err != nil {
	// 	log.Fatal(err)
	// }
	// if err = sqlitex.ExecuteScript(conn, fixups, nil); err != nil {
	// 	log.Fatal(err)
	// }
	// if err = conn.Close(); err != nil {
	// 	log.Fatal(err)
	// }
}
