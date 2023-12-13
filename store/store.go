package store

import (
	"context"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/util/log"
)

// Store is a handle for access to the data store.  Each Store handle can be
// used only in a single goroutine.
type Store = phys.Store

// Storer is an interface to anything that has access to a Store.
type Storer = phys.Storer

// Connect connects to the physical data store and runs the supplied function in
// it.  The function must not pass the connection handle to any other goroutine.
func Connect(ctx context.Context, logentry *log.Entry, fn func(*Store)) error {
	return phys.Connect(ctx, logentry, fn)
}
