package authz

import (
	"sunnyvaleserv.org/portal/model"
)

// Save saves the authorizer data to the database.
func (a *Authorizer) Save() {
	var (
		data []byte
		err  error
	)
	if data, err = a.Marshal(); err != nil {
		panic(err)
	}
	a.tx.SaveAuthorizer(data)
}

// SetPersonRoles sets the list of roles held by a person.
func (a *Authorizer) SetPersonRoles(person model.PersonID, roles []model.RoleID) {
	for by := 0; by < a.bytesPerPerson; by++ {
		a.personRoles[int(person)*a.bytesPerPerson+by] = 0
	}
	for _, role := range roles {
		if int(role) >= len(a.roles) || role < 1 {
			panic("role out of range")
		}
		a.personRoles[int(person)*a.bytesPerPerson+int(role/8)] |= 1 << int(role%8)
	}
}

// AddPerson adjusts the authorizer to account for a newly-added person.
func (a *Authorizer) AddPerson(person model.PersonID) {
	if int(person) < int(a.numPeople) {
		for by := 0; by < a.bytesPerPerson; by++ {
			a.personRoles[int(person)*a.bytesPerPerson+by] = 0
		}
		return
	}
	npr := make([]byte, a.bytesPerPerson*(int(person)+1))
	copy(npr, a.personRoles)
	a.personRoles = npr
	a.numPeople = person + 1
}

// CreateGroup creates a new group.
func (a *Authorizer) CreateGroup() *model.Group {
	var id int
	for id = 1; id < len(a.groups); id++ {
		if a.groups[id].ID == 0 {
			a.groups[id] = model.Group{ID: model.GroupID(id)}
			return &a.groups[id]
		}
	}
	newlen := len(a.groups) + 1
	nrp := make([]model.Privilege, len(a.roles)*newlen)
	for rid := range a.roles {
		for gid := range a.groups {
			nrp[rid*(newlen)+gid] = a.rolePrivs[rid*len(a.groups)+gid]
		}
	}
	a.groups = append(a.groups, model.Group{ID: model.GroupID(id)})
	a.rolePrivs = nrp
	return &a.groups[id]
}

// CreateRole creates a new role.
func (a *Authorizer) CreateRole() *model.Role {
	var id int
	for id = 1; int(id) < len(a.roles); id++ {
		if a.roles[id].ID == 0 {
			a.roles[id] = model.Role{ID: model.RoleID(id)}
			return &a.roles[id]
		}
	}
	nrp := make([]model.Privilege, (len(a.roles)+1)*len(a.groups))
	copy(nrp, a.rolePrivs)
	a.rolePrivs = nrp
	a.roles = append(a.roles, model.Role{ID: model.RoleID(id)})
	if len(a.roles)%8 != 1 {
		return &a.roles[id]
	}
	nbpp := a.bytesPerPerson + 1
	npr := make([]byte, int(a.numPeople)*nbpp)
	for pid := 0; pid < int(a.numPeople); pid++ {
		copy(npr[pid*nbpp:], a.personRoles[pid*a.bytesPerPerson:(pid+1)*a.bytesPerPerson])
	}
	a.bytesPerPerson = nbpp
	a.personRoles = npr
	return &a.roles[id]
}

// DeleteGroup deletes a group.
func (a *Authorizer) DeleteGroup(group model.GroupID) {
	if group < 1 || int(group) >= len(a.groups) {
		panic("group out of range")
	}
	for rid := range a.roles {
		a.rolePrivs[rid*len(a.groups)+int(group)] = 0
	}
	a.groups[group] = model.Group{}
}

// DeleteRole deletes a role.
func (a *Authorizer) DeleteRole(role model.RoleID) {
	if role < 1 || int(role) >= len(a.roles) {
		panic("role out of range")
	}
	for gid := range a.groups {
		a.rolePrivs[int(role)*len(a.groups)+gid] = 0
	}
	for pid := 0; pid < int(a.numPeople); pid++ {
		a.personRoles[pid*a.bytesPerPerson+int(role/8)] &^= 1 << int(role%8)
	}
	a.roles[role] = model.Role{}
}

// SetPrivileges sets a role's privileges on a group.
func (a *Authorizer) SetPrivileges(role model.RoleID, actions model.Privilege, group model.GroupID) {
	a.rolePrivs[int(role)*len(a.groups)+int(group)] = actions
}
