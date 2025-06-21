package maillist

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/store/task"
	"zombiezen.com/go/sqlite"
)

var eventListRE = regexp.MustCompile(`^(event|task)-(\d+)-(signedup|signedin|invited)$`)

func getEventList(dbconn *sqlite.Conn, listname string) (list *List) {
	var (
		match       []string
		id          int
		signupsOpen bool
		eventName   string
		taskName    string
		eventDate   string
		verb        string
		orgs        sets.Set[int64]
		sql         string
		stmt        *sqlite.Stmt
	)
	if match = eventListRE.FindStringSubmatch(listname); match == nil {
		return nil
	}
	id, _ = strconv.Atoi(match[2])
	orgs = sets.New[int64]()
	switch match[1] {
	case "event":
		stmt := dbconn.Prep("SELECT name, start FROM event WHERE id=?")
		stmt.BindInt64(1, int64(id))
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: event lookup: %s", err)
		} else if !found {
			return nil
		}
		eventName = stmt.ColumnText(0)
		eventDate = stmt.ColumnText(1)[:10]
		stmt.Reset()
		stmt = dbconn.Prep("SELECT org, flags Open FROM task WHERE event=?")
		stmt.BindInt64(1, int64(id))
		for {
			if found, err := stmt.Step(); err != nil {
				log.Fatalf("ERROR: task org lookup: %s", err)
			} else if !found {
				break
			}
			orgs.Insert(stmt.ColumnInt64(0))
			signupsOpen = signupsOpen || (task.Flag(stmt.ColumnInt64(1))&task.SignupsOpen != 0)
		}
		stmt.Reset()
	case "task":
		stmt := dbconn.Prep("SELECT name, event, org, flags FROM task WHERE id=?")
		stmt.BindInt64(1, int64(id))
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: event lookup: %s", err)
		} else if !found {
			return nil
		}
		if taskName = stmt.ColumnText(0); taskName == "Tracking" {
			taskName = ""
		}
		eventID := stmt.ColumnInt64(1)
		orgs.Insert(stmt.ColumnInt64(2))
		signupsOpen = task.Flag(stmt.ColumnInt64(3))&task.SignupsOpen != 0
		stmt.Reset()
		stmt = dbconn.Prep("SELECT name, start FROM event WHERE id=?")
		stmt.BindInt64(1, int64(eventID))
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: event lookup: %s", err)
		} else if !found {
			return nil
		}
		eventName = stmt.ColumnText(0)
		eventDate = stmt.ColumnText(1)[:10]
		stmt.Reset()
	}
	// Set up the list.
	list = &List{
		Name:          listname,
		DisplayName:   "SunnyvaleSERV Scheduler",
		Recipients:    map[string]*RecipientData{},
		NoUnsubscribe: true,
	}
	verb = "are"
	if eventDate < time.Now().Format("2006-01-02") {
		verb = "were"
	}
	switch match[3] {
	case "signedup":
		list.Reason = "because they " + verb + " signed up for "
	case "signedin":
		list.Reason = "because they signed in at "
	case "invited":
		list.Reason = "because they " + verb + " invited to "
		if signupsOpen {
			list.Reason += "sign up for "
		}
	}
	if taskName != "" {
		list.Reason += taskName + " at "
	}
	list.Reason += eventName + " on " + eventDate
	for i, org := range orgs.UnsortedList() {
		if i == 0 {
			list.Senders = getPrivLeaderEmails(dbconn, org)
		} else {
			list.Senders = list.Senders.Intersection(getPrivLeaderEmails(dbconn, org))
		}
		list.addLeaderRecipients(dbconn, org)
	}
	// Next, get the recipients.
	sql = "SELECT p.informal_name, p.email, p.email2, p.unsubscribe_token, p.flags FROM person p, "
	switch match[1] {
	case "event":
		switch match[3] {
		case "signedup":
			sql += "shift_person sp, shift s, task t WHERE sp.person=p.id AND sp.shift=s.id AND s.task=t.id AND sp.signed_up>0 AND t.event=?"
		case "signedin":
			sql += "task_person tp, task t WHERE tp.person=p.id AND tp.task=t.id AND tp.flags AND t.event=?"
		case "invited":
			sql += "task t, task_role tr, person_role pr WHERE pr.person=p.id AND pr.role=tr.role AND tr.task=t.id AND t.event=?"
		}
	case "task":
		switch match[3] {
		case "signedup":
			sql += "shift_person sp, shift s WHERE sp.person=p.id AND sp.shift=s.id AND sp.signed_up>0 AND s.task=?"
		case "signedin":
			sql += "task_person tp WHERE tp.person=p.id AND tp.task=? AND tp.flags"
		case "invited":
			sql += "task_role tr, person_role pr WHERE pr.person=p.id AND pr.role=tr.role AND tr.task=?"
		}
	}
	stmt = dbconn.Prep(sql)
	stmt.BindInt64(1, int64(id))
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: event person lookup: %s", err)
		} else if !found {
			break
		}
		informalName, email, email2, unsubToken, flags := stmt.ColumnText(0), stmt.ColumnText(1), stmt.ColumnText(2), stmt.ColumnText(3), stmt.ColumnInt64(4)
		list.addRecipient(informalName, email, email2, unsubToken, flags)
	}
	stmt.Reset()
	return list
}
