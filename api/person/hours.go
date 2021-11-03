package person

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// Regular expression for a valid month string.
var monthRE = regexp.MustCompile(`^20\d\d-(?:0[1-9]|1[0-2])$`)

// GPPersonHoursMonth handles GET and POST /api/people/$id/hours/$month
// requests.
func GPPersonHoursMonth(r *util.Request, idstr, month string) error {
	var (
		person        *model.Person
		out           jwriter.Writer
		today         string
		editableMonth bool
		first         = true
	)
	// idstr could be the ID of a person, when used in a regular session, or
	// it could be the HoursToken of a person, when used outside a session.
	// We assume that anything longer than 5 characters must be an
	// HoursToken and anything less must be an ID.
	if len(idstr) <= 5 {
		if person = r.Tx.FetchPerson(model.PersonID(util.ParseID(idstr))); person == nil {
			return util.NotFound
		}
	} else {
		if person = r.Tx.FetchPersonByHoursToken(idstr); person == nil {
			return util.NotFound
		}
		r.Person = person
	}
	if r.Person == nil || (person != r.Person && !r.Person.HasPrivLevel(model.PrivLeader)) {
		return util.Forbidden
	}
	if !monthRE.MatchString(month) {
		return util.NotFound
	}
	today, editableMonth = isEditableMonth(month)
	if r.Person.Roles[model.Webmaster] {
		editableMonth = true
	}
	if person == r.Person && person.HoursReminder {
		r.Tx.WillUpdatePerson(person)
		person.HoursReminder = false
		r.Tx.UpdatePerson(person)
	}
	if r.Method == http.MethodGet {
		out.RawString(`{"id":`)
		out.Int(int(person.ID))
		out.RawString(`,"name":`)
		out.String(person.InformalName)
		if person.VolgisticsID == 0 {
			out.RawString(`,"needsVolgistics":true`)
		}
		out.RawString(`,"events":[`)
	}
	// Since we're just doing a <= comparison on strings, it doesn't matter
	// how many days there are in the month.
	for _, e := range r.Tx.FetchEvents(month+"-01", month+"-31") {
		var (
			amap                          = r.Tx.FetchAttendanceByEvent(e)
			attend                        = amap[person.ID]
			canView, canViewType, canEdit = hoursPermissions(r.Person, person, e, today, editableMonth, attend.Minutes != 0)
		)
		if !canView {
			continue
		}
		switch r.Method {
		case http.MethodGet:
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
			out.RawString(`,"minutes":`)
			out.Uint16(attend.Minutes)
			out.RawString(`,"type":`)
			out.String(attend.Type.String())
			if e.Type == model.EventHours {
				out.RawString(`,"placeholder":true`)
			}
			if canViewType {
				out.RawString(`,"canViewType":true`)
			}
			if canEdit {
				out.RawString(`,"canEdit":true`)
			}
			if e.RenewsDSW && attend.Type == model.AttendAsVolunteer && attend.Minutes > 0 {
				out.RawString(`,"renewsDSW":true`)
			}
			out.RawByte('}')
		case http.MethodPost:
			var (
				value   string
				parts   []string
				minutes int
				err     error
			)
			if !canEdit {
				continue
			}
			if value = r.FormValue(fmt.Sprintf("e%d", e.ID)); value == "" {
				return fmt.Errorf("missing e%d", e.ID)
			}
			parts = strings.Split(value, ":")
			if len(parts) != 2 {
				return fmt.Errorf("invalid e%d", e.ID)
			}
			if minutes, err = strconv.Atoi(parts[0]); err != nil || minutes < 0 {
				return fmt.Errorf("invalid e%d", e.ID)
			}
			if canViewType {
				if attend.Type, err = model.ParseAttendanceType(parts[1]); err != nil {
					return fmt.Errorf("invalid e%d", e.ID)
				}
			} else if attend.Minutes == 0 {
				attend.Type = model.AttendAsAbsent
			}
			if minutes == 0 && attend.Minutes != 0 {
				delete(amap, person.ID)
				r.Tx.SaveEventAttendance(e, amap)
			} else if minutes != 0 {
				attend.Minutes = uint16(minutes)
				amap[person.ID] = attend
				r.Tx.SaveEventAttendance(e, amap)
			}
		}
	}
	r.Tx.Commit()
	if r.Method == http.MethodGet {
		out.RawString(`]}`)
		r.Header().Set("Content-Type", "application/json; charset=utf-8")
		out.DumpTo(r)
	}
	return nil
}

// isEditableMonth returns whether attendance for events in the month are
// editable.  A month's event attendance is editable from the start of that
// month through the 10th of the following month.
func isEditableMonth(month string) (today string, editableMonth bool) {
	var now = time.Now()
	today = now.Format("2006-01-02")
	if month > today[:7] {
		return today, false
	}
	if now.Day() > 10 && month < today[:7] {
		return today, false
	}
	lastMonth := time.Date(now.Year(), now.Month()-1, 1, 0, 0, 0, 0, time.Local).Format("2006-01")
	return today, month >= lastMonth
}

// hoursPermissions returns the permissions the current caller has for viewing
// or editing the target person's attendance at the specified event.  today and
// editableMonth are memoized data returned from isEditableMonth, above.
// hasHours indicates whether the target person already has recorded hours for
// the event.
func hoursPermissions(caller, target *model.Person, event *model.Event, today string, editableMonth, hasHours bool) (canView, canViewType, canEdit bool) {
	// We need to know whether attendance at the event itself is editable.
	eventAttendanceEditable := editableMonth && event.Type != model.EventSocial
	if eventAttendanceEditable && event.Date > today && event.Type != model.EventHours {
		eventAttendanceEditable = false
	}
	callerIsLeader := caller.Orgs[event.Org].PrivLevel >= model.PrivLeader
	// The type can be edited only if the caller is an event org leader,
	// attendance for the event itself is editable, and the event is not a
	// placeholder.
	if eventAttendanceEditable && callerIsLeader && (event.Type != model.EventHours || caller.Roles[model.Webmaster]) {
		return true, true, true
	}
	// The hours can be edited if the caller is the target person, the
	// attendance for the event itself is editable, and either the caller
	// attended the event, the caller is a member of the event org, or the
	// event is a placeholder for an open org.
	if eventAttendanceEditable && caller == target &&
		(hasHours || target.Orgs[event.Org].PrivLevel >= model.PrivMember || isOpenPlaceholder(event)) {
		return true, false, true
	}
	// The event can be seen, with the attendance type, if it has hours
	// recorded and the caller is an event org leader.
	if hasHours && callerIsLeader {
		return true, true, false
	}
	// The event can be seen, without the attendance type, if it has hours
	// recorded and the caller is the target person.
	if hasHours && caller == target {
		return true, false, false
	}
	// Otherwise, no permissions.
	return false, false, false
}

// isOpenPlaceholder returns whether the event is a placeholder (i.e., "Other
// XXX Hours") for an org that anyone can report hours for, whether they belong
// to it or not.
func isOpenPlaceholder(event *model.Event) bool {
	// Admin is the only one that you have to belong to in order to report
	// hours.
	return event.Type == model.EventHours && event.Org != model.OrgAdmin
}
