package signups

import (
	"sort"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/shiftperson"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const ShowTaskSignupsTaskFields = task.FID | shiftperson.EligibilityCheckerTaskFields
const ShowTaskSignupsPersonFields = person.FID | shiftperson.EligibilityCheckerPersonFields

// ShowTaskSignups shows the task signup form for the specified event, task, and
// person.  (It does not include heading information such as task name or
// details.)  The privileged flag indicates whether the caller (which might not
// be the same person) has leader privileges on the task organization.  The
// editable flag indicates that an Edit button should be placed next to each
// shift, and remove links next to each listed person.
//
// In order for this to work, and for HandleShiftSignup to work, the supplied
// form element must have hidden input elements in it named 'shift', 'signedup',
// and if editable is true, 'person'.
func ShowTaskSignups(r *request.Request, form *htmlb.Element, t *task.Task, p *person.Person, privileged, editable bool) {
	const shiftFields = shift.FID | shift.FStart | shift.FEnd | shift.FMin | shift.FMax | shiftperson.EligibilityCheckerShiftFields
	const venueFields = venue.FID | venue.FName | venue.FURL
	var (
		tdiv  *htmlb.Element
		ec    *shiftperson.EligibilityChecker
		lastv = venue.ID(9999999999)
	)
	shift.AllForTask(r, t.ID(), shiftFields, venueFields, func(s *shift.Shift, v *venue.Venue) {
		if tdiv == nil {
			tdiv = form.E("div class=signupTaskShifts")
			ec = shiftperson.NewEligibilityChecker(r, t, p, privileged)
			lastv = 9999999999
		}
		if v.ID() != lastv {
			if v == nil {
				tdiv.E("div class=signupVenueName").R(r.Loc("Location TBD"))
			} else if v.URL() != "" {
				tdiv.E("div class=signupVenueName").
					E("a href=%s target=_blank>%s", v.URL(), v.Name())
			} else {
				tdiv.E("div class=signupVenueName>%s", v.Name())
			}
			lastv = v.ID()
		}
		signedup := shiftperson.Get(r, s.ID(), p.ID()) > 0
		var people []*person.Person
		shiftperson.PeopleForShift(r, s.ID(), person.FID|person.FSortName, func(p *person.Person) {
			pclone := *p
			people = append(people, &pclone)
		})
		sort.Slice(people, func(i, j int) bool { return people[i].SortName() < people[j].SortName() })
		var ineligibleReason shiftperson.IneligibleReason
		if signedup {
			ineligibleReason = ec.CanCancel(s)
		} else {
			ineligibleReason = ec.CanSignUp(s)
		}
		ineligibleReason = shiftperson.IneligibleReason(r.Loc(string(ineligibleReason)))
		label := s.Start()[11:]
		if s.End() != s.Start() {
			label += "–" + s.End()[11:]
		}
		tdiv.E("input type=checkbox class='s-check signupShiftCheck' label=%s data-shift=%d", label, s.ID(),
			signedup, "checked",
			ineligibleReason != "", "disabled", ineligibleReason != "", "title=%s", string(ineligibleReason))
		hdiv := tdiv.E("div class=signupShiftHave")
		if len(people) != 0 {
			hdiv.E("a href=#").TF(r.Loc("Have %d,"), len(people))
		} else {
			hdiv.TF(r.Loc("Have %d,"), len(people))
		}
		if s.Min() != 0 && len(people) < int(s.Min()) {
			tdiv.E("div class=signupShiftNeed").TF(r.Loc("need %d"), s.Min())
		} else if s.Max() != 0 {
			tdiv.E("div class=signupShiftMax").TF(r.Loc("limit %d"), s.Max())
		} else {
			tdiv.E("div class=signupShiftMax").R(r.Loc("no limit"))
		}
		if privileged && editable {
			tdiv.E("div class=signupShiftEdit").
				E("a href=/events/edshift/%d up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Edit", s.ID())
		}
		if ineligibleReason != "" {
			tdiv.E("div class=signupShiftDisabled hidden>%s", string(ineligibleReason))
		}
		if len(people) != 0 {
			list := tdiv.E("div class=signupShiftList hidden")
			for _, p := range people {
				pdiv := list.E("div>%s", p.SortName())
				if privileged && editable {
					pdiv.E("a class=signupShiftRemove data-shift=%d data-person=%d href=#>remove", s.ID(), p.ID())
				}
			}
		}
	})
}

const HandleShiftSignupPersonFields = shiftperson.EligibilityCheckerPersonFields | shiftperson.SignUpPersonFields

// HandleShiftSignup applies a user request to sign up, or cancel signup, a
// person for a shift.  It expects the request to have a shift=%d parameter,
// indicating which shift is being changed, and a signedup=true|false parameter,
// indicating whether the specified person should be signed up for that shift.
// It returns the date of the shift (i.e., the date all of whose shift signups
// should be refreshed in the UI to reflect new eligibility), or an empty string
// if the shift was not found.
//
// The caller is expected to have validated the user and checked CSRF.
func HandleShiftSignup(r *request.Request, user, p *person.Person) (date string) {
	const (
		eventFields = shiftperson.SignUpEventFields
		taskFields  = task.FEvent | task.FOrg | shiftperson.EligibilityCheckerTaskFields | shiftperson.SignUpTaskFields
		shiftFields = shift.FTask | shift.FStart | shiftperson.EligibilityCheckerShiftFields | shiftperson.SignUpShiftFields
	)
	var (
		e        *event.Event
		t        *task.Task
		s        *shift.Shift
		want     bool
		have     bool
		editable bool
		ec       *shiftperson.EligibilityChecker
	)
	if s = shift.WithID(r, shift.ID(util.ParseID(r.FormValue("shift"))), shiftFields); s == nil {
		return
	}
	date = s.Start()[:10]
	t = task.WithID(r, s.Task(), taskFields)
	e = event.WithID(r, t.Event(), eventFields)
	editable = user.HasPrivLevel(t.Org(), enum.PrivLeader)
	if pid := person.ID(util.ParseID(r.FormValue("person"))); pid >= 0 {
		if !editable {
			return
		}
		if p = person.WithID(r, pid, HandleShiftSignupPersonFields); p == nil {
			return
		}
	}
	want = r.FormValue("signedup") == "true"
	have = shiftperson.Get(r, s.ID(), p.ID()) > 0
	if want == have {
		return
	}
	ec = shiftperson.NewEligibilityChecker(r, t, p, editable)
	if (want && ec.CanSignUp(s) != "") || (!want && ec.CanCancel(s) != "") {
		return
	}
	r.Transaction(func() {
		if want {
			shiftperson.SignUp(r, e, t, s, p)
		} else {
			shiftperson.Decline(r, e, t, s, p)
		}
	})
	return
}
