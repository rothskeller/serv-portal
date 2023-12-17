package files

import (
	"fmt"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// FolderFields is the set of Folder fields that need to be retrieved in Folder
// objects passed to the GetFolder function.
const FolderFields = folder.FID | folder.FName | folder.FURLName | folder.FViewer | folder.FEditor

// Handle handles /files/${path} requests.
func Handle(r *request.Request) {
	var (
		flist   []*folder.Folder
		docname string
		user    *person.Person
	)
	user = auth.SessionUser(r, 0, false)
	// Try to find the folder with that path.
	if flist, docname = folder.WithPath(r, r.Path, FolderFields); flist == nil {
		errpage.NotFound(r, user)
		return
	}
	// Make sure we have view privilege on the leaf folder.
	if !user.HasPrivLevel(flist[len(flist)-1].Viewer()) {
		if user == nil { // not logged in, so send a 401
			auth.SessionUser(r, 0, true)
		} else { // logged in, so send a 403
			errpage.Forbidden(r, user)
		}
		return
	}
	// Honor their document or folder request.
	if docname != "" {
		getDocument(r, user, flist[len(flist)-1], docname)
	} else if r.Method == http.MethodPost {
		postFolder(r, user, flist)
	} else {
		GetFolder(r, user, flist, 0, nil)
	}
}

// getDocument handles GET /files/${path}/${docname}.  Permissions have already
// been checked.
func getDocument(r *request.Request, user *person.Person, f *folder.Folder, docname string) {
	var (
		doc  *document.Document
		fh   *os.File
		stat fs.FileInfo
		err  error
	)
	if doc = document.WithName(r, f.ID(), docname); doc == nil || doc.URL != "" {
		errpage.NotFound(r, user)
		return
	}
	fh = document.Open(doc.ID)
	if stat, err = fh.Stat(); err != nil {
		panic(err)
	}
	switch strings.ToLower(filepath.Ext(docname)) {
	case ".jpeg", ".jpg", ".png", ".pdf":
		r.Header().Set("Content-Disposition", "inline")
	default:
		r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", docname))
	}
	r.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(r, r.Request, doc.Name, stat.ModTime(), fh)
}

// postFolder handles POST /files/${path} when the path is a folder.  This is a
// post from a hidden form, used when something is dragged and dropped on the
// folder.
func postFolder(r *request.Request, user *person.Person, flist []*folder.Folder) {
	var (
		newfolder folder.ID
		newdocs   map[document.ID]bool
		f         = flist[len(flist)-1]
	)
	if !user.HasPrivLevel(f.Editor()) {
		errpage.Forbidden(r, user)
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if idstr := r.FormValue("doc"); idstr != "" {
		// They dropped a document from another folder onto this one.
		docid, ok := handleDropDocument(r, user, f, idstr)
		if !ok {
			return
		}
		if docid != 0 {
			newdocs = make(map[document.ID]bool)
			newdocs[docid] = true
		}
	}
	if idstr := r.FormValue("deldoc"); idstr != "" {
		// They dropped a document onto the trash target.
		if !handleDeleteDocument(r, user, idstr) {
			return
		}
	}
	if idstr := r.FormValue("folder"); idstr != "" {
		// They dropped another folder onto this one.
		folderid, ok := handleDropFolder(r, user, f, idstr)
		if !ok {
			return
		}
		newfolder = folderid
	}
	if idstr := r.FormValue("delfolder"); idstr != "" {
		// They dropped a folder onto the trash target.
		folderid, ok := handleDeleteFolder(r, user, idstr)
		if !ok {
			return
		}
		if folderid == f.ID() {
			flist = flist[:len(flist)-1]
		}
	}
	if u := r.FormValue("url"); u != "" {
		// They dropped a URL onto this folder.
		docid := handleDropURL(r, f, u)
		newdocs = map[document.ID]bool{docid: true}
	}
	if r.MultipartForm != nil && r.MultipartForm.File != nil && len(r.MultipartForm.File["file"]) != 0 {
		// They dropped one or more files onto this folder.
		files := r.MultipartForm.File["file"]
		newdocs = make(map[document.ID]bool)
		for _, file := range files {
			docid := handleDropFile(r, f, file)
			newdocs[docid] = true
		}
	}
	GetFolder(r, user, flist, newfolder, newdocs)
}

// GetFolder handles GET /files/${path} when the path is a folder.  Permissions
// have already been checked.
func GetFolder(r *request.Request, user *person.Person, flist []*folder.Folder, newfolder folder.ID, newdocs map[document.ID]bool) {
	var (
		location         string
		anythingEditable bool
	)
	// Notify the unpoly library in the client that this response is
	// GETtable, even if it is being sent as the response to a POST, i.e.,
	// after a drag/drop operation or an edit dialog submission.
	for _, f := range flist {
		if location == "" {
			location = "/files"
		} else {
			location = path.Join(location, f.URLName())
		}
	}
	r.Header().Set("X-Up-Method", "GET")
	r.Header().Set("X-Up-Location", location)
	// Now display the folder page.
	ui.Page(r, user, ui.PageOpts{Title: r.Loc("Files")}, func(main *htmlb.Element) {
		var fdiv = main
		var fpath string

		// Display each of the ancestor folders, up to and including the
		// target folder.
		for i, f := range flist {
			var label string

			if fpath == "" && i == len(flist)-1 && user == nil {
				fpath = "/"
				label = "Sunnyvale SERV"
			} else if fpath == "" {
				fpath = "/files"
				label = f.Name()
			} else {
				fpath = path.Join(fpath, f.URLName())
				label = f.Name()
			}
			canEdit := user.HasPrivLevel(f.Editor())
			canDelete := i == len(flist)-1 && canEdit && !folder.ExistsWithParent(r, f.ID()) && !document.ExistInFolder(r, f.ID())
			canDrag := canDelete || (i > 0 && user.HasPrivLevel(flist[i-1].Editor()))
			anythingEditable = anythingEditable || canEdit || canDelete
			fline := fdiv.E("div class=folder data-id=%d", f.ID(),
				canEdit, "editable data-path=%s", fpath,
				canDelete, "deletable",
				canDrag, "draggable=true")
			fline.E("s-icon icon=folder-open")
			if user.HasPrivLevel(f.Viewer()) {
				fline.E("a href=%s up-target=main", fpath, canDrag, "draggable=false").T(label)
			} else {
				fline.T(label)
			}
			fdiv = fdiv.E("div class=folderContents")
			if fpath == "/" {
				fpath = "/files"
			}
		}
		f := flist[len(flist)-1]
		canEdit := user.HasPrivLevel(f.Editor())
		fdiv.E("input type=hidden id=folderpath value=%s", fpath)
		fdiv.E("input type=hidden id=folderid value=%d", int(f.ID()))
		// Display the child folders of the target folder.
		folder.AllWithParent(r, f.ID(), FolderFields, func(cf *folder.Folder) {
			if user.HasPrivLevel(cf.Viewer()) {
				canEditChild := user.HasPrivLevel(cf.Editor())
				canDeleteChild := canEditChild && !folder.ExistsWithParent(r, cf.ID()) && !document.ExistInFolder(r, cf.ID())
				anythingEditable = anythingEditable || canEditChild || canDeleteChild
				fdiv.E("div class=folder data-id=%d", cf.ID(),
					cf.ID() == newfolder, "class=folderItem-new",
					canEditChild, "editable data-path=%s", path.Join(fpath, cf.URLName()),
					canDeleteChild, "deletable",
					canEdit, "draggable=true").
					E("s-icon icon=folder").P().
					E("a href=%s up-target=main", path.Join(fpath, cf.URLName()), canEdit, "draggable=false").
					T(cf.Name())
			}
		})
		// Display the documents in the target folder.
		document.AllInFolder(r, f.ID(), func(doc *document.Document) {
			var newtab bool
			ddiv := fdiv.E("div class=document data-id=%d", doc.ID,
				newdocs != nil && newdocs[doc.ID], "class=folderItem-new",
				canEdit, "editable deletable draggable=true")
			if doc.URL != "" {
				ddiv.E("s-icon icon=link")
				ddiv.E("a href=%s target=_blank", doc.URL, canEdit, "draggable=false").T(doc.Name)
				return
			}
			switch strings.ToLower(filepath.Ext(doc.Name)) {
			case ".pdf":
				newtab = true
				ddiv.E("s-icon icon=pdf")
			case ".png", ".jpeg", ".jpg":
				newtab = true
				ddiv.E("s-icon icon=image")
			case ".docx", ".doc":
				ddiv.E("s-icon icon=word")
			case ".pptx", ".ppt":
				ddiv.E("s-icon icon=powerpoint")
			case ".xlsx", ".xls":
				ddiv.E("s-icon icon=excel")
			default:
				ddiv.E("s-icon icon=file")
			}
			ddiv.E("a href=%s", path.Join(fpath, url.PathEscape(doc.Name)),
				newtab, "target=_blank",
				canEdit, "draggable=false").
				T(doc.Name)
		})
		// If the target folder is editable, display the add buttons.
		if canEdit {
			buttons := main.E("div id=folderButtons class=formButtons")
			buttons.E("a href=/docedit/%d/NEWFILE class='sbtn sbtn-primary' up-layer=new up-size=grow up-dismissable=key up-history=false>Add File", f.ID())
			buttons.E("a href=/docedit/%d/NEWURL class='sbtn sbtn-primary' up-layer=new up-size=grow up-dismissable=key up-history=false>Add Web Link", f.ID())
			buttons.E("a href=/folderedit/NEW?parent=%d class='sbtn sbtn-primary' up-layer=new up-size=grow up-dismissable=key up-history=false>Add Folder", f.ID())
			deltarget := main.E("div id=folderDelete style=display:none")
			deltarget.E("s-icon icon=trash")
			deltarget.R("Drag here to delete.")
		}
		// If any folder on the page is editable, add the hidden form
		// that's used by the client code to handle drag and drop.
		if anythingEditable {
			form := main.E("form method=POST enctype=multipart/form-data up-target=main style=display:none")
			form.E("input type=hidden name=csrf value=%s", r.CSRF)
			form.E("input type=file id=folderDropFiles name=file")
			form.E("input type=hidden id=folderDropURL name=url")
			form.E("input type=hidden id=folderDropDoc name=doc")
			form.E("input type=hidden id=folderDropFolder name=folder")
			form.E("input type=hidden id=folderTrashDoc name=deldoc")
			form.E("input type=hidden id=folderTrashFolder name=delfolder")
		}
	})
}
