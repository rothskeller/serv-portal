package text

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
)

// GetSMSNew handles GET /api/sms/NEW requests.
func GetSMSNew(r *util.Request) error {
	var (
		out jwriter.Writer
	)
	if !r.Auth.CanA(model.PrivSendTextMessages) {
		return util.Forbidden
	}
	out.RawString(`{"groups":[`)
	first := true
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivSendTextMessages)) {
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g.ID))
		out.RawString(`,"name":`)
		out.String(g.Name)
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

type twilioMessage struct {
	DateUpdated  string `json:"date_updated"`
	ErrorMessage string `json:"error_message"`
	Status       string `json:"status"`
}

// PostSMS handles POST /api/sms requests.
func PostSMS(r *util.Request) error {
	var (
		message  model.TextMessage
		request  *http.Request
		response *http.Response
		tmessage twilioMessage
		err      error
		params   = url.Values{}
		groups   = map[*model.Group]bool{}
	)
	if !r.Auth.CanA(model.PrivSendTextMessages) {
		return util.Forbidden
	}
	message.Sender = r.Person.ID
	if message.Message = r.FormValue("message"); message.Message == "" {
		return errors.New("missing message")
	}
	for _, g := range r.Form["group"] {
		if group := r.Auth.FetchGroup(model.GroupID(util.ParseID(g))); group != nil && r.Auth.CanAG(model.PrivSendTextMessages, group.ID) {
			groups[group] = true
			message.Groups = append(message.Groups, group.ID)
		} else {
			return errors.New("invalid group")
		}
	}
	if len(groups) == 0 {
		return errors.New("no groups selected")
	}
PEOPLE:
	for _, p := range r.Tx.FetchPeople() {
		var blocked bool
		var added bool
	GROUPS:
		for group := range groups {
			if r.Auth.MemberPG(p.ID, group.ID) {
				if p.NoText {
					blocked = true
					break GROUPS
				}
				for _, nt := range group.NoText {
					if p.ID == nt {
						blocked = true
						continue GROUPS
					}
				}
				added = true
				if p.CellPhone != "" {
					message.Recipients = append(message.Recipients, &model.TextRecipient{
						Recipient: p.ID,
						Number:    formatPhoneForText(p.CellPhone),
					})
				} else {
					message.Recipients = append(message.Recipients, &model.TextRecipient{
						Recipient: p.ID,
						Status:    "No Cell Phone",
					})
				}
				continue PEOPLE
			}
		}
		if blocked && !added {
			message.Recipients = append(message.Recipients, &model.TextRecipient{
				Recipient: p.ID,
				Status:    "Texting Blocked",
			})
		}
	}
	r.Tx.CreateTextMessage(&message)
	params.Set("From", config.Get("twilioPhoneNumber"))
	params.Set("Body", message.Message)
	params.Set("StatusCallback", "https://sunnyvaleserv.org/text-status-hook")
	href := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", config.Get("twilioAccountSID"))
	for _, recip := range message.Recipients {
		if recip.Number == "" {
			continue
		}
		params.Set("To", recip.Number)
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
		recip.Status = tmessage.Status
		if tmessage.ErrorMessage != "" {
			recip.Status += ": " + tmessage.ErrorMessage
		}
		recip.Timestamp, _ = time.Parse(time.RFC1123Z, tmessage.DateUpdated)
	}
	message.Timestamp = time.Now()
	r.Tx.UpdateTextMessage(&message)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(r, `{"id":%d}`, message.ID)
	return nil
}

func formatPhoneForText(s string) string {
	return "+1" + strings.Map(util.KeepDigits, s)
}
