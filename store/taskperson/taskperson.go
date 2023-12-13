package taskperson

// Flag is a flag (or bitmask of flags) describing a Person's relationship to a
// Task.
type Flag uint

// Values for Flag:
const (
	// Attended is a flag indicating that the Person attended the Task.
	Attended Flag = 1 << iota
	// Credited is a flag indicating that the Person was credited for
	// participation in the Task.
	Credited
)
