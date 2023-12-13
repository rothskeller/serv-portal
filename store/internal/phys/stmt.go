package phys

import (
	"log"
	"sync/atomic"

	"zombiezen.com/go/sqlite"
)

var (
	trace           = false
	statementNumber uint64
)

// Stmt is a database statement.
type Stmt struct {
	conn   *sqlite.Conn
	stmt   *sqlite.Stmt
	number uint64
	row    int
	column int
}

// NewStmt creates a new statement.  The statement must be Reset when no longer
// needed.  (In most cases, it's better to call SQL or Exec rather than
// NewStmt.)
func NewStmt(storer Storer, sql string) *Stmt {
	var (
		stmt Stmt
		err  error
	)
	if trace {
		stmt.number = atomic.AddUint64(&statementNumber, 1)
		log.Printf("%d = %q\n", stmt.number, sql)
	}
	if stmt.stmt, err = storer.AsStore().conn.Prepare(sql); err != nil {
		panic(err)
	}
	stmt.conn = storer.AsStore().conn
	return &stmt
}

// BindBool binds the next statement parameter to the specified bool value.
func (stmt *Stmt) BindBool(v bool) {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = %v\n", stmt.number, stmt.column, v)
	}
	stmt.stmt.BindBool(stmt.column, v)
}

// BindNullBool binds the next statement parameter to the specified bool value
// if it's true, or to NULL if it's false.
func (stmt *Stmt) BindNullBool(v bool) {
	if !v {
		stmt.BindNull()
	} else {
		stmt.BindBool(v)
	}
}

// BindFloat binds the next statement parameter to the specified float value.
func (stmt *Stmt) BindFloat(v float64) {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = %f\n", stmt.number, stmt.column, v)
	}
	stmt.stmt.BindFloat(stmt.column, v)
}

// BindNullFloat binds the next statement parameter to the specified float value
// if it's non-zero, or to NULL if it's zero.
func (stmt *Stmt) BindNullFloat(v float64) {
	if v == 0 {
		stmt.BindNull()
	} else {
		stmt.BindFloat(v)
	}
}

// BindHexInt binds the next statement parameter to the specified int value.  If
// tracing is enabled, the value is formatted in hexadecimal.
func (stmt *Stmt) BindHexInt(v int) {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = 0x%x\n", stmt.number, stmt.column, v)
	}
	stmt.stmt.BindInt64(stmt.column, int64(v))
}

// BindNullHexInt binds the next statement parameter to the specified int value
// if it's non-zero, or to NULL if it's zero.  If tracing is enabled, the value
// is formatted in hexadecimal.
func (stmt *Stmt) BindNullHexInt(v int) {
	if v == 0 {
		stmt.BindNull()
	} else {
		stmt.BindHexInt(v)
	}
}

// BindInt binds the next statement parameter to the specified int value.
func (stmt *Stmt) BindInt(v int) {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = %d\n", stmt.number, stmt.column, v)
	}
	stmt.stmt.BindInt64(stmt.column, int64(v))
}

// BindNullInt binds the next statement parameter to the specified int value if
// it's non-zero, or to NULL if it's zero.
func (stmt *Stmt) BindNullInt(v int) {
	if v == 0 {
		stmt.BindNull()
	} else {
		stmt.BindInt(v)
	}
}

// BindText binds the next statement parameter to the specified text value.
func (stmt *Stmt) BindText(v string) {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = %q\n", stmt.number, stmt.column, v)
	}
	stmt.stmt.BindText(stmt.column, v)
}

// BindNullText binds the next statement parameter to the specified text value
// if it's non-empty, or to NULL if it's empty.
func (stmt *Stmt) BindNullText(v string) {
	if v == "" {
		stmt.BindNull()
	} else {
		stmt.BindText(v)
	}
}

// BindNull binds the next statement parameter to NULL.
func (stmt *Stmt) BindNull() {
	stmt.column++
	if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.stmt.BindNull(stmt.column)
}

// Step reads the next row returned by the statement, if there is one, and
// returns whether there was one.
func (stmt *Stmt) Step() (found bool) {
	var err error

	stmt.column = 0
	if found, err = stmt.stmt.Step(); err != nil {
		panic(err)
	}
	if trace {
		if found {
			stmt.row++
			log.Printf("%d ROW %d\n", stmt.number, stmt.row)
		} else {
			log.Printf("%d DONE\n", stmt.number)
		}
	}
	return found
}

// ColumnBool reads the next statement column as a bool value.  If the column is
// NULL, it returns false.
func (stmt *Stmt) ColumnBool() (v bool) {
	if !stmt.ColumnIsNull() {
		v = stmt.stmt.ColumnBool(stmt.column)
		if trace {
			log.Printf("  %d[%d] = %v\n", stmt.number, stmt.column, v)
		}
	} else if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.column++
	return v
}

// ColumnFloat reads the next statement column as a float value.  If the column
// is NULL, it returns 0.
func (stmt *Stmt) ColumnFloat() (v float64) {
	if !stmt.ColumnIsNull() {
		v = stmt.stmt.ColumnFloat(stmt.column)
		if trace {
			log.Printf("  %d[%d] = %f\n", stmt.number, stmt.column, v)
		}
	} else if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.column++
	return v
}

// ColumnHexInt reads the next statement column as an int value.  If the column
// is NULL, it returns 0.  If tracing is enabled, the value is formatted in
// hexadecimal.
func (stmt *Stmt) ColumnHexInt() (v int) {
	if !stmt.ColumnIsNull() {
		v = stmt.stmt.ColumnInt(stmt.column)
		if trace {
			log.Printf("  %d[%d] = 0x%x\n", stmt.number, stmt.column, v)
		}
	} else if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.column++
	return v
}

// ColumnInt reads the next statement column as an int value.  If the column is
// NULL, it returns 0.
func (stmt *Stmt) ColumnInt() (v int) {
	if !stmt.ColumnIsNull() {
		v = stmt.stmt.ColumnInt(stmt.column)
		if trace {
			log.Printf("  %d[%d] = %d\n", stmt.number, stmt.column, v)
		}
	} else if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.column++
	return v
}

// ColumnText reads the next statement column as a string value.  If the column
// is NULL, it returns an empty string.
func (stmt *Stmt) ColumnText() (v string) {
	if !stmt.ColumnIsNull() {
		v = stmt.stmt.ColumnText(stmt.column)
		if trace {
			log.Printf("  %d[%d] = %q\n", stmt.number, stmt.column, v)
		}
	} else if trace {
		log.Printf("  %d[%d] = NULL\n", stmt.number, stmt.column)
	}
	stmt.column++
	return v
}

// ColumnIsNull returns whether the next statement column has a NULL value.  It
// does not consume the column; the next ColumnXXX method will read it.
func (stmt *Stmt) ColumnIsNull() bool {
	return stmt.stmt.ColumnType(stmt.column) == sqlite.TypeNull
}

// Reset resets the statement so that it can be run again.  It does not clear
// out the prior bound parameters, but it does reset the sequence, so if any
// BindXXX calls are going to be made after Reset, they all should be.
func (stmt *Stmt) Reset() {
	if trace {
		log.Printf("%d RESET [ra %d] [liid %d]\n", stmt.number, stmt.conn.Changes(), stmt.conn.LastInsertRowID())
	}
	stmt.column = 0
	stmt.row = 0
	if err := stmt.stmt.Reset(); err != nil {
		panic(err)
	}
}

// A Separator is a function that returns a null string when first called, and a
// separator string for every subsequent call.
type Separator func() string

// NewSeparator returns a new Separator with the specified separator string.
func NewSeparator(sep string) Separator {
	var first = true

	return func() string {
		if first {
			first = false
			return ""
		}
		return sep
	}
}
