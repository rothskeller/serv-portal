package event

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/venue"
)

// IndexKey returns the index key for the event.
func (e *Event) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("E%d", e.ID())
}

// IndexEntry returns a complete index entry for the event.
func (e *Event) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	var label = e.Start()[:10] + " " + e.Name()
	return &phys.IndexEntry{
		Key:     e.IndexKey(nil),
		Type:    "Event",
		Label:   label,
		Name:    e.Name(),
		Date:    e.Start()[:10],
		Context: e.Details(),
	}
}

// IndexAll indexes all people.
func IndexAll(storer phys.Storer) {
	AllBetween(storer, "0000-00-00", "2999-99-99", FID|FStart|FName|FDetails, 0, func(e *Event, _ *venue.Venue) {
		phys.Index(storer, e)
	})
}
