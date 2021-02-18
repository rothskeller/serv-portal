// send-signups sends emails to people to let them know of new shifts available
// that they can sign up for.
package main

import (
	"bytes"
	"fmt"
	htemplate "html/template"
	"mime/quotedprintable"
	"net/mail"
	"net/url"
	"os"
	"sort"
	ttemplate "text/template"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/sendmail"
)

func main() {
	var (
		entry  *log.Entry
		tx     *store.Tx
		people []*model.Person
		mailer *sendmail.Mailer
		toSend = map[*model.Person]map[*model.Event]bool{}
		err    error
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
	entry = log.New("", "send-signups")
	defer entry.Log()
	tx = store.Begin(entry)
	people = tx.FetchPeople()
	for _, e := range tx.FetchEvents(time.Now().Format("2006-01-02"), "2099-12-31") {
		if len(e.Shifts) == 0 {
			continue
		}
		var touched = false
		for _, s := range e.Shifts {
			if !s.NewOpen {
				continue
			}
			for _, p := range people {
				if shouldSend(tx, e, s, p) {
					if toSend[p] == nil {
						toSend[p] = make(map[*model.Event]bool)
					}
					toSend[p][e] = true
				}
			}
			s.NewOpen = false
			touched = true
		}
		if touched {
			tx.UpdateEvent(e)
		}
	}
	if len(toSend) == 0 {
		os.Exit(0)
	}
	if mailer, err = sendmail.OpenMailer(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't open mailer: %s\n", err)
		os.Exit(1)
	}
	tx.Commit()
	for p, es := range toSend {
		sendEmail(mailer, p, es)
	}
}

func shouldSend(tx *store.Tx, e *model.Event, s *model.Shift, p *model.Person) bool {
	var found bool

	// People with email disabled, or no address, are obvious no-sends.
	if p.NoEmail || (p.Email == "" && p.Email2 == "") {
		return false
	}
	// Event org leaders always get notice.
	if p.Orgs[e.Org].PrivLevel >= model.PrivLeader {
		return true
	}
	// Make sure this person is holds an invited role.
	for _, r := range e.Roles {
		if _, ok := p.Roles[r]; ok {
			found = true
			break
		}
	}
	if !found {
		return false
	}
	// Make sure this person hasn't already declined this shift.
	for _, dp := range s.Declined {
		if dp == p.ID {
			return false
		}
	}
	// Make sure this person isn't signed up for any shift that overlaps
	// with this one.  (This also ensures they aren't already signed up for
	// this shift itself.)
	for _, os := range e.Shifts {
		if !(s.Start < os.End && s.End > os.Start) {
			continue // no overlap
		}
		for _, sp := range s.SignedUp {
			if sp == p.ID {
				return false
			}
		}
	}
	return true
}

func sendEmail(mailer *sendmail.Mailer, p *model.Person, emap map[*model.Event]bool) {
	var data struct {
		Name   string
		Token  string
		Events []string
	}
	var (
		buf     bytes.Buffer
		qpw     *quotedprintable.Writer
		crlf    sendmail.CRLFWriter
		toaddrs []string
	)
	data.Name = p.InformalName
	data.Token = url.PathEscape(p.UnsubscribeToken)
	data.Events = make([]string, 0, len(emap))
	for e := range emap {
		data.Events = append(data.Events, fmt.Sprintf("%s %s", e.Date, e.Name))
	}
	sort.Strings(data.Events)
	crlf = sendmail.NewCRLFWriter(&buf)
	if p.Email != "" && p.Email2 != "" {
		var ma1 = mail.Address{Name: p.InformalName, Address: p.Email}
		var ma2 = mail.Address{Name: p.InformalName, Address: p.Email2}
		fmt.Fprintf(crlf, "To: %s, %s\n", &ma1, &ma2)
		toaddrs = []string{p.Email, p.Email2}
	} else if p.Email != "" {
		var ma = mail.Address{Name: p.InformalName, Address: p.Email}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{p.Email}
	} else {
		var ma = mail.Address{Name: p.InformalName, Address: p.Email2}
		fmt.Fprintf(crlf, "To: %s\n", &ma)
		toaddrs = []string{p.Email2}
	}
	fmt.Fprintf(crlf, `From: SunnyvaleSERV.org <serv@sunnyvale.ca.gov>
Subject: SERV Volunteer Shifts Available
Content-Type: multipart/alternative; boundary="BOUNDARY"
MIME-Version: 1.0

--BOUNDARY
Content-Type: text/plain; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`)
	qpw = quotedprintable.NewWriter(&buf)
	shiftsAvailText.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY
Content-Type: text/html; charset=utf-8
Content-Transfer-Encoding: quoted-printable

`)
	qpw = quotedprintable.NewWriter(&buf)
	shiftsAvailHTML.Execute(qpw, data)
	qpw.Close()
	fmt.Fprint(crlf, `

--BOUNDARY--
`)
	if err := mailer.SendMessage("serv@sunnyvale.ca.gov", toaddrs, buf.Bytes()); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: can't send email to %s: %s\n", p.InformalName, err)
	}
}

var shiftsAvailText = ttemplate.Must(ttemplate.New("").Parse(`Hello, {{ .Name }},

The following upcoming {{ if eq (len .Events) 1 }}event has{{ else }}events have{{ end }} open shifts that you can sign up for:

{{ range .Events }}{{ . }}
{{ end }}
If you wish to sign up, please visit the Signups page on SunnyvaleSERV.org:
    https://sunnyvaleserv.org/events/signups/{{ .Token }}

Many thanks,
Sunnyvale OES
`))

var shiftsAvailHTML = htemplate.Must(htemplate.New("").Parse(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head><body><div>Hello, {{ .Name }},</div><div><br></div><div>The following upcoming {{ if eq (len .Events) 1 }}event has{{ else }}events have{{ end }} open shifts that you can sign up for:</div><div><br></div>{{ range .Events }}<div style="margin-left:2em">{{ . }}</div>{{ end }}<div><br></div><div>If you wish to sign up, please visit the <a href="https://sunnyvaleserv.org/events/signups/{{ .Token }}">Signups</a> page on SunnyvaleSERV.org.</div><div><br></div><div>Many thanks,<br>Sunnyvale OES</div></body></html>`))
