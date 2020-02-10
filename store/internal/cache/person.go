package cache

import (
	"sort"

	"sunnyvaleserv.org/portal/model"
)

// FetchPerson retrieves a single person from the database by ID.  It returns
// nil if no such person exists.
func (tx *Tx) FetchPerson(id model.PersonID) (p *model.Person) {
	if p = tx.people[id]; p != nil {
		return p
	}
	if p = tx.Tx.FetchPerson(id); p != nil {
		tx.people[id] = p
	}
	return p
}

// FetchPersonByUsername retrieves a single person from the database, given
// their username.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByUsername(username string) (p *model.Person) {
	if p = tx.Tx.FetchPersonByUsername(username); p != nil {
		if p2 := tx.people[p.ID]; p2 != nil {
			return p2
		}
		tx.people[p.ID] = p
	}
	return p
}

// FetchPersonByPWResetToken retrieves a single person from the database, given
// a password reset token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByPWResetToken(token string) (p *model.Person) {
	if p = tx.Tx.FetchPersonByPWResetToken(token); p != nil {
		if p2 := tx.people[p.ID]; p2 != nil {
			return p2
		}
		tx.people[p.ID] = p
	}
	return p
}

// FetchPersonByCellPhone retrieves a single person from the database, given a
// cell phone number.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByCellPhone(number string) (p *model.Person) {
	if p = tx.Tx.FetchPersonByCellPhone(number); p != nil {
		if p2 := tx.people[p.ID]; p2 != nil {
			return p2
		}
		tx.people[p.ID] = p
	}
	return p
}

// FetchPersonByUnsubscribe retrieves a single person from the database, given
// an unsubscribe token.  It returns nil if no such person exists.
func (tx *Tx) FetchPersonByUnsubscribe(token string) (p *model.Person) {
	if p = tx.Tx.FetchPersonByUnsubscribe(token); p != nil {
		if p2 := tx.people[p.ID]; p2 != nil {
			return p2
		}
		tx.people[p.ID] = p
	}
	return p
}

// FetchPeople returns all of the people in the database, in order by sortname.
func (tx *Tx) FetchPeople() (people []*model.Person) {
	if tx.personList == nil {
		tx.personList = tx.Tx.FetchPeople()
		for i, p := range tx.personList {
			if p2 := tx.people[p.ID]; p2 != nil {
				tx.personList[i] = p2
			} else {
				tx.people[p.ID] = p
			}
		}
	}
	return tx.personList
}

// CreatePerson creates a new person in the database.
func (tx *Tx) CreatePerson(p *model.Person) {
	tx.Tx.CreatePerson(p)
	tx.people[p.ID] = p
	if tx.personList != nil {
		tx.personList = append(tx.personList, p)
		sort.Sort(model.PersonSort(tx.personList))
	}
}

// UpdatePerson updates a person in the database.
func (tx *Tx) UpdatePerson(p *model.Person) {
	if tx.people[p.ID] != p {
		panic("must modify people in place")
	}
	tx.Tx.UpdatePerson(p)
	if tx.personList != nil {
		sort.Sort(model.PersonSort(tx.personList))
	}
}
