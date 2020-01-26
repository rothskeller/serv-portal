package main

import (
	"fmt"
	"os"
	"strconv"

	"sunnyvaleserv.org/portal/model"
)

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
