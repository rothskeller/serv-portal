// Package db contains the database access code for the SERV portal.
package db

import (
	"database/sql"
	"database/sql/driver"
	"fmt"
	"time"

	_ "github.com/scholacantorum/go-sqlite3" // SQLite driver
	"rothskeller.net/serv/model"
)

var dbh *sql.DB

// Open opens the database.
func Open(path string) {
	var (
		url string
		err error
	)
	url = "file:" + path + "?mode=rw&_busy_timeout=1000&_txlock=immediate&_foreign_keys=1&_journal_mode=TRUNCATE"
	if dbh, err = sql.Open("sqlite3", url); err != nil {
		panic(err)
	}
}

// Time is a wrapper around time.Time that stores in the database as a string.
// go-sqlite3 would do that for us, but it stores the timestamps with
// fractional seconds and a time zone indicator.  Ours stores them as integral
// seconds and no time zone indicator (local time assumed).  This makes the
// database easier to work with manually.  Note that when applying sqlite3 date
// and time functions to these columns, one must add the 'localtime' modifier to
// them, otherwise SQLite will assume they are in UTC.
type Time time.Time

// Value converts the time to a string, for storage into the database.
func (t Time) Value() (driver.Value, error) {
	if time.Time(t).IsZero() {
		return "", nil
	}
	return time.Time(t).In(time.Local).Format("2006-01-02 15:04:05"), nil
}

// Scan converts the string from the database into a Time.
func (t *Time) Scan(value interface{}) error {
	tt, ok := value.(string)
	if !ok {
		tb, ok := value.([]byte)
		if !ok {
			return fmt.Errorf("scanning %T into db.Time, should be string", value)
		}
		tt = string(tb)
	}
	if tt == "" {
		*t = Time(time.Time{})
		return nil
	}
	ft, err := time.ParseInLocation("2006-01-02 15:04:05", tt, time.Local)
	if err == nil {
		*t = Time(ft)
	}
	return err
}

// ID is a wrapper around int that stores 0 in the database as NULL.
type ID int

// Value converts the ID to database format.
func (id ID) Value() (driver.Value, error) {
	if id == 0 {
		return nil, nil
	}
	return int64(id), nil
}

// Scan converts the ID from database format.
func (id *ID) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*id = 0
	case int64:
		*id = ID(value)
	default:
		return fmt.Errorf("scanning %T into db.ID, should be int64 or nil", value)
	}
	return nil
}

// IDStr is a wrapper around string that stores "" in the database as NULL.
type IDStr string

// Value converts the ID to database format.
func (id IDStr) Value() (driver.Value, error) {
	if id == "" {
		return nil, nil
	}
	return string(id), nil
}

// Scan converts the ID from database format.
func (id *IDStr) Scan(value interface{}) error {
	switch value := value.(type) {
	case nil:
		*id = ""
	case string:
		*id = IDStr(value)
	case []byte:
		*id = IDStr(string(value))
	default:
		return fmt.Errorf("scanning %T into db.IDStr, should be string or nil", value)
	}
	return nil
}

// Tx is a wrapper around sql.Tx that implements all of the database package
// functions for access to our tables.
type Tx struct {
	tx        *sql.Tx
	roles     map[model.RoleID]*model.Role
	roleTags  map[model.RoleTag]*model.Role
	roleList  []*model.Role
	maxRoleID model.RoleID
	username  string
	request   string
}

// Begin starts a transaction, returning our Tx wrapper instead of a raw sql.Tx.
func Begin() (tx *Tx) {
	var err error

	tx = new(Tx)
	if tx.tx, err = dbh.Begin(); err != nil {
		panic(err)
	}
	tx.cacheRoles()
	return tx
}

// Commit commits a transaction.
func (tx *Tx) Commit() {
	panicOnError(tx.tx.Commit())
}

// Rollback rolls back a transaction.
func (tx *Tx) Rollback() error {
	if tx.tx != nil {
		return tx.tx.Rollback()
	}
	return nil
}

// SetUsername sets the username used in audit logging of database changes.
func (tx *Tx) SetUsername(username string) { tx.username = username }

// SetRequest sets the request used in audit logging of database changes.
func (tx *Tx) SetRequest(request string) { tx.request = request }

func panicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnExecError(_ sql.Result, err error) {
	if err != nil {
		panic(err)
	}
}

func panicOnNoRows(res sql.Result, err error) {
	var rows int64

	if err != nil {
		panic(err)
	}
	if rows, err = res.RowsAffected(); err != nil {
		panic(err)
	}
	if rows == 0 {
		panic("affected no rows")
	}
}

func lastInsertID(res sql.Result) int {
	nid, err := res.LastInsertId()
	panicOnError(err)
	return int(nid)
}
