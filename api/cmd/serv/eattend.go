package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"rothskeller.net/serv/model"
)

func listEventAttendance(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, e := range matchEvents(args[0]) {
		attend := tx.FetchAttendanceByEvent(e)
		for _, p := range tx.FetchPeople() {
			ai, ok := attend[p.ID]
			if !ok {
				continue
			}
			cw.Write([]string{strconv.Itoa(int(e.ID)), e.Date, e.Name, strconv.Itoa(int(p.ID)), p.InformalName, model.AttendanceTypeNames[ai.Type], strconv.Itoa(int(ai.Minutes))})
		}
	}
	cw.Flush()
}

func removeEventAttendance(args []string, _ map[string]string) {
	var changes int
	var ematches = matchEvents(args[0])
	var pmatches = matchPeople(args[1])
	for _, e := range ematches {
		attend := tx.FetchAttendanceByEvent(e)
		for _, p := range pmatches {
			if _, ok := attend[p.ID]; ok {
				changes++
				delete(attend, p.ID)
			}
		}
		tx.SaveEventAttendance(e, attend)
	}
	fmt.Printf("Matched %d events times %d people; changed %d attendance flags.\n", len(ematches), len(pmatches), changes)
}

func setEventAttendance(args []string, fields map[string]string) {
	var adds, changes int
	var ematches = matchEvents(args[0])
	var pmatches = matchPeople(args[1])
	var ai = model.AttendanceInfo{Type: 0xFF, Minutes: 10000}
	for f, v := range fields {
		switch f {
		case "type":
			switch v {
			case "volunteer":
				ai.Type = model.AttendAsVolunteer
			case "student":
				ai.Type = model.AttendAsStudent
			case "audit", "auditor":
				ai.Type = model.AttendAsVolunteer
			default:
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid attendance type.  Valid types are volunteer, student, and audit.\n", v)
				os.Exit(1)
			}
		case "minutes", "min":
			if m, err := strconv.Atoi(v); err == nil && m >= 0 && m <= 1440 {
				ai.Minutes = uint16(m)
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid number of minutes.\n", v)
				os.Exit(1)
			}
		case "hours", "hour":
			if h, err := strconv.ParseFloat(v, 64); err == nil && h >= 0 && h <= 24 {
				ai.Minutes = uint16(h * 60)
			} else {
				fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid number of minutes.\n", v)
				os.Exit(1)
			}
		default:
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a valid field for event attendance.  Valid fields are type, minutes, and hours.\n", f)
			os.Exit(1)
		}
	}
	for _, e := range ematches {
		attend := tx.FetchAttendanceByEvent(e)
		for _, p := range pmatches {
			changed := false
			var nai model.AttendanceInfo
			nai, ok := attend[p.ID]
			if !ok {
				changed = true
				adds++
			}
			if nai.Type != ai.Type && ai.Type != 0xFF {
				changed = true
				nai.Type = ai.Type
			}
			if nai.Minutes != ai.Minutes && ai.Minutes != 10000 {
				changed = true
				nai.Minutes = ai.Minutes
			}
			attend[p.ID] = nai
			if changed {
				changes++
			}
		}
		tx.SaveEventAttendance(e, attend)
	}
	fmt.Printf("Matched %d events times %d people; added %d attendances and changed %d.\n", len(ematches), len(pmatches), adds, changes-adds)
}
