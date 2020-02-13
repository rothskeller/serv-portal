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
		out      jwriter.Writer
	)
	if idstr != "0" {
		folderID = model.FolderID(util.ParseID(idstr))
		if folder = r.Tx.FetchFolder(folderID); folder == nil {
			return util.NotFound
		}
		if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.IsWebmaster() {
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
	}
	first := true
	for _, f := range r.Tx.FetchFolders() {
		if f.Parent != folderID {
			continue
		}
		if f.Group != 0 && !r.Auth.MemberG(f.Group) && !r.Auth.IsWebmaster() {
			continue
		}
		switch {
		case first && folder == nil:
			out.RawString(`"children":[`)
		case first && folder != nil:
			out.RawString(`,"children":[`)
		case !first:
			out.RawByte(',')
		}
		first = false
		out.RawString(`{"id":`)
		out.Int(int(f.ID))
		out.RawString(`,"name":`)
		out.String(f.Name)
		out.RawByte('}')
	}
	if !first {
		out.RawByte(']')
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
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
	if folder.Group != 0 && !r.Auth.MemberG(folder.Group) && !r.Auth.IsWebmaster() {
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
