package main

import (
	"encoding/csv"
	"os"
	"strconv"
)

func listPersonRoles(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, p := range matchPeople(args[0]) {
		for _, rn := range p.Roles {
			r := tx.FetchRole(rn)
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FullName, strconv.Itoa(int(r.ID)), r.Name})
		}
	}
	cw.Flush()
}
