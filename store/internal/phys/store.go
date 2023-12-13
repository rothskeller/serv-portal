// Package phys contains code that interacts directly with physical storage:
// the database, the file system (including audit log), and the index service.
package phys

import (
	"context"
	"fmt"
	"os"
	"runtime/debug"
	"sync"
	"time"

	"github.com/algolia/algoliasearch-client-go/v3/algolia/search"
	"zombiezen.com/go/sqlite"

	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/problem"
)

// opened ensures that we open the physical storage layer exactly once.
var opened sync.Once

// open opens the physical storage layer.
func open() {
	// Actually, the search index is the only thing that needs to be opened.
	openSearch()
}

// Store is a handle for access to the data store.  Each Store handle can be
// used only in a single goroutine.
type Store struct {
	conn     *sqlite.Conn
	logentry *log.Entry
	tx       *tx
}

// Connect connects to the data store and calls the supplied function with a
// handle to it.  The connection is released after the function returns (or
// panics).  The connection must be used only within a single goroutine.  If the
// maximum number of connections is already in use, Connect will block until one
// frees up.  If the context passed into Connect is canceled, the connection is
// canceled with it, and all subsequent actions on the connection will fail.
// Changed made through the connection are recorded in the supplied log entry.
func Connect(ctx context.Context, logentry *log.Entry, fn func(*Store)) (err error) {
	var store Store

	opened.Do(open)
	store.logentry = logentry
	if store.conn, err = dbconnect(ctx); err != nil {
		return err
	}
	defer func() {
		// This catches panics in fn.
		if panicked := recover(); panicked != nil {
			err = fmt.Errorf("PANIC: %v", panicked)
			store.logentry.Problems.AddError(err)
			store.logentry.Stack = debug.Stack()
		}
		dbrelease(store.conn, err != nil)
	}()
	fn(&store)
	return // err as set in defer function
}

// Storer is an interface for an object that can return a Store.
type Storer interface {
	AsStore() *Store
	Problems() *problem.List
}

// AsStore satisfies the Storer interface, returning itself.
func (store *Store) AsStore() *Store { return store }

// tx represents a physical layer transaction.
type tx struct {
	parent         *tx
	audit          []string
	Problems       problem.List
	removeOnFail   []string
	removeOnCommit []string
	searchOps      []search.BatchOperationIndexed
	nocommit       bool
}

// Transaction executes the supplied function in a transaction.  It may call the
// function multiple times, in the face of retriable errors, so the function
// must be idempotent.  If a non-retriable error occurs, it panics.  If the
// function calls the DoNotCommit method, the transaction rolls back rather than
// committing, but does not panic.
//
// Transactions can nest.  If an inner transaction succeeds, its changes become
// part of its parent transaction.  If an inner transaction rolls back due to a
// call to DoNotCommit, its changes are discarded but the outer transaction
// continues.  Note that if an inner transaction fails because of a panic,
// the panic will propagate upward and cause all parent transactions to fail.
func (store *Store) Transaction(fn func()) {
	var delay time.Duration = 8 * time.Millisecond
	var start time.Time = time.Now()
	const maxDelay = 1024 * time.Millisecond
	const maxWait = 30 * time.Second

	for {
		if store.attemptTransaction(fn) {
			return // success
		}
		if time.Since(start) > maxWait {
			panic("too many retries")
		}
		time.Sleep(delay)
		if delay < maxDelay {
			delay *= 2
		}
	}
}

// attemptTransaction is a single transaction attempt.  It returns true if the
// transaction succeeded, false if it failed with a retriable error.  It panics
// if the transaction fails with a non-retriable error.  Note that a rollback
// due to the no-commit flag is considered "success" in this context.
func (store *Store) attemptTransaction(fn func()) (success bool) {
	var tx tx

	// Set up the transaction and put it in the map.
	tx.parent = store.tx
	store.tx = &tx
	// Set up the transaction cleanup handler.
	defer func() {
		success = store.cleanupTransaction(recover())
	}()
	// Run the transaction.
	if err := store.createSavepoint(); err != nil {
		panic(err)
	}
	fn()
	return // return value is set in defer above.
}

// cleanup finalizes a transaction started with Transaction.
func (store *Store) cleanupTransaction(panicked interface{}) bool {
	// Roll back the transaction if it failed or its no-commit flag is set.
	if panicked != nil || store.tx.nocommit {
		goto ROLLBACK
	}
	if err := store.releaseSavepoint(); err != nil {
		panicked = err
		goto ROLLBACK
	}
	// We have successfully released the database savepoint.  If we have a
	// parent transaction, propagate the audit log entries, errors, file
	// removals, and search ops to it, and we're done.
	if tx := store.tx; tx.parent != nil {
		tx.parent.audit = append(tx.parent.audit, tx.audit...)
		tx.parent.Problems.AddList(&tx.Problems)
		tx.parent.removeOnFail = append(tx.parent.removeOnFail, tx.removeOnFail...)
		tx.parent.removeOnCommit = append(tx.parent.removeOnCommit, tx.removeOnCommit...)
		tx.parent.searchOps = append(tx.parent.searchOps, tx.searchOps...)
		store.tx = tx.parent
		return true
	}
	// We don't have a parent transaction, so releasing the savepoint
	// actually committed the database changes.  Now we commit the audit log
	// entries, errors, file removals, and search ops.
	store.logentry.Changes = append(store.logentry.Changes, store.tx.audit...)
	store.logentry.Problems.AddError(&store.tx.Problems)
	for _, r := range store.tx.removeOnCommit {
		os.RemoveAll(r)
	}
	if err := store.applySearchOps(); err != nil {
		store.tx = nil
		panic(err)
	}
	store.tx = nil
	return true

ROLLBACK:
	// Remove any files that are supposed to be removed if the transaction
	// fails.
	for _, r := range store.tx.removeOnFail {
		os.RemoveAll(r)
	}
	// Roll back the savepoint.  If we aren't already panicking and the
	// rollback fails, we'll panic with that cause.
	if err := store.rollbackSavepoint(); err != nil && panicked == nil {
		panicked = err
	}
	// If we rolled back due to the no-commit flag and had no other error,
	// return true.
	if panicked == nil {
		store.tx = store.tx.parent
		return true
	}
	// If we have a parent transaction, we always re-raise the panic.  Retry
	// only happens for the outermost transaction.
	if store.tx.parent != nil {
		store.tx = store.tx.parent
		panic(panicked)
	}
	store.tx = nil
	// If we rolled back due to a retriable error, return false.
	if err, ok := panicked.(error); ok {
		if code := sqlite.ErrCode(err).ToPrimary(); code == sqlite.ResultBusy || code == sqlite.ResultLocked {
			return false
		}
	}
	// In all other cases, re-raise the panic.
	panic(panicked)
}

// DoNotCommit marks the current transaction as invalid; when it completes, it
// will be rolled back instead of being committed.
func (store *Store) DoNotCommit() {
	store.tx.nocommit = true
}
