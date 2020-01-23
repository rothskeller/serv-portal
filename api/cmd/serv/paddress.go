package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"regexp"
	"strconv"

	"rothskeller.net/serv/model"
)

func listPersonAddresses(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, p := range matchPeople(args[0]) {
		for i, a := range p.Addresses {
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FormalName, strconv.Itoa(i + 1), a.Label, a.Address, a.City, a.State, a.Zip, strconv.FormatBool(a.WorkHours), strconv.FormatBool(a.HomeHours), strconv.FormatBool(a.MailingAddress), strconv.FormatFloat(a.Latitude, 'f', -1, 64), strconv.FormatFloat(a.Longitude, 'f', -1, 64)})

		}
	}
	cw.Flush()
}

func addPersonAddress(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FormalName)
		}
		os.Exit(1)
	}
	var address model.PersonAddress
	address.Address = args[1]
	address.City = args[2]
	address.State = args[3]
	address.Zip = args[4]
	applyPersonAddressFields(&address, fields)
	for _, e := range people[0].Addresses {
		if e.Address == address.Address && e.City == address.City && e.State == address.State && e.Zip == address.Zip {
			fmt.Fprintf(os.Stderr, "ERROR: %s already has this address.\n", people[0].FormalName)
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
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FormalName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Addresses) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid address index for %s.\n", args[1], people[0].FormalName)
	}
	applyPersonAddressFields(people[0].Addresses[index-1], fields)
	tx.SavePerson(people[0])
}

func removePersonAddress(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FormalName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Addresses) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid address index for %s.\n", args[1], people[0].FormalName)
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
			address.City = v
		case "state":
			address.State = v
		case "zip":
			address.Zip = v
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

func setPersonHomeAddress(args []string, fields map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		var changed bool
		for f, v := range fields {
			switch f {
			case "addr", "address":
				if p.HomeAddress.Address != v {
					changed = true
					p.HomeAddress.Address = v
				}
			case "lat", "latitude":
				if l, err := strconv.ParseFloat(v, 64); err == nil && l >= -90 && l <= 90 {
					if p.HomeAddress.Latitude != l {
						changed = true
						p.HomeAddress.Latitude = l
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for latitude.\n", v)
					os.Exit(1)
				}
			case "long", "longitude":
				if l, err := strconv.ParseFloat(v, 64); err == nil && l >= -180 && l <= 180 {
					if p.HomeAddress.Longitude != l {
						changed = true
						p.HomeAddress.Longitude = l
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for longitude.\n", v)
					os.Exit(1)
				}
			default:
				fmt.Fprintf(os.Stderr, "ERROR: there is no field %q on a person's home address.  Valid fields are address, latitude, and longitude.\n", v)
				os.Exit(1)
			}
		}
		if p.HomeAddress.Address == "" && (p.HomeAddress.Latitude != 0 || p.HomeAddress.Longitude != 0) {
			fmt.Fprintf(os.Stderr, "ERROR: latitude and longitude cannot be set without setting address.\n")
			os.Exit(1)
		}
		if changed {
			changes++
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}

func setPersonMailAddress(args []string, fields map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		var changed bool
		for f, v := range fields {
			switch f {
			case "same", "same_as_home":
				if v != "true" {
					fmt.Fprintf(os.Stderr, "ERROR: same_as_home can only be set to true.  (To unset it, set the address field.)\n")
					os.Exit(1)
				}
				if !p.MailAddress.SameAsHome {
					changed = true
					p.MailAddress = model.Address{SameAsHome: true}
				}
			case "addr", "address":
				if p.MailAddress.Address != v {
					changed = true
					p.MailAddress.Address = v
					p.MailAddress.SameAsHome = false
				}
			default:
				fmt.Fprintf(os.Stderr, "ERROR: there is no field %q on a person's mailing address.  Valid fields are same_as_home and address.\n", v)
				os.Exit(1)
			}
		}
		if changed {
			changes++
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}

func setPersonWorkAddress(args []string, fields map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		var changed bool
		for f, v := range fields {
			switch f {
			case "same", "same_as_home":
				if v != "true" {
					fmt.Fprintf(os.Stderr, "ERROR: same_as_home can only be set to true.  (To unset it, set the address field.)\n")
					os.Exit(1)
				}
				if !p.WorkAddress.SameAsHome {
					changed = true
					p.WorkAddress = model.Address{SameAsHome: true}
				}
			case "addr", "address":
				if p.WorkAddress.Address != v {
					changed = true
					p.WorkAddress.Address = v
					p.WorkAddress.SameAsHome = false
				}
			case "lat", "latitude":
				if l, err := strconv.ParseFloat(v, 64); err == nil && l >= -90 && l <= 90 {
					if p.WorkAddress.Latitude != l {
						changed = true
						p.WorkAddress.Latitude = l
						p.WorkAddress.SameAsHome = false
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for latitude.\n", v)
					os.Exit(1)
				}
			case "long", "longitude":
				if l, err := strconv.ParseFloat(v, 64); err == nil && l >= -180 && l <= 180 {
					if p.WorkAddress.Longitude != l {
						changed = true
						p.WorkAddress.Longitude = l
						p.WorkAddress.SameAsHome = false
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for longitude.\n", v)
					os.Exit(1)
				}
			default:
				fmt.Fprintf(os.Stderr, "ERROR: there is no field %q on a person's work address.  Valid fields are same_as_home, address, latitude, and longitude.\n", v)
				os.Exit(1)
			}
		}
		if p.WorkAddress.Address == "" && (p.WorkAddress.Latitude != 0 || p.WorkAddress.Longitude != 0) {
			fmt.Fprintf(os.Stderr, "ERROR: latitude and longitude cannot be set without setting address.\n")
			os.Exit(1)
		}
		if changed {
			changes++
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}

func removePersonHomeAddress(args []string, _ map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		if p.HomeAddress.Address != "" {
			changes++
			p.HomeAddress = model.Address{}
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}

func removePersonMailAddress(args []string, _ map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		if p.MailAddress.Address != "" || p.MailAddress.SameAsHome {
			changes++
			p.MailAddress = model.Address{}
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}

func removePersonWorkAddress(args []string, _ map[string]string) {
	var changes int
	matches := matchPeople(args[0])
	for _, p := range matches {
		if p.WorkAddress.Address != "" || p.WorkAddress.SameAsHome {
			changes++
			p.WorkAddress = model.Address{}
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", len(matches), changes)
}
