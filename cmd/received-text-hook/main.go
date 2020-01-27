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
		println("received-text-hook: ", err)
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	cgi.Serve(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			number  = r.FormValue("From")
			body    = r.FormValue("Body")
			message *model.TextMessage
		)
		db.Open("serv.db")
		tx := db.Begin()
		if message = tx.FetchTextMessageByNumber(number); message == nil {
			println("received-text-hook: incoming message from unknown phone number: ", number)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		for _, r := range message.Recipients {
			if r.Number == number {
				r.Responses = append(r.Responses, &model.TextResponse{
					Response:  body,
					Timestamp: time.Now(),
				})
				break
			}
		}
		tx.SaveTextMessage(message)
		tx.Commit()
		w.WriteHeader(http.StatusNoContent)
	}))
}
