// This program dumps all or part of the SERV database contents in JSON format.
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
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

func dumpAudit(tx *db.Tx) {
	tx.FetchAudit(func(timestamp time.Time, username, request, otype string, id, data interface{}) {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		out.RawString(`{"timestamp":`)
		out.Raw(timestamp.MarshalJSON())
		out.RawString(`,"username":`)
		out.String(username)
		out.RawString(`,"request":`)
		out.String(request)
		out.RawString(`,"type":`)
		out.String(otype)
		switch otype {
		case "authz":
			out.RawString(`,"groups":[`)
			for i, g := range data.(model.AuthzData).Groups {
				if i != 0 {
					out.RawByte(',')
				}
				g.MarshalEasyJSON(&out)
			}
			out.RawString(`],"roles":[`)
			for i, r := range data.(model.AuthzData).Roles {
				if i != 0 {
					out.RawByte(',')
				}
				dumpRole(tx, &out, r)
			}
			out.RawByte(']')
		case "event":
			out.RawString(`,"event":`)
			var event = data.(model.Event)
			dumpEvent(tx, &out, &event)
		case "person":
			out.RawString(`,"person":`)
			var person = data.(model.Person)
			dumpPerson(tx, &out, &person)
		case "session":
			out.RawString(`,"session":`)
			var session = data.(model.Session)
			dumpSession(tx, &out, &session)
		case "text_message":
			out.RawString(`,"textMessage":`)
			var textMessage = data.(model.TextMessage)
			dumpTextMessage(tx, &out, &textMessage)
		case "venues":
			out.RawString(`,"venues":[`)
			for i, v := range data.(model.Venues).Venues {
				if i != 0 {
					out.RawByte(',')
				}
				v.MarshalEasyJSON(&out)
			}
		default:
			panic("unknown object type in audit record: " + otype)
		}
		out.RawByte('}')
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	})
}

func dumpEvents(tx *db.Tx) {
	for _, e := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpEvent(tx, &out, e)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpEvent(tx *db.Tx, out *jwriter.Writer, e *model.Event) {
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
		out.RawString(`,"venue":{"id":`)
		out.Int(int(e.Venue))
		out.RawString(`,"name":`)
		out.String(venueName(tx, e.Venue))
		out.RawByte('}')
	}
	if e.Details != "" {
		out.RawString(`,"details":`)
		out.String(e.Details)
	}
	out.RawString(`,"types":[`)
	first := true
	for _, t := range model.AllEventTypes {
		if e.Type&t != 0 {
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(model.EventTypeNames[t])
		}
	}
	out.RawString(`],"groups":[`)
	for i, g := range e.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g))
		out.RawString(`,"name":`)
		out.String(groupName(tx, g))
		out.RawByte('}')
	}
	out.RawByte(']')
	if e.SccAresID != "" {
		out.RawString(`,"sccAresID":`)
		out.String(e.SccAresID)
	}
	out.RawString(`,"attendance":[`)
	first = true
	for p, ai := range tx.FetchAttendanceByEvent(e) {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"person":`)
		out.Int(int(p))
		out.RawString(`,"sortName":`)
		out.String(personName(tx, p))
		out.RawString(`,"type":`)
		out.String(model.AttendanceTypeNames[ai.Type])
		out.RawString(`,"minutes":`)
		out.Uint16(ai.Minutes)
		out.RawByte('}')
	}
	out.RawString(`]}`)
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
		dumpPerson(tx, &out, p)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpPerson(tx *db.Tx, out *jwriter.Writer, p *model.Person) {
	out.RawString(`{"id":`)
	out.Int(int(p.ID))
	if p.Username != "" {
		out.RawString(`,"username":`)
		out.String(p.Username)
	}
	out.RawString(`,"informalName":`)
	out.String(p.InformalName)
	out.RawString(`,"formalName":`)
	out.String(p.FormalName)
	out.RawString(`,"sortName":`)
	out.String(p.SortName)
	if p.CallSign != "" {
		out.RawString(`,"callSign":`)
		out.String(p.CallSign)
	}
	if p.Email != "" {
		out.RawString(`,"email":`)
		out.String(p.Email)
	}
	if p.Email2 != "" {
		out.RawString(`,"email2":`)
		out.String(p.Email2)
	}
	if len(p.Emails) != 0 {
		out.RawString(`,"emails":[`)
		for i, e := range p.Emails {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"email":`)
			out.String(e.Email)
			if e.Label != "" {
				out.RawString(`,"label":`)
				out.String(e.Label)
			}
			if e.Bad {
				out.RawString(`,"bad":`)
				out.Bool(e.Bad)
			}
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	if p.HomeAddress.Address != "" {
		out.RawString(`,"homeAddress":`)
		p.HomeAddress.MarshalEasyJSON(out)
	}
	if p.WorkAddress.Address != "" || p.WorkAddress.SameAsHome {
		out.RawString(`,"workAddress":`)
		p.WorkAddress.MarshalEasyJSON(out)
	}
	if p.MailAddress.Address != "" || p.MailAddress.SameAsHome {
		out.RawString(`,"mailAddress":`)
		p.MailAddress.MarshalEasyJSON(out)
	}
	if p.CellPhone != "" {
		out.RawString(`,"cellPhone":`)
		out.String(p.CellPhone)
	}
	if p.HomePhone != "" {
		out.RawString(`,"homePhone":`)
		out.String(p.HomePhone)
	}
	if p.WorkPhone != "" {
		out.RawString(`,"workPhone":`)
		out.String(p.WorkPhone)
	}
	if len(p.Password) != 0 {
		out.RawString(`,"password":`)
		out.RawText(p.Password, nil)
	}
	if p.BadLoginCount != 0 {
		out.RawString(`,"badLoginCount":`)
		out.Int(p.BadLoginCount)
	}
	if !p.BadLoginTime.IsZero() {
		out.RawString(`,"badLoginTime":`)
		out.Raw(p.BadLoginTime.MarshalJSON())
	}
	if p.PWResetToken != "" {
		out.RawString(`,"pwresetToken":`)
		out.String(string(p.PWResetToken))
	}
	if !p.PWResetTime.IsZero() {
		out.RawString(`,"pwresetTime":`)
		out.Raw(p.PWResetTime.MarshalJSON())
	}
	if len(p.Roles) != 0 {
		out.RawString(`,"roles":[`)
		for i, r := range p.Roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(r))
			out.RawString(`,"name":`)
			out.String(roleName(tx, r))
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	out.RawString(`,"privileges":`)
	dumpPrivilegeMap(tx, out, p.Privileges)
	if len(p.Archive) != 0 {
		out.RawString(`,"archive":[`)
		for i, a := range p.Archive {
			if i != 0 {
				out.RawByte(',')
			}
			parts := strings.SplitN(a, "=", 2)
			out.RawByte('[')
			out.String(parts[0])
			if len(parts) > 1 {
				out.RawByte(',')
				out.String(parts[1])
			}
			out.RawByte(']')
		}
		out.RawByte(']')
	}
	out.RawByte('}')
}

func dumpPrivilegeMap(tx *db.Tx, out *jwriter.Writer, privileges model.PrivilegeMap) {
	first := true
	out.RawByte('[')
	for _, g := range tx.FetchGroups() {
		privs := privileges.Get(g)
		for _, p := range model.AllPrivileges {
			if privs&p != 0 {
				if first {
					first = false
				} else {
					out.RawByte(',')
				}
				out.RawString(`{"id":`)
				out.Int(int(g.ID))
				out.RawString(`,"name":`)
				out.String(g.Name)
				out.RawString(`,"privilege":`)
				out.String(model.PrivilegeNames[p])
				out.RawByte('}')
			}
		}
	}
	out.RawByte(']')
}

func dumpRoles(tx *db.Tx) {
	for _, r := range tx.FetchRoles() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpRole(tx, &out, r)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpRole(tx *db.Tx, out *jwriter.Writer, r *model.Role) {
	out.RawString(`{"id":`)
	out.Int(int(r.ID))
	if r.Tag != "" {
		out.RawString(`,"tag":`)
		out.String(string(r.Tag))
	}
	out.RawString(`,"name":`)
	out.String(r.Name)
	if r.Individual {
		out.RawString(`,"individual":`)
		out.Bool(r.Individual)
	}
	out.RawString(`,"privileges":`)
	dumpPrivilegeMap(tx, out, r.Privileges)
	out.RawByte('}')
}

func dumpSessions(tx *db.Tx) {
	for _, s := range tx.FetchSessions() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpSession(tx, &out, s)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpSession(tx *db.Tx, out *jwriter.Writer, s *model.Session) {
	out.RawString(`{"token":`)
	out.String(string(s.Token))
	out.RawString(`,"person":`)
	out.Int(int(s.Person.ID))
	out.RawString(`,"sortName":`)
	out.String(s.Person.SortName)
	out.RawString(`,"expires":`)
	out.Raw(s.Expires.MarshalJSON())
	out.RawByte('}')
}

func dumpTextMessages(tx *db.Tx) {
	for _, t := range tx.FetchTextMessages() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpTextMessage(tx, &out, t)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpTextMessage(tx *db.Tx, out *jwriter.Writer, t *model.TextMessage) {
	out.RawString(`{"id":`)
	out.Int(int(t.ID))
	out.RawString(`,"sender":{"id":`)
	out.Int(int(t.Sender))
	out.RawString(`,"sortName":`)
	out.String(personName(tx, t.Sender))
	out.RawString(`},"groups":[`)
	for i, g := range t.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g))
		out.RawString(`,"name":`)
		out.String(groupName(tx, g))
		out.RawByte('}')
	}
	out.RawString(`],"timestamp":`)
	out.Raw(t.Timestamp.MarshalJSON())
	out.RawString(`,"message":`)
	out.String(t.Message)
	out.RawString(`,"recipients":[`)
	for i, r := range t.Recipients {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(r.Recipient))
		out.RawString(`,"sortName":`)
		out.String(personName(tx, r.Recipient))
		out.RawString(`,"number":`)
		out.String(r.Number)
		out.RawString(`,"status":`)
		out.String(r.Status)
		out.RawString(`,"timestamp":`)
		out.Raw(r.Timestamp.MarshalJSON())
		out.RawString(`,"responses":[`)
		for i, resp := range r.Responses {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"response":`)
			out.String(resp.Response)
			out.RawString(`,"timestamp":`)
			out.Raw(resp.Timestamp.MarshalJSON())
			out.RawByte('}')
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
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

func groupName(tx *db.Tx, id model.GroupID) string {
	if v := tx.FetchGroup(id); v != nil {
		return v.Name
	}
	return ""
}
func personName(tx *db.Tx, id model.PersonID) string {
	if v := tx.FetchPerson(id); v != nil {
		return v.SortName
	}
	return ""
}
func roleName(tx *db.Tx, id model.RoleID) string {
	if v := tx.FetchRole(id); v != nil {
		return v.Name
	}
	return ""
}
func venueName(tx *db.Tx, id model.VenueID) string {
	if v := tx.FetchVenue(id); v != nil {
		return v.Name
	}
	return ""
}
