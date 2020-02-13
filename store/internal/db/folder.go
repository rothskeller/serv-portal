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

// FetchFolders returns all of the folders in the database, sorted by name.
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
}

// CreateDocument creates a new document in the specified folder, with the
// specified ID and contents.
func (tx *Tx) CreateDocument(folder *model.Folder, document model.DocumentID, contents io.Reader) {
	var (
		fh  *os.File
		err error
	)
	fh, err = os.OpenFile(fmt.Sprintf("folders/%d/%d", folder.ID, document), os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666)
	panicOnError(err)
	_, err = io.Copy(fh, contents)
	panicOnError(err)
	panicOnError(fh.Close())
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
}

// DeleteDocument deletes a document from the specified folder.
func (tx *Tx) DeleteDocument(folder *model.Folder, document model.DocumentID) {
	panicOnError(os.Remove(fmt.Sprintf("folders/%d/%d", folder.ID, document)))
}

// DeleteFolder deletes a folder from the database.  This includes deleting all
// of its documents.
func (tx *Tx) DeleteFolder(f *model.Folder) {
	panicOnError(os.RemoveAll(fmt.Sprintf("folders/%d", f.ID)))
	panicOnNoRows(tx.tx.Exec(`DELETE FROM Folder WHERE id=?`, f.ID))
}
