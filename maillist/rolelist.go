package maillist

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"zombiezen.com/go/sqlite"
)

var roleListRE = regexp.MustCompile(`^role-(\d+)([-+].*)?$`)

func getRoleList(dbconn *sqlite.Conn, listname string) (list *List) {
	match := roleListRE.FindStringSubmatch(listname)
	if match == nil {
		return nil
	}
	// Parse the add-ons and build the SQL query now, so that we detect any
	// errors in it before we hit the database.
	sql := "SELECT p.informal_name, p.email, p.email2, p.unsubscribe_token, p.flags FROM person p, person_role pr"
	where := " WHERE pr.person=p.id AND pr.role=?"
	criteria := match[2]
	var critLabels []string
	for criteria != "" {
		var crit string
		if idx := strings.IndexAny(criteria[1:], "-+"); idx < 0 {
			crit, criteria = criteria, ""
		} else {
			crit, criteria = criteria[:idx+1], criteria[idx+1:]
		}
		not, crit := crit[0] == '-', crit[1:]
		switch strings.ToLower(crit) {
		case "bgcheck":
			if not {
				sql += " LEFT JOIN person_bgcheck bg0 ON bg0.person=p.id AND bg0.type=0 LEFT JOIN person_bgcheck bg1 ON bg1.person=p.id AND bg1.type=1"
				where += " AND (bg0.person IS NULL OR bg0.nli<='" + time.Now().Format("2006-01-02") + "')"
				where += " AND (bg1.person IS NULL OR bg1.nli<='" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "without background check")
			} else {
				sql += ", person_bgcheck bg0, person_bgcheck bg1"
				where += " AND bg0.person=p.id AND bg0.type=0 AND (bg0.nli IS NULL OR bg0.nli>'" + time.Now().Format("2006-01-02") + "')"
				where += " AND bg1.person=p.id AND bg1.type=1 AND (bg1.nli IS NULL OR bg1.nli>'" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "with background check")
			}
		case "cardkey":
			if not {
				where += " AND NOT p.identification&2"
				critLabels = append(critLabels, "without card key")
			} else {
				where += " AND p.identification&2"
				critLabels = append(critLabels, "with card key")
			}
		case "dswcert":
			if not {
				sql += " LEFT JOIN person_dswreg pd3 ON pd3.person=p.id AND pd3.class=3"
				where += " AND (pd3.person IS NULL OR pd3.expiration<='" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "not registered as DSW for CERT")
			} else {
				sql += ", person_dswreg pd3"
				where += " AND pd3.person=p.id AND pd3.class=3 AND (pd3.expiration IS NULL OR pd3.expiration>'" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "registered as DSW for CERT")
			}
		case "dswcomm":
			if not {
				sql += " LEFT JOIN person_dswreg pd2 ON pd2.person=p.id AND pd2.class=2"
				where += " AND (pd2.person IS NULL OR pd2.expiration<='" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "not registered as DSW for Communications")
			} else {
				sql += ", person_dswreg pd2"
				where += " AND pd2.person=p.id AND pd2.class=2 AND (pd2.expiration IS NULL OR pd2.expiration>'" + time.Now().Format("2006-01-02") + "')"
				critLabels = append(critLabels, "registered as DSW for Communications")
			}
		case "photoid":
			if not {
				where += " AND NOT p.identification&1"
				critLabels = append(critLabels, "without photo ID")
			} else {
				where += " AND p.identification&1"
				critLabels = append(critLabels, "with photo ID")
			}
		case "volreg":
			if not {
				where += " AND p.volgistics_id IS NULL"
				critLabels = append(critLabels, "not registered as a volunteer")
			} else {
				where += " AND p.volgistics_id"
				critLabels = append(critLabels, "registered as a volunteer")
			}
		default:
			return nil
		}
		// The result still may not be valid SQL, if they added
		// duplicate criteria.  We'll catch that below when we prepare
		// the statement.
	}
	sql += where
	roleID, _ := strconv.Atoi(match[1])
	// First get the role organization.  This also verifies that the role
	// exists.
	stmt := dbconn.Prep("SELECT name, title, org FROM role WHERE id=?")
	stmt.BindInt64(1, int64(roleID))
	if found, err := stmt.Step(); err != nil {
		log.Fatalf("ERROR: role org lookup: %s", err)
	} else if !found {
		return nil
	}
	title := stmt.ColumnText(1)
	if title == "" {
		title = stmt.ColumnText(0)
	}
	org := stmt.ColumnInt64(2)
	stmt.Reset()
	// Set up the list.
	list = &List{
		Name:          listname,
		DisplayName:   "SunnyvaleSERV",
		Senders:       getPrivLeaderEmails(dbconn, org),
		Recipients:    map[string]*RecipientData{},
		Reason:        "because they are a " + title,
		NoUnsubscribe: true,
	}
	switch len(critLabels) {
	case 0:
		break
	case 1:
		list.Reason += " " + critLabels[0]
	default:
		list.Reason += " " + strings.Join(critLabels, ", ")
	}
	list.addLeaderRecipients(dbconn, org)
	// Next, get the recipients.
	var err error
	if stmt, err = dbconn.Prepare(sql); err != nil {
		// most likely a duplicate field that caused a duplicate join
		return nil
	}
	stmt.BindInt64(1, int64(roleID))
	for {
		if found, err := stmt.Step(); err != nil {
			log.Fatalf("ERROR: role person lookup: %s", err)
		} else if !found {
			break
		}
		informalName, email, email2, unsubToken, flags := stmt.ColumnText(0), stmt.ColumnText(1), stmt.ColumnText(2), stmt.ColumnText(3), stmt.ColumnInt64(4)
		list.addRecipient(informalName, email, email2, unsubToken, flags)
	}
	stmt.Reset()
	return list
}
