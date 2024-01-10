package class

import "sunnyvaleserv.org/portal/store/enum"

// A Type is a type of class (i.e., a curriculum).
type Type uint8

// Values for Type:
const (
	_ Type = iota
	// CERTBasic is a CERT basic training class.
	CERTBasic
	// PEP is a Personal Emergency Preparedness class (or its Spanish
	// equivalent, Preparación para desastres y emergencias).
	PEP
)

// String returns the name of the specified Type.
func (ctype Type) String() string {
	switch ctype {
	case CERTBasic:
		return "CERT Basic Training"
	case PEP:
		return "Personal Emergency Preparedness"
	default:
		return ""
	}
}

// Org returns the organization for the specified class type.
func (ctype Type) Org() enum.Org {
	switch ctype {
	case CERTBasic:
		return enum.OrgCERTT
	case PEP:
		return enum.OrgListos
	default:
		return 0
	}
}

// Int returns the specified Type as an integer.
func (ctype Type) Int() int { return int(ctype) }

// AllTypes is the list of all class types.
var AllTypes = []Type{CERTBasic, PEP}
