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

func listPersonEmails(args []string, _ map[string]string) {
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
		for i, e := range p.Emails {
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FullName, strconv.Itoa(i + 1), e.Label, e.Email, strconv.FormatBool(e.Bad)})

		}
	}
	cw.Flush()
}

func addPersonEmail(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	var email model.PersonEmail
	email.Email = strings.ToLower(args[1])
	applyPersonEmailFields(&email, fields)
	for _, e := range people[0].Emails {
		if e.Email == email.Email {
			fmt.Fprintf(os.Stderr, "ERROR: %s already has this email address.\n", people[0].FullName)
			os.Exit(1)
		}
	}
	people[0].Emails = append(people[0].Emails, &email)
	tx.SavePerson(people[0])
}

func setPersonEmail(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Emails) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid email index for %s.\n", args[1], people[0].FullName)
	}
	applyPersonEmailFields(people[0].Emails[index-1], fields)
	tx.SavePerson(people[0])
}

func removePersonEmail(args []string, fields map[string]string) {
	people := matchPeople(args[0])
	if len(people) != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple people:\n", args[0])
		for _, p := range people {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", p.ID, p.FullName)
		}
		os.Exit(1)
	}
	index, err := strconv.Atoi(args[1])
	if err != nil || index < 1 || index > len(people[0].Emails) {
		fmt.Fprintf(os.Stderr, "ERROR: %s is not a valid email index for %s.\n", args[1], people[0].FullName)
	}
	people[0].Emails = append(people[0].Emails[:index-1], people[0].Emails[index:]...)
	tx.SavePerson(people[0])
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func applyPersonEmailFields(email *model.PersonEmail, fields map[string]string) {
	for f, v := range fields {
		switch f {
		case "email":
			email.Email = strings.ToLower(v)
		case "label":
			email.Label = v
		case "bad":
			if b, err := strconv.ParseBool(v); err == nil {
				email.Bad = b
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for bad.\n", v)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "ERROR: there is no %q field on a person's email.  Valid fields are email, label, and bad.\n", f)
			os.Exit(1)
		}
	}
	if !emailRE.MatchString(email.Email) {
		fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid email address.\n", email.Email)
		os.Exit(1)
	}
}
