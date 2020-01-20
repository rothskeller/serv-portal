package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"rothskeller.net/serv/model"
)

func listPersonPhones(args []string, _ map[string]string) {
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
		for i, pp := range p.Phones {
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FullName, strconv.Itoa(i + 1), pp.Label, pp.Phone, strconv.FormatBool(pp.SMS)})

		}
	}
	cw.Flush()
}

func addPersonPhone(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	var phone model.PersonPhone
	phone.Phone = formatPhone(args[1])
	applyPersonPhoneFields(&phone, fields)
	for _, e := range people[0].Phones {
		if e.Phone == phone.Phone {
			fmt.Fprintf(os.Stderr, "ERROR: %s already has this phone number.\n", people[0].FullName)
			os.Exit(1)
		}
	}
	people[0].Phones = append(people[0].Phones, &phone)
	tx.SavePerson(people[0])
}

func setPersonPhone(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Phones) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid phone index for %s.\n", args[1], people[0].FullName)
	}
	applyPersonPhoneFields(people[0].Phones[index-1], fields)
	tx.SavePerson(people[0])
}

func removePersonPhone(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Phones) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid phone index for %s.\n", args[1], people[0].FullName)
	}
	people[0].Phones = append(people[0].Phones[:index-1], people[0].Phones[index:]...)
	tx.SavePerson(people[0])
}

func applyPersonPhoneFields(phone *model.PersonPhone, fields map[string]string) {
	for f, v := range fields {
		switch f {
		case "phone":
			phone.Phone = formatPhone(v)
		case "label":
			phone.Label = v
		case "sms":
			if b, err := strconv.ParseBool(v); err == nil {
				phone.SMS = b
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for sms.\n", v)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "ERROR: there is no %q field on a person's phone.  Valid fields are phone, label, and sms.\n", f)
			os.Exit(1)
		}
	}
}

func formatPhone(phone string) string {
	ph := strings.Map(keepDigits, phone)
	if len(ph) != 10 {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid phone number.\n", phone)
		os.Exit(1)
	}
	return ph[0:3] + "-" + ph[3:6] + "-" + ph[6:10]
}
func keepDigits(r rune) rune {
	if r >= '0' && r <= '9' {
		return r
	}
	return -1
}
