package maillist

import (
	"fmt"
	"log"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"zombiezen.com/go/sqlite"
)

func getNamedList(dbconn *sqlite.Conn, listname string) (list *List) {
	stmt := dbconn.Prep("SELECT id FROM list WHERE name=? AND type=1")
	stmt.BindText(1, listname)
	if found, err := stmt.Step(); err != nil {
		log.Fatalf("ERROR: list lookup: %s", err)
	} else if !found {
		return nil
	}
	id := stmt.ColumnInt64(0)
	stmt.Reset()
	list = &List{
		Name:        listname,
		DisplayName: listname,
		Senders:     sets.New[string](),
		Recipients:  map[string]*RecipientData{},
		Reason:      fmt.Sprintf("via the %s@SunnyvaleSERV.org mailing list", listname),
	}
	stmt = dbconn.Prep("SELECT p.informal_name, p.email, p.email2, p.unsubscribe_token, p.flags, lp.sender, lp.sub, lp.unsub FROM list_person lp, person p WHERE lp.person=p.id AND lp.list=?")
	stmt.BindInt64(1, id)
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: list person lookup: %s", err)
		} else if !found {
			break
		}
		informalName, email, email2, unsubToken, flags := stmt.ColumnText(0), stmt.ColumnText(1), stmt.ColumnText(2), stmt.ColumnText(3), stmt.ColumnInt64(4)
		sender, sub, unsub := stmt.ColumnBool(5), stmt.ColumnBool(6), stmt.ColumnBool(7)
		email, email2 = strings.ToLower(email), strings.ToLower(email2)
		if sender {
			if email != "" {
				list.Senders.Insert(email)
			}
			if email2 != "" {
				list.Senders.Insert(email2)
			}
		}
		if !sub || unsub {
			continue
		}
		list.addRecipient(informalName, email, email2, unsubToken, flags)
	}
	stmt.Reset()
	return list
}
