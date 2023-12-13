// Package folder defines the Folder type, which describes a folder of
// documents.
package folder

import "sunnyvaleserv.org/portal/store/enum"

// ID uniquely identifies a folder.
type ID int

// RootID is the ID of the root folder.
const RootID ID = 1

// Fields is a bitmask of flags identifying specified fields of the Folder
// structure.
type Fields uint64

// Values for Fields:
const (
	FID Fields = 1 << iota
	FParent
	FName
	FURLName
	FViewer
	FEditor
)

// Folder describes a folder containing documents.
type Folder struct {
	// NOTE: documentation of the fields is on the getter functions in
	// getters.go.

	fields   Fields
	id       ID
	parent   ID
	name     string
	urlName  string
	viewOrg  enum.Org
	viewPriv enum.PrivLevel
	editOrg  enum.Org
	editPriv enum.PrivLevel
}
