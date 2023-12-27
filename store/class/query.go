package class

import (
	"strings"
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

var withIDSQLCache map[Fields]string

// WithID returns the class with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (c *Class) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM class c WHERE c.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			c = new(Class)
			c.Scan(stmt, fields)
			c.id = id
			c.fields |= FID
		}
	})
	return c
}

var allSQLCache map[Fields]string

// All reads each class from the database, in descending date order.
func All(storer phys.Storer, fields Fields, fn func(*Class)) {
	if allSQLCache == nil {
		allSQLCache = make(map[Fields]string)
	}
	if _, ok := allSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM class c ORDER BY c.start DESC")
		allSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allSQLCache[fields], func(stmt *phys.Stmt) {
		var c Class
		for stmt.Step() {
			c.Scan(stmt, fields)
			fn(&c)
		}
	})
}

var allFutureSQLCache map[Fields]string

// AllFuture reads each future class of the specified type from the database, in
// ascending date order.
func AllFuture(storer phys.Storer, ctype Type, fields Fields, fn func(*Class)) {
	if allFutureSQLCache == nil {
		allFutureSQLCache = make(map[Fields]string)
	}
	if _, ok := allFutureSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM class c WHERE type=? AND c.start>=? ORDER BY c.start")
		allFutureSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allFutureSQLCache[fields], func(stmt *phys.Stmt) {
		var c Class
		stmt.BindInt(int(ctype))
		stmt.BindText(time.Now().Format("2006-01-02"))
		for stmt.Step() {
			c.Scan(stmt, fields)
			fn(&c)
		}
	})
}
