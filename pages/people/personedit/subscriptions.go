package personedit

import (
	"fmt"
	"net/http"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/server/l10n"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/recalc"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const subscriptionsPersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | person.FFlags

type listdata struct {
	warnroles []string
	list      *list.List
}

// HandleSubscriptions handles requests for /people/$id/edsubscriptions.
func HandleSubscriptions(r *request.Request, idstr string) {
	var (
		user *person.Person
		p    *person.Person
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), subscriptionsPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() != p.ID() && !user.IsWebmaster() {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		postSubscriptions(r, user, p)
	} else {
		getSubscriptions(r, p)
	}
}

func getSubscriptions(r *request.Request, p *person.Person) {
	var (
		lists      []*listdata
		subscribed = make(map[list.ID]bool)
	)
	listperson.SubscriptionsByPerson(r, p.ID(), func(l *list.List) {
		subscribed[l.ID] = true
	})
	listperson.SubscriptionRights(r, p.ID(), role.FName|role.FTitle, func(l *list.List, rl *role.Role, sm listrole.SubscriptionModel) {
		var found *listdata

		for _, ld := range lists {
			if ld.list.ID == l.ID {
				found = ld
				break
			}
		}
		if found == nil {
			clone := *l
			found = &listdata{list: &clone}
			lists = append(lists, found)
		}
		if sm == listrole.WarnOnUnsubscribe {
			if rl.Title() != "" {
				found.warnroles = append(found.warnroles, fmt.Sprintf("“%s”", rl.Title()))
			} else {
				found.warnroles = append(found.warnroles, fmt.Sprintf("“%s”", rl.Name()))
			}
		}
	})
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewSubscriptions")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Edit List Subscriptions"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	ldiv := form.E("div class='formRow-3col personeditSubscriptions'")
	for _, ld := range lists {
		var name string
		switch ld.list.Type {
		case list.Email:
			name = ld.list.Name + "@SunnyvaleSERV.org"
		case list.SMS:
			name = "SMS: " + ld.list.Name
		}
		switch len(ld.warnroles) {
		case 0:
			ldiv.E("div").E("input type=checkbox class=s-check name=list%d label=%s", ld.list.ID, name, subscribed[ld.list.ID], "checked")
		case 1:
			ldiv.E("div").E("input type=checkbox class=s-check name=list%d label=%s", ld.list.ID, name, subscribed[ld.list.ID], "checked",
				"data-warnroles=%s", r.Loc("Messages sent to %s are considered required for the %s role.  Unsubscribing from it may cause you to lose that role."),
				name, ld.warnroles[0])
		default:
			rolelist := l10n.Conjoin(ld.warnroles, "and", r.Language)
			ldiv.E("div").E("input type=checkbox class=s-check name=list%d label=%s", ld.list.ID, name, subscribed[ld.list.ID], "checked",
				"data-warnroles=%s", r.Loc("Messages sent to %s are considered required for the %s roles.  Unsubscribing from it may cause you to lose those roles."),
				name, rolelist)
		}
	}
	form.E("div id=personeditSubscriptionsWarnings class=formRow-3col")
	buttons := form.E("div class=formButtons")
	buttons.E("div class=formButtonSpace")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss").R(r.Loc("Cancel"))
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Save"))
	// This button comes last in the tree order so that it is not
	// the default.  But it comes first in the visual order because
	// of the formButton-beforeAll class.
	buttons.E("input type=submit name=unsuball class='sbtn sbtn-secondary formButton-beforeAll' value=%s", r.Loc("Unsubscribe All"))
}

func postSubscriptions(r *request.Request, user, p *person.Person) {
	r.Transaction(func() {
		if r.FormValue("unsuball") != "" {
			up := p.Updater()
			up.Flags |= person.NoEmail | person.NoText
			p.Update(r, up, person.FFlags)
			recalc.Recalculate(r)
			return
		}
		var heldemail, heldsms bool
		listperson.SubscriptionRights(r, p.ID(), 0, func(l *list.List, _ *role.Role, _ listrole.SubscriptionModel) {
			if r.FormValue(fmt.Sprintf("list%d", l.ID)) != "" {
				if l.Type == list.Email {
					heldemail = true
				} else {
					heldsms = true
				}
				listperson.Subscribe(r, l, p)
			} else {
				listperson.Unsubscribe(r, l, p)
			}
		})
		if (heldemail && p.Flags()&person.NoEmail != 0) || (heldsms && p.Flags()&person.NoText != 0) {
			up := p.Updater()

			if heldemail {
				up.Flags &^= person.NoEmail
			}
			if heldsms {
				up.Flags &^= person.NoText
			}
			p.Update(r, up, person.FFlags)
		}
		recalc.Recalculate(r)
	})
	personview.Render(r, user, p, person.ViewFull, "subscriptions")
}
