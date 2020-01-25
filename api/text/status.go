package text

import (
	"time"

	"github.com/mailru/easyjson/jwriter"
	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetTextMessage handles GET /api/textMessage/$id requests.
func GetTextMessage(r *util.Request, idstr string) (err error) {
	var (
		message    *model.TextMessage
		deliveries []*model.TextDelivery
		out        jwriter.Writer
	)
	if !auth.CanSendTextMessages(r) {
		return util.Forbidden
	}
	if message = r.Tx.FetchTextMessage(model.TextMessageID(util.ParseID(idstr))); message == nil {
		return util.NotFound
	}
	deliveries = r.Tx.FetchTextDeliveries(message.ID)
	out.RawString(`{"id":`)
	out.Int(int(message.ID))
	out.RawString(`,"sender":`)
	out.String(r.Tx.FetchPerson(message.Sender).InformalName)
	out.RawString(`,"groups":[`)
	for i, gid := range message.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.String(r.Tx.FetchGroup(gid).Name)
	}
	out.RawString(`],"timestamp":`)
	out.String(message.Timestamp.In(time.Local).Format("2006-01-02 15:04:05"))
	out.RawString(`,"message":`)
	out.String(message.Message)
	out.RawString(`,"deliveries":[`)
	for i, d := range deliveries {
		if i != 0 {
			out.RawByte(',')
		}
		p := r.Tx.FetchPerson(d.Recipient)
		out.RawString(`{"recipient":`)
		out.String(p.InformalName)
		out.RawString(`,"number":`)
		out.String(p.CellPhone)
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
