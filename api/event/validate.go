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
	var seenRoles = map[model.Role2ID]bool{}

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
	var matched bool
	for _, et := range model.AllEventTypes {
		if et == event.Type {
			matched = true
			break
		}
	}
	if !matched {
		return errors.New("invalid type")
	}
	if !event.Org.Valid() {
		return errors.New("invalid org")
	}
	for _, r := range event.Roles {
		if seenRoles[r] {
			return errors.New("duplicate role in roles list")
		}
		seenRoles[r] = true
		if tx.FetchRole(r) == nil {
			return errors.New("invalid role")
		}
	}
	for _, e := range tx.FetchEvents(event.Date, event.Date) {
		if e.ID != event.ID && e.Name == event.Name {
			return errors.New("duplicate name")
		}
	}
	return nil
}
