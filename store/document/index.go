package document

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
)

// IndexKey returns the index key for the document.
func (d *Document) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("D%d", d.ID)
}

// IndexEntry returns a complete index entry for the document.
func (d *Document) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	return &phys.IndexEntry{
		Key:   d.IndexKey(nil),
		Type:  "Document",
		Label: d.Name,
		Name:  d.Name,
	}
}

// IndexAll indexes all people.
func IndexAll(storer phys.Storer) {
	indexAllWithParent(storer, folder.RootID)
}
func indexAllWithParent(storer phys.Storer, parent folder.ID) {
	AllInFolder(storer, parent, func(d *Document) {
		phys.Index(storer, d)
	})
	folder.AllWithParent(storer, parent, folder.FID|folder.FName, func(f *folder.Folder) {
		indexAllWithParent(storer, f.ID())
	})
}
