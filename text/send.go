package text

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
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

// PostSMS handles POST /api/sms requests.
func PostSMS(r *util.Request) error {
	var (
		message    model.TextMessage
		request    *http.Request
		response   *http.Response
		deliveries []*model.TextDelivery
		recipients []string
		err        error
		params     = url.Values{}
		groups     = map[*model.Group]bool{}
	)
	if !auth.CanSendTextMessages(r) {
		return util.Forbidden
	}
	message.Sender = r.Person.ID
	if message.Message = r.FormValue("message"); message.Message == "" {
		return errors.New("missing message")
	} else if !messageFits(message.Message) {
		return errors.New("message too long")
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
					recipients = append(recipients, formatPhoneForText(p.CellPhone))
					deliveries = append(deliveries, &model.TextDelivery{Recipient: p.ID})
				} else {
					deliveries = append(deliveries, &model.TextDelivery{Recipient: p.ID, Status: "No Cell Phone"})
				}
				continue PEOPLE
			}
		}
	}
	if len(recipients) > 50 {
		return errors.New("too many numbers in one batch")
	}
	r.Tx.SaveTextMessage(&message, deliveries)
	params.Set("originator", "inbox")
	params.Set("reference", strconv.Itoa(int(message.ID)))
	params.Set("body", message.Message)
	params.Set("recipients", strings.Join(recipients, ","))
	request, _ = http.NewRequest(http.MethodPost, "https://rest.messagebird.com/messages", strings.NewReader(params.Encode()))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Authorization", "AccessKey "+config.Get("messageBirdAccessKey"))
	if response, err = http.DefaultClient.Do(request); err != nil {
		panic(err)
	}
	if response.StatusCode >= 400 {
		panic(response.Status)
	}
	response.Body.Close()
	message.Timestamp = time.Now()
	r.Tx.SaveTextMessage(&message, nil)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(r, `{"id":%d}`, message.ID)
	return nil
}

func formatPhoneForText(s string) string {
	return "1" + strings.Map(util.KeepDigits, s)
}

var doubleWidthCharacters = "\f\n^{}\\[~]|€"
var singleWidthCharacters = "£$¥èéùìòÇ\nØø\rÅåΔ_ΦΓΛΩΠΨΣΘΞ\x1bÆæßÉ !\"#¤%&'()*+,-./0123456789:;<=>?¡ABCDEFGHIJKLMNOPQRSTUVWXYZÄÖÑÜ§¿abcdefghijklmnopqrstuvwxyzäöñüà"

func messageFits(s string) bool {
	var runes, chars int
	var unicode bool
	for _, r := range s {
		runes++
		if strings.ContainsRune(singleWidthCharacters, r) {
			chars++
		} else if strings.ContainsRune(doubleWidthCharacters, r) {
			chars += 2
		} else {
			unicode = true
		}
	}
	if unicode {
		return runes <= 70
	}
	return chars <= 160
}
