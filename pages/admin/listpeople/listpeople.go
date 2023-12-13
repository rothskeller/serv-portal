package listpeople

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles GET /admin/lists/$id/$type requests, where $type is "sub",
// "unsub", or "sender".
func Get(r *request.Request, idstr, typestr string) {
	var (
		user *person.Person
		l    *list.List
		name string
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if l = list.WithID(r, list.ID(util.ParseID(idstr))); l == nil {
		errpage.NotFound(r, user)
		return
	}
	if l.Type == list.SMS {
		name = "SMS: " + l.Name
	} else {
		name = l.Name + "@sunnyvaleserv.org"
	}
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("div class='form form-2col' up-main")
	form.E("div class='formTitle formTitle-primary'>%s", name)
	form = form.E("div class=formRow-3col")
	switch typestr {
	case "sender":
		form.E("div class=listpeopleHeading>Authorized Senders")
	case "unsub":
		form.E("div class=listpeopleHeading>Unsubscribed")
	default:
		form.E("div class=listpeopleHeading>Subscribed")
	}
	listperson.All(r, l.ID, person.FSortName|person.FEmail|person.FEmail2|person.FCellPhone|person.FFlags, func(p *person.Person, sender, sub, unsub bool) {
		if (l.Type == list.SMS && p.Flags()&person.NoText != 0) || (l.Type == list.Email && p.Flags()&person.NoEmail != 0) {
			unsub = true
		}
		switch typestr {
		case "sender":
			if !sender {
				return
			}
		case "unsub":
			if !unsub {
				return
			}
		default:
			if !sub || unsub {
				return
			}
		}
		div := form.E("div>%s", p.SortName())
		if typestr == "sender" || typestr == "unsub" {
			return
		}
		if l.Type == list.SMS && p.CellPhone() == "" {
			div.E("span class=listpeopleWarn> (no cell phone)")
		}
		if l.Type == list.Email && p.Email() == "" && p.Email2() == "" {
			div.E("span class=listpeopleWarn> (no email address)")
		}
	})
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-primary' up-dismiss>OK")
}
