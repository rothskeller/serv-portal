package db

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"sort"

	"sunnyvaleserv.org/portal/model"
)

// FetchFolder retrieves a single Folder from the database by ID.  It returns
// nil if no such Folder exists.
func (tx *Tx) FetchFolder(id model.FolderID) (f *model.Folder) {
	var data []byte
	f = new(model.Folder)
	switch err := tx.tx.QueryRow(`SELECT data FROM folder WHERE id=?`, id).Scan(&data); err {
	case nil:
		panicOnError(f.Unmarshal(data))
		return f
	case sql.ErrNoRows:
		return nil
	default:
		panic(err)
	}
}

// FetchFolders returns all of the folders in the database, sorted by order and
// name.
func (tx *Tx) FetchFolders() (folders []*model.Folder) {
	var (
		rows *sql.Rows
		err  error
	)
	rows, err = tx.tx.Query(`SELECT data FROM folder`)
	panicOnError(err)
	for rows.Next() {
		var data []byte
		var f model.Folder
		panicOnError(rows.Scan(&data))
		panicOnError(f.Unmarshal(data))
		folders = append(folders, &f)
	}
	panicOnError(rows.Err())
	sort.Sort(model.FolderSort(folders))
	return folders
}

// FetchDocument returns an open file handle to the specified document in the
// specified folder.  The caller should close the file handle when finished.
func (tx *Tx) FetchDocument(folder *model.Folder, document model.DocumentID) (fh *os.File) {
	var err error

	if fh, err = os.Open(fmt.Sprintf("folders/%d/%d", folder.ID, document)); err != nil {
		panic(err)
	}
	return fh
}

// CreateFolder creates a new folder in the database, assigning it the next
// available ID.
func (tx *Tx) CreateFolder(f *model.Folder) {
	var (
		data []byte
		err  error
	)
	panicOnError(tx.tx.QueryRow(`SELECT coalesce(max(id), 0) FROM folder`).Scan(&f.ID))
	f.ID++
	data, err = f.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`INSERT INTO folder (id, data) VALUES (?,?)`, f.ID, data))
	panicOnError(os.Mkdir(fmt.Sprintf("folders/%d", f.ID), 0777))
	tx.indexFolder(f, false)
}

// CreateDocument creates a new document in the specified folder, with the
// specified ID and contents.
func (tx *Tx) CreateDocument(folder *model.Folder, document *model.Document, contents io.Reader) {
	var (
		fh  *os.File
		err error
	)
	fh, err = os.OpenFile(fmt.Sprintf("folders/%d/%d", folder.ID, document.ID), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	panicOnError(err)
	_, err = io.Copy(fh, contents)
	panicOnError(err)
	panicOnError(fh.Close())
	tx.indexDocument(folder, document, false)
}

// UpdateFolder updates an existing Folder in the database.
func (tx *Tx) UpdateFolder(f *model.Folder) {
	var (
		data []byte
		err  error
	)
	data, err = f.Marshal()
	panicOnError(err)
	panicOnExecError(tx.tx.Exec(`UPDATE folder SET data=? WHERE id=?`, data, f.ID))
	tx.indexFolder(f, true)
}

// DeleteDocument deletes a document from the specified folder.
func (tx *Tx) DeleteDocument(folder *model.Folder, document model.DocumentID) {
	panicOnError(os.Remove(fmt.Sprintf("folders/%d/%d", folder.ID, document)))
	panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='document' AND id=? AND id2=?`, folder.ID, document))
}

// DeleteFolder deletes a folder from the database.  This includes deleting all
// of its documents.
func (tx *Tx) DeleteFolder(f *model.Folder) {
	panicOnError(os.RemoveAll(fmt.Sprintf("folders/%d", f.ID)))
	panicOnNoRows(tx.tx.Exec(`DELETE FROM folder WHERE id=?`, f.ID))
	panicOnNoRows(tx.tx.Exec(`DELETE FROM search WHERE type='folder' AND id=?`, f.ID))
	panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='document' AND id=?`, f.ID))
}

func (tx *Tx) indexFolder(f *model.Folder, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='folder' AND id=?`, f.ID))
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, folderName) VALUES ('folder',?,?)`, f.ID, f.Name))
}

func (tx *Tx) indexDocument(f *model.Folder, d *model.Document, replace bool) {
	if replace {
		panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='document' AND id=? and id2=?`, f.ID, d.ID))
	}
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id, id2, documentName) VALUES ('document',?,?,?)`, f.ID, d.ID, d.Name))
}
