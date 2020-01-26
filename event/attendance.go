package event

import (
	"errors"
	"strconv"

	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// PostEventAttendance handles POST /api/events/$id/attendance requests.
func PostEventAttendance(r *util.Request, idstr string) error {
	var (
		event  *model.Event
		person *model.Person
		attend = map[model.PersonID]model.AttendanceInfo{}
	)
	if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
		return util.NotFound
	}
	if !auth.CanRecordAttendanceAtEvent(r, event) {
		return util.Forbidden
	}
	attend = r.Tx.FetchAttendanceByEvent(event)
	for pid := range attend {
		if auth.CanViewPerson(r, r.Tx.FetchPerson(pid)) {
			delete(attend, pid)
		}
	}
	r.ParseMultipartForm(1048576)
	for i, idstr := range r.Form["person"] {
		var ai model.AttendanceInfo
		var typestr string
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return errors.New("invalid person")
		}
		if !auth.CanViewPerson(r, person) {
			return errors.New("invalid person")
		}
		if len(r.Form["type"]) > i {
			typestr = r.Form["type"][i]
		}
		switch typestr {
		case model.AttendanceTypeNames[model.AttendAsVolunteer]:
			ai.Type = model.AttendAsVolunteer
		case model.AttendanceTypeNames[model.AttendAsStudent]:
			ai.Type = model.AttendAsStudent
		case model.AttendanceTypeNames[model.AttendAsAuditor]:
			ai.Type = model.AttendAsAuditor
		default:
			return errors.New("invalid type")
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
