package text

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetSMS handles GET /api/sms requests.
func GetSMS(r *util.Request) error {
	var (
		out jwriter.Writer
	)
	if !r.Auth.CanA(model.PrivSendTextMessages) {
		return util.Forbidden
	}
	out.RawString(`{"messages":[`)
	for i, m := range r.Tx.FetchTextMessages() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(m.ID))
		out.RawString(`,"timestamp":`)
		out.String(m.Timestamp.In(time.Local).Format("2006-01-02 15:04"))
		out.RawString(`,"sender":`)
		out.String(r.Tx.FetchPerson(m.Sender).InformalName)
		out.RawString(`,"groups":[`)
		for i, g := range m.Groups {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(r.Auth.FetchGroup(g).Name)
		}
		out.RawString(`],"message":`)
		out.String(m.Message)
		out.RawByte('}')
	}
	out.RawString(`],"groups":[`)
	first := true
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivSendTextMessages)) {
		if r.Auth.CanAG(model.PrivSendTextMessages, g.ID) {
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