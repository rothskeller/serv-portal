package class

import (
	"fmt"
	"slices"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// UpdaterFields are the fields that must be fetched prior to creating an
// Updater.
const UpdaterFields = FID | FType | FStart | FEnDesc | FEsDesc | FLimit | FReferrals

// Updater is a structure that can be filled with data for a new or changed
// class, and then later applied.  For creating new classes, it can simply be
// instantiated with new().  For updating existing classes, either *every* field
// in it must be set, or it should be instantiated with the Updater method of
// the class being changed.
type Updater struct {
	ID        ID
	Type      Type
	Start     string
	EnDesc    string
	EsDesc    string
	Limit     uint
	Referrals []uint
}

// Updater returns a new Updater for the specified class, with its data matching
// the current data for the class.  The class must have fetched UpdaterFields.
func (c *Class) Updater() *Updater {
	if c.fields&UpdaterFields != UpdaterFields {
		panic("Class.Updater called without fetching UpdaterFields")
	}
	return &Updater{
		ID:        c.id,
		Type:      c.ctype,
		Start:     c.start,
		EnDesc:    c.enDesc,
		EsDesc:    c.esDesc,
		Limit:     c.limit,
		Referrals: slices.Clone(c.referrals),
	}
}

const createSQL = `INSERT INTO class (id, type, start, en_desc, es_desc, elimit, referrals) VALUES (?,?,?,?,?,?,?)`

// Create creates a new class, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (c *Class) {
	c = new(Class)
	c.fields = UpdaterFields
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			c.id = u.ID
		} else {
			c.id = ID(phys.LastInsertRowID(storer))
		}
	})
	c.auditAndUpdate(storer, u, true)
	return c
}

const updateSQL = `UPDATE class SET type=?, start=?, en_desc=?, es_desc=?, elimit=?, referrals=? WHERE id=?`

// Update updates the existing class, with the data in the Updater.
func (c *Class) Update(storer phys.Storer, u *Updater) {
	if c.fields&UpdaterFields != UpdaterFields {
		panic("Class.Update called without fetching UpdaterFields")
	}
	phys.SQL(storer, updateSQL, func(stmt *phys.Stmt) {
		bindUpdater(stmt, u)
		stmt.BindInt(int(c.id))
		stmt.Step()
	})
	c.auditAndUpdate(storer, u, false)
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Type))
	stmt.BindText(u.Start)
	stmt.BindText(u.EnDesc)
	stmt.BindText(u.EsDesc)
	stmt.BindInt(int(u.Limit))
	var refmask uint64
	if u.Referrals != nil {
		for _, ref := range AllReferrals {
			refmask |= uint64(u.Referrals[ref]) << (ref * 8)
		}
	}
	stmt.BindInt(int(refmask))
}

func (c *Class) auditAndUpdate(storer phys.Storer, u *Updater, create bool) {
	context := fmt.Sprintf("Class %s %s [%d]", u.Type, u.Start, c.id)
	if create {
		context = "ADD " + context
	}
	if u.Type != c.ctype {
		phys.Audit(storer, "%s:: type = %s", context, u.Type)
		c.ctype = u.Type
	}
	if u.Start != c.start {
		phys.Audit(storer, "%s:: start = %s", context, u.Start)
		c.start = u.Start
	}
	if u.EnDesc != c.enDesc {
		phys.Audit(storer, "%s:: enDesc = %s", context, u.EnDesc)
		c.enDesc = u.EnDesc
	}
	if u.EsDesc != c.esDesc {
		phys.Audit(storer, "%s:: esDesc = %s", context, u.EsDesc)
		c.esDesc = u.EsDesc
	}
	if u.Limit != c.limit {
		phys.Audit(storer, "%s:: limit = %d", context, u.Limit)
		c.limit = u.Limit
	}
	for _, ref := range AllReferrals {
		var ur, cr uint
		if u.Referrals != nil {
			ur = u.Referrals[ref]
		}
		if c.referrals != nil {
			cr = c.referrals[ref]
		}
		if ur != cr {
			phys.Audit(storer, "%s:: referrals[%s] = %d", context, ref, ur)
			if c.referrals == nil {
				c.referrals = make([]uint, len(AllReferrals))
			}
			c.referrals[ref] = ur
		}
	}
}

const duplicateStartSQL = `SELECT 1 FROM class WHERE id!=? AND type=? AND start=?`

// DuplicateStart returns whether the type and start date specified in the
// Updater would be a duplicate if applied.
func (u *Updater) DuplicateStart(storer phys.Storer) (found bool) {
	phys.SQL(storer, duplicateStartSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(u.ID))
		stmt.BindInt(int(u.Type))
		stmt.BindText(u.Start)
		found = stmt.Step()
	})
	return found
}

// Delete deletes the receiver class.
func (c *Class) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM class WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(c.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE Class %s %s [%d]", c.Type(), c.Start(), c.ID())
}
