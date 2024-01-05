package docedit

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/files"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /docedit/${folderid}/${docid} requests.  ${docid} may be a
// document ID, or it may be the word "NEWURL" or "NEWFILE".
func Handle(r *request.Request, fidstr, didstr string) {
	var (
		user *person.Person
		doc  *document.Document
		ud   *document.Updater
		f    *folder.Folder
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if f = folder.WithID(r, folder.ID(util.ParseID(fidstr)), files.FolderFields|folder.FParent); f == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(f.Editor()) {
		errpage.Forbidden(r, user)
		return
	}
	if didstr == "NEWURL" || didstr == "NEWFILE" {
		ud = &document.Updater{Folder: f}
	} else {
		if doc = document.WithID(r, document.ID(util.ParseID(didstr))); doc == nil || doc.Archived || doc.Folder != f.ID() {
			errpage.NotFound(r, user)
			return
		}
		if r.FormValue("delete") != "" {
			handleDelete(r, user, f, doc)
			return
		}
		ud = doc.Updater(r, f)
	}
	if ud.URL != "" || didstr == "NEWURL" {
		handleURL(r, user, f, doc, ud)
	} else {
		handleFile(r, user, f, doc, ud)
	}
}

func handleFile(r *request.Request, user *person.Person, f *folder.Folder, doc *document.Document, ud *document.Updater) {
	var nameError, fileError string

	if r.Method == http.MethodPost {
		nameError = readNameForFile(r, ud)
		fileError = readFile(r, ud)
		if nameError == "" && fileError == "" {
			r.Transaction(func() {
				if ud.ID == 0 {
					document.Create(r, ud)
				} else {
					doc.Update(r, ud)
				}
			})
			files.GetFolder(r, user, f.FolderPath(r, files.FolderFields), 0, map[document.ID]bool{ud.ID: true})
			return
		}
	}
	r.HTMLNoCache()
	if nameError != "" || fileError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST enctype=multipart/form-data up-main up-layer=parent up-target=main")
	if ud.ID == 0 {
		form.E("div class='formTitle formTitle-primary'>Add File")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit File")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitNameForFile(form, ud, nameError != "" || fileError == "", nameError)
	emitFile(form, fileError != "", fileError)
	emitButtons(form, ud.ID != 0)
}

func handleURL(r *request.Request, user *person.Person, f *folder.Folder, doc *document.Document, ud *document.Updater) {
	var nameError, urlError string
	var validate = strings.Fields(r.Request.Header.Get("X-Up-Validate"))

	if r.Method == http.MethodPost {
		nameError = readNameForURL(r, ud)
		urlError = readURL(r, user, ud)
		if len(validate) == 0 && nameError == "" && urlError == "" {
			r.Transaction(func() {
				if ud.ID == 0 {
					document.Create(r, ud)
				} else {
					doc.Update(r, ud)
				}
			})
			files.GetFolder(r, user, f.FolderPath(r, files.FolderFields), 0, map[document.ID]bool{ud.ID: true})
			return
		}
	}
	r.HTMLNoCache()
	if nameError != "" || urlError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	if ud.ID == 0 {
		form.E("div class='formTitle formTitle-primary'>Add Web Link")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit Web Link")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "name") {
		emitNameForURL(form, ud, nameError != "" || urlError == "", nameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "url") {
		emitURL(form, ud, urlError != "", urlError)
	}
	if len(validate) == 0 {
		emitButtons(form, ud.ID != 0)
	}
}

func readNameForFile(r *request.Request, ud *document.Updater) string {
	if ud.Name = strings.TrimSpace(r.FormValue("name")); ud.Name == "" {
		return "The file name is required."
	}
	if ud.Name[0] == '.' || strings.ContainsAny(ud.Name, "/:") {
		return fmt.Sprintf("%q is not a valid name.  Names may not start with a period, and may not contain slashes or colons.", ud.Name)
	}
	if ud.DuplicateName(r) {
		return fmt.Sprintf("The name %q is in use by another document in this folder.", ud.Name)
	}
	return ""
}

func emitNameForFile(form *htmlb.Element, ud *document.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=doceditName>File Name")
	row.E("input id=doceditName name=name value=%s", ud.Name, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readFile(r *request.Request, ud *document.Updater) string {
	var files []*multipart.FileHeader

	if r.MultipartForm.File != nil {
		files = r.MultipartForm.File["file"]
	}
	if len(files) > 1 {
		return "Multiple-file uploads are not allowed through this dialog.  (They can be performed with drag-and-drop, however.)"
	}
	if ud.ID == 0 && len(files) == 0 {
		return "The file to be added must be provided here."
	}
	if len(files) != 0 {
		var mf multipart.File
		var err error

		if mf, err = files[0].Open(); err != nil {
			return "File was not uploaded correctly: " + err.Error()
		}
		defer mf.Close()
		if ud.Contents, err = io.ReadAll(mf); err != nil {
			return "File was not uploaded correctly: " + err.Error()
		}
	}
	return ""
}

func emitFile(form *htmlb.Element, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=doceditFile>File Contents")
	row.E("input type=file id=doceditFile name=file", focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readNameForURL(r *request.Request, ud *document.Updater) string {
	if ud.Name = strings.TrimSpace(r.FormValue("name")); ud.Name == "" {
		return "The web link name is required."
	}
	if ud.DuplicateName(r) {
		return fmt.Sprintf("The name %q is in use by another document in this folder.", ud.Name)
	}
	return ""
}

func emitNameForURL(form *htmlb.Element, ud *document.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=doceditName>Link Name")
	row.E("input id=doceditName name=name s-validate value=%s", ud.Name, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readURL(r *request.Request, user *person.Person, ud *document.Updater) string {
	if ud.URL = strings.TrimSpace(r.FormValue("url")); ud.URL == "" {
		return "The web link URL is required."
	}
	if u, err := url.Parse(ud.URL); err != nil {
		return fmt.Sprintf("%q is not a valid URL.", ud.URL)
	} else if u.Scheme != "http" && u.Scheme != "https" && (u.Scheme != "" || !user.IsAdminLeader()) {
		return fmt.Sprintf("The %q URL scheme is not supported.  Only \"http\" and \"https\" are supported.", u.Scheme)
	}
	return ""
}

func emitURL(form *htmlb.Element, ud *document.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=doceditURL>Link URL")
	row.E("input id=doceditURL name=url s-validate value=%s", ud.URL, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func emitButtons(form *htmlb.Element, canDelete bool) {
	buttons := form.E("div class=formButtons")
	if canDelete {
		buttons.E("div class=formButtonSpace")
	}
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	if canDelete {
		// This button comes last in the tree order so that it is not
		// the default.  But it comes first in the visual order because
		// of the formButton-beforeAll class.
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
}

func handleDelete(r *request.Request, user *person.Person, f *folder.Folder, doc *document.Document) {
	r.Transaction(func() {
		if doc.URL != "" {
			// URL documents we simply delete.  If we ever need to
			// restore them, there's enough information in the log
			// to reconstruct them.
			doc.Delete(r)
		} else {
			// Files we don't actually delete; instead we archive
			// them.  This allows the webmaster to restore them if
			// they were deleted accidentally (or maliciously).  The
			// only way to actually delete an archived file (and
			// therefore, any folder that contains archived files)
			// is through the use of offline tools.
			ud := doc.Updater(r, f)
			ud.Archived = true
			doc.Update(r, ud)
		}
	})
	files.GetFolder(r, user, f.FolderPath(r, files.FolderFields), 0, nil)
}
