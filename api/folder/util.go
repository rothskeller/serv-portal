package folder

import (
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"

	"sunnyvaleserv.org/portal/model"
)

// CanViewFolder returns whether the specified person (which may be nil, for
// someone not logged in) can view the specified folder and its documents.
func CanViewFolder(p *model.Person, f *model.Folder) bool {
	if p == nil {
		return f.Visibility == model.FolderVisibleToPublic
	}
	if f.Visibility != model.FolderVisibleToOrg {
		return true
	}
	return p.Orgs[f.Org].PrivLevel >= model.PrivMember2
}

// canEditFolder returns whether the specified person (which may be nil, for
// someone not logged in) can edit the specified folder.
func canEditFolder(person *model.Person, folder *model.Folder) bool {
	return canEditVO(person, folder.Visibility, folder.Org)
}

// canEditVO returns whether the specified person (which may be nil, for someone
// not logged in) can edit folders with the specified visibility and org.
func canEditVO(person *model.Person, vis model.FolderVisibility, org model.Org) bool {
	if person == nil {
		return false
	}
	switch vis {
	case model.FolderVisibleToPublic:
		return person.IsAdminLeader()
	case model.FolderVisibleToSERV:
		return person.HasPrivLevel(model.PrivLeader)
	default:
		return person.Orgs[org].PrivLevel >= model.PrivLeader
	}
}

// CanShowInBrowser returns whether the document type is one that browsers can
// universally, natively display: in other words, whether the document should be
// opened in a new tab rather than downloaded.
func CanShowInBrowser(doc *model.Document) bool {
	switch {
	case doc.URL != "":
	case strings.HasSuffix(doc.Name, ".pdf"):
	case strings.HasSuffix(doc.Name, ".png"):
	case strings.HasSuffix(doc.Name, ".jpeg"):
	case strings.HasSuffix(doc.Name, ".jpg"):
		break
	default:
		return false
	}
	return true
}

// canParent returns whether a parent folder with the specified org can be a
// legal parent of a child folder with the specified org.
func canParent(pv model.FolderVisibility, po model.Org, cv model.FolderVisibility, co model.Org) bool {
	if pv > cv {
		return false
	}
	if pv == cv && pv == model.FolderVisibleToOrg {
		return po == co
	}
	return true
}

// allowedVisibilities returns the ordered list of visibility/org pairs that are
// valid for a folder with the specified parent and children, and that can be
// assigned by the specified person.
func allowedVisibilities(person *model.Person, parent *model.Folder, children []*model.Folder) (allowed []vo) {
	for _, vo := range allVOs {
		if !allowedVisibility(person, parent, children, vo.v, vo.o) {
			continue
		}
		allowed = append(allowed, vo)
	}
	return
}

// allowedVisibility returns whether the specified visibility/org pair is valid
// for a folder with the specified parent and children, and can be assigned by
// the specified person.
func allowedVisibility(
	person *model.Person, parent *model.Folder, children []*model.Folder, vis model.FolderVisibility, org model.Org,
) bool {
	if !canEditVO(person, vis, org) {
		return false
	}
	if !canParent(parent.Visibility, parent.Org, vis, org) {
		return false
	}
	for _, child := range children {
		if !canParent(vis, org, child.Visibility, child.Org) {
			return false
		}
	}
	return true
}

// nameToURL translates a folder name into a URL-sanitized version.  It removes
// accent marks and diacritics, translates spaces to hyphens, translates to
// lower case, and removes everything except letters, digits, hyphens, and
// underscores.
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

// parentURL returns the URL of the parent of the specified folder.
func parentURL(folder *model.Folder) (url string) {
	url = filepath.Dir(folder.URL)
	if url == "/" {
		url = ""
	}
	return url
}

// getURLTitle fetches the specified URL and gets the title of the resulting
// HTML document.  For any form of error, it returns the basename of the URL,
// with any name-illegal characters removed.
func getURLTitle(url string) string {
	var (
		resp *http.Response
		mt   string
		node *html.Node
		err  error
	)
	resp, err = http.Get(url)
	if err != nil {
		goto FAIL
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		goto FAIL
	}
	mt, _, err = mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil || mt != "text/html" {
		goto FAIL
	}
	node, err = html.Parse(resp.Body)
	if err != nil {
		goto FAIL
	}
	if node.Type != html.DocumentNode {
		goto FAIL
	}
	for n1 := node.FirstChild; n1 != nil; n1 = n1.NextSibling {
		if n1.Type != html.ElementNode || n1.DataAtom != atom.Html {
			continue
		}
		for n2 := n1.FirstChild; n2 != nil; n2 = n2.NextSibling {
			if n2.Type != html.ElementNode || n2.DataAtom != atom.Head {
				continue
			}
			for n3 := n2.FirstChild; n3 != nil; n3 = n3.NextSibling {
				if n3.Type != html.ElementNode || n3.DataAtom != atom.Title {
					continue
				}
				for n4 := n3.FirstChild; n4 != nil; n4 = n4.NextSibling {
					if n4.Type != html.TextNode || strings.TrimSpace(n4.Data) == "" {
						continue
					}
					return strings.Map(func(r rune) rune {
						if r == ':' || r == '/' {
							return -1
						}
						return r
					}, strings.TrimSpace(n4.Data))
				}
			}
		}
	}
FAIL:
	url = strings.Replace(filepath.Base(url), ":", "", -1)
	for len(url) > 0 && url[0] == '.' {
		url = url[1:]
	}
	if url == "" {
		url = "UNKNOWN"
	}
	return url
}
