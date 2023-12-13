package person

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

const joinFields = FAddresses | FBGChecks | FDSWRegistrations | FNotes | FEmContacts | FPrivLevels

var withIDSQLCache map[Fields]string

// WithID returns the person with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (p *Person) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p WHERE p.id=?")
		withIDSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			p = new(Person)
			p.Scan(stmt, fields&^joinFields)
			p.id = id
			p.fields |= FID
			p.readJoins(storer, fields)
		}
	})
	return p
}

var withEmailSQLCache map[Fields]string

// WithEmail returns the person with the specified email address, or nil if it
// does not exist.  Note that it only checks people's primary email addresses,
// not their secondary email addresses if any.
func WithEmail(storer phys.Storer, email string, fields Fields) (p *Person) {
	if withEmailSQLCache == nil {
		withEmailSQLCache = make(map[Fields]string)
	}
	if _, ok := withEmailSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p WHERE p.email=?")
		withEmailSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, withEmailSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		stmt.BindText(email)
		if stmt.Step() {
			p = new(Person)
			p.Scan(stmt, fields&^joinFields)
			p.email = email
			p.fields |= FEmail
			p.readJoins(storer, fields)
		}
	})
	return p
}

var withHoursTokenSQLCache map[Fields]string

// WithHoursToken returns the person with the specified password reset token,
// or nil if it does not exist.
func WithHoursToken(storer phys.Storer, token string, fields Fields) (p *Person) {
	if withHoursTokenSQLCache == nil {
		withHoursTokenSQLCache = make(map[Fields]string)
	}
	if _, ok := withHoursTokenSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p WHERE p.hours_token=?")
		withHoursTokenSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, withHoursTokenSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		stmt.BindText(token)
		if stmt.Step() {
			p = new(Person)
			p.Scan(stmt, fields&^joinFields)
			p.hoursToken = token
			p.fields |= FHoursToken
			p.readJoins(storer, fields)
		}
	})
	return p
}

var withPWResetTokenSQLCache map[Fields]string

// WithPWResetToken returns the person with the specified password reset token,
// or nil if it does not exist.
func WithPWResetToken(storer phys.Storer, token string, fields Fields) (p *Person) {
	if withPWResetTokenSQLCache == nil {
		withPWResetTokenSQLCache = make(map[Fields]string)
	}
	if _, ok := withPWResetTokenSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p WHERE p.pwreset_token=?")
		withPWResetTokenSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, withPWResetTokenSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		stmt.BindText(token)
		if stmt.Step() {
			p = new(Person)
			p.Scan(stmt, fields&^joinFields)
			p.pwresetToken = token
			p.fields |= FPWResetToken
			p.readJoins(storer, fields)
		}
	})
	return p
}

var withUnsubscribeTokenSQLCache map[Fields]string

// WithUnsubscribeToken returns the person with the specified password reset token,
// or nil if it does not exist.
func WithUnsubscribeToken(storer phys.Storer, token string, fields Fields) (p *Person) {
	if withUnsubscribeTokenSQLCache == nil {
		withUnsubscribeTokenSQLCache = make(map[Fields]string)
	}
	if _, ok := withUnsubscribeTokenSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p WHERE p.pwreset_token=?")
		withUnsubscribeTokenSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, withUnsubscribeTokenSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		stmt.BindText(token)
		if stmt.Step() {
			p = new(Person)
			p.Scan(stmt, fields&^joinFields)
			p.unsubscribeToken = token
			p.fields |= FUnsubscribeToken
			p.readJoins(storer, fields)
		}
	})
	return p
}

var allSQLCache map[Fields]string

// All reads each person from the database, in order by sortName.
func All(storer phys.Storer, fields Fields, fn func(*Person)) {
	if fields&joinFields != 0 {
		fields |= FID
	}
	if allSQLCache == nil {
		allSQLCache = make(map[Fields]string)
	}
	if _, ok := allSQLCache[fields&^joinFields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^joinFields)
		sb.WriteString(" FROM person p ORDER BY p.sort_name")
		allSQLCache[fields&^joinFields] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields&^joinFields], func(stmt *phys.Stmt) {
		var p Person
		for stmt.Step() {
			p.Scan(stmt, fields&^joinFields)
			p.readJoins(storer, fields)
			fn(&p)
		}
	})
}
