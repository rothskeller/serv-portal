package listperson

import (
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
)

const subscribeSQL = `
INSERT INTO list_person (list, person, sender, sub, unsub) VALUES (?,?,FALSE,TRUE,FALSE)
ON CONFLICT DO UPDATE SET sub=TRUE, unsub=FALSE WHERE unsub OR NOT sub`

// Subscribe subscribes a person to a list.  It does not check whether they are
// allowed to subscribe to it; the caller should do that.  If the person is
// explicitly unsubscribed from the list, it removes that notation.
func Subscribe(storer phys.Storer, l *list.List, p *person.Person) {
	phys.SQL(storer, subscribeSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(l.ID))
		stmt.BindInt(int(p.ID()))
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "List %q [%d]:: SUBSCRIBE Person %q [%d]", l.Name, l.ID, p.InformalName(), p.ID())
	}
}

const unsubscribeSQL = `
INSERT INTO list_person (list, person, sender, sub, unsub) VALUES (?,?,FALSE,FALSE,TRUE)
ON CONFLICT DO UPDATE SET unsub=TRUE WHERE NOT unsub`

// Unsubscribe unsubscribes a person from a list.
func Unsubscribe(storer phys.Storer, l *list.List, p *person.Person) {
	phys.SQL(storer, unsubscribeSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(l.ID))
		stmt.BindInt(int(p.ID()))
		stmt.Step()
	})
	if phys.RowsAffected(storer) != 0 {
		phys.Audit(storer, "List %q [%d]:: UNSUBSCRIBE Person %q [%d]", l.Name, l.ID, p.InformalName(), p.ID())
	}
}
