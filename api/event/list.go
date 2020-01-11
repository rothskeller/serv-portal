package event

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetEvents handles GET /events requests.
func GetEvents(r *util.Request) error {
	var (
		year   int
		events []*model.Event
		out    jwriter.Writer
		first  = true
	)
	if year, _ = strconv.Atoi(r.FormValue("year")); year < 2000 || year > 2099 {
		year = time.Now().Year()
	}
	events = r.Tx.FetchEvents(fmt.Sprintf("%d-01-01", year), fmt.Sprintf("%d-12-31", year))
	r.Tx.Commit()
	out.RawString(`{"canAdd":`)
	out.Bool(auth.CanCreateEvents(r))
	out.RawString(`,"events":[`)
	for _, e := range events {
		if !auth.CanViewEvent(r, e) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(e.ID))
		out.RawString(`,"date":`)
		out.String(e.Date)
		out.RawString(`,"name":`)
		out.String(e.Name)
		out.RawString(`,"servGroup":`)
		out.String(string(servGroupForEvent(e)))
		out.RawString(`,"roles":[`)
		for i, r := range e.Roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(r.Name)
		}
		out.RawString(`]}`)
	}
	out.RawString(`]}`)
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

func servGroupForEvent(e *model.Event) (group model.SERVGroup) {
	for _, role := range e.Roles {
		if group == "" {
			group = role.SERVGroup
		} else if group != role.SERVGroup {
			group = model.GroupSERV
		}
	}
	return group
}
