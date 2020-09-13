package store

import (
	"time"

	"sunnyvaleserv.org/portal/model"
)

// CreateFolder creates a new folder in the database, assigning it the next
// available ID.
func (tx *Tx) CreateFolder(f *model.FolderNode) {
	tx.Tx.CreateFolder(f)
	tx.entry.Change("create folder [%d]", f.ID)
	if f.Parent != 0 {
		tx.entry.Change("set folder [%d] parent to %q [%d]", f.ID, tx.FetchFolder(f.Parent).Folder.Name, f.Parent)
	}
	tx.entry.Change("set folder [%d] name to %q", f.ID, f.Name)
	if f.Group != 0 {
		tx.entry.Change("set folder [%d] group to %q [%d]", f.ID, tx.auth.FetchGroup(f.Group).Name, f.Group)
	}
}

// WillUpdateFolder saves a copy of a folder's data prior to updating it, so
// that audit logs can be generated.
func (tx *Tx) WillUpdateFolder(f *model.FolderNode) {
	if tx.originalFolders[f.ID] == nil {
		of := *f.Folder
		of.Documents = make([]*model.Document, len(f.Documents))
		for i := range f.Documents {
			od := *f.Documents[i]
			of.Documents[i] = &od
		}
		tx.originalFolders[f.ID] = &of
	}
}

// UpdateFolder updates an existing Folder in the database.
func (tx *Tx) UpdateFolder(f *model.FolderNode) {
	var of = tx.originalFolders[f.ID]
	if of == nil {
		panic("must call WillUpdateFolder before calling UpdateFolder")
	}
	tx.Tx.UpdateFolder(f)
	if f.Name != of.Name {
		tx.entry.Change("set folder [%d] name to %q", f.ID, f.Name)
	}
	if f.Parent != of.Parent {
		if f.Parent != 0 {
			tx.entry.Change("set folder %q [%d] parent to %q [%d]", f.Name, f.ID, tx.FetchFolder(f.Parent).Folder.Name, f.Parent)
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
				if od.PostedBy != d.PostedBy {
					tx.entry.Change("set folder %q [%d] document %q [%d] postedBy to %q [%d]", f.Name, f.ID, d.Name, d.ID, tx.FetchPerson(d.PostedBy).InformalName, d.PostedBy)
				}
				if od.PostedAt != d.PostedAt {
					tx.entry.Change("set folder %q [%d] document %q [%d] postedAt to %s", f.Name, f.ID, d.Name, d.ID, d.PostedAt.In(time.Local).Format("2006-01-02 15:04:05"))
				}
				if od.NeedsApproval != d.NeedsApproval {
					if d.NeedsApproval {
						tx.entry.Change("set folder %q [%d] document %q [%d] needsApproval flag", f.Name, f.ID, d.Name, d.ID)
					} else {
						tx.entry.Change("clear folder %q [%d] document %q [%d] needsApproval flag", f.Name, f.ID, d.Name, d.ID)
					}
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
		tx.entry.Change("add folder %q [%d] document [%d]", f.Name, f.ID, d.ID)
		tx.entry.Change("set folder %q [%d] document [%d] name to %q", f.Name, f.ID, d.ID, d.Name)
		tx.entry.Change("set folder %q [%d] document %q [%d] postedBy to %q [%d]", f.Name, f.ID, d.Name, d.ID, tx.FetchPerson(d.PostedBy).InformalName, d.PostedBy)
		tx.entry.Change("set folder %q [%d] document %q [%d] postedAt to %s", f.Name, f.ID, d.Name, d.ID, d.PostedAt.In(time.Local).Format("2006-01-02 15:04:05"))
		if d.NeedsApproval {
			tx.entry.Change("set folder %q [%d] document %q [%d] needsApproval flag", f.Name, f.ID, d.Name, d.ID)
		}
	}
	if f.Approvals != of.Approvals {
		tx.entry.Change("set folder %q [%d] approvals to %d", f.Name, f.ID, f.Approvals)
	}
}

// DeleteFolder deletes a folder from the database.  This includes deleting all
// of its documents.
func (tx *Tx) DeleteFolder(f *model.FolderNode) {
	tx.Tx.DeleteFolder(f)
	tx.entry.Change("delete folder %q [%d]", f.Name, f.ID)
}
