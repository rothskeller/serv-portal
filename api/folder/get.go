package folder

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetFolder handles GET /api/folders/$id requests, where $id may be 0 to get
// the virtual parent of the top-level folders.
func GetFolder(r *util.Request, idstr string) (err error) {
	var (
		folderID model.FolderID
		folder   *model.Folder
		canEdit  bool
		out      jwriter.Writer
	)
	if idstr != "0" {
		folderID = model.FolderID(util.ParseID(idstr))
		if folder = r.Tx.FetchFolder(folderID); folder == nil {
			return util.NotFound
		}
		if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
			return util.Forbidden
		}
	}
	out.RawByte('{')
	if folder != nil {
		out.RawString(`"id":`)
		out.Int(int(folder.ID))
		if folder.Parent != 0 {
			out.RawString(`,"parent":{"id":`)
			out.Int(int(folder.Parent))
			out.RawString(`,"name":`)
			out.String(r.Tx.FetchFolder(folder.Parent).Name)
			out.RawByte('}')
		}
		out.RawString(`,"group":`)
		out.Int(int(folder.Group))
		out.RawString(`,"name":`)
		out.String(folder.Name)
		out.RawString(`,"documents":[`)
		for i, d := range folder.Documents {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(d.ID))
			out.RawString(`,"name":`)
			out.String(d.Name)
			out.RawByte('}')
		}
		out.RawByte(']')
	} else {
		out.RawString(`"group":0`)
	}
	first := true
	for _, f := range r.Tx.FetchFolders() {
		if f.Parent != folderID {
			continue
		}
		if f.Group != 0 && !r.Auth.MemberG(f.Group) && !r.Auth.CanAG(model.PrivManageFolders, f.Group) {
			continue
		}
		if first {
			out.RawString(`,"children":[`)
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(f.ID))
		out.RawString(`,"name":`)
		out.String(f.Name)
		out.RawString(`,"group":`)
		out.Int(int(f.Group))
		out.RawByte('}')
	}
	if !first {
		out.RawByte(']')
	}
	if (folder == nil || folder.Group == 0) && r.Auth.IsWebmaster() {
		canEdit = true
	} else if folder != nil && folder.Group != 0 && r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		canEdit = true
	}
	if canEdit {
		out.RawString(`,"canEdit":true,"allowedGroups":[`)
		first := true
		if r.Auth.IsWebmaster() {
			out.RawString(`{"id":0,"name":"Public"}`)
			first = false
		}
		for _, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivManageFolders)) {
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(g.ID))
			out.RawString(`,"name":`)
			out.String(g.Name)
			out.RawByte('}')
		}
		out.RawString(`],"allowedParents":[`)
		folders := r.Tx.FetchFolders()
		if r.Auth.IsWebmaster() {
			out.RawString(`{"id":0,"name":"Files"}`)
			emitAllowedParents(r, folders, &out, 0, 1, false)
		} else {
			emitAllowedParents(r, folders, &out, 0, 0, true)
		}
		out.RawByte(']')
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

func emitAllowedParents(r *util.Request, folders []*model.Folder, out *jwriter.Writer, parent model.FolderID, indent int, first bool) bool {
	for _, f := range folders {
		if f.Parent != parent {
			continue
		}
		if f.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, f.Group) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(f.ID))
		out.RawString(`,"name":`)
		out.String(strings.Repeat("\u00A0", indent*4) + f.Name)
		out.RawString(`,"indent":`)
		out.Int(indent)
		if f.Group == 0 && !r.Auth.IsWebmaster() {
			out.RawString(`,"disabled":true`)
		}
		out.RawByte('}')
		emitAllowedParents(r, folders, out, f.ID, indent+1, first)
	}
	return first
}

// GetDocument handles GET /api/folders/$fid/$did requests.
func GetDocument(r *util.Request, fidstr, didstr string) (err error) {
	var (
		folder *model.Folder
		docID  model.DocumentID
		doc    *model.Document
		fh     *os.File
		stat   os.FileInfo
	)
	if folder = r.Tx.FetchFolder(model.FolderID(util.ParseID(fidstr))); folder == nil {
		return util.NotFound
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
	if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	fh = r.Tx.FetchDocument(folder, docID)
	r.Tx.Commit()
	if stat, err = fh.Stat(); err != nil {
		return err
	}
	if strings.HasSuffix(doc.Name, ".pdf") {
		r.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", doc.Name))
	} else {
		r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.Name))
	}
	http.ServeContent(r, r.Request, doc.Name, stat.ModTime(), fh)
	return nil
}
