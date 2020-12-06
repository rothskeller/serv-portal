package model

import (
	"errors"
)

// A NoteVisibility value specifies the visibility of a note on a person.
type NoteVisibility uint8

// Values for NoteVisibility
const (
	NoteVisiblityNone NoteVisibility = iota
	// NoteVisibleWithPerson specifies that the note can be seen by anyone
	// who can see the person.
	NoteVisibleWithPerson
	// NoteVisibleWithContact specifies that the note can be seen by anyone
	// who can see the person's contact information.
	NoteVisibleWithContact
	// NoteVisibleToLeaders specifies that the note can be seen by org
	// leaders only.
	NoteVisibleToLeaders
	// NoteVisibleToAdmins specifies that the note can be seen only by
	// admin leaders.
	NoteVisibleToAdmins
	// NoteVisibleToWebmaster specifies that the note can be seen only by
	// webmasters
	NoteVisibleToWebmaster
)

// String returns the string form of the NoteVisibility, as used in APIs.  It is
// not the display label.
func (v NoteVisibility) String() string {
	switch v {
	case NoteVisibleWithPerson:
		return "person"
	case NoteVisibleWithContact:
		return "contact"
	case NoteVisibleToLeaders:
		return "leader"
	case NoteVisibleToAdmins:
		return "admin"
	case NoteVisibleToWebmaster:
		return "webmaster"
	default:
		return ""
	}
}

// ParseNoteVisibility translates the string form of a NoteVisibility (as
// returned by String(), above) into a NoteVisibility value.  It returns an
// error if the string is not recognized.
func ParseNoteVisibility(s string) (NoteVisibility, error) {
	switch s {
	case "person":
		return NoteVisibleWithPerson, nil
	case "contact":
		return NoteVisibleWithContact, nil
	case "leader":
		return NoteVisibleToLeaders, nil
	case "admin":
		return NoteVisibleToAdmins, nil
	case "webmaster":
		return NoteVisibleToWebmaster, nil
	default:
		return NoteVisiblityNone, errors.New("invalid visibility")
	}
}

// Valid returns whether a NoteVisibility value is valid.
func (v NoteVisibility) Valid() bool {
	switch v {
	case NoteVisibleWithPerson, NoteVisibleWithContact, NoteVisibleToLeaders, NoteVisibleToAdmins, NoteVisibleToWebmaster:
		return true
	default:
		return false
	}
}

// AllNoteVisibilities is the ordered list of NoteVisibility values, which is
// used for iteration.
var AllNoteVisibilities = []NoteVisibility{
	NoteVisibleToWebmaster,
	NoteVisibleToAdmins,
	NoteVisibleToLeaders,
	NoteVisibleWithContact,
	NoteVisibleWithPerson,
}
