package listrole

import (
	"sunnyvaleserv.org/portal/pages/admin/roleselect"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Get handles /admin/lists/$lid/roleedit/$rid requests, where either ID may be
// "NEW".
func Get(r *request.Request, lidstr, ridstr string) {
	var (
		user     *person.Person
		sender   bool
		submodel listrole.SubscriptionModel
		treedesc string
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	sender, submodel = listrole.Get(r, list.ID(util.ParseID(lidstr)), role.ID(util.ParseID(ridstr)))
	treedesc = roleselect.MakeRoleTree(r, role.FID|role.FName, nil)
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col listeditRoleForm' up-main")
	form.E("div class='formTitle formTitle-primary'>Role Privileges")
	row := form.E("div class=formRow")
	row.E("label>Role(s)")
	row.E("s-seltree name=roles class=formInput value=%s", ridstr).R(treedesc)
	row = form.E("div class=formRow")
	row.E("label>Subscription")
	box := row.E("div class=formInput")
	box.E("s-radio name=submodel value=0 label='Not allowed'", submodel == 0, "checked")
	box.E("s-radio name=submodel value=%d label='Manual'", listrole.AllowSubscription, submodel == listrole.AllowSubscription, "checked")
	box.E("s-radio name=submodel value=%d label='Automatic'", listrole.AutoSubscribe, submodel == listrole.AutoSubscribe, "checked")
	box.E("s-radio name=submodel value=%d label='Automatic, warn on unsubscribe'", listrole.WarnOnUnsubscribe, submodel == listrole.WarnOnUnsubscribe, "checked")
	box.E("div> ")
	row = form.E("div class=formRow")
	row.E("label>Sender")
	row.E("div class=formInput").E("input type=checkbox name=sender class=s-check label='Can send without moderation'", sender, "checked")
	row = form.E("div class=formButtons")
	row.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	row.E("input type=submit name=save class='sbtn sbtn-primary' value=OK")
}
