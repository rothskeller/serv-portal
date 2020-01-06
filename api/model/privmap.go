package model

import (
	"database/sql/driver"
	"encoding/hex"
	"errors"
	"reflect"
	"unsafe"
)

// A PrivilegeMap is a map from role to privilege bitmask.
type PrivilegeMap []Privilege

// Add adds a privilege to a map.
func (pm PrivilegeMap) Add(role RoleID, priv Privilege) PrivilegeMap {
	pm = pm.enlargeFor(role)
	pm[role] |= priv
	return pm
}

// Remove removes a privilege from a map.
func (pm PrivilegeMap) Remove(role RoleID, priv Privilege) PrivilegeMap {
	if int(role) < len(pm) {
		pm[role] &^= priv
	}
	return pm
}

// Set sets the privilege on a role in a map.
func (pm PrivilegeMap) Set(role RoleID, priv Privilege) PrivilegeMap {
	pm = pm.enlargeFor(role)
	pm[role] = priv
	return pm
}

// Get returns the privileges on a role in a map.
func (pm PrivilegeMap) Get(role RoleID) Privilege {
	if int(role) >= len(pm) {
		return 0
	}
	return pm[role]
}

// Has returns whether the specified privilege(s) exist in the receiver map for
// the specified role.
func (pm PrivilegeMap) Has(role RoleID, priv Privilege) bool {
	if int(role) >= len(pm) {
		return false
	}
	return pm[role]&priv == priv
}

// HasAny returns whether the receiver map has the specified privilege(s) on any
// role.
func (pm PrivilegeMap) HasAny(priv Privilege) bool {
	for _, p := range pm {
		if p&priv == priv {
			return true
		}
	}
	return false
}

// Merge merges all of the privileges in the argument map into the receiver map.
func (pm PrivilegeMap) Merge(other PrivilegeMap) PrivilegeMap {
	pm = pm.enlargeFor(RoleID(len(other) - 1))
	for r, p := range other {
		pm[r] |= p
	}
	return pm
}

// RolesWith returns an unsorted list of all roles for which the map contains
// the specified privilege(s).
func (pm PrivilegeMap) RolesWith(privs Privilege) (roles []RoleID) {
	for r, p := range pm {
		if p&privs == privs {
			roles = append(roles, RoleID(r))
		}
	}
	return roles
}

func (pm PrivilegeMap) String() string {
	var bytes []byte
	phdr := (*reflect.SliceHeader)(unsafe.Pointer(&pm))
	bhdr := (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	*bhdr = *phdr
	return hex.EncodeToString(bytes)
}

// enlargeFor enlarges the map, if needed, so that it can accommodate the
// specified role ID.  It returns the new map.
func (pm PrivilegeMap) enlargeFor(role RoleID) PrivilegeMap {
	switch {
	case int(role) < len(pm):
		return pm
	case int(role) == len(pm):
		return append(pm, 0)
	default:
		npm := make(PrivilegeMap, int(role)+1)
		copy(npm, pm)
		return npm
	}
}

// Value translates the map into a blob for database storage.
func (pm PrivilegeMap) Value() (driver.Value, error) {
	var buf = make([]byte, len(pm))
	var bytes []byte
	var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&pm))
	bhdr.Data = phdr.Data
	bhdr.Len = phdr.Len
	bhdr.Cap = phdr.Cap
	copy(buf, bytes)
	return buf, nil
}

// Scan translates a database blob into a map.
func (pm *PrivilegeMap) Scan(value interface{}) error {
	buf, ok := value.([]byte)
	if !ok {
		return errors.New("PrivilegeMap.Scan expects []byte")
	}
	*pm = make(PrivilegeMap, len(buf))
	var bytes []byte
	var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	var phdr = (*reflect.SliceHeader)(unsafe.Pointer(pm))
	bhdr.Data = phdr.Data
	bhdr.Len = phdr.Len
	bhdr.Cap = phdr.Cap
	copy(bytes, buf)
	return nil
}

// RecalculateTransitivePrivilegeMaps recalculates the transitive privilege maps
// in the specified list of roles (which must be all defined roles, plus
// possibly a new role being defined).  It returns true if successful, false if
// there was a cycle in the role graph.
func RecalculateTransitivePrivilegeMaps(roles []*Role) bool {
	var (
		maxID    RoleID
		visiting = map[*Role]bool{}
	)
	for _, role := range roles {
		if role.ID > maxID {
			maxID = role.ID
		}
		role.TransPrivs = nil
	}
	for _, role := range roles {
		if role.TransPrivs == nil {
			if !calcTransPrivs(roles, role, maxID, visiting) {
				return false
			}
		}
	}
	for _, target := range roles {
		for _, implied := range roles {
			if target != implied && target.TransPrivs.Has(implied.ID, PrivHoldsRole) {
				for _, actor := range roles {
					privs := actor.TransPrivs.Get(implied.ID) &^ PrivHoldsRole
					actor.TransPrivs = actor.TransPrivs.Add(target.ID, privs)
				}
			}
		}
	}
	return true
}
func calcTransPrivs(roles []*Role, role *Role, maxID RoleID, visiting map[*Role]bool) bool {
	if visiting[role] {
		return false // cycle in role graph
	}
	visiting[role] = true
	role.TransPrivs = make(PrivilegeMap, maxID+1).Merge(role.PrivMap)
	for _, other := range roles {
		if other == role {
			continue
		}
		if !role.PrivMap.Has(other.ID, PrivHoldsRole) {
			continue
		}
		if other.TransPrivs == nil {
			if !calcTransPrivs(roles, other, maxID, visiting) {
				return false
			}
		}
		role.TransPrivs = role.TransPrivs.Merge(other.TransPrivs)
	}
	visiting[role] = false
	return true
}
