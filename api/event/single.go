package event

import (
	"errors"
	"math"
	"regexp"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/auth"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetEvent handles GET /api/events/$id requests (where $id may be "NEW").
func GetEvent(r *util.Request, idstr string) error {
	var (
		event   *model.Event
		canEdit bool
		out     jwriter.Writer
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
	}
	r.Tx.Commit()
	out.RawString(`{"canEdit":`)
	out.Bool(canEdit)
	if event != nil {
		out.RawString(`,"event":{"id":`)
		out.Int(int(event.ID))
		out.RawString(`,"date":`)
		out.String(event.Date)
		out.RawString(`,"name":`)
		out.String(event.Name)
		out.RawString(`,"hours":`)
		out.Float64(event.Hours)
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
		out.RawString(`,"event":{"id":0,"date":"","name":"","hours":1.0,"type":"","roles":[]}`)
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
		out.RawByte(']')
	}
	out.RawByte('}')
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

var dateRE = regexp.MustCompile(`^20\d\d-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`)
var yearRE = regexp.MustCompile(`^20\d\d$`)

// PostEvent handles POST /events/$id requests (where $id may be "NEW").
func PostEvent(r *util.Request, idstr string) error {
	var (
		event *model.Event
		err   error
	)
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
	event.Date = r.FormValue("date")
	if event.Date == "" {
		return errors.New("missing date")
	} else if !dateRE.MatchString(event.Date) {
		return errors.New("invalid date (YYYY-MM-DD)")
	}
	if event.Name = strings.TrimSpace(r.FormValue("name")); event.Name == "" {
		return errors.New("missing name")
	}
	event.Hours, err = strconv.ParseFloat(r.FormValue("hours"), 64)
	if err != nil || event.Hours < 0.0 || event.Hours > 24.0 || math.IsNaN(event.Hours) {
		return errors.New("invalid hours")
	}
	event.Hours = math.Round(event.Hours*2.0) / 2.0
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
