package model

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"unsafe"

	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// chunkSize is the granularity of size increases to the map.
const chunkSize = 32

// A PrivilegeMap is functionally a map[*Role]Privilege, but stored and handled
// more efficiently than that.
type PrivilegeMap struct {
	p []Privilege
}

// Add adds a privilege to a map.
func (pm *PrivilegeMap) Add(group *Group, priv Privilege) {
	pm.enlargeFor(group.ID)
	pm.p[group.ID] |= priv
}

// Remove removes a privilege from a map.
func (pm *PrivilegeMap) Remove(group *Group, priv Privilege) {
	if int(group.ID) < len(pm.p) {
		pm.p[group.ID] &^= priv
	}
}

// Set sets the privilege bitmask on a group in a map.
func (pm *PrivilegeMap) Set(group *Group, priv Privilege) {
	pm.enlargeFor(group.ID)
	pm.p[group.ID] = priv
}

// Clear clears all privileges out of the map.
func (pm *PrivilegeMap) Clear() {
	for i := range pm.p {
		pm.p[i] = 0
	}
}

// Merge merges all of the privileges in the argument map into the receiver map.
func (pm *PrivilegeMap) Merge(other *PrivilegeMap) {
	pm.enlargeFor(GroupID(len(other.p) - 1))
	for r, p := range other.p {
		pm.p[r] |= p
	}
}

// Get returns the privilege bitmask on a group in a map.
func (pm *PrivilegeMap) Get(group *Group) Privilege {
	if int(group.ID) >= len(pm.p) {
		return 0
	}
	return pm.p[group.ID]
}

// Has returns whether the specified privilege(s) exist in the receiver map for
// the specified group.
func (pm *PrivilegeMap) Has(group *Group, priv Privilege) bool {
	if int(group.ID) >= len(pm.p) {
		return false
	}
	return pm.p[group.ID]&priv == priv
}

// HasAny returns whether the receiver map has the specified privilege(s) on any
// group.
func (pm *PrivilegeMap) HasAny(priv Privilege) bool {
	for _, p := range pm.p {
		if p&priv == priv {
			return true
		}
	}
	return false
}

// RolesWith returns an unsorted list of the IDs of all groups for which the map
// contains the specified privilege(s).
func (pm *PrivilegeMap) RolesWith(privs Privilege) (groupIDs []GroupID) {
	for id, p := range pm.p {
		if p&privs == privs {
			groupIDs = append(groupIDs, GroupID(id))
		}
	}
	return groupIDs
}

func (pm *PrivilegeMap) String() string {
	var sb strings.Builder
	sb.WriteByte('{')
	first := true
	for id, p := range pm.p {
		if p == 0 {
			continue
		}
		if first {
			first = false
		} else {
			sb.WriteByte(',')
		}
		fmt.Fprintf(&sb, "%d: %02X", id, p)
	}
	sb.WriteByte('}')
	return sb.String()
}

// enlargeFor enlarges the map, if needed, so that it can accommodate the
// specified group ID.  It returns the new map.
func (pm *PrivilegeMap) enlargeFor(group GroupID) {
	if int(group) < len(pm.p) {
		return
	}
	npm := make([]Privilege, (int(group)+chunkSize)/chunkSize*chunkSize)
	copy(npm, pm.p)
	pm.p = npm
}

// Value translates the map into a blob for database storage.
func (pm *PrivilegeMap) Value() (driver.Value, error) {
	var buf = make([]byte, len(pm.p))
	var bytes []byte
	var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&pm.p))
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
	pm.p = make([]Privilege, len(buf))
	var bytes []byte
	var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&pm.p))
	bhdr.Data = phdr.Data
	bhdr.Len = phdr.Len
	bhdr.Cap = phdr.Cap
	copy(bytes, buf)
	return nil
}

// Size is used by protocol buffers.
func (pm *PrivilegeMap) Size() int {
	if pm == nil || len(pm.p) == 0 {
		return 0
	}
	return 1 + len(pm.p) + sovModel(uint64(len(pm.p)))
}

// MarshalToSizedBuffer is used by protocol buffers.
func (pm *PrivilegeMap) MarshalToSizedBuffer(buf []byte) (int, error) {
	i := len(buf)
	if pm != nil && len(pm.p) > 0 {
		i -= len(pm.p)
		var bytes []byte
		var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
		var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&pm.p))
		bhdr.Data = phdr.Data
		bhdr.Len = phdr.Len
		bhdr.Cap = phdr.Cap
		copy(buf[i:], bytes)
		i = encodeVarintModel(buf, i, uint64(len(pm.p)))
		i--
		buf[i] = 0x0a
	}
	return len(buf) - i, nil
}

// Unmarshal is used by protocol buffers.
func (pm *PrivilegeMap) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	if l == 0 {
		return nil
	}
	if dAtA[0] != 0x0a {
		return errors.New("wrong field number or type for PrivilegeMap")
	}
	iNdEx := 1
	msglen := 0
	for shift := uint(0); ; shift += 7 {
		if shift >= 64 {
			return ErrIntOverflowModel
		}
		if iNdEx >= l {
			return io.ErrUnexpectedEOF
		}
		b := dAtA[iNdEx]
		iNdEx++
		msglen |= int(b&0x7F) << shift
		if b < 0x80 {
			break
		}
	}
	if msglen < 0 {
		return ErrInvalidLengthModel
	}
	postIndex := iNdEx + msglen
	if postIndex < 0 {
		return ErrInvalidLengthModel
	}
	if postIndex > l {
		return io.ErrUnexpectedEOF
	}
	pm.p = make([]Privilege, msglen)
	var bytes []byte
	var bhdr = (*reflect.SliceHeader)(unsafe.Pointer(&bytes))
	var phdr = (*reflect.SliceHeader)(unsafe.Pointer(&pm.p))
	bhdr.Data = phdr.Data
	bhdr.Len = phdr.Len
	bhdr.Cap = phdr.Cap
	copy(bytes, dAtA[iNdEx:postIndex])
	return nil
}

// MarshalEasyJSON encodes the privilege map into JSON.
func (pm PrivilegeMap) MarshalEasyJSON(w *jwriter.Writer) {
	w.RawByte('{')
	first := true
	for id, p := range pm.p {
		if p == 0 {
			continue
		}
		if first {
			first = false
		} else {
			w.RawByte(',')
		}
		w.IntStr(id)
		w.RawByte(':')
		w.Uint8(uint8(p))
	}
	w.RawByte('}')
}

// UnmarshalEasyJSON decodes the privilege map from JSON.
func (pm *PrivilegeMap) UnmarshalEasyJSON(l *jlexer.Lexer) {
	isTopLevel := l.IsStart()
	if l.IsNull() {
		if isTopLevel {
			l.Consumed()
		}
		l.Skip()
		return
	}
	l.Delim('{')
	for !l.IsDelim('}') {
		key := l.UnsafeString()
		l.WantColon()
		if l.IsNull() {
			l.Skip()
			l.WantComma()
			continue
		}
		id, err := strconv.Atoi(key)
		if err != nil {
			l.AddError(err)
		}
		privs := Privilege(l.Uint8())
		pm.enlargeFor(GroupID(id))
		pm.p[id] = privs
		l.WantComma()
	}
	l.Delim('}')
	if isTopLevel {
		l.Consumed()
	}
}
