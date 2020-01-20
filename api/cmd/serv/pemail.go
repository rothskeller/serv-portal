package main

import (
	"encoding/csv"
	"os"
	"strconv"
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
