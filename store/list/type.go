package list

// Type identifies a type of list.
type Type uint

// Values for Type:
const (
	_ Type = iota
	// Email is a list whose messages are distributed by email.
	Email
	// SMS is a list whose messages are distributed by SMS (i.e., cell phone
	// text messages).
	SMS
)

// AllTypes is the list of all known list types.
var AllTypes = []Type{Email, SMS}

func (t Type) String() string {
	switch t {
	case Email:
		return "Email"
	case SMS:
		return "SMS"
	default:
		return ""
	}
}

func (t Type) Int() int { return int(t) }
