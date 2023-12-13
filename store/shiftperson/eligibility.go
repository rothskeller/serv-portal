package shiftperson

import (
	"time"

	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/taskperson"
)

// An EligibilityChecker can be used to check the eligibility of a person to
// sign up for, or cancel a signup for, a shift.  It maintains cached
// information to streamline multiple such checks for different shifts of the
// same task.
type EligibilityChecker struct {
	storer         phys.Storer
	t              *task.Task
	p              *person.Person
	privileged     bool
	hasRoleCache   bool
	hasRoleChecked bool
}

const EligibilityCheckerPersonFields = person.FID | person.FDSWRegistrations | person.FBGChecks
const EligibilityCheckerTaskFields = task.FID | task.FFlags | task.FOrg
const EligibilityCheckerShiftFields = shift.FID | shift.FStart | shift.FEnd | shift.FMax

// NewEligibilityChecker creates a new EligibilityChecker for the specified task
// and person.  The task must have retrieved EligibilityCheckerTaskFields; the
// person must have retrieved EligibilityCheckerPersonFields.  The privileged
// flag indicates whether the calling user (who may not be the person identified
// by pid) has leader privileges on the task organization.
func NewEligibilityChecker(storer phys.Storer, t *task.Task, p *person.Person, privileged bool) (ec *EligibilityChecker) {
	return &EligibilityChecker{storer: storer, t: t, p: p, privileged: privileged}
}

// hasRole returns whether the person has one of the task roles.  The answer is
// cached to avoid multiple queries.
func (ec *EligibilityChecker) hasRole() bool {
	if ec.p == nil {
		ec.hasRoleChecked = true
	}
	if !ec.hasRoleChecked {
		ec.hasRoleCache = taskperson.HasRoleForTask(ec.storer, ec.t.ID(), ec.p.ID())
		ec.hasRoleChecked = true
	}
	return ec.hasRoleCache
}

type IneligibleReason string

var (
	ErrOverlapping IneligibleReason = "Already signed up for a conflicting shift."
	ErrClosed      IneligibleReason = "Signups are closed."
	ErrIneligible  IneligibleReason = "Not eligible to sign up."
	ErrNoDSW       IneligibleReason = "DSW registration is required."
	ErrNoBGCheck   IneligibleReason = "A background check is required."
	ErrEnded       IneligibleReason = "The shift has ended."
	ErrStarted     IneligibleReason = "The shift has already started."
	ErrFull        IneligibleReason = "The shift is full."
	ErrNoPerson    IneligibleReason = "No person selected."
)

func (ec *EligibilityChecker) CanSignUp(s *shift.Shift) IneligibleReason {
	if ec.p == nil {
		return ErrNoPerson
	}
	if OverlappingSignup(ec.storer, ec.p.ID(), s.Start(), s.End()) {
		return ErrOverlapping
	}
	if ec.privileged {
		return ""
	}
	if ec.t.Flags()&task.SignupsOpen == 0 {
		return ErrClosed
	}
	if !ec.hasRole() {
		return ErrIneligible
	}
	if ec.t.Flags()&task.RequiresBGCheck != 0 {
		if !ec.p.BGChecks().DOJ.Valid() || !ec.p.BGChecks().FBI.Valid() {
			return ErrNoBGCheck
		}
	}
	if ec.t.Flags()&task.CoveredByDSW != 0 {
		if reg, _ := ec.p.DSWRegistrationForOrg(ec.t.Org()); !reg.Valid() {
			return ErrNoDSW
		}
	}
	var end, _ = time.ParseInLocation("2006-01-02T15:04", s.End(), time.Local)
	var now = time.Now()
	if !end.After(now) {
		return ErrEnded
	}
	var start, _ = time.ParseInLocation("2006-01-02T15:04", s.Start(), time.Local)
	if start.Before(now) {
		return ErrStarted
	}
	var signupCount uint
	phys.SQL(ec.storer, "SELECT COUNT(*) FROM shift_person WHERE shift=? AND signed_up>0", func(stmt *phys.Stmt) {
		stmt.BindInt(int(s.ID()))
		stmt.Step()
		signupCount = uint(stmt.ColumnInt())
	})
	if s.Max() != 0 && signupCount >= s.Max() {
		return ErrFull
	}
	return ""
}

func (ec *EligibilityChecker) CanCancel(s *shift.Shift) IneligibleReason {
	if ec.p == nil {
		return ErrNoPerson
	}
	if ec.privileged {
		return ""
	}
	var end, _ = time.ParseInLocation("2006-01-02T15:04", s.End(), time.Local)
	var now = time.Now()
	if !end.After(now) {
		return ErrEnded
	}
	var start, _ = time.ParseInLocation("2006-01-02T15:04", s.Start(), time.Local)
	if start.Before(now) {
		return ErrStarted
	}
	return ""
}
