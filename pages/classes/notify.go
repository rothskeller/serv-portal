package classes

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const notifyPersonFields = person.FID | person.FInformalName | person.FSortName | person.FEmail | person.FEmail2 | person.FCellPhone | person.FCallSign | person.FFlags

// HandleNotify handles /cert-basic/notify and /pep/notify requests.
func HandleNotify(r *request.Request, cname string) {
	var (
		user  *person.Person
		ld    *list.List
		sub   bool
		unsub bool
	)
	// Get the user information.
	if user = auth.SessionUser(r, notifyPersonFields, false); user == nil {
		if user = handleRegisterNotLoggedIn(r, true); user == nil {
			return
		}
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	// Get the list.
	if ld = list.WithName(r, cnameToListName[cname]); ld == nil {
		println("no such list", cnameToListName[cname])
		errpage.NotFound(r, user)
		return
	}
	// Check whether this person is already subscribed.
	if sub, unsub = listperson.Subscribed(r, user, ld); unsub || user.Flags()&person.NoEmail != 0 {
		sub = false
	}
	// If they're not already subscribed, it's an unconditional subscribe.
	if !sub {
		r.Transaction(func() {
			listperson.Subscribe(r, ld, user)
			if user.Flags()&person.NoEmail != 0 {
				uu := user.Updater()
				uu.Flags &^= person.NoEmail
				user.Update(r, uu, person.FFlags)
			}
		})
		r.HTMLNoCache()
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form method=POST class='form form-2col' up-main up-layer=parent up-target=main")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Notifications"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class=formRow-3col").TF(r.Loc("You are now subscribed to the %s@SunnyvaleSERV.org notification list."), ld.Name)
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-primary' up-dismiss>%s", r.Loc("OK"))
		return
	}
	// If they submitted the unsubscribe form, save the change.
	if r.FormValue("unsubscribe") != "" {
		r.Transaction(func() {
			listperson.Unsubscribe(r, ld, user)
		})
		r.HTMLNoCache()
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form method=POST class='form form-2col' up-main up-layer=parent up-target=main")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Notifications"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class=formRow-3col").TF(r.Loc("You have been removed from the %s@SunnyvaleSERV.org notification list."), ld.Name)
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-primary' up-dismiss>%s", r.Loc("OK"))
		return
	}
	// Give them an unsubscribe form.
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form method=POST class='form form-2col' up-main up-target=main")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Notifications"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("div class=formRow-3col").TF(r.Loc("You are already subscribed to the %s@SunnyvaleSERV.org notification list.  To remove yourself from the list, click “Unsubscribe”."), ld.Name)
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
	buttons.E("input type=submit name=unsubscribe class='sbtn sbtn-primary' value=%s", r.Loc("Unsubscribe"))
}

var cnameToListName = map[string]string{
	"cert-basic": "cert-notify",
	"pep":        "pep-notify",
}
