package main

import (
	"context"
	"fmt"
	"os"
	"time"

	ics "github.com/arran4/golang-ical"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	const eventFields = event.FID | event.FStart | event.FEnd | event.FName | event.FDetails
	const venueFields = venue.FName
	var (
		entry *log.Entry
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
	cal = ics.NewCalendar()
	cal.SetProductId("SunnyvaleSERV.org")
	cal.SetVersion("2.0")
	now = time.Now()
	start = time.Date(now.Year(), now.Month()-6, now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")
	entry = log.New("", "gen-ical")
	store.Connect(context.Background(), entry, func(st *store.Store) {
		event.AllBetween(st, start, "2099-12-31", eventFields, venueFields, func(e *event.Event, v *venue.Venue) {
			ie := cal.AddEvent(fmt.Sprintf("%d@sunnyvaleserv.org", e.ID()))
			ie.SetDtStampTime(time.Now())
			ie.SetSummary(e.Name())
			start, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
			ie.SetStartAt(start)
			end, _ := time.ParseInLocation("2006-01-02T15:04", e.End(), time.Local)
			ie.SetEndAt(end)
			ie.SetURL(fmt.Sprintf("https://sunnyvaleserv.org/events/%d", e.ID()))
			if e.Details() != "" {
				ie.SetDescription(e.Details())
				// TODO render links in plain text
			}
			if v != nil {
				ie.SetLocation(v.Name())
			}
		})
	})
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
