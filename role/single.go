package role

/*
import (
	"errors"
	"strconv"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/auth"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetRole handles GET /api/roles/$id requests (where $id may be "NEW").
func GetRole(r *util.Request, idstr string) error {
	var (
		role     *model.Role
		allRoles []*model.Role
		forced   map[*model.Role]model.PrivilegeMap
		out      jwriter.Writer
	)
	allRoles = r.Tx.FetchRoles()
	if idstr == "NEW" {
		role = &model.Role{
			ID:      r.Tx.UnusedRoleID(),
			PrivMap: model.PrivilegeMap{},
		}
		role.PrivMap.Add(role.ID, model.PrivHoldsRole)
		allRoles = append([]*model.Role{}, allRoles...)
		allRoles = append(allRoles, role)
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
	}
	if !auth.IsWebmaster(r) {
		return util.Forbidden
	}
	r.Tx.Commit()
	if forced = enforceRoleConstraints(allRoles); forced == nil {
		return errors.New("uncorrectable error in role graph (cycle?)")
	}
	out.RawString(`{"role":{"id":`)
	out.Int(int(role.ID))
	out.RawString(`,"name":`)
	out.String(role.Name)
	out.RawString(`,"memberLabel":`)
	out.String(role.MemberLabel)
	out.RawString(`,"implyOnly":`)
	out.Bool(role.ImplyOnly)
	out.RawString(`,"individual":`)
	out.Bool(role.Individual)
	out.RawString(`},"canDelete":`)
	out.Bool(role.Tag == "")
	out.RawString(`,"privs":`)
	emitRolePrivs(&out, role, allRoles, forced)
	out.RawByte('}')
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRoleReloadPrivs handles POST /api/roles/$id/reloadPrivs requests (where
// $id may be "NEW").
func PostRoleReloadPrivs(r *util.Request, idstr string) error {
	var (
		role     *model.Role
		allRoles []*model.Role
		forced   map[*model.Role]model.PrivilegeMap
		out      jwriter.Writer
	)
	allRoles = r.Tx.FetchRoles()
	if idstr == "NEW" {
		role = &model.Role{
			ID:      r.Tx.UnusedRoleID(),
			PrivMap: model.PrivilegeMap{},
		}
		role.PrivMap.Add(role.ID, model.PrivHoldsRole)
		allRoles = append([]*model.Role{}, allRoles...)
		allRoles = append(allRoles, role)
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
	}
	if !auth.IsWebmaster(r) {
		return util.Forbidden
	}
	readPrivilegesFromRequest(r, role, allRoles)
	if forced = enforceRoleConstraints(allRoles); forced == nil {
		return errors.New("uncorrectable error in role graph (cycle?)")
	}
	emitRolePrivs(&out, role, allRoles, forced)
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostRole handles POST /api/roles/$id requests (where $id may be "NEW").
func PostRole(r *util.Request, idstr string) error {
	var (
		role     *model.Role
		allRoles []*model.Role
	)
	allRoles = r.Tx.FetchRoles()
	if idstr == "NEW" {
		role = &model.Role{
			ID:      r.Tx.UnusedRoleID(),
			PrivMap: model.PrivilegeMap{},
		}
		role.PrivMap.Add(role.ID, model.PrivHoldsRole)
		allRoles = append([]*model.Role{}, allRoles...)
		allRoles = append(allRoles, role)
	} else {
		if role = r.Tx.FetchRole(model.RoleID(util.ParseID(idstr))); role == nil {
			return util.NotFound
		}
	}
	if !auth.IsWebmaster(r) {
		return util.Forbidden
	}
	if r.FormValue("delete") != "" && role.Tag == "" && idstr != "NEW" {
		return deleteRole(r, role)
	}
	if role.Name = strings.TrimSpace(r.FormValue("name")); role.Name == "" {
		return errors.New("missing name")
	}
	for _, r2 := range allRoles {
		if r2 != role && r2.Name == role.Name {
			r.Header().Set("Content-Type", "application/json")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		}
	}
	role.MemberLabel = strings.TrimSpace(r.FormValue("memberLabel"))
	role.ImplyOnly = r.FormValue("implyOnly") == "true"
	role.Individual = r.FormValue("individual") == "true"
	readPrivilegesFromRequest(r, role, allRoles)
	if enforceRoleConstraints(allRoles) == nil {
		return errors.New("uncorrectable error in role graph (cycle?)")
	}
	r.Tx.SaveRole(role)
	r.Tx.SavePrivileges()
	r.Tx.Commit()
	return nil
}

func emitRolePrivs(out *jwriter.Writer, role *model.Role, roles []*model.Role, forced map[*model.Role]model.PrivilegeMap) {
	out.RawByte('[')
	for i, or := range roles {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(or.ID))
		out.RawString(`,"name":`)
		out.String(or.Name)
		out.RawString(`,"actor":{"holdsRole":`)
		out.Bool(role.TransPrivs.Has(or.ID, model.PrivHoldsRole))
		out.RawString(`,"holdsRoleEnabled":`)
		out.Bool(!forced[role].Has(or.ID, model.PrivHoldsRole))
		out.RawString(`,"viewHolders":`)
		out.Bool(role.TransPrivs.Has(or.ID, model.PrivViewHolders))
		out.RawString(`,"viewHoldersEnabled":`)
		out.Bool(!forced[role].Has(or.ID, model.PrivViewHolders))
		out.RawString(`,"assignRole":`)
		out.Bool(role.TransPrivs.Has(or.ID, model.PrivAssignRole))
		out.RawString(`,"assignRoleEnabled":`)
		out.Bool(!forced[role].Has(or.ID, model.PrivAssignRole))
		out.RawString(`,"manageEvents":`)
		out.Bool(role.TransPrivs.Has(or.ID, model.PrivManageEvents))
		out.RawString(`,"manageEventsEnabled":`)
		out.Bool(!forced[role].Has(or.ID, model.PrivManageEvents))
		out.RawString(`},"target":{"holdsRole":`)
		out.Bool(or.TransPrivs.Has(role.ID, model.PrivHoldsRole))
		out.RawString(`,"holdsRoleEnabled":`)
		out.Bool(!forced[or].Has(role.ID, model.PrivHoldsRole))
		out.RawString(`,"viewHolders":`)
		out.Bool(or.TransPrivs.Has(role.ID, model.PrivViewHolders))
		out.RawString(`,"viewHoldersEnabled":`)
		out.Bool(!forced[or].Has(role.ID, model.PrivViewHolders))
		out.RawString(`,"assignRole":`)
		out.Bool(or.TransPrivs.Has(role.ID, model.PrivAssignRole))
		out.RawString(`,"assignRoleEnabled":`)
		out.Bool(!forced[or].Has(role.ID, model.PrivAssignRole))
		out.RawString(`,"manageEvents":`)
		out.Bool(or.TransPrivs.Has(role.ID, model.PrivManageEvents))
		out.RawString(`,"manageEventsEnabled":`)
		out.Bool(!forced[or].Has(role.ID, model.PrivManageEvents))
		out.RawString(`}}`)
	}
	out.RawByte(']')
}

func deleteRole(r *util.Request, role *model.Role) error {
	r.Tx.DeleteRole(role)
	r.Tx.Commit()
	return nil
}

func holdsRoleEnabled(actor, target *model.Role) bool {
	if actor == target {
		return false // can't turn off a role "holding" itself
	}
	if target.TransPrivs.Has(actor.ID, model.PrivHoldsRole) {
		return false // adding this implication would create a cycle
	}
	if actor.TransPrivs.Has(target.ID, model.PrivHoldsRole) && !actor.PrivMap.Has(target.ID, model.PrivHoldsRole) {
		return false // it's forced by some other role implication
	}
	return true
}

func privilegeEnabled(actor, target, webmaster *model.Role, priv model.Privilege) bool {
	if actor.TransPrivs.Has(webmaster.ID, model.PrivHoldsRole) {
		return false // webmaster privileges can't be disabled
	}
	if actor.TransPrivs.Has(target.ID, priv) && !actor.PrivMap.Has(target.ID, priv) {
		return false // it's forced by some other role implication
	}
	return true
}

func readPrivilegesFromRequest(r *util.Request, role *model.Role, allRoles []*model.Role) {
	for _, or := range allRoles {
		var orid = strconv.Itoa(int(or.ID))
		switch r.FormValue("a:holdsRole-" + orid) {
		case "true":
			role.PrivMap = role.PrivMap.Add(or.ID, model.PrivHoldsRole)
		case "false":
			role.PrivMap = role.PrivMap.Remove(or.ID, model.PrivHoldsRole)
		}
		switch r.FormValue("a:viewHolders-" + orid) {
		case "true":
			role.PrivMap = role.PrivMap.Add(or.ID, model.PrivViewHolders)
		case "false":
			role.PrivMap = role.PrivMap.Remove(or.ID, model.PrivViewHolders)
		}
		switch r.FormValue("a:assignRole-" + orid) {
		case "true":
			role.PrivMap = role.PrivMap.Add(or.ID, model.PrivAssignRole)
		case "false":
			role.PrivMap = role.PrivMap.Remove(or.ID, model.PrivAssignRole)
		}
		switch r.FormValue("a:manageEvents-" + orid) {
		case "true":
			role.PrivMap = role.PrivMap.Add(or.ID, model.PrivManageEvents)
		case "false":
			role.PrivMap = role.PrivMap.Remove(or.ID, model.PrivManageEvents)
		}
		if or == role {
			continue
		}
		switch r.FormValue("t:holdsRole-" + orid) {
		case "true":
			or.PrivMap = or.PrivMap.Add(role.ID, model.PrivHoldsRole)
		case "false":
			or.PrivMap = or.PrivMap.Remove(role.ID, model.PrivHoldsRole)
		}
		switch r.FormValue("t:viewHolders-" + orid) {
		case "true":
			or.PrivMap = or.PrivMap.Add(role.ID, model.PrivViewHolders)
		case "false":
			or.PrivMap = or.PrivMap.Remove(role.ID, model.PrivViewHolders)
		}
		switch r.FormValue("t:assignRole-" + orid) {
		case "true":
			or.PrivMap = or.PrivMap.Add(role.ID, model.PrivAssignRole)
		case "false":
			or.PrivMap = or.PrivMap.Remove(role.ID, model.PrivAssignRole)
		}
		switch r.FormValue("t:manageEvents-" + orid) {
		case "true":
			or.PrivMap = or.PrivMap.Add(role.ID, model.PrivManageEvents)
		case "false":
			or.PrivMap = or.PrivMap.Remove(role.ID, model.PrivManageEvents)
		}
	}
}

// enforceRoleConstraints adjusts the privileges on all roles so that they
// conform to the role constraints.  It returns nil if the constraints cannot be
// met due to an uncorrectable error.  Otherwise, it returns a set of shadow
// privilege maps whose bits indicate which privileges are forced.  (Doing this
// here keeps that logic aligned with the constraint enforcement.)
func enforceRoleConstraints(roles []*model.Role) (forced map[*model.Role]model.PrivilegeMap) {
	// 0, Create the map.
	forced = make(map[*model.Role]model.PrivilegeMap, len(roles))
	var maxID model.RoleID
	for _, r := range roles {
		if r.ID > maxID {
			maxID = r.ID
		}
	}
	for _, r := range roles {
		forced[r] = make(model.PrivilegeMap, maxID+1)
	}
	// 1. All roles must hold themselves.
	for _, r := range roles {
		r.PrivMap = r.PrivMap.Add(r.ID, model.PrivHoldsRole)
		forced[r].Add(r.ID, model.PrivHoldsRole)
	}
	// 2. No cycles in the role graph.
	if !model.RecalculateTransitivePrivilegeMaps(roles) {
		return nil
	}
	for _, r := range roles {
		for _, or := range roles {
			if r != or && r.TransPrivs.Has(or.ID, model.PrivHoldsRole) {
				forced[or].Add(r.ID, model.PrivHoldsRole)
			}
		}
	}
	// 3. No roles may explicitly specify privileges that they also inherit
	// from other roles that they hold.
	for _, r := range roles {
		for _, imp := range roles {
			if r != imp && r.PrivMap.Has(imp.ID, model.PrivHoldsRole) {
				for _, or := range roles {
					disallowed := imp.TransPrivs.Get(or.ID) &^ model.PrivHoldsRole
					r.PrivMap = r.PrivMap.Remove(or.ID, disallowed)
					forced[r].Add(or.ID, disallowed)
				}
			}
		}
	}
	// 4. No roles may directly imply other roles that they also inherit by
	// implication of other roles.
	for _, r := range roles {
		for _, imp := range roles {
			if r != imp && r.PrivMap.Has(imp.ID, model.PrivHoldsRole) {
				for _, or := range roles {
					if imp != or && imp.TransPrivs.Has(or.ID, model.PrivHoldsRole) {
						r.PrivMap.Remove(or.ID, model.PrivHoldsRole)
						forced[r].Add(or.ID, model.PrivHoldsRole)
					}
				}
			}
		}
	}
	// 5. Webmaster has all privileges on all roles.
	var webmaster *model.Role
	for _, r := range roles {
		if r.Tag == model.RoleWebmaster {
			webmaster = r
			break
		}
	}
	if webmaster == nil {
		return nil
	}
	for _, r := range roles {
		webmaster.PrivMap = webmaster.PrivMap.Add(r.ID, model.PrivViewHolders|model.PrivAssignRole|model.PrivManageEvents)
		webmaster.TransPrivs = webmaster.TransPrivs.Add(r.ID, model.PrivViewHolders|model.PrivAssignRole|model.PrivManageEvents)
		forced[webmaster].Add(r.ID, model.PrivViewHolders|model.PrivAssignRole|model.PrivManageEvents)
	}
	// 6. No roles R0 may explicitly specify privileges on R1 that they have
	// specified on another role R2, directly or indirectly implied by R1.
	for _, r0 := range roles {
		for _, r1 := range roles {
			for _, r2 := range roles {
				if r1 != r2 && r1.TransPrivs.Has(r2.ID, model.PrivHoldsRole) {
					disallowed := r0.PrivMap.Get(r2.ID) &^ model.PrivHoldsRole
					r0.PrivMap = r0.PrivMap.Remove(r1.ID, disallowed)
					forced[r0].Add(r1.ID, disallowed)
				}
			}
		}
	}
	return forced
}
*/
