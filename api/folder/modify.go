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
		folder   = model.FolderNode{Folder: new(model.Folder)}
	)
	if idstr != "0" {
		parentID = model.FolderID(util.ParseID(idstr))
		if folder.ParentNode = r.Tx.FetchFolder(parentID); folder.ParentNode == nil {
			return util.NotFound
		}
		if folder.ParentNode.Group == 0 && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
		if folder.ParentNode.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.ParentNode.Group) {
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
		folder     *model.FolderNode
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
	r.Tx.WillUpdateFolder(folder)
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
	if folder.ParentNode == nil && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.ParentNode != nil {
		if folder.ParentNode.Group == 0 && !r.Auth.IsWebmaster() {
			return util.Forbidden
		}
		if folder.ParentNode.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.ParentNode.Group) {
			return util.Forbidden
		}
	}
	r.Tx.UpdateFolder(folder)
	propagateApprovalCounts(r, folder.ParentNode)
	return GetFolder(r, strconv.Itoa(int(prevParent)))
}

// DeleteFolder handles DELETE /api/folders/$id requests.
func DeleteFolder(r *util.Request, idstr string) (err error) {
	var folder *model.FolderNode

	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(idstr))); folder == nil {
		return util.NotFound
	}
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	pn := folder.ParentNode
	r.Tx.DeleteFolder(folder)
	propagateApprovalCounts(r, pn)
	return GetFolder(r, strconv.Itoa(int(folder.Parent)))
}

// PostDocument handles POST /api/folders/$fid/$did requests (but not $did="NEW").
func PostDocument(r *util.Request, fidstr, didstr string) (err error) {
	var (
		folder      *model.FolderNode
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
	r.Tx.WillUpdateFolder(folder)
	doc.Name = r.FormValue("name")
	doc.NeedsApproval = false
	newFolderID = model.FolderID(util.ParseID(r.FormValue("folder")))
	if newFolderID != folder.ID {
		var (
			newFolder *model.FolderNode
			newDoc    model.Document
			maxDocID  model.DocumentID
			contents  *os.File
		)
		if newFolder = r.Tx.FetchFolder(newFolderID); newFolder == nil {
			return errors.New("nonexistent folder")
		}
		r.Tx.WillUpdateFolder(newFolder)
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
		propagateApprovalCounts(r, newFolder.ParentNode)
		r.Tx.DeleteDocument(folder, doc.ID)
		j := 0
		for _, d := range folder.Documents {
			if d != doc {
				folder.Documents[j] = d
				j++
			}
		}
		folder.Documents = folder.Documents[:j]
		if err = ValidateFolder(r.Tx, folder); err != nil {
			return err
		}
		r.Tx.UpdateFolder(folder)
		propagateApprovalCounts(r, folder.ParentNode)
	} else {
		j := 0
		for _, d := range folder.Documents {
			if d.Name != doc.Name || d.ID == doc.ID {
				folder.Documents[j] = d
				j++
			} else {
				r.Tx.DeleteDocument(folder, d.ID)
			}
		}
		folder.Documents = folder.Documents[:j]
		if err = ValidateFolder(r.Tx, folder); err != nil {
			return err
		}
		sort.Slice(folder.Documents, func(i, j int) bool { return folder.Documents[i].Name < folder.Documents[j].Name })
		r.Tx.UpdateFolder(folder)
		propagateApprovalCounts(r, folder.ParentNode)
	}
	return GetFolder(r, fidstr)
}

// PostNewDocuments handles POST /api/folders/$id/NEW requests.
func PostNewDocuments(r *util.Request, idstr string) (err error) {
	var (
		folder        *model.FolderNode
		files         []*multipart.FileHeader
		docs          []*model.Document
		maxDocID      model.DocumentID
		needsApproval bool
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(idstr))); folder == nil {
		return util.NotFound
	}
	if !r.Auth.CanA(model.PrivManageFolders) {
		return util.Forbidden
	}
	r.Tx.WillUpdateFolder(folder)
	if folder.Group == 0 && !r.Auth.IsWebmaster() {
		needsApproval = true
	}
	if folder.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		needsApproval = true
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
		doc.NeedsApproval = needsApproval
		if fh, err = file.Open(); err != nil {
			goto ERROR
		}
		r.Tx.CreateDocument(folder, doc.ID, fh)
		docs = append(docs, &doc)
	}
	if !needsApproval {
		j := 0
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
		folder.Documents = folder.Documents[:j]
	}
	folder.Documents = append(folder.Documents, docs...)
	if err = ValidateFolder(r.Tx, folder); err != nil {
		goto ERROR
	}
	sort.Slice(folder.Documents, func(i, j int) bool { return folder.Documents[i].Name < folder.Documents[j].Name })
	r.Tx.UpdateFolder(folder)
	propagateApprovalCounts(r, folder.ParentNode)
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
		folder *model.FolderNode
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
	r.Tx.WillUpdateFolder(folder)
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
	if err = ValidateFolder(r.Tx, folder); err != nil {
		return err
	}
	r.Tx.UpdateFolder(folder)
	propagateApprovalCounts(r, folder.ParentNode)
	return GetFolder(r, fidstr)
}

func propagateApprovalCounts(r *util.Request, folder *model.FolderNode) {
	if folder == nil {
		return
	}
	newApprovals := 0
	for _, cf := range folder.ChildNodes {
		newApprovals += cf.Approvals
	}
	for _, d := range folder.Documents {
		if d.NeedsApproval {
			newApprovals++
		}
	}
	if folder.Approvals != newApprovals {
		r.Tx.WillUpdateFolder(folder)
		folder.Approvals = newApprovals
		r.Tx.UpdateFolder(folder)
		propagateApprovalCounts(r, folder.ParentNode)
	}
}
