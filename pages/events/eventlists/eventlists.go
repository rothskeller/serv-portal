package eventlists

import (
	"fmt"
	"slices"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /events/eventlists/$eid requests.
func Handle(r *request.Request, eidstr string) {
	const taskFields = task.FOrg
	var (
		user      *person.Person
		eventID   event.ID
		exists    bool
		forbidden bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	eventID = event.ID(util.ParseID(eidstr))
	task.AllForEvent(r, eventID, taskFields, func(t *task.Task) {
		exists = true
		if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
			forbidden = true
		}
	})
	if !exists {
		errpage.NotFound(r, user)
		return
	}
	if forbidden {
		errpage.Forbidden(r, user)
		return
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form' up-main")
	form.E("div class='formTitle formTitle-primary'>Event Email Lists")
	main := form.E("div class=formRow-3col")
	// First, the signed-in list.
	list := maillist.GetList(r.DBConn(), fmt.Sprintf("event-%d-signedin", eventID))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventListsList>task-%d-signedin@SunnyvaleSERV.org", eventID)
		main.E("div>This list goes to all volunteers who are recorded as having signed in for, and/or credited with participation in, any of the event tasks.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventListsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	// Next, the signed-up list.
	list = maillist.GetList(r.DBConn(), fmt.Sprintf("event-%d-signedup", eventID))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventListsList>task-%d-signedup@SunnyvaleSERV.org", eventID)
		main.E("div>This list goes to all volunteers signed up for the event (any task, any shift).")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventListsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	// Finally, the invited list.
	list = maillist.GetList(r.DBConn(), fmt.Sprintf("event-%d-invited", eventID))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventListsList>task-%d-invited@SunnyvaleSERV.org", eventID)
		main.E("div>This list goes to all volunteers in the role(s) invited to any of the event tasks.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventListsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	form.E("div class=formButtons").E("button type=button class='sbtn sbtn-primary' up-dismiss>OK")
}
