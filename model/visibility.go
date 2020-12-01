package model

import (
	"errors"
)

// A FolderVisibility value specifies the visibility of a folder.
type FolderVisibility uint8

// Values for FolderVisibility
const (
	FolderVisiblityNone FolderVisibility = iota
	// FolderVisibleToPublic specifies that the folder is visible to the
	// general public.
	FolderVisibleToPublic
	// FolderVisibleToSERV specifies that the folder is visible to anyone
	// logged into the site.
	FolderVisibleToSERV
	// FolderVisibleToOrg specifies that the folder is only visible to
	// members of one SERV organization, which is specified by the
	// folder.Org value.
	FolderVisibleToOrg
)

// String returns the string form of the FolderVisibility, as used in APIs.  It
// is not the display label.
func (v FolderVisibility) String() string {
	switch v {
	case FolderVisibleToPublic:
		return "public"
	case FolderVisibleToSERV:
		return "serv"
	case FolderVisibleToOrg:
		return "org"
	default:
		return ""
	}
}

// Label returns the display label of the FolderVisibility.
func (v FolderVisibility) Label() string {
	switch v {
	case FolderVisibleToPublic:
		return "Public"
	case FolderVisibleToSERV:
		return "SERV Volunteers"
	case FolderVisibleToOrg:
		return "Specific Organization"
	default:
		return ""
	}
}

// ParseFolderVisibility translates the string form of a FolderVisibility (as
// returned by String(), above) into a FolderVisibility value.  It returns an
// error if the string is not recognized.
func ParseFolderVisibility(s string) (FolderVisibility, error) {
	switch s {
	case "public":
		return FolderVisibleToPublic, nil
	case "serv":
		return FolderVisibleToSERV, nil
	case "org":
		return FolderVisibleToOrg, nil
	default:
		return FolderVisiblityNone, errors.New("invalid visibility")
	}
}

// Valid returns whether a FolderVisibility value is valid.
func (v FolderVisibility) Valid() bool {
	switch v {
	case FolderVisibleToPublic, FolderVisibleToSERV, FolderVisibleToOrg:
		return true
	default:
		return false
	}
}

// AllFolderVisibilities is the ordered list of FolderVisibility values, which
// is used for iteration.
var AllFolderVisibilities = []FolderVisibility{FolderVisibleToPublic, FolderVisibleToSERV, FolderVisibleToOrg}
