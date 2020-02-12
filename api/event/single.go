package event

import (
	"errors"
	"fmt"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetEvent handles GET /api/events/$id requests (where $id may be "NEW").
func GetEvent(r *util.Request, idstr string) error {
	var (
		event          *model.Event
		canView        bool
		canEdit        bool
		canAttendance  bool
		out            jwriter.Writer
		wantAttendance = r.FormValue("attendance") != ""
		wantEdit       = r.FormValue("edit") != ""
	)
	if idstr == "NEW" {
		if !r.Auth.CanA(model.PrivManageEvents) {
			return util.Forbidden
		}
		event = new(model.Event)
		canEdit = true
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		canEdit = true
		for _, group := range event.Groups {
			if r.Auth.MemberG(group) {
				canView = true
			}
			if r.Auth.CanAG(model.PrivManageEvents, group) {
				canView = true
				canAttendance = true
			} else {
				canEdit = false
			}
		}
		if !canView {
			return util.Forbidden
		}
		if event.SccAresID != "" {
			canEdit = false
		}
	}
	out.RawString(`{"event":{"id":`)
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
	if event.Venue != 0 {
		venue := r.Tx.FetchVenue(event.Venue)
		out.RawString(`{"id":`)
		out.Int(int(event.Venue))
		out.RawString(`,"name":`)
		out.String(venue.Name)
		out.RawString(`,"address":`)
		out.String(venue.Address)
		out.RawString(`,"city":`)
		out.String(venue.City)
		out.RawString(`,"url":`)
		out.String(venue.URL)
		out.RawByte('}')
	} else {
		out.RawString(`{"id":0,"name":"","address":"","city":"","url":""}`)
	}
	out.RawString(`,"details":`)
	out.String(event.Details)
	out.RawString(`,"organization":`)
	out.String(model.OrganizationNames[event.Organization])
	out.RawString(`,"types":[`)
	first := true
	for _, et := range model.AllEventTypes {
		if event.Type&et != 0 {
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(model.EventTypeNames[et])
		}
	}
	out.RawString(`],"groups":[`)
	for i, g := range event.Groups {
		if i != 0 {
			out.RawByte(',')
		}
		out.Int(int(g))
	}
	out.RawString(`],"canEdit":`)
	out.Bool(canEdit)
	out.RawString(`,"canAttendance":`)
	out.Bool(canAttendance)
	out.RawByte('}')
	if canEdit && wantEdit {
		out.RawString(`,"types":[`)
		for i, et := range model.AllEventTypes {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(model.EventTypeNames[et])
		}
		out.RawString(`],"groups":[`)
		for i, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivManageEvents)) {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(g.ID))
			out.RawString(`,"name":`)
			out.String(g.Name)
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
		out.RawString(`],"organizations":[`)
		for i, o := range model.AllOrganizations {
			if i != 0 {
				out.RawByte(',')
			}
			out.String(model.OrganizationNames[o])
		}
		out.RawByte(']')
	}
	if canAttendance && wantAttendance {
		var (
			attended = r.Tx.FetchAttendanceByEvent(event)
			first    = true
		)
		out.RawString(`,"people":[`)
		for _, p := range r.Auth.FetchPeople(r.Auth.PeopleA(model.PrivViewMembers)) {
			ai, att := attended[p.ID]
			canView := att
			for _, group := range event.Groups {
				if r.Auth.MemberPG(p.ID, group) {
					canView = true
					break
				}
			}
			if !canView {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(p.ID))
			out.RawString(`,"sortName":`)
			out.String(p.SortName)
			if att {
				out.RawString(`,"attended":{"type":`)
				out.String(model.AttendanceTypeNames[ai.Type])
				out.RawString(`,"minutes":`)
				out.Uint16(ai.Minutes)
				out.RawByte('}')
			} else {
				out.RawString(`,"attended":false`)
			}
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

// PostEvent handles POST /events/$id requests (where $id may be "NEW").
func PostEvent(r *util.Request, idstr string) error {
	var (
		event *model.Event
		org   string
	)
	if idstr == "NEW" {
		if !r.Auth.CanA(model.PrivManageEvents) {
			return util.Forbidden
		}
		event = new(model.Event)
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		for _, group := range event.Groups {
			if !r.Auth.CanAG(model.PrivManageEvents, group) {
				return util.Forbidden
			}
		}
	}
	if r.FormValue("delete") != "" && event.ID != 0 {
		r.Tx.DeleteEvent(event)
		r.Tx.Commit()
		return nil
	}
	event.Name = r.FormValue("name")
	event.Date = r.FormValue("date")
	event.Start = r.FormValue("start")
	event.End = r.FormValue("end")
	vidstr := r.FormValue("venue")
	if vidstr == "NEW" {
		venue := &model.Venue{
			Name:    strings.TrimSpace(r.FormValue("venueName")),
			Address: strings.TrimSpace(r.FormValue("venueAddress")),
			City:    strings.TrimSpace(r.FormValue("venueCity")),
			URL:     strings.TrimSpace(r.FormValue("venueURL")),
		}
		if venue.Name == "" || venue.Address == "" || venue.City == "" {
			return errors.New("missing venue name, address, or city")
		}
		if venue.URL != "" && !strings.HasPrefix(venue.URL, "https://www.google.com/maps/") {
			return errors.New("invalid venue URL")
		}
		r.Tx.CreateVenue(venue)
		event.Venue = venue.ID
	} else if vidstr == "0" {
		event.Venue = 0
	} else if event.Venue = model.VenueID(util.ParseID(vidstr)); event.Venue == 0 {
		return errors.New("invalid venue")
	}
	event.Details = r.FormValue("details")
	org = r.FormValue("organization")
	event.Organization = model.OrgNone
	if org != "" {
		for _, o := range model.AllOrganizations {
			if org == model.OrganizationNames[o] {
				event.Organization = o
			}
		}
		if event.Organization == 0 {
			return errors.New("invalid organization")
		}
	}
	event.Type = 0
	for _, v := range r.Form["type"] {
		var matched bool
		for _, et := range model.AllEventTypes {
			if model.EventTypeNames[et] == v {
				event.Type |= et
				matched = true
				break
			}
		}
		if !matched {
			return errors.New("invalid type")
		}
	}
	event.Groups = event.Groups[:0]
	for _, idstr := range r.Form["group"] {
		var gid = model.GroupID(util.ParseID(idstr))
		if !r.Auth.CanAG(model.PrivManageEvents, gid) {
			return util.Forbidden
		}
		event.Groups = append(event.Groups, gid)
	}
	if err := ValidateEvent(r.Tx, event); err != nil {
		if err.Error() == "duplicate name" {
			r.Header().Set("Content-Type", "application/json; charset=utf-8")
			r.Write([]byte(`{"nameError":true}`))
			return nil
		}
		return err
	}
	if event.ID != 0 {
		r.Tx.UpdateEvent(event)
	} else {
		r.Tx.CreateEvent(event)
	}
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(r, `{"id":%d}`, event.ID)
	return nil
}
