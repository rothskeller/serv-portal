package listperson

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
)

const canSendSQL = `SELECT sender FROM list_person WHERE list=? AND person=?`

// CanSend returns whether the specified person can send messages to the
// specified list.
func CanSend(storer phys.Storer, pid person.ID, lid list.ID) (canSend bool) {
	phys.SQL(storer, canSendSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(lid))
		stmt.BindInt(int(pid))
		if stmt.Step() {
			canSend = stmt.ColumnBool()
		}
	})
	return canSend
}

const canSendTextSQL = `SELECT 1 FROM list l, list_person lp WHERE lp.sender AND lp.list=l.id AND l.type=2 AND lp.person=?`

// CanSendText returns whether the specified person can send to any text (SMS)
// list.
func CanSendText(storer phys.Storer, pid person.ID) (canSend bool) {
	phys.SQL(storer, canSendTextSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(pid))
		canSend = stmt.Step()
	})
	return
}

const subscribedSQL = `SELECT sub, unsub FROM list_person WHERE list=? AND person=?`

// Subscribed returns whether the specified person is subscribed to the
// specified, and also whether they have explicitly unsubscribed from it.  If
// both flags are true, the unsubscribe takes precedence.
func Subscribed(storer phys.Storer, p *person.Person, l *list.List) (subscribed, unsubscribed bool) {
	phys.SQL(storer, subscribedSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(l.ID))
		stmt.BindInt(int(p.ID()))
		if stmt.Step() {
			subscribed = stmt.ColumnBool()
			unsubscribed = stmt.ColumnBool()
		}
	})
	return subscribed, unsubscribed
}

var allSQLCache map[person.Fields]string

// All fetches the people who are subscribed to, unsubscribed from, and/or
// allowed to send to the specified list, in order by name.
func All(storer phys.Storer, lid list.ID, fields person.Fields, fn func(p *person.Person, sender, sub, unsub bool)) {
	if allSQLCache == nil {
		allSQLCache = make(map[person.Fields]string)
	}
	if allSQLCache[fields] == "" {
		var sb strings.Builder
		sb.WriteString("SELECT lp.sender, lp.sub, lp.unsub")
		if fields != 0 {
			sb.WriteString(", ")
			person.ColumnList(&sb, fields)
		}
		sb.WriteString(" FROM list_person lp, person p WHERE lp.list=? AND lp.person=p.id ORDER BY p.sort_name")
		allSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields], func(stmt *phys.Stmt) {
		var p person.Person
		var sender, sub, unsub bool

		stmt.BindInt(int(lid))
		for stmt.Step() {
			sender = stmt.ColumnBool()
			sub = stmt.ColumnBool()
			unsub = stmt.ColumnBool()
			p.Scan(stmt, fields)
			fn(&p, sender, sub, unsub)
		}
	})
}

const subscriptionsByPersonSQL = `SELECT ` + list.ColumnList + ` FROM list l, list_person lp WHERE l.id=lp.list AND lp.person=? AND lp.sub AND NOT lp.unsub ORDER BY l.type, l.name`

// SubscriptionsByPerson fetches each of the lists to which the specified person
// is subscribed, in order by type and then name.
func SubscriptionsByPerson(storer phys.Storer, pid person.ID, fn func(*list.List)) {
	phys.SQL(storer, subscriptionsByPersonSQL, func(stmt *phys.Stmt) {
		var l list.List

		stmt.BindInt(int(pid))
		for stmt.Step() {
			l.Scan(stmt)
			fn(&l)
		}
	})
}

const sendersByPersonSQL = `SELECT ` + list.ColumnList + ` FROM list l, list_person lp WHERE l.id=lp.list AND lp.person=? AND lp.sender ORDER BY l.type, l.name`

// SendersByPerson fetches each of the lists to which the specified person
// is allowed to send without moderation, in order by type and then name.
func SendersByPerson(storer phys.Storer, pid person.ID, fn func(*list.List)) {
	phys.SQL(storer, sendersByPersonSQL, func(stmt *phys.Stmt) {
		var l list.List

		stmt.BindInt(int(pid))
		for stmt.Step() {
			l.Scan(stmt)
			fn(&l)
		}
	})
}

const subscriptionRightsSQL1 = `SELECT l.id, l.type, l.name, `
const subscriptionRightsSQL2 = `, lr.submodel FROM list l, role r, list_role lr, person_role pr WHERE pr.person=? AND pr.role=r.id AND r.id=lr.role AND lr.submodel!=0 AND lr.list=l.id ORDER BY l.type, l.name, r.priority`

var subscriptionRightsSQLCache map[role.Fields]string

// SubscriptionRights fetches each of the lists to which the specified person is
// allowed to subscribe.  For each, it returns the list, the role held by the
// person that grants subscription rights to the list, and the subscription
// model for that role.  Note that if a person is granted subscription rights to
// a list by multiple roles held by that person, the same list will be fetched
// more than once.  Lists are fetched in order by type and then name.
func SubscriptionRights(storer phys.Storer, pid person.ID, roleFields role.Fields, fn func(*list.List, *role.Role, listrole.SubscriptionModel)) {
	if roleFields == 0 {
		roleFields = role.FID
	}
	if subscriptionRightsSQLCache == nil {
		subscriptionRightsSQLCache = make(map[role.Fields]string)
	}
	if _, ok := subscriptionRightsSQLCache[roleFields]; !ok {
		var sb strings.Builder
		sb.WriteString(subscriptionRightsSQL1)
		role.ColumnList(&sb, roleFields)
		sb.WriteString(subscriptionRightsSQL2)
		subscriptionRightsSQLCache[roleFields] = sb.String()
	}
	phys.SQL(storer, subscriptionRightsSQLCache[roleFields], func(stmt *phys.Stmt) {
		var (
			l  list.List
			r  role.Role
			sm listrole.SubscriptionModel
		)
		stmt.BindInt(int(pid))
		for stmt.Step() {
			l.Scan(stmt)
			r.Scan(stmt, roleFields)
			sm = listrole.SubscriptionModel(stmt.ColumnInt())
			fn(&l, &r, sm)
		}
	})
}
