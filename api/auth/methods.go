package auth

import "rothskeller.net/serv/util"

import "rothskeller.net/serv/model"

/*
// CanAssignAnyRole returns whether the caller is allowed to assign people to
// (or unassign them from) any Role.
func CanAssignAnyRole(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.HasAny(model.PrivAssignRole)
}

// CanAssignRole returns whether the caller is allowed to assign people to (or
// unassign them from) the specified Role.
func CanAssignRole(r *util.Request, role *model.Role) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.Has(role.ID, model.PrivAssignRole)
}
*/

// CanCreateEvents returns whether the caller is allowed to create new Event
// entries.
func CanCreateEvents(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.HasAny(model.PrivManageEvents) || IsWebmaster(r)
}

// CanCreatePeople returns whether the caller is allowed to create new Person
// entries.
func CanCreatePeople(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.HasAny(model.PrivManageMembers) || IsWebmaster(r)
}

// CanManageEvent returns whether the caller is allowed to edit or delete the
// specified Event.
func CanManageEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	if IsWebmaster(r) {
		return true
	}
	for _, gid := range event.Groups {
		if !r.Person.Privileges.Has(r.Tx.FetchGroup(gid), model.PrivManageEvents) {
			return false
		}
	}
	return true
}

// CanManageEvents returns whether the caller is allowed to edit or delete
// events to which the specified Group is invited.
func CanManageEvents(r *util.Request, group *model.Group) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.Has(group, model.PrivManageEvents) || IsWebmaster(r)
}

// CanRecordAttendanceAtEvent returns whether the caller is allowed to record
// attendance at the specified Event.
func CanRecordAttendanceAtEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	for _, gid := range event.Groups {
		if r.Person.Privileges.Has(r.Tx.FetchGroup(gid), model.PrivManageEvents) {
			return true
		}
	}
	return IsWebmaster(r)
}

// CanViewContactInfo returns whether the caller can view contact information
// for the argument Person.
func CanViewContactInfo(r *util.Request, p *model.Person) bool {
	if r.Person == nil {
		return false
	}
	for _, group := range r.Tx.FetchGroups() {
		if p.Privileges.Has(group, model.PrivMember) && r.Person.Privileges.Has(group, model.PrivViewContactInfo) {
			return true
		}
	}
	return IsWebmaster(r)
}

// CanViewEvent returns whether the caller is allowed to see the specified
// Event.
func CanViewEvent(r *util.Request, event *model.Event) bool {
	if r.Person == nil {
		return false
	}
	for _, gid := range event.Groups {
		group := r.Tx.FetchGroup(gid)
		if r.Person.Privileges.Has(group, model.PrivMember) || r.Person.Privileges.Has(group, model.PrivManageEvents) {
			return true
		}
	}
	return IsWebmaster(r)
}

// CanViewEventP returns whether the specified Person is allowed to see the
// specified Event.
func CanViewEventP(r *util.Request, person *model.Person, event *model.Event) bool {
	for _, gid := range event.Groups {
		group := r.Tx.FetchGroup(gid)
		if person.Privileges.Has(group, model.PrivMember) || person.Privileges.Has(group, model.PrivManageEvents) {
			return true
		}
	}
	return false
}

// CanViewGroup returns whether the caller can view the roster of the argument
// Group.
func CanViewGroup(r *util.Request, g *model.Group) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.Has(g, model.PrivViewMembers) || IsWebmaster(r)
}

// CanViewPerson returns whether the caller can view the argument Person.
func CanViewPerson(r *util.Request, p *model.Person) bool {
	if r.Person == nil {
		return false
	}
	for _, group := range r.Tx.FetchGroups() {
		if p.Privileges.Has(group, model.PrivMember) && r.Person.Privileges.Has(group, model.PrivViewMembers) {
			return true
		}
	}
	return IsWebmaster(r)
}

/*
// CanViewRole returns whether the caller can view the holders of the specified
// role.
func CanViewRole(r *util.Request, role *model.Role) bool {
	if r.Person == nil {
		return false
	}
	return r.Person.Privileges.Has(role.ID, model.PrivViewHolders) || r.Person.Privileges.Has(role.ID, model.PrivAssignRole)
}

*/

// GroupsCanManageEvents returns the list of groups whose events the caller is
// allowed to manage.
func GroupsCanManageEvents(r *util.Request) (groups []*model.Group) {
	if r.Person == nil {
		return nil
	}
	for _, group := range r.Tx.FetchGroups() {
		if r.Person.Privileges.Has(group, model.PrivManageEvents) || IsWebmaster(r) {
			groups = append(groups, group)
		}
	}
	return groups
}

// HasRole returns whether the specified Person holds the specified Role.
func HasRole(p *model.Person, r *model.Role) bool {
	if p == nil {
		return false
	}
	for _, role := range p.Roles {
		if role == r.ID {
			return true
		}
	}
	return false
}

// IsMember returns whether the specified Person is a member of the specified
// group.
func IsMember(p *model.Person, g *model.Group) bool {
	if p == nil {
		return false
	}
	return p.Privileges.Has(g, model.PrivMember)
}

// IsEnabled returns whether the specified Person is enabled.
func IsEnabled(r *util.Request, p *model.Person) bool {
	if p.Privileges.Has(r.Tx.FetchGroupByTag(model.GroupDisabled), model.PrivMember) {
		return false
	}
	if p.Privileges.HasAny(model.PrivMember) {
		return true
	}
	return HasRole(p, r.Tx.FetchRoleByTag(model.RoleWebmaster))
}

// IsWebmaster returns whether the caller is a webmaster.
func IsWebmaster(r *util.Request) bool {
	if r.Person == nil {
		return false
	}
	return HasRole(r.Person, r.Tx.FetchRoleByTag(model.RoleWebmaster))
}
