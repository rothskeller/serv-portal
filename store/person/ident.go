package person

// An IdentType is a type of identification (or a bitmask of multiple types)
// that has been issued to a volunteer (and potentially should be retrieved if
// they leave).
type IdentType uint8

// Values for IdentType.
const (
	// IDPhoto is a regular DPS Volunteer photo ID badge.
	IDPhoto IdentType = 1 << iota
	// IDCardKey is a photo badge with card key access to DPS buildings.
	IDCardKey
	// IDSERVShirt is a tan long-sleeved button-down shirt identifying the
	// person as a SERV leader and/or class instructor.
	IDSERVShirt
	// IDCERTShirt is a green long-sleeved tee shirt identifying the person
	// as a member of a CERT deployment team.
	IDCERTShirtLS
	// IDCERTShirt is a green short-sleeved tee shirt identifying the person
	// as a member of a CERT deployment team.
	IDCERTShirtSS
)

// String returns the name of the specified IdentType.  String accepts only a
// single IdentType, not a bitmask of multiple IdentTypes.
func (id IdentType) String() string {
	switch id {
	case IDPhoto:
		return "photo ID"
	case IDCardKey:
		return "access card"
	case IDSERVShirt:
		return "tan SERV shirt"
	case IDCERTShirtLS:
		return "green CERT shirt (LS)"
	case IDCERTShirtSS:
		return "green CERT shirt (SS)"
	default:
		return ""
	}
}

// AllIdentTypes is the list of all identification types.
var AllIdentTypes = []IdentType{IDPhoto, IDCardKey, IDSERVShirt, IDCERTShirtLS, IDCERTShirtSS}
