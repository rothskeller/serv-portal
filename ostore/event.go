package store

import (
	"fmt"
	"strings"

	"sunnyvaleserv.org/portal/model"
)

// CreateEvent creates a new event in the database.
func (tx *Tx) CreateEvent(e *model.Event) {
	var rstr []string

	tx.Tx.CreateEvent(e)
	tx.entry.Change("create event [%d]", e.ID)
	tx.entry.Change("set event [%d] name to %q", e.ID, e.Name)
	tx.entry.Change("set event [%d] date to %s", e.ID, e.Date)
	tx.entry.Change("set event [%d] start to %s", e.ID, e.Start)
	tx.entry.Change("set event [%d] end to %s", e.ID, e.End)
	if e.Venue != 0 {
		tx.entry.Change("set event [%d] venue to %q [%d]", e.ID, tx.FetchVenue(e.Venue).Name, e.Venue)
	}
	if e.Details != "" {
		tx.entry.Change("set event [%d] details to %q", e.ID, e.Details)
	}
	tx.entry.Change("set event [%d] type to %s", e.ID, model.EventTypeNames[e.Type])
	if e.RenewsDSW {
		tx.entry.Change("set event [%d] renewsDSW", e.ID)
	}
	if e.CoveredByDSW {
		tx.entry.Change("set event [%d] coveredByDSW", e.ID)
	}
	tx.entry.Change("set event %s %q [%d] org to %s", e.Date, e.Name, e.ID, e.Org.String())
	if len(e.Roles) != 0 {
		for _, r := range e.Roles {
			rstr = append(rstr, fmt.Sprintf("%q [%d]", tx.FetchRole(r).Name, r))
		}
		tx.entry.Change("set event [%d] roles to %s", e.ID, strings.Join(rstr, ", "))
	}
	for _, s := range e.Shifts {
		tx.entry.Change("add event [%d] shift %s-%s task %q min %d max %d announce %v", e.ID, s.Start, s.End, s.Task, s.Min, s.Max, s.Announce)
		for _, p := range s.SignedUp {
			tx.entry.Change("add event [%d] shift %s-%s task %q signedUp person %q [%d]", e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
		}
		for _, p := range s.Declined {
			tx.entry.Change("add event [%d] shift %s-%s task %q declined person %q [%d]", e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
		}
	}
	if e.SignupText != "" {
		tx.entry.Change("set event [%d] signupText to %q", e.ID, e.SignupText)
	}
}

// UpdateEvent updates an existing event in the database.
func (tx *Tx) UpdateEvent(e *model.Event) {
	var oe *model.Event

	oe = tx.Tx.FetchEvent(e.ID)
	tx.Tx.UpdateEvent(e)
	if e.Name != oe.Name {
		tx.entry.Change("set event %s %q [%d] name to %q", e.Date, e.Name, e.ID, e.Name)
	}
	if e.Date != oe.Date {
		tx.entry.Change("set event %s %q [%d] date to %s", e.Date, e.Name, e.ID, e.Date)
	}
	if e.Start != oe.Start {
		tx.entry.Change("set event %s %q [%d] start to %s", e.Date, e.Name, e.ID, e.Start)
	}
	if e.End != oe.End {
		tx.entry.Change("set event %s %q [%d] end to %s", e.Date, e.Name, e.ID, e.End)
	}
	if e.Venue != oe.Venue {
		if e.Venue != 0 {
			tx.entry.Change("set event %s %q [%d] venue to %q [%d]", e.Date, e.Name, e.ID, tx.FetchVenue(e.Venue).Name, e.Venue)
		} else {
			tx.entry.Change("remove event %s %q [%d] venue", e.Date, e.Name, e.ID)
		}
	}
	if e.Details != oe.Details {
		tx.entry.Change("set event %s %q [%d] details to %q", e.Date, e.Name, e.ID, e.Details)
	}
	if e.Type != oe.Type {
		tx.entry.Change("set event %s %q [%d] type to %s", e.Date, e.Name, e.ID, model.EventTypeNames[e.Type])
	}
	if e.RenewsDSW != oe.RenewsDSW {
		if e.RenewsDSW {
			tx.entry.Change("set event %s %q [%d] renewsDSW", e.Date, e.Name, e.ID)
		} else {
			tx.entry.Change("clear event %s %q [%d] renewsDSW", e.Date, e.Name, e.ID)
		}
	}
	if e.CoveredByDSW != oe.CoveredByDSW {
		if e.CoveredByDSW {
			tx.entry.Change("set event %s %q [%d] coveredByDSW", e.Date, e.Name, e.ID)
		} else {
			tx.entry.Change("clear event %s %q [%d] coveredByDSW", e.Date, e.Name, e.ID)
		}
	}
	if e.Org != oe.Org {
		tx.entry.Change("set event %s %q [%d] org to %s", e.Date, e.Name, e.ID, e.Org.String())
	}
ROLES1:
	for _, or := range oe.Roles {
		for _, r := range e.Roles {
			if or == r {
				continue ROLES1
			}
		}
		tx.entry.Change("remove event %s %q [%d] role %q [%d]", e.Date, e.Name, e.ID, tx.FetchRole(or).Name, or)
	}
ROLES2:
	for _, r := range e.Roles {
		for _, or := range oe.Roles {
			if or == r {
				continue ROLES2
			}
		}
		tx.entry.Change("add event %s %q [%d] role %q [%d]", e.Date, e.Name, e.ID, tx.FetchRole(r).Name, r)
	}
	if oe.SignupText != e.SignupText {
		if e.SignupText != "" {
			tx.entry.Change("set event %s %q [%d] signupText to %q", e.Date, e.Name, e.ID, e.SignupText)
		} else {
			tx.entry.Change("clear event %s %q [%d] signupText", e.Date, e.Name, e.ID)
		}
	}
SHIFT1:
	for _, os := range oe.Shifts {
		for _, s := range e.Shifts {
			if os.Start == s.Start && os.Task == s.Task {
				tx.updateShift(e, s, os)
				continue SHIFT1
			}
		}
		tx.entry.Change("remove event %s %q [%d] shift %s-%s task %q", e.Date, e.Name, e.ID, os.Start, os.End, os.Task)
	}
SHIFT2:
	for _, s := range e.Shifts {
		for _, os := range oe.Shifts {
			if os.Start == s.Start && os.Task == s.Task {
				continue SHIFT2
			}
		}
		tx.entry.Change("add event %s %q [%d] shift %s-%s task %q min %d max %d announce %v", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, s.Min, s.Max, s.Announce)
		for _, p := range s.SignedUp {
			tx.entry.Change("add event %s %q [%d] shift %s-%s task %q signedUp person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
		}
		for _, p := range s.Declined {
			tx.entry.Change("add event %s %q [%d] shift %s-%s task %q declined person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
		}
	}
}
func (tx *Tx) updateShift(e *model.Event, s, os *model.Shift) {
	if s.End != os.End {
		tx.entry.Change("set event %s %q [%d] shift at %s task %q end to %s", e.Date, e.Name, e.ID, s.Start, s.Task, s.End)
	}
	if s.Min != os.Min {
		tx.entry.Change("set event %s %q [%d] shift %s-%s task %q min to %d", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, s.Min)
	}
	if s.Max != os.Max {
		tx.entry.Change("set event %s %q [%d] shift %s-%s task %q max to %d", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, s.Max)
	}
	if s.Announce != os.Announce {
		if s.Announce {
			tx.entry.Change("set event %s %q [%d] shift %s-%s task %q announce", e.Date, e.Name, e.ID, s.Start, s.End, s.Task)
		} else {
			tx.entry.Change("clear event %s %q [%d] shift %s-%s task %q announce", e.Date, e.Name, e.ID, s.Start, s.End, s.Task)
		}
	}
SIGNED1:
	for _, op := range os.SignedUp {
		for _, p := range s.SignedUp {
			if op == p {
				continue SIGNED1
			}
		}
		tx.entry.Change("remove event %s %q [%d] shift %s-%s task %q signedUp person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(op).InformalName, op)
	}
SIGNED2:
	for _, p := range s.SignedUp {
		for _, op := range os.SignedUp {
			if op == p {
				continue SIGNED2
			}
		}
		tx.entry.Change("add event %s %q [%d] shift %s-%s task %q signedUp person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
	}
DECLINE1:
	for _, op := range os.Declined {
		for _, p := range s.Declined {
			if op == p {
				continue DECLINE1
			}
		}
		tx.entry.Change("remove event %s %q [%d] shift %s-%s task %q declined person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(op).InformalName, op)
	}
DECLINE2:
	for _, p := range s.Declined {
		for _, op := range os.Declined {
			if op == p {
				continue DECLINE2
			}
		}
		tx.entry.Change("add event %s %q [%d] shift %s-%s task %q declined person %q [%d]", e.Date, e.Name, e.ID, s.Start, s.End, s.Task, tx.FetchPerson(p).InformalName, p)
	}
}

// DeleteEvent deletes an event from the database.
func (tx *Tx) DeleteEvent(e *model.Event) {
	tx.Tx.DeleteEvent(e)
	tx.entry.Change("delete event %s %q [%d]", e.Date, e.Name, e.ID)
}

// SaveEventAttendance saves the attendance for a specific event.
func (tx *Tx) SaveEventAttendance(e *model.Event, attend map[model.PersonID]model.AttendanceInfo) {
	var oattend = tx.Tx.FetchAttendanceByEvent(e)
	tx.Tx.SaveEventAttendance(e, attend)
	for pid, ai := range attend {
		oai := oattend[pid]
		if ai.Minutes != oai.Minutes || ai.Type != oai.Type {
			tx.entry.Change("set event %s %q [%d] person %q [%d] attendance to %s %d min", e.Date, e.Name, e.ID, tx.FetchPerson(pid).InformalName, pid, ai.Type, ai.Minutes)
		}
	}
	for pid := range oattend {
		if _, ok := attend[pid]; !ok {
			tx.entry.Change("remove event %s %q [%d] person %q [%d] attendance", e.Date, e.Name, e.ID, tx.FetchPerson(pid).InformalName, pid)
		}
	}
}
