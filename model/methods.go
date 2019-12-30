package model

// CanCreateEvents returns whether the receiver Person is allowed to create
// events.
func (p *Person) CanCreateEvents() bool {
	if p == nil {
		return false
	}
	return p.PrivMap.HasAny(PrivAdmin)
}

// CanLogIn returns whether the receiver Person is allowed to log in.  This does
// not check whether a password reset is in progress or mhether the account has
// a temporary lockout due to bad password attempts; it only checks whether they
// belong to a group that allows them to log in.
func (p *Person) CanLogIn() bool {
	if p == nil {
		return false
	}
	return p.PrivMap.HasTag(TeamLogin, PrivMember)
}

// CanManageEvent returns whether the receiver Person is allowed to manage
// (i.e., edit or delete) the specified Event.
func (p *Person) CanManageEvent(e *Event) bool {
	if p == nil {
		return false
	}
	for _, t := range e.Teams {
		if !p.PrivMap.Has(t, PrivAdmin) {
			return false
		}
	}
	return true
}

// CanRecordAttendanceAtEvent returns whether the receiver Person is allowed to
// record attendance at the specified Event.
func (p *Person) CanRecordAttendanceAtEvent(e *Event) bool {
	if p == nil {
		return false
	}
	for _, t := range e.Teams {
		if p.PrivMap.Has(t, PrivAdmin) {
			return true
		}
	}
	return false
}

// CanViewEvent returns whether the receiver Person is allowed to see the
// specified Event.
func (p *Person) CanViewEvent(e *Event) bool {
	for _, t := range e.Teams {
		if p.IsMember(t) {
			return true
		}
	}
	if p.IsWebmaster() {
		return true
	}
	return false
}

// CanViewPerson returns whether the receiver Person can view the argument
// Person.
func (p *Person) CanViewPerson(p2 *Person) bool {
	if p == nil {
		return false
	}
	for _, r := range p2.Roles {
		if p.PrivMap.Has(r.Team, PrivView) {
			return true
		}
	}
	if p.IsWebmaster() { // in case p2 has no roles
		return true
	}
	return false
}

// IsMember returns whether the receiver Person is a member of the specified
// Team (directly or indirectly).
func (p *Person) IsMember(t *Team) bool {
	if p == nil {
		return false
	}
	return p.PrivMap.Has(t, PrivMember)
}

// IsWebmaster returns whether the receiver Person is a webmaster.
func (p *Person) IsWebmaster() bool {
	if p == nil {
		return false
	}
	return p.PrivMap.HasTag(TeamWebmasters, PrivMember)
}

// ManagedTeams returns the list of teams that the receiver Person is allowed
// to manage.
func (p *Person) ManagedTeams() (teams []*Team) {
	if p == nil {
		return nil
	}
	teams = p.PrivMap.TeamsWith(PrivManage)
	j := 0
	for _, t := range teams {
		if t.Type == TeamNormal {
			teams[j] = t
			j++
		}
	}
	return teams[:j]
}

// SchedulableTeams returns the list of teams for which the receiver Person is
// allowed to create new Events.
func (p *Person) SchedulableTeams() []*Team {
	if p == nil {
		return nil
	}
	return p.PrivMap.TeamsWith(PrivAdmin)
}

// ViewableTeams returns a list of teams that the receiver Person is allowed to
// see.  It returns nil if they are not allowed to see any teams.
func (p *Person) ViewableTeams() []*Team {
	if p == nil {
		return nil
	}
	return p.PrivMap.TeamsWith(PrivView)
}

// Depth returns the depth of the receiver Team in the team hierarchy.  A
// parentless Team has depth 0.
func (t *Team) Depth() (indent int) {
	for ; t.Parent != nil; t = t.Parent {
		indent++
	}
	return indent
}
