package store

import (
	"io"

	"sunnyvaleserv.org/portal/model"
)

// CreateFolder creates a new folder.
func (tx *Tx) CreateFolder(folder *model.Folder) {
	tx.Tx.CreateFolder(folder)
	tx.entry.Change("create folder %s name %q visibility %s %s", folder.URL, folder.Name, folder.Visibility, folder.Org)
}

// UpdateFolder updates the existing folder at the specified path, to have the
// specified details.
func (tx *Tx) UpdateFolder(from, to *model.Folder) {
	tx.Tx.UpdateFolder(from, to)
	if from.Name != to.Name {
		tx.entry.Change("set folder %s name to %q", to.URL, to.Name)
	}
	if from.Visibility != to.Visibility || from.Org != to.Org {
		tx.entry.Change("set folder %s visibility to %s %s", to.URL, to.Visibility, to.Org)
	}
}

// DeleteFolder deletes an existing folder.
func (tx *Tx) DeleteFolder(folder *model.Folder) {
	tx.Tx.DeleteFolder(folder)
	tx.entry.Change("delete folder %s", folder.URL)
}

// CreateLink creates a link document.
func (tx *Tx) CreateLink(folder *model.Folder, link *model.Document) {
	tx.Tx.CreateLink(folder, link)
	tx.entry.Change("add folder %s link %q to url %q", folder.URL, link.Name, link.URL)
}

// UpdateLink updates the existing link document at the specified name.
func (tx *Tx) UpdateLink(fromf, tof *model.Folder, from, to *model.Document) {
	tx.Tx.UpdateLink(fromf, tof, from, to)
	if fromf.URL != tof.URL {
		tx.entry.Change("move link %q from folder %s to folder %s", from.Name, fromf.URL, tof.URL)
	}
	if from.Name != to.Name {
		tx.entry.Change("rename folder %s link %q to %q", tof.URL, from.Name, to.Name)
	}
	if from.URL != to.URL {
		tx.entry.Change("set folder %s link %q url to %q", tof.URL, to.Name, to.URL)
	}
}

// DeleteLink deletes an existing link document.
func (tx *Tx) DeleteLink(folder *model.Folder, link *model.Document) {
	tx.Tx.DeleteLink(folder, link)
	tx.entry.Change("delete folder %s link %q", folder.URL, link.Name)
}

// CreateFile creates a file document.
func (tx *Tx) CreateFile(folder *model.Folder, file *model.Document, contents io.Reader) {
	tx.Tx.CreateFile(folder, file, contents)
	tx.entry.Change("add folder %s file %q", folder.URL, file.Name)
}

// UpdateFile updates the existing file document at the specified name.
func (tx *Tx) UpdateFile(fromf, tof *model.Folder, from, to *model.Document, contents io.Reader) {
	tx.Tx.UpdateFile(fromf, tof, from, to, contents)
	if fromf.URL != tof.URL {
		tx.entry.Change("move file %q from folder %s to folder %s", from.Name, fromf.URL, tof.URL)
	}
	if from.Name != to.Name {
		tx.entry.Change("rename folder %s file %q to %q", tof.URL, from.Name, to.Name)
	}
}

// DeleteFile deletes an existing file document.
func (tx *Tx) DeleteFile(folder *model.Folder, file *model.Document) {
	tx.Tx.DeleteFile(folder, file)
	tx.entry.Change("delete folder %s file %q", folder.URL, file.Name)
}
