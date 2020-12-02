// volunteer-hours handles tasks related to reporting volunteer hours.  It is
// normally invoked as a cron job at various different times for different
// tasks.
//
// usage: volunteer-hours [-m YYYY-MM] [-p people] [-d] request|remind|submit|report...
//     -m YYYY-MM specifies the target month (default "last month")
//     -p person specifies a target person (ID or username)
//     -d specifies debug mode; emails to go admin only
//     "request" means to send an email requesting hours
//     "remind" means to send an email reminder for submitting hours
//     "submit" means to submit hours to Volgistics
//     "report" means to email a summary report
package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
)

var tx *store.Tx
var mflag monthArg
var dflag = flag.Bool("d", false, "debug (emails to admin only)")
var kflag = flag.Bool("k", false, "keep existing HoursTokens")
var pflags = peopleList{}

func main() {
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
	store.Open("serv.db")
	entry := log.New("", "volunteer-hours")
	defer entry.Log()
	tx = store.Begin(entry)
	makePlaceholders(tx)
	flag.Var(&mflag, "m", "target month (YYYY-MM, default last month)")
	flag.Var(pflags, "p", "target person (ID or username)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `usage: volunteer-hours [-m YYYY-MM] [-p people] [-dk] request|remind|submit|report...
     -m YYYY-MM specifies the target month (default "last month")
     -p person specifies a target person (ID or username)
     -d specifies debug mode; emails to go admin only
     -k keeps existing HoursTokens and HoursReminders
     "request" means to send an email requesting hours
     "remind" means to send an email reminder for submitting hours
     "submit" means to submit hours to Volgistics
     "report" means to email a summary report
`)
		os.Exit(2)
	}
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: no operation specified\n")
		flag.Usage()
	}
	for _, op := range flag.Args() {
		switch op {
		case "request":
			sendRequests(tx)
		case "remind":
			sendReminders(tx)
		case "submit":
			submitHours(tx)
		case "report":
			reportHours(tx)
		default:
			fmt.Fprintf(os.Stderr, "ERROR: invalid operation %q\n", op)
			flag.Usage()
		}
	}
	tx.Commit()
}

type peopleList map[model.PersonID]bool

func (pl peopleList) Set(v string) (err error) {
	var pid int
	var person *model.Person

	if pid, err = strconv.Atoi(v); err == nil && pid > 0 {
		person = tx.FetchPerson(model.PersonID(pid))
	} else {
		person = tx.FetchPersonByUsername(v)
	}
	if person == nil {
		return fmt.Errorf("no such person: %s", v)
	}
	pl[person.ID] = true
	return nil
}

func (pl peopleList) String() string {
	var sb strings.Builder
	var first = true
	for pid := range pl {
		if first {
			first = false
		} else {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d", pid)
	}
	return sb.String()
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

func makePlaceholders(tx *store.Tx) {
	var (
		now   = time.Now()
		mstr  = now.Format("2006-01")
		found = make(map[model.Org]bool)
	)
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		if e.Type == model.EventHours {
			found[e.Org] = true
		}
	}
	mstr = time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Second).Format("2006-01-02")
	for _, o := range model.AllOrgs {
		if !found[o] {
			var e = model.Event{
				Date:  mstr,
				Start: "23:59",
				End:   "23:59",
				Name:  fmt.Sprintf("Other %s Hours", orgNames[o]),
				Org:   o,
				Type:  model.EventHours,
			}
			tx.CreateEvent(&e)
		}
	}
}

var orgNames = map[model.Org]string{
	model.OrgAdmin:  "Admin",
	model.OrgCERTD:  "CERT Deployment",
	model.OrgCERTT:  "CERT Training",
	model.OrgListos: "Listos",
	model.OrgSARES:  "SARES",
	model.OrgSNAP:   "SNAP",
}
