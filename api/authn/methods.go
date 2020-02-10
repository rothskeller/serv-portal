package authn

import (
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// IsEnabled returns whether the specified Person is enabled.
func IsEnabled(r *util.Request, p *model.Person) bool {
	if r.Auth.MemberPG(p.ID, r.Auth.FetchGroupByTag(model.GroupDisabled).ID) {
		return false // person is explicitly disabled
	}
	if !r.Auth.CanPA(p.ID, model.PrivMember) {
		return false // person belongs to no groups
	}
	return true
}
