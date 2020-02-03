package group

import (
	"errors"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/authz"
	"sunnyvaleserv.org/portal/model"
)

var groupEmailRE = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

// ValidateGroup verifies that the group passes all consistency checks.
func ValidateGroup(auth *authz.Authorizer, group *model.Group) error {
	if group.Name = strings.TrimSpace(group.Name); group.Name == "" {
		return errors.New("missing name")
	}
	group.Email = strings.ToLower(strings.TrimSpace(group.Email))
	if group.Email != "" && !groupEmailRE.MatchString(group.Email) {
		return errors.New("invalid email")
	}
	for _, g := range auth.FetchGroups(auth.AllGroups()) {
		if g.ID != group.ID && g.Name == group.Name {
			return errors.New("duplicate name")
		}
		if g.ID != group.ID && g.Tag != "" && g.Tag == group.Tag {
			return errors.New("duplicate tag")
		}
		if g.ID != group.ID && g.Email != "" && g.Email == group.Email {
			return errors.New("duplicate email")
		}
	}
	return nil
}
