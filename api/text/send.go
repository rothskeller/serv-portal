package text

import (
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/config"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// PostTextMessage handles POST /api/textMessage requests.
func PostTextMessage(r *util.Request) error {
	var (
		message    model.TextMessage
		request    *http.Request
		response   *http.Response
		err        error
		params     = url.Values{}
		deliveries = map[string]*model.TextDelivery{}
		groups     = map[*model.Group]bool{}
		numbers    = map[string]model.PersonID{}
	)
	if !auth.CanSendTextMessages(r) {
		return util.Forbidden
	}
	message.Sender = r.Person.ID
	if message.Message = r.FormValue("message"); message.Message == "" {
		return errors.New("missing message")
	} else if len(message.Message) > 140 {
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
	for _, p := range r.Tx.FetchPeople() {
		for group := range groups {
			if auth.IsMember(p, group) {
				if p.CellPhone != "" {
					numbers[formatPhoneForText(p.CellPhone)] = p.ID
				}
			}
		}
	}
	if len(numbers) > 50 {
		return errors.New("too many numbers in one batch")
	}
	for number, id := range numbers {
		deliveries[number] = &model.TextDelivery{Recipient: id}
		params.Add("recipients", number)
	}
	r.Tx.SaveTextMessage(&message, deliveries)
	params.Set("originator", "inbox")
	params.Set("reference", strconv.Itoa(int(message.ID)))
	params.Set("body", message.Message)
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
	return nil
}

func formatPhoneForText(s string) string {
	return "1" + strings.Map(util.KeepDigits, s)
}
