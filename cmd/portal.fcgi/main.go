package main

import (
	"fmt"
	"net/http/fcgi"
	"os"

	"sunnyvaleserv.org/portal/server"
	"sunnyvaleserv.org/portal/ui"
)

func main() {
	var (
		fh  *os.File
		err error
	)
	if fh, err = os.OpenFile("/tmp/serv-portal.err", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666); err == nil {
		os.Stderr = fh
	}
	if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if len(os.Args) == 2 && os.Args[1] == "-writeassets" {
		if err = ui.WriteAssetFiles(); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
		return
	}
	if err = fcgi.Serve(nil, server.Server); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: fcgi.Serve: %s\n", err)
		os.Exit(1)
	}
}
