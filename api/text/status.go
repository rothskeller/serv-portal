package text

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetSMS1 handles GET /api/sms/$id requests.
func GetSMS1(r *util.Request, idstr string) (err error) {
	var (
		message *model.TextMessage
		visible bool
		out     jwriter.Writer
	)
	if !r.Auth.CanA(model.PrivSendTextMessages) {
		return util.Forbidden
	}
	if message = r.Tx.FetchTextMessage(model.TextMessageID(util.ParseID(idstr))); message == nil {
		return util.NotFound
	}
	for _, gid := range message.Groups {
		if r.Auth.CanAG(model.PrivSendTextMessages, gid) {
			visible = true
			break
		}
	}
	if !visible {
		return util.Forbidden
	}
	out.RawString(`{"id":`)
	out.Int(int(message.ID))
	out.RawString(`,"sender":`)
	out.String(r.Tx.FetchPerson(message.Sender).InformalName)
	out.RawString(`,"groups":[`)
	for i, gid := range message.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(r.Auth.FetchGroup(gid).Name)
	}
	out.RawString(`],"timestamp":`)
	out.String(message.Timestamp.In(time.Local).Format("2006-01-02 15:04:05"))
	out.RawString(`,"message":`)
	out.String(message.Message)
	out.RawString(`,"deliveries":[`)
	for i, d := range message.Recipients {
		if i != 0 {
			out.RawByte(',')
		}
		p := r.Tx.FetchPerson(d.Recipient)
		out.RawString(`{"id":`)
		out.Int(int(d.Recipient))
		out.RawString(`,"recipient":`)
		out.String(p.SortName)
		out.RawString(`,"number":`)
		out.String(d.Number)
		out.RawString(`,"status":`)
		out.String(d.Status)
		out.RawString(`,"timestamp":`)
		out.String(d.Timestamp.In(time.Local).Format("2006-01-02 15:04:05"))
		out.RawString(`,"responses":[`)
		for i, r := range d.Responses {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"response":`)
			out.String(r.Response)
			out.RawString(`,"timestamp":`)
			out.String(r.Timestamp.In(time.Local).Format("2006-01-02 15:04:05"))
			out.RawByte('}')
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
