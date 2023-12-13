package main

import (
	ostore "sunnyvaleserv.org/portal/ostore"
	nstore "sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
)

func convertTextMessages(tx *ostore.Tx, st *nstore.Store) {
	otextmsgs := tx.FetchTextMessages()
	for i := len(otextmsgs) - 1; i >= 0; i-- {
		otextmsg := otextmsgs[i]
		ntextmsg := textmsg.Updater{
			ID:        textmsg.ID(otextmsg.ID),
			Sender:    person.WithID(st, person.ID(otextmsg.Sender), person.FID|person.FInformalName),
			Timestamp: otextmsg.Timestamp,
			Message:   otextmsg.Message,
		}
		for _, otlid := range otextmsg.Lists {
			ntl := list.WithID(st, list.ID(otlid))
			ntextmsg.Lists = append(ntextmsg.Lists, ntl)
		}
		nt := textmsg.Create(st, &ntextmsg)
		for _, orecip := range otextmsg.Recipients {
			nrecip := person.WithID(st, person.ID(orecip.Recipient), person.FID|person.FInformalName)
			textrecip.AddRecipient(st, nt, nrecip, orecip.Number, orecip.Status, orecip.Timestamp)
			for _, oreply := range orecip.Responses {
				textrecip.AddReply(st, nt, nrecip, oreply.Response, oreply.Timestamp)
			}
		}
	}
}
