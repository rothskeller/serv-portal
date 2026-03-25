package texthook

import (
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
	"sunnyvaleserv.org/portal/util/request"
)

func StatusHook(r *request.Request) {
	var (
		number  = r.FormValue("To")
		status  = r.FormValue("MessageStatus")
		message *textmsg.TextMessage
		p       *person.Person
	)
	if message = textmsg.WithNumber(r, number, textmsg.FID); message == nil {
		r.LogEntry.Problems.Add("unknown recipient phone number " + number)
		http.Error(r, "Invalid recipient phone number.", http.StatusBadRequest)
		return
	}
	if p = textrecip.WithNumber(r, message.ID(), number, person.FID|person.FInformalName); p == nil {
		r.LogEntry.Problems.Add("unmatched recipient phone number " + number)
		http.Error(r, "Invalid recipient phone number.", http.StatusBadRequest)
		return
	}
	r.Transaction(func() {
		textrecip.UpdateStatus(r, message, p, status, time.Now())
	})
	r.WriteHeader(http.StatusNoContent)
}
