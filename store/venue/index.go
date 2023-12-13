package venue

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// IndexKey returns the index key for the venue.
func (v *Venue) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("V%d", v.ID())
}

// IndexEntry returns a complete index entry for the venue.
func (v *Venue) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	return &phys.IndexEntry{
		Key:   v.IndexKey(nil),
		Type:  "Venue",
		Label: v.Name(),
		Name:  v.Name(),
	}
}

// IndexAll indexes all venues.
func IndexAll(storer phys.Storer) {
	All(storer, FID|FName, func(v *Venue) {
		phys.Index(storer, v)
	})
}
