package store

import (
	"sunnyvaleserv.org/portal/model"
)

// CreateFolder creates a new folder in the database, assigning it the next
// available ID.
func (tx *Tx) CreateFolder(f *model.Folder) {
	tx.Tx.CreateFolder(f)
	tx.entry.Change("create folder [%d]", f.ID)
	if f.Parent != 0 {
		tx.entry.Change("set folder [%d] parent to %q [%d]", f.ID, tx.FetchFolder(f.Parent).Name, f.Parent)
	}
	tx.entry.Change("set folder [%d] name to %q", f.ID, f.Name)
	if f.Group != 0 {
		tx.entry.Change("set folder [%d] group to %q [%d]", f.ID, tx.auth.FetchGroup(f.Group).Name, f.Group)
	}
}

// UpdateFolder updates an existing Folder in the database.
func (tx *Tx) UpdateFolder(f *model.Folder) {
	var of = tx.Tx.FetchFolder(f.ID)
	tx.Tx.UpdateFolder(f)
	if f.Name != of.Name {
		tx.entry.Change("set folder [%d] name to %q", f.ID, f.Name)
	}
	if f.Parent != of.Parent {
		if f.Parent != 0 {
			tx.entry.Change("set folder %q [%d] parent to %q [%d]", f.Name, f.ID, tx.FetchFolder(f.Parent).Name, f.Parent)
		} else {
			tx.entry.Change("remove folder %q [%d] parent", f.Name, f.ID)
		}
	}
	if f.Group != of.Group {
		if f.Group != 0 {
			tx.entry.Change("set folder %q [%d] group to %q [%d]", f.Name, f.ID, tx.auth.FetchGroup(f.Group).Name, f.Group)
		} else {
			tx.entry.Change("remove folder %q [%d] group", f.Name, f.ID)
		}
	}
DOCS1:
	for _, od := range of.Documents {
		for _, d := range f.Documents {
			if od.ID == d.ID {
				if od.Name != d.Name {
					tx.entry.Change("set folder %q [%d] document [%d] name to %q", f.Name, f.ID, d.ID, d.Name)
				}
				continue DOCS1
			}
		}
		tx.entry.Change("remove folder %q [%d] document %q [%d]", f.Name, f.ID, od.Name, od.ID)
	}
DOCS2:
	for _, d := range f.Documents {
		for _, od := range of.Documents {
			if od.ID == d.ID {
				continue DOCS2
			}
		}
		tx.entry.Change("add folder %q [%d] document %q [%d]", f.Name, f.ID, d.Name, d.ID)
	}
}

// DeleteFolder deletes a folder from the database.  This includes deleting all
// of its documents.
func (tx *Tx) DeleteFolder(f *model.Folder) {
	tx.Tx.DeleteFolder(f)
	tx.entry.Change("delete folder %q [%d]", f.Name, f.ID)
}
