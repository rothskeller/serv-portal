package maillist

import (
	"log"
	"net/mail"
	"strings"

	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/config"

	"k8s.io/apimachinery/pkg/util/sets"
	"zombiezen.com/go/sqlite"
)

// A List gives the characteristics and membership of a mailing list.
type List struct {
	// Name is the name of the mailing list, as seen in its email address.
	// (LHS only, no "@" or domain name.)
	Name string
	// DisplayName is the name of the mailing list as it appears after
	// "via" in the From line of messages sent out by the list server.
	// For defined email lists, it's usually the same as Name; for dynamic
	// lists, it's usually a canned string like "SERV Scheduler".
	DisplayName string
	// Senders is the set of email addresses that are allowed to send to
	// the list without moderation.
	Senders sets.Set[string]
	// Recipients is a map keyed by the recipient email addresses of the
	// list.
	Recipients map[string]*RecipientData
	// Bcc is a map keyed by email address of the people who should be
	// bcc'd on the list emails.  Usually this is the leaders of the
	// organization(s) relevant to the list.
	Bcc map[string]*RecipientData
	// Reason is a description of who's on the mailing list, i.e., why a
	// person received a mail addressed to it.  It should complete the
	// sentence "This email was sent to «name» <»email»> ..."
	Reason string
	// NoUnsubscribe is a flag indicating that people cannot unsubscribe
	// from the list.  Instructions of how to unsubscribe are omitted from
	// the footer.
	NoUnsubscribe bool
}

// RecipientData contains information about a recipient email address.
type RecipientData struct {
	// Name is the full name of the person with this email address.
	Name string
	// UnsubscribeToken is the unsubscribe token for the person, if any.
	UnsubscribeToken string
}

// GetList gets the list definition and membership for the named list.  It
// returns nil if there is no such list.
func GetList(dbconn *sqlite.Conn, listname string) (list *List) {
	if list = getNamedList(dbconn, listname); list != nil {
		return list
	}
	if list = getEventList(dbconn, listname); list != nil {
		return list
	}
	if list = getRoleList(dbconn, listname); list != nil {
		return list
	}
	if list = getClassList(dbconn, listname); list != nil {
		return list
	}
	return nil
}

func (list *List) addRecipient(informalName, email, email2, unsubToken string, flags int64) {
	list.addRecipient2(informalName, email, email2, unsubToken, flags, list.Recipients)
}

func (list *List) addRecipient2(informalName, email, email2, unsubToken string, flags int64, set map[string]*RecipientData) {
	if person.Flags(flags)&person.NoEmail != 0 {
		return
	}
	email, email2 = strings.ToLower(email), strings.ToLower(email2)
	if email != "" && list.Recipients[email] == nil && set[email] == nil {
		set[email] = &RecipientData{
			Name:             informalName,
			UnsubscribeToken: unsubToken,
		}
		if list.Recipients[email] != nil && list.Bcc[email] != nil {
			delete(list.Bcc, email)
		}
	}
	if email2 != "" && list.Recipients[email2] == nil && set[email] == nil {
		set[email2] = &RecipientData{
			Name:             informalName,
			UnsubscribeToken: unsubToken,
		}
		if list.Recipients[email2] != nil && list.Bcc[email2] != nil {
			delete(list.Bcc, email2)
		}
	}
}

func getPrivLeaderEmails(dbconn *sqlite.Conn, org int64) (senders sets.Set[string]) {
	senders = sets.New[string]()
	stmt := dbconn.Prep("SELECT p.email, p.email2 FROM person p, person_privlevel pp WHERE pp.person=p.id AND pp.privlevel>=3 AND pp.org=?")
	stmt.BindInt64(1, org)
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: org leader lookup: %s", err)
		} else if !found {
			break
		}
		email := strings.ToLower(stmt.ColumnText(0))
		email2 := strings.ToLower(stmt.ColumnText(1))
		senders.Insert(email, email2)
	}
	stmt.Reset()
	senders.Delete("")
	return senders
}

func (list *List) addLeaderRecipients(dbconn *sqlite.Conn, org int64) {
	if list.Bcc == nil {
		list.Bcc = make(map[string]*RecipientData)
	}
	stmt := dbconn.Prep("SELECT p.informal_name, p.email, p.email2, p.unsubscribe_token, p.flags FROM person p, person_role pr, role r WHERE pr.person=p.id AND pr.role=r.id AND pr.explicit AND r.privlevel>=3 AND r.org=?")
	stmt.BindInt64(1, org)
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: org leader lookup: %s", err)
		} else if !found {
			break
		}
		informalName, email, email2, unsubToken, flags := stmt.ColumnText(0), stmt.ColumnText(1), stmt.ColumnText(2), stmt.ColumnText(3), stmt.ColumnInt64(4)
		list.addRecipient2(informalName, email, email2, unsubToken, flags, list.Bcc)
	}
	stmt.Reset()
	// Also add the wiretap addresses if any.
	addrs, _ := mail.ParseAddressList(config.Get("listWiretap"))
	for _, addr := range addrs {
		list.addRecipient2(addr.Name, addr.Address, "", "", 0, list.Bcc)
	}
}
