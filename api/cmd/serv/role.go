package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"rothskeller.net/serv/model"
)

func listRoles(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, r := range matchRoles(args[0]) {
		cw.Write([]string{strconv.Itoa(int(r.ID)), string(r.Tag), r.Name, strconv.FormatBool(r.Individual)})
	}
	cw.Flush()
}

func matchRoles(pattern string) (roles []*model.Role) {
	id, re, single := parsePattern(pattern)
	for _, r := range tx.FetchRoles() {
		if id != 0 && id != int(r.ID) {
			continue
		}
		if re != nil && !re.MatchString(r.Name) && !re.MatchString(string(r.Tag)) {
			continue
		}
		roles = append(roles, r)
	}
	if single && len(roles) > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple roles:\n", pattern)
		for _, g := range roles {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", g.ID, g.Name)
		}
		os.Exit(1)
	}
	if pattern != "" && len(roles) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q does not match any role.\n", pattern)
		os.Exit(1)
	}
	return roles
}
