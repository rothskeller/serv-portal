package clearrep

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type parameters struct {
	role      *role.Role
	with      string
	without   string
	renderCSV bool
	// Not really a parameter, but cached here for convenience:
	allowedRoles        []*role.Role
	allowedRestrictions []string
}

func readParameters(r *request.Request, user *person.Person) (params parameters) {
	if params.allowedRoles = allowedRoles(r, user); len(params.allowedRoles) == 0 {
		return params // caller will raise error
	}
	params.allowedRestrictions = validRestrictions(user)
	if rolestr := r.FormValue("role"); rolestr != "0" {
		params.role = role.WithID(r, role.ID(util.ParseID(rolestr)), role.FOrg|role.FFlags)
		if params.role == nil || !allowedRole(user, params.role) {
			params.role = params.allowedRoles[0]
		}
	}
	if params.with = r.FormValue("with"); !validRestriction(user, params.with) {
		params.with = ""
	}
	if params.without = r.FormValue("without"); !validRestriction(user, params.without) {
		params.without = ""
	}
	if r.FormValue("format") == "csv" {
		params.renderCSV = true
	}
	return params
}

func personMatch(r *request.Request, p *person.Person, params parameters) bool {
	if params.role != nil {
		held, _ := personrole.PersonHasRole(r, p.ID(), params.role.ID())
		return held
	}
	for _, role := range params.allowedRoles {
		if held, _ := personrole.PersonHasRole(r, p.ID(), role.ID()); held {
			return true
		}
	}
	return false
}

func renderParams(main *htmlb.Element, params parameters) {
	form := main.E("form class=clearrepForm")
	form.E("div>Show")
	sel := form.E("select name=role")
	sel.E("option value=0", params.role == nil, "selected").R("Everyone")
	for _, ar := range params.allowedRoles {
		sel.E("option value=%d", ar.ID(), params.role != nil && params.role.ID() == ar.ID(), "selected").T(ar.Name())
	}
	form.E("div>with")
	sel = form.E("select name=with")
	sel.E("option value=''", params.with == "", "selected").R("&mdash;")
	for _, ar := range params.allowedRestrictions {
		sel.E("option value=%s", ar, params.with == ar, "selected").T(restrictionLabels[ar])
	}
	form.E("div>and without")
	sel = form.E("select name=without")
	sel.E("option value=''", params.without == "", "selected").R("&mdash;")
	for _, ar := range params.allowedRestrictions {
		sel.E("option value=%s", ar, params.without == ar, "selected").T(restrictionLabels[ar])
	}
}

func allowedRole(user *person.Person, rl *role.Role) bool {
	return rl.Flags()&role.Filter != 0 && user.HasPrivLevel(rl.Org(), enum.PrivLeader)
}
func allowedRoles(r *request.Request, user *person.Person) (roles []*role.Role) {
	const roleFields = role.FFlags | role.FOrg | role.FID | role.FName
	role.All(r, roleFields, func(rl *role.Role) {
		if allowedRole(user, rl) {
			clone := *rl
			roles = append(roles, &clone)
		}
	})
	return roles
}
