// Webhook that receives notification of status changes for outbound text
// messages.  It is invoked as a CGI "script" by the Dreamhost web server.
//
// This program expects to be run in the web root directory, which must contain
// a mode-700 "data" subdirectory.  The data subdirectory must contain the
// serv.db database and the config.json configuration file.  The audit.log and
// request.log log files will be created there.
package main

import (
	"fmt"
	"net/http"
	"net/http/cgi"
	"os"
	"time"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

func main() {
	// Change working directory to the data subdirectory of the CGI script
	// location.  This directory should be mode 700 so that it not directly
	// readable by the web server.
	if err := os.Chdir("data"); err != nil {
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			number     = r.FormValue("recipient")
			message    = model.TextMessageID(util.ParseID(r.FormValue("reference")))
			status     = r.FormValue("status")
			statusTime time.Time
			delivery   *model.TextDelivery
			err        error
		)
		if statusTime, err = time.Parse(time.RFC3339, r.FormValue("statusDatetime")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid timestamp.")
			return
		}
		db.Open("serv.db")
		tx := db.Begin()
		if delivery = tx.FetchTextDelivery(message, number); delivery == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Status on unknown message.")
			return
		}
		delivery.Status = status
		delivery.Timestamp = statusTime.In(time.Local)
		tx.SaveTextDelivery(delivery, message, number)
		tx.Commit()
		w.WriteHeader(http.StatusNoContent)
	}))
}
