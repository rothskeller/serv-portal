package folder

import (
	"errors"
	"mime/multipart"
	"os"
	"sort"
	"strconv"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// PostFolder handles POST /api/folders/$id requests.
func PostFolder(r *util.Request, idstr string) (err error) {
	var (
		parentID model.FolderID
		parent   *model.Folder
		folder   model.Folder
	)
	if idstr != "0" {
		parentID = model.FolderID(util.ParseID(idstr))
		if parent = r.Tx.FetchFolder(parentID); parent == nil {
			return util.NotFound
		}
		if parent.Group == 0 && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
		if parent.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, parent.Group) {
			return util.Forbidden
		}
	} else if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	folder.Name = r.FormValue("name")
	folder.Group = model.GroupID(util.ParseID(r.FormValue("group")))
	folder.Parent = parentID
	if err = ValidateFolder(r.Tx, &folder); err != nil {
		return err
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	r.Tx.CreateFolder(&folder)
	return GetFolder(r, idstr)
}

// PutFolder handles PUT /api/folders/$id requests.
func PutFolder(r *util.Request, idstr string) (err error) {
	var (
		folder     *model.Folder
		parent     *model.Folder
		prevParent model.FolderID
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(idstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	prevParent = folder.Parent
	folder.Name = r.FormValue("name")
	folder.Group = model.GroupID(util.ParseID(r.FormValue("group")))
	folder.Parent = model.FolderID(util.ParseID(r.FormValue("parent")))
	if err = ValidateFolder(r.Tx, folder); err != nil {
		return err
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	if folder.Parent == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Parent != 0 {
		parent = r.Tx.FetchFolder(folder.Parent)
		if parent.Group == 0 && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
		if parent.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, parent.Group) {
			return util.Forbidden
		}
	}
	r.Tx.UpdateFolder(folder)
	return GetFolder(r, strconv.Itoa(int(prevParent)))
}

// DeleteFolder handles DELETE /api/folders/$id requests.
func DeleteFolder(r *util.Request, idstr string) (err error) {
	var (
		folder *model.Folder
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(idstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	r.Tx.DeleteFolder(folder)
	return GetFolder(r, strconv.Itoa(int(folder.Parent)))
}

// PostDocument handles POST /api/folders/$fid/$did requests (but not $did="NEW").
func PostDocument(r *util.Request, fidstr, didstr string) (err error) {
	var (
		folder      *model.Folder
		newFolderID model.FolderID
		docID       model.DocumentID
		doc         *model.Document
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(fidstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	docID = model.DocumentID(util.ParseID(didstr))
	for _, d := range folder.Documents {
		if d.ID == docID {
			doc = d
			break
		}
	}
	if doc == nil {
		return util.NotFound
	}
	doc.Name = r.FormValue("name")
	doc.PostedBy = r.Person.ID
	doc.PostedAt = time.Now()
	newFolderID = model.FolderID(util.ParseID(r.FormValue("folder")))
	if newFolderID != folder.ID {
		var (
			newFolder *model.Folder
			newDoc    model.Document
			maxDocID  model.DocumentID
			contents  *os.File
		)
		if newFolder = r.Tx.FetchFolder(newFolderID); newFolder == nil {
			return errors.New("nonexistent folder")
		}
		for _, d := range newFolder.Documents {
			if d.ID > maxDocID {
				maxDocID = d.ID
			}
		}
		newDoc.ID = maxDocID + 1
		newDoc.Name = doc.Name
		newFolder.Documents = append(newFolder.Documents, &newDoc)
		if err = ValidateFolder(r.Tx, newFolder); err != nil {
			return err
		}
		sort.Slice(newFolder.Documents, func(i, j int) bool { return newFolder.Documents[i].Name < newFolder.Documents[j].Name })
		contents = r.Tx.FetchDocument(folder, doc.ID)
		r.Tx.CreateDocument(newFolder, newDoc.ID, contents)
		contents.Close()
		r.Tx.UpdateFolder(newFolder)
		r.Tx.DeleteDocument(folder, doc.ID)
		j := 0
		for _, d := range folder.Documents {
			if d != doc {
				folder.Documents[j] = d
				j++
			}
		}
		folder.Documents = folder.Documents[:j]
		r.Tx.UpdateFolder(folder)
	} else {
		if err = ValidateFolder(r.Tx, folder); err != nil {
			return err
		}
		sort.Slice(folder.Documents, func(i, j int) bool { return folder.Documents[i].Name < folder.Documents[j].Name })
		r.Tx.UpdateFolder(folder)
	}
	return GetFolder(r, fidstr)
}

// PostNewDocuments handles POST /api/folders/$id/NEW requests.
func PostNewDocuments(r *util.Request, idstr string) (err error) {
	var (
		folder   *model.Folder
		files    []*multipart.FileHeader
		docs     []*model.Document
		maxDocID model.DocumentID
		j        int
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(idstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	for _, doc := range folder.Documents {
		if doc.ID > maxDocID {
			maxDocID = doc.ID
		}
	}
	r.ParseMultipartForm(1048576)
	files = r.MultipartForm.File["file"]
	if files == nil || len(files) == 0 {
		return errors.New("missing files")
	}
	for _, file := range files {
		var (
			doc model.Document
			fh  multipart.File
		)
		maxDocID++
		doc.ID = maxDocID
		doc.Name = file.Filename
		doc.PostedBy = r.Person.ID
		doc.PostedAt = time.Now()
		if fh, err = file.Open(); err != nil {
			goto ERROR
		}
		r.Tx.CreateDocument(folder, doc.ID, fh)
		docs = append(docs, &doc)
	}
	j = 0
	for _, doc := range folder.Documents {
		found := false
		for _, file := range files {
			if file.Filename == doc.Name {
				found = true
				break
			}
		}
		if found {
			r.Tx.DeleteDocument(folder, doc.ID)
		} else {
			folder.Documents[j] = doc
			j++
		}
	}
	folder.Documents = append(folder.Documents[:j], docs...)
	sort.Slice(folder.Documents, func(i, j int) bool { return folder.Documents[i].Name < folder.Documents[j].Name })
	r.Tx.UpdateFolder(folder)
	return GetFolder(r, idstr)
ERROR:
	for _, doc := range docs {
		r.Tx.DeleteDocument(folder, doc.ID)
	}
	return err
}

// DeleteDocument handles DELETE /api/folders/$fid/$did requests.
func DeleteDocument(r *util.Request, fidstr, didstr string) (err error) {
	var (
		folder *model.Folder
		docID  model.DocumentID
		j      int
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(fidstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	docID = model.DocumentID(util.ParseID(didstr))
	j = 0
	for _, doc := range folder.Documents {
		if doc.ID == docID {
			r.Tx.DeleteDocument(folder, docID)
		} else {
			folder.Documents[j] = doc
			j++
		}
	}
	if j == len(folder.Documents) {
		return util.NotFound
	}
	folder.Documents = folder.Documents[:j]
	r.Tx.UpdateFolder(folder)
	return GetFolder(r, fidstr)
}
