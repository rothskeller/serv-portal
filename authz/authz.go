package authz

import (
	"errors"
	"fmt"
	"io"
	"math/bits"
	"reflect"
	"unsafe"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
)

// NewAuthorizer returns a new Authorizer backed by the specified database
// transaction, which must be open throughout the lifetime of the Authorizer.
// This call fetches all of the authorizer data from the database.
func NewAuthorizer(tx *db.Tx) (a *Authorizer) {
	var data []byte

	a = &Authorizer{tx: tx}
	data = tx.FetchAuthorizer()
	if err := a.Unmarshal(data); err != nil {
		panic(fmt.Sprintf("error reading authorizer data: %s", err))
	}
	return a
}

// An Authorizer encapsulates all of the algorithms and data used to determine
// which people (users) are authorized to perform which actions on which target
// groups.
//
// The query methods on this type follow a naming pattern.  They start with a
// name that indicates what type of information is returned:
//     Can = boolean, action allowed or not
//     Groups = list of groups meeting criteria
//     Actions = bitmask of actions meeting criteria
//     Roles = list of roles meeting criteria
//     People = list of people meeting criteria
//     Member = boolean, member of group or not
// That name is then followed by zero or more characters indicating the
// arguments passed to the function:
//     P = a person, the actor in the query
//     R = a role, the actor in the query
//     A = a bitmask of actions
//     G = a group, the target in the query
// In a non-People, non-Roles call without a P or R argument, the actor is
// assumed to be the caller of the API.  In a non-Groups call without a G
// argument, the query is satisfied by any target group.  In a non-Actions call
// without an A argument, the action is assumed to be PrivMember.
type Authorizer struct {
	tx             *db.Tx
	me             model.PersonID
	roles          []model.Role
	groups         []model.Group
	numPeople      model.PersonID
	rolePrivs      []model.Privilege
	personRoles    []byte
	bytesPerPerson int
}

// SetMe sets the identity of the API caller in the authorizer.
func (a *Authorizer) SetMe(me *model.Person) {
	if me != nil {
		a.me = me.ID
	} else {
		a.me = 0
	}
}

/*
func (a *Authorizer) Transition() {
	var (
		roles       []*model.Role
		maxRoleID   model.RoleID
		groups      []*model.Group
		maxGroupID  model.GroupID
		people      []*model.Person
		maxPersonID model.PersonID
		data        []byte
		err         error
	)
	roles = a.tx.FetchRoles()
	groups = a.tx.FetchGroups()
	people = a.tx.FetchPeople()
	for _, role := range roles {
		if role.ID > maxRoleID {
			maxRoleID = role.ID
		}
	}
	for _, group := range groups {
		if group.ID > maxGroupID {
			maxGroupID = group.ID
		}
	}
	for _, person := range people {
		if person.ID > maxPersonID {
			maxPersonID = person.ID
		}
	}
	a.roles = make([]model.Role, maxRoleID+1)
	a.groups = make([]model.Group, maxGroupID+1)
	a.rolePrivs = make([]model.Privilege, len(a.roles)*len(a.groups))
	a.bytesPerPerson = (len(a.roles) + 7) / 8
	a.personRoles = make([]byte, int(maxPersonID+1)*a.bytesPerPerson)
	for _, role := range roles {
		a.roles[role.ID] = *role
	}
	for _, group := range groups {
		a.groups[group.ID] = *group
	}
	for _, role := range roles {
		for _, group := range groups {
			a.rolePrivs[int(role.ID)*len(a.groups)+int(group.ID)] = role.Privileges.Get(group)
		}
		role.Privileges = model.PrivilegeMap{}
	}
	for _, person := range people {
		for _, rid := range person.Roles {
			a.personRoles[int(person.ID)*a.bytesPerPerson+int(rid/8)] |= 1 << int(rid%8)
		}
	}
	if data, err = a.Marshal(); err != nil {
		panic(err)
	}
	a.tx.CreateAuthorizer(data)
}
*/

// Marshal renders the authorizer data in protocol buffer format.
func (a *Authorizer) Marshal() (buf []byte, err error) {
	size := a.Size()
	buf = make([]byte, size)
	n, err := a.MarshalToSizedBuffer(buf[:size])
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

// MarshalTo renders the data in protocol buffer format.
func (a *Authorizer) MarshalTo(buf []byte) (int, error) {
	size := a.Size()
	return a.MarshalToSizedBuffer(buf[:size])
}

// Size computes the size of the protocol buffer format of the authorizer data.
func (a *Authorizer) Size() (n int) {
	var l int
	{
		l = len(a.personRoles)
		n += 1 + l + sovModel(uint64(l))
	}
	{
		l = len(a.rolePrivs) * int(unsafe.Sizeof(model.PrivMember))
		n += 1 + l + sovModel(uint64(l))
	}
	for _, e := range a.groups {
		l = e.Size()
		n += 1 + l + sovModel(uint64(l))
	}
	for _, e := range a.roles {
		l = e.Size()
		n += 1 + l + sovModel(uint64(l))
	}
	return n
}

// MarshalToSizedBuffer renders the authorizer data in protocol buffer format.
func (a *Authorizer) MarshalToSizedBuffer(buf []byte) (int, error) {
	i := len(buf)
	// Encode the personRoles array.
	{
		i -= len(a.personRoles)
		copy(buf[i:], a.personRoles)
		i = encodeVarintModel(buf, i, uint64(len(a.personRoles)))
		i--
		buf[i] = 0x22
	}
	// Encode the rolePrivs array.
	{
		var bytes []byte
		var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
		var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&a.rolePrivs))
		bhdr.Data = phdr.Data
		bhdr.Len = phdr.Len * int(unsafe.Sizeof(model.PrivMember))
		bhdr.Cap = bhdr.Len
		i -= bhdr.Len
		copy(buf[i:], bytes)
		i = encodeVarintModel(buf, i, uint64(bhdr.Len))
		i--
		buf[i] = 0x1a
	}
	// Encode the roles.
	for idx := len(a.roles) - 1; idx >= 0; idx-- {
		size, err := a.roles[idx].MarshalToSizedBuffer(buf[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintModel(buf, i, uint64(size))
		i--
		buf[i] = 0x12
	}
	// Encode the groups.
	for idx := len(a.groups) - 1; idx >= 0; idx-- {
		size, err := a.groups[idx].MarshalToSizedBuffer(buf[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintModel(buf, i, uint64(size))
		i--
		buf[i] = 0xa
	}
	return len(buf) - i, nil
}

func encodeVarintModel(buf []byte, offset int, v uint64) int {
	offset -= sovModel(v)
	base := offset
	for v >= 1<<7 {
		buf[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	buf[offset] = uint8(v)
	return base
}
func sovModel(v uint64) int { return (bits.Len64(v|1) + 6) / 7 }

// Unmarshal decodes the authorizer data from protocol buffer format.
func (a *Authorizer) Unmarshal(buf []byte) error {
	l := len(buf)
	idx := 0
	for idx < l {
		var wire uint64 = uint64(buf[idx])
		idx++
		if wire > 0x7F {
			return errors.New("fieldnum > 15 not implemented for Authorizer")
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field groups", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return errors.New("groups size too large")
				}
				if idx >= l {
					return io.ErrUnexpectedEOF
				}
				b := buf[idx]
				idx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return errors.New("negative groups size")
			}
			postIndex := idx + msglen
			if postIndex < 0 {
				return errors.New("groups size too large")
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			a.groups = append(a.groups, model.Group{})
			if err := a.groups[len(a.groups)-1].Unmarshal(buf[idx:postIndex]); err != nil {
				return err
			}
			idx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field roles", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return errors.New("roles size too large")
				}
				if idx >= l {
					return io.ErrUnexpectedEOF
				}
				b := buf[idx]
				idx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return errors.New("negative roles size")
			}
			postIndex := idx + msglen
			if postIndex < 0 {
				return errors.New("roles size too large")
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			a.roles = append(a.roles, model.Role{})
			if err := a.roles[len(a.roles)-1].Unmarshal(buf[idx:postIndex]); err != nil {
				return err
			}
			idx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field rolePrivs", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return errors.New("rolePrivs size too large")
				}
				if idx >= l {
					return io.ErrUnexpectedEOF
				}
				b := buf[idx]
				idx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return errors.New("negative rolePrivs size")
			}
			postIndex := idx + msglen
			if postIndex < 0 {
				return errors.New("rolePrivs size too large")
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			a.rolePrivs = make([]model.Privilege, msglen/int(unsafe.Sizeof(model.PrivMember)))
			var bytes []byte
			var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
			var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&a.rolePrivs))
			bhdr.Data = phdr.Data
			bhdr.Len = msglen
			bhdr.Cap = msglen
			copy(bytes, buf[idx:])
			idx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field personRoles", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return errors.New("personRoles size too large")
				}
				if idx >= l {
					return io.ErrUnexpectedEOF
				}
				b := buf[idx]
				idx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return errors.New("negative personRoles size")
			}
			postIndex := idx + msglen
			if postIndex < 0 {
				return errors.New("personRoles size too large")
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			a.personRoles = make([]byte, msglen)
			copy(a.personRoles, buf[idx:])
			idx = postIndex
		default:
			return errors.New("unexpected field number in Authorizer")
		}
	}
	if idx > l {
		return io.ErrUnexpectedEOF
	}
	if a.roles == nil || a.groups == nil || a.rolePrivs == nil || a.personRoles == nil {
		return errors.New("missing required fields for Authorizer")
	}
	a.bytesPerPerson = (len(a.roles) + 7) / 8
	a.numPeople = model.PersonID(len(a.personRoles) / a.bytesPerPerson)
	return nil
}
