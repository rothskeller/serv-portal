package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"sunnyvaleserv.org/portal/model"
)

func listEvents(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, e := range matchEvents(args[0]) {
		var vname string
		if e.Venue != 0 {
			vname = tx.FetchVenue(e.Venue).Name
		}
		cw.Write([]string{strconv.Itoa(int(e.ID)), e.Name, e.Date, e.Start, e.End, strconv.Itoa(int(e.Venue)), vname, e.Details, eventTypesToString(e.Type), e.SccAresID})
	}
	cw.Flush()
}

func matchEvents(pattern string) (events []*model.Event) {
	id, re, single := parsePattern(pattern)
	for _, e := range tx.FetchEvents("2000-01-01", "2099-12-31") {
		if id != 0 && id != int(e.ID) {
			continue
		}
		if re != nil && !re.MatchString(e.Name) && !re.MatchString(string(e.Date)) {
			continue
		}
		events = append(events, e)
	}
	if single && len(events) > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple events:\n", pattern)
		for _, e := range events {
			fmt.Fprintf(os.Stderr, "%d\t%s\t%s\n", e.ID, e.Date, e.Name)
		}
		os.Exit(1)
	}
	if pattern != "" && len(events) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q does not match any event.\n", pattern)
		os.Exit(1)
	}
	return events
}

func eventTypesToString(etype model.EventType) string {
	var sb strings.Builder
	first := true
	for _, et := range model.AllEventTypes {
		if etype&et == 0 {
			continue
		}
		if first {
			first = false
		} else {
			sb.WriteByte(',')
		}
		sb.WriteString(model.EventTypeNames[et])
	}
	return sb.String()
}

var dateRE = regexp.MustCompile(`^20\d\d-(?:0[1-9]|1[0-2])-(?:0[1-9]|[12][0-9]|3[01])$`)
var timeRE = regexp.MustCompile(`^(?:[01]\d|2[0-3]):[0-5]\d$`)

func createEvent(_ []string, fields map[string]string) {
	var event model.Event
	for f, v := range fields {
		switch f {
		case "name":
			event.Name = v
		case "date":
			event.Date = v
		case "start":
			event.Start = v
		case "end":
			event.End = v
		case "venue":
			venues := matchVenues(v, true)
			event.Venue = venues[0].ID
		case "details":
			event.Details = v
		case "type":
			event.Type |= stringToEventTypes(v)
		case "scc_ares_id":
			event.SccAresID = v
		default:
			fmt.Fprintf(os.Stderr, "ERROR: unknown attribute %q for event\n", f)
			os.Exit(1)
		}
	}
	if event.Name == "" {
		fmt.Fprintf(os.Stderr, "ERROR: event name is required\n")
		os.Exit(1)
	}
	if event.Date == "" {
		fmt.Fprintf(os.Stderr, "ERROR: event date is required\n")
		os.Exit(1)
	}
	if !dateRE.MatchString(event.Date) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid date\n", event.Date)
		os.Exit(1)
	}
	if event.Start == "" {
		fmt.Fprintf(os.Stderr, "ERROR: event start time is required\n")
		os.Exit(1)
	}
	if !timeRE.MatchString(event.Start) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid time\n", event.Start)
		os.Exit(1)
	}
	if event.End == "" {
		fmt.Fprintf(os.Stderr, "ERROR: event end time is required\n")
		os.Exit(1)
	}
	if !timeRE.MatchString(event.End) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid time\n", event.End)
		os.Exit(1)
	}
	if event.End < event.Start {
		fmt.Fprintf(os.Stderr, "ERROR: event end time must not be before event start time")
		os.Exit(1)
	}
	if event.Type == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: event type is required\n")
		os.Exit(1)
	}
	for _, e := range tx.FetchEvents(event.Date, event.Date) {
		if e.Name == event.Name {
			fmt.Fprintf(os.Stderr, "ERROR: an event named %q already exists on %s.\n", event.Name, event.Date)
			os.Exit(1)
		}
	}
	if event.SccAresID != "" && tx.FetchEventBySccAresID(event.SccAresID) != nil {
		fmt.Fprintf(os.Stderr, "ERROR: an event with SCC ARES ID %q already exists.\n", event.SccAresID)
		os.Exit(1)
	}
	tx.SaveEvent(&event)
	fmt.Printf("Created event has ID %d.\n", event.ID)
}

func stringToEventTypes(str string) (etype model.EventType) {
	for _, s := range strings.Split(str, ",") {
		found := false
		s = strings.TrimSpace(s)
		if et, err := strconv.Atoi(s); err == nil && et > 0 {
			etype |= model.EventType(et)
			continue
		}
		for _, et := range model.AllEventTypes {
			if s == model.EventTypeNames[et] {
				etype |= et
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid event type.\n", s)
			os.Exit(1)
		}
	}
	return etype
}

func setEvent(args []string, fields map[string]string) {
	ematches := matchEvents(args[0])
	for _, event := range ematches {
		for f, v := range fields {
			switch f {
			case "name":
				event.Name = v
			case "date":
				event.Date = v
			case "start":
				event.Start = v
			case "end":
				event.End = v
			case "venue":
				venues := matchVenues(v, true)
				event.Venue = venues[0].ID
			case "details":
				event.Details = v
			case "type":
				event.Type |= stringToEventTypes(v)
			case "scc_ares_id":
				event.SccAresID = v
			default:
				fmt.Fprintf(os.Stderr, "ERROR: unknown attribute %q for event\n", f)
				os.Exit(1)
			}
		}
		if event.Name == "" {
			fmt.Fprintf(os.Stderr, "ERROR: event name is required\n")
			os.Exit(1)
		}
		if event.Date == "" {
			fmt.Fprintf(os.Stderr, "ERROR: event date is required\n")
			os.Exit(1)
		}
		if !dateRE.MatchString(event.Date) {
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid date\n", event.Date)
			os.Exit(1)
		}
		if event.Start == "" {
			fmt.Fprintf(os.Stderr, "ERROR: event start time is required\n")
			os.Exit(1)
		}
		if !timeRE.MatchString(event.Start) {
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid time\n", event.Start)
			os.Exit(1)
		}
		if event.End == "" {
			fmt.Fprintf(os.Stderr, "ERROR: event end time is required\n")
			os.Exit(1)
		}
		if !timeRE.MatchString(event.End) {
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid time\n", event.End)
			os.Exit(1)
		}
		if event.End < event.Start {
			fmt.Fprintf(os.Stderr, "ERROR: event end time must not be before event start time")
			os.Exit(1)
		}
		if event.Type == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: event type is required\n")
			os.Exit(1)
		}
		for _, e := range tx.FetchEvents(event.Date, event.Date) {
			if e.Name == event.Name && e.ID != event.ID {
				fmt.Fprintf(os.Stderr, "ERROR: an event named %q already exists on %s.\n", event.Name, event.Date)
				os.Exit(1)
			}
		}
		if event.SccAresID != "" && tx.FetchEventBySccAresID(event.SccAresID) != nil {
			fmt.Fprintf(os.Stderr, "ERROR: an event with SCC ARES ID %q already exists.\n", event.SccAresID)
			os.Exit(1)
		}
		tx.SaveEvent(event)
	}
	fmt.Printf("Matched and set %d events.\n", len(ematches))
}

func deleteEvents(args []string, _ map[string]string) {
	ematches := matchEvents(args[0])
	for _, e := range ematches {
		tx.DeleteEvent(e)
	}
	fmt.Printf("Matched and deleted %d events.", len(ematches))
}
