package class

// A Referral is a way that registrants learn about the class, as provided by
// them when they register.
type Referral uint8

// Values for Referral:
const (
	_ Referral = iota
	WordOfMouth
	Flyer
	EventBooth
	SocialMedia
)

// String returns the name of the specified Referral.
func (ref Referral) String() string {
	switch ref {
	case WordOfMouth:
		return "Word of mouth"
	case EventBooth:
		return "Information table at an event"
	case Flyer:
		return "Printed advertisement"
	case SocialMedia:
		return "Online advertisement"
	default:
		return ""
	}
}

// AllReferrals is the list of all class types.
var AllReferrals = []Referral{WordOfMouth, EventBooth, Flyer, SocialMedia}

func (ref Referral) Valid() bool {
	return ref == WordOfMouth || ref == EventBooth || ref == Flyer || ref == SocialMedia
}
