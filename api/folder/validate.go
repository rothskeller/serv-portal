package folder

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

// ValidateFolder verifies the consistency of the folder.
func ValidateFolder(tx *store.Tx, folder *model.Folder) (err error) {
	if folder.Parent < 0 || (folder.ID != 0 && folder.Parent == folder.ID) {
		return errors.New("invalid parent")
	}
	// TODO check for folder ancestry loop
	if folder.Parent > 0 && tx.FetchFolder(folder.Parent) == nil {
		return errors.New("nonexistent parent")
	}
	if folder.Name = strings.TrimSpace(folder.Name); folder.Name == "" {
		return errors.New("missing name")
	}
	for _, f := range tx.FetchFolders() {
		if f.ID != folder.ID && f.Parent == folder.Parent && f.Name == folder.Name {
			return errors.New("duplicate name")
		}
	}
	if folder.Group < 0 {
		return errors.New("invalid group")
	}
	if folder.Group > 0 && tx.Authorizer().FetchGroup(folder.Group) == nil {
		return errors.New("nonexistent group")
	}
	for _, doc := range folder.Documents {
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
			if doc.Name == doc2.Name {
				return errors.New("duplicate document name")
			}
		}
	}
	return nil
}
