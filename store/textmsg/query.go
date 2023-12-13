package textmsg

import (
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

var withIDSQLCache map[Fields]string

// WithID returns the text message with the specified ID, or nil if it does not
// exist.
func WithID(storer phys.Storer, id ID, fields Fields) (t *TextMessage) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields&^FLists]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FLists)
		sb.WriteString(" FROM textmsg t WHERE t.id=?")
		withIDSQLCache[fields&^FLists] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields&^FLists], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			t = new(TextMessage)
			t.Scan(stmt, fields&^FLists)
			t.id = id
			t.fields |= FID
			if fields&FLists != 0 {
				t.readLists(storer)
			}
		}
	})
	return t
}

var withNumberSQLCache map[Fields]string

// WithNumber returns the most recent text message sent to the specified number,
// or nil if none has been sent to it.
func WithNumber(storer phys.Storer, number string, fields Fields) (t *TextMessage) {
	if withNumberSQLCache == nil {
		withNumberSQLCache = make(map[Fields]string)
	}
	if _, ok := withNumberSQLCache[fields&^FLists]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FLists)
		sb.WriteString(" FROM textmsg t, textmsg_number tn WHERE tn.number=? AND tn.textmsg=t.id")
		withNumberSQLCache[fields&^FLists] = sb.String()
	}
	phys.SQL(storer, withNumberSQLCache[fields&^FLists], func(stmt *phys.Stmt) {
		stmt.BindText(number)
		if stmt.Step() {
			t = new(TextMessage)
			t.Scan(stmt, fields&^FLists)
			if fields&FLists != 0 {
				t.readLists(storer)
			}
		}
	})
	return t
}

var allSQLCache map[Fields]string

// All reads each text message from the database, in descending chronological
// order.
func All(storer phys.Storer, fields Fields, fn func(*TextMessage)) {
	if fields&FLists != 0 {
		fields |= FID
	}
	if allSQLCache == nil {
		allSQLCache = make(map[Fields]string)
	}
	if _, ok := allSQLCache[fields&^FLists]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields&^FLists)
		sb.WriteString(" FROM textmsg t ORDER BY t.timestamp DESC")
		allSQLCache[fields&^FLists] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields&^FLists], func(stmt *phys.Stmt) {
		var t TextMessage
		for stmt.Step() {
			t.lists = t.lists[:0]
			t.Scan(stmt, fields&^FLists)
			if fields&FLists != 0 {
				t.readLists(storer)
			}
			fn(&t)
		}
	})
}
