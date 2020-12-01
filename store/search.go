package store

import (
	"sunnyvaleserv.org/portal/model"
)

// Search executes a search and returns the matched objects.  It returns an
// error if the search string syntax is invalid.
func (tx *Tx) Search(query string, handler func(interface{}) bool) (err error) {
	err = tx.Tx.Search(query, func(typ string, id int, id2 string) bool {
		switch typ {
		case "document":
			var fd FolderAndDocument

			fd.Folder, fd.Document = tx.FetchDocument(id2)
			return handler(fd)
		case "event":
			return handler(tx.FetchEvent(model.EventID(id)))
		case "folder":
			return handler(tx.FetchFolder(id2))
		case "person":
			return handler(tx.FetchPerson(model.PersonID(id)))
		case "role":
			return handler(tx.FetchRole(model.Role2ID(id)))
		case "textMessage":
			return handler(tx.FetchTextMessage(model.TextMessageID(id)))
		default:
			panic("unexpected search entry type: " + typ)
		}
	})
	return err
}

// FolderAndDocument is the object returned by Search when the search matches a
// document in a folder.
type FolderAndDocument struct {
	Folder   *model.Folder
	Document *model.Document
}

// RebuildSearchIndex rebuilds the entire search index.
func (tx *Tx) RebuildSearchIndex() {
	tx.Tx.RebuildSearchIndex(tx.Authorizer().FetchGroups(tx.Authorizer().AllGroups()))
}
