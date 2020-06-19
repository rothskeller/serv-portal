// volunteer-hours is a cron job that handles periodic tasks related to
// reporting volunteer hours.  It normally decides what to do based on the day
// of month it's invoked, but that can be overridden by putting a day of month
// on the command line.
package main

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"mime/quotedprintable"
	"os"
	"strconv"
	ttemplate "text/template"
	"time"

	"sunnyvaleserv.org/portal/api/email"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/sendmail"
)

func main() {
	var now = time.Now()
	var dayOfMonth = now.Day()

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
	if len(os.Args) > 1 {
		if dom, err := strconv.Atoi(os.Args[1]); err == nil {
			dayOfMonth = dom
		}
	}
	store.Open("serv.db")
	entry := log.New("", "volunteer-hours")
	defer entry.Log()
	tx := store.Begin(entry)
	switch dayOfMonth {
	case 1:
		// Verify that placeholder events exist in the current month.
		makePlaceholders(tx, now)
		// Send hours requests for the past month.
		sendRequests(tx, time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local))
	case 8:
		// Send hours reminders for the past month.
		sendReminders(tx, time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local))
	case 11:
		// Send the accumulated hours to Volgistics.
		// submitHours(tx, time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local))
	default:
		fmt.Fprintf(os.Stderr, "ERROR: nothing to do on the %dth", dayOfMonth)
		os.Exit(1)
	}
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

func sendRequests(tx *store.Tx, month time.Time) {
	var (
		mstr             string
		events           []*model.Event
		eatt             = make(map[model.EventID]map[model.PersonID]model.AttendanceInfo)
		people           = make(map[model.PersonID]*model.Person)
		requestFromRoles = make(map[model.RoleID]bool)
		auth             = tx.Authorizer()
		mailer           *sendmail.Mailer
		err              error
	)
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	// First, get a list of the events in the month, and the people with
	// hours at those events.
	mstr = month.Format("2006-01")
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		if e.Organization == model.OrgNone {
			continue
		}
		events = append(events, e)
		eatt[e.ID] = tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt[e.ID] {
			if ai.Minutes == 0 || ai.Type == model.AttendAsAuditor || ai.Type == model.AttendAsStudent {
				delete(eatt[e.ID], pid)
				continue
			}
			if people[pid] == nil {
				people[pid] = tx.FetchPerson(pid)
			}
		}
	}
	// If we have people with minutes reported, but who don't have a
	// Volgistics ID, we need to send them an email to that effect.
	for _, p := range people {
		if p.VolgisticsID == 0 {
			notifyNotInVolgistics(p, mailer, month)
			delete(people, p.ID)
		}
	}
	// Next, we need to add in all of the people who are in groups that we
	// request volunteer hours from (and who are in Volgistics).  For
	// starters, let's get a list of those groups.
	for _, g := range auth.FetchGroups(auth.AllGroups()) {
		if g.GetHours {
			for _, r := range auth.RolesG(g.ID) {
				requestFromRoles[r] = true
			}
		}
	}
	for _, p := range tx.FetchPeople() {
		if people[p.ID] == nil && p.VolgisticsID != 0 && !p.NoEmail && (p.Email != "" || p.Email2 != "") {
			for _, r := range auth.RolesP(p.ID) {
				if requestFromRoles[r] {
					people[p.ID] = p
					break
				}
			}
		}
		if people[p.ID] == nil && p.HoursReminder {
			tx.WillUpdatePerson(p)
			p.HoursReminder = false
			tx.UpdatePerson(p)
		}
	}
	// Send an email to each of those people.
	for _, p := range people {
		tx.WillUpdatePerson(p)
		p.HoursToken = util.RandomToken()
		p.HoursReminder = true
		tx.UpdatePerson(p)
		sendRequest(p, mailer, month, events, eatt, false)
	}
}

func notifyNotInVolgistics(person *model.Person, mailer *sendmail.Mailer, month time.Time) {
	var data struct {
		Name  string
		Month string
	}
	var (
		buf     bytes.Buffer
		qpw     *quotedprintable.Writer
		crlf    email.CRLFWriter
		toaddrs []string
	)
	data.Name = person.InformalName
	data.Month = month.Format("January 2006")
	crlf = email.NewCRLFWriter(&buf)
	if person.Email != "" && person.Email2 != "" {
		fmt.Fprintf(crlf, "To: %s <%s>, %s <%s>\n", person.InformalName, person.Email, person.InformalName, person.Email2)
		toaddrs = []string{person.Email, person.Email2}
	} else if person.Email != "" {
		fmt.Fprintf(crlf, "To: %s <%s>\n", person.InformalName, person.Email)
		toaddrs = []string{person.Email}
	} else {
		fmt.Fprintf(crlf, "To: %s <%s>\n", person.InformalName, person.Email2)
		toaddrs = []string{person.Email2}
	}
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>
Subject: SERV Volunteer Hours for %s
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0

--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`, data.Month)
	qpw = quotedprintable.NewWriter(&buf)
	noVolgisticsText.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`)
	qpw = quotedprintable.NewWriter(&buf)
	noVolgisticsHTML.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY--
`)
	if err := mailer.SendMessage("admin@sunnyvaleserv.org", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", person.InformalName, err)
	}
}

var noVolgisticsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, LISTOS, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please fill out this form:
    https://www.volgistics.com/ex/portal.dll/ap?AP=929478828
In the “City employee status or referral” box, please enter
    Rebecca Elizondo
    Department of Public Safety
and the names of the organizations you're volunteering for (CERT, LISTOS, SNAP, and/or SARES).  If you have any difficulties with this, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var noVolgisticsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, LISTOS, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please fill out <a href="https://www.volgistics.com/ex/portal.dll/ap?AP=929478828">this form</a>.  In the “City employee status or referral” box, please enter</div><div><br></div><div style="margin-left:2em">Rebecca Elizondo<br>Department of Public Safety</div><div><br></div><div>and the names of the organizations you're volunteering for (CERT, LISTOS, SNAP, and/or SARES).  If you have any difficulties with this, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))

func sendRequest(
	person *model.Person, mailer *sendmail.Mailer, month time.Time, events []*model.Event,
	eatt map[model.EventID]map[model.PersonID]model.AttendanceInfo, reminder bool,
) {
	var data struct {
		Name     string
		Month    string
		URL      string
		Deadline string
		Events   []*model.Event
		Hours    map[model.EventID]string
		Total    string
	}
	var (
		remindstr string
		total     float64
		buf       bytes.Buffer
		qpw       *quotedprintable.Writer
		crlf      email.CRLFWriter
		toaddrs   []string
		pevents   []*model.Event
	)
	if reminder {
		remindstr = "Reminder: "
	}
	for _, e := range events {
		if m := eatt[e.ID][person.ID].Minutes; m != 0 {
			pevents = append(pevents, e)
		}
	}
	data.Name = person.InformalName
	data.Month = month.Format("January 2006")
	data.URL = "https://sunnyvaleserv.org/volunteer-hours/" + person.HoursToken
	data.Deadline = time.Date(month.Year(), month.Month()+1, 10, 0, 0, 0, 0, time.Local).Format("January 2")
	data.Events = pevents
	data.Hours = make(map[model.EventID]string)
	for _, e := range pevents {
		data.Hours[e.ID] = fmt.Sprintf("%.1f Hours", float64(eatt[e.ID][person.ID].Minutes)/60)
		total += float64(eatt[e.ID][person.ID].Minutes)
	}
	data.Total = fmt.Sprintf("%.1f Hours", total/60)
	crlf = email.NewCRLFWriter(&buf)
	if person.Email != "" && person.Email2 != "" {
		fmt.Fprintf(crlf, "To: %s <%s>, %s <%s>\n", person.InformalName, person.Email, person.InformalName, person.Email2)
		toaddrs = []string{person.Email, person.Email2}
	} else if person.Email != "" {
		fmt.Fprintf(crlf, "To: %s <%s>\n", person.InformalName, person.Email)
		toaddrs = []string{person.Email}
	} else {
		fmt.Fprintf(crlf, "To: %s <%s>\n", person.InformalName, person.Email2)
		toaddrs = []string{person.Email2}
	}
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>
Subject: %sSERV Volunteer Hours for %s
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0

--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`, remindstr, data.Month)
	qpw = quotedprintable.NewWriter(&buf)
	withEventsText.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`)
	qpw = quotedprintable.NewWriter(&buf)
	withEventsHTML.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY--
`)
	if err := mailer.SendMessage("admin@sunnyvaleserv.org", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", person.InformalName, err)
	}
}

var withEventsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, LISTOS, and/or SNAP during {{ .Month }}.  Currently, {{ if .Events }}our records show:
{{ range .Events }}    {{ .Date }} {{ .Name }}: {{ index $.Hours .ID }}
{{ end }}    Total Hours: {{ .Total }}
If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit
    {{ .URL }}
and report it prior to {{ .Deadline }}.  If you have questions about reporting volunteer hours, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var withEventsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, LISTOS, and/or SNAP during {{ .Month }}.  Currently, {{ if .Events }}our records show:</div><div><br></div><table style="margin-left:2em">{{ range .Events }}<tr><td>{{ .Date }} {{ .Name }}</td><td style="text-align:right;padding-left:1em">{{ index $.Hours .ID }}</td></tr>{{ end }}<tr><td style="text-align:right">Total</td><td style="text-align:right;padding-left:1em">{{ .Total }}</td></tr></table><div><br></div><div>If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit <a href="{{ .URL }}">our web site</a> and report it prior to {{ .Deadline }}.  If you have questions about reporting volunteer hours, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))

func sendReminders(tx *store.Tx, month time.Time) {
	var (
		mstr   string
		events []*model.Event
		eatt   = make(map[model.EventID]map[model.PersonID]model.AttendanceInfo)
		people = make(map[model.PersonID]*model.Person)
		mailer *sendmail.Mailer
		err    error
	)
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	// First, get a list of the events in the month, and the attendance at
	// each.
	mstr = month.Format("2006-01")
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		if e.Organization == model.OrgNone {
			continue
		}
		events = append(events, e)
		eatt[e.ID] = tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt[e.ID] {
			if ai.Minutes == 0 || ai.Type == model.AttendAsAuditor || ai.Type == model.AttendAsStudent {
				delete(eatt[e.ID], pid)
				continue
			}
		}
	}
	// Next, get all the people to whom we need to send reminders.
	for _, p := range tx.FetchPeople() {
		if p.HoursReminder {
			people[p.ID] = p
		}
	}
	// Send an email to each of those people.
	for _, p := range people {
		sendRequest(p, mailer, month, events, eatt, true)
	}
}
