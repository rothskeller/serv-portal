package tasklists

import (
	"fmt"
	"slices"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /events/tasklists/$tid requests.
func Handle(r *request.Request, tidstr string) {
	const taskFields = task.FOrg
	var (
		user *person.Person
		t    *task.Task
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if t = task.WithID(r, task.ID(util.ParseID(tidstr)), taskFields); t == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form' up-main")
	form.E("div class='formTitle formTitle-primary'>Task Email Lists")
	main := form.E("div class=formRow-3col")
	// First, the signed-in list.
	list := maillist.GetList(r.DBConn(), fmt.Sprintf("task-%d-signedin", t.ID()))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventTasklistsList>task-%d-signedin@SunnyvaleSERV.org", t.ID())
		main.E("div>This list goes to all volunteers who are recorded as having signed in for the task and/or credited with participation in it.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventTasklistsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	// Next, the signed-up list.
	list = maillist.GetList(r.DBConn(), fmt.Sprintf("task-%d-signedup", t.ID()))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventTasklistsList>task-%d-signedup@SunnyvaleSERV.org", t.ID())
		main.E("div>This list goes to all volunteers signed up for the task (any shift).")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventTasklistsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	// Finally, the invited list.
	list = maillist.GetList(r.DBConn(), fmt.Sprintf("task-%d-invited", t.ID()))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=eventTasklistsList>task-%d-invited@SunnyvaleSERV.org", t.ID())
		main.E("div>This list goes to all volunteers in the role(s) invited to the task.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=eventTasklistsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	form.E("div class=formButtons").E("button type=button class='sbtn sbtn-primary' up-dismiss>OK")
}
