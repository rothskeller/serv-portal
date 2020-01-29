package role

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
)

// ValidateRole verifies that the role passes all consistency checks.
func ValidateRole(tx *db.Tx, role *model.Role) error {
	if role.Name = strings.TrimSpace(role.Name); role.Name == "" {
		return errors.New("missing name")
	}
	for _, r := range tx.FetchRoles() {
		if r.ID != role.ID && r.Name == role.Name {
			return errors.New("duplicate name")
		}
	}
	for _, g := range role.Privileges.GroupsWithAny() {
		var group = tx.FetchGroup(g)
		if group == nil {
			return errors.New("invalid group in privileges")
		}
		privs := role.Privileges.Get(group)
		if privs&(model.PrivViewMembers|model.PrivViewContactInfo) == model.PrivViewContactInfo {
			return errors.New("can't have priv contact without having priv roster")
		}
		if privs&(model.PrivSendTextMessages|model.PrivViewContactInfo) == model.PrivSendTextMessages {
			return errors.New("can't have priv texts without having priv contact")
		}
		if privs&(model.PrivManageMembers|model.PrivViewMembers) == model.PrivManageMembers {
			return errors.New("can't have priv admin without having priv roster")
		}
		if privs&(model.PrivManageEvents|model.PrivViewMembers) == model.PrivManageEvents {
			return errors.New("can't have priv events without having priv roster")
		}
	}
	return nil
}
