package event

import (
	"errors"
	"regexp"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"github.com/microcosm-cc/bluemonday"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

var htmlSanitizer = bluemonday.NewPolicy().
	RequireParseableURLs(true).
	AllowURLSchemes("http", "https").
	RequireNoFollowOnLinks(true).
	AllowAttrs("href").OnElements("a").
	AddTargetBlankToFullyQualifiedLinks(true)

// GetEvent handles GET /api/events/$id requests (where $id may be "NEW").
func GetEvent(r *util.Request, idstr string) error {
	var (
		event         *model.Event
		canEdit       bool
		canAttendance bool
		out           jwriter.Writer
	)
	if idstr == "NEW" {
		if !auth.CanCreateEvents(r) {
			return util.Forbidden
		}
		canEdit = true
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		if !auth.CanViewEvent(r, event) {
			return util.Forbidden
		}
		canEdit = auth.CanManageEvent(r, event)
		canAttendance = auth.CanRecordAttendanceAtEvent(r, event)
	}
	out.RawString(`{"canEdit":`)
	out.Bool(canEdit)
	out.RawString(`,"canAttendance":`)
	out.Bool(canAttendance)
	if event != nil {
		out.RawString(`,"event":{"id":`)
		out.Int(int(event.ID))
		out.RawString(`,"name":`)
		out.String(event.Name)
		out.RawString(`,"date":`)
		out.String(event.Date)
		out.RawString(`,"start":`)
		out.String(event.Start)
		out.RawString(`,"end":`)
		out.String(event.End)
		out.RawString(`,"venue":`)
		if event.Venue != nil {
			out.RawString(`{"id":`)
			out.Int(int(event.Venue.ID))
			out.RawString(`,"name":`)
			out.String(event.Venue.Name)
			out.RawString(`,"address":`)
			out.String(event.Venue.Address)
			out.RawString(`,"city":`)
			out.String(event.Venue.City)
			out.RawString(`,"url":`)
			out.String(event.Venue.URL)
			out.RawByte('}')
		} else {
			out.RawString(`{"id":0,"name":"","address":"","city":"","url":""}`)
		}
		out.RawString(`,"details":`)
		out.String(event.Details)
		out.RawString(`,"type":`)
		out.String(string(event.Type))
		out.RawString(`,"roles":[`)
		for i, r := range event.Roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.Int(int(r.ID))
		}
		out.RawString(`]}`)
	} else {
		out.RawString(`,"event":{"id":0,"name":"","date":"","start":"","end":"","venue":null,"details":"","type":"","roles":[]}`)
	}
	if canEdit {
		out.RawString(`,"roles":[`)
		for i, t := range auth.RolesCanManageEvents(r) {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(t.ID))
			out.RawString(`,"name":`)
			out.String(t.Name)
			out.RawByte('}')
		}
		out.RawString(`],"venues":[`)
		for i, v := range r.Tx.FetchVenues() {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(v.ID))
			out.RawString(`,"name":`)
			out.String(v.Name)
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	if canAttendance {
		var (
			attended = r.Tx.FetchAttendanceByEvent(event)
			first    = true
		)
		out.RawString(`,"people":[`)
		for _, p := range r.Tx.FetchPeople() {
			if !auth.CanViewEventP(r, p, event) {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(p.ID))
			out.RawString(`,"lastName":`)
			out.String(p.LastName)
			out.RawString(`,"firstName":`)
			out.String(p.FirstName)
			out.RawString(`,"attended":`)
			out.Bool(attended[p.ID])
			out.RawByte('}')
		}
		out.RawByte(']')
	}
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

var dateRE = regexp.MustCompile(`^20\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
var timeRE = regexp.MustCompile(`^(?:[01][0-9]|2[0-3]):[0-5][0-9]$`)
var yearRE = regexp.MustCompile(`^20\d\d$`)

// PostEvent handles POST /events/$id requests (where $id may be "NEW").
func PostEvent(r *util.Request, idstr string) error {
	var event *model.Event

	if idstr == "NEW" {
		if !auth.CanCreateEvents(r) {
			return util.Forbidden
		}
		event = new(model.Event)
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		if !auth.CanManageEvent(r, event) {
			return util.Forbidden
		}
	}
	if r.FormValue("delete") != "" && event.ID != 0 {
		r.Tx.DeleteEvent(event)
		r.Tx.Commit()
		return nil
	}
	if event.Name = strings.TrimSpace(r.FormValue("name")); event.Name == "" {
		return errors.New("missing name")
	}
	event.Date = r.FormValue("date")
	if event.Date == "" {
		return errors.New("missing date")
	} else if !dateRE.MatchString(event.Date) {
		return errors.New("invalid date (YYYY-MM-DD)")
	}
	if event.Start == "" {
		return errors.New("missing start")
	} else if !timeRE.MatchString(event.Start) {
		return errors.New("invalid start (HH:MM)")
	}
	if event.End == "" {
		return errors.New("missing end")
	} else if !timeRE.MatchString(event.End) {
		return errors.New("invalid end (HH:MM)")
	}
	if event.End < event.Start {
		return errors.New("end before start")
	}
	vidstr := r.FormValue("venue")
	if vidstr == "NEW" {
		event.Venue = &model.Venue{
			Name:    strings.TrimSpace(r.FormValue("venueName")),
			Address: strings.TrimSpace(r.FormValue("venueAddress")),
			City:    strings.TrimSpace(r.FormValue("venueCity")),
			URL:     strings.TrimSpace(r.FormValue("venueURL")),
		}
		if event.Venue.Name == "" || event.Venue.Address == "" || event.Venue.City == "" {
			return errors.New("missing venue name, address, or city")
		}
		if event.Venue.URL != "" && !strings.HasPrefix(event.Venue.URL, "https://www.google.com/maps/") {
			return errors.New("invalid venue URL")
		}
		r.Tx.SaveVenue(event.Venue)
	} else if vidstr == "0" {
		event.Venue = nil
	} else if event.Venue = r.Tx.FetchVenue(model.VenueID(util.ParseID(vidstr))); event.Venue == nil {
		return errors.New("nonexistent venue")
	}
	event.Details = htmlSanitizer.Sanitize(strings.TrimSpace(r.FormValue("details")))
	if event.Type = model.EventType(r.FormValue("type")); event.Type == "" {
		return errors.New("missing type")
	}
	found := false
	for _, t := range model.AllEventTypes {
		if event.Type == t {
			found = true
			break
		}
	}
	if !found {
		return errors.New("invalid type")
	}
	event.Roles = event.Roles[:0]
	for _, idstr := range r.Form["role"] {
		role := r.Tx.FetchRole(model.RoleID(util.ParseID(idstr)))
		if role == nil {
			return errors.New("invalid role")
		}
		if !auth.CanManageEvents(r, role) {
			return util.Forbidden
		}
		event.Roles = append(event.Roles, role)
	}
	if len(event.Roles) == 0 {
		return errors.New("missing role")
	}
	for _, e := range r.Tx.FetchEvents(event.Date, event.Date) {
		if e.ID != event.ID && e.Name == event.Name {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"nameError":true}`))
			return nil
		}
	}
	r.Tx.SaveEvent(event)
	r.Tx.Commit()
	return nil
}
