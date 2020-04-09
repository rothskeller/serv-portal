package store

import (
	"sunnyvaleserv.org/portal/model"
)

// Search executes a search and returns the matched objects.  It returns an
// error if the search string syntax is invalid.
func (tx *Tx) Search(query string, handler func(interface{}) bool) (err error) {
	err = tx.Tx.Search(query, func(typ string, id, id2 int) bool {
		switch typ {
		case "document":
			var fd FolderAndDocument
			fd.Folder = tx.FetchFolder(model.FolderID(id))
			for _, d := range fd.Folder.Documents {
				if d.ID == model.DocumentID(id2) {
					fd.Document = d
					break
				}
			}
			return handler(fd)
		case "event":
			return handler(tx.FetchEvent(model.EventID(id)))
		case "folder":
			return handler(tx.FetchFolder(model.FolderID(id)))
		case "group":
			return handler(tx.Authorizer().FetchGroup(model.GroupID(id)))
		case "person":
			return handler(tx.FetchPerson(model.PersonID(id)))
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
	Folder   *model.FolderNode
	Document *model.Document
}

// RebuildSearchIndex rebuilds the entire search index.
func (tx *Tx) RebuildSearchIndex() {
	tx.Tx.RebuildSearchIndex(tx.Authorizer().FetchGroups(tx.Authorizer().AllGroups()))
}
