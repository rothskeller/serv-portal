package model

import (
	"errors"
)

// An Org identifies one of the SERV volunteer organizations.
type Org uint8

// Values for Org.
const (
	OrgNone2 Org = iota
	OrgAdmin2
	OrgCERTD2
	OrgCERTT2
	OrgListos
	OrgSARES2
	OrgSNAP2
	numOrgs
)

// NumOrgs is the number of defined Org values, used to size slices indexed by
// Org.
const NumOrgs = int(numOrgs)

// String returns the string form of the Org, as used in APIs.  It is not the
// display label.
func (o Org) String() string {
	switch o {
	case OrgAdmin2:
		return "admin"
	case OrgCERTD2:
		return "cert-d"
	case OrgCERTT2:
		return "cert-t"
	case OrgListos:
		return "listos"
	case OrgSARES2:
		return "sares"
	case OrgSNAP2:
		return "snap"
	default:
		return ""
	}
}

// Label returns the display label of the Org.
func (o Org) Label() string {
	switch o {
	case OrgAdmin2:
		return "Admin"
	case OrgCERTD2:
		return "CERT-D"
	case OrgCERTT2:
		return "CERT-T"
	case OrgListos:
		return "Listos"
	case OrgSARES2:
		return "SARES"
	case OrgSNAP2:
		return "SNAP"
	default:
		return ""
	}
}

// ParseOrg translates the string form of an Org (as returned by String(),
// above) into an Org value.  It returns an error if the string is not
// recognized.
func ParseOrg(s string) (Org, error) {
	switch s {
	case "admin":
		return OrgAdmin2, nil
	case "cert-d":
		return OrgCERTD2, nil
	case "cert-t":
		return OrgCERTT2, nil
	case "listos":
		return OrgListos, nil
	case "sares":
		return OrgSARES2, nil
	case "snap":
		return OrgSNAP2, nil
	default:
		return OrgNone2, errors.New("invalid org")
	}
}

// Valid returns whether an Org value is valid.
func (o Org) Valid() bool {
	return o > 0 && o < numOrgs
}

// AllOrgs is the ordered list of Orgs, which is used for iteration.
var AllOrgs = []Org{OrgAdmin2, OrgCERTD2, OrgCERTT2, OrgListos, OrgSARES2, OrgSNAP2}

// MembersCanViewContactInfo returns whether Members of the receiver Org can
// view contact info for members of the Org.
func (o Org) MembersCanViewContactInfo() bool { return o != OrgSARES2 }

// DSWClass returns the DSW class for the receiver Org.
func (o Org) DSWClass() DSWClass {
	switch o {
	case OrgSARES2:
		return DSWComm
	case OrgCERTD2, OrgCERTT2:
		return DSWCERT
	default:
		return 0
	}
}
