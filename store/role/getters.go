package role

import "sunnyvaleserv.org/portal/store/enum"

// Fields returns the set of fields that have been retrieved for this role.
func (r *Role) Fields() Fields {
	return r.fields
}

// ID is the unique identifier of the Role.
func (r *Role) ID() ID {
	if r == nil {
		return 0
	}
	if r.fields&FID == 0 {
		panic("Role.ID called without having fetched FID")
	}
	return r.id
}

// Name is the name of the Role, in title case.  If there can be more than one
// person with this Role, Name should be a plural noun phrase describing the
// group.  If only one person can hold the Role, it will be a singular noun
// phrase.
func (r *Role) Name() string {
	if r.fields&FName == 0 {
		panic("Role.Name called without having fetched FName")
	}
	return r.name
}

// Title is the title shown for a person holding the Role, as a singular noun
// phrase in title case.  If no title should be shown for this Role, Title is an
// empty string.
func (r *Role) Title() string {
	if r.fields&FTitle == 0 {
		panic("Role.Title called without having fetched FTitle")
	}
	return r.title
}

// Priority gives the placement of this Role in the list of roles:  lower values
// are higher in the list.
func (r *Role) Priority() uint {
	if r.fields&FPriority == 0 {
		panic("Role.Priority called without having fetched FPriority")
	}
	return r.priority
}

// Org is the organization to which this Role conveys the privilege identified
// by PrivLevel.
func (r *Role) Org() enum.Org {
	if r.fields&FOrg == 0 {
		panic("Role.Org called without having fetched FOrg")
	}
	return r.org
}

// PrivLevel is the level of privilege conveyed by this Role to the organization
// identified by Org.
func (r *Role) PrivLevel() enum.PrivLevel {
	if r.fields&FPrivLevel == 0 {
		panic("Role.PrivLevel called without having fetched FPrivLevel")
	}
	return r.privLevel
}

// Flags returns the flags associated with this Role.
func (r *Role) Flags() Flags {
	if r.fields&FFlags == 0 {
		panic("Role.Flags called without having fetched FFlags")
	}
	return r.flags
}

// Implies is the list of IDs of other Roles implied by this one, i.e., if this
// Role is assigned to a Person, the listed Roles are also assigned to that
// Person.
func (r *Role) Implies() []ID {
	if r.fields&FImplies == 0 {
		panic("Role.Implies called without having fetched FImplies")
	}
	return r.implies
}
