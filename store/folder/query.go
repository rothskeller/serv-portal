package folder

import (
	"net/url"
	"strings"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// ExistsWithParent returns whether any folder exists with the specified folder
// as its parent.
func ExistsWithParent(storer phys.Storer, id ID) (found bool) {
	phys.SQL(storer, "SELECT 1 FROM folder WHERE parent=?", func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		found = stmt.Step()
	})
	return found
}

// Path returns the URL path of the receiver folder.  It is an absolute path
// without hostname.
func (f *Folder) Path(storer phys.Storer) string {
	var parts = []string{f.URLName()}

	for f.Parent() != RootID {
		f = WithID(storer, f.Parent(), FParent|FURLName)
		parts = append(parts, f.URLName())
	}
	parts = append(parts, "/files")
	for i := 0; i < len(parts)/2; i++ {
		parts[i], parts[len(parts)-i-1] = parts[len(parts)-i-1], parts[i]
	}
	return strings.Join(parts, "/")
}

// FolderPath returns the list of Folder structures on the path to the receiver
// folder.
func (f *Folder) FolderPath(storer phys.Storer, fields Fields) (flist []*Folder) {
	fields |= FParent
	flist = []*Folder{f}
	for f.Parent() != f.ID() {
		f = WithID(storer, f.Parent(), fields)
		flist = append(flist, f)
	}
	for i := 0; i < len(flist)/2; i++ {
		flist[i], flist[len(flist)-i-1] = flist[len(flist)-i-1], flist[i]
	}
	return flist
}

var withIDSQLCache map[Fields]string

// WithID returns the folder with the specified ID, or nil if it does not exist.
func WithID(storer phys.Storer, id ID, fields Fields) (f *Folder) {
	if withIDSQLCache == nil {
		withIDSQLCache = make(map[Fields]string)
	}
	if _, ok := withIDSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM folder f WHERE f.id=?")
		withIDSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, withIDSQLCache[fields], func(stmt *phys.Stmt) {
		stmt.BindInt(int(id))
		if stmt.Step() {
			f = new(Folder)
			f.Scan(stmt, fields)
			f.id = id
			f.fields |= FID
		}
	})
	return f
}

var withPathSQLCache map[Fields]string

// WithPath returns the folders on the specified path, in order.  The path must
// be an absolute path without hostname (i.e., starting with "/files"), and must
// refer to either a folder or a document.
//
// If the path refers to an existing folder, flist will contain a list of Folder
// structures mirroring the path, starting with the root folder and ending with
// the referenced folder, and docname will be an empty string.
//
// If the path, up through its second-to-last component, refers to an existing
// folder, flist will contain a list of Folder structures mirroring the path to
// that point, and docname will contain the final component of the path.  The
// existence of a document by that name is not checked.
//
// If neither of those cases applies, WithPath returns nil, "".
func WithPath(storer phys.Storer, path string, fields Fields) (flist []*Folder, docname string) {
	var (
		parts []string
		id    = RootID
	)
	fields |= FID
	if withPathSQLCache == nil {
		withPathSQLCache = make(map[Fields]string)
	}
	if _, ok := withPathSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString(`SELECT `)
		ColumnList(&sb, fields)
		sb.WriteString(` FROM folder f WHERE f.parent=? AND f.url_name=?`)
		withPathSQLCache[fields] = sb.String()
	}
	if !strings.HasPrefix(path, "/files") {
		return nil, ""
	}
	path = strings.Trim(path[6:], "/")
	phys.SQL(storer, withPathSQLCache[fields], func(stmt *phys.Stmt) {
		var f Folder
		stmt.BindInt(int(RootID))
		stmt.BindText("")
		stmt.Step()
		f.Scan(stmt, fields)
		f.fields |= FParent
		f.parent = RootID
		flist = append(flist, &f)
	})
	if path == "" {
		return flist, ""
	}
	parts = strings.Split(path, "/")
	for i, part := range parts {
		found := false
		phys.SQL(storer, withPathSQLCache[fields], func(stmt *phys.Stmt) {
			stmt.BindInt(int(id))
			stmt.BindText(part)
			if stmt.Step() {
				var f Folder
				f.Scan(stmt, fields)
				f.parent, f.fields = id, f.fields|FParent
				flist = append(flist, &f)
				found = true
				id = f.ID()
			}
		})
		if !found {
			if i == len(parts)-1 {
				if docname, err := url.PathUnescape(part); err == nil {
					return flist, docname
				}
			}
			return nil, ""
		}
	}
	return flist, ""
}

var allWithParentSQLCache map[Fields]string

// AllWithParent fetches each of the folders with the specified parent folder,
// in order by name.
func AllWithParent(storer phys.Storer, parent ID, fields Fields, fn func(*Folder)) {
	if allWithParentSQLCache == nil {
		allWithParentSQLCache = make(map[Fields]string)
	}
	if _, ok := allWithParentSQLCache[fields]; !ok {
		var sb strings.Builder
		sb.WriteString("SELECT ")
		ColumnList(&sb, fields)
		sb.WriteString(" FROM folder f WHERE parent=?1 AND id!=?1 ORDER BY name")
		// The id!=?1 prevents retrieving the root folder when the
		// requested parent is the root folder.
		allWithParentSQLCache[fields] = sb.String()
	}
	phys.SQL(storer, allWithParentSQLCache[fields], func(stmt *phys.Stmt) {
		var f Folder

		stmt.BindInt(int(parent))
		for stmt.Step() {
			f.Scan(stmt, fields)
			f.parent = parent
			f.fields |= FParent
			fn(&f)
		}
	})
}
