package store

import (
	"fmt"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/model"
)

// CreateEvent creates a new event in the database.
func (tx *Tx) CreateEvent(e *model.Event) {
	var gstr []string

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
	tx.entry.Change("set event %s %q [%d] organization to %s", e.Date, e.Name, e.ID, model.OrganizationNames[e.Organization])
	if e.Private {
		tx.entry.Change("set event %s %q [%d] private flag", e.Date, e.Name, e.ID)
	}
	tx.entry.Change("set event [%d] type to %s", e.ID, model.EventTypeNames[e.Type])
	if len(e.Groups) != 0 {
		for _, g := range e.Groups {
			gstr = append(gstr, fmt.Sprintf("%q [%d]", tx.Authorizer().FetchGroup(g).Name, g))
		}
		tx.entry.Change("set event [%d] groups to %s", e.ID, strings.Join(gstr, ", "))
	}
	if e.RenewsDSW {
		tx.entry.Change("set event [%d] renewsDSW", e.ID)
	}
	if e.CoveredByDSW {
		tx.entry.Change("set event [%d] coveredByDSW", e.ID)
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
	if e.Type != oe.Type {
		tx.entry.Change("set event %s %q [%d] type to %s", e.Date, e.Name, e.ID, model.EventTypeNames[e.Type])
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
	if e.RenewsDSW != oe.RenewsDSW {
		if e.RenewsDSW {
			tx.entry.Change("set event %s %q [%d] renewsDSW", e.Date, e.Name, e.ID)
		} else {
			tx.entry.Change("clear event %s %q [%d] renewsDSW", e.Date, e.Name, e.ID)
		}
		tx.recalculateDSWUntil(model.OrganizationToDSWClass[e.Organization], tx.Tx.FetchAttendanceByEvent(e), nil, nil)
	}
	if e.CoveredByDSW != oe.CoveredByDSW {
		if e.CoveredByDSW {
			tx.entry.Change("set event %s %q [%d] coveredByDSW", e.Date, e.Name, e.ID)
		} else {
			tx.entry.Change("clear event %s %q [%d] coveredByDSW", e.Date, e.Name, e.ID)
		}
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
	tx.recalculateDSWUntil(model.OrganizationToDSWClass[e.Organization], attend, oattend, nil)
}

// recalculateDSWUntil recalculates the DSWUntil values, for the specified DSW
// classification, for all people listed in either of the two provided maps or
// the single provided pointer.
func (tx *Tx) recalculateDSWUntil(class model.DSWClass, a1, a2 map[model.PersonID]model.AttendanceInfo, one *model.Person) {
	var (
		oldest time.Time
		people = make(map[model.PersonID]*model.Person)
	)
	if a1 != nil {
		for p := range a1 {
			people[p] = nil
		}
	}
	if a2 != nil {
		for p := range a2 {
			people[p] = nil
		}
	}
	if one != nil {
		people[one.ID] = one
	}
	for pid := range people {
		if people[pid] == nil {
			people[pid] = tx.FetchPerson(pid)
			tx.WillUpdatePerson(people[pid])
		}
		if people[pid].DSWRegistrations == nil || people[pid].DSWRegistrations[class].IsZero() {
			if people[pid].DSWUntil != nil && !people[pid].DSWUntil[class].IsZero() {
				delete(people[pid].DSWUntil, class)
				tx.UpdatePerson(people[pid])
			}
			delete(people, pid)
			continue
		}
		if people[pid].DSWUntil == nil {
			people[pid].DSWUntil = make(map[model.DSWClass]time.Time)
		}
		people[pid].DSWUntil[class] = people[pid].DSWRegistrations[class].AddDate(1, 0, 0)
		if oldest.IsZero() || people[pid].DSWRegistrations[class].Before(oldest) {
			oldest = people[pid].DSWRegistrations[class]
		}
	}
	if oldest.IsZero() || len(people) == 0 {
		return
	}
	for _, e := range tx.FetchEvents(oldest.Format("2006-01-02"), time.Now().Format("2006-01-02")) {
		if !e.RenewsDSW || model.OrganizationToDSWClass[e.Organization] != class {
			continue
		}
		edate, _ := time.ParseInLocation("2006-01-02", e.Date, time.Local)
		extendTo := time.Date(edate.Year()+1, edate.Month(), edate.Day(), 0, 0, 0, 0, time.Local)
		eatt := tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt {
			if ai.Minutes == 0 {
				continue
			}
			if _, ok := people[pid]; !ok {
				continue
			}
			if !people[pid].DSWUntil[class].Before(edate) && people[pid].DSWUntil[class].Before(extendTo) {
				people[pid].DSWUntil[class] = extendTo
			}
		}
	}
	for _, p := range people {
		if p != one {
			tx.UpdatePerson(p)
		}
	}
}
