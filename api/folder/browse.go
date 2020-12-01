package folder

import (
	"net/url"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// getBrowseFolder returns the information needed when browsing to the specified
// folder.
func getBrowseFolder(r *util.Request, folder *model.Folder) (err error) {
	var (
		out       jwriter.Writer
		ancestors []*model.Folder
	)
	if !CanViewFolder(r.Person, folder) {
		return util.Forbidden
	}
	for af := folder; af != nil; af = r.Tx.FetchParentFolder(af) {
		ancestors = append(ancestors, af)
	}
	out.RawString(`{"ancestors":[`)
	for i := len(ancestors) - 1; i >= 0; i-- {
		if i != len(ancestors)-1 {
			out.RawByte(',')
		}
		getBrowseFolderEmitFolder(r, &out, ancestors[i])
	}
	out.RawString(`],"documents":[`)
	first := true
	for i, d := range r.Tx.FetchDocuments(folder) {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"name":`)
		out.String(d.Name)
		out.RawString(`,"url":`)
		if d.URL != "" {
			out.String(d.URL)
		} else {
			out.String("/dl/" + folder.URL + "/" + url.PathEscape(d.Name))
		}
		if CanShowInBrowser(d) {
			out.RawString(`,"newtab":true`)
		}
		out.RawByte('}')
	}
	out.RawString(`],"children":[`)
	first = true
	for _, cf := range r.Tx.FetchSubFolders(folder) {
		if !CanViewFolder(r.Person, cf) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		getBrowseFolderEmitFolder(r, &out, cf)
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

func getBrowseFolderEmitFolder(r *util.Request, out *jwriter.Writer, folder *model.Folder) {
	out.RawString(`{"name":`)
	out.String(folder.Name)
	out.RawString(`,"url":`)
	out.String(folder.URL)
	if canEditFolder(r.Person, folder) {
		out.RawString(`,"canEdit":true`)
	}
	out.RawByte('}')
}
