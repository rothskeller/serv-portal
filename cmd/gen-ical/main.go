package main

import (
	"fmt"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func main() {
	var (
		tx    *store.Tx
		cal   *ics.Calendar
		now   time.Time
		start string
		fh    *os.File
		err   error
	)
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err = os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err = os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	store.Open("serv.db")
	tx = store.Begin(nil)
	cal = ics.NewCalendar()
	cal.SetProductId("SunnyvaleSERV.org")
	cal.SetVersion("2.0")
	now = time.Now()
	start = time.Date(now.Year(), now.Month()-6, now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")
	for _, e := range tx.FetchEvents(start, "2099-12-31") {
		ie := cal.AddEvent(fmt.Sprintf("%d@sunnyvaleserv.org", e.ID))
		ie.SetDtStampTime(time.Now())
		ie.SetSummary(e.Name)
		start, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Date+" "+e.Start, time.Local)
		ie.SetStartAt(start)
		end, _ := time.ParseInLocation("2006-01-02 15:04:05", e.Date+" "+e.End, time.Local)
		ie.SetEndAt(end)
		ie.SetURL(fmt.Sprintf("https://sunnyvaleserv.org/events/%d", e.ID))
		if e.Organization != model.OrgNone {
			ie.AddProperty(ics.ComponentPropertyCategories, model.OrganizationNames[e.Organization])
		}
		if e.Details != "" {
			ie.SetDescription(e.Details)
			// TODO render links in plain text
		}
		if e.Venue != 0 {
			ie.SetLocation(tx.FetchVenue(e.Venue).Name)
		}
	}
	if fh, err = os.Create("../calendar.ics.new"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	fmt.Fprint(fh, cal.Serialize())
	fh.Close()
	if err = os.Rename("../calendar.ics.new", "../calendar.ics"); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
}
