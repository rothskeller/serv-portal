package main

import (
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/role"
)

func convertRoles(tx *ostore.Tx, st *nstore.Store) {
	for _, orole := range tx.FetchRoles() {
		var flags role.Flags
		if orole.ShowRoster {
			flags |= role.Filter
		}
		if orole.ImplicitOnly {
			flags |= role.ImplicitOnly
		}
		var nrole = role.Updater{
			ID:        role.ID(orole.ID),
			Name:      orole.Name,
			Title:     orole.Title,
			Priority:  uint(orole.Priority),
			Org:       enum.Org(orole.Org),
			PrivLevel: enum.PrivLevel(orole.PrivLevel),
			Flags:     flags,
		}
		nr := role.Create(st, &nrole)
		for lid, rtl := range orole.Lists {
			if rtl.Sender() || rtl.SubModel() != 0 {
				listrole.SetListRole(st, list.WithID(st, list.ID(lid)), nr, rtl.Sender(), listrole.SubscriptionModel(rtl.SubModel()))
			}
		}
	}
	// Add role implications as a second pass to avoid foreign key violations.
	for _, orole := range tx.FetchRoles() {
		if len(orole.Implies) == 0 {
			continue
		}
		var nrole = role.WithID(st, role.ID(orole.ID), role.UpdaterFields)
		var updater = nrole.Updater()
		for implies, explicit := range orole.Implies {
			if explicit {
				updater.Implies = append(updater.Implies, role.ID(implies))
			}
		}
		nrole.Update(st, updater)
	}
}
