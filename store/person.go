package store

import (
	"bytes"

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
		tx.entry.Change("set person %q [%d] username to %q", p.ID, p.InformalName, p.Username)
	}
	if p.InformalName != op.InformalName {
		tx.entry.Change("set person %q [%d] informalName to %q", p.ID, p.InformalName, p.InformalName)
	}
	if p.FormalName != op.FormalName {
		tx.entry.Change("set person %q [%d] formalName to %q", p.ID, p.InformalName, p.FormalName)
	}
	if p.SortName != op.SortName {
		tx.entry.Change("set person %q [%d] sortName to %q", p.ID, p.InformalName, p.SortName)
	}
	if p.CallSign != op.CallSign {
		tx.entry.Change("set person %q [%d] callSign to %q", p.ID, p.InformalName, p.CallSign)
	}
	if p.Email != op.Email {
		tx.entry.Change("set person %q [%d] email to %q", p.ID, p.InformalName, p.Email)
	}
	if p.Email2 != op.Email2 {
		tx.entry.Change("set person %q [%d] email2 to %q", p.ID, p.InformalName, p.Email2)
	}
	if p.NoEmail != op.NoEmail {
		if p.NoEmail {
			tx.entry.Change("set person %q [%d] noEmail flag", p.ID, p.InformalName)
		} else {
			tx.entry.Change("clear person %q [%d] noEmail flag", p.ID, p.InformalName)
		}
	}
	if p.HomeAddress.Address != op.HomeAddress.Address {
		if p.HomeAddress.FireDistrict != 0 {
			tx.entry.Change("set person %q [%d] homeAddress to %q (%f, %f) (district %d)", p.ID, p.InformalName, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude, p.HomeAddress.FireDistrict)
		} else if p.HomeAddress.Latitude != 0 {
			tx.entry.Change("set person %q [%d] homeAddress to %q (%f, %f)", p.ID, p.InformalName, p.HomeAddress.Address, p.HomeAddress.Latitude, p.HomeAddress.Longitude)
		} else {
			tx.entry.Change("set person %q [%d] homeAddress to %q", p.ID, p.InformalName, p.HomeAddress.Address)
		}
	}
	if p.WorkAddress.SameAsHome {
		if !op.WorkAddress.SameAsHome {
			tx.entry.Change("set person %q [%d] workAddress sameAsHome flag", p.ID, p.InformalName)
		}
	} else if op.WorkAddress.SameAsHome && p.WorkAddress.Address == "" {
		tx.entry.Change("clear person %q [%d] workAddress sameAsHome flag", p.InformalName, p.ID)
	} else if p.WorkAddress.Address != op.WorkAddress.Address {
		if p.WorkAddress.FireDistrict != 0 {
			tx.entry.Change("set person %q [%d] workAddress to %q (%f, %f) (district %d)", p.ID, p.InformalName, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude, p.WorkAddress.FireDistrict)
		} else if p.WorkAddress.Latitude != 0 {
			tx.entry.Change("set person %q [%d] workAddress to %q (%f, %f)", p.ID, p.InformalName, p.WorkAddress.Address, p.WorkAddress.Latitude, p.WorkAddress.Longitude)
		} else {
			tx.entry.Change("set person %q [%d] workAddress to %q", p.ID, p.InformalName, p.WorkAddress.Address)
		}
	}
	if p.MailAddress.SameAsHome {
		if !op.MailAddress.SameAsHome {
			tx.entry.Change("set person %q [%d] mailAddress sameAsHome flag", p.ID, p.InformalName)
		}
	} else if op.MailAddress.SameAsHome && p.MailAddress.Address == "" {
		tx.entry.Change("clear person %q [%d] mailAddress sameAsHome flag", p.InformalName, p.ID)
	} else if p.MailAddress.Address != op.MailAddress.Address {
		tx.entry.Change("set person %q [%d] mailAddress to %q", p.ID, p.InformalName, p.MailAddress.Address)
	}
	if p.CellPhone != op.CellPhone {
		tx.entry.Change("set person %q [%d] cellPhone to %q", p.ID, p.InformalName, p.CellPhone)
	}
	if p.HomePhone != op.HomePhone {
		tx.entry.Change("set person %q [%d] homePhone to %q", p.ID, p.InformalName, p.HomePhone)
	}
	if p.WorkPhone != op.WorkPhone {
		tx.entry.Change("set person %q [%d] workPhone to %q", p.ID, p.InformalName, p.WorkPhone)
	}
	if p.NoText != op.NoText {
		if p.NoText {
			tx.entry.Change("set person %q [%d] noText flag", p.ID, p.InformalName)
		} else {
			tx.entry.Change("clear person %q [%d] noText flag", p.ID, p.InformalName)
		}
	}
	if len(op.Password) == 0 && len(p.Password) != 0 {
		tx.entry.Change("set person %q [%d] password", p.ID, p.InformalName)
	} else if !bytes.Equal(p.Password, op.Password) {
		tx.entry.Change("change person %q [%d] password", p.ID, p.InformalName)
	}
NOTES1:
	for _, on := range op.Notes {
		for _, n := range p.Notes {
			if n.Date == on.Date && n.Note == on.Note && n.Privilege == on.Privilege {
				continue NOTES1
			}
		}
		tx.entry.Change("remove person %q [%d] note %q at %s with privilege %s", p.ID, p.InformalName, on.Note, on.Date, model.PrivilegeNames[on.Privilege])
	}
NOTES2:
	for _, n := range p.Notes {
		for _, on := range op.Notes {
			if n.Date == on.Date && n.Note == on.Note && n.Privilege == on.Privilege {
				continue NOTES2
			}
		}
		tx.entry.Change("add person %q [%d] note %q at %s with privilege %s", p.ID, p.InformalName, n.Note, n.Date, model.PrivilegeNames[n.Privilege])
	}
	delete(tx.originalPeople, p.ID)
}
