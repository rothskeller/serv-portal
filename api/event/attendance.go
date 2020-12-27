package event

import (
	"errors"
	"strconv"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// PostEventAttendance handles POST /api/events/$id/attendance requests.
func PostEventAttendance(r *util.Request, idstr string) error {
	var (
		event  *model.Event
		person *model.Person
		attend = make(map[model.PersonID]model.AttendanceInfo)
	)
	if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
		return util.NotFound
	}
	if r.Person.Orgs[event.Org].PrivLevel < model.PrivLeader {
		return util.Forbidden
	}
	if attendanceFinalized(event.Date) {
		return util.Forbidden
	}
	r.ParseMultipartForm(1048576)
	for i, idstr := range r.Form["person"] {
		var ai model.AttendanceInfo
		var typestr string
		var err error
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return errors.New("invalid person")
		}
		if len(r.Form["type"]) > i {
			typestr = r.Form["type"][i]
		}
		if ai.Type, err = model.ParseAttendanceType(typestr); err != nil {
			return err
		}
		if len(r.Form["minutes"]) > i {
			if min, err := strconv.Atoi(r.Form["minutes"][i]); err == nil && min >= 0 && min <= 1440 {
				ai.Minutes = uint16(min)
			} else {
				return errors.New("invalid minutes")
			}
		} else {
			return errors.New("invalid minutes")
		}
		attend[person.ID] = ai
	}
	r.Tx.SaveEventAttendance(event, attend)
	r.Tx.Commit()
	return nil
}

// attendanceFinalized returns true if the attendance for the event has been
// finalized (i.e., the hours for it have been submitted to Volgistics).  This
// becomes true at midnight on the start of the 11th of the month after the
// event.
func attendanceFinalized(date string) bool {
	var event, threshold time.Time
	event, _ = time.ParseInLocation("2006-01-02", date, time.Local)
	threshold = time.Date(event.Year(), event.Month()+1, 11, 0, 0, 0, 0, time.Local)
	return threshold.Before(time.Now())
}
