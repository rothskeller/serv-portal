package main

import (
	"os"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	os.Chdir("/home/snyserv/sunnyvaleserv.org/data")
	store.Open("serv.db")
	entry := log.New("", "add-unsub")
	tx := store.Begin(entry)
	for _, p := range tx.FetchPeople() {
		tx.WillUpdatePerson(p)
		p.UnsubscribeToken = util.RandomToken()
		tx.UpdatePerson(p)
	}
	tx.Commit()
}
