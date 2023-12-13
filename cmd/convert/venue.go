package main

import (
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/venue"
)

func convertVenues(tx *ostore.Tx, st *nstore.Store) {
	for _, ovenue := range tx.FetchVenues() {
		var name = ovenue.Name
		if nn := renameVenues[name]; nn != "" {
			name = nn
		}
		var flags venue.Flag
		if ovenue.ID == 23 /* via Zoom */ || ovenue.ID == 17 /* see details */ {
			flags |= venue.CanOverlap
		}
		var nvenue = venue.Updater{
			ID:    venue.ID(ovenue.ID),
			Name:  name,
			URL:   ovenue.URL,
			Flags: flags,
		}
		venue.Create(st, &nvenue)
	}
}

var renameVenues = map[string]string{
	"DPS EOC Training Room (105)":                            "EOC Training Room (105)",
	"DPS Headquarters Classroom (2028)":                      "HQ Classroom (2028)",
	"DPS Headquarters Downstairs Conference Room (1090)":     "HQ Downstairs Conference Room (1090)",
	"Mountain View Police/Fire Headquarters":                 "Mountain View Police/Fire HQ",
	"Santa Clara County Sheriff's Office Auditorium":         "Sheriff's Office Auditorium",
	"Santa Clara Fire Department Training Center":            "SCFD Training Center",
	"Sunnyvale Baylands Park":                                "Baylands Park",
	"Sunnyvale Community Center":                             "Community Center",
	"Sunnyvale Community Center, Recreation Center Ballroom": "Recreation Center Ballroom",
	"Sunnyvale Elementary School District Office":            "Sunnyvale School District Office",
	"via Zoom (see details)":                                 "Zoom Meeting",
}
