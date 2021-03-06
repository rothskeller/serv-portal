package role

import (
	"errors"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

// ValidateRole verifies that the role passes all consistency checks.
func ValidateRole(tx *store.Tx, role *model.Role) error {
	var maxprio int

	if role.ID == model.Webmaster {
		return errors.New("webmaster role cannot be changed")
	}
	if role.Name = strings.TrimSpace(role.Name); role.Name == "" {
		return errors.New("missing name")
	}
	role.Title = strings.TrimSpace(role.Title)
	if !role.Org.Valid() {
		return errors.New("invalid org")
	}
	if !role.PrivLevel.Valid() && role.PrivLevel != model.PrivNone {
		return errors.New("invalid privLevel")
	}
	if role.PrivLevel != model.PrivNone && role.Org == model.OrgNone {
		return errors.New("missing org")
	}
	if role.Priority < 0 {
		return errors.New("invalid priority")
	}
	if _, ok := role.Implies[role.ID]; ok {
		return errors.New("role cannot imply itself")
	}
	for irid, direct := range role.Implies {
		if !direct {
			continue // they'll get recalculated later
		}
		if irid == model.Webmaster {
			return errors.New("webmaster role cannot be implied")
		}
		ir := tx.FetchRole(irid)
		if ir == nil {
			return errors.New("nonexistent role in implies")
		}
		if _, ok := ir.Implies[role.ID]; ok {
			return errors.New("cycle in role implies")
		}
	}
	for lid, rtl := range role.Lists {
		if rtl == 0 {
			delete(role.Lists, lid)
			continue
		}
		list := tx.FetchList(lid)
		if list == nil {
			return errors.New("nonexistent list in lists")
		}
		if _, ok := model.ListSubModelNames[rtl.SubModel()]; !ok && rtl.SubModel() != model.ListNoSub {
			return errors.New("invalid subModel in lists")
		}
	}
	for _, r := range tx.FetchRoles() {
		if r == role {
			continue
		}
		if r.Name == role.Name {
			return errors.New("duplicate name")
		}
		if r.Title == role.Title && r.Title != "" {
			return errors.New("duplicate title")
		}
		if r.Priority == role.Priority {
			return errors.New("duplicate priority")
		}
		if r.Priority > maxprio {
			maxprio = r.Priority
		}
	}
	if role.Priority == 0 {
		role.Priority = maxprio + 1
	}
	return nil
}
