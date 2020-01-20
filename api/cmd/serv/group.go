package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"rothskeller.net/serv/model"
)

func listGroups(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, g := range matchGroups(args[0]) {
		cw.Write([]string{strconv.Itoa(int(g.ID)), string(g.Tag), g.Name})
	}
	cw.Flush()
}

func matchGroups(pattern string) (groups []*model.Group) {
	id, re, single := parsePattern(pattern)
	for _, g := range tx.FetchGroups() {
		if id != 0 && id != int(g.ID) {
			continue
		}
		if re != nil && !re.MatchString(g.Name) && !re.MatchString(string(g.Tag)) {
			continue
		}
		groups = append(groups, g)
	}
	if single && len(groups) > 1 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q matches multiple groups:\n", pattern)
		for _, g := range groups {
			fmt.Fprintf(os.Stderr, "%d\t%s\n", g.ID, g.Name)
		}
		os.Exit(1)
	}
	if pattern != "" && len(groups) == 0 {
		fmt.Fprintf(os.Stderr, "ERROR: pattern %q does not match any group.\n", pattern)
		os.Exit(1)
	}
	return groups
}
