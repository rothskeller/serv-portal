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

	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/config"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetSMSNew handles GET /api/sms/NEW requests.
func GetSMSNew(r *util.Request) error {
	var (
		out jwriter.Writer
	)
	if !auth.CanSendTextMessages(r) {
		return util.Forbidden
	}
	out.RawString(`{"groups":[`)
	first := true
	for _, g := range r.Tx.FetchGroups() {
		if g.AllowTextMessages {
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
	if !auth.CanSendTextMessages(r) {
		return util.Forbidden
	}
	message.Sender = r.Person.ID
	if message.Message = r.FormValue("message"); message.Message == "" {
		return errors.New("missing message")
	}
	for _, g := range r.Form["group"] {
		if group := r.Tx.FetchGroup(model.GroupID(util.ParseID(g))); group != nil && group.AllowTextMessages {
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
		for group := range groups {
			if auth.IsMember(p, group) {
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
	}
	r.Tx.SaveTextMessage(&message)
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
	r.Tx.SaveTextMessage(&message)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(r, `{"id":%d}`, message.ID)
	return nil
}

func formatPhoneForText(s string) string {
	return "+1" + strings.Map(util.KeepDigits, s)
}
