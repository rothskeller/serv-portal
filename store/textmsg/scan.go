package textmsg

import (
	"strings"
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
)

const timestampFormat = "2006-01-02T15:04:05"

// ColumnList generates a comma-separated list of column names for the specified
// text message fields.  It is used in constructing SQL SELECT statements.
func ColumnList(sb *strings.Builder, fields Fields) {
	sep := phys.NewSeparator(", ")
	if fields&FID != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.id")
	}
	if fields&FSender != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.sender")
	}
	if fields&FTimestamp != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.timestamp")
	}
	if fields&FMessage != 0 {
		sb.WriteString(sep())
		sb.WriteString("t.message")
	}
	if fields&FLists != 0 {
		panic("TextMessage.Lists cannot be retrieved with ColumnList/Scan")
	}
}

// Scan reads columns corresponding to the specified fields from the specified
// statement into the receiver.
func (t *TextMessage) Scan(stmt *phys.Stmt, fields Fields) {
	if fields&FID != 0 {
		t.id = ID(stmt.ColumnInt())
	}
	if fields&FSender != 0 {
		t.sender = person.ID(stmt.ColumnInt())
	}
	if fields&FTimestamp != 0 {
		t.timestamp, _ = time.ParseInLocation(timestampFormat, stmt.ColumnText(), time.Local)
	}
	if fields&FMessage != 0 {
		t.message = stmt.ColumnText()
	}
	t.fields |= fields &^ FLists
}

const readListsSQL = `SELECT list, name FROM textmsg_list WHERE textmsg=?`

func (t *TextMessage) readLists(store phys.Storer) {
	phys.SQL(store, readListsSQL, func(stmt *phys.Stmt) {
		stmt.BindInt(int(t.ID()))
		for stmt.Step() {
			var tl TextToList
			tl.ID = list.ID(stmt.ColumnInt())
			tl.Name = stmt.ColumnText()
			t.lists = append(t.lists, tl)
		}
	})
	t.fields |= FLists
}
