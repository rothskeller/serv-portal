package auth

import (
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// CanAssignAnyRole returns whether the caller is allowed to assign people to
// (or unassign them from) any Role.
func CanAssignAnyRole(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.HasAny(model.PrivAssignRole)
}

// CanAssignRole returns whether the caller is allowed to assign people to (or
// unassign them from) the specified Role.
func CanAssignRole(r *util.Request, role *model.Role) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.Has(role.ID, model.PrivAssignRole)
}

// CanCreateEvents returns whether the caller is allowed to create new Event
// entries.
func CanCreateEvents(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.HasAny(model.PrivManageEvents)
}

// CanCreatePeople returns whether the caller is allowed to create new Person
// entries.
func CanCreatePeople(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.HasAny(model.PrivAssignRole)
}

// CanManageEvent returns whether the caller is allowed to edit or delete the
// specified Event.
func CanManageEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	for _, role := range event.Roles {
		if !r.Person.PrivMap.Has(role.ID, model.PrivManageEvents) {
			return false
		}
	}
	return true
}

// CanManageEvents returns whether the caller is allowed to edit or delete
// events to which the specified Role is invited.
func CanManageEvents(r *util.Request, role *model.Role) bool {
	if r.Person == nil || role.Individual {
		return false
	}
	return r.Person.PrivMap.Has(role.ID, model.PrivManageEvents)
}

// CanRecordAttendanceAtEvent returns whether the caller is allowed to record
// attendance at the specified Event.
func CanRecordAttendanceAtEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	for _, role := range event.Roles {
		if r.Person.PrivMap.Has(role.ID, model.PrivManageEvents) {
			return true
		}
	}
	return false
}

// CanViewEvent returns whether the caller is allowed to see the specified
// Event.
func CanViewEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	for _, role := range event.Roles {
		if r.Person.PrivMap.Has(role.ID, model.PrivHoldsRole) || r.Person.PrivMap.Has(role.ID, model.PrivManageEvents) {
			return true
		}
	}
	return false
}

// CanViewEventP returns whether the specified Person is allowed to see the
// specified Event.
func CanViewEventP(r *util.Request, person *model.Person, event *model.Event) bool {
	for _, role := range event.Roles {
		if person.PrivMap.Has(role.ID, model.PrivHoldsRole) || person.PrivMap.Has(role.ID, model.PrivManageEvents) {
			return true
		}
	}
	return false
}

// CanViewPerson returns whether the caller can view the argument Person.
func CanViewPerson(r *util.Request, p *model.Person) bool {
	if r.Person == nil {
		return false
	}
	for _, role := range p.Roles {
		if r.Person.PrivMap.Has(role.ID, model.PrivViewHolders) || r.Person.PrivMap.Has(role.ID, model.PrivAssignRole) {
			return true
		}
	}
	return false
}

// CanViewRole returns whether the caller can view the holders of the specified
// role.
func CanViewRole(r *util.Request, role *model.Role) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.Has(role.ID, model.PrivViewHolders) || r.Person.PrivMap.Has(role.ID, model.PrivAssignRole)
}

// HasRole returns whether the specified Person holds the specified Role
// (directly or indirectly).
func HasRole(p *model.Person, r *model.Role) bool {
	if p == nil {
		return false
	}
	return p.PrivMap.Has(r.ID, model.PrivHoldsRole)
}

// IsEnabled returns whether the specified Person is enabled.
func IsEnabled(r *util.Request, p *model.Person) bool {
	if p.PrivMap.Has(r.Tx.FetchRoleByTag(model.RoleDisabled).ID, model.PrivHoldsRole) {
		return false
	}
	return p.PrivMap.HasAny(model.PrivHoldsRole)
}

// IsWebmaster returns whether the caller is a webmaster.
func IsWebmaster(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.PrivMap.Has(r.Tx.FetchRoleByTag(model.RoleWebmaster).ID, model.PrivHoldsRole)
}

// RolesCanManageEvents returns the list of roles whose events the caller is
// allowed to manage.
func RolesCanManageEvents(r *util.Request) (roles []*model.Role) {
	if r.Person == nil {
		return nil
	}
	for _, role := range r.Tx.FetchRoles() {
		if !role.Individual && r.Person.PrivMap.Has(role.ID, model.PrivManageEvents) {
			roles = append(roles, role)
		}
	}
	return roles
}
