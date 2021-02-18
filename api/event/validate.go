package event

import (
	"errors"
	"regexp"
	"sort"
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
	var seenRoles = map[model.RoleID]bool{}
	var seenShifts = map[string]map[string]bool{}
	var seenSignups = map[model.PersonID][]*model.Shift{}

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
	for _, s := range event.Shifts {
		if m := seenShifts[s.Start]; m != nil {
			if m[s.Task] {
				return errors.New("duplicate shift")
			}
			m[s.Task] = true
		} else {
			seenShifts[s.Start] = make(map[string]bool)
			seenShifts[s.Start][s.Task] = true
		}
		if !timeRE.MatchString(s.Start) {
			return errors.New("invalid shift start time")
		}
		if !timeRE.MatchString(s.End) {
			return errors.New("invalid shift end time")
		}
		if s.End <= s.Start {
			return errors.New("shift end time is not after shift start time")
		}
		if s.Min < 0 {
			return errors.New("invalid shift min")
		}
		if s.Max < 0 || (s.Max != 0 && s.Max < s.Min) {
			return errors.New("invalid shift max")
		}
		if s.Max > 0 && len(s.SignedUp) > s.Max {
			return errors.New("too many signups for shift")
		}
		var seenPerson = make(map[model.PersonID]*model.Person)
		for _, p := range s.SignedUp {
			if seenPerson[p] != nil {
				return errors.New("duplicate person in shift")
			}
			seenPerson[p] = tx.FetchPerson(p)
			if seenPerson[p] == nil {
				return errors.New("nonexistent person in shift signup")
			}
			seenSignups[p] = append(seenSignups[p], s)
		}
		sort.Slice(s.SignedUp, func(i, j int) bool {
			return seenPerson[s.SignedUp[i]].SortName < seenPerson[s.SignedUp[j]].SortName
		})
		for _, ss := range seenSignups {
			if hasShiftOverlaps(ss) {
				return errors.New("person signed up for overlapping shifts")
			}
		}
		for _, p := range s.Declined {
			if seenPerson[p] != nil {
				return errors.New("duplicate person in shift")
			}
			seenPerson[p] = tx.FetchPerson(p)
			if seenPerson[p] == nil {
				return errors.New("nonexistent person in shift decline")
			}
		}
		sort.Slice(s.Declined, func(i, j int) bool {
			return seenPerson[s.Declined[i]].SortName < seenPerson[s.Declined[j]].SortName
		})
	}
	for _, m := range seenShifts {
		if len(m) > 1 && m[""] {
			return errors.New("empty shift task in shift with duplicate start")
		}
	}
	sort.Slice(event.Shifts, func(i, j int) bool {
		if event.Shifts[i].Start != event.Shifts[j].Start {
			return event.Shifts[i].Start < event.Shifts[j].Start
		}
		return event.Shifts[i].Task < event.Shifts[j].Task
	})
	for _, e := range tx.FetchEvents(event.Date, event.Date) {
		if e.ID != event.ID && e.Name == event.Name {
			return errors.New("duplicate name")
		}
	}
	return nil
}

func hasShiftOverlaps(ss []*model.Shift) bool {
	for i := range ss {
		for j := i + 1; j < len(ss); j++ {
			if ss[i].Start < ss[j].End && ss[i].End > ss[j].Start {
				return true
			}
		}
	}
	return false
}
