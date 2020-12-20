package event

import (
	"errors"
	"strconv"

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
