package db

// Folders and documents are stored in the file system, as follows:
//
// Each folder is stored as a directory underneath data/folders, with the
// directory location modeling the folder hierarchy.  data/folders is the root
// folder.  The directory names are a lower-kebab-case version of the folder
// name.  Each such directory contains a file ".folder.json" with metadata about
// the folder: its original name and its visibility.
//
// Files are stored under their own names as plain files in the directory
// corresponding to the folder that contains them.  Files have mode 0666 (or
// possibly 0644 depending on umask).
//
// Links are stored as plain files in the directory corresponding to the folder
// that contains them.  The name of the file is the link title, and the contents
// of the file are the link URL (with a trailing newline).  Links have mode
// 0600; the mode is how they are distinguished from files.
//
// File and link names cannot start or end with a dot or a space, or contain any
// of a variety of characters that are unsafe on either Unix or Windows.
//
// The code in this module will not allow silent overwriting of files via
// rename.  It is up to higher level code to remove files prior to attempting to
// overwrite them.
import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/mailru/easyjson/jlexer"
	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
)

const foldersRoot = "folders"
const folderMetadataFile = "/.folder.json"

// FetchFolder retrieves the folder at the specified path, or nil if there is no
// folder at that path.  Use an empty string for the path to get the root
// folder.
func (tx *Tx) FetchFolder(path string) *model.Folder {
	var (
		in     jlexer.Lexer
		err    error
		folder = model.Folder{URL: path}
	)
	if in.Data, err = ioutil.ReadFile(foldersRoot + path + folderMetadataFile); err != nil {
		return nil
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeString()
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "name":
			folder.Name = in.String()
		case "visibility":
			folder.Visibility, err = model.ParseFolderVisibility(in.String())
			in.AddError(err)
		case "org":
			folder.Org, err = model.ParseOrg(in.String())
			in.AddError(err)
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	in.Consumed()
	if err = in.Error(); err != nil {
		panic(err)
	}
	return &folder
}

// FetchFolders retrieves the list of all folders.
func (tx *Tx) FetchFolders() (list []*model.Folder) {
	filepath.Walk(foldersRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil || !info.IsDir() {
			return nil
		}
		list = append(list, tx.FetchFolder(path[7:]))
		return nil
	})
	return list
}

// FetchSubFolders retrieves the list of subfolders of the specified folder.
func (tx *Tx) FetchSubFolders(parent *model.Folder) (list []*model.Folder) {
	var (
		fh       *os.File
		children []os.FileInfo
		err      error
	)
	if fh, err = os.Open(foldersRoot + parent.URL); err != nil {
		panic(err)
	}
	if children, err = fh.Readdir(0); err != nil {
		panic(err)
	}
	fh.Close()
	for _, child := range children {
		if child.IsDir() {
			if cf := tx.FetchFolder(parent.URL + "/" + child.Name()); cf == nil {
				panic(parent.URL + "/" + child.Name())
			} else {
				list = append(list, cf)
			}
		}
	}
	sort.Sort(model.FolderSort(list))
	return list
}

// FetchParentFolder retrieves the parent of the specified folder.  It returns
// nil if the specified folder is the root.
func (tx *Tx) FetchParentFolder(child *model.Folder) *model.Folder {
	if child.URL == "" {
		return nil
	}
	var parentPath = filepath.Dir(child.URL)
	if parentPath == "/" {
		parentPath = ""
	}
	return tx.FetchFolder(parentPath)
}

// FetchDocuments retrieves the list of documents in the specified folder.
func (tx *Tx) FetchDocuments(folder *model.Folder) (list []*model.Document) {
	var (
		fh       *os.File
		children []os.FileInfo
		err      error
	)
	if fh, err = os.Open(foldersRoot + folder.URL); err != nil {
		panic(err)
	}
	if children, err = fh.Readdir(0); err != nil {
		panic(err)
	}
	fh.Close()
	for _, child := range children {
		var (
			fname string
			doc   model.Document
			stat  os.FileInfo
		)
		if child.IsDir() {
			continue
		}
		if child.Name()[0] == '.' {
			continue
		}
		doc.Name = child.Name()
		fname = foldersRoot + folder.URL + "/" + doc.Name
		if stat, err = os.Stat(fname); err != nil {
			panic(err)
		}
		if stat.Mode()&0004 == 0 { // Link
			var by []byte
			if by, err = ioutil.ReadFile(fname); err != nil {
				panic(err)
			}
			doc.URL = strings.TrimSpace(string(by))
		}
		list = append(list, &doc)
	}
	sort.Sort(model.DocumentSort(list))
	return list
}

// PathExists returns whether a specified pathname exists, and its type.
func (tx *Tx) PathExists(path string) (exists, folder bool) {
	var (
		stat os.FileInfo
		err  error
	)
	if stat, err = os.Stat(foldersRoot + path); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err != nil {
		return false, false
	}
	return true, stat.IsDir()
}

// FetchFile opens a handle to the specified file.  It returns nil if there is
// no such file.
func (tx *Tx) FetchFile(path string) *os.File {
	var (
		fh   *os.File
		stat os.FileInfo
		err  error
	)
	if fh, err = os.Open(foldersRoot + path); err != nil && !os.IsNotExist(err) {
		panic(err)
	}
	if err != nil {
		return nil
	}
	if stat, err = fh.Stat(); err != nil {
		panic(err)
	}
	if stat.IsDir() || stat.Mode()&0004 == 0 {
		fh.Close()
		return nil
	}
	return fh
}

// FetchDocument returns the folder and document identified by the specified
// path, or nils if there is no such document.
func (tx *Tx) FetchDocument(path string) (folder *model.Folder, document *model.Document) {
	var (
		stat os.FileInfo
		by   []byte
		err  error
	)
	if folder = tx.FetchFolder(filepath.Dir(path)); folder == nil {
		return nil, nil
	}
	if stat, err = os.Stat(foldersRoot + path); err != nil || stat.IsDir() {
		return nil, nil
	}
	document = &model.Document{Name: filepath.Base(path)}
	if stat.Mode()&0004 != 0 {
		return
	}
	if by, err = ioutil.ReadFile(foldersRoot + path); err != nil {
		panic(err)
	}
	document.URL = strings.TrimSpace(string(by))
	return
}

// CreateFolder creates a new folder.
func (tx *Tx) CreateFolder(folder *model.Folder) {
	if err := os.Mkdir(foldersRoot+folder.URL, 0777); err != nil {
		panic(err)
	}
	tx.writeFolderMetadata(folder)
	tx.indexFolder(folder)
}

// UpdateFolder updates the existing folder at the specified path, to have the
// specified details.
func (tx *Tx) UpdateFolder(from, to *model.Folder) {
	tx.unindexFolder(from)
	if from.URL != to.URL {
		if _, err := os.Stat(foldersRoot + to.URL); err == nil {
			panic("destination exists")
		}
		if err := os.Rename(foldersRoot+from.URL, foldersRoot+to.URL); err != nil {
			panic(err)
		}
	}
	tx.writeFolderMetadata(to)
	tx.indexFolder(to)
}

func (tx *Tx) writeFolderMetadata(folder *model.Folder) {
	var (
		fh  *os.File
		out jwriter.Writer
		err error
	)
	out.RawString(`{"name":`)
	out.String(folder.Name)
	out.RawString(`,"visibility":`)
	out.String(folder.Visibility.String())
	if folder.Visibility == model.FolderVisibleToOrg {
		out.RawString(`,"org":`)
		out.String(folder.Org.String())
	}
	out.RawString("}\n")
	if fh, err = os.Create(foldersRoot + folder.URL + folderMetadataFile); err != nil {
		panic(err)
	}
	if _, err = out.DumpTo(fh); err != nil {
		panic(err)
	}
	if err = fh.Close(); err != nil {
		panic(err)
	}
}

// DeleteFolder deletes an existing folder.
func (tx *Tx) DeleteFolder(folder *model.Folder) {
	tx.unindexFolder(folder)
	if err := os.RemoveAll(foldersRoot + folder.URL); err != nil {
		panic(err)
	}
}

// CreateLink creates a link document.
func (tx *Tx) CreateLink(folder *model.Folder, link *model.Document) {
	var (
		fh  *os.File
		err error
	)
	if fh, err = os.OpenFile(foldersRoot+folder.URL+"/"+link.Name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err != nil {
		panic(err)
	}
	if _, err = fmt.Fprintln(fh, link.URL); err != nil {
		panic(err)
	}
	if err = fh.Close(); err != nil {
		panic(err)
	}
	tx.indexDocument(folder, link)
}

// UpdateLink updates the existing link document at the specified name.
func (tx *Tx) UpdateLink(fromf, tof *model.Folder, from, to *model.Document) {
	var (
		fh    *os.File
		stat  os.FileInfo
		err   error
		fname = foldersRoot + tof.URL + "/" + to.Name
	)
	tx.unindexDocument(fromf, from)
	if fromf.URL != tof.URL || from.Name != to.Name {
		if fh, err = os.OpenFile(fname, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0600); err != nil {
			panic(err)
		}
		if err = os.Remove(foldersRoot + fromf.URL + "/" + from.Name); err != nil {
			panic(err)
		}
	} else {
		if stat, err = os.Stat(fname); err != nil {
			panic(err)
		} else if stat.IsDir() || stat.Mode()&0004 != 0 {
			panic("not a link")
		}
		if fh, err = os.OpenFile(fname, os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
			panic(err)
		}
	}
	if _, err = fmt.Fprintln(fh, to.URL); err != nil {
		panic(err)
	}
	if err = fh.Close(); err != nil {
		panic(err)
	}
	tx.indexDocument(tof, to)
}

// DeleteLink deletes an existing link document.
func (tx *Tx) DeleteLink(folder *model.Folder, link *model.Document) {
	tx.unindexDocument(folder, link)
	if err := os.Remove(foldersRoot + folder.URL + "/" + link.Name); err != nil {
		panic(err)
	}
}

// CreateFile creates a file document.
func (tx *Tx) CreateFile(folder *model.Folder, file *model.Document, contents io.Reader) {
	var (
		fh  *os.File
		err error
	)
	if fh, err = os.OpenFile(foldersRoot+folder.URL+"/"+file.Name, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0666); err != nil {
		panic(err)
	}
	if _, err = io.Copy(fh, contents); err != nil {
		panic(err)
	}
	if err = fh.Close(); err != nil {
		panic(err)
	}
	tx.indexDocument(folder, file)
}

// UpdateFile updates the existing file document at the specified name.
func (tx *Tx) UpdateFile(fromf, tof *model.Folder, from, to *model.Document, contents io.Reader) {
	var (
		fh    *os.File
		stat  os.FileInfo
		err   error
		fname = foldersRoot + tof.URL + "/" + to.Name
	)
	tx.unindexDocument(fromf, from)
	if fromf.URL != tof.URL || from.Name != to.Name {
		if _, err = os.Stat(fname); err == nil {
			panic("destination exists")
		}
		if err = os.Rename(foldersRoot+fromf.URL+"/"+from.Name, fname); err != nil {
			panic(err)
		}
	}
	if contents != nil {
		if stat, err = os.Stat(fname); err != nil {
			panic(err)
		} else if stat.IsDir() || stat.Mode()&0004 == 0 {
			panic("not a file")
		}
		if fh, err = os.OpenFile(fname, os.O_TRUNC|os.O_WRONLY, 0600); err != nil {
			panic(err)
		}
		if _, err = io.Copy(fh, contents); err != nil {
			panic(err)
		}
		if err = fh.Close(); err != nil {
			panic(err)
		}
	}
	tx.indexDocument(tof, to)
}

// DeleteFile deletes an existing file document.
func (tx *Tx) DeleteFile(folder *model.Folder, file *model.Document) {
	tx.unindexDocument(folder, file)
	if err := os.Remove(foldersRoot + folder.URL + "/" + file.Name); err != nil {
		panic(err)
	}
}

// indexFolder adds information to the index about a folder and its contents.
func (tx *Tx) indexFolder(folder *model.Folder) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id2, folderName) VALUES ('folder',?,?)`, folder.URL, folder.Name))
	for _, sf := range tx.FetchSubFolders(folder) {
		tx.indexFolder(sf)
	}
	for _, d := range tx.FetchDocuments(folder) {
		tx.indexDocument(folder, d)
	}
}

// unindexFolder removes information from the index about a folder and its
// contents.
func (tx *Tx) unindexFolder(folder *model.Folder) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='folder' AND id2=?`, folder.URL))
	for _, sf := range tx.FetchSubFolders(folder) {
		tx.unindexFolder(sf)
	}
	for _, d := range tx.FetchDocuments(folder) {
		tx.unindexDocument(folder, d)
	}
}

// indexDocument updates the search index with information about a document.
func (tx *Tx) indexDocument(folder *model.Folder, document *model.Document) {
	panicOnExecError(tx.tx.Exec(`INSERT INTO search (type, id2, documentName, documentContents) VALUES ('document',?,?,?)`, folder.URL+"/"+document.Name, document.Name, document.URL))
}

// unindexDocument updates the search index with information about a document.
func (tx *Tx) unindexDocument(folder *model.Folder, document *model.Document) {
	panicOnExecError(tx.tx.Exec(`DELETE FROM search WHERE type='document' AND id2=?`, folder.URL+"/"+document.Name))
}
