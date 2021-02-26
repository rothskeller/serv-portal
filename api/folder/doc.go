package folder

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// downloadFile downloads the specified file.
func downloadFile(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	var (
		fh   *os.File
		stat os.FileInfo
		path = r.Path[3:]
	)
	if !CanViewFolder(r.Person, folder) {
		return util.Forbidden
	}
	if doc.URL != "" {
		return errors.New("attempt to download a link")
	}
	fh = r.Tx.FetchFile(path)
	r.Tx.Commit()
	if stat, err = fh.Stat(); err != nil {
		return err
	}
	if CanShowInBrowser(doc) {
		r.Header().Set("Content-Disposition", "inline")
	} else {
		r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.Name))
	}
	r.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(r, r.Request, doc.Name, stat.ModTime(), fh)
	return nil
}

// postNewFiles adds files to the specified folder.
func postNewFiles(r *util.Request, folder *model.Folder) (err error) {
	var (
		children []*model.Folder
		docs     []*model.Document
		replace  []*model.Document
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	if len(r.MultipartForm.File["file"]) == 0 {
		return errors.New("missing files")
	}
	children = r.Tx.FetchSubFolders(folder)
	docs = r.Tx.FetchDocuments(folder)
	for _, f := range r.MultipartForm.File["file"] {
		if f.Filename == "" {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not be empty.", f.Filename)
		}
		if f.Filename[0] == ' ' {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not start with a space.", f.Filename)
		}
		if f.Filename[0] == '.' {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not start with a dot.", f.Filename)
		}
		if f.Filename[len(f.Filename)-1] == ' ' {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not end with a space.", f.Filename)
		}
		if f.Filename[len(f.Filename)-1] == '.' {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not end with a dot.", f.Filename)
		}
		if strings.IndexAny(f.Filename, "\x00\x01\x02\x03\x04\x05\x06\x07\x08\x09\x0a\x0b\x0c\x0d\x0e\x0f\x10\x11\x12\x13\x14\x15\x16\x17\x18\x19\x1a\x1b\x1c\x1d\x1e\x1f\x7f<>:\"/\\|?*") >= 0 {
			return util.SendConflict(r, "The filename “%s” is not valid.  The filename must not contain < > : \" / \\ | ? * characters or unprintable characters.", f.Filename)
		}
		url := nameToURL(f.Filename)
		for _, child := range children {
			if nameToURL(child.Name) == url {
				return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", folder.Name, child.Name)
			}
		}
		for _, d := range docs {
			if d.Name == f.Filename {
				if d.URL != "" {
					return util.SendConflict(r, "The folder “%s” already contains a link named “%s”.", folder.Name, f.Filename)
				}
				replace = append(replace, d)
			}
		}
	}
	for _, rep := range replace {
		r.Tx.DeleteFile(folder, rep)
	}
	for _, f := range r.MultipartForm.File["file"] {
		var doc = model.Document{Name: f.Filename}
		var contents multipart.File
		if contents, err = f.Open(); err != nil {
			return err
		}
		r.Tx.CreateFile(folder, &doc, contents)
		contents.Close()
	}
	r.Tx.Commit()
	return nil
}

// postNewLink adds a link to the specified folder.
func postNewLink(r *util.Request, folder *model.Folder) (err error) {
	var (
		doc     model.Document
		url     string
		replace *model.Document
	)
	if !canEditFolder(r.Person, folder) || folder.URL == "" {
		return util.Forbidden
	}
	if doc.URL = r.FormValue("url"); doc.URL == "" {
		return errors.New("missing url")
	}
	if !strings.HasPrefix(doc.URL, "http://") && !strings.HasPrefix(doc.URL, "https://") {
		return errors.New("invalid URL")
	}
	if doc.Name = r.FormValue("name"); doc.Name == "" {
		doc.Name = getURLTitle(doc.URL)
	}
	if doc.Name[0] == '.' || strings.ContainsAny(doc.Name, "/:") {
		return errors.New("invalid name")
	}
	url = nameToURL(doc.Name)
	for _, child := range r.Tx.FetchSubFolders(folder) {
		if nameToURL(child.Name) == url {
			return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", folder.Name, child.Name)
		}
	}
	for _, d := range r.Tx.FetchDocuments(folder) {
		if d.Name == doc.Name {
			if d.URL == "" {
				return util.SendConflict(r, "The folder “%s” already contains a file named “%s”.", folder.Name, doc.Name)
			}
			if r.FormValue("replace") != "true" {
				r.Header().Set("X-Can-Replace", "true")
				return util.SendConflict(r, "The folder “%s” already contains a link named “%s”. Do you want to replace it?", folder.Name, doc.Name)
			}
			replace = d
		}
	}
	if replace != nil {
		r.Tx.DeleteLink(folder, replace)
	}
	r.Tx.CreateLink(folder, &doc)
	r.Tx.Commit()
	return nil
}

// getEditDocument returns the information needed when starting to edit the
// specified document.
func getEditDocument(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	var out jwriter.Writer

	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	out.RawString(`{"name":`)
	out.String(doc.Name)
	out.RawString(`,"url":`)
	out.String(doc.URL)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// postEditFile updates a file.
func postEditFile(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	var (
		newDoc   model.Document
		url      string
		contents multipart.File
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	if doc.URL != "" {
		return errors.New("document is not a file")
	}
	if newDoc.Name = r.FormValue("name"); newDoc.Name == "" {
		return errors.New("missing name")
	}
	if newDoc.Name[0] == '.' || strings.ContainsAny(newDoc.Name, "/:") {
		return errors.New("invalid name")
	}
	url = nameToURL(newDoc.Name)
	for _, child := range r.Tx.FetchSubFolders(folder) {
		if nameToURL(child.Name) == url {
			return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", folder.Name, child.Name)
		}
	}
	for _, d := range r.Tx.FetchDocuments(folder) {
		if d.Name != doc.Name && d.Name == newDoc.Name {
			if d.URL == "" {
				return util.SendConflict(r, "The folder “%s” already contains a file named “%s”.", folder.Name, doc.Name)
			}
			return util.SendConflict(r, "The folder “%s” already contains a link named “%s”.", folder.Name, doc.Name)
		}
	}
	switch len(r.MultipartForm.File["file"]) {
	case 0:
		break
	case 1:
		if contents, err = r.MultipartForm.File["file"][0].Open(); err != nil {
			return err
		}
		defer contents.Close()
	default:
		return errors.New("multiple files supplied")
	}
	r.Tx.UpdateFile(folder, folder, doc, &newDoc, contents)
	r.Tx.Commit()
	return nil
}

// postEditLink updates a link.
func postEditLink(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	var (
		newDoc model.Document
		url    string
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	if doc.URL == "" {
		return errors.New("document is not a link")
	}
	if newDoc.URL = r.FormValue("url"); newDoc.URL == "" {
		return errors.New("missing url")
	}
	if !strings.HasPrefix(newDoc.URL, "http://") && !strings.HasPrefix(newDoc.URL, "https://") {
		return errors.New("invalid URL")
	}
	if newDoc.Name = r.FormValue("name"); newDoc.Name == "" {
		return errors.New("missing name")
	}
	if newDoc.Name[0] == '.' || strings.ContainsAny(newDoc.Name, "/:") {
		return errors.New("invalid name")
	}
	url = nameToURL(newDoc.Name)
	for _, child := range r.Tx.FetchSubFolders(folder) {
		if nameToURL(child.Name) == url {
			return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", folder.Name, child.Name)
		}
	}
	for _, d := range r.Tx.FetchDocuments(folder) {
		if d != doc && d.Name == newDoc.Name {
			if d.URL == "" {
				return util.SendConflict(r, "The folder “%s” already contains a file named “%s”.", folder.Name, doc.Name)
			}
			return util.SendConflict(r, "The folder “%s” already contains a link named “%s”.", folder.Name, doc.Name)
		}
	}
	r.Tx.UpdateLink(folder, folder, doc, &newDoc)
	r.Tx.Commit()
	return nil
}

// postMoveDocument moves a document to a new folder.
func postMoveDocument(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	var (
		newParent *model.Folder
		url       string
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	if r.Form["parent"] == nil {
		return errors.New("missing parent")
	}
	if newParent = r.Tx.FetchFolder(r.FormValue("parent")); newParent == nil || newParent.URL == "" {
		return errors.New("invalid parent")
	}
	if !canEditFolder(r.Person, newParent) {
		return errors.New("forbidden parent")
	}
	url = nameToURL(doc.Name)
	for _, child := range r.Tx.FetchSubFolders(newParent) {
		if nameToURL(child.Name) == url {
			return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", newParent.Name, child.Name)
		}
	}
DOCS:
	for _, d := range r.Tx.FetchDocuments(newParent) {
		if d.Name == doc.Name {
			switch {
			case d.URL == "" && doc.URL == "":
				if r.FormValue("replace") == "true" {
					r.Tx.DeleteFile(newParent, d)
					break DOCS
				}
				r.Header().Set("X-Can-Replace", "true")
				return util.SendConflict(r, "The folder “%s” already contains a file named “%s”. Do you want to replace it?", newParent.Name, doc.Name)
			case d.URL == "" && doc.URL != "":
				return util.SendConflict(r, "The folder “%s” already contains a file named “%s”.", newParent.Name, doc.Name)
			case d.URL != "" && doc.URL == "":
				return util.SendConflict(r, "The folder “%s” already contains a link named “%s”.", newParent.Name, doc.Name)
			case d.URL != "" && doc.URL != "":
				if r.FormValue("replace") == "true" {
					r.Tx.DeleteLink(newParent, d)
					break DOCS
				}
				r.Header().Set("X-Can-Replace", "true")
				return util.SendConflict(r, "The folder “%s” already contains a link named “%s”. Do you want to replace it?", newParent.Name, doc.Name)
			}
		}
	}
	if doc.URL == "" {
		r.Tx.UpdateFile(folder, newParent, doc, doc, nil)
	} else {
		r.Tx.UpdateLink(folder, newParent, doc, doc)
	}
	r.Tx.Commit()
	return nil
}

// deleteDocument deletes a document.
func deleteDocument(r *util.Request, folder *model.Folder, doc *model.Document) (err error) {
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	if doc.URL != "" {
		r.Tx.DeleteLink(folder, doc)
	} else {
		r.Tx.DeleteFile(folder, doc)
	}
	r.Tx.Commit()
	return nil
}
