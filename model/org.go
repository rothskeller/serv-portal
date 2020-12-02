package model

import (
	"errors"
)

// An Org identifies one of the SERV volunteer organizations.
type Org uint8

// Values for Org.
const (
	OrgNone Org = iota
	OrgAdmin
	OrgCERTD
	OrgCERTT
	OrgListos
	OrgSARES
	OrgSNAP
	numOrgs
)

// NumOrgs is the number of defined Org values, used to size slices indexed by
// Org.
const NumOrgs = int(numOrgs)

// String returns the string form of the Org, as used in APIs.  It is not the
// display label.
func (o Org) String() string {
	switch o {
	case OrgAdmin:
		return "admin"
	case OrgCERTD:
		return "cert-d"
	case OrgCERTT:
		return "cert-t"
	case OrgListos:
		return "listos"
	case OrgSARES:
		return "sares"
	case OrgSNAP:
		return "snap"
	default:
		return ""
	}
}

// Label returns the display label of the Org.
func (o Org) Label() string {
	switch o {
	case OrgAdmin:
		return "SERV Leads"
	case OrgCERTD:
		return "CERT Deployment Teams"
	case OrgCERTT:
		return "CERT Training Committee"
	case OrgListos:
		return "Listos Team"
	case OrgSARES:
		return "SARES Members"
	case OrgSNAP:
		return "SNAP Team"
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
		return OrgAdmin, nil
	case "cert-d":
		return OrgCERTD, nil
	case "cert-t":
		return OrgCERTT, nil
	case "listos":
		return OrgListos, nil
	case "sares":
		return OrgSARES, nil
	case "snap":
		return OrgSNAP, nil
	default:
		return OrgNone, errors.New("invalid org")
	}
}

// Valid returns whether an Org value is valid.
func (o Org) Valid() bool {
	return o > 0 && o < numOrgs
}

// AllOrgs is the ordered list of Orgs, which is used for iteration.
var AllOrgs = []Org{OrgAdmin, OrgCERTD, OrgCERTT, OrgListos, OrgSARES, OrgSNAP}

// MembersCanViewContactInfo returns whether Members of the receiver Org can
// view contact info for members of the Org.
func (o Org) MembersCanViewContactInfo() bool { return o != OrgSARES }

// DSWClass returns the DSW class for the receiver Org.
func (o Org) DSWClass() DSWClass {
	switch o {
	case OrgSARES:
		return DSWComm
	case OrgCERTD, OrgCERTT:
		return DSWCERT
	default:
		return 0
	}
}
