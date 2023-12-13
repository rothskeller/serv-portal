package search

import (
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/venue"
)

// EmptyEntireIndex empties the entire search index.  It's for use by the
// convert program.
func EmptyEntireIndex() {
	phys.EmptyEntireIndex()
}

// RebuildSearchIndex rebuilds the entire search index.  It's used by the
// helper program of the same name.
func RebuildSearchIndex(storer phys.Storer) {
	phys.EmptyEntireIndex()
	document.IndexAll(storer)
	event.IndexAll(storer)
	folder.IndexAll(storer)
	person.IndexAll(storer)
	role.IndexAll(storer)
	venue.IndexAll(storer)
}
