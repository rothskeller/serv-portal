package person

import "sunnyvaleserv.org/portal/store/enum"

// ViewLevel describes what details of a Person can be viewed by another Person,
// as a return from the CanView method.
type ViewLevel uint

const (
	// ViewNone indicates that the target person should not be shown at all.
	ViewNone ViewLevel = iota
	// ViewNoContact indicates that the target person can be shown in lists,
	// but without any contact information.
	ViewNoContact
	// ViewWorkContact indicates that the target person can be shown in
	// lists, but only work contact information (emails and work phone)
	// should be visible.
	ViewWorkContact
	// ViewFull indicates that the target person should be visible with all
	// contact information.
	ViewFull
)

const CanViewViewerFields = FPrivLevels
const CanViewTargetFields = FID | FPrivLevels | FFlags

// CanView returns whether the receiver Person is allowed to view the target
// Person.
func (p *Person) CanView(target *Person) (view ViewLevel) {
	// Any person can always view themselves.
	if p.ID() == target.ID() {
		return ViewFull
	}
	// Administrator can be seen only by webmaster.
	if target.ID() == AdminID {
		if p.IsWebmaster() {
			return ViewFull
		} else {
			return ViewNone
		}
	}
	// Anyone can see the work info of visible-to-all folks.
	if target.Flags()&VisibleToAll != 0 {
		view = ViewWorkContact
	} else if target.IsAdminLeader() {
		// Admin leaders (i.e., OES paid staff) are fully visible to
		// each other, viewable with work info only to org leaders, and
		// not visible to anyone else.
		if p.IsAdminLeader() {
			return ViewFull
		} else if p.HasPrivLevel(0, enum.PrivLeader) {
			return ViewWorkContact
		} else {
			return ViewNone
		}
	} else if p.HasPrivLevel(0, enum.PrivLeader) {
		// Org leaders can view anyone else fully.
		return ViewFull
	}
	// Walk through the organizations looking for matches.
	for _, o := range enum.AllOrgs {
		switch p.PrivLevels()[o] {
		case enum.PrivMember:
			switch target.PrivLevels()[o] {
			case 0:
				// nothing
			case enum.PrivStudent:
				// Members can see the list of students but no
				// contact info.
				return ViewNoContact
			default:
				// Members can see leaders with full contact
				// info.  But only direct leaders, not those
				// who inherited from admin leader.
				if !target.IsAdminLeader() {
					return ViewFull
				}
				fallthrough
			case enum.PrivMember:
				// Members can see other members, with full
				// contact info except for SARES.
				if o == enum.OrgSARES {
					view = max(view, ViewNoContact)
				} else {
					return ViewFull
				}
			}
		case enum.PrivStudent:
			switch target.PrivLevels()[o] {
			case 0, enum.PrivStudent:
				// nothing
			default:
				// Students can see members and leaders with no
				// contact info.
				view = max(view, ViewNoContact)
			}
		}
	}
	return view
}
