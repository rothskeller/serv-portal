package event

import (
	"errors"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// PostEventAttendance handles POST /api/events/$id/attendance requests.
func PostEventAttendance(r *util.Request, idstr string) error {
	var (
		event  *model.Event
		person *model.Person
		attend = map[model.PersonID]bool{}
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
		attend[person.ID] = true
	}
	r.Tx.SaveEventAttendance(event, attend)
	r.Tx.Commit()
	return nil
}
