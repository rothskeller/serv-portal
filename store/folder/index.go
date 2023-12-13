package folder

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// IndexKey returns the index key for the folder.
func (f *Folder) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("F%d", f.ID())
}

// IndexEntry returns a complete index entry for the folder.
func (f *Folder) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	return &phys.IndexEntry{
		Key:   f.IndexKey(nil),
		Type:  "Folder",
		Label: f.Name(),
		Name:  f.Name(),
	}
}

// IndexAll indexes all people.
func IndexAll(storer phys.Storer) {
	indexAllWithParent(storer, RootID)
}
func indexAllWithParent(storer phys.Storer, parent ID) {
	AllWithParent(storer, parent, FID|FName, func(f *Folder) {
		phys.Index(storer, f)
		indexAllWithParent(storer, f.ID())
	})
}
