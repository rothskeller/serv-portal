package person

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPersonHours handles GET /api/people/$id/hours requests.
func GetPersonHours(r *util.Request, idstr string) error {
	var (
		person *model.Person
		out    jwriter.Writer
		now    = time.Now()
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if person != r.Person && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if person.VolgisticsID == 0 {
		r.Header().Set("Content-Type", "application/json; charset=utf-8")
		r.Write([]byte(`false`))
		return nil
	}
	out.RawByte('[')
	if now.Day() <= 10 {
		getPersonHours(r, &out, person, time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local))
		out.RawByte(',')
	}
	getPersonHours(r, &out, person, now)
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
func getPersonHours(r *util.Request, out *jwriter.Writer, person *model.Person, month time.Time) {
	var (
		mstr    = month.Format("2006-01")
		first   = true
		pgroups = r.Auth.FetchGroups(r.Auth.GroupsP(person.ID))
	)
	out.RawString(`{"month":`)
	out.String(month.Format("January 2006"))
	out.RawString(`,"events":[`)
	// Since we're just doing a <= comparison on strings, it doesn't matter
	// how many days there are in the month.
	for _, e := range r.Tx.FetchEvents(mstr+"-01", mstr+"-31") {
		// Show this event if the person belongs to the relevant org or
		// if they have hours already recorded for it.
		var (
			amap = r.Tx.FetchAttendanceByEvent(e)
			show = amap[person.ID].Minutes != 0
		)
		if !show && e.Organization != model.OrgNone {
			for _, g := range pgroups {
				if g.Organization == e.Organization {
					show = true
					break
				}
			}
		}
		if !show {
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
		out.RawString(`,"minutes":`)
		out.Uint16(amap[person.ID].Minutes)
		if e.Type == model.EventHours {
			out.RawString(`,"placeholder":true`)
		}
		out.RawByte('}')
	}
	out.RawString(`]}`)
}

// PostPersonHours handles POST /api/people/$id/hours requests.
func PostPersonHours(r *util.Request, idstr string) (err error) {
	var (
		person *model.Person
		now    = time.Now()
	)
	if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
		return util.NotFound
	}
	if person != r.Person && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if person.VolgisticsID == 0 {
		return errors.New("can't report hours for person not registered as volunteer")
	}
	r.ParseMultipartForm(1048576)
	for k, v := range r.Form {
		var (
			eid     int
			event   *model.Event
			minutes int
		)
		if !strings.HasPrefix(k, "e") {
			continue
		}
		if eid, err = strconv.Atoi(k[1:]); err != nil || eid < 1 {
			continue
		}
		if event = r.Tx.FetchEvent(model.EventID(eid)); event == nil {
			return errors.New("nonexistent event")
		}
		if now.Day() <= 10 {
			if event.Date < time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Format("2006-01-02") {
				return errors.New("event is too old")
			}
		} else {
			if event.Date < now.Format("2006-01")+"-01" {
				return errors.New("event is too old")
			}
		}
		if minutes, err = strconv.Atoi(v[0]); err != nil || minutes < 0 {
			return errors.New("invalid minutes")
		}
		emap := r.Tx.FetchAttendanceByEvent(event)
		if att, ok := emap[person.ID]; ok {
			att.Minutes = uint16(minutes)
			emap[person.ID] = att
		} else {
			emap[person.ID] = model.AttendanceInfo{Type: model.AttendAsAbsent, Minutes: uint16(minutes)}
		}
		r.Tx.SaveEventAttendance(event, emap)
	}
	r.Tx.Commit()
	return nil
}
