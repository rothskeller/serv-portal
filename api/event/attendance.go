package event

import (
	"errors"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetEventAttendance handles GET /api/events/$id/attendance requests.
func GetEventAttendance(r *util.Request, idstr string) error {
	var (
		event    *model.Event
		attended map[model.PersonID]bool
		out      jwriter.Writer
	)
	if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
		return util.NotFound
	}
	if !auth.CanRecordAttendanceAtEvent(r, event) {
		return util.Forbidden
	}
	attended = r.Tx.FetchAttendanceByEvent(event)
	out.RawString(`{"event":{"id":`)
	out.Int(int(event.ID))
	out.RawString(`,"date":`)
	out.String(event.Date)
	out.RawString(`,"name":`)
	out.String(event.Name)
	out.RawString(`},"people":[`)
	first := true
	for _, p := range r.Tx.FetchPeople() {
		if !auth.CanViewEventP(r, p, event) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(p.ID))
		out.RawString(`,"lastName":`)
		out.String(p.LastName)
		out.RawString(`,"firstName":`)
		out.String(p.FirstName)
		out.RawString(`,"attended":`)
		out.Bool(attended[p.ID])
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostEventAttendance handles POST /api/events/$id/attendance requests.
func PostEventAttendance(r *util.Request, idstr string) error {
	var (
		event  *model.Event
		person *model.Person
		people []*model.Person
	)
	if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
		return util.NotFound
	}
	if !auth.CanRecordAttendanceAtEvent(r, event) {
		return util.Forbidden
	}
	r.ParseMultipartForm(1048576)
	for _, idstr := range r.Form["person"] {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return errors.New("invalid person")
		}
		if !auth.CanViewEventP(r, person, event) {
			return errors.New("illegal person")
		}
		people = append(people, person)
	}
	r.Tx.SaveEventAttendance(event, people)
	r.Tx.Commit()
	return nil
}
