package store

import (
	"fmt"
	"strings"

	"sunnyvaleserv.org/portal/model"
)

// CreateEvent creates a new event in the database.
func (tx *Tx) CreateEvent(e *model.Event) {
	var (
		etstr []string
		gstr  []string
	)
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
	if e.Organization != model.OrgNone {
		tx.entry.Change("set event %s %q [%d] organization to %s", e.Date, e.Name, e.ID, model.OrganizationNames[e.Organization])
	}
	if e.Private {
		tx.entry.Change("set event %s %q [%d] private flag", e.Date, e.Name, e.ID)
	}
	for _, et := range model.AllEventTypes {
		if e.Type&et != 0 {
			etstr = append(etstr, model.EventTypeNames[et])
		}
	}
	if len(etstr) == 1 {
		tx.entry.Change("set event [%d] type to %s", e.ID, etstr[0])
	} else if len(etstr) > 1 {
		tx.entry.Change("set event [%d] types to %s", e.ID, strings.Join(etstr, ", "))
	}
	if len(e.Groups) != 0 {
		for _, g := range e.Groups {
			gstr = append(gstr, fmt.Sprintf("%q [%d]", tx.Authorizer().FetchGroup(g).Name, g))
		}
		tx.entry.Change("set event [%d] groups to %s", e.ID, strings.Join(gstr, ", "))
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
	if e.Organization != oe.Organization {
		tx.entry.Change("set event %s %q [%d] organization to %s", e.Date, e.Name, e.ID, model.OrganizationNames[e.Organization])
	}
	if e.Private != oe.Private {
		if e.Private {
			tx.entry.Change("set event %s %q [%d] private flag", e.Date, e.Name, e.ID)
		} else {
			tx.entry.Change("clear event %s %q [%d] private flag", e.Date, e.Name, e.ID)
		}
	}
	for _, et := range model.AllEventTypes {
		if e.Type&et != oe.Type&et {
			if e.Type&et != 0 {
				tx.entry.Change("add event %s %q [%d] type %s", e.Date, e.Name, e.ID, model.EventTypeNames[et])
			} else {
				tx.entry.Change("remove event %s %q [%d] type %s", e.Date, e.Name, e.ID, model.EventTypeNames[et])
			}
		}
	}
GROUPS1:
	for _, og := range oe.Groups {
		for _, g := range e.Groups {
			if og == g {
				continue GROUPS1
			}
		}
		tx.entry.Change("remove event %s %q [%d] group %q [%d]", e.Date, e.Name, e.ID, tx.Authorizer().FetchGroup(og).Name, og)
	}
GROUPS2:
	for _, g := range e.Groups {
		for _, og := range oe.Groups {
			if og == g {
				continue GROUPS2
			}
		}
		tx.entry.Change("add event %s %q [%d] group %q [%d]", e.Date, e.Name, e.ID, tx.Authorizer().FetchGroup(g).Name, g)
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
			tx.entry.Change("set event %s %q [%d] person %q [%d] attendance to %s %d min", e.Date, e.Name, e.ID, tx.FetchPerson(pid).InformalName, pid, model.AttendanceTypeNames[ai.Type], ai.Minutes)
		}
	}
	for pid := range oattend {
		if _, ok := attend[pid]; !ok {
			tx.entry.Change("remove event %s %q [%d] person %q [%d] attendance", e.Date, e.Name, e.ID, tx.FetchPerson(pid).InformalName, pid)
		}
	}
}
