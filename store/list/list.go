// Package list defines the List type, which describes an email or SMS list to
// which people can be subscribed and to which messages can be sent.
package list

import "k8s.io/apimachinery/pkg/util/sets"

// ID uniquely identifies a list.
type ID int

// List describes an email or SMS list to which people can subscribe and to
// which messages can be sent.
type List struct {
	// ID is the unique identifier of the List.
	ID ID
	// Type is the type of the List.
	Type Type
	// Name is the name of the list.  For email lists, it is also the
	// local-part of the email address of the list.
	Name string
	// Moderators is the set of moderator email addresses for the list.
	// It is nil for unmoderated lists (including all SMS lists).
	Moderators sets.Set[string]
}

func (l *List) Clone() (c *List) {
	c = new(List)
	*c = *l
	if l.Moderators != nil {
		c.Moderators = l.Moderators.Clone()
	}
	return c
}
