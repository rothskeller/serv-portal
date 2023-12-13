package phys

import (
	"context"
	"fmt"
	"sync"

	"zombiezen.com/go/sqlite"

	"sunnyvaleserv.org/portal/util/config"
)

// poolSize is the maximum number of concurrent database connections.
const poolSize = 4

// openFlags is the flags controlling how database connections are opened.
const openFlags = sqlite.OpenReadWrite | sqlite.OpenNoMutex | sqlite.OpenSharedCache

// pool is the pool of open database connections.
var pool = make(chan *sqlite.Conn, poolSize)

// numOpened is the number of connections opened, access controlled by
// numOpenedMutex.
var numOpened int
var numOpenedMutex sync.Mutex

// dbconnect returns an open database connection.  It may wait until one becomes
// available if the maximum number of connections have been created.  The
// resulting conn will be canceled when the provided context is canceled.  It
// returns an error only if a connection cannot be created.
func dbconnect(ctx context.Context) (conn *sqlite.Conn, err error) {
	var stmt *sqlite.Stmt

	// If there's an open connection sitting in the free pool, grab it and
	// use it.
	select {
	case conn = <-pool:
		conn.SetInterrupt(ctx.Done())
		return conn, nil
	default:
		break
	}
	// If we've opened the maximum number of connections, wait for one to be
	// put into the free pool.
	numOpenedMutex.Lock()
	if numOpened == poolSize {
		numOpenedMutex.Unlock()
		conn = <-pool // will block
		conn.SetInterrupt(ctx.Done())
		return conn, nil
	}
	// Open a new connection and return it.
	numOpened++
	numOpenedMutex.Unlock()
	if conn, err = sqlite.OpenConn(config.Get("databaseFilename"), openFlags); err != nil {
		return nil, err
	}
	conn.SetInterrupt(ctx.Done())
	// Set the journal mode to truncate.
	if stmt, _, err = conn.PrepareTransient("PRAGMA journal_mode = TRUNCATE"); err != nil {
		conn.Close()
		return nil, err
	}
	if _, err = stmt.Step(); err != nil {
		conn.Close()
		return nil, err
	}
	if err = stmt.Finalize(); err != nil {
		conn.Close()
		return nil, err
	}
	// Turn on foreign key checking.
	if stmt, _, err = conn.PrepareTransient("PRAGMA foreign_keys = ON"); err != nil {
		conn.Close()
		return nil, err
	}
	if _, err = stmt.Step(); err != nil {
		conn.Close()
		return nil, err
	}
	if err = stmt.Finalize(); err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

// dbrelease releases a database connection back to the pool.
func dbrelease(conn *sqlite.Conn, withPanic bool) {
	if !withPanic && conn.CheckReset() != "" {
		panic("connection returned to pool with busy statement")
	}
	conn.SetInterrupt(nil)
	pool <- conn
}

// SQL prepares an SQL statement and passes it to the supplied function.  When
// the function returns, the statement is reset.  Any errors in preparing the
// statement cause a panic, since they are programming errors rather than
// runtime issues.  Note that SQL takes the store as an argument, not a
// receiver, which means it is accessible only within the store packages.
func SQL(storer Storer, sql string, fn func(*Stmt)) {
	var stmt = NewStmt(storer, sql)
	fn(stmt)
	stmt.Reset()
}

// Exec is a shortcut for a SQL statement that has no bindings and returns no
// rows.
func Exec(storer Storer, sql string) {
	var stmt = NewStmt(storer, sql)
	stmt.Step()
	stmt.Reset()
}

// RowsAffected returns the number of rows affected by the most recent
// statement.
func RowsAffected(storer Storer) int { return storer.AsStore().conn.Changes() }

// LastInsertRowID returns the row ID of the last row inserted.
func LastInsertRowID(storer Storer) int64 { return storer.AsStore().conn.LastInsertRowID() }

// createSavepoint creates a savepoint in the database.
func (store *Store) createSavepoint() (err error) {
	var stmt *sqlite.Stmt

	if stmt, err = store.conn.Prepare(fmt.Sprintf("SAVEPOINT x")); err != nil {
		return err
	}
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}

// releaseSavepoint releases a savepoint in the database.
func (store *Store) releaseSavepoint() (err error) {
	var stmt *sqlite.Stmt

	if stmt, err = store.conn.Prepare(fmt.Sprintf("RELEASE x")); err != nil {
		return err
	}
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}

// rollbackSavepoint rolls back a savepoint in the database.
func (store *Store) rollbackSavepoint() (err error) {
	var stmt *sqlite.Stmt

	if store.conn.AutocommitEnabled() {
		// Already rolled back, probably automatically because of
		// whatever error occurred.
		return nil
	}
	if stmt, err = store.conn.Prepare(fmt.Sprintf("ROLLBACK TO x")); err == nil {
		return err
	}
	if _, err = stmt.Step(); err != nil {
		return err
	}
	return stmt.Reset()
}
