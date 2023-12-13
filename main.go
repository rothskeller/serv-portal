package main

import (
	"net/http"
	"os"

	"sunnyvaleserv.org/portal/server"
)

func main() {
	os.Chdir("/Users/stever/src/serv-portal/data")
	if err := http.ListenAndServe("localhost:3000", server.Server); err != nil {
		panic(err)
	}
}
