package store

import (
	"bytes"
	"strconv"
	"time"

	"sunnyvaleserv.org/portal/model"
)

// CreatePerson creates a new person in the database.
func (tx *Tx) CreatePerson(p *model.Person) {
	tx.Tx.CreatePerson(p)
	if tx.auth != nil {
		tx.auth.AddPerson(p.ID)
	}
	tx.entry.Change("create person [%d]", p.ID)
	if p.Username != "" {
		tx.entry.Change("set person [%d] username to %q", p.ID, p.Username)
	}
	tx.entry.Change("set person [%d] informalName to %q", p.ID, p.InformalName)
	tx.entry.Change("set person [%d] formalName to %q", p.ID, p.FormalName)
	tx.entry.Change("set person [%d] sortName to %q", p.ID, p.SortName)
	if p.CallSign != "" {
		tx.entry.Change("set person [%d] callSign to %q", p.ID, p.CallSign)
	}
	if p.Email != "" {
		tx.entry.Change("set person [%d] email to %q", p.ID, p.Email)
	}
	if p.Email2 != "" {
		tx.entry.Change("set person [%d] email2 to %q", p.ID, p.Email2)
	}
	if p.NoEmail {
		tx.entry.Change("set person [%d] noEmail flag", p.ID)
	}
	if p.HomeAddress.Address != "" {
		if p.HomeAddress.FireDistrict != 0 {
			tx.entry.Change("set person [%d] homeAddress to %q (%f, %f) (district %d)", p.ID, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude, p.HomeAddress.FireDistrict)
		} else if p.HomeAddress.Latitude != 0 {
			tx.entry.Change("set person [%d] homeAddress to %q (%f, %f)", p.ID, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude)
		} else {
			tx.entry.Change("set person [%d] homeAddress to %q", p.ID, p.HomeAddress.Address)
		}
	}
	if p.WorkAddress.Address != "" {
		if p.WorkAddress.FireDistrict != 0 {
			tx.entry.Change("set person [%d] workAddress to %q (%f, %f) (district %d)", p.ID, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude, p.WorkAddress.FireDistrict)
		} else if p.WorkAddress.Latitude != 0 {
			tx.entry.Change("set person [%d] workAddress to %q (%f, %f)", p.ID, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude)
		} else {
			tx.entry.Change("set person [%d] workAddress to %q", p.ID, p.WorkAddress.Address)
		}
	} else if p.WorkAddress.SameAsHome {
		tx.entry.Change("set person [%d] workAddress sameAsHome flag", p.ID)
	}
	if p.MailAddress.Address != "" {
		tx.entry.Change("set person [%d] mailAddress to %q", p.ID, p.MailAddress.Address)
	} else if p.MailAddress.SameAsHome {
		tx.entry.Change("set person [%d] mailAddress sameAsHome flag", p.ID)
	}
	if p.CellPhone != "" {
		tx.entry.Change("set person [%d] cellPhone to %q", p.ID, p.CellPhone)
	}
	if p.HomePhone != "" {
		tx.entry.Change("set person [%d] homePhone to %q", p.ID, p.HomePhone)
	}
	if p.WorkPhone != "" {
		tx.entry.Change("set person [%d] workPhone to %q", p.ID, p.WorkPhone)
	}
	if p.NoText {
		tx.entry.Change("set person [%d] noText flag", p.ID)
	}
	if len(p.Password) != 0 {
		tx.entry.Change("set person [%d] password", p.ID)
	}
	for _, n := range p.Notes {
		tx.entry.Change("add person [%d] note %q at %s with privilege %s", p.ID, n.Note, n.Date, model.PrivilegeNames[n.Privilege])
	}
	tx.entry.Change("set person [%d] unsubscribeToken to %s", p.ID, p.UnsubscribeToken)
	for _, f := range p.DSWForms {
		var invalid string
		if f.Invalid != "" {
			invalid = " INVALID " + strconv.Quote(f.Invalid)
		}
		tx.entry.Change("add person [%d] dsw from %s to %s for %q%s", f.From.Format("2006-01-02"), f.To.Format("2006-01-02"), f.For, invalid)
	}
	if p.VolgisticsID != 0 {
		tx.entry.Change("set person [%d] volgisticsID to %d", p.ID, p.VolgisticsID)
	}
	if p.BackgroundCheck != "" {
		tx.entry.Change("set person [%d] backgroundCheck to %s", p.ID, p.BackgroundCheck)
	}
	if p.DSWRegistrations != nil {
		for c, r := range p.DSWRegistrations {
			tx.entry.Change("set person [%d] dswRegistration[%s] to %s", p.ID, model.DSWClassNames[c], r.Format("2006-01-02"))
		}
	}
	if p.DSWUntil != nil {
		for c, r := range p.DSWUntil {
			tx.entry.Change("set person [%d] dswUntil[%s] to %s", p.ID, model.DSWClassNames[c], r.Format("2006-01-02"))
		}
	}
}

// WillUpdatePerson saves a copy of a person before it's updated, so that we can
// compare against it to generate audit log entries.
func (tx *Tx) WillUpdatePerson(p *model.Person) {
	if tx.originalPeople[p.ID] != nil {
		return
	}
	var op = *p
	if p.Notes != nil {
		op.Notes = make([]*model.PersonNote, len(p.Notes))
		for i := range p.Notes {
			opn := *p.Notes[i]
			op.Notes[i] = &opn
		}
	}
	if p.DSWForms != nil {
		forms := make([]model.DSWForm, len(p.DSWForms))
		op.DSWForms = make([]*model.DSWForm, len(p.DSWForms))
		for i := range p.DSWForms {
			forms[i] = *p.DSWForms[i]
			op.DSWForms[i] = &forms[i]
		}
	}
	if p.DSWRegistrations != nil {
		op.DSWRegistrations = make(map[model.DSWClass]time.Time, len(p.DSWRegistrations))
		for c, r := range p.DSWRegistrations {
			op.DSWRegistrations[c] = r
		}
	}
	if p.DSWUntil != nil {
		op.DSWUntil = make(map[model.DSWClass]time.Time, len(p.DSWUntil))
		for c, r := range p.DSWUntil {
			op.DSWUntil[c] = r
		}
	}
	tx.originalPeople[p.ID] = &op
}

// UpdatePerson updates a person in the database.
func (tx *Tx) UpdatePerson(p *model.Person) {
	var op = tx.originalPeople[p.ID]

	if op == nil {
		panic("must call WillUpdatePerson before UpdatePerson")
	}
	tx.Tx.UpdatePerson(p)
	if p.Username != op.Username {
		tx.entry.Change("set person %q [%d] username to %q", p.InformalName, p.ID, p.Username)
	}
	if p.InformalName != op.InformalName {
		tx.entry.Change("set person %q [%d] informalName to %q", p.InformalName, p.ID, p.InformalName)
	}
	if p.FormalName != op.FormalName {
		tx.entry.Change("set person %q [%d] formalName to %q", p.InformalName, p.ID, p.FormalName)
	}
	if p.SortName != op.SortName {
		tx.entry.Change("set person %q [%d] sortName to %q", p.InformalName, p.ID, p.SortName)
	}
	if p.CallSign != op.CallSign {
		tx.entry.Change("set person %q [%d] callSign to %q", p.InformalName, p.ID, p.CallSign)
	}
	if p.Email != op.Email {
		tx.entry.Change("set person %q [%d] email to %q", p.InformalName, p.ID, p.Email)
	}
	if p.Email2 != op.Email2 {
		tx.entry.Change("set person %q [%d] email2 to %q", p.InformalName, p.ID, p.Email2)
	}
	if p.NoEmail != op.NoEmail {
		if p.NoEmail {
			tx.entry.Change("set person %q [%d] noEmail flag", p.InformalName, p.ID)
		} else {
			tx.entry.Change("clear person %q [%d] noEmail flag", p.InformalName, p.ID)
		}
	}
	if p.HomeAddress.Address != op.HomeAddress.Address {
		if p.HomeAddress.FireDistrict != 0 {
			tx.entry.Change("set person %q [%d] homeAddress to %q (%f, %f) (district %d)", p.InformalName, p.ID, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude, p.HomeAddress.FireDistrict)
		} else if p.HomeAddress.Latitude != 0 {
			tx.entry.Change("set person %q [%d] homeAddress to %q (%f, %f)", p.InformalName, p.ID, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude)
		} else {
			tx.entry.Change("set person %q [%d] homeAddress to %q", p.InformalName, p.ID, p.HomeAddress.Address)
		}
	}
	if p.WorkAddress.SameAsHome {
		if !op.WorkAddress.SameAsHome {
			tx.entry.Change("set person %q [%d] workAddress sameAsHome flag", p.InformalName, p.ID)
		}
	} else if op.WorkAddress.SameAsHome && p.WorkAddress.Address == "" {
		tx.entry.Change("clear person %q [%d] workAddress sameAsHome flag", p.InformalName, p.ID)
	} else if p.WorkAddress.Address != op.WorkAddress.Address {
		if p.WorkAddress.FireDistrict != 0 {
			tx.entry.Change("set person %q [%d] workAddress to %q (%f, %f) (district %d)", p.InformalName, p.ID, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude, p.WorkAddress.FireDistrict)
		} else if p.WorkAddress.Latitude != 0 {
			tx.entry.Change("set person %q [%d] workAddress to %q (%f, %f)", p.InformalName, p.ID, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude)
		} else {
			tx.entry.Change("set person %q [%d] workAddress to %q", p.InformalName, p.ID, p.WorkAddress.Address)
		}
	}
	if p.MailAddress.SameAsHome {
		if !op.MailAddress.SameAsHome {
			tx.entry.Change("set person %q [%d] mailAddress sameAsHome flag", p.InformalName, p.ID)
		}
	} else if op.MailAddress.SameAsHome && p.MailAddress.Address == "" {
		tx.entry.Change("clear person %q [%d] mailAddress sameAsHome flag", p.InformalName, p.ID)
	} else if p.MailAddress.Address != op.MailAddress.Address {
		tx.entry.Change("set person %q [%d] mailAddress to %q", p.InformalName, p.ID, p.MailAddress.Address)
	}
	if p.CellPhone != op.CellPhone {
		tx.entry.Change("set person %q [%d] cellPhone to %q", p.InformalName, p.ID, p.CellPhone)
	}
	if p.HomePhone != op.HomePhone {
		tx.entry.Change("set person %q [%d] homePhone to %q", p.InformalName, p.ID, p.HomePhone)
	}
	if p.WorkPhone != op.WorkPhone {
		tx.entry.Change("set person %q [%d] workPhone to %q", p.InformalName, p.ID, p.WorkPhone)
	}
	if p.NoText != op.NoText {
		if p.NoText {
			tx.entry.Change("set person %q [%d] noText flag", p.InformalName, p.ID)
		} else {
			tx.entry.Change("clear person %q [%d] noText flag", p.InformalName, p.ID)
		}
	}
	if len(op.Password) == 0 && len(p.Password) != 0 {
		tx.entry.Change("set person %q [%d] password", p.InformalName, p.ID)
	} else if !bytes.Equal(p.Password, op.Password) {
		tx.entry.Change("change person %q [%d] password", p.InformalName, p.ID)
	}
NOTES1:
	for _, on := range op.Notes {
		for _, n := range p.Notes {
			if n.Date == on.Date && n.Note == on.Note && n.Privilege == on.Privilege {
				continue NOTES1
			}
		}
		tx.entry.Change("remove person %q [%d] note %q at %s with privilege %s", p.InformalName, p.ID, on.Note, on.Date, model.PrivilegeNames[on.Privilege])
	}
NOTES2:
	for _, n := range p.Notes {
		for _, on := range op.Notes {
			if n.Date == on.Date && n.Note == on.Note && n.Privilege == on.Privilege {
				continue NOTES2
			}
		}
		tx.entry.Change("add person %q [%d] note %q at %s with privilege %s", p.InformalName, p.ID, n.Note, n.Date, model.PrivilegeNames[n.Privilege])
	}
	if p.UnsubscribeToken != op.UnsubscribeToken {
		tx.entry.Change("change person %q [%d] unsubscribeToken to %s", p.InformalName, p.ID, p.UnsubscribeToken)
	}
DSW1:
	for _, of := range op.DSWForms {
		for _, f := range p.DSWForms {
			if of.From.Equal(f.From) {
				continue DSW1
			}
		}
		var invalid string
		if of.Invalid != "" {
			invalid = " INVALID " + strconv.Quote(of.Invalid)
		}
		tx.entry.Change("remove person [%d] dsw from %s to %s for %q%s", of.From.Format("2006-01-02"), of.To.Format("2006-01-02"), of.For, invalid)
	}
DSW2:
	for _, f := range p.DSWForms {
		for _, of := range op.DSWForms {
			if of.From.Equal(f.From) {
				if !of.To.Equal(f.To) || of.For != f.For || of.Invalid != f.Invalid {

					var invalid string
					if f.Invalid != "" {
						invalid = " INVALID " + strconv.Quote(f.Invalid)
					}
					tx.entry.Change("change person [%d] dsw from %s to %s for %q%s", f.From.Format("2006-01-02"), f.To.Format("2006-01-02"), f.For, invalid)
				}
			}
			continue DSW2
		}
		var invalid string
		if f.Invalid != "" {
			invalid = " INVALID " + strconv.Quote(f.Invalid)
		}
		tx.entry.Change("add person [%d] dsw from %s to %s for %q%s", f.From.Format("2006-01-02"), f.To.Format("2006-01-02"), f.For, invalid)
	}
	if p.VolgisticsID != op.VolgisticsID {
		tx.entry.Change("set person %q [%d] volgisticsID to %d", p.InformalName, p.ID, p.VolgisticsID)
	}
	if p.BackgroundCheck != op.BackgroundCheck {
		if p.BackgroundCheck == "" {
			tx.entry.Change("clear person %q [%d] backgroundCheck", p.InformalName, p.ID)
		} else {
			tx.entry.Change("set person %q [%d] backgroundCheck to %s", p.InformalName, p.ID, p.BackgroundCheck)
		}
	}
	if p.HoursToken != op.HoursToken {
		if p.HoursToken == "" {
			tx.entry.Change("clear person %q [%d] hoursToken", p.InformalName, p.ID)
		} else {
			tx.entry.Change("set person %q [%d] hoursToken to %s", p.InformalName, p.ID, p.HoursToken)
		}
	}
	if p.HoursReminder != op.HoursReminder {
		if p.HoursReminder {
			tx.entry.Change("set person %q [%d] hoursReminder", p.InformalName, p.ID)
		} else {
			tx.entry.Change("clear person %q [%d] hoursReminder", p.InformalName, p.ID)
		}
	}
	{
		var nm, om = p.DSWRegistrations, op.DSWRegistrations
		if nm == nil {
			nm = make(map[model.DSWClass]time.Time)
		}
		if om == nil {
			om = make(map[model.DSWClass]time.Time)
		}
		for c, nr := range nm {
			if or := om[c]; or != nr {
				if nr.IsZero() {
					tx.entry.Change("clear person %q [%d] dswRegistrations[%s]", p.InformalName, p.ID, model.DSWClassNames[c])
				} else {
					tx.entry.Change("set person %q [%d] dswRegistrations[%s] to %s", p.InformalName, p.ID, model.DSWClassNames[c], nr.Format("2006-01-02"))
				}
			}
		}
		for c := range om {
			if _, ok := nm[c]; !ok {
				tx.entry.Change("clear person %q [%d] dswRegistrations[%s]", p.InformalName, p.ID, model.DSWClassNames[c])
			}
		}
	}
	{
		var nm, om = p.DSWUntil, op.DSWUntil
		if nm == nil {
			nm = make(map[model.DSWClass]time.Time)
		}
		if om == nil {
			om = make(map[model.DSWClass]time.Time)
		}
		for c, nr := range nm {
			if or := om[c]; or != nr {
				if nr.IsZero() {
					tx.entry.Change("clear person %q [%d] dswUntil[%s]", p.InformalName, p.ID, model.DSWClassNames[c])
				} else {
					tx.entry.Change("set person %q [%d] dswUntil[%s] to %s", p.InformalName, p.ID, model.DSWClassNames[c], nr.Format("2006-01-02"))
				}
			}
		}
		for c := range om {
			if _, ok := nm[c]; !ok {
				tx.entry.Change("clear person %q [%d] dswUntil[%s]", p.InformalName, p.ID, model.DSWClassNames[c])
			}
		}
	}
	delete(tx.originalPeople, p.ID)
}
