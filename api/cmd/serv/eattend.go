package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"rothskeller.net/serv/model"
)

func listEventAttendance(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, e := range matchEvents(args[0]) {
		var people []*model.Person
		for pid := range tx.FetchAttendanceByEvent(e) {
			people = append(people, tx.FetchPerson(pid))
		}
		sort.Sort(model.PersonSort(people))
		for _, p := range people {
			cw.Write([]string{strconv.Itoa(int(e.ID)), e.Date, e.Name, strconv.Itoa(int(p.ID)), p.FormalName})
		}
	}
	cw.Flush()
}

func addEventAttendance(args []string, _ map[string]string) {
	var changes int
	var ematches = matchEvents(args[0])
	var pmatches []*model.Person
	for _, ppatt := range args[1:] {
		pmatches = append(pmatches, matchPeople(ppatt)...)
	}
	for _, e := range ematches {
		attend := tx.FetchAttendanceByEvent(e)
		for _, p := range pmatches {
			if !attend[p.ID] {
				changes++
				attend[p.ID] = true
			}
		}
		tx.SaveEventAttendance(e, attend)
	}
	fmt.Printf("Matched %d events times %d people; changed %d attendance flags.\n", len(ematches), len(pmatches), changes)
}

func removeEventAttendance(args []string, _ map[string]string) {
	var changes int
	var ematches = matchEvents(args[0])
	var pmatches []*model.Person
	for _, ppatt := range args[1:] {
		pmatches = append(pmatches, matchPeople(ppatt)...)
	}
	for _, e := range ematches {
		attend := tx.FetchAttendanceByEvent(e)
		for _, p := range pmatches {
			if attend[p.ID] {
				changes++
				attend[p.ID] = false
			}
		}
		tx.SaveEventAttendance(e, attend)
	}
	fmt.Printf("Matched %d events times %d people; changed %d attendance flags.\n", len(ematches), len(pmatches), changes)
}

func setEventAttendance(args []string, _ map[string]string) {
	var ematches = matchEvents(args[0])
	var pmatches []*model.Person
	for _, ppatt := range args[1:] {
		pmatches = append(pmatches, matchPeople(ppatt)...)
	}
	for _, e := range ematches {
		attend := make(map[model.PersonID]bool)
		for _, p := range pmatches {
			attend[p.ID] = true
		}
		tx.SaveEventAttendance(e, attend)
	}
	fmt.Printf("Matched and set %d events times %d people.\n", len(ematches), len(pmatches))
}
