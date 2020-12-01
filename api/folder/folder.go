package folder

import (
	"errors"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// getNewFolder returns the information needed when starting to create a folder
// as a child of the specified folder.
func getNewFolder(r *util.Request, parent *model.Folder) (err error) {
	var (
		allowed []vo
		out     jwriter.Writer
	)
	if !canEditFolder(r.Person, parent) {
		return util.Forbidden
	}
	allowed = allowedVisibilities(r.Person, parent, nil)
	out.RawString(`{"name":"","visibility":`)
	out.String(allowed[0].String())
	out.RawString(`,"allowedVisibilities":[`)
	for i, vo := range allowed {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(vo.String())
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// postNewFolder creates a new folder as a child of the specified folder.
func postNewFolder(r *util.Request, parent *model.Folder) (err error) {
	var (
		folder model.Folder
		vo     vo
	)
	if !canEditFolder(r.Person, parent) {
		return util.Forbidden
	}
	if folder.Name = strings.TrimSpace(r.FormValue("name")); folder.Name == "" {
		return errors.New("missing name")
	}
	folder.URL = parent.URL + "/" + nameToURL(folder.Name)
	for _, child := range r.Tx.FetchSubFolders(parent) {
		if child.URL == folder.URL {
			return util.SendConflict(r, "The parent folder “%s” already contains a folder named “%s”.", parent.Name, child.Name)
		}
	}
	for _, doc := range r.Tx.FetchDocuments(parent) {
		if doc.Name == folder.Name {
			return util.SendConflict(r, "The parent folder “%s” already contains a document named “%s”.", parent.Name, doc.Name)
		}
	}
	if vo, err = parseVO(r.FormValue("visibility")); err != nil {
		return err
	}
	if !allowedVisibility(r.Person, parent, nil, vo.v, vo.o) {
		return errors.New("invalid visibility")
	}
	folder.Visibility, folder.Org = vo.v, vo.o
	r.Tx.CreateFolder(&folder)
	r.Tx.Commit()
	return nil
}

// getEditFolder returns the information needed when starting to edit the
// specified folder.
func getEditFolder(r *util.Request, folder *model.Folder) (err error) {
	var (
		parent   *model.Folder
		children []*model.Folder
		allowed  []vo
		out      jwriter.Writer
	)
	if !canEditFolder(r.Person, folder) || folder.URL == "" {
		return util.Forbidden
	}
	parent = r.Tx.FetchParentFolder(folder)
	children = r.Tx.FetchSubFolders(folder)
	allowed = allowedVisibilities(r.Person, parent, children)
	out.RawString(`{"name":`)
	out.String(folder.Name)
	out.RawString(`,"visibility":`)
	out.String(folderVOString(folder))
	out.RawString(`,"allowedVisibilities":[`)
	for i, vo := range allowed {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(vo.String())
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// postEditFolder updates the metadata for the specified folder.
func postEditFolder(r *util.Request, folder *model.Folder) (err error) {
	var (
		parent *model.Folder
		newf   model.Folder
		vo     vo
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	parent = r.Tx.FetchParentFolder(folder)
	if newf.Name = strings.TrimSpace(r.FormValue("name")); newf.Name == "" {
		return errors.New("missing name")
	}
	newf.URL = parent.URL + "/" + nameToURL(newf.Name)
	if newf.URL != folder.URL {
		for _, child := range r.Tx.FetchSubFolders(parent) {
			if child.URL == newf.URL {
				return util.SendConflict(r, "The parent folder “%s” already contains a folder named “%s”.", parent.Name, child.Name)
			}
		}
	}
	if newf.Name != folder.Name {
		for _, doc := range r.Tx.FetchDocuments(parent) {
			if doc.Name == newf.Name {
				return util.SendConflict(r, "The parent folder “%s” already contains a document named “%s”.", parent.Name, doc.Name)
			}
		}
	}
	if vo, err = parseVO(r.FormValue("visibility")); err != nil {
		return err
	}
	if !allowedVisibility(r.Person, parent, r.Tx.FetchSubFolders(folder), vo.v, vo.o) {
		return errors.New("invalid visibility")
	}
	newf.Visibility, newf.Org = vo.v, vo.o
	r.Tx.UpdateFolder(folder, &newf)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "text/plain; utf-8")
	r.Write([]byte(newf.URL))
	return nil
}

// postMoveFolder moves the specified folder to a new parent.
func postMoveFolder(r *util.Request, folder *model.Folder) (err error) {
	var (
		pURL   string
		parent *model.Folder
		newf   model.Folder
	)
	if !canEditFolder(r.Person, folder) {
		return util.Forbidden
	}
	newf = *folder
	if r.Form["parent"] == nil {
		return errors.New("missing parent")
	}
	if pURL = r.FormValue("parent"); pURL == folder.URL || pURL == parentURL(folder) { // no-op
		r.Header().Set("Content-Type", "text/plain; utf-8")
		r.Write([]byte(folder.URL))
		return nil
	}
	if parent = r.Tx.FetchFolder(r.FormValue("parent")); parent == nil {
		return errors.New("nonexistent parent")
	}
	if !canEditFolder(r.Person, parent) {
		return util.Forbidden
	}
	newf.URL = parent.URL + "/" + nameToURL(newf.Name)
	for _, child := range r.Tx.FetchSubFolders(parent) {
		if child.URL == newf.URL {
			return util.SendConflict(r, "The destination folder “%s” already contains a folder named “%s”.", parent.Name, newf.Name)
		}
	}
	for _, doc := range r.Tx.FetchDocuments(parent) {
		if doc.Name == newf.Name {
			return util.SendConflict(r, "The destination folder “%s” already contains a document named “%s”.", parent.Name, newf.Name)
		}
	}
	if !allowedVisibility(r.Person, parent, r.Tx.FetchSubFolders(folder), newf.Visibility, newf.Org) {
		return util.SendConflict(r, "The destination folder “%s” is only visible to %s. It cannot contain the folder “%s”, which is visible to %s.", parent.Name, folderVOLabel(parent), newf.Name, folderVOLabel(&newf))
	}
	r.Tx.UpdateFolder(folder, &newf)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "text/plain; utf-8")
	r.Write([]byte(newf.URL))
	return nil
}

// deleteFolder deletes the specified folder.
func deleteFolder(r *util.Request, folder *model.Folder) (err error) {
	if !canEditFolder(r.Person, folder) || folder.URL == "" {
		return util.Forbidden
	}
	r.Tx.DeleteFolder(folder)
	r.Tx.Commit()
	return nil
}

// checkFolderNames returns whether the specified name(s) are in use in the
// specified folder.  It is used when dropping one or more file system files
// onto a folder, to see whether they would replace files already in that
// folder.
func checkFolderNames(r *util.Request, folder *model.Folder) (err error) {
	var (
		children    []*model.Folder
		docs        []*model.Document
		wantReplace bool
		replace     []string
	)
	if !CanViewFolder(r.Person, folder) {
		return util.Forbidden
	}
	children = r.Tx.FetchSubFolders(folder)
	docs = r.Tx.FetchDocuments(folder)
	wantReplace = r.FormValue("replace") == "true"
	for _, name := range r.Form["name"] {
		url := nameToURL(name)
		for _, child := range children {
			if nameToURL(child.Name) == url {
				return util.SendConflict(r, "The folder “%s” already contains a folder named “%s”.", folder.Name, child.Name)
			}
		}
		for _, doc := range docs {
			if doc.Name == name {
				if doc.URL != "" {
					return util.SendConflict(r, "The folder “%s” already contains a link named “%s”.", folder.Name, doc.Name)
				}
				if !wantReplace {
					return util.SendConflict(r, "The folder “%s” already contains a file named “%s”.", folder.Name, doc.Name)
				}
				replace = append(replace, doc.Name)
				break
			}
		}
	}
	r.Tx.Commit()
	if len(replace) == 0 {
		return nil
	}
	r.Header().Set("X-Can-Replace", "true")
	switch len(replace) {
	case 1:
		return util.SendConflict(r, "The folder “%s” already contains a file named “%s”. Do you want to replace it?", folder.Name, replace[0])
	case 2:
		return util.SendConflict(r, "The folder “%s” already contains files named “%s” and “%s”. Do you want to replace them?", folder.Name, replace[0], replace[1])
	default:
		return util.SendConflict(r, "The folder “%s” already contains files named “%s”, and “%s”. Do you want to replace them?", folder.Name, strings.Join(replace[:len(replace)-1], "”, “"), replace[len(replace)-1])
	}
}
