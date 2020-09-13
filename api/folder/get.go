package folder

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetPath handles GET /api/folders?path= requests.  The supplied path may
// address either a folder or a document within a folder.  If it addresses a
// folder, this API returns the folder details.  If it addresses a document,
// this API returns the URL for retrieval of the document.
func GetPath(r *util.Request) (err error) {
	var (
		folder *model.FolderNode
		doc    *model.Document
		path   = strings.Split(r.FormValue("path"), "/")
	)
	folder = r.Tx.FetchRootFolder()
PATH:
	for len(path) > 0 {
		for _, c := range folder.ChildNodes {
			if nameToURL(c.Name) == path[0] {
				if c.Group != 0 && r.Person == nil {
					return util.Forbidden
				}
				if c.Group != 0 && !r.Auth.MemberG(c.Group) && !r.Auth.CanAG(model.PrivManageFolders, c.Group) {
					return util.Forbidden
				}
				folder = c
				path = path[1:]
				continue PATH
			}
		}
		if len(path) == 1 {
			for _, d := range folder.Documents {
				if d.Name == path[0] {
					doc = d
					path = path[1:]
					continue PATH
				}
			}
		}
		return util.NotFound
	}
	return getFolder(r, folder, doc)
}

// GetFolder handles GET /api/folders/$id requests, where $id may be 0 to get
// the virtual parent of the top-level folders.
func GetFolder(r *util.Request, idstr string) (err error) {
	var (
		folderID model.FolderID
		folder   *model.FolderNode
	)
	folderID = model.FolderID(util.ParseID(idstr))
	if folder = r.Tx.FetchFolder(folderID); folder == nil {
		return util.NotFound
	}
	if folder.Group != 0 && r.Person == nil {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	return getFolder(r, folder, nil)
}

func getFolder(r *util.Request, folder *model.FolderNode, doc *model.Document) (err error) {
	var (
		canEdit    bool
		canApprove bool
		out        jwriter.Writer
	)
	out.RawByte('{')
	canApprove = r.Auth.CanA(model.PrivManageFolders)
	out.RawString(`"id":`)
	out.Int(int(folder.ID))
	if folder.ParentNode != nil {
		out.RawString(`,"parent":{"id":`)
		out.Int(int(folder.Parent))
		out.RawString(`,"name":`)
		out.String(folder.ParentNode.Name)
		out.RawString(`,"url":`)
		out.String(folderURL(folder.ParentNode))
		out.RawByte('}')
	}
	out.RawString(`,"group":`)
	out.Int(int(folder.Group))
	out.RawString(`,"name":`)
	out.String(folder.Name)
	out.RawString(`,"url":`)
	out.String(folderURL(folder))
	out.RawString(`,"documents":[`)
	first := true
	for _, d := range folder.Documents {
		if d.NeedsApproval && !canApprove {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(d.ID))
		out.RawString(`,"name":`)
		out.String(d.Name)
		if d.NeedsApproval {
			out.RawString(`,"needsApproval":true`)
		}
		out.RawByte('}')
	}
	out.RawByte(']')
	first = true
	for _, cf := range folder.ChildNodes {
		if cf.Group != 0 && !r.Auth.MemberG(cf.Group) && !r.Auth.CanAG(model.PrivManageFolders, cf.Group) {
			continue
		}
		if first {
			out.RawString(`,"children":[`)
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(cf.ID))
		out.RawString(`,"name":`)
		out.String(cf.Name)
		out.RawString(`,"url":`)
		out.String(folderURL(cf))
		out.RawString(`,"group":`)
		out.Int(int(cf.Group))
		if cf.Approvals > 0 && canApprove {
			out.RawString(`,"approvals":`)
			out.Int(cf.Approvals)
		}
		out.RawByte('}')
	}
	if !first {
		out.RawByte(']')
	}
	if folder.Group == 0 && r.Auth.IsWebmaster() {
		canEdit = true
	} else if folder.Group != 0 && r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
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
		emitAllowedParents(r, r.Tx.FetchRootFolder(), &out, 0, true)
		out.RawByte(']')
	}
	if r.Auth.CanA(model.PrivManageFolders) {
		out.RawString(`,"canAdd":true`)
	}
	if doc != nil {
		out.RawString(`,"docDownload":`)
		out.String(fmt.Sprintf("/api/folders/%d/%d", folder.ID, doc.ID))
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

func emitAllowedParents(r *util.Request, root *model.FolderNode, out *jwriter.Writer, indent int, first bool) bool {
	for _, f := range root.ChildNodes {
		if f.Group != 0 && !r.Auth.MemberG(f.Group) && !r.Auth.CanAG(model.PrivManageFolders, f.Group) {
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
		} else if f.Group != 0 && !r.Auth.CanAG(model.PrivManageFolders, f.Group) {
			out.RawString(`,"disabled":true`)
		}
		out.RawByte('}')
		emitAllowedParents(r, f, out, indent+1, first)
	}
	return first
}

// GetDocument handles GET /api/folders/$fid/$did requests.
func GetDocument(r *util.Request, fidstr, didstr string) (err error) {
	var (
		folder *model.FolderNode
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
	if folder.Group != 0 && r.Person == nil {
		return util.Forbidden
	}
	if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.CanAG(model.PrivManageFolders, folder.Group) {
		return util.Forbidden
	}
	if doc.NeedsApproval {
		switch {
		case doc.PostedBy == r.Person.ID:
		case folder.Group != 0 && r.Auth.CanAG(model.PrivManageFolders, folder.Group):
		case r.Auth.IsWebmaster():
			break
		default:
			return util.Forbidden
		}
	}
	fh = r.Tx.FetchDocument(folder, docID)
	r.Tx.Commit()
	if stat, err = fh.Stat(); err != nil {
		return err
	}
	r.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", doc.Name))
	r.Header().Set("Cache-Control", "no-cache")
	http.ServeContent(r, r.Request, doc.Name, stat.ModTime(), fh)
	return nil
}

func folderURL(folder *model.FolderNode) string {
	if folder.ParentNode == nil {
		return ""
	}
	return folderURL(folder.ParentNode) + "/" + nameToURL(folder.Name)
}

func nameToURL(name string) string {
	result, _, _ := transform.String(norm.NFD, name)
	return strings.Map(func(r rune) rune {
		if r == ' ' {
			return '-'
		}
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '-' && r != '_' {
			return -1
		}
		return r
	}, strings.ToLower(result))
}
