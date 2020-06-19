// This program dumps all or part of the SERV database contents in JSON format.
package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func usage() {
	fmt.Fprintf(os.Stderr, `usage: serv-dump object-type
    where object-type is one of:
        audit
	email_message
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
	var tx *store.Tx

	if len(os.Args) != 2 || len(os.Args[1]) == 0 {
		usage()
	}
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		store.Open("/home/snyserv/sunnyvaleserv.org/data/serv.db")
	case "/Users/stever":
		store.Open("/Users/stever/src/serv-portal/data/serv.db")
	default:
		store.Open("./serv.db")
	}
	tx = store.Begin(nil)
	switch {
	case (strings.HasPrefix("email_messages", os.Args[1]) && len(os.Args[1]) > 1) || os.Args[1] == "emails":
		dumpEmailMessages(tx)
	case strings.HasPrefix("events", os.Args[1]) && len(os.Args[1]) > 1:
		dumpEvents(tx)
	case strings.HasPrefix("folders", os.Args[1]):
		dumpFolders(tx)
	case strings.HasPrefix("groups", os.Args[1]):
		dumpGroups(tx)
	case strings.HasPrefix("person", os.Args[1]) || strings.HasPrefix("people", os.Args[1]):
		dumpPeople(tx)
	case strings.HasPrefix("roles", os.Args[1]):
		dumpRoles(tx)
	case strings.HasPrefix("sessions", os.Args[1]):
		dumpSessions(tx)
	case strings.HasPrefix("text_messages", os.Args[1]) || os.Args[1] == "texts":
		dumpTextMessages(tx)
	case strings.HasPrefix("venues", os.Args[1]):
		dumpVenues(tx)
	default:
		usage()
	}
	tx.Rollback()
}

func dumpEmailMessages(tx *store.Tx) {
	tx.FetchEmailMessages(func(em *model.EmailMessage) bool {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpEmailMessage(tx, &out, em)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
		return true
	})
}

func dumpEmailMessage(tx *store.Tx, out *jwriter.Writer, em *model.EmailMessage) {
	out.RawString(`{"id":`)
	out.Int(int(em.ID))
	out.RawString(`,"messageID":`)
	out.String(em.MessageID)
	out.RawString(`,"timestamp":`)
	out.Raw(em.Timestamp.MarshalJSON())
	out.RawString(`,"type":`)
	out.String(model.EmailMessageTypeNames[em.Type])
	out.RawString(`,"attention":`)
	out.Bool(em.Attention)
	out.RawString(`,"groups":[`)
	for i, g := range em.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g))
		out.RawString(`,"name":`)
		out.String(groupName(tx, g))
		out.RawByte('}')
	}
	out.RawString(`],"from":`)
	out.String(em.From)
	out.RawString(`,"subject":`)
	out.String(em.Subject)
	out.RawByte('}')
}

func dumpEvents(tx *store.Tx) {
	for _, e := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpEvent(tx, &out, e)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpEvent(tx *store.Tx, out *jwriter.Writer, e *model.Event) {
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
	if e.Organization != model.OrgNone {
		out.RawString(`,"organization":`)
		out.String(model.OrganizationNames[e.Organization])
	}
	if e.Private {
		out.RawString(`,"private":true`)
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
	out.RawString(`,"attendance":[`)
	var eattend = tx.FetchAttendanceByEvent(e)
	var pids = make([]model.PersonID, 0, len(eattend))
	for p := range eattend {
		pids = append(pids, p)
	}
	sort.Slice(pids, func(i, j int) bool { return pids[i] < pids[j] })
	for i, pid := range pids {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"person":`)
		out.Int(int(pid))
		out.RawString(`,"sortName":`)
		out.String(personName(tx, pid))
		out.RawString(`,"type":`)
		var ai = eattend[pid]
		out.String(model.AttendanceTypeNames[ai.Type])
		out.RawString(`,"minutes":`)
		out.Uint16(ai.Minutes)
		out.RawByte('}')
	}
	out.RawString(`]}`)
}

func dumpFolders(tx *store.Tx) {
	dumpFolder(tx, tx.FetchRootFolder())
}

func dumpFolder(tx *store.Tx, f *model.FolderNode) {
	var out jwriter.Writer
	out.NoEscapeHTML = true
	out.RawString(`{"id":`)
	out.Int(int(f.ID))
	if f.Parent != 0 {
		out.RawString(`,"parent":{"id":`)
		out.Int(int(f.Parent))
		out.RawString(`,"name":`)
		out.String(tx.FetchFolder(f.Parent).Name)
		out.RawByte('}')
	}
	out.RawString(`,"name":`)
	out.String(f.Name)
	if f.Group != 0 {
		out.RawString(`,"group":{"id":`)
		out.Int(int(f.Group))
		out.RawString(`,"name":`)
		out.String(tx.Authorizer().FetchGroup(f.Group).Name)
		out.RawByte('}')
	}
	if f.Approvals != 0 {
		out.RawString(`,"approvals":`)
		out.Int(f.Approvals)
	}
	out.RawString(`,"documents":[`)
	for i, d := range f.Documents {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(d.ID))
		out.RawString(`,"name":`)
		out.String(d.Name)
		if d.PostedBy != 0 {
			out.RawString(`,"postedBy":{"id":`)
			out.Int(int(d.PostedBy))
			out.RawString(`,"sortName":`)
			out.String(tx.FetchPerson(d.PostedBy).SortName)
			out.RawByte('}')
		}
		if !d.PostedAt.IsZero() {
			out.RawString(`,"postedAt":`)
			out.Raw(d.PostedAt.MarshalJSON())
		}
		if d.NeedsApproval {
			out.RawString(`,"needsApproval":true`)
		}
		out.RawByte('}')
	}
	out.RawString("]}\n")
	out.DumpTo(os.Stdout)
	for _, cf := range f.ChildNodes {
		dumpFolder(tx, cf)
	}
}

func dumpGroups(tx *store.Tx) {
	for _, g := range tx.Authorizer().FetchGroups(tx.Authorizer().AllGroups()) {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		g.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpPeople(tx *store.Tx) {
	for _, p := range tx.FetchPeople() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpPerson(tx, &out, p)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpPerson(tx *store.Tx, out *jwriter.Writer, p *model.Person) {
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
	if len(p.Notes) != 0 {
		out.RawString(`,"notes":[`)
		for i, n := range p.Notes {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"note":`)
			out.String(n.Note)
			out.RawString(`,"date":`)
			out.String(n.Date)
			out.RawString(`,"privilege":`)
			out.String(model.PrivilegeNames[n.Privilege])
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	if p.NoEmail {
		out.RawString(`,"noEmail":true`)
	}
	if p.NoText {
		out.RawString(`,"noText":true`)
	}
	out.RawString(`,"unsubscribeToken":`)
	out.String(p.UnsubscribeToken)
	out.RawString(`,"roles":[`)
	for i, r := range tx.Authorizer().FetchRoles(tx.Authorizer().RolesP(p.ID)) {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(r.ID))
		out.RawString(`,"name":`)
		out.String(r.Name)
		out.RawByte('}')
	}
	out.RawString(`],"groups":[`) // not read by serv-load, but useful for filtering
	for i, g := range tx.Authorizer().FetchGroups(tx.Authorizer().GroupsP(p.ID)) {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g.ID))
		out.RawString(`,"name":`)
		out.String(g.Name)
		out.RawByte('}')
	}
	out.RawByte(']')
	if len(p.DSWForms) != 0 {
		out.RawString(`,"dswForms":[`)
		for i, f := range p.DSWForms {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"from":`)
			out.String(f.From.Format("2006-01-02"))
			if !f.To.IsZero() {
				out.RawString(`,"to":`)
				out.String(f.To.Format("2006-01-02"))
			}
			out.RawString(`,"for":`)
			out.String(f.For)
			if f.Invalid != "" {
				out.RawString(`,"invalid":`)
				out.String(f.Invalid)
			}
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	if p.VolgisticsID != 0 {
		out.RawString(`,"volgisticsID":`)
		out.Int(p.VolgisticsID)
	}
	if p.BackgroundCheck != "" {
		out.RawString(`,"backgroundCheck":`)
		out.String(p.BackgroundCheck)
	}
	if p.HoursToken != "" {
		out.RawString(`,"hoursToken":`)
		out.String(p.HoursToken)
	}
	if p.HoursReminder {
		out.RawString(`,"hoursReminder":true`)
	}
	out.RawByte('}')
}

func dumpRoles(tx *store.Tx) {
	for _, r := range tx.Authorizer().FetchRoles(tx.Authorizer().AllRoles()) {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpRole(tx, &out, r)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpRole(tx *store.Tx, out *jwriter.Writer, r *model.Role) {
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
	if r.Detail {
		out.RawString(`,"detail":`)
		out.Bool(r.Detail)
	}
	out.RawString(`,"privileges":[`)
	first := true
	for _, g := range tx.Authorizer().FetchGroups(tx.Authorizer().AllGroups()) {
		privs := tx.Authorizer().ActionsRG(r.ID, g.ID)
		if privs == 0 {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g.ID))
		out.RawString(`,"name":`)
		out.String(g.Name)
		out.RawString(`,"privileges":[`)
		first2 := true
		for _, p := range model.AllPrivileges {
			if privs&p == 0 {
				continue
			}
			if first2 {
				first2 = false
			} else {
				out.RawByte(',')
			}
			out.String(model.PrivilegeNames[p])
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
}

func dumpSessions(tx *store.Tx) {
	for _, s := range tx.FetchSessions() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpSession(tx, &out, s)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpSession(tx *store.Tx, out *jwriter.Writer, s *model.Session) {
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

func dumpTextMessages(tx *store.Tx) {
	for _, t := range tx.FetchTextMessages() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		dumpTextMessage(tx, &out, t)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func dumpTextMessage(tx *store.Tx, out *jwriter.Writer, t *model.TextMessage) {
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

func dumpVenues(tx *store.Tx) {
	for _, v := range tx.FetchVenues() {
		var out jwriter.Writer
		out.NoEscapeHTML = true
		v.MarshalEasyJSON(&out)
		out.DumpTo(os.Stdout)
		os.Stdout.Write([]byte{'\n'})
	}
}

func groupName(tx *store.Tx, id model.GroupID) string {
	if v := tx.Authorizer().FetchGroup(id); v != nil {
		return v.Name
	}
	return ""
}
func personName(tx *store.Tx, id model.PersonID) string {
	if v := tx.FetchPerson(id); v != nil {
		return v.SortName
	}
	return ""
}
func roleName(tx *store.Tx, id model.RoleID) string {
	if v := tx.Authorizer().FetchRole(id); v != nil {
		return v.Name
	}
	return ""
}
func venueName(tx *store.Tx, id model.VenueID) string {
	if v := tx.FetchVenue(id); v != nil {
		return v.Name
	}
	return ""
}
