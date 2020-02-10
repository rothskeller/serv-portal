package email

import (
	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetEmails handles GET /api/emails requests.
func GetEmails(r *util.Request) error {
	var out jwriter.Writer
	out.RawByte('[')
	first := true
	r.Tx.FetchEmailMessages(func(msg *model.EmailMessage) bool {
		if !msg.Attention {
			return true
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(msg.ID))
		out.RawString(`,"timestamp":`)
		out.Raw(msg.Timestamp.MarshalJSON())
		out.RawString(`,"type":`)
		out.String(model.EmailMessageTypeNames[msg.Type])
		if msg.Attention {
			out.RawString(`,"attention":true`)
		}
		if msg.Error != "" {
			out.RawString(`,"error":`)
			out.String(msg.Error)
		}
		out.RawString(`,"from":`)
		out.String(msg.From)
		out.RawString(`,"to":[`)
		for i, group := range r.Auth.FetchGroups(msg.Groups) {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(group.Name)
		}
		out.RawString(`],"subject":`)
		out.String(msg.Subject)
		out.RawByte('}')
		return true
	})
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
