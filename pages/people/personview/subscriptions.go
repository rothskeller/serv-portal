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
		section.E("div class=personviewSubscriptionsUnsubscribed").R(r.LangString("Unsubscribed from all email.", "Se ha dado de baja de todos los correos electrónicos."))
	}
	if p.Flags()&person.NoText != 0 {
		section = startSubscriptions(r, main, section, p, editable)
		section.E("div class=personviewSubscriptionsUnsubscribed").R(r.LangString("Unsubscribed from all text messaging.", "Se ha dado de baja de todos los mensajes de texto."))
	}
	if section == nil {
		if editable {
			section = startSubscriptions(r, main, section, p, editable)
			section.E("div").R(r.LangString("Not subscribed to any email or text messaging.", "No suscrito a ningún correo electrónico o mensaje de texto."))
		}
	}
}

func startSubscriptions(r *request.Request, main *htmlb.Element, section *htmlb.Element, p *person.Person, editable bool) *htmlb.Element {
	if section == nil {
		section = main.E("div class=personviewSection")
		sheader := section.E("div class=personviewSectionHeader")
		sheader.E("div class=personviewSectionHeaderText").R(r.LangString("Subscriptions", "Suscripciones"))
		if editable {
			sheader.E("div class=personviewSectionHeaderEdit").
				E("a href=/people/%d/edsubscriptions up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'", p.ID()).R(r.LangString("Edit", "Editar"))
		}
		section = section.E("div class=personviewSubscriptions")
	}
	return section
}
