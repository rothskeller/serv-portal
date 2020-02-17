package cache

import (
	"io"
	"os"
	"sort"

	"sunnyvaleserv.org/portal/model"
)

// cacheFolders reads all Folders and builds the resulting tree.
func (tx *Tx) cacheFolders() {
	if tx.folders != nil {
		return
	}
	tx.folders = make(map[model.FolderID]*model.FolderNode)
	tx.rootFolder = &model.FolderNode{Folder: &model.Folder{Name: "Files"}}
	tx.folders[0] = tx.rootFolder
	for _, f := range tx.Tx.FetchFolders() {
		var fn = model.FolderNode{Folder: f}
		tx.folders[f.ID] = &fn
	}
	tx.relinkFolderCache()
}
func (tx *Tx) relinkFolderCache() {
	for _, f := range tx.folders {
		f.ChildNodes = nil
	}
	for _, f := range tx.folders {
		if f != tx.rootFolder {
			f.ParentNode = tx.folders[f.Parent]
			f.ParentNode.ChildNodes = append(f.ParentNode.ChildNodes, f)
		}
	}
	for _, f := range tx.folders {
		sort.Sort(model.FolderNodeSort(f.ChildNodes))
	}
}

// FetchFolder retrieves a single FolderNode from the database by ID.  It
// returns nil if no such Folder exists.
func (tx *Tx) FetchFolder(id model.FolderID) (f *model.FolderNode) {
	tx.cacheFolders()
	return tx.folders[id]
}

// FetchRootFolder returns the root folder.
func (tx *Tx) FetchRootFolder() (f *model.FolderNode) {
	tx.cacheFolders()
	return tx.rootFolder
}

// CreateFolder creates a new folder in the database, assigning it the next
// available ID.
func (tx *Tx) CreateFolder(f *model.FolderNode) {
	tx.cacheFolders()
	tx.Tx.CreateFolder(f.Folder)
	tx.folders[f.ID] = f
	f.ParentNode = tx.folders[f.Parent]
	f.ParentNode.ChildNodes = append(f.ParentNode.ChildNodes, f)
	sort.Sort(model.FolderNodeSort(f.ParentNode.ChildNodes))
}

// UpdateFolder updates an existing Folder in the database.
func (tx *Tx) UpdateFolder(f *model.FolderNode) {
	tx.cacheFolders()
	tx.Tx.UpdateFolder(f.Folder)
	tx.relinkFolderCache()
	if f.ParentNode != nil {
		sort.Sort(model.FolderNodeSort(f.ParentNode.ChildNodes))
	}
}

// DeleteFolder deletes a folder from the database.  This includes deleting all
// of its documents.
func (tx *Tx) DeleteFolder(f *model.FolderNode) {
	tx.Tx.DeleteFolder(f.Folder)
	delete(tx.folders, f.ID)
	tx.relinkFolderCache()
}

// FetchDocument returns an open file handle to the specified document in the
// specified folder.  The caller should close the file handle when finished.
func (tx *Tx) FetchDocument(folder *model.FolderNode, document model.DocumentID) (fh *os.File) {
	return tx.Tx.FetchDocument(folder.Folder, document)
}

// CreateDocument creates a new document in the specified folder, with the
// specified ID and contents.
func (tx *Tx) CreateDocument(folder *model.FolderNode, document model.DocumentID, contents io.Reader) {
	tx.Tx.CreateDocument(folder.Folder, document, contents)
}

// DeleteDocument deletes a document from the specified folder.
func (tx *Tx) DeleteDocument(folder *model.FolderNode, document model.DocumentID) {
	tx.Tx.DeleteDocument(folder.Folder, document)
}
