package folder

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

// ValidateFolder verifies the consistency of the folder.
func ValidateFolder(tx *store.Tx, folder *model.FolderNode) (err error) {
	if folder.Parent < 0 {
		return errors.New("invalid parent")
	}
	if folder.Parent == 0 {
		folder.ParentNode = nil
	} else if folder.ParentNode == nil || folder.ParentNode.ID != folder.Parent {
		if folder.ParentNode = tx.FetchFolder(folder.Parent); folder.ParentNode == nil {
			return errors.New("nonexistent parent")
		}
	}
	if folder.ID != 0 {
		for parent := folder.ParentNode; parent != nil; parent = parent.ParentNode {
			if parent.ID == folder.ID {
				return errors.New("loop in folder ancestry")
			}
		}
	}
	if folder.Name = strings.TrimSpace(folder.Name); folder.Name == "" {
		return errors.New("missing name")
	}
	if folder.ParentNode == nil {
		for _, f := range tx.FetchFolders() {
			if f.ID != folder.ID && f.Name == folder.Name {
				return errors.New("duplicate name")
			}
		}
	} else {
		for _, f := range folder.ParentNode.ChildNodes {
			if f.ID != folder.ID && f.Name == folder.Name {
				return errors.New("duplicate name")
			}
		}
	}
	if folder.Group < 0 {
		return errors.New("invalid group")
	}
	if folder.Group > 0 && tx.Authorizer().FetchGroup(folder.Group) == nil {
		return errors.New("nonexistent group")
	}
	if folder.Org != model.OrgNone2 {
		if folder.Public {
			return errors.New("public folder with org")
		}
		var found = false
		for _, o := range model.AllOrgs {
			if folder.Org == o {
				found = true
				break
			}
		}
		if !found {
			return errors.New("invalid org")
		}
	}
	if folder.ParentNode != nil && folder.ParentNode.Org != model.OrgNone2 && folder.Org != folder.ParentNode.Org {
		return errors.New("folder has different org than parent")
	}
	if folder.ParentNode != nil && !folder.ParentNode.Public && folder.Public {
		return errors.New("public folder under non-public parent")
	}
	folder.Approvals = 0
	for _, cf := range folder.ChildNodes {
		folder.Approvals += cf.Approvals
	}
	for _, doc := range folder.Documents {
		if doc.NeedsApproval {
			folder.Approvals++
		}
		if doc.ID < 1 {
			return errors.New("invalid document ID")
		}
		if doc.Name = strings.TrimSpace(doc.Name); doc.Name == "" {
			return errors.New("missing document name")
		}
		if p := tx.FetchPerson(doc.PostedBy); p == nil {
			return errors.New("nonexistent postedBy person")
		}
		if doc.PostedAt.IsZero() {
			return errors.New("missing postedAt")
		}
		for _, doc2 := range folder.Documents {
			if doc == doc2 {
				continue
			}
			if doc.ID == doc2.ID {
				return errors.New("duplicate document ID")
			}
			if doc.Name == doc2.Name && doc.NeedsApproval == doc2.NeedsApproval && (doc.PostedBy == doc2.PostedBy || !doc.NeedsApproval) {
				return errors.New("duplicate document name")
			}
		}
	}
	return nil
}
