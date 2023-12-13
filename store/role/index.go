package role

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// IndexKey returns the index key for the venue.
func (r *Role) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("R%d", r.ID())
}

// IndexEntry returns a complete index entry for the venue.
func (r *Role) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	return &phys.IndexEntry{
		Key:   r.IndexKey(nil),
		Type:  "Role",
		Label: r.Name(),
		Name:  r.Name(),
	}
}

// IndexAll indexes all venues.
func IndexAll(storer phys.Storer) {
	All(storer, FID|FName, func(r *Role) {
		phys.Index(storer, r)
	})
}
