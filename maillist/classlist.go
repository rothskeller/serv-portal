package maillist

import (
	"fmt"
	"log"
	"regexp"
	"strconv"

	"sunnyvaleserv.org/portal/store/class"
	"zombiezen.com/go/sqlite"
)

var classListRE = regexp.MustCompile(`^class-(\d+)-(registered|waitlist)$`)

func getClassList(dbconn *sqlite.Conn, listname string) (list *List) {
	match := classListRE.FindStringSubmatch(listname)
	if match == nil {
		return nil
	}
	// First get the class type and the enrollment limit.  This also
	// verifies that the class exists.
	id, _ := strconv.Atoi(match[1])
	stmt := dbconn.Prep("SELECT type, start, elimit FROM class WHERE id=?")
	stmt.BindInt64(1, int64(id))
	if found, err := stmt.Step(); err != nil {
		log.Fatalf("ERROR: class type lookup: %s", err)
	} else if !found {
		return nil
	}
	ctype := class.Type(stmt.ColumnInt64(0))
	start := stmt.ColumnText(1)
	elimit := stmt.ColumnInt(2)
	if stmt.ColumnType(1) == sqlite.TypeNull {
		elimit = 1000
	}
	stmt.Reset()
	// Set up the list.
	list = &List{
		Name:          listname,
		DisplayName:   "SunnyvaleSERV Registrar",
		Senders:       getPrivLeaderEmails(dbconn, int64(ctype.Org())),
		Recipients:    map[string]*RecipientData{},
		NoUnsubscribe: true,
	}
	if match[2] == "registered" {
		list.Reason = "because they are registered for the "
	} else {
		list.Reason = "because they are on the waiting list for the "
	}
	list.Reason += fmt.Sprintf("%s %s class", start, ctype)
	list.addLeaderRecipients(dbconn, int64(ctype.Org()))
	// Next, get the recipients.
	stmt = dbconn.Prep("SELECT first_name, last_name, email FROM classreg WHERE class=? ORDER BY id")
	stmt.BindInt64(1, int64(id))
	var count int
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: class registration lookup: %s", err)
		} else if !found {
			break
		}
		count++
		firstName := stmt.ColumnText(0)
		lastName := stmt.ColumnText(1)
		email := stmt.ColumnText(2)
		if email == "" {
			continue
		}
		if (match[2] == "registered" && count <= elimit) || (match[2] != "registered" && count > elimit) {
			list.addRecipient(firstName+" "+lastName, email, "", "", 0)
		}
	}
	stmt.Reset()
	return list
}
