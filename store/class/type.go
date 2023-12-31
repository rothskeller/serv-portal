package class

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

// AllTypes is the list of all class types.
var AllTypes = []Type{CERTBasic, PEP}
