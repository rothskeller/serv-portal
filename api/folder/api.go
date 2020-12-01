package folder

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetFolder handles GET /api/folders/... requests.  These encompass a variety
// of queries, depending on the required "op" URL query parameter.
func GetFolder(r *util.Request) (err error) {
	var (
		folder *model.Folder
	)
	if folder = r.Tx.FetchFolder(r.Path[12:]); folder == nil {
		// 12 is the length of "/api/folders".
		return util.NotFound
	}
	switch r.FormValue("op") {
	case "browse":
		return getBrowseFolder(r, folder)
	case "editFolder":
		return getEditFolder(r, folder)
	case "newFolder":
		return getNewFolder(r, folder)
	case "checkNames":
		return checkFolderNames(r, folder)
	}
	return errors.New("invalid op")
}

// PostFolder handles POST /api/folders/... requests.  These encompass a variety
// of queries and actions, depending on the required "op" form parameter.
func PostFolder(r *util.Request) (err error) {
	var (
		folder *model.Folder
	)
	if folder = r.Tx.FetchFolder(r.Path[12:]); folder == nil {
		// 12 is the length of "/api/folders".
		return util.NotFound
	}
	switch r.FormValue("op") {
	case "editFolder":
		return postEditFolder(r, folder)
	case "newFolder":
		return postNewFolder(r, folder)
	case "move":
		return postMoveFolder(r, folder)
	case "checkNames":
		return checkFolderNames(r, folder)
	case "newFiles":
		return postNewFiles(r, folder)
	case "newLink":
		return postNewLink(r, folder)
	}
	return errors.New("invalid op")
}

// DeleteFolder handles DELETE /api/folders/... requests, which delete folders
// (and the documents within them).
func DeleteFolder(r *util.Request) (err error) {
	var (
		folder *model.Folder
	)
	if folder = r.Tx.FetchFolder(r.Path[12:]); folder == nil {
		// 12 is the length of "/api/folders".
		return util.NotFound
	}
	return deleteFolder(r, folder)
}

// GetDocument handles GET /api/document/... requests, and also GET /dl/...
// requests.  These encompass a variety of queries, depending on the required
// "op" URL query parameter.
func GetDocument(r *util.Request) (err error) {
	var (
		folder *model.Folder
		doc    *model.Document
		path   string
		op     string
	)
	if strings.HasPrefix(r.Path, "/dl") {
		path = r.Path[3:]
		op = "download"
	} else {
		path = r.Path[13:] // length of "/api/document"
		op = r.FormValue("op")
	}
	if folder, doc = r.Tx.FetchDocument(path); folder == nil {
		return util.NotFound
	}
	switch op {
	case "download":
		return downloadFile(r, folder, doc)
	case "edit":
		return getEditDocument(r, folder, doc)
	}
	return errors.New("invalid op")
}

// PostDocument handles POST /api/document/... requests.  These encompass a
// variety of queries and actions, depending on the required "op" form
// parameter.
func PostDocument(r *util.Request) (err error) {
	var (
		folder *model.Folder
		doc    *model.Document
		path   = r.Path[13:] // length of "/api/document"
	)
	if folder, doc = r.Tx.FetchDocument(path); folder == nil {
		return util.NotFound
	}
	switch r.FormValue("op") {
	case "editFile":
		return postEditFile(r, folder, doc)
	case "editLink":
		return postEditLink(r, folder, doc)
	case "move":
		return postMoveDocument(r, folder, doc)
	}
	return errors.New("invalid op")
}

// DeleteDocument handles DELETE /api/document/... requests, which delete
// documents.
func DeleteDocument(r *util.Request) (err error) {
	var (
		folder *model.Folder
		doc    *model.Document
		path   = r.Path[13:] // length of "/api/document"ng
	)
	if folder, doc = r.Tx.FetchDocument(path); folder == nil {
		return util.NotFound
	}
	return deleteDocument(r, folder, doc)
}
