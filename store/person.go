package store

import (
	"bytes"
	"time"

	"sunnyvaleserv.org/portal/model"
)

// CreatePerson creates a new person in the database.
func (tx *Tx) CreatePerson(p *model.Person) {
	tx.Tx.CreatePerson(p)
	tx.entry.Change("create person [%d]", p.ID)
	tx.entry.Change("set person [%d] informalName to %q", p.ID, p.InformalName)
	tx.entry.Change("set person [%d] formalName to %q", p.ID, p.FormalName)
	tx.entry.Change("set person [%d] sortName to %q", p.ID, p.SortName)
	if p.CallSign != "" {
		tx.entry.Change("set person [%d] callSign to %q", p.ID, p.CallSign)
	}
	if p.Birthdate != "" {
		tx.entry.Change("set person [%d] birthdate to %s", p.ID, p.Birthdate)
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
		tx.entry.Change("add person [%d] note %q at %s visibility %s", p.ID, n.Note, n.Date, n.Visibility)
	}
	tx.entry.Change("set person [%d] unsubscribeToken to %s", p.ID, p.UnsubscribeToken)
	if p.VolgisticsID != 0 {
		tx.entry.Change("set person [%d] volgisticsID to %d", p.ID, p.VolgisticsID)
	}
	if p.VolgisticsPending {
		tx.entry.Change("set person [%d] volgisticsPending to true", p.ID)
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
	if p.Identification != 0 {
		for _, t := range model.AllIdentTypes {
			if p.Identification&t != 0 {
				tx.entry.Change("add person [%d] identification %s", p.ID, model.IdentTypeNames[t])
			}
		}
	}
	for r, direct := range p.Roles {
		if direct {
			tx.entry.Change("add person [%d] role %q [%d]", p.ID, tx.FetchRole(r).Name, r)
		}
	}
	for _, bc := range p.BGChecks {
		if bc.Date != "" {
			tx.entry.Change("add person [%d] bgCheck %s on %s", p.ID, bc.Type.MaskString(), bc.Date)
		} else if bc.Assumed {
			tx.entry.Change("add person [%d] bgCheck %s assumed", p.ID, bc.Type.MaskString())
		} else {
			tx.entry.Change("add person [%d] bgCheck %s", p.ID, bc.Type.MaskString())
		}
	}
	for _, em := range p.EmContacts {
		tx.entry.Change("add person [%d] emContact name %q home %s cell %s rel %s", p.ID, em.Name, em.HomePhone, em.CellPhone, em.Relationship)
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
	if p.EmContacts != nil {
		op.EmContacts = make([]*model.EmContact, len(p.EmContacts))
		for i := range p.EmContacts {
			oem := *p.EmContacts[i]
			op.EmContacts[i] = &oem
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
	if p.Birthdate != op.Birthdate {
		tx.entry.Change("set person %q [%d] birthdate to %s", p.InformalName, p.ID, p.Birthdate)
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
			if n.Date == on.Date && n.Note == on.Note && n.Visibility == on.Visibility {
				continue NOTES1
			}
		}
		tx.entry.Change("remove person %q [%d] note %q at %s", p.InformalName, p.ID, on.Note, on.Date)
	}
NOTES2:
	for _, n := range p.Notes {
		for _, on := range op.Notes {
			if n.Date == on.Date && n.Note == on.Note && n.Visibility == on.Visibility {
				continue NOTES2
			}
		}
		tx.entry.Change("add person %q [%d] note %q at %s visibility %s", p.InformalName, p.ID, n.Note, n.Date, n.Visibility)
	}
	if p.UnsubscribeToken != op.UnsubscribeToken {
		tx.entry.Change("change person %q [%d] unsubscribeToken to %s", p.InformalName, p.ID, p.UnsubscribeToken)
	}
	if p.VolgisticsID != op.VolgisticsID {
		tx.entry.Change("set person %q [%d] volgisticsID to %d", p.InformalName, p.ID, p.VolgisticsID)
	}
	if p.VolgisticsPending != op.VolgisticsPending {
		tx.entry.Change("set person %q [%d] volgisticsPending to %v", p.InformalName, p.ID, p.VolgisticsPending)
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
	for _, t := range model.AllIdentTypes {
		if op.Identification&t != 0 && p.Identification&t == 0 {
			tx.entry.Change("remove person %q [%d] identification %s", p.InformalName, p.ID, model.IdentTypeNames[t])
		} else if op.Identification&t == 0 && p.Identification&t != 0 {
			tx.entry.Change("add person %q [%d] identification %s", p.InformalName, p.ID, model.IdentTypeNames[t])
		}
	}
	for r, direct := range p.Roles {
		if !direct {
			continue
		}
		var found = false
		for or, odirect := range op.Roles {
			if r == or && odirect {
				found = true
				break
			}
		}
		if !found {
			tx.entry.Change("add person %q [%d] role %q [%d]", p.InformalName, p.ID, tx.FetchRole(r).Name, r)
		}
	}
	for or, odirect := range op.Roles {
		if !odirect {
			continue
		}
		var found = false
		for r, direct := range p.Roles {
			if r == or && direct {
				found = true
				break
			}
		}
		if !found {
			tx.entry.Change("remove person %q [%d] role %q [%d]", p.InformalName, p.ID, tx.FetchRole(or).Name, or)
		}
	}
	for _, c := range p.BGChecks {
		var found = false
		for _, oc := range op.BGChecks {
			if c.Equal(oc) {
				found = true
				break
			}
		}
		if !found {
			if c.Date != "" {
				tx.entry.Change("add person [%d] bgCheck %s on %s", p.ID, c.Type.MaskString(), c.Date)
			} else if c.Assumed {
				tx.entry.Change("add person [%d] bgCheck %s assumed", p.ID, c.Type.MaskString())
			} else {
				tx.entry.Change("add person [%d] bgCheck %s", p.ID, c.Type.MaskString())
			}
		}
	}
	for _, oc := range op.BGChecks {
		var found = false
		for _, c := range p.BGChecks {
			if oc.Equal(c) {
				found = true
				break
			}
		}
		if !found {
			if oc.Date != "" {
				tx.entry.Change("remove person [%d] bgCheck %s on %s", p.ID, oc.Type.MaskString(), oc.Date)
			} else if oc.Assumed {
				tx.entry.Change("remove person [%d] bgCheck %s assumed", p.ID, oc.Type.MaskString())
			} else {
				tx.entry.Change("remove person [%d] bgCheck %s", p.ID, oc.Type.MaskString())
			}
		}
	}
EMCONTACT1:
	for _, oem := range op.EmContacts {
		for _, em := range p.EmContacts {
			if em.Name == oem.Name && em.HomePhone == oem.HomePhone && em.CellPhone == oem.CellPhone && em.Relationship == oem.Relationship {
				continue EMCONTACT1
			}
		}
		tx.entry.Change("remove person %q [%d] emContact name %q home %s cell %s rel %s", p.InformalName, p.ID, oem.Name, oem.HomePhone, oem.CellPhone, oem.Relationship)
	}
EMCONTACT2:
	for _, em := range p.EmContacts {
		for _, oem := range op.EmContacts {
			if em.Name == oem.Name && em.HomePhone == oem.HomePhone && em.CellPhone == oem.CellPhone && em.Relationship == oem.Relationship {
				continue EMCONTACT2
			}
		}
		tx.entry.Change("add person %q [%d] emContact name %q home %s cell %s rel %s", p.InformalName, p.ID, em.Name, em.HomePhone, em.CellPhone, em.Relationship)
	}
	tx.Tx.UpdatePerson(p)
	delete(tx.originalPeople, p.ID)
}
