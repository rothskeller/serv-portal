package files

import (
	"context"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	xhtml "golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

// handleDropDocument handles a document (identified by docidstr) being dropped
// onto folder f.  It returns the ID of the document (or zero if the drop was
// rejected) and whether or not the folder page should be emitted.  (The latter
// will be false if an error page has been emitted.)
func handleDropDocument(r *request.Request, user *person.Person, f *folder.Folder, docidstr string) (document.ID, bool) {
	var (
		doc   *document.Document
		fromf *folder.Folder
	)
	if doc = document.WithID(r, document.ID(util.ParseID(docidstr))); doc == nil || doc.Archived {
		errpage.NotFound(r, user)
		return 0, false
	}
	if fromf = folder.WithID(r, doc.Folder, folder.FEditor); !user.HasPrivLevel(fromf.Editor()) {
		errpage.Forbidden(r, user)
		return 0, false
	}
	if fromf.ID() == f.ID() {
		return 0, true // nothing to do
	}
	var ud = doc.Updater(r, fromf)
	ud.Folder = f
	cleanDocName(r, ud)
	r.Transaction(func() {
		doc.Update(r, ud)
	})
	return doc.ID, true
}

// handleDeleteDocument handles a document (identified by docidstr) being
// dropped onto the trash target.  It returns whether or not the folder page
// should be emitted.  (The return will be false if an error page has been
// emitted.)
func handleDeleteDocument(r *request.Request, user *person.Person, docidstr string) bool {
	var (
		doc   *document.Document
		fromf *folder.Folder
	)
	if doc = document.WithID(r, document.ID(util.ParseID(docidstr))); doc == nil || doc.Archived {
		errpage.NotFound(r, user)
		return false
	}
	if fromf = folder.WithID(r, doc.Folder, folder.FEditor); !user.HasPrivLevel(fromf.Editor()) {
		errpage.Forbidden(r, user)
		return false
	}
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
			ud := doc.Updater(r, fromf)
			ud.Archived = true
			doc.Update(r, ud)
		}
	})
	return true
}

// handleDropFolder handles a folder (identified by fidstr) being dropped onto
// folder toparent.  It returns the ID of the folder that moved (or zero if the
// drop was rejected) and whether or not the toparent page should be emitted.
// (The latter will be false if an error page has been emitted.)
func handleDropFolder(r *request.Request, user *person.Person, toparent *folder.Folder, fidstr string) (folder.ID, bool) {
	var (
		fromparent *folder.Folder
		tomove     *folder.Folder
	)
	if tomove = folder.WithID(r, folder.ID(util.ParseID(fidstr)), folder.UpdaterFields); tomove == nil {
		errpage.NotFound(r, user)
		return 0, false
	}
	fromparent = folder.WithID(r, tomove.Parent(), folder.FName|folder.FEditor)
	if !user.HasPrivLevel(fromparent.Editor()) {
		errpage.Forbidden(r, user)
		return 0, false
	}
	// Make sure we're not moving a folder into itself or one of its
	// descendants.
	f := toparent
	for {
		if f.ID() == tomove.ID() {
			return 0, true
		}
		if f.ID() == f.Parent() {
			break
		}
		f = folder.WithID(r, f.Parent(), folder.FParent)
	}
	var uf = tomove.Updater(r, fromparent)
	uf.Parent = toparent
	cleanFolderName(r, uf)
	r.Transaction(func() {
		tomove.Update(r, uf)
	})
	return tomove.ID(), true
}

// handleDeleteFolder handles a folder (identified by fidstr) being dropped onto
// the trash target.  It returns the ID of the folder that was deleted (or zero
// if the drop was rejected) and whether or not the folder page should be
// emitted.  (The latter will be false if an error page has been emitted.)
func handleDeleteFolder(r *request.Request, user *person.Person, fidstr string) (folder.ID, bool) {
	var todel *folder.Folder

	if todel = folder.WithID(r, folder.ID(util.ParseID(fidstr)), folder.UpdaterFields); todel == nil {
		errpage.NotFound(r, user)
		return 0, false
	}
	if !user.HasPrivLevel(todel.Editor()) || folder.ExistsWithParent(r, todel.ID()) || document.ExistInFolder(r, todel.ID()) {
		errpage.Forbidden(r, user)
		return 0, false
	}
	r.Transaction(func() {
		todel.Delete(r)
	})
	return todel.ID(), true
}

// handleDropURL handles a URL (u) being dropped onto the folder f.  It returns
// the document ID of the document created for the URL, or zero if the URL is
// invalid.
func handleDropURL(r *request.Request, f *folder.Folder, u string) (did document.ID) {
	var ud = document.Updater{Folder: f, URL: u}
	if parse, err := url.Parse(u); err == nil && (parse.Scheme == "http" || parse.Scheme == "https") {
		ud.Name = getLinkTitle(u)
		cleanDocName(r, &ud)
		r.Transaction(func() {
			doc := document.Create(r, &ud)
			did = doc.ID
		})
	}
	return
}

// handleDropFile handles a file being dropped onto the folder f.  It returns
// the document ID of the document created for the file, or zero if the file
// could not be uploaded.
func handleDropFile(r *request.Request, f *folder.Folder, file *multipart.FileHeader) (did document.ID) {
	var ud = document.Updater{Folder: f, Name: file.Filename}
	var mf multipart.File
	var err error

	cleanDocName(r, &ud)
	if mf, err = file.Open(); err != nil {
		r.LogEntry.Problems.AddF("file upload failed: %s", err)
		return
	}
	defer mf.Close()
	if ud.Contents, err = io.ReadAll(mf); err != nil {
		r.LogEntry.Problems.AddF("file upload failed: %s", err)
		return
	}
	r.Transaction(func() {
		doc := document.Create(r, &ud)
		did = doc.ID
	})
	return
}

// cleanDocName ensures that the supplied document name is valid:  non-empty, no
// illegal characters, and unique within its folder.
func cleanDocName(r *request.Request, ud *document.Updater) {
	ud.Name = strings.TrimSpace(ud.Name)
	if ud.URL == "" {
		ud.Name = strings.NewReplacer("/", "", ":", "").Replace(ud.Name)
		ud.Name = strings.TrimLeft(ud.Name, ". \t\f\r\n")
	}
	if ud.Name == "" {
		ud.Name = "File"
	}
	for ud.DuplicateName(r) {
		ud.Name = incrementName(ud.Name)
	}
}

// cleanFolderName ensures that the supplied folder name is valid:  non-empty,
// no illegal characters, and unique within its parent folder.
func cleanFolderName(r *request.Request, uf *folder.Updater) {
	uf.Name = strings.TrimSpace(uf.Name)
	uf.Name = strings.NewReplacer("/", "", ":", "").Replace(uf.Name)
	uf.Name = strings.TrimLeft(uf.Name, ". \t\f\r\n")
	if uf.Name == "" {
		uf.Name = "Folder"
	}
	for {
		MakeURLName(uf)
		if uf.URLName == "" {
			uf.Name, uf.URLName = "Folder", "folder"
		}
		if !uf.DuplicateURLName(r) {
			break
		}
		uf.Name = incrementName(uf.Name)
	}
}

// MakeURLName fills in uf.URLName, generated from uf.Name.  The caller must
// verify that the result is non-empty and not a duplicate.
func MakeURLName(uf *folder.Updater) {
	uf.URLName, _, _ = transform.String(norm.NFD, uf.Name)
	uf.URLName = strings.Map(func(r rune) rune {
		if r == ' ' {
			return '-'
		}
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return -1
		}
		return r
	}, strings.ToLower(uf.URLName))
}

var seqRE = regexp.MustCompile(`(.*) (\d+)$`)

// incrementName adds one to the sequence number at the end of the supplied
// filename (prior to its extension).  If the supplied filename doesn't have a
// sequence number at the end, it adds the sequence number "2".
func incrementName(name string) string {
	var ext string
	if idx := strings.LastIndexByte(name, '.'); idx >= 0 {
		name, ext = name[:idx], name[idx:]
	}
	if match := seqRE.FindStringSubmatch(name); match != nil {
		seq, _ := strconv.Atoi(match[2])
		return fmt.Sprintf("%s %d%s", match[1], seq+1, ext)
	}
	return fmt.Sprintf("%s 2%s", name, ext)
}

// getLinkTitle fetches the specified URL and returns the contents of the
// <title> element of the resulting HTML.  If it cannot be fetched, isn't HTML,
// or doesn't have a <title> element, it returns the input URL.
func getLinkTitle(url string) string {
	ctx, cncl := context.WithTimeout(context.Background(), time.Second)
	defer cncl()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return url
	}
	req.Header.Set("Accept", "text/html")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return url
	}
	defer resp.Body.Close()
	if m, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type")); err != nil || m != "text/html" {
		return url
	}
	doc, err := xhtml.Parse(resp.Body)
	if err != nil {
		return url
	}
	for doc != nil {
		println(doc.Type)
		switch {
		case doc.Type == xhtml.DocumentNode:
			doc = doc.FirstChild
		case doc.Type != xhtml.ElementNode:
			doc = doc.NextSibling
		case doc.DataAtom == atom.Title:
			if doc.FirstChild != nil && doc.FirstChild.Type == xhtml.TextNode {
				return doc.FirstChild.Data
			}
			return url
		case doc.DataAtom == atom.Html:
			doc = doc.FirstChild
		case doc.DataAtom == atom.Head:
			doc = doc.FirstChild
		default:
			doc = doc.NextSibling
		}
	}
	return url
}
