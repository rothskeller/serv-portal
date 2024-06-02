package main

import (
	"bytes"
	"context"
	"fmt"
	htemplate "html/template"
	"mime/quotedprintable"
	"net/mail"
	"os"
	ttemplate "text/template"
	"time"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/sendmail"
)

type minfo struct {
	tid     task.ID
	minutes uint
}

func sendRequests(st *store.Store) {
	const personFields = person.FID | person.FEmail | person.FEmail2 | person.FFlags | person.FVolgisticsID | person.FPrivLevels | person.FHoursToken | person.FInformalName | person.FCallSign
	var (
		mstr   string
		people = make(map[person.ID]*person.Person)
		mailer *sendmail.Mailer
		pm     = make(map[person.ID][]*minfo)
		enames = make(map[event.ID]string)
		tnames = make(map[task.ID]string)
		err    error
	)
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	// First, get a list of the events in the month, and the people with
	// hours at those events.
	mstr = time.Time(mflag).Format("2006-01")
	taskperson.MinutesBetween(st, mstr+"-01", mstr+"-32", func(eid event.ID, tid task.ID, pid person.ID, org enum.Org, minutes uint) {
		pm[pid] = append(pm[pid], &minfo{tid, minutes})
		enames[eid] = ""
		tnames[tid] = ""
	})
	for eid := range enames {
		e := event.WithID(st, eid, event.FStart|event.FName)
		enames[eid] = e.Start()[:10] + " " + e.Name()
	}
	for tid := range tnames {
		t := task.WithID(st, tid, task.FEvent|task.FName)
		if t.Name() == "Tracking" {
			tnames[tid] = enames[t.Event()]
		} else {
			tnames[tid] = enames[t.Event()] + " (" + t.Name() + ")"
		}
	}
	for pid := range pm {
		people[pid] = person.WithID(st, pid, personFields)
		if people[pid].Flags()&person.NoEmail != 0 || (people[pid].Email() == "" && people[pid].Email2() == "") {
			delete(people, pid)
		}
	}
	// If we have people with minutes reported, but who don't have a
	// Volgistics ID, we need to send them an email to that effect.
	for pid, p := range people {
		if p.VolgisticsID() == 0 {
			if len(pflag) == 0 || pflag[p.ID()] {
				notifyNotInVolgistics(p, mailer)
			}
			delete(people, pid)
		}
	}
	// Next, we need to add in all of the people who are members of orgs,
	// and also in Volgistics.
	person.All(st, personFields, func(p *person.Person) {
		if people[p.ID()] != nil {
			return // already have them
		}
		if p.VolgisticsID() == 0 || !p.HasPrivLevel(0, enum.PrivMember) {
			return // not a volunteer
		}
		if p.Flags()&person.NoEmail != 0 || (p.Email() == "" && p.Email2() == "") {
			return // can't reach them
		}
		people[p.ID()] = p.Clone()
	})
	// Set the hours reminder flag for those people and clear it on anyone
	// else.
	if !*kflag {
		st.Transaction(func() {
			person.ClearAllHoursReminders(st)
			for _, p := range people {
				if len(pflag) == 0 || pflag[p.ID()] {
					up := p.Updater()
					up.HoursToken = util.RandomToken()
					up.Flags |= person.HoursReminder
					p.Update(st, up, person.FHoursToken|person.FFlags)
				}
			}
		})
	}
	// Send an email to each of those people.
	for _, p := range people {
		if len(pflag) == 0 || pflag[p.ID()] {
			sendRequest(p, mailer, pm[p.ID()], tnames, false)
		}
	}
}

func notifyNotInVolgistics(person *person.Person, mailer *sendmail.Mailer) {
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
	data.Name = person.InformalName()
	data.Month = time.Time(mflag).Format("January 2006")
	crlf = sendmail.NewCRLFWriter(&buf)
	if person.Email() != "" && person.Email2() != "" {
		var ma1 = mail.Address{Name: person.InformalName(), Address: person.Email()}
		var ma2 = mail.Address{Name: person.InformalName(), Address: person.Email2()}
		fmt.Fprintf(crlf, "To: %s, %s\n", &ma1, &ma2)
		toaddrs = []string{person.Email(), person.Email2()}
	} else if person.Email() != "" {
		var ma = mail.Address{Name: person.InformalName(), Address: person.Email()}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email()}
	} else {
		var ma = mail.Address{Name: person.InformalName(), Address: person.Email2()}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{person.Email2()}
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
	if err := mailer.SendMessage(context.Background(), "admin@sunnyvaleserv.org", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", person.InformalName(), err)
	}
}

var noVolgisticsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please visit your profile on the Sunnyvale SERV website (https://SunnyvaleSERV.org) and click the "Register" button in the "Volunteer Status" area.  If you have any difficulties with this, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var noVolgisticsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you did for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  However, you are not currently registered as a City of Sunnyvale volunteer, so we cannot record your volunteer work.  To register as a City of Sunnyvale volunteer, please visit your profile on the <a href="https://sunnyvaleserv.org/">Sunnyvale SERV website</a> and click the “Register” button in the “Volunteer Status” area.  If you have any difficulties with this, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))

func sendReminders(st *store.Store) {
	const personFields = person.FID | person.FEmail | person.FEmail2 | person.FFlags | person.FInformalName | person.FHoursToken
	var (
		mstr   string
		mailer *sendmail.Mailer
		pm     = make(map[person.ID][]*minfo)
		enames = make(map[event.ID]string)
		tnames = make(map[task.ID]string)
		err    error
	)
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	// First, get a list of the events in the month, and the people with
	// hours at those events.
	mstr = time.Time(mflag).Format("2006-01")
	taskperson.MinutesBetween(st, mstr+"-01", mstr+"-32", func(eid event.ID, tid task.ID, pid person.ID, org enum.Org, minutes uint) {
		pm[pid] = append(pm[pid], &minfo{tid, minutes})
		enames[eid] = ""
		tnames[tid] = ""
	})
	for eid := range enames {
		e := event.WithID(st, eid, event.FStart|event.FName)
		enames[eid] = e.Start()[:10] + " " + e.Name()
	}
	for tid := range tnames {
		t := task.WithID(st, tid, task.FEvent|task.FName)
		if t.Name() == "Tracking" {
			tnames[tid] = enames[t.Event()]
		} else {
			tnames[tid] = enames[t.Event()] + " (" + t.Name() + ")"
		}
	}
	// Next, send to all the people who need reminders.
	person.All(st, personFields, func(p *person.Person) {
		if p.Flags()&person.HoursReminder != 0 && p.Flags()&person.NoEmail == 0 && (p.Email() != "" || p.Email2() != "") {
			if len(pflag) == 0 || pflag[p.ID()] {
				sendRequest(p, mailer, pm[p.ID()], tnames, true)
			}
		}
	})
}

func sendRequest(p *person.Person, mailer *sendmail.Mailer, tasks []*minfo, tnames map[task.ID]string, reminder bool) {
	type tinfo struct {
		Name  string
		Hours string
	}
	var data struct {
		Name     string
		Month    string
		URL      string
		Deadline string
		Tasks    []*tinfo
		Total    string
	}
	var (
		remindstr string
		total     float64
		buf       bytes.Buffer
		qpw       *quotedprintable.Writer
		crlf      sendmail.CRLFWriter
		toaddrs   []string
	)
	if reminder {
		remindstr = "Reminder: "
	}
	for _, t := range tasks {
		data.Tasks = append(data.Tasks, &tinfo{tnames[t.tid], fmt.Sprintf("%.1f Hours", float64(t.minutes)/60)})
		total += float64(t.minutes)
	}
	data.Name = p.InformalName()
	data.Month = time.Time(mflag).Format("January 2006")
	data.URL = "https://sunnyvaleserv.org/volunteer-hours/" + p.HoursToken()
	data.Deadline = time.Date(time.Time(mflag).Year(), time.Time(mflag).Month()+1, 10, 0, 0, 0, 0, time.Local).Format("January 2")
	data.Total = fmt.Sprintf("%.1f Hours", total/60)
	crlf = sendmail.NewCRLFWriter(&buf)
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0
Subject: %sSERV Volunteer Hours for %s
Date: %s
`, remindstr, data.Month, time.Now().Format(time.RFC1123Z))
	if p.Email() != "" && p.Email2() != "" {
		var ma1 = mail.Address{Name: p.InformalName(), Address: p.Email()}
		var ma2 = mail.Address{Name: p.InformalName(), Address: p.Email2()}
		fmt.Fprintf(crlf, "To: %s, %s\n", &ma1, &ma2)
		toaddrs = []string{p.Email(), p.Email2()}
	} else if p.Email() != "" {
		var ma = mail.Address{Name: p.InformalName(), Address: p.Email()}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{p.Email()}
	} else {
		var ma = mail.Address{Name: p.InformalName(), Address: p.Email2()}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{p.Email2()}
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
	if err := mailer.SendMessage(context.Background(), "admin@sunnyvaleserv.org", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", p.InformalName(), err)
	}
}

var withEventsText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  Currently, {{ if .Tasks }}our records show:
{{ range .Tasks }}    {{ .Name }}: {{ .Hours }}
{{ end }}    Total Hours: {{ .Total }}
If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit
    {{ .URL }}
and report it prior to {{ .Deadline }}.  If you have questions about reporting volunteer hours, just reply to this email.

Many thanks,
Sunnyvale OES
`))

var withEventsHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The Sunnyvale Office of Emergency Services would like to be sure that you get credit for all of the volunteer work you may have done for SERV, CERT, Listos, SARES, and/or SNAP during {{ .Month }}.  Currently, {{ if .Tasks }}our records show:</div><div><br></div><table style="margin-left:2em">{{ range .Tasks }}<tr><td>{{ .Name }}</td><td style="text-align:right;padding-left:1em">{{ .Hours }}</td></tr>{{ end }}<tr><td style="text-align:right">Total</td><td style="text-align:right;padding-left:1em">{{ .Total }}</td></tr></table><div><br></div><div>If these records are incorrect, or you have any additional volunteer time to report,{{ else }}our records do not show any volunteer time from you during {{ .Month }}.  If you have volunteer time to report,{{ end }} please visit <a href="{{ .URL }}">our web site</a> and report it prior to {{ .Deadline }}.</div><div style="margin:1em 0"><a style="color:#fff;background-color:#007bff;border:1px solid #007bff;border-radius:4px;padding:6px 12px;line-height:1.5;text-align:center;vertical-align:middle;display:inline-block;cursor:pointer;user-select:none;text-decoration:none" href="{{ .URL }}">Report Hours</a><div style="display:inline-block;margin-left:16px;color:#888">This button takes you directly to your Activity page without needing to log in.</div></div><div>If you have questions about reporting volunteer hours, just reply to this email.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))
