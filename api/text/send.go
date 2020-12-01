package text

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sort"
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
	out.RawString(`{"lists":[`)
	first := true
	for _, l := range r.Tx.FetchLists() {
		if l.Type != model.ListSMS {
			continue
		}
		if l.People[r.Person.ID]&model.ListSender == 0 {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(l.ID))
		out.RawString(`,"name":`)
		out.String(l.Name)
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
		message    model.TextMessage
		request    *http.Request
		response   *http.Response
		tmessage   twilioMessage
		err        error
		params     = url.Values{}
		recipients = map[model.PersonID]string{}
	)
	if !r.Auth.CanA(model.PrivSendTextMessages) {
		return util.Forbidden
	}
	message.Sender = r.Person.ID
	if message.Message = r.FormValue("message"); message.Message == "" {
		return errors.New("missing message")
	}
	if len(r.Form["list"]) == 0 {
		return errors.New("no lists selected")
	}
	for _, l := range r.Form["list"] {
		if list := r.Tx.FetchList(model.ListID(util.ParseID(l))); list != nil && list.Type == model.ListSMS && list.People[r.Person.ID]&model.ListSender != 0 {
			message.Lists = append(message.Lists, list.ID)
			for pid, lps := range list.People {
				if lps&model.ListSubscribed != 0 {
					recipients[pid] = ""
				}
			}
		} else {
			return errors.New("invalid list")
		}
	}
	for pid := range recipients {
		p := r.Tx.FetchPerson(pid)
		recipients[pid] = p.SortName
		if p.NoText {
			continue
		}
		if p.CellPhone != "" {
			message.Recipients = append(message.Recipients, &model.TextRecipient{
				Recipient: pid,
				Number:    formatPhoneForText(p.CellPhone),
			})
		} else {
			message.Recipients = append(message.Recipients, &model.TextRecipient{
				Recipient: pid,
				Status:    "No Cell Phone",
			})
		}
	}
	sort.Slice(message.Recipients, func(i, j int) bool {
		return recipients[message.Recipients[i].Recipient] < recipients[message.Recipients[j].Recipient]
	})
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
