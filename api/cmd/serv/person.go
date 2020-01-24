package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
)

func listPeople(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, p := range matchPeople(args[0]) {
		cw.Write([]string{strconv.Itoa(int(p.ID)), p.Username, p.InformalName, p.FormalName, p.SortName, p.CallSign, strconv.Itoa(p.BadLoginCount), formatTime(p.BadLoginTime), p.PWResetToken, formatTime(p.PWResetTime), p.HomeAddress.Address, strconv.FormatFloat(p.HomeAddress.Latitude, 'f', -1, 64), strconv.FormatFloat(p.HomeAddress.Longitude, 'f', -1, 64), strconv.Itoa(p.HomeAddress.FireDistrict), strconv.FormatBool(p.MailAddress.SameAsHome), p.MailAddress.Address, strconv.FormatBool(p.WorkAddress.SameAsHome), p.WorkAddress.Address, strconv.FormatFloat(p.WorkAddress.Latitude, 'f', -1, 64), strconv.FormatFloat(p.WorkAddress.Longitude, 'f', -1, 64), strconv.Itoa(p.WorkAddress.FireDistrict), p.CellPhone, p.HomePhone, p.WorkPhone})
	}
	cw.Flush()
}

func matchPeople(pattern string) (people []*model.Person) {
	id, re, single := parsePattern(pattern)
	for _, p := range tx.FetchPeople() {
		if id != 0 && id != int(p.ID) {
			continue
		}
		if re != nil && !re.MatchString(p.InformalName) && !re.MatchString(p.FormalName) && !re.MatchString(p.CallSign) && !re.MatchString(p.Username) {
			continue
		}
		people = append(people, p)
	}
	if single && len(people) > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", pattern)
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FormalName)
		}
		os.Exit(1)
	}
	if pattern != "" && len(people) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q does not match any event.\n", pattern)
		os.Exit(1)
	}
	return people
}

func createPerson(_ []string, fields map[string]string) {
	var person model.Person
	applyPersonFields(&person, fields)
	tx.SavePerson(&person)
	fmt.Printf("Created person has ID %d.\n", person.ID)
}

func setPeople(args []string, fields map[string]string) {
	var matches, changes int
	for _, p := range matchPeople(args[0]) {
		matches++
		if applyPersonFields(p, fields) {
			changes++
			tx.SavePerson(p)
		}
	}
	fmt.Printf("Matched %d people and changed %d.\n", matches, changes)
}

func applyPersonFields(person *model.Person, fields map[string]string) (changed bool) {
	var err error
	for f, v := range fields {
		switch f {
		case "username":
			if person.Username != v {
				changed = true
				person.Username = v
			}
		case "informal", "informal_name":
			if person.InformalName != v {
				changed = true
				person.InformalName = v
			}
		case "formalname", "formal_name", "formal":
			if person.FormalName != v {
				changed = true
				person.FormalName = v
			}
		case "sortname", "sort_name", "sort":
			if person.SortName != v {
				changed = true
				person.SortName = v
			}
		case "callsign", "call_sign", "call":
			if person.CallSign != v {
				changed = true
				person.CallSign = v
			}
		case "cell", "cell_phone":
			if person.CellPhone != v {
				changed = true
				person.CellPhone = v
			}
		case "home_phone":
			if person.HomePhone != v {
				changed = true
				person.HomePhone = v
			}
		case "work_phone":
			if person.WorkPhone != v {
				changed = true
				person.WorkPhone = v
			}
		case "password", "pw":
			if v == "" {
				if len(person.Password) != 0 {
					changed = true
					person.Password = nil
				}
			} else if strings.HasPrefix(v, "$2a$") {
				if string(person.Password) != v {
					changed = true
					person.Password = []byte(v)
				}
			} else {
				changed = true
				person.Password = auth.EncryptPassword(v)
			}
		case "bad_login_count":
			if count, err := strconv.Atoi(v); err == nil && count >= 0 {
				if count != person.BadLoginCount {
					changed = true
					person.BadLoginCount = count
				}
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for bad_login_count.\n", v)
				os.Exit(1)
			}
		case "bad_login_time":
			var t time.Time
			if v == "" {
				t = time.Time{}
			} else if t, err = time.ParseInLocation(time.RFC3339, v, time.Local); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for bad_login_time. Use an RFC3339 timestamp.\n", v)
				os.Exit(1)
			}
			if !t.Equal(person.BadLoginTime) {
				changed = true
				person.BadLoginTime = t
			}
		case "pwreset_token":
			if person.PWResetToken != v {
				changed = true
				person.PWResetToken = v
			}
		case "pwreset_time":
			var t time.Time
			if v == "" {
				t = time.Time{}
			} else if t, err = time.ParseInLocation(time.RFC3339, v, time.Local); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for pwreset_time. Use an RFC3339 timestamp.\n", v)
				os.Exit(1)
			}
			if !t.Equal(person.PWResetTime) {
				changed = true
				person.PWResetTime = t
			}
		default:
			fmt.Fprintf(os.Stderr, "ERROR: there is no %q field on a person.  Valid fields are username, informal_name, formal_name, sort_name, call_sign, cell_phone, home_phone, work_phone, password, bad_login_count, bad_login_time, pwreset_token, and pwreset_time, and abbreviations of those.\n", f)
			os.Exit(1)
		}
	}
	if person.Username != "" {
		if p := tx.FetchPersonByUsername(person.Username); p != nil && p.ID != person.ID {
			fmt.Fprintf(os.Stderr, "ERROR: another person has this username.\n")
			os.Exit(1)
		}
	}
	if person.InformalName == "" {
		fmt.Fprintf(os.Stderr, "ERROR: informal_name is required.\n")
		os.Exit(1)
	}
	if person.FormalName == "" {
		fmt.Fprintf(os.Stderr, "ERROR: formal_name is required.\n")
		os.Exit(1)
	}
	if person.SortName == "" {
		fmt.Fprintf(os.Stderr, "ERROR: sort_name is required.\n")
		os.Exit(1)
	}
	if person.PWResetToken != "" {
		if p := tx.FetchPersonByPWResetToken(person.PWResetToken); p != nil && p.ID != person.ID {
			fmt.Fprintf(os.Stderr, "ERROR: another person has this pwreset_token.\n")
			os.Exit(1)
		}
	}
	return changed
}
