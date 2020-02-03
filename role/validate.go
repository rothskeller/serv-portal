package role

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/authz"
	"sunnyvaleserv.org/portal/model"
)

// ValidateRole verifies that the role passes all consistency checks.
func ValidateRole(auth *authz.Authorizer, role *model.Role) error {
	if role.Name = strings.TrimSpace(role.Name); role.Name == "" {
		return errors.New("missing name")
	}
	for _, r := range auth.FetchRoles(auth.AllRoles()) {
		if r.ID != role.ID && r.Name == role.Name {
			return errors.New("duplicate name")
		}
		if r.ID != role.ID && r.Tag != "" && r.Tag == role.Tag {
			return errors.New("duplicate tag")
		}
	}
	return nil
}
