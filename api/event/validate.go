package event

import (
	"errors"
	"regexp"
	"strings"

	"github.com/microcosm-cc/bluemonday"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

var htmlSanitizer = bluemonday.NewPolicy().
	RequireParseableURLs(true).
	AllowURLSchemes("http", "https").
	RequireNoFollowOnLinks(true).
	AllowAttrs("href").OnElements("a").
	AddTargetBlankToFullyQualifiedLinks(true)

var dateRE = regexp.MustCompile(`^20\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
var timeRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)
var yearRE = regexp.MustCompile(`^20\d\d$`)

// ValidateEvent validates the details of an event.
func ValidateEvent(tx *store.Tx, event *model.Event) error {
	var (
		etype      model.EventType
		seenGroups = map[model.GroupID]bool{}
	)
	if event.Name = strings.TrimSpace(event.Name); event.Name == "" {
		return errors.New("missing name")
	}
	if event.Date == "" {
		return errors.New("missing date")
	} else if !dateRE.MatchString(event.Date) {
		return errors.New("invalid date (YYYY-MM-DD)")
	}
	if event.Start == "" {
		return errors.New("missing start")
	} else if !timeRE.MatchString(event.Start) {
		return errors.New("invalid start (HH:MM)")
	}
	if event.End == "" {
		return errors.New("missing end")
	} else if !timeRE.MatchString(event.End) {
		return errors.New("invalid end (HH:MM)")
	}
	if event.End < event.Start {
		return errors.New("end before start")
	}
	if event.Venue != 0 && tx.FetchVenue(event.Venue) == nil {
		return errors.New("nonexistent venue")
	}
	event.Details = htmlSanitizer.Sanitize(strings.TrimSpace(event.Details))
	if event.Organization != model.OrgNone {
		var matched bool
		for _, o := range model.AllOrganizations {
			if o == event.Organization {
				matched = true
				break
			}
		}
		if !matched {
			return errors.New("invalid organization")
		}
	}
	for _, et := range model.AllEventTypes {
		etype &^= et
	}
	if etype != 0 {
		return errors.New("invalid types")
	}
	if len(event.Groups) == 0 {
		return errors.New("missing group")
	}
	for _, g := range event.Groups {
		if seenGroups[g] {
			return errors.New("duplicate group in groups list")
		}
		seenGroups[g] = true
		if tx.Authorizer().FetchGroup(g) == nil {
			return errors.New("invalid group")
		}
	}
	for _, e := range tx.FetchEvents(event.Date, event.Date) {
		if e.ID != event.ID && e.Name == event.Name {
			return errors.New("duplicate name")
		}
	}
	return nil
}
