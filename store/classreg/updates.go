package classreg

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FClass | FPerson | FRegisteredBy | FFirstName | FLastName | FEmail | FCellPhone

// Updater is a structure that can be filled with data for a new or changed
// class, and then later applied.  For creating new classes, it can simply be
// instantiated with new().  For updating existing classes, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the class being changed.
type Updater struct {
	ID           ID
	Class        *class.Class
	Person       *person.Person
	RegisteredBy *person.Person
	FirstName    string
	LastName     string
	Email        string
	CellPhone    string
}

// Updater returns a new Updater for the specified class, with its data matching
// the current data for the class.  The class must have fetched UpdaterFields.
// The class, person, and registeredBy pointers *may* be provided to avoid
// lookups.
func (cr *ClassReg) Updater(storer phys.Storer, c *class.Class, p, rb *person.Person) *Updater {
	if cr.fields&UpdaterFields != UpdaterFields {
		panic("ClassReg.Updater called without fetching UpdaterFields")
	}
	if c == nil {
		c = class.WithID(storer, cr.class, class.FID|class.FType|class.FStart)
	}
	if rb == nil {
		rb = person.WithID(storer, cr.registeredBy, person.FID|person.FInformalName)
	}
	if p == nil && cr.person != 0 {
		if cr.person == cr.registeredBy {
			p = rb
		} else {
			p = person.WithID(storer, cr.person, person.FID|person.FInformalName)
		}
	}
	return &Updater{
		ID:           cr.id,
		Class:        c,
		Person:       p,
		RegisteredBy: rb,
		FirstName:    cr.firstName,
		LastName:     cr.lastName,
		Email:        cr.email,
		CellPhone:    cr.cellPhone,
	}
}

const createSQL = `INSERT INTO classreg (id, class, person, registered_by, first_name, last_name, email, cell_phone) VALUES (?,?,?,?,?,?,?,?)`

// Create creates a new class registration with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (cr *ClassReg) {
	cr = new(ClassReg)
	cr.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			cr.id = u.ID
		} else {
			cr.id = ID(phys.LastInsertRowID(storer))
		}
	})
	cr.auditAndUpdate(storer, u, true)
	return cr
}

const updateSQL = `UPDATE classreg SET class=?, person=?, registered_by=?, first_name=?, last_name=?, email=?, cell_phone=? WHERE id=?`

// Update updates the existing class, with the data in the Updater.
func (cr *ClassReg) Update(storer phys.Storer, u *Updater) {
	if cr.fields&UpdaterFields != UpdaterFields {
		panic("ClassReg.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(cr.id))
		stmt.Step()
	})
	cr.auditAndUpdate(storer, u, false)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Class.ID()))
	stmt.BindNullInt(int(u.Person.ID()))
	stmt.BindInt(int(u.RegisteredBy.ID()))
	stmt.BindText(u.FirstName)
	stmt.BindText(u.LastName)
	stmt.BindText(u.Email)
	stmt.BindText(u.CellPhone)
}

func (cr *ClassReg) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	var context string
	if create {
		context = fmt.Sprintf("Class %s %s [%d]:: ADD Registration %d", u.Class.Type(), u.Class.Start(), u.Class.ID(), cr.id)
	} else {
		context = fmt.Sprintf("Class %s %s [%d]:: Registration %d", u.Class.Type(), u.Class.Start(), u.Class.ID(), cr.id)
	}
	if u.Class.ID() != cr.class {
		phys.Audit(storer, "%s:: class = %s %s [%d]", context, u.Class.Type(), u.Class.Start(), u.Class.ID())
		cr.class = u.Class.ID()
	}
	if u.Person.ID() != cr.person {
		if u.Person != nil {
			phys.Audit(storer, "%s:: person = %s [%d]", context, u.Person.InformalName(), u.Person.ID())
		} else {
			phys.Audit(storer, "%s:: person = nil", context)
		}
		cr.person = u.Person.ID()
	}
	if u.RegisteredBy.ID() != cr.registeredBy {
		phys.Audit(storer, "%s:: registeredBy = %s [%d]", context, u.RegisteredBy.InformalName(), u.RegisteredBy.ID())
		cr.registeredBy = u.RegisteredBy.ID()
	}
	if u.FirstName != cr.firstName {
		phys.Audit(storer, "%s:: firstName = %s", context, u.FirstName)
		cr.firstName = u.FirstName
	}
	if u.LastName != cr.lastName {
		phys.Audit(storer, "%s:: lastName = %s", context, u.LastName)
		cr.lastName = u.LastName
	}
	if u.Email != cr.email {
		phys.Audit(storer, "%s:: email = %s", context, u.Email)
		cr.email = u.Email
	}
	if u.CellPhone != cr.cellPhone {
		phys.Audit(storer, "%s:: cellPhone = %s", context, u.CellPhone)
		cr.cellPhone = u.CellPhone
	}
}

// Delete deletes the receiver class registration.  The class *may* be provided
// to avoid a lookup.
func (cr *ClassReg) Delete(storer phys.Storer, c *class.Class) {
	if c == nil {
		c = class.WithID(storer, cr.class, class.FID|class.FType|class.FStart)
	}
	phys.SQL(storer, `DELETE FROM classreg WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(cr.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "Class %s %s [%d]:: DELETE Registration %s %s [%d]", c.Type(), c.Start(), c.ID(), cr.FirstName(), cr.LastName(), cr.ID())
}
