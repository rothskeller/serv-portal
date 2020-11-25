package event

import (
	"fmt"
	"strconv"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
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
	out.RawString(`{"canAdd":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader))
	out.RawString(`,"events":[`)
	for _, e := range events {
		if e.Type == model.EventHours {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(e.ID))
		out.RawString(`,"name":`)
		out.String(e.Name)
		out.RawString(`,"date":`)
		out.String(e.Date)
		out.RawString(`,"start":`)
		out.String(e.Start)
		out.RawString(`,"venue":`)
		if e.Venue != 0 {
			venue := r.Tx.FetchVenue(e.Venue)
			out.RawString(`{"id":`)
			out.Int(int(e.Venue))
			out.RawString(`,"name":`)
			out.String(venue.Name)
			out.RawString(`,"url":`)
			out.String(venue.URL)
			out.RawByte('}')
		} else {
			out.RawString(`null`)
		}
		out.RawString(`,"org":`)
		out.String(model.OrgNames[e.Org])
		out.RawString(`,"type":`)
		out.String(model.EventTypeNames[e.Type])
		out.RawString(`}`)
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
