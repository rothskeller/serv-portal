package enum

import (
	"errors"
)

// An Org identifies one of the SERV volunteer organizations.
type Org uint

// Values for Org.
const (
	_ Org = iota
	// OrgAdmin is the pseudo-organization for SERV administration.
	OrgAdmin
	// OrgCERTD is the CERT Deployment Team.
	OrgCERTD
	// OrgCERTT is the CERT Training Committee.
	OrgCERTT
	// OrgListos is the Listos Team.
	OrgListos
	// OrgSARES is the SARES organization.
	OrgSARES
	// OrgSNAP is the SNAP Team.
	OrgSNAP
	// NumOrgs is the number of defined organizations.
	NumOrgs
)

// String returns the string form of the Org, as used in APIs.  It is not the
// display label.
func (o Org) String() string {
	switch o {
	case OrgAdmin:
		return "Admin"
	case OrgCERTD:
		return "CERT-D"
	case OrgCERTT:
		return "CERT-T"
	case OrgListos:
		return "Listos"
	case OrgSARES:
		return "SARES"
	case OrgSNAP:
		return "SNAP"
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
		return "CERT Deployment"
	case OrgCERTT:
		return "CERT Training"
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

// ParseOrg translates the string form of an Org (as returned by String()) into
// an Org value.  It returns an error if the string is not recognized.
func ParseOrg(s string) (Org, error) {
	switch s {
	case "Admin":
		return OrgAdmin, nil
	case "CERT-D":
		return OrgCERTD, nil
	case "CERT-T":
		return OrgCERTT, nil
	case "Listos":
		return OrgListos, nil
	case "SARES":
		return OrgSARES, nil
	case "SNAP":
		return OrgSNAP, nil
	default:
		return 0, errors.New("invalid org")
	}
}

func (o Org) Int() int { return int(o) }

// Valid returns whether an Org value is valid.
func (o Org) Valid() bool {
	return o > 0 && o < NumOrgs
}

// AllOrgs is the ordered list of Orgs, which is used for iteration.
var AllOrgs = []Org{OrgAdmin, OrgCERTD, OrgCERTT, OrgListos, OrgSARES, OrgSNAP}

// MembersCanViewContactInfo returns whether Members of the receiver Org can
// view contact info for members of the Org.
func (o Org) MembersCanViewContactInfo() bool { return o != OrgSARES }
