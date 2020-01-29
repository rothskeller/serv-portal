package group

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
)

// ValidateGroup verifies that the group passes all consistency checks.
func ValidateGroup(tx *db.Tx, group *model.Group) error {
	if group.Name = strings.TrimSpace(group.Name); group.Name == "" {
		return errors.New("missing name")
	}
	for _, g := range tx.FetchGroups() {
		if g.ID != group.ID && g.Name == group.Name {
			return errors.New("duplicate name")
		}
		if g.ID != group.ID && g.Tag != "" && g.Tag == group.Tag {
			return errors.New("duplicate tag")
		}
	}
	return nil
}
