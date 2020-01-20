package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"rothskeller.net/serv/model"
)

func listEventGroups(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, e := range matchEvents(args[0]) {
		for _, g := range e.Groups {
			gobj := tx.FetchGroup(g)
			cw.Write([]string{strconv.Itoa(int(e.ID)), e.Date, e.Name, strconv.Itoa(int(g)), gobj.Name})
		}
	}
	cw.Flush()
}

func addEventGroups(args []string, _ map[string]string) {
	var matches, changes int
	var gids = make(map[model.GroupID]bool)
	for _, gp := range args[1:] {
		for _, g := range matchGroups(gp) {
			gids[g.ID] = true
		}
	}
	for _, e := range matchEvents(args[0]) {
		matches++
		var already = make(map[model.GroupID]bool)
		for _, g := range e.Groups {
			already[g] = true
		}
		for g := range gids {
			if !already[g] {
				changes++
				e.Groups = append(e.Groups, g)
			}
		}
		tx.SaveEvent(e)
	}
	fmt.Printf("Matched %d events and changed %d.\n", matches, changes)
}

func setEventGroups(args []string, _ map[string]string) {
	var matches int
	var gids = make(map[model.GroupID]bool)
	for _, gp := range args[1:] {
		for _, g := range matchGroups(gp) {
			gids[g.ID] = true
		}
	}
	for _, e := range matchEvents(args[0]) {
		matches++
		e.Groups = e.Groups[:0]
		for g := range gids {
			e.Groups = append(e.Groups, g)
		}
		tx.SaveEvent(e)
	}
	fmt.Printf("Matched and changed %d events.\n", matches)
}

func removeEventGroups(args []string, _ map[string]string) {
	var matches, changes int
	var gids = make(map[model.GroupID]bool)
	for _, gp := range args[1:] {
		for _, g := range matchGroups(gp) {
			gids[g.ID] = true
		}
	}
	for _, e := range matchEvents(args[0]) {
		matches++
		var already = make(map[model.GroupID]bool)
		for _, g := range e.Groups {
			already[g] = true
		}
		e.Groups = e.Groups[:0]
		for g := range already {
			if !gids[g] {
				e.Groups = append(e.Groups, g)
			} else {
				changes++
			}
		}
		if len(e.Groups) == 0 {
			fmt.Fprintf(os.Stderr, "ERROR: group removal would leave %d %s %q without any groups.\n", e.ID, e.Date, e.Name)
			os.Exit(1)
		}
		tx.SaveEvent(e)
	}
	fmt.Printf("Matched %d events and changed %d.\n", matches, changes)
}
