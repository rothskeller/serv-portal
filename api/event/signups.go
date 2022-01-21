package event

import (
	"errors"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetEventSignups handles GET /api/events/signups and GET
// /api/events/signups/${token} requests.
func GetEventSignups(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
		first  = true
	)
	if idstr != "" {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	} else {
		if person = r.Person; person == nil {
			return util.Forbidden
		}
	}
	out.RawString(`{"id":`)
	out.Int(int(person.ID))
	out.RawString(`,"sortName":`)
	out.String(person.SortName)
	out.RawString(`,"events":[`)
	for _, e := range r.Tx.FetchEvents(time.Now().Format("2006-01-02"), "2099-12-31") {
		if len(e.Shifts) == 0 {
			continue
		}
		var canSignUp = person.Orgs[e.Org].PrivLevel >= model.PrivLeader
		if !canSignUp {
			for _, r := range e.Roles {
				if _, ok := person.Roles[r]; ok {
					canSignUp = true
					break
				}
			}
		}
		if !canSignUp {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(e.ID))
		out.RawString(`,"date":`)
		out.String(e.Date)
		out.RawString(`,"name":`)
		out.String(e.Name)
		out.RawString(`,"signupText":`)
		out.String(e.SignupText)
		out.RawString(`,"shifts":[`)
		for i, s := range e.Shifts {
			var signedUp bool
			var first = true
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"start":`)
			out.String(s.Start)
			out.RawString(`,"end":`)
			out.String(s.End)
			out.RawString(`,"task":`)
			out.String(s.Task)
			out.RawString(`,"min":`)
			out.Int(s.Min)
			out.RawString(`,"max":`)
			out.Int(s.Max)
			out.RawString(`,"count":`)
			out.Int(len(s.SignedUp))
			out.RawString(`,"list":[`)
			for _, p := range s.SignedUp {
				if p == person.ID {
					signedUp = true
					continue
				} else {
					if first {
						first = false
					} else {
						out.RawByte(',')
					}
					out.RawString(`{"id":`)
					out.Int(int(p))
					out.RawString(`,"sortName":`)
					out.String(r.Tx.FetchPerson(p).SortName)
					out.RawByte('}')
				}
			}
			out.RawString(`],"signedUp":`)
			out.Bool(signedUp)
			out.RawByte('}')
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostEventSignups handles POST /api/events/signups and POST
// /api/events/signups/${token} requests.
func PostEventSignups(r *util.Request, idstr string) error {
	var (
		person *model.Person
	)
	if idstr != "" {
		if person = r.Tx.FetchPersonByUnsubscribe(idstr); person == nil {
			return util.NotFound
		}
	} else {
		if person = r.Person; person == nil {
			return util.Forbidden
		}
	}
	r.ParseMultipartForm(1048576)
	for key, values := range r.Form {
		var (
			parts    []string
			event    *model.Event
			shift    *model.Shift
			signedUp bool
			declined bool
			found    bool
		)
		if len(values) != 1 {
			return errors.New("duplicate form parameter")
		}
		if parts = strings.Split(key, "."); len(parts) != 3 {
			return errors.New("invalid form parameter name")
		}
		switch values[0] {
		case "declined":
			declined = true
		case "true":
			signedUp = true
		case "false":
			break
		default:
			return errors.New("invalid form parameter value")
		}
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(parts[0]))); event == nil {
			return errors.New("nonexistent event")
		}
		if event.Date < time.Now().Format("2006-01-02") {
			return errors.New("event is in the past")
		}
		if person.Orgs[event.Org].PrivLevel < model.PrivLeader {
			for _, role := range event.Roles {
				if _, ok := person.Roles[role]; ok {
					found = true
					break
				}
			}
			if !found {
				return errors.New("not invited to event")
			}
		}
		for _, s := range event.Shifts {
			if s.Start == parts[1] && s.Task == parts[2] {
				shift = s
				break
			}
		}
		if shift == nil {
			return errors.New("event has no such shift")
		}
		found = false
		for i, p := range shift.SignedUp {
			if p == person.ID {
				found = true
				if !signedUp {
					shift.SignedUp = append(shift.SignedUp[:i], shift.SignedUp[i+1:]...)
				}
				break
			}
		}
		if !found && signedUp {
			shift.SignedUp = append(shift.SignedUp, person.ID)
		}
		found = false
		for i, p := range shift.Declined {
			if p == person.ID {
				found = true
				if !declined {
					shift.Declined = append(shift.Declined[:i], shift.Declined[i+1:]...)
				}
				break
			}
		}
		if !found && declined {
			shift.Declined = append(shift.Declined, person.ID)
		}
		r.Tx.UpdateEvent(event)
	}
	r.Tx.Commit()
	return nil
}
