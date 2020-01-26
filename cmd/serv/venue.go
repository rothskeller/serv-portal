package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"sunnyvaleserv.org/portal/model"
)

func listVenues(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, v := range matchVenues(args[0], false) {
		cw.Write([]string{strconv.Itoa(int(v.ID)), v.Name, v.Address, v.City, v.URL})
	}
	cw.Flush()
}

func createVenue(_ []string, fields map[string]string) {
	var venue model.Venue
	for f, v := range fields {
		switch f {
		case "name":
			venue.Name = v
		case "address":
			venue.Address = v
		case "city":
			venue.City = v
		case "url":
			venue.URL = v
		default:
			fmt.Fprintf(os.Stderr, "ERROR: unknown attribute %q for venue\n", f)
			os.Exit(1)
		}
	}
	if venue.Name == "" {
		fmt.Fprintf(os.Stderr, "ERROR: venue name is required\n")
		os.Exit(1)
	}
	for _, v := range tx.FetchVenues() {
		if v.Name == venue.Name {
			fmt.Fprintf(os.Stderr, "ERROR: a venue named %q already exists.\n", venue.Name)
			os.Exit(1)
		}
	}
	tx.SaveVenue(&venue)
	fmt.Printf("Created venue has ID %d.\n", venue.ID)
}

func matchVenues(pattern string, single bool) (venues []*model.Venue) {
	id, re, _ := parsePattern(pattern)
	for _, v := range tx.FetchVenues() {
		if id != 0 && id != int(v.ID) {
			continue
		}
		if re != nil && !re.MatchString(v.Name) && !re.MatchString(string(v.City)) && !re.MatchString(string(v.Address)) {
			continue
		}
		venues = append(venues, v)
	}
	if single && len(venues) > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple venues:\n", pattern)
		for _, v := range venues {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", v.ID, v.Name)
		}
		os.Exit(1)
	}
	if pattern != "" && len(venues) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q does not match any venue.\n", pattern)
		os.Exit(1)
	}
	return venues
}
