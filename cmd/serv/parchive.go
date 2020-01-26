package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func listPersonArchive(args []string, _ map[string]string) {
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, p := range matchPeople(args[0]) {
		for _, a := range p.Archive {
			pair := strings.SplitN(a, "=", 2)
			cw.Write([]string{strconv.Itoa(int(p.ID)), p.FormalName, pair[0], pair[1]})
		}
	}
	cw.Flush()
}

func addPersonArchive(args []string, _ map[string]string) {
	people := matchPeople(args[0])
	args = args[1:]
	if len(args)%2 != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: expecting even number of arguments (field value pairs)\n")
		os.Exit(1)
	}
	for _, p := range people {
		for i := 0; i < len(args); i += 2 {
			p.Archive = append(p.Archive, args[i]+"="+args[i+1])
		}
		tx.SavePerson(p)
	}
	fmt.Printf("Matched and changed %d people.\n", len(people))
}

func setPersonArchive(args []string, _ map[string]string) {
	people := matchPeople(args[0])
	args = args[1:]
	if len(args)%2 != 0 {
		fmt.Fprintf(os.Stderr, "ERROR: expecting even number of arguments (field value pairs)\n")
		os.Exit(1)
	}
	keys := make(map[string]bool)
	for i := 0; i < len(args); i += 2 {
		keys[args[i]] = true
	}
	for _, p := range people {
		j := 0
		for _, a := range p.Archive {
			key := strings.SplitN(a, "=", 2)[0]
			if !keys[key] {
				p.Archive[j] = a
				j++
			}
		}
		p.Archive = p.Archive[:j]
		for i := 0; i < len(args); i += 2 {
			p.Archive = append(p.Archive, args[i]+"="+args[i+1])
		}
		tx.SavePerson(p)
	}
	fmt.Printf("Matched and changed %d people.\n", len(people))
}

func removePersonArchive(args []string, _ map[string]string) {
	people := matchPeople(args[0])
	args = args[1:]
	keys := make(map[string]bool)
	for _, k := range args {
		keys[k] = true
	}
	for _, p := range people {
		j := 0
		for _, a := range p.Archive {
			key := strings.SplitN(a, "=", 2)[0]
			if !keys[key] {
				p.Archive[j] = a
				j++
			}
		}
		p.Archive = p.Archive[:j]
		tx.SavePerson(p)
	}
	fmt.Printf("Matched and changed %d people.\n", len(people))
}
