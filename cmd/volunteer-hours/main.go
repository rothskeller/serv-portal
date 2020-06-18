// volunteer-hours is a cron job that handles periodic tasks related to
// reporting volunteer hours.
package main

import (
	"fmt"
	"os"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	var now = time.Now()

	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err := os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err := os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	store.Open("serv.db")
	entry := log.New("", "volunteer-hours")
	defer entry.Log()
	tx := store.Begin(entry)
	// Verify that placeholder events exist in the current month and the
	// subsequent month.
	makePlaceholders(tx, now)
	makePlaceholders(tx, time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local))
	tx.Commit()
}

func makePlaceholders(tx *store.Tx, month time.Time) {
	var (
		mstr  = month.Format("2006-01")
		found = make(map[model.Organization]bool)
	)
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		if e.Type == model.EventHours {
			found[e.Organization] = true
		}
	}
	mstr = time.Date(month.Year(), month.Month()+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Second).Format("2006-01-02")
	for _, o := range model.CurrentOrganizations {
		month = month.Add(-time.Second)
		if !found[o] {
			var e = model.Event{
				Date:         mstr,
				Start:        "23:59",
				End:          "23:59",
				Name:         fmt.Sprintf("Other %s Hours", model.OrganizationNames[o]),
				Organization: o,
				Type:         model.EventHours,
			}
			tx.CreateEvent(&e)
		}
	}
}
