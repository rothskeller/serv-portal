// volunteer-hours handles tasks related to reporting volunteer hours.  It is
// normally invoked as a cron job at various different times for different
// tasks.
//
// usage: volunteer-hours [-m YYYY-MM] [-p people] [-d] request|remind|submit|report...
//
//	-m YYYY-MM specifies the target month (default "last month")
//	-p person specifies a target person ID
//	-d specifies debug mode; emails to go admin only
//	"request" means to send an email requesting hours
//	"remind" means to send an email reminder for submitting hours
//	"submit" means to submit hours to Volgistics
//	"report" means to email a summary report
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util/log"
)

var mflag monthArg
var dflag = flag.Bool("d", false, "debug (emails to admin only)")
var kflag = flag.Bool("k", false, "keep existing HoursTokens")

func main() {
	var (
		loginID string
		entry   *log.Entry
	)
	mflag = monthArg(time.Now().AddDate(0, -1, 0))
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
	flag.Var(&mflag, "m", "target month (YYYY-MM, default last month)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: volunteer-hours [-m YYYY-MM] [-dk] request|remind|submit|report|status...
     -m YYYY-MM specifies the target month (default "last month")
     -d specifies debug mode; emails to go admin only
     -k keeps existing HoursTokens and HoursReminders
     "request" means to send an email requesting hours
     "remind" means to send an email reminder for submitting hours
     "submit" means to submit hours to Volgistics
     "report" means to email a summary report
     "status" means to update volunteer status in Volgistics
`)
		os.Exit(2)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: no operation specified\n")
		flag.Usage()
	}
	entry = log.New("", "volunteer-hours")
	defer entry.Log()
	store.Connect(context.Background(), entry, func(st *store.Store) {
		makePlaceholders(st)
		for _, op := range flag.Args() {
			switch op {
			case "request":
				// sendRequests(st)
			case "remind":
				// sendReminders(st)
			case "submit":
				if loginID == "" {
					loginID = logInToVolgistics()
				}
				// submitHours(st, loginID)
			case "report":
				reportHours(st)
			case "status":
				if loginID == "" {
					loginID = logInToVolgistics()
				}
				markActive(st, loginID)
			default:
				fmt.Fprintf(os.Stderr, "ERROR: invalid operation %q\n", op)
				flag.Usage()
			}
		}
	})
}

type monthArg time.Time

func (m *monthArg) Set(v string) (err error) {
	var t time.Time

	if t, err = time.ParseInLocation("2006-01", v, time.Local); err != nil {
		return err
	}
	*m = monthArg(t)
	return nil
}

func (m monthArg) String() string {
	if time.Time(m).IsZero() {
		return ""
	}
	return time.Time(m).Format("2006-01")
}

func makePlaceholders(st *store.Store) {
	var (
		ue    event.Updater
		ut    task.Updater
		found bool
		now   = time.Now()
		mstr  = now.Format("2006-01")
	)
	event.AllBetween(st, mstr+"-28", mstr+"-31", event.FFlags, 0, func(e *event.Event, _ *venue.Venue) {
		if e.Flags()*event.OtherHours != 0 {
			found = true
		}
	})
	if found {
		return
	}
	mstr = time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, time.Local).Format("2006-01-02T15:04")
	ue.Start, ue.End = mstr, mstr
	ue.Name = "Other Volunteer Hours"
	ue.Flags = event.OtherHours
	st.Transaction(func() {
		ut.Event = event.Create(st, &ue)
		for _, o := range enum.AllOrgs {
			ut.Name = o.Label()
			ut.Org = o
			ut.Flags = task.RecordHours
			task.Create(st, &ut)
		}
	})
}
