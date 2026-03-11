package icsfile

import (
	"fmt"
	"io"
	"time"

	ics "github.com/arran4/golang-ical"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util/request"
)

func Get(r *request.Request) {
	const eventFields = event.FID | event.FStart | event.FEnd | event.FName | event.FDetails | event.FFlags
	const venueFields = venue.FName
	var (
		cal   *ics.Calendar
		now   time.Time
		start string
	)
	cal = ics.NewCalendar()
	cal.SetProductId("SunnyvaleSERV.org")
	cal.SetVersion("2.0")
	now = time.Now()
	start = time.Date(now.Year(), now.Month()-6, now.Day(), 0, 0, 0, 0, time.Local).Format("2006-01-02")
	event.AllBetween(r, start, "2099-12-31", eventFields, venueFields, func(e *event.Event, v *venue.Venue) {
		if e.Flags()&event.OtherHours != 0 {
			return
		}
		ie := cal.AddEvent(fmt.Sprintf("%d@sunnyvaleserv.org", e.ID()))
		ie.SetDtStampTime(time.Now())
		ie.SetSummary(e.Name())
		start, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
		ie.SetStartAt(start)
		end, _ := time.ParseInLocation("2006-01-02T15:04", e.End(), time.Local)
		ie.SetEndAt(end)
		ie.SetURL(fmt.Sprintf("https://sunnyvaleserv.org/events/%d", e.ID()))
		if e.Details() != "" {
			ie.SetDescription(e.Details())
			// TODO render links in plain text
		}
		if v != nil {
			ie.SetLocation(v.Name())
		}
	})
	r.Header().Set("Content-Type", "text/calendar")
	r.Header().Set("Cache-Control", "public, max-age=43200") // 12 hours
	io.WriteString(r, cal.Serialize())
}
