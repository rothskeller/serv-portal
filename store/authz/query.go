package authz

import (
	"math/bits"
	"sort"

	"sunnyvaleserv.org/portal/model"
)

// ActionsRG returns the bitmask of privileges that the specified role has on
// the specified group.
func (a *Authorizer) ActionsRG(role model.RoleID, group model.GroupID) model.Privilege {
	return a.rolePrivs[int(role)*len(a.groups)+int(group)]
}

// AllGroups returns the list of all roles.
func (a *Authorizer) AllGroups() (groups []model.GroupID) {
	for gid := range a.groups {
		if gid != 0 && a.groups[model.GroupID(gid)].ID == model.GroupID(gid) {
			groups = append(groups, model.GroupID(gid))
		}
	}
	return groups
}

// AllRoles returns the list of all roles.
func (a *Authorizer) AllRoles() (roles []model.RoleID) {
	for rid := range a.roles {
		if rid != 0 && a.roles[model.RoleID(rid)].ID == model.RoleID(rid) {
			roles = append(roles, model.RoleID(rid))
		}
	}
	return roles
}

// CanA returns whether the API caller has the specified privilege(s) on any
// group.
func (a *Authorizer) CanA(actions model.Privilege) bool {
	return a.CanPA(a.me, actions)
}

// CanAG returns whether the API caller has the specified privilege(s) on the
// specified target group.
func (a *Authorizer) CanAG(actions model.Privilege, group model.GroupID) bool {
	return a.CanPAG(a.me, actions, group)
}

// CanAP returns whether the API caller has the specified privilege(s) on any
// group to which the specified target person belongs.
func (a *Authorizer) CanAP(actions model.Privilege, person model.PersonID) (found bool) {
	a.forPersonRoles(person, func(role model.RoleID) bool {
		if a.CanAR(actions, role) {
			found = true
			return false
		}
		return true
	})
	return found
}

// CanAR returns whether the API caller has the specified privilege(s) on any
// group to which the specified role grants membership.
func (a *Authorizer) CanAR(actions model.Privilege, role model.RoleID) (found bool) {
	for gid := range a.groups {
		if a.CanRAG(role, model.PrivMember, model.GroupID(gid)) && a.CanPAG(a.me, actions, model.GroupID(gid)) {
			return true
		}
	}
	return false
}

// CanPA returns whether the specified person has the specified privilege(s) on
// any group.
func (a *Authorizer) CanPA(person model.PersonID, actions model.Privilege) (found bool) {
	a.forPersonRoles(person, func(rid model.RoleID) bool {
		for gid := range a.groups {
			if a.CanRAG(rid, actions, model.GroupID(gid)) {
				found = true
				return false
			}
		}
		return true
	})
	return found
}

// CanPAG returns whether the specified person has the specified privilege(s) on
// the specified target group.
func (a *Authorizer) CanPAG(person model.PersonID, actions model.Privilege, group model.GroupID) (found bool) {
	a.forPersonRoles(person, func(rid model.RoleID) bool {
		if a.CanRAG(rid, actions, group) {
			found = true
			return false
		}
		return true
	})
	return found
}

// CanRAG returns whether the specified role has the specified privilege(s) on
// the specified target group.
func (a *Authorizer) CanRAG(role model.RoleID, actions model.Privilege, group model.GroupID) bool {
	return a.ActionsRG(role, group)&actions == actions
}

// FetchGroup returns the group with the specified ID, or nil if there is none.
func (a *Authorizer) FetchGroup(id model.GroupID) *model.Group {
	if id > 0 && int(id) < len(a.groups) && a.groups[id].ID == id {
		return &a.groups[id]
	}
	return nil
}

// FetchGroupByEmail returns the group with the specified email, or nil if there
// is none.
func (a *Authorizer) FetchGroupByEmail(email string) *model.Group {
	for gid := range a.groups {
		if a.groups[gid].Email == email {
			return &a.groups[gid]
		}
	}
	return nil
}

// FetchGroupByTag returns the group with the specified tag, or nil if there is
// none.
func (a *Authorizer) FetchGroupByTag(tag model.GroupTag) *model.Group {
	for gid := range a.groups {
		if a.groups[gid].Tag == tag {
			return &a.groups[gid]
		}
	}
	return nil
}

// FetchGroups takes a list of group IDs (usually from one of the GroupsXXX
// functions) and converts it to a sorted list of group objects.
func (a *Authorizer) FetchGroups(gids []model.GroupID) (groups []*model.Group) {
	groups = make([]*model.Group, len(gids))
	for i, gid := range gids {
		groups[i] = &a.groups[gid]
	}
	sort.Sort(model.GroupSort(groups))
	return groups
}

// FetchPeople takes a list of person IDs (usually from one of the PeopleXXX
// functions) and converts it to a sorted list of person objects.
func (a *Authorizer) FetchPeople(pids []model.PersonID) (people []*model.Person) {
	people = make([]*model.Person, len(pids))
	for i, pid := range pids {
		people[i] = a.tx.FetchPerson(pid)
	}
	sort.Sort(model.PersonSort(people))
	return people
}

// FetchRole returns the role with the specified ID, or nil if there is none.
func (a *Authorizer) FetchRole(id model.RoleID) *model.Role {
	if id > 0 && int(id) < len(a.roles) && a.roles[id].ID == id {
		return &a.roles[id]
	}
	return nil
}

// FetchRoleByTag returns the role with the specified tag, or nil if there is
// none.
func (a *Authorizer) FetchRoleByTag(tag model.RoleTag) *model.Role {
	for rid := range a.roles {
		if a.roles[rid].Tag == tag {
			return &a.roles[rid]
		}
	}
	return nil
}

// FetchRoles takes a list of role IDs (usually from one of the RolesXXX
// functions) and converts it to a sorted list of role objects.
func (a *Authorizer) FetchRoles(rids []model.RoleID) (roles []*model.Role) {
	roles = make([]*model.Role, len(rids))
	for i, rid := range rids {
		roles[i] = &a.roles[rid]
	}
	sort.Sort(model.RoleSort(roles))
	return roles
}

// forPersonRoles calls the supplied function for each role held by the
// specified person.
func (a *Authorizer) forPersonRoles(pid model.PersonID, f func(model.RoleID) bool) {
	for by := 0; by < a.bytesPerPerson; by++ {
		b := a.personRoles[int(pid)*a.bytesPerPerson+by]
		for l := bits.Len8(b); l > 0; l = bits.Len8(b) {
			if !f(model.RoleID(by*8 + l - 1)) {
				return
			}
			b &^= 1 << (l - 1)
		}
	}
}

// GroupsA returns the list of groups on which the API caller has the specified
// privilege(s).
func (a *Authorizer) GroupsA(actions model.Privilege) (groups []model.GroupID) {
	return a.groupsPA(a.me, actions)
}

// GroupsP returns the list of groups to which the specified person belongs.
func (a *Authorizer) GroupsP(person model.PersonID) (groups []model.GroupID) {
	return a.groupsPA(person, model.PrivMember)
}

// groupsPA returns the list of groups on which the specified person has the
// specified privilege(s).
func (a *Authorizer) groupsPA(person model.PersonID, actions model.Privilege) (groups []model.GroupID) {
	var privs = make([]model.Privilege, len(a.groups))
	a.forPersonRoles(person, func(role model.RoleID) bool {
		for gid := range a.groups {
			privs[gid] |= a.ActionsRG(role, model.GroupID(gid))
		}
		return true
	})
	for gid := range a.groups {
		if privs[gid]&actions == actions {
			groups = append(groups, model.GroupID(gid))
		}
	}
	return groups
}

// GroupsR returns the list of groups to which the specified role conveys
// membership.
func (a *Authorizer) GroupsR(role model.RoleID) (groups []model.GroupID) {
	return a.groupsRA(role, model.PrivMember)
}

// groupsRA returns the list of groups on which the specified role has the
// specified privilege(s).
func (a *Authorizer) groupsRA(role model.RoleID, actions model.Privilege) (groups []model.GroupID) {
	for group := range a.groups {
		if a.CanRAG(role, actions, model.GroupID(group)) {
			groups = append(groups, model.GroupID(group))
		}
	}
	return groups
}

// holdsPR returns whether the specified person holds the specified role.
func (a *Authorizer) holdsPR(person model.PersonID, role model.RoleID) bool {
	return a.personRoles[int(person)*a.bytesPerPerson+int(role/8)]&(1<<int(role%8)) != 0
}

// holdsR returns whether the API caller holds the specified role.
func (a *Authorizer) holdsR(role model.RoleID) bool {
	return a.holdsPR(a.me, role)
}

// IsWebmaster returns whether the API caller is a webmaster.
func (a *Authorizer) IsWebmaster() bool {
	for rid := range a.roles {
		if a.roles[rid].Tag == model.RoleWebmaster {
			return a.holdsR(model.RoleID(rid))
		}
	}
	return false
}

// May returns whether the API caller has the specified permission.
func (a *Authorizer) May(perm model.Permission) (hasPerm bool) {
	a.forPersonRoles(a.me, func(rid model.RoleID) bool {
		if a.FetchRole(rid).Permissions&perm != 0 {
			hasPerm = true
			return false
		}
		return true
	})
	return hasPerm
}

// MemberG returns whether the API caller is a member of the specified group.
func (a *Authorizer) MemberG(group model.GroupID) bool {
	return a.MemberPG(a.me, group)
}

// MemberPG returns whether the specified person is a member of the specified
// group.
func (a *Authorizer) MemberPG(person model.PersonID, group model.GroupID) (found bool) {
	a.forPersonRoles(person, func(rid model.RoleID) bool {
		if a.memberRG(rid, group) {
			found = true
			return false
		}
		return true
	})
	return found
}

// memberRG returns whether the specified role conveys membership in the
// specified group.
func (a *Authorizer) memberRG(role model.RoleID, group model.GroupID) bool {
	return a.CanRAG(role, model.PrivMember, group)
}

// PeopleA returns the list of people on whom the API caller has the specified
// privilege(s).
func (a *Authorizer) PeopleA(actions model.Privilege) (people []model.PersonID) {
	var roleMask = make([]byte, a.bytesPerPerson)
	a.forPersonRoles(a.me, func(role model.RoleID) bool {
		for gid := range a.groups {
			if a.CanRAG(role, actions, model.GroupID(gid)) {
				for role := range a.roles {
					if a.CanRAG(model.RoleID(role), model.PrivMember, model.GroupID(gid)) {
						roleMask[int(role/8)] |= 1 << int(role%8)
					}
				}
			}
		}
		return true
	})
	for person := 0; person < int(a.numPeople); person++ {
		for by := 0; by < a.bytesPerPerson; by++ {
			if a.personRoles[person*a.bytesPerPerson+by]&roleMask[by] != 0 {
				people = append(people, model.PersonID(person))
				break
			}
		}
	}
	return people
}

// PeopleG returns the list of people who are in the specified group.
func (a *Authorizer) PeopleG(group model.GroupID) (people []model.PersonID) {
	var roleMask = make([]byte, a.bytesPerPerson)
	for role := range a.roles {
		if a.CanRAG(model.RoleID(role), model.PrivMember, group) {
			roleMask[int(role/8)] |= 1 << int(role%8)
		}
	}
	for person := 0; person < int(a.numPeople); person++ {
		for by := 0; by < a.bytesPerPerson; by++ {
			if a.personRoles[person*a.bytesPerPerson+by]&roleMask[by] != 0 {
				people = append(people, model.PersonID(person))
				break
			}
		}
	}
	return people
}

// PeopleR returns the list of people who have the specified role.
func (a *Authorizer) PeopleR(role model.RoleID) (people []model.PersonID) {
	for person := 0; person < int(a.numPeople); person++ {
		if a.holdsPR(model.PersonID(person), role) {
			people = append(people, model.PersonID(person))
		}
	}
	return people
}

// RolesIndividuallyHeld returns the roles that are individual, and currently
// held by a person.  They are returned in a map from role ID to person ID of
// the person holding it.
func (a *Authorizer) RolesIndividuallyHeld() (roles map[model.RoleID]model.PersonID) {
	roles = make(map[model.RoleID]model.PersonID)
	for rid := range a.roles {
		if !a.roles[rid].Individual {
			continue
		}
		for person := 0; person < int(a.numPeople); person++ {
			if a.holdsPR(model.PersonID(person), model.RoleID(rid)) {
				roles[model.RoleID(rid)] = model.PersonID(person)
			}
		}
	}
	return roles
}

// RolesAG returns the list of roles that grant the specified privilege(s) on
// the specified group.
func (a *Authorizer) RolesAG(actions model.Privilege, group model.GroupID) (roles []model.RoleID) {
	for role := model.RoleID(1); role < model.RoleID(len(a.roles)); role++ {
		if a.CanRAG(role, actions, group) {
			roles = append(roles, role)
		}
	}
	return roles
}

// RolesG returns the list of roles that grant membership to the specified
// group.
func (a *Authorizer) RolesG(group model.GroupID) (roles []model.RoleID) {
	for role := model.RoleID(1); role < model.RoleID(len(a.roles)); role++ {
		if a.memberRG(role, group) {
			roles = append(roles, role)
		}
	}
	return roles
}

// RolesP returns the list of roles held by a person.
func (a *Authorizer) RolesP(person model.PersonID) (roles []model.RoleID) {
	roles = []model.RoleID{}
	a.forPersonRoles(person, func(role model.RoleID) bool {
		roles = append(roles, role)
		return true
	})
	return roles
}
