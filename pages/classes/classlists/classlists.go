package classlists

import (
	"fmt"
	"slices"

	"sunnyvaleserv.org/portal/maillist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /classes/$cid/lists requests.
func Handle(r *request.Request, cidstr string) {
	const classFields = class.FType
	var (
		user *person.Person
		c    *class.Class
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), classFields); c == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(c.Type().Org(), enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form' up-main")
	form.E("div class='formTitle formTitle-primary'>Class Registration Lists")
	main := form.E("div class=formRow-3col")
	// First, the registered list.
	list := maillist.GetList(r.DBConn(), fmt.Sprintf("class-%d-registered", c.ID()))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=classListsList>class-%d-registered@SunnyvaleSERV.org", c.ID())
		main.E("div>This list goes to all students whose registration in the class is confirmed.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=classListsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	// Next, the waiting list.
	list = maillist.GetList(r.DBConn(), fmt.Sprintf("class-%d-waitlist", c.ID()))
	if list != nil && len(list.Recipients) != 0 {
		main.E("div class=classListsList>class-%d-waitlist@SunnyvaleSERV.org", c.ID())
		main.E("div>This list goes to all prospective students on the waiting list for the class.")
		emails := make([]string, 0, len(list.Recipients))
		for email, recip := range list.Recipients {
			emails = append(emails, fmt.Sprintf("%s <%s>", recip.Name, email))
		}
		slices.Sort(emails)
		addrs := main.E("div class=classListsAddrs")
		for _, email := range emails {
			addrs.E("div").T(email)
		}
	}
	form.E("div class=formButtons").E("button type=button class='sbtn sbtn-primary' up-dismiss>OK")
}
