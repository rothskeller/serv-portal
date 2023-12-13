package textmsg

import (
	"fmt"
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
)

// Updater is a structure that can be filled with data for a new text
// message, and then later applied.
type Updater struct {
	ID        ID
	Sender    *person.Person
	Timestamp time.Time
	Message   string
	Lists     []*list.List
}

const createSQL = `INSERT INTO textmsg (id, sender, timestamp, message) VALUES (?,?,?,?)`

// Create creates a new text message, with the data in the Updater.
func Create(storer phys.Storer, u *Updater) (t *TextMessage) {
	t = new(TextMessage)
	t.fields = FID | FSender | FTimestamp | FMessage | FLists
	phys.SQL(storer, createSQL, func(stmt *phys.Stmt) {
		stmt.BindNullInt(int(u.ID))
		bindUpdater(stmt, u)
		stmt.Step()
		if u.ID != 0 {
			t.id = u.ID
		} else {
			t.id = ID(phys.LastInsertRowID(storer))
		}
	})
	storeLists(storer, t.id, u.Lists)
	t.auditAndUpdate(storer, u)
	return t
}

func bindUpdater(stmt *phys.Stmt, u *Updater) {
	stmt.BindInt(int(u.Sender.ID()))
	stmt.BindText(u.Timestamp.In(time.Local).Format(timestampFormat))
	stmt.BindText(u.Message)
}

const addListSQL = `INSERT INTO textmsg_list (textmsg, list, name) VALUES (?,?,?)`

func storeLists(storer phys.Storer, tid ID, lists []*list.List) {
	if len(lists) == 0 {
		return
	}
	phys.SQL(storer, addListSQL, func(stmt *phys.Stmt) {
		for _, tl := range lists {
			stmt.BindInt(int(tid))
			stmt.BindInt(int(tl.ID))
			stmt.BindText(tl.Name)
			stmt.Step()
			stmt.Reset()
		}
	})
}

func (t *TextMessage) auditAndUpdate(storer phys.Storer, u *Updater) {
	context := fmt.Sprintf("ADD TextMessage %d", t.id)
	phys.Audit(storer, "%s:: sender = %q [%d]", context, u.Sender.InformalName(), u.Sender.ID())
	t.sender = u.Sender.ID()
	phys.Audit(storer, "%s:: timestamp = %s", context, u.Timestamp.In(time.Local).Format(timestampFormat))
	t.timestamp = u.Timestamp.In(time.Local)
	phys.Audit(storer, "%s:: message = %q", context, u.Message)
	t.message = u.Message
	for i, tl := range u.Lists {
		phys.Audit(storer, "%s:: Lists[%d] = %q [%d]", context, i+1, tl.Name, tl.ID)
		t.lists = append(t.lists, TextToList{ID: tl.ID, Name: tl.Name})
	}
}

// Delete deletes the receiver text message.
func (t *TextMessage) Delete(storer phys.Storer) {
	phys.SQL(storer, `DELETE FROM textmsg WHERE id=?`, func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		stmt.Step()
	})
	phys.Audit(storer, "DELETE TextMessage %d", t.ID())
}
