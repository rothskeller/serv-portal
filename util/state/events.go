package state

import (
	"net/http"
	"time"

	"sunnyvaleserv.org/portal/util/request"
)

var (
	page  string
	month string
)

// GetEventsMonth returns the 2006-01 string describing the most recent month
// visited in the calendar (extracted from a session cookie).  If no month has
// been visited, it returns the current month.
func GetEventsMonth(r *request.Request) string {
	initialize(r)
	return month
}

// GetEventsURL returns the URL to the page the user should see when they click
// on the Events menu item.  It has whichever view they used last (calendar or
// list), and whichever year and month they viewed last.
func GetEventsURL(r *request.Request) string {
	initialize(r)
	if page == "list" {
		return "/events/list/" + month[0:4]
	}
	return "/events/calendar/" + month
}

// SetEventsPage sets the events page (list or calendar) the user prefers.
func SetEventsPage(r *request.Request, pg string) {
	page = pg
	http.SetCookie(r, &http.Cookie{
		Name:    "serv-events-page",
		Value:   page,
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

// SetEventsMonth sets what month should be visited when we next go to the
// calendar.
func SetEventsMonth(r *request.Request, mon string) {
	month = mon
	http.SetCookie(r, &http.Cookie{
		Name:  "serv-events-month",
		Value: month,
	})
}

// initialize ensures that page and month have values.  If they don't, they are
// read from request cookies, and if those cookies don't exist, default values
// are applied.
func initialize(r *request.Request) {
	if page == "" {
		page = "calendar"
		if c, err := r.Cookie("serv-events-page"); err == nil && c != nil {
			if c.Value == "list" {
				page = "list"
			}
		}
	}
	if month == "" {
		month = time.Now().Format("2006-01")
		if c, err := r.Cookie("serv-events-month"); err == nil && c != nil {
			if _, err := time.ParseInLocation("2006-01", c.Value, time.Local); err == nil {
				month = c.Value
			}
		}
	}
}
