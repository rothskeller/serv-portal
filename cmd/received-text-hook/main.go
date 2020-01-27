// Webhook that receives notification of incoming text messages.  It is invoked
// as a CGI "script" by the Dreamhost web server.
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
			number      = r.FormValue("originator")
			message     = r.FormValue("body")
			createdTime time.Time
			person      *model.Person
			delivery    *model.TextDelivery
			err         error
		)
		if createdTime, err = time.Parse(time.RFC3339, r.FormValue("createdDatetime")); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid timestamp.")
			return
		}
		if len(number) != 11 || number[0] != '1' {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Invalid originator phone number.")
			return
		}
		number = number[1:4] + "-" + number[4:7] + "-" + number[7:11]
		db.Open("serv.db")
		tx := db.Begin()
		if person = tx.FetchPersonByCellPhone(number); person == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Incoming message from unknown phone number.")
			return
		}
		if delivery = tx.FetchNewestTextDelivery(person.ID); delivery == nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintln(w, "Incoming reply without a corresponding outgoing message.")
			return
		}
		delivery.Responses = append(delivery.Responses, &model.TextResponse{Response: message, Timestamp: createdTime.In(time.Local)})
		tx.SaveTextDelivery(delivery)
		tx.Commit()
		w.WriteHeader(http.StatusNoContent)
	}))
}
