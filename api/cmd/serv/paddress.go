package main

import (
	"encoding/csv"
	"os"
	"strconv"
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
