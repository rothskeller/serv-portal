// Webhook that receives notification of status changes for outbound text
// messages.  It is invoked as a CGI "script" by the Dreamhost web server.
//
// This program expects to be run in the web root directory, which must contain
// a mode-700 "data" subdirectory.  The data subdirectory must contain the
// serv.db database and the config.json configuration file.  The audit.log and
// request.log log files will be created there.
package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"time"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	// Change working directory to the data subdirectory of the CGI script
	// location.  This directory should be mode 700 so that it not directly
	// readable by the web server.
	if err := os.Chdir("data"); err != nil {
		println("text-status-hook: ", err)
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			number  = r.FormValue("To")
			status  = r.FormValue("MessageStatus")
			message *textmsg.TextMessage
			p       *person.Person
		)
		entry := log.New("", "text-status-hook")
		defer entry.Log()
		store.Connect(context.Background(), entry, func(st *store.Store) {
			if message = textmsg.WithNumber(st, number, textmsg.FID); message == nil {
				println("text-status-hook: unknown recipient phone number: ", number)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid recipient phone number.")
				return
			}
			if p = textrecip.WithNumber(st, message.ID(), number, person.FID|person.FInformalName); p == nil {
				println("text-status-hook: unmatched recipient phone number: ", number)
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintln(w, "Invalid recipient phone number.")
				return
			}
			st.Transaction(func() {
				textrecip.UpdateStatus(st, message, p, status, time.Now())
			})
		})
		w.WriteHeader(http.StatusNoContent)
	}))
}
