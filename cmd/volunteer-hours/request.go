package main

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"mime/quotedprintable"
	"net/mail"
	"os"
	ttemplate "text/template"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/sendmail"
)

func sendRequests(tx *store.Tx) {
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
	// First, get a list of the events in the month, and the people with
	// hours at those events.
	mstr = time.Time(mflag).Format("2006-01")
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		events = append(events, e)
		eatt[e.ID] = tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt[e.ID] {
			if ai.Minutes == 0 || ai.Type == model.AttendAsAuditor || ai.Type == model.AttendAsStudent {
				delete(eatt[e.ID], pid)
				continue
			}
			if people[pid] == nil && (len(pflags) == 0 || pflags[pid]) {
				person := tx.FetchPerson(pid)
				if !person.NoEmail && (person.Email != "" || person.Email2 != "") {
					people[pid] = person
				}
			}
		}
	}
	// If we have people with minutes reported, but who don't have a
	// Volgistics ID, we need to send them an email to that effect.
	for _, p := range people {
		if p.VolgisticsID == 0 {
			notifyNotInVolgistics(p, mailer)
			delete(people, p.ID)
		}
	}
	// Next, we need to add in all of the people who are members of orgs,
	// and also in Volgistics.
	for _, p := range tx.FetchPeople() {
		if len(pflags) != 0 && !pflags[p.ID] {
			continue
		}
		if people[p.ID] == nil && p.VolgisticsID != 0 && !p.NoEmail && (p.Email != "" || p.Email2 != "") &&
			p.HasPrivLevel(model.PrivMember) {
			people[p.ID] = p
		}
	}
	// Send an email to each of those people.
	for _, p := range people {
		if !*kflag {
			tx.WillUpdatePerson(p)
			p.HoursToken = util.RandomToken()
			p.HoursReminder = true
			tx.UpdatePerson(p)
		}
		sendRequest(p, mailer, events, eatt, false)
	}
}

func notifyNotInVolgistics(person *model.Person, mailer *sendmail.Mailer) {
	var data struct {
		Name  string
		Month string
	}
	var (
		buf     bytes.Buffer
		qpw     *quotedprintable.Writer
		crlf    sendmail.CRLFWriter
		toaddrs []string
	)
	data.Name = person.InformalName
	data.Month = time.Time(mflag).Format("January 2006")
	crlf = sendmail.NewCRLFWriter(&buf)
	if person.Email != "" && person.Email2 != "" {
		var ma1 = mail.Address{Name: person.InformalName, Address: person.Email}
		var ma2 = mail.Address{Name: person.InformalName, Address: person.Email2}
		fmt.Fprintf(crlf, "To: %s, %s\n", &ma1, &ma2)
		toaddrs = []string{person.Email, person.Email2}
	} else if person.Email != "" {
		var ma = mail.Address{Name: person.InformalName, Address: person.Email}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email}
	} else {
		var ma = mail.Address{Name: person.InformalName, Address: person.Email2}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email2}
	}
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <serv@sunnyvale.ca.gov>
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
	if err := mailer.SendMessage("serv@sunnyvale.ca.gov", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", person.InformalName, err)
	}
}

var noVolgisticsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please fill out this form:
    https://www.volgistics.com/ex/portal.dll/ap?AP=929478828
In the “City employee status or referral” box, please enter
    Rebecca Elizondo
    Department of Public Safety
and the names of the organizations you're volunteering for (CERT, LISTOS, SNAP, and/or SARES).  If you have any difficulties with this, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var noVolgisticsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please fill out <a href="https://www.volgistics.com/ex/portal.dll/ap?AP=929478828">this form</a>.  In the “City employee status or referral” box, please enter</div><div><br></div><div style="margin-left:2em">Rebecca Elizondo<br>Department of Public Safety</div><div><br></div><div>and the names of the organizations you're volunteering for (CERT, LISTOS, SNAP, and/or SARES).  If you have any difficulties with this, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))

func sendReminders(tx *store.Tx) {
	var (
		mstr   string
		events []*model.Event
		eatt   = make(map[model.EventID]map[model.PersonID]model.AttendanceInfo)
		mailer *sendmail.Mailer
		err    error
	)
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	// First, get a list of the events in the month, and the attendance at
	// each.
	mstr = time.Time(mflag).Format("2006-01")
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		events = append(events, e)
		eatt[e.ID] = tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt[e.ID] {
			if ai.Minutes == 0 || ai.Type == model.AttendAsAuditor || ai.Type == model.AttendAsStudent {
				delete(eatt[e.ID], pid)
				continue
			}
		}
	}
	// Next, send to all the people who need reminders.
	for _, p := range tx.FetchPeople() {
		if len(pflags) != 0 && !pflags[p.ID] {
			continue
		}
		if p.HoursReminder && !p.NoEmail && (p.Email != "" || p.Email2 != "") {
			sendRequest(p, mailer, events, eatt, true)
		}
	}
}

func sendRequest(
	person *model.Person, mailer *sendmail.Mailer, events []*model.Event,
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
		crlf      sendmail.CRLFWriter
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
	data.Month = time.Time(mflag).Format("January 2006")
	data.URL = "https://sunnyvaleserv.org/volunteer-hours/" + person.HoursToken
	data.Deadline = time.Date(time.Time(mflag).Year(), time.Time(mflag).Month()+1, 10, 0, 0, 0, 0, time.Local).Format("January 2")
	data.Events = pevents
	data.Hours = make(map[model.EventID]string)
	for _, e := range pevents {
		data.Hours[e.ID] = fmt.Sprintf("%.1f Hours", float64(eatt[e.ID][person.ID].Minutes)/60)
		total += float64(eatt[e.ID][person.ID].Minutes)
	}
	data.Total = fmt.Sprintf("%.1f Hours", total/60)
	crlf = sendmail.NewCRLFWriter(&buf)
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <serv@sunnyvale.ca.gov>
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0
Subject: %sSERV Volunteer Hours for %s
Date: %s
`, remindstr, data.Month, time.Now().Format(time.RFC1123Z))
	if person.Email != "" && person.Email2 != "" {
		var ma1 = mail.Address{Name: person.InformalName, Address: person.Email}
		var ma2 = mail.Address{Name: person.InformalName, Address: person.Email2}
		fmt.Fprintf(crlf, "To: %s, %s\n", &ma1, &ma2)
		toaddrs = []string{person.Email, person.Email2}
	} else if person.Email != "" {
		var ma = mail.Address{Name: person.InformalName, Address: person.Email}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email}
	} else {
		var ma = mail.Address{Name: person.InformalName, Address: person.Email2}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email2}
	}
	if *dflag {
		toaddrs = []string{"admin@sunnyvaleserv.org"}
	}
	fmt.Fprint(crlf, `
--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`)
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
	if err := mailer.SendMessage("serv@sunnyvale.ca.gov", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", person.InformalName, err)
	}
}

var withEventsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  Currently, {{ if .Events }}our records show:
{{ range .Events }}    {{ .Date }} {{ .Name }}: {{ index $.Hours .ID }}
{{ end }}    Total Hours: {{ .Total }}
If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit
    {{ .URL }}
and report it prior to {{ .Deadline }}.  If you have questions about reporting volunteer hours, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var withEventsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  Currently, {{ if .Events }}our records show:</div><div><br></div><table style="margin-left:2em">{{ range .Events }}<tr><td>{{ .Date }} {{ .Name }}</td><td style="text-align:right;padding-left:1em">{{ index $.Hours .ID }}</td></tr>{{ end }}<tr><td style="text-align:right">Total</td><td style="text-align:right;padding-left:1em">{{ .Total }}</td></tr></table><div><br></div><div>If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit <a href="{{ .URL }}">our web site</a> and report it prior to {{ .Deadline }}.  If you have questions about reporting volunteer hours, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))
