package phys

import (
	"fmt"

	"sunnyvaleserv.org/portal/util/problem"
)

// Audit records an audit log entry, which is logged when the transaction
// commits.
func Audit(storer Storer, f string, a ...interface{}) {
	store := storer.AsStore()
	store.tx.audit = append(store.tx.audit, fmt.Sprintf(f, a...))
}

// Problems returns the store's problem logger, which the caller can use to log
// problems.
func (store *Store) Problems() *problem.List {
	if store.tx != nil {
		return &store.tx.Problems
	}
	return &store.logentry.Problems
}
