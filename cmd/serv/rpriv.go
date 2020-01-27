package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"sunnyvaleserv.org/portal/model"
)

func listRolePrivileges(args []string, _ map[string]string) {
	rpattern, gpattern := args[0], args[1]
	cw := csv.NewWriter(os.Stdout)
	cw.Comma = '\t'
	for _, r := range matchRoles(rpattern) {
		for _, g := range matchGroups(gpattern) {
			privs := r.Privileges.Get(g)
			if privs&model.PrivMember != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "member"})
			}
			if privs&model.PrivViewMembers != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "roster"})
			}
			if privs&model.PrivViewContactInfo != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "contact"})
			}
			if privs&model.PrivManageMembers != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "admin"})
			}
			if privs&model.PrivManageEvents != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "events"})
			}
			if privs&model.PrivSendTextMessages != 0 {
				cw.Write([]string{strconv.Itoa(int(r.ID)), r.Name, strconv.Itoa(int(g.ID)), g.Name, "texts"})
			}
		}
	}
	cw.Flush()
}

func addRolePrivileges(args []string, _ map[string]string) {
	updateRolePrivileges(args, func(orig, mask model.Privilege) model.Privilege { return orig | mask })
}

func setRolePrivileges(args []string, _ map[string]string) {
	updateRolePrivileges(args, func(_, mask model.Privilege) model.Privilege { return mask })
}

func removeRolePrivileges(args []string, _ map[string]string) {
	updateRolePrivileges(args, func(orig, mask model.Privilege) model.Privilege { return orig &^ mask })
}

func updateRolePrivileges(args []string, updater func(orig, mask model.Privilege) model.Privilege) {
	var rmatches, gmatches, changes int
	rpattern, gpattern := args[0], args[1]
	var toChange model.Privilege
	for _, t := range args[2:] {
		switch t {
		case "member":
			toChange |= model.PrivMember
		case "roster":
			toChange |= model.PrivViewMembers
		case "contact":
			toChange |= model.PrivViewContactInfo
		case "admin":
			toChange |= model.PrivManageMembers
		case "events":
			toChange |= model.PrivManageEvents
		case "texts":
			toChange |= model.PrivSendTextMessages
		default:
			fmt.Fprintf(os.Stderr, "ERROR: %q is not a known privilege.  Privileges are member, roster, contact, admin, events, and texts.\n", t)
			os.Exit(1)
		}
	}
	for _, r := range matchRoles(rpattern) {
		rmatches++
		for _, g := range matchGroups(gpattern) {
			gmatches++
			orig := r.Privileges.Get(g)
			updated := updater(orig, toChange)
			if orig != updated {
				r.Privileges.Set(g, updated)
				changes++
			}
		}
	}
	if changes > 0 {
		tx.SaveAuthz()
	}
	fmt.Printf("Matched %d roles times %d groups, changed %d privileges.\n", rmatches, gmatches/rmatches, changes)
}
