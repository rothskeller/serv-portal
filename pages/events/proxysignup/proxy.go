package proxysignup

import (
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/pages/events/signups"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /events/proxysignup/$tid requests.
func Handle(r *request.Request, tidstr string) {
	const eventFields = event.FFlags | eventview.EventFields
	const taskFields = task.FEvent | task.FOrg | signups.ShowTaskSignupsTaskFields
	var (
		user       *person.Person
		e          *event.Event
		t          *task.Task
		pid        person.ID
		p          *person.Person
		pname      string
		proxyError string
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if t = task.WithID(r, task.ID(util.ParseID(tidstr)), taskFields); t == nil {
		errpage.NotFound(r, user)
		return
	}
	e = event.WithID(r, t.Event(), eventFields)
	if !user.HasPrivLevel(t.Org(), enum.PrivLeader) || e.Flags()&event.OtherHours != 0 {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		pid, p, pname, proxyError = readProxy(r, t)
		if p != nil && p.ID() == pid {
			if r.FormValue("shift") != "" {
				signups.HandleShiftSignup(r, user, p)
			} else {
				eventview.Render(r, user, e, "")
				return
			}
		}
	}
	r.HTMLNoCache()
	if r.Method == http.MethodPost {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.eventview")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("input type=hidden name=shift")
	form.E("input type=hidden name=signedup")
	if p != nil {
		form.E("input type=hidden name=pid value=%d", p.ID())
	}
	form.E("div class='formTitle formTitle-primary'>Proxy Signup")
	row := form.E("div class=formRow")
	row.E("label for=proxy>Sign up for")
	row.E("input id=proxy name=proxy class='formInput s-search' s-filter=type:Person autofocus value=%s", pname,
		p != nil, "s-value=P%d", p.ID())
	if proxyError != "" {
		row.E("div class=formError>%s", proxyError)
	}
	row = form.E("div class=formRow")
	row.E("label>Signups")
	box := row.E("div class=formInput")
	signups.ShowTaskSignups(r, box, t, p, true, false)
	buttons := form.E("div class=formButtons")
	buttons.E("input type=submit name=ok class='sbtn sbtn-primary' value=OK")
}

func readProxy(r *request.Request, t *task.Task) (pid person.ID, p *person.Person, pname string, proxyError string) {
	var pkey string

	pid = person.ID(util.ParseID(r.FormValue("pid")))
	pname = r.FormValue("proxy")
	if pname != "" && len(r.Form["proxy"]) > 1 {
		if pkey = r.Form["proxy"][1]; strings.HasPrefix(pkey, "P") {
			p = person.WithID(r, person.ID(util.ParseID(pkey[1:])), signups.HandleShiftSignupPersonFields)
		}
	}
	if pname == "" {
		proxyError = "Search for a person by name or call sign to continue."
	} else if p == nil {
		proxyError = "No such person."
	} else {
		if t.Flags()&task.RequiresBGCheck != 0 && (!p.BGChecks().DOJ.Valid() || !p.BGChecks().FBI.Valid()) {
			proxyError = "WARNING: not background checked. "
		}
		if t.Flags()&task.CoveredByDSW != 0 {
			if reg, _ := p.DSWRegistrationForOrg(t.Org()); !reg.Valid() {
				proxyError += "WARNING: not registered as DSW."
			}
		}
	}
	return
}
