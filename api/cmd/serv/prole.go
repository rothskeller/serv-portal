package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"rothskeller.net/serv/model"
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

func addPersonRoles(args []string, _ map[string]string) {
	changePersonRoles(args, func(hasRole, changeRequested bool) bool {
		return hasRole || changeRequested
	})
}
func removePersonRoles(args []string, _ map[string]string) {
	changePersonRoles(args, func(hasRole, changeRequested bool) bool {
		return hasRole && !changeRequested
	})
}
func setPersonRoles(args []string, _ map[string]string) {
	changePersonRoles(args, func(hasRole, changeRequested bool) bool {
		return changeRequested
	})
}

func changePersonRoles(args []string, shouldHave func(has, change bool) bool) {
	var (
		people        []*model.Person
		changes       int
		changeRoleIDs = map[model.RoleID]bool{}
	)
	people = matchPeople(args[0])
	for _, rpatt := range args[1:] {
		for _, role := range matchRoles(rpatt) {
			changeRoleIDs[role.ID] = true
		}
	}
	for _, p := range people {
		var hasRoleIDs = map[model.RoleID]bool{}
		for _, rid := range p.Roles {
			hasRoleIDs[rid] = true
		}
		p.Roles = p.Roles[:0]
		for _, role := range tx.FetchRoles() {
			shouldHaveRole := shouldHave(hasRoleIDs[role.ID], changeRoleIDs[role.ID])
			if shouldHaveRole != hasRoleIDs[role.ID] {
				changes++
			}
			if shouldHaveRole {
				p.Roles = append(p.Roles, role.ID)
			}
		}
		tx.SavePerson(p)
	}
	fmt.Printf("Matched %d people times %d roles; made %d changes.\n", len(people), len(changeRoleIDs), changes)
}
