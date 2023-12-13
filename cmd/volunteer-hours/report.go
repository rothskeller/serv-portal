package main

import (
	"bytes"
	"fmt"
	"html/template"
	"mime/quotedprintable"
	"sort"
	"time"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/taskperson"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/sendmail"
)

var orgToAssignment = map[enum.Org]int{
	enum.OrgAdmin:  1052,
	enum.OrgCERTD:  1047,
	enum.OrgCERTT:  1047,
	enum.OrgListos: 1048,
	enum.OrgSARES:  399,
	enum.OrgSNAP:   373,
}
var assnToName = map[int]string{
	0:    "TOTAL", // as used in this program
	373:  "SNAP",
	399:  "SARES",
	1047: "CERT",
	1048: "LISTOS",
	1052: "Admin",
}
var assnToLabel = map[int]string{ // as shown in Volgistics
	373:  "SNAP Volunteer [EMERGENCY PREPAREDNESS]",
	399:  "SARES Volunteer [EMERGENCY PREPAREDNESS]",
	1047: "CERT [EMERGENCY PREPAREDNESS]",
	1048: "LISTOS [EMERGENCY PREPAREDNESS]",
	1052: "SERV Admin [EMERGENCY PREPAREDNESS]",
}

type einfo struct {
	Date       string
	Name       string
	Volunteers int
	Hours      uint
	Assignment string
}

type rdata struct {
	Month        string
	ByGroup      map[int]uint
	Events       []*einfo
	Leaders      []*pinfo
	Unregistered []*pinfo
}

func reportHours(st *store.Store) {
	var (
		mstr   string
		people = make(map[person.ID]*pinfo)
		events = make(map[event.ID]map[string]*einfo)
		report rdata
	)
	mstr = time.Time(mflag).Format("2006-01")
	report.Month = time.Time(mflag).Format("January 2006")
	report.ByGroup = make(map[int]uint)
	taskperson.MinutesBetween(st, mstr+"-01", mstr+"-32", func(eid event.ID, pid person.ID, org enum.Org, minutes uint) {
		var ei *einfo
		var pi *pinfo

		assn := orgToAssignment[org]
		if assn == 0 {
			return
		}
		aname := assnToName[assn]
		if events[eid] == nil {
			events[eid] = make(map[string]*einfo)
		}
		if ei = events[eid][aname]; ei == nil {
			ei = &einfo{Assignment: aname}
			events[eid][aname] = ei
		}
		if pi = people[pid]; pi == nil {
			pi = new(pinfo)
			people[pid] = pi
		}
		pi.Total += minutes
		report.ByGroup[assn] += minutes
		report.ByGroup[0] += minutes
		ei.Volunteers++
		ei.Hours += minutes
	})
	for eid, emap := range events {
		for _, ei := range emap {
			e := event.WithID(st, eid, event.FName|event.FStart)
			ei.Name = e.Name()
			ei.Date = e.Start()[:10]
			ei.Hours = (ei.Hours + 59) / 60
			report.Events = append(report.Events, ei)
		}
	}
	for pid, pi := range people {
		p := person.WithID(st, pid, person.FInformalName|person.FVolgisticsID)
		pi.Name = p.InformalName()
		pi.VolgisticsID = p.VolgisticsID()
		pi.Total = (pi.Total + 59) / 60
		report.Leaders = append(report.Leaders, pi)
	}
	for assn := range report.ByGroup {
		report.ByGroup[assn] = (report.ByGroup[assn] + 59) / 60
	}
	sort.Slice(report.Events, func(i, j int) bool {
		if report.Events[i].Date != report.Events[j].Date {
			return report.Events[i].Date < report.Events[j].Date
		}
		if report.Events[i].Name != report.Events[j].Name {
			return report.Events[i].Name < report.Events[j].Name
		}
		return report.Events[i].Assignment < report.Events[j].Assignment
	})
	sort.Slice(report.Leaders, func(i, j int) bool {
		if report.Leaders[i].Total != report.Leaders[j].Total {
			return report.Leaders[i].Total > report.Leaders[j].Total
		}
		return report.Leaders[i].Name < report.Leaders[j].Name
	})
	if len(report.Leaders) > 10 {
		report.Leaders = report.Leaders[:10]
	}
	for _, pi := range people {
		if pi.VolgisticsID == 0 {
			report.Unregistered = append(report.Unregistered, pi)
		}
	}
	sort.Slice(report.Unregistered, func(i, j int) bool {
		return report.Unregistered[i].Name < report.Unregistered[j].Name
	})
	sendReport(&report)
}

func sendReport(report *rdata) {
	var (
		buf    bytes.Buffer
		toaddr string
	)
	if *dflag {
		toaddr = "admin@sunnyvaleserv.org"
	} else {
		toaddr = "volunteer-hours@sunnyvaleserv.org"
	}
	crlf := sendmail.NewCRLFWriter(&buf)
	fmt.Fprintf(crlf, `From: Sunnyvale SERV <cert@sunnyvale.ca.gov>
To: volunteer-hours@sunnyvaleserv.org
Date: %s
Subject: Sunnyvale SERV Volunteer Hours for %s
Content-Type: text/html
Content-Transfer-Encoding: quoted-printable

`, time.Now().Format(time.RFC1123Z), report.Month)
	qpw := quotedprintable.NewWriter(&buf)
	if err := reportTemplate.Execute(qpw, report); err != nil {
		panic(err)
	}
	qpw.Close()
	if err := sendmail.SendMessage(config.Get("fromAddr"), []string{toaddr}, buf.Bytes()); err != nil {
		panic(err)
	}
}

var reportTemplate = template.Must(template.New("").Parse(`
<html>
  <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
  </head>
  <body style="font-family:Arial,Helvetica,sans-serif">
    <h1>SunnyvaleSERV.org Volunteer Hours for {{ .Month }}</h1>
    <p>
      Volunteer hours for {{ .Month }} that were reported on SunnyvaleSERV.org have been automatically recorded in Volgistics.
    </p>
    <table>
      <tr>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Group</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Hours</td>
      </tr>
      {{ $even := true }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="padding:0.2em">CERT</td>
        <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 1047 }}</td>
      </tr>{{ $even = not $even }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="padding:0.2em">LISTOS</td>
        <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 1048 }}</td>
      </tr>{{ $even = not $even }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="padding:0.2em">SARES</td>
        <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 399 }}</td>
      </tr>{{ $even = not $even }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="padding:0.2em">SNAP</td>
        <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 373 }}</td>
      </tr>{{ $even = not $even }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="padding:0.2em">Admin</td>
        <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 1052 }}</td>
      </tr>{{ $even = not $even }}
      <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
        <td style="font-weight:bold;padding:0.2em">TOTAL</td>
        <td style="font-weight:bold;text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ index .ByGroup 0 }}</td>
      </tr>
    </table>
    <h2>Leader Board</h2>
    <table>
      <tr>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Volunteer</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Hours</td>
      </tr>
      {{ $even := true }}
      {{ range .Leaders }}
        <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
          <td style="padding:0.2em">{{ .Name }}</td>
          <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ .Total }}</td>
        </tr>
        {{ $even = not $even }}
      {{ end }}
    </table>
    <h2>Events</h2>
    <table>
      <tr>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Date</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Event</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Group</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Volunteers</td>
        <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Hours</td>
      </tr>
      {{ $even := true }}
      {{ range .Events }}
        <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
          <td style="padding:0.2em">{{ .Date }}</td>
          <td style="padding:0.2em">{{ .Name }}</td>
          <td style="padding:0.2em">{{ .Assignment }}</td>
          <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ .Volunteers }}</td>
          <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ .Hours }}</td>
        </tr>
        {{ $even = not $even }}
      {{ end }}
    </table>
    {{ if .Unregistered }}
      <h2>Unregistered Volunteers (hours not recorded in Volgistics)</h2>
      <table>
        <tr>
          <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Volunteer</td>
          <td style="background-color:#538135;color:#FFFFFF;font-weight:bold;padding:0.2em">Hours</td>
        </tr>
        {{ $even := true }}
        {{ range .Unregistered }}
          <tr style="background-color:{{ if $even }}#A8D08D{{ else }}#BFBFBF{{ end }}">
            <td style="padding:0.2em">{{ .Name }}</td>
            <td style="text-align:right;padding:0.2em 0.2em 0.2em 1em">{{ .Total }}</td>
          </tr>
        {{ end }}
      </table>
    {{ end }}
  </body>
</html>`))
