package list

import (
	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetLists handles GET /api/lists requests.
func GetLists(r *util.Request) error {
	var out jwriter.Writer

	if !r.Person.Roles[model.Webmaster] {
		return util.Forbidden
	}
	out.RawByte('[')
	for i, l := range r.Tx.FetchLists() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(l.ID))
		out.RawString(`,"type":`)
		out.String(model.ListTypeNames[l.Type])
		out.RawString(`,"name":`)
		out.String(l.Name)
		var subs, unsubs, senders int
		for _, lps := range l.People {
			if lps&model.ListSubscribed != 0 {
				subs++
			}
			if lps&model.ListUnsubscribed != 0 {
				unsubs++
			}
			if lps&model.ListSender != 0 {
				senders++
			}
		}
		out.RawString(`,"subscribed":`)
		out.Int(subs)
		out.RawString(`,"unsubscribed":`)
		out.Int(unsubs)
		out.RawString(`,"senders":`)
		out.Int(senders)
		out.RawByte('}')
	}
	out.RawByte(']')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
