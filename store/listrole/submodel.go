package listrole

// SubscriptionModel describes the subscription model a given role grants to a
// given list.
type SubscriptionModel uint

// Values for SubscriptionModel.
const (
	// NoSubscription indicates that holders of the role are not granted any
	// subscription privileges on the list.
	NoSubscription SubscriptionModel = iota
	// AllowSubscription indicates that holders of the role are allowed to
	// manually subscribe to the list.
	AllowSubscription
	// AutoSubscribe indicates that holders of the role are automatically
	// subscribed to the list (although they can manually unsubscribe from
	// it).
	AutoSubscribe
	// WarnOnUnsubscribe is like AutoSubscribe, but people trying to
	// manually unsubscribe from the list are warned that they may lose the
	// role if they do.
	WarnOnUnsubscribe
)

// AllSubscriptionModels is a list of all list subscription models values.
var AllSubscriptionModels = []SubscriptionModel{NoSubscription, AllowSubscription, AutoSubscribe, WarnOnUnsubscribe}

func (sm SubscriptionModel) String() string {
	switch sm {
	case AllowSubscription:
		return "allow"
	case AutoSubscribe:
		return "auto"
	case WarnOnUnsubscribe:
		return "warn"
	default:
		return ""
	}
}

func (sm SubscriptionModel) LongString() string {
	switch sm {
	case 0:
		return "Not allowed"
	case AllowSubscription:
		return "Manual"
	case AutoSubscribe:
		return "Automatic"
	case WarnOnUnsubscribe:
		return "Automatic, warn on unsubscribe"
	default:
		return ""
	}
}

func (sm SubscriptionModel) Int() int { return int(sm) }

func (sm SubscriptionModel) Valid() bool {
	return sm == 0 || sm == AllowSubscription || sm == AutoSubscribe || sm == WarnOnUnsubscribe
}
