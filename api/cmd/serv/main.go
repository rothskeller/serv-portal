package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"rothskeller.net/serv/db"
)

func usage() {
	fmt.Fprint(os.Stderr, `usage: serv object-type [selectors] [command] [arguments]
    serv backup
    serv event [«pattern»] [list]
    serv event create name «name» date «date» start «start» end «end» type «type» [field «value»]...
    serv event [«pattern»|«id»] set [field «value»]...
    serv event [«pattern»|«id»] delete
    serv event [«pattern»|«id»] attendance [list]
    serv event [«pattern»|«id»] attendance add «pattern»|«id»...
    serv event [«pattern»|«id»] attendance remove «pattern»|«id»...
    serv event [«pattern»|«id»] attendance set «pattern»|«id»...
    serv event [«pattern»|«id»] group [list]
    serv event [«pattern»|«id»] group add «pattern»|«id»...
    serv event [«pattern»|«id»] group remove «pattern»|«id»...
    serv event [«pattern»|«id»] group set «pattern»|«id»...
    serv group [«pattern»] [list]
    serv group create [tag «tag»] name «name»
    serv group [«pattern»|«id»] set [tag «tag»] [name «name»]
    serv group [«pattern»|«id»] delete
    serv person [«pattern»] [list]
    serv person [«pattern»|«id»] set [field «value»]...
    serv person [«pattern»|«id»] address [list]
    serv person [«pattern»|«id»] address add «address» «city» «state» «zip» [field «value»]...
    serv person [«pattern»|«id»] address «index» set [field «value»]...
    serv person [«pattern»|«id»] address «index» remove
    serv person [«pattern»|«id»] email [list]
    serv person [«pattern»|«id»] email add «email» [field «value»]...
    serv person [«pattern»|«id»] email «index» set [field «value»]...
    serv person [«pattern»|«id»] email «index» remove
    serv person [«pattern»|«id»] phone [list]
    serv person [«pattern»|«id»] phone add «phone» [field «value»]...
    serv person [«pattern»|«id»] phone «index» set [field «value»]...
    serv person [«pattern»|«id»] phone «index» remove
    serv person [«pattern»|«id»] role [list]
    serv person [«pattern»|«id»] role add «pattern»|«id»...
    serv person [«pattern»|«id»] role remove «pattern»|«id»...
    serv person [«pattern»|«id»] role set «pattern»|«id»...
    serv role [«pattern»] [list]
    serv role create [tag «tag»] name «name» [individual «individual»]
    serv role [«pattern»|«id»] set [tag «tag»] [name «name»] [individual «individual»]
    serv role [«pattern»|«id»] privileges [«pattern»|«id»] [list]
    serv role [«pattern»|«id»] privileges [«pattern»|«id»] {set|add|remove} [member] [view] [contactInfo] [admin] [events]
    serv role [«pattern»|«id»] delete
    serv venue [«pattern»] [list]
    serv venue create name «name» [field «value»]...
    serv venue [«pattern»|«id»] set [field «value»]...
    serv venue [«pattern»|«id»] delete
`)
	os.Exit(2)
}

// A command structure describes one command that the user can invoke, and its
// argument pattern.  In the list of arguments, words are taken as required
// literal words, except for "list" which is optional.  Other symbols are:
//     "?" optional argument
//     "." required argument
//     "+" one or more arguments
//     "=" any even number of arguments, stored as a map[field]value
// The "?", ".", and "+" arguments are passed to the handler in the string slice
// (with an empty string passed for an optional "?" argument that was not
// present).  The "=" arguments are passed to the handler in the string map.
type command struct {
	arguments []string
	handler   func([]string, map[string]string)
}

var commands = []command{
	{[]string{"backup"}, makeBackup},
	{[]string{"event", "?", "list"}, listEvents},
	{[]string{"event", "create", "="}, createEvent},
	{[]string{"event", "?", "set", "="}, setEvent},
	{[]string{"event", "?", "delete"}, deleteEvents},
	{[]string{"event", "?", "attendance", "list"}, listEventAttendance},
	{[]string{"event", "?", "attendance", "add", "+"}, addEventAttendance},
	{[]string{"event", "?", "attendance", "remove", "+"}, removeEventAttendance},
	{[]string{"event", "?", "attendance", "set", "+"}, setEventAttendance},
	{[]string{"event", "?", "group", "list"}, listEventGroups},
	{[]string{"event", "?", "group", "add", "+"}, addEventGroups},
	{[]string{"event", "?", "group", "remove", "+"}, removeEventGroups},
	{[]string{"event", "?", "group", "set", "+"}, setEventGroups},
	{[]string{"group", "?", "list"}, listGroups},
	// {[]string{"group", "create", "="}, handler},
	// {[]string{"group", "?", "set", "="}, handler},
	// {[]string{"group", "?", "delete"}, handler},
	{[]string{"person", "?", "list"}, listPeople},
	{[]string{"person", "create", "="}, createPerson},
	{[]string{"person", "?", "set", "="}, setPeople},
	{[]string{"person", "?", "address", "list"}, listPersonAddresses},
	// {[]string{"person", "?", "address", "add", ".", ".", ".", ".", "="}, handler},
	// {[]string{"person", "?", "address", "«index»", "set", "="}, handler},
	// {[]string{"person", "?", "address", "«index»", "remove"}, handler},
	{[]string{"person", "?", "email", "list"}, listPersonEmails},
	{[]string{"person", "?", "email", "add", ".", "="}, addPersonEmail},
	{[]string{"person", "?", "email", "«index»", "set", "="}, setPersonEmail},
	{[]string{"person", "?", "email", "«index»", "remove"}, removePersonEmail},
	{[]string{"person", "?", "phone", "list"}, listPersonPhones},
	// {[]string{"person", "?", "phone", "add", ".", "="}, handler},
	// {[]string{"person", "?", "phone", "«index»", "set", "="}, handler},
	// {[]string{"person", "?", "phone", "«index»", "remove"}, handler},
	{[]string{"person", "?", "role", "list"}, listPersonRoles},
	{[]string{"person", "?", "role", "add", "+"}, addPersonRoles},
	{[]string{"person", "?", "role", "remove", "+"}, removePersonRoles},
	{[]string{"person", "?", "role", "set", "+"}, setPersonRoles},
	{[]string{"role", "?", "list"}, listRoles},
	// {[]string{"role", "create", "="}, handler},
	// {[]string{"role", "?", "set", "="}, handler},
	{[]string{"role", "?", "privilege", "?", "list"}, listRolePrivileges},
	{[]string{"role", "?", "privilege", "?", "set", "+"}, setRolePrivileges},
	{[]string{"role", "?", "privilege", "?", "add", "+"}, addRolePrivileges},
	{[]string{"role", "?", "privilege", "?", "remove", "+"}, removeRolePrivileges},
	// {[]string{"role", "?", "delete"}, handler},
	{[]string{"venue", "?", "list"}, listVenues},
	{[]string{"venue", "create", "="}, createVenue},
	// {[]string{"venue", "?", "set", "="}, handler},
	// {[]string{"venue", "?", "delete"}, handler},
}
var abbreviations = map[string]string{
	"-":          "remove",
	"+":          "add",
	"=":          "set",
	"addr":       "address",
	"addresses":  "address",
	"addrs":      "address",
	"attend":     "attendance",
	"emails":     "email",
	"events":     "event",
	"groups":     "group",
	"people":     "person",
	"phones":     "phone",
	"priv":       "privilege",
	"privileges": "privilege",
	"privs":      "privilege",
	"roles":      "role",
	"venues":     "venue",
}

var tx *db.Tx

func main() {
	maybeMakeBackup()
	if len(os.Args) < 2 {
		usage()
	}
	db.Open("data/serv.db")
	tx = db.Begin()
	for _, c := range commands {
		if match, argslice, argmap := matchCommand(c, os.Args[1:]); match {
			c.handler(argslice, argmap)
			tx.Commit()
			return
		}
	}
	tx.Rollback()
	usage()
}

func matchCommand(c command, args []string) (match bool, argslice []string, argmap map[string]string) {
	cmd := c.arguments
	for len(cmd) != 0 && len(args) != 0 {
		switch cmd[0] {
		case "?":
			if len(cmd) > 1 && abbrMatch(cmd[1], args[0]) {
				argslice = append(argslice, "")
				cmd = cmd[1:]
			} else {
				argslice = append(argslice, args[0])
				cmd, args = cmd[1:], args[1:]
			}
		case ".":
			argslice = append(argslice, args[0])
			cmd, args = cmd[1:], args[1:]
		case "+":
			argslice = append(argslice, args...)
			cmd, args = cmd[1:], args[:0]
		case "=":
			if len(args)%2 != 0 {
				fmt.Fprintf(os.Stderr, "ERROR: odd number of arguments in field/value pair list\n")
				usage()
			}
			argmap = make(map[string]string, len(args)/2)
			for i := 0; i < len(args); i += 2 {
				if _, ok := argmap[args[i]]; ok {
					fmt.Fprintf(os.Stderr, "ERROR: multiple values for field %q\n", args[i])
					usage()
				}
				argmap[args[i]] = args[i+1]
			}
			cmd, args = cmd[1:], args[:0]
		default:
			if !abbrMatch(cmd[0], args[0]) {
				return false, nil, nil
			}
			cmd, args = cmd[1:], args[1:]
		}
	}
	if len(cmd) == 2 && cmd[0] == "?" && cmd[1] == "list" && len(args) == 0 {
		argslice = append(argslice, "")
		return true, argslice, argmap
	}
	if len(cmd) == 1 && cmd[0] == "list" && len(args) == 0 {
		return true, argslice, argmap
	}
	if len(cmd) == 0 && len(args) == 0 {
		return true, argslice, argmap
	}
	return false, nil, nil
}

func abbrMatch(cmd, arg string) bool {
	if abbr := abbreviations[arg]; abbr != "" {
		return cmd == abbr
	}
	return cmd == arg
}

func parsePattern(pattern string) (id int, re *regexp.Regexp, single bool) {
	var err error
	if pattern == "" {
		return 0, nil, false
	}
	if id, err = strconv.Atoi(pattern); err == nil && id > 0 {
		return id, nil, true
	}
	if pattern[0] == '=' {
		single = true
		pattern = pattern[1:]
	}
	if re, err = regexp.Compile(pattern); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: invalid pattern %q: %s\n", pattern, err)
		os.Exit(1)
	}
	return 0, re, single
}

func formatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}
