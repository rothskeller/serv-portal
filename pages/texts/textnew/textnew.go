package textnew

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"slices"
	"sort"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/texts/textview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func Handle(r *request.Request) {
	var (
		user         *person.Person
		utm          textmsg.Updater
		messageError string
		listsError   string
		haveErrors   bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if !listperson.CanSendText(r, user.ID()) {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		messageError = readMessage(r, &utm)
		listsError = readLists(r, user, &utm)
		haveErrors = messageError != "" || listsError != ""
		if !haveErrors {
			utm.Sender = user
			var tm = sendMessage(r, &utm)
			r.Header().Set("X-Up-Location", fmt.Sprintf("/texts/%d", tm.ID()))
			textview.Render(r, user, tm)
			return
		}
	}
	r.HTMLNoCache()
	if haveErrors {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.pageCanvas")
	form.E("div class='formTitle formTitle-primary'>New Text Message")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitMessage(form, &utm, messageError)
	emitLists(r, form, user, &utm, listsError)
	emitButtons(form)
}

func readMessage(r *request.Request, utm *textmsg.Updater) string {
	utm.Message = strings.TrimSpace(r.FormValue("message"))
	if utm.Message == "" {
		return "The message text is required."
	}
	return ""
}
func emitMessage(form *htmlb.Element, utm *textmsg.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=textnewMessage>Message")
	row.E("textarea id=textnewMessage name=message class=formInput rows=5>%s", utm.Message)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readLists(r *request.Request, user *person.Person, utm *textmsg.Updater) string {
	utm.Lists = utm.Lists[:0]
	list.All(r, func(l *list.List) {
		if l.Type != list.SMS || !listperson.CanSend(r, user.ID(), l.ID) {
			return
		}
		if r.FormValue(fmt.Sprintf("list%d", l.ID)) != "" {
			utm.Lists = append(utm.Lists, l.Clone())
		}
	})
	if len(utm.Lists) == 0 {
		return "At least one list must be selected."
	}
	return ""
}
func emitLists(r *request.Request, form *htmlb.Element, user *person.Person, utm *textmsg.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label>Recipients")
	box := row.E("div class=formInput")
	list.All(r, func(l *list.List) {
		if l.Type != list.SMS || !listperson.CanSend(r, user.ID(), l.ID) {
			return
		}
		box.E("div").E("input type=checkbox name=list%d class=s-check label=%s", l.ID, l.Name,
			slices.ContainsFunc(utm.Lists, func(tl *list.List) bool { return tl.ID == l.ID }), "checked")
	})
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func emitButtons(form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit class='sbtn sbtn-primary' value=Send")
}

type twilioMessage struct {
	DateUpdated  string `json:"date_updated"`
	ErrorMessage string `json:"error_message"`
	Status       string `json:"status"`
}

func sendMessage(r *request.Request, utm *textmsg.Updater) (tm *textmsg.TextMessage) {
	const personFields = person.FID | person.FInformalName | person.FSortName | person.FCellPhone | person.FFlags
	var (
		ids    []person.ID
		href   string
		rmap   = make(map[person.ID]*person.Person)
		params = make(url.Values)
	)
	utm.Timestamp = time.Now()
	for _, l := range utm.Lists {
		listperson.All(r, l.ID, personFields, func(p *person.Person, _, sub, unsub bool) {
			if sub && !unsub && p.Flags()&person.NoText == 0 {
				rmap[p.ID()] = p.Clone()
			}
		})
	}
	ids = make([]person.ID, 0, len(rmap))
	for id := range rmap {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(i, j int) bool { return rmap[ids[i]].SortName() < rmap[ids[j]].SortName() })
	r.Transaction(func() {
		tm = textmsg.Create(r, utm)
		for _, id := range ids {
			var (
				status string
				p      = rmap[id]
			)
			if p.CellPhone() == "" {
				status = "No Cell Phone"
			} else {
				status = "Not Sent Yet"
			}
			textrecip.AddRecipient(r, tm, p, p.CellPhone(), status, utm.Timestamp)
		}
	})
	params.Set("From", config.Get("twilioPhoneNumber"))
	params.Set("Body", utm.Message)
	params.Set("StatusCallback", "https://sunnyvaleserv.org/text-status-hook")
	href = fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.Get("twilioAccountSID"))
	for _, id := range ids {
		var (
			request   *http.Request
			response  *http.Response
			tmessage  twilioMessage
			status    string
			timestamp time.Time
			err       error
			p         = rmap[id]
		)
		if p.CellPhone() == "" {
			return
		}
		params.Set("To", p.CellPhone())
		request, _ = http.NewRequest(http.MethodPost, href, strings.NewReader(params.Encode()))
		request.SetBasicAuth(config.Get("twilioAccountSID"), config.Get("twilioAuthToken"))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		if response, err = http.DefaultClient.Do(request); err != nil {
			panic(err)
		}
		if response.StatusCode >= 400 {
			by, _ := httputil.DumpResponse(response, true)
			panic(string(by))
		}
		if err = json.NewDecoder(response.Body).Decode(&tmessage); err != nil {
			panic(err.Error())
		}
		response.Body.Close()
		status = tmessage.Status
		if tmessage.ErrorMessage != "" {
			status += ": " + tmessage.ErrorMessage
		}
		timestamp, _ = time.Parse(time.RFC1123Z, tmessage.DateUpdated)
		r.Transaction(func() {
			textrecip.UpdateStatus(r, tm, p, status, timestamp)
		})
	}
	return tm
}
