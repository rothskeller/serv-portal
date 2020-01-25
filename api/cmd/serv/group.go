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

func setGroups(args []string, fields map[string]string) {
	var changes int
	matches := matchGroups(args[0])
	for _, g := range matches {
		var changed bool
		for f, v := range fields {
			switch f {
			case "tag":
				if g.Tag != model.GroupTag(v) {
					changed = true
					g.Tag = model.GroupTag(v)
				}
			case "name":
				if g.Name != v {
					changed = true
					g.Name = v
				}
			case "text", "allow_text", "allow_text_messages":
				if b, err := strconv.ParseBool(v); err == nil {
					if g.AllowTextMessages != b {
						changed = true
						g.AllowTextMessages = b
					}
				} else {
					fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid value for allow_text_messages.\n", v)
					os.Exit(1)
				}
			default:
				fmt.Fprintf(os.Stderr, "ERROR: there is no field %q for a group.  Valid fields are tag, name, and allow_text_messages.\n", f)
				os.Exit(1)
			}
		}
		if changed {
			changes++
		}
	}
	if changes != 0 {
		tx.SaveAuthz()
	}
	fmt.Printf("Matched %d groups, changed %d.\n", len(matches), changes)
}
