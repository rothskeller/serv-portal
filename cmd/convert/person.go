package main

import (
	"fmt"
	"sort"
	"time"

	"sunnyvaleserv.org/portal/model"
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
)

func convertPeople(tx *ostore.Tx, st *nstore.Store) {
	for _, operson := range tx.FetchPeople() {
		var nperson = person.Updater{
			ID:               person.ID(operson.ID),
			VolgisticsID:     uint(operson.VolgisticsID),
			InformalName:     operson.InformalName,
			FormalName:       operson.FormalName,
			SortName:         operson.SortName,
			CallSign:         operson.CallSign,
			Email:            operson.Email,
			Email2:           operson.Email2,
			CellPhone:        operson.CellPhone,
			HomePhone:        operson.HomePhone,
			WorkPhone:        operson.WorkPhone,
			Password:         string(operson.Password),
			BadLoginCount:    uint(operson.BadLoginCount),
			BadLoginTime:     operson.BadLoginTime,
			PWResetToken:     operson.PWResetToken,
			PWResetTime:      operson.PWResetTime,
			UnsubscribeToken: operson.UnsubscribeToken,
			HoursToken:       operson.HoursToken,
			Identification:   person.IdentType(operson.Identification),
			Birthdate:        operson.Birthdate,
		}
		if operson.NoEmail {
			nperson.Flags |= person.NoEmail
		}
		if operson.NoText {
			nperson.Flags |= person.NoText
		}
		if operson.HoursReminder {
			nperson.Flags |= person.HoursReminder
		}
		if operson.VolgisticsPending {
			nperson.Flags |= person.VolgisticsPending
		}
		if operson.ID == 635 {
			nperson.Flags |= person.VisibleToAll
		}
		if operson.HomeAddress.Address != "" {
			nperson.Addresses.Home = &person.Address{
				Address:      operson.HomeAddress.Address,
				Latitude:     operson.HomeAddress.Latitude,
				Longitude:    operson.HomeAddress.Longitude,
				FireDistrict: uint(operson.HomeAddress.FireDistrict),
			}
		}
		if operson.WorkAddress.SameAsHome {
			nperson.Addresses.Work = &person.Address{SameAsHome: true}
		} else if operson.WorkAddress.Address != "" {
			nperson.Addresses.Work = &person.Address{
				Address:      operson.WorkAddress.Address,
				Latitude:     operson.WorkAddress.Latitude,
				Longitude:    operson.WorkAddress.Longitude,
				FireDistrict: uint(operson.WorkAddress.FireDistrict),
			}
		}
		if operson.MailAddress.SameAsHome {
			nperson.Addresses.Mail = &person.Address{SameAsHome: true}
		} else if operson.MailAddress.Address != "" {
			nperson.Addresses.Mail = &person.Address{
				Address:      operson.MailAddress.Address,
				Latitude:     operson.MailAddress.Latitude,
				Longitude:    operson.MailAddress.Longitude,
				FireDistrict: uint(operson.MailAddress.FireDistrict),
			}
		}
		sort.Slice(operson.BGChecks, func(i, j int) bool {
			return operson.BGChecks[i].Date < operson.BGChecks[j].Date
		})
		for _, obg := range operson.BGChecks {
			for _, obgt := range model.AllBGCheckTypes {
				if obg.Type&obgt != 0 {
					var nbg person.BGCheck
					switch obgt {
					case model.BGCheckDOJ:
						nperson.BGChecks.DOJ = &nbg
					case model.BGCheckFBI:
						nperson.BGChecks.FBI = &nbg
					case model.BGCheckPHS:
						nperson.BGChecks.PHS = &nbg
					default:
						panic(fmt.Sprintf("unknown background check type %d", obgt))
					}
					if obg.Date != "" {
						nbg.Cleared, _ = time.ParseInLocation("2006-01-02", obg.Date, time.Local)
					}
					nbg.Assumed = obg.Assumed
				}
			}
		}
		if reg, ok := operson.DSWRegistrations[model.DSWCERT]; ok && !reg.IsZero() {
			nperson.DSWRegistrations.CERT = &person.DSWRegistration{
				Registered: reg,
				Expiration: operson.DSWUntil[model.DSWCERT],
			}
		}
		if reg, ok := operson.DSWRegistrations[model.DSWComm]; ok && !reg.IsZero() {
			nperson.DSWRegistrations.Communications = &person.DSWRegistration{
				Registered: reg,
				Expiration: operson.DSWUntil[model.DSWComm],
			}
		}
		for _, on := range operson.Notes {
			var nn = person.Note{
				Note:       on.Note,
				Visibility: person.NoteVisibility(on.Visibility),
			}
			nn.Date, _ = time.ParseInLocation("2006-01-02", on.Date, time.Local)
			nperson.Notes = append(nperson.Notes, &nn)
		}
		for _, oec := range operson.EmContacts {
			var nec = person.EmContact{
				Name:         oec.Name,
				HomePhone:    oec.HomePhone,
				CellPhone:    oec.CellPhone,
				Relationship: oec.Relationship,
			}
			nperson.EmContacts = append(nperson.EmContacts, &nec)
		}
		sort.Slice(nperson.Notes, func(i, j int) bool {
			return nperson.Notes[i].Date.Before(nperson.Notes[j].Date)
		})
		np := person.Create(st, &nperson)
		for rid, exp := range operson.Roles {
			if exp {
				personrole.AddRole(st, np, role.WithID(st, role.ID(rid), role.FID|role.FName))
			}
		}
	}
}
