// This program dumps all or part of the SERV database contents in JSON format.
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/db"
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: serv-dump object-type
    where object-type is one of:
        audit
	event
	group
	person
	role
	session
	text_message
	venue
    or an abbreviation of one of those.
`)
	os.Exit(2)
}

func main() {
	var (
		tx *db.Tx
	)
	if len(os.Args) != 2 || len(os.Args[1]) == 0 {
		usage()
	}
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		db.Open("/home/snyserv/sunnyvaleserv.org/data/serv.db")
	case "/Users/stever":
		db.Open("/Users/stever/src/serv-portal/data/serv.db")
	default:
		db.Open("./serv.db")
	}
	tx = db.Begin()
	switch {
	case strings.HasPrefix("audit", os.Args[1]):
		dumpAudit(tx)
	case strings.HasPrefix("events", os.Args[1]):
		dumpEvents(tx)
	case strings.HasPrefix("groups", os.Args[1]):
		dumpGroups(tx)
	case strings.HasPrefix("person", os.Args[1]) || strings.HasPrefix("people", os.Args[1]):
		dumpPeople(tx)
	case strings.HasPrefix("roles", os.Args[1]):
		dumpRoles(tx)
	case strings.HasPrefix("sessions", os.Args[1]):
		dumpSessions(tx)
	case strings.HasPrefix("text_messages", os.Args[1]):
		dumpTextMessages(tx)
	case strings.HasPrefix("venues", os.Args[1]):
		dumpVenues(tx)
	default:
		usage()
	}
	tx.Rollback()
}

func dumpAudit(tx *db.Tx) {}

func dumpEvents(tx *db.Tx) {
	for _, e := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		out.RawString(`{"id":`)
		out.Int(int(e.ID))
		out.RawString(`,"name":`)
		out.String(e.Name)
		out.RawString(`,"date":`)
		out.String(e.Date)
		out.RawString(`,"start":`)
		out.String(e.Start)
		out.RawString(`,"end":`)
		out.String(e.End)
		if e.Venue != 0 {
			out.RawString(`,"venue":`)
			out.Int(int(e.Venue))
		}
		if e.Details != "" {
			out.RawString(`,"details":`)
			out.String(e.Details)
		}
		out.RawString(`,"type":`)
		out.Int(int(e.Type))
		out.RawString(`,"groups":[`)
		for i, g := range e.Groups {
			if i != 0 {
				out.RawByte(',')
			}
			out.Int(int(g))
		}
		out.RawByte(']')
		if e.SccAresID != "" {
			out.RawString(`,"sccAresID":`)
			out.String(e.SccAresID)
		}
		out.RawString(`,"attendance":[`)
		first := true
		for p, ai := range tx.FetchAttendanceByEvent(e) {
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"person":`)
			out.Int(int(p))
			out.RawString(`,"type":`)
			out.Int(int(ai.Type))
			out.RawString(`,"minutes":`)
			out.Uint16(ai.Minutes)
			out.RawByte('}')
		}
		out.RawString(`]}`)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpGroups(tx *db.Tx) {
	for _, g := range tx.FetchGroups() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		g.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpPeople(tx *db.Tx) {
	for _, p := range tx.FetchPeople() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		p.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpRoles(tx *db.Tx) {
	for _, r := range tx.FetchRoles() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		r.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpSessions(tx *db.Tx) {
	for _, s := range tx.FetchSessions() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		out.RawString(`{"token":`)
		out.String(string(s.Token))
		out.RawString(`,"person":`)
		out.Int(int(s.Person.ID))
		out.RawString(`,"expires":`)
		out.Raw(s.Expires.MarshalJSON())
		out.RawByte('}')
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpTextMessages(tx *db.Tx) {
	for _, t := range tx.FetchTextMessages() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		out.RawString(`{"id":`)
		out.Int(int(t.ID))
		out.RawString(`,"sender":`)
		out.Int(int(t.Sender))
		out.RawString(`,"groups":[`)
		for i, g := range t.Groups {
			if i != 0 {
				out.RawByte(',')
			}
			out.Int(int(g))
		}
		out.RawString(`],"timestamp":`)
		out.Raw(t.Timestamp.MarshalJSON())
		out.RawString(`,"message":`)
		out.String(t.Message)
		out.RawString(`,"deliveries":[`)
		for i, d := range tx.FetchTextDeliveries(t.ID) {
			if i != 0 {
				out.RawByte(',')
			}
			d.MarshalEasyJSON(&out)
		}
		out.RawString(`]}`)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpVenues(tx *db.Tx) {
	for _, v := range tx.FetchVenues() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		v.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}
