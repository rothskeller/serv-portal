package texthook

import (
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
	"sunnyvaleserv.org/portal/util/request"
)

func ReceivedHook(r *request.Request) {
	var (
		number  = r.FormValue("From")
		body    = r.FormValue("Body")
		message *textmsg.TextMessage
		p       *person.Person
	)
	if message = textmsg.WithNumber(r, number, textmsg.FID); message == nil {
		r.LogEntry.Problems.Add("incoming message from unknown phone number: " + number)
		r.WriteHeader(http.StatusNoContent)
		return
	}
	if p = textrecip.WithNumber(r, message.ID(), number, person.FID|person.FInformalName); p == nil {
		r.LogEntry.Problems.Add("no recipient with phone number: " + number)
		r.WriteHeader(http.StatusNoContent)
		return
	}
	r.Transaction(func() {
		textrecip.AddReply(r, message, p, body, time.Now())
	})
	r.WriteHeader(http.StatusNoContent)
}
