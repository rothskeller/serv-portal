package model

import (
	"errors"
)

// A FolderOrg identifies the visibility of a folder.
type FolderOrg uint8

// All of the values of Org are valid for FolderOrg.  In addition:
const (
	// OrgSERV represents folders that are visible to all SERV volunteers.
	OrgSERV FolderOrg = 14
	// OrgPublic represents folders that are visible to the general public.
	OrgPublic FolderOrg = 15
)

// String returns the string form of the FolderOrg, as used in APIs.  It is not
// the display label.
func (o FolderOrg) String() string {
	switch o {
	case FolderOrg(OrgAdmin2):
		return "admin"
	case FolderOrg(OrgCERTD2):
		return "cert-d"
	case FolderOrg(OrgCERTT2):
		return "cert-t"
	case FolderOrg(OrgListos):
		return "listos"
	case FolderOrg(OrgSARES2):
		return "sares"
	case FolderOrg(OrgSNAP2):
		return "snap"
	case OrgSERV:
		return "serv"
	case OrgPublic:
		return "public"
	default:
		return ""
	}
}

// ParseFolderOrg translates the string form of an Org (as returned by String(),
// above) into a FolderOrg value.  It returns an error if the string is not
// recognized.
func ParseFolderOrg(s string) (FolderOrg, error) {
	switch s {
	case "admin":
		return FolderOrg(OrgAdmin2), nil
	case "cert-d":
		return FolderOrg(OrgCERTD2), nil
	case "cert-t":
		return FolderOrg(OrgCERTT2), nil
	case "listos":
		return FolderOrg(OrgListos), nil
	case "sares":
		return FolderOrg(OrgSARES2), nil
	case "snap":
		return FolderOrg(OrgSNAP2), nil
	case "serv":
		return OrgSERV, nil
	case "public":
		return OrgPublic, nil
	default:
		return FolderOrg(OrgNone2), errors.New("invalid org")
	}
}

// Valid returns whether a FolderOrg value is valid.
func (o FolderOrg) Valid() bool {
	return (o > 0 && o < FolderOrg(numOrgs)) || o == OrgSERV || o == OrgPublic
}

// AllFolderOrgs is the ordered list of Orgs, which is used for iteration.
var AllFolderOrgs = []FolderOrg{OrgPublic, OrgSERV, FolderOrg(OrgAdmin2), FolderOrg(OrgCERTD2), FolderOrg(OrgCERTT2), FolderOrg(OrgListos), FolderOrg(OrgSARES2), FolderOrg(OrgSNAP2)}
