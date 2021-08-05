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
		event             *model.Event
		canEdit           bool
		canViewAttendance bool
		canEditAttendance bool
		out               jwriter.Writer
		wantAttendance    = r.FormValue("attendance") != ""
		wantEdit          = r.FormValue("edit") != ""
	)
	if idstr == "NEW" {
		if !r.Person.HasPrivLevel(model.PrivLeader) {
			return util.Forbidden
		}
		event = new(model.Event)
		canEdit = true
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		if r.Person.Orgs[event.Org].PrivLevel >= model.PrivLeader {
			canEdit = true
			canViewAttendance = true
			if !attendanceFinalized(event.Date) || r.Person.Roles[model.Webmaster] {
				canEditAttendance = true
			}
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
	out.RawString(`,"renewsDSW":`)
	out.Bool(event.RenewsDSW)
	out.RawString(`,"coveredByDSW":`)
	out.Bool(event.CoveredByDSW)
	out.RawString(`,"org":`)
	out.String(event.Org.String())
	out.RawString(`,"type":`)
	out.String(model.EventTypeNames[event.Type])
	out.RawString(`,"roles":[`)
	for i, r := range event.Roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.Int(int(r))
	}
	out.RawString(`],"canEdit":`)
	out.Bool(canEdit)
	out.RawString(`,"canViewAttendance":`)
	out.Bool(canViewAttendance)
	out.RawString(`,"canEditAttendance":`)
	out.Bool(canEditAttendance)
	out.RawString(`,"canEditDSWFlags":`)
	out.Bool(r.Person.IsAdminLeader())
	out.RawByte('}')
	if canEdit && wantEdit {
		out.RawString(`,"types":[`)
		var first = true
		for _, et := range model.AllEventTypes {
			if et == model.EventHours {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(model.EventTypeNames[et])
		}
		out.RawString(`],"roles":[`)
		first = true
		for _, role := range r.Tx.FetchRoles() {
			if !role.ShowRoster || r.Person.Orgs[role.Org].PrivLevel < model.PrivLeader {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.Int(int(role.ID))
			out.RawString(`,"name":`)
			out.String(role.Name)
			out.RawString(`,"org":`)
			out.String(role.Org.String())
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
		out.RawString(`],"orgs":[`)
		first = true
		for _, o := range model.AllOrgs {
			if r.Person.Orgs[o].PrivLevel < model.PrivLeader {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.String(o.String())
		}
		out.RawByte(']')
	}
	if canViewAttendance && wantAttendance {
		var (
			attended = r.Tx.FetchAttendanceByEvent(event)
			first    = true
		)
		out.RawString(`,"people":[`)
		for _, p := range r.Tx.FetchPeople() {
			ai, att := attended[p.ID]
			show := att
			for _, role := range event.Roles {
				if _, ok := p.Roles[role]; ok {
					show = true
					break
				}
			}
			if !show {
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
			out.RawString(`,"callSign":`)
			out.String(p.CallSign)
			if att {
				out.RawString(`,"attended":{"type":`)
				out.String(ai.Type.String())
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
		err   error
	)
	if idstr == "NEW" {
		if !r.Person.HasPrivLevel(model.PrivLeader) {
			return util.Forbidden
		}
		event = new(model.Event)
	} else {
		if event = r.Tx.FetchEvent(model.EventID(util.ParseID(idstr))); event == nil {
			return util.NotFound
		}
		if r.Person.Orgs[event.Org].PrivLevel < model.PrivLeader {
			return util.Forbidden
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
	if r.Person.IsAdminLeader() {
		event.RenewsDSW = r.FormValue("renewsDSW") == "true"
		event.CoveredByDSW = r.FormValue("coveredByDSW") == "true"
	}
	if event.Org, err = model.ParseOrg(r.FormValue("org")); err != nil {
		return err
	}
	if r.Person.Orgs[event.Org].PrivLevel < model.PrivLeader {
		return errors.New("forbidden org")
	}
	event.Type = 0
	for _, et := range model.AllEventTypes {
		if model.EventTypeNames[et] == r.FormValue("type") {
			event.Type = et
			break
		}
	}
	if event.Type == 0 {
		return errors.New("invalid type")
	}
	event.Roles = event.Roles[:0]
	for _, idstr := range r.Form["role"] {
		var rid = model.RoleID(util.ParseID(idstr))
		event.Roles = append(event.Roles, rid)
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
