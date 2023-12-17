package personview

import (
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const subscriptionsPersonFields = person.FFlags

func showSubscriptions(r *request.Request, main *htmlb.Element, user, p *person.Person) {
	var (
		section  *htmlb.Element
		editable = user.ID() == p.ID() || user.IsWebmaster()
	)
	listperson.SubscriptionsByPerson(r, p.ID(), func(l *list.List) {
		switch l.Type {
		case list.Email:
			if p.Flags()&person.NoEmail == 0 {
				section = startSubscriptions(r, main, section, p, editable)
				section.E("div>%s@SunnyvaleSERV.org", l.Name)
			}
		case list.SMS:
			if p.Flags()&person.NoText == 0 {
				section = startSubscriptions(r, main, section, p, editable)
				section.E("div>SMS: %s", l.Name)
			}
		}
	})
	if p.Flags()&person.NoEmail != 0 {
		section = startSubscriptions(r, main, section, p, editable)
		section.E("div class=personviewSubscriptionsUnsubscribed").R(r.Loc("Unsubscribed from all email."))
	}
	if p.Flags()&person.NoText != 0 {
		section = startSubscriptions(r, main, section, p, editable)
		section.E("div class=personviewSubscriptionsUnsubscribed").R(r.Loc("Unsubscribed from all text messaging."))
	}
	if section == nil {
		if editable {
			section = startSubscriptions(r, main, section, p, editable)
			section.E("div").R(r.Loc("Not subscribed to any email or text messaging."))
		}
	}
}

func startSubscriptions(r *request.Request, main *htmlb.Element, section *htmlb.Element, p *person.Person, editable bool) *htmlb.Element {
	if section == nil {
		section = main.E("div class=personviewSection")
		sheader := section.E("div class=personviewSectionHeader")
		sheader.E("div class=personviewSectionHeaderText").R(r.Loc("Subscriptions"))
		if editable {
			sheader.E("div class=personviewSectionHeaderEdit").
				E("a href=/people/%d/edsubscriptions up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'", p.ID()).R(r.Loc("Edit"))
		}
		section = section.E("div class=personviewSubscriptions")
	}
	return section
}
