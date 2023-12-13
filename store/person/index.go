package person

import (
	"fmt"

	"sunnyvaleserv.org/portal/store/internal/phys"
)

// IndexKey returns the index key for the person.
func (p *Person) IndexKey(_ phys.Storer) string {
	return fmt.Sprintf("P%d", p.ID())
}

// IndexEntry returns a complete index entry for the person.
func (p *Person) IndexEntry(_ phys.Storer) *phys.IndexEntry {
	var label = p.InformalName()
	if cs := p.CallSign(); cs != "" {
		label += " " + cs
	}
	return &phys.IndexEntry{
		Key:      p.IndexKey(nil),
		Type:     "Person",
		Label:    label,
		Name:     p.InformalName(),
		CallSign: p.CallSign(),
	}
}

// IndexAll indexes all people.
func IndexAll(storer phys.Storer) {
	All(storer, FID|FInformalName, func(p *Person) {
		phys.Index(storer, p)
	})
}
