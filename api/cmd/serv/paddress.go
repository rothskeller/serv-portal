package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"rothskeller.net/serv/model"
)

func listPersonAddresses(args []string, _ map[string]string) {
	id, re, _ := parsePattern(args[0])
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, p := range tx.FetchPeople() {
		if id != 0 && id != int(p.ID) {
			continue
		}
		if re != nil && !re.MatchString(p.FullName) && !re.MatchString(p.Nickname) && !re.MatchString(p.CallSign) && !re.MatchString(p.Username) {
			continue
		}
		for i, a := range p.Addresses {
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FullName, strconv.Itoa(i + 1), a.Label, a.Address, a.City, a.State, a.Zip, strconv.FormatBool(a.WorkHours), strconv.FormatBool(a.HomeHours), strconv.FormatBool(a.MailingAddress), strconv.FormatFloat(a.Latitude, 'f', -1, 64), strconv.FormatFloat(a.Longitude, 'f', -1, 64)})

		}
	}
	cw.Flush()
}

func addPersonAddress(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	var address model.PersonAddress
	address.Address = strings.ToLower(args[1])
	address.City = strings.ToLower(args[2])
	address.State = strings.ToLower(args[3])
	address.Zip = strings.ToLower(args[4])
	applyPersonAddressFields(&address, fields)
	for _, e := range people[0].Addresses {
		if e.Address == address.Address && e.City == address.City && e.State == address.State && e.Zip == address.Zip {
			fmt.Fprintf(os.Stderr, "ERROR: %s already has this address.\n", people[0].FullName)
			os.Exit(1)
		}
	}
	people[0].Addresses = append(people[0].Addresses, &address)
	tx.SavePerson(people[0])
}

func setPersonAddress(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Addresses) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid address index for %s.\n", args[1], people[0].FullName)
	}
	applyPersonAddressFields(people[0].Addresses[index-1], fields)
	tx.SavePerson(people[0])
}

func removePersonAddress(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Addresses) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid address index for %s.\n", args[1], people[0].FullName)
	}
	people[0].Addresses = append(people[0].Addresses[:index-1], people[0].Addresses[index:]...)
	tx.SavePerson(people[0])
}

var stateRE = regexp.MustCompile(`^[A-Z][A-Z]$`)
var zipRE = regexp.MustCompile(`^[0-9]{5}$`)

func applyPersonAddressFields(address *model.PersonAddress, fields map[string]string) {
	for f, v := range fields {
		switch f {
		case "address":
			address.Address = v
		case "city":
			address.Address = v
		case "state":
			address.Address = v
		case "zip":
			address.Address = v
		case "label":
			address.Label = v
		case "work", "work_hours":
			if b, err := strconv.ParseBool(v); err == nil {
				address.WorkHours = b
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for work_hours.\n", v)
				os.Exit(1)
			}
		case "home", "home_hours":
			if b, err := strconv.ParseBool(v); err == nil {
				address.HomeHours = b
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for home_hours.\n", v)
				os.Exit(1)
			}
		case "mailing", "mailing_address":
			if b, err := strconv.ParseBool(v); err == nil {
				address.MailingAddress = b
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for mailing_address.\n", v)
				os.Exit(1)
			}
		case "latitude":
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				address.Latitude = f
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for latitude.\n", v)
				os.Exit(1)
			}
		case "longitude":
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				address.Longitude = f
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for longitude.\n", v)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "ERROR: there is no %q field on a person's address.  Valid fields are address, city, state, zip, label, work_hours, home_hours, mailing_address, latitude, and longitude.\n", f)
			os.Exit(1)
		}
	}
	if !stateRE.MatchString(address.State) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid state.\n", address.State)
		os.Exit(1)
	}
	if !zipRE.MatchString(address.Zip) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid ZIP code.\n", address.Zip)
		os.Exit(1)
	}
}
