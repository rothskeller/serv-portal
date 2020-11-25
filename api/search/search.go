package search

import (
	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util"
)

// GetSearch handles GET /api/search?q= requests.
func GetSearch(r *util.Request) error {
	var (
		count int
		out   jwriter.Writer
		err   error
		first = true
	)
	out.RawString(`{"results":[`)
	err = r.Tx.Search(r.FormValue("q"), func(obj interface{}) bool {
		// The first switch just checks the user's rights to see the
		// object.  It returns from the function if they can't see it.
		switch tobj := obj.(type) {
		case store.FolderAndDocument:
			if tobj.Folder.Group != 0 && !r.Auth.MemberG(tobj.Folder.Group) && !r.Auth.CanAG(model.PrivManageFolders, tobj.Folder.Group) {
				return true
			}
			if tobj.Document.NeedsApproval {
				return true
			}
		case *model.Event:
			break
		case *model.FolderNode:
			if tobj.Group != 0 && !r.Auth.MemberG(tobj.Group) && !r.Auth.CanAG(model.PrivManageFolders, tobj.Group) {
				return true
			}
		case *model.Group:
			if !r.Auth.CanAG(model.PrivViewMembers, tobj.ID) {
				return true
			}
		case *model.Person:
			if tobj == r.Person || r.Person.HasPrivLevel(model.PrivLeader) {
				break
			}
			found := false
			for _, o := range model.AllOrgs {
				if tobj.Orgs[o].PrivLevel >= model.PrivMember2 && r.Person.Orgs[o].PrivLevel >= model.PrivMember2 {
					found = true
					break
				}
			}
			if !found {
				return true
			}
		case *model.Role2:
			switch {
			case r.Person.Orgs[tobj.Org].PrivLevel >= model.PrivMember2:
			case r.Person.HasPrivLevel(model.PrivLeader):
				break
			default:
				return true
			}
		case *model.TextMessage:
			var canView bool
			for _, list := range r.Tx.FetchLists() {
				if list.Type == model.ListSMS && list.People[r.Person.ID]&model.ListSender != 0 {
					canView = true
					break
				}
			}
			if !canView {
				return true
			}
		default:
			panic("unexpected type of search result")
		}
		// OK, now we know they can see the object.
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		// The second switch writes out the details of the object.
		switch tobj := obj.(type) {
		case store.FolderAndDocument:
			var path []string
			for p := tobj.Folder; p != nil; p = p.ParentNode {
				path = append(path, p.Name)
			}
			out.RawString(`{"type":"document","path":[`)
			for i := len(path) - 1; i >= 0; i-- {
				if i != len(path)-1 {
					out.RawByte(',')
				}
				out.String(path[i])
			}
			out.RawString(`],"folderID":`)
			out.Int(int(tobj.Folder.ID))
			out.RawString(`,"documentID":`)
			out.Int(int(tobj.Document.ID))
			out.RawString(`,"name":`)
			out.String(tobj.Document.Name)
			out.RawByte('}')
		case *model.Event:
			out.RawString(`{"type":"event","id":`)
			out.Int(int(tobj.ID))
			out.RawString(`,"date":`)
			out.String(tobj.Date)
			out.RawString(`,"name":`)
			out.String(tobj.Name)
			out.RawByte('}')
		case *model.FolderNode:
			var path []string
			for p := tobj.ParentNode; p != nil; p = p.ParentNode {
				path = append(path, p.Name)
			}
			out.RawString(`{"type":"folder","path":[`)
			for i := len(path) - 1; i >= 0; i-- {
				if i != len(path)-1 {
					out.RawByte(',')
				}
				out.String(path[i])
			}
			out.RawString(`],"id":`)
			out.Int(int(tobj.ID))
			out.RawString(`,"name":`)
			out.String(tobj.Name)
			out.RawByte('}')
		case *model.Person:
			out.RawString(`{"type":"person","id":`)
			out.Int(int(tobj.ID))
			out.RawString(`,"informalName":`)
			out.String(tobj.InformalName)
			out.RawByte('}')
		case *model.Role2:
			out.RawString(`{"type":"role","id":`)
			out.Int(int(tobj.ID))
			out.RawString(`,"name":`)
			out.String(tobj.Name)
			out.RawByte('}')
		case *model.TextMessage:
			out.RawString(`{"type":"textMessage","id":`)
			out.Int(int(tobj.ID))
			out.RawString(`,"sender":`)
			out.String(r.Tx.FetchPerson(tobj.Sender).InformalName)
			out.RawString(`,"timestamp":`)
			out.Raw(tobj.Timestamp.MarshalJSON())
			out.RawByte('}')
		}
		// Keep going until we have 100 results.
		count++
		return count < 100
	})
	out.RawByte(']')
	if err != nil {
		out.RawString(`,"error":`)
		out.String(err.Error())
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
