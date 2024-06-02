package personedit

import (
	"fmt"
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/recalc"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleRoles handles requests for /people/$id/edroles.
func HandleRoles(r *request.Request, idstr string) {
	var (
		user *person.Person
		p    *person.Person
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), personview.PersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(0, enum.PrivLeader) || p.ID() == person.AdminID {
		errpage.Forbidden(r, user)
		return
	}
	if r.Method == http.MethodPost {
		handlePostRoles(r, user, p)
	} else {
		handleGetRoles(r, user, p)
	}
}

func handleGetRoles(r *request.Request, user, p *person.Person) {
	var held = make(map[role.ID]bool)

	personrole.RolesForPerson(r, p.ID(), role.FID, func(r *role.Role, explicit bool) {
		if explicit {
			held[r.ID()] = true
		}
	})
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col personeditRoles' method=POST up-main up-layer=parent up-target=main")
	form.E("div class='formTitle formTitle-primary'>Edit Roles")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	for _, org := range enum.AllOrgs {
		if org != enum.OrgAdmin && (user.HasPrivLevel(org, enum.PrivLeader) || user.IsAdminLeader()) {
			handleGetOrgRoles(r, form, p, held, org)
		}
	}
	if user.IsAdminLeader() {
		handleGetOrgRoles(r, form, p, held, enum.OrgAdmin)
	}
	emitButtons(r, form)
}

var orgLabels = map[enum.Org]string{
	enum.OrgAdmin:  "Administrative Roles",
	enum.OrgCERTD:  "CERT Deployment Team Roles",
	enum.OrgCERTT:  "CERT Training Roles",
	enum.OrgListos: "Listos Roles",
	enum.OrgSARES:  "SARES Roles",
	enum.OrgSNAP:   "SNAP Roles",
}

func handleGetOrgRoles(r *request.Request, form *htmlb.Element, p *person.Person, held map[role.ID]bool, org enum.Org) {
	row := form.E("div class=formRow-3col")
	row.E("div class=personeditRolesOrg>%s", orgLabels[org])
	role.AllWithOrg(r, role.FID|role.FFlags|role.FImplies|role.FName, org, func(rl *role.Role) {
		var implies string

		if rl.Flags()&role.ImplicitOnly != 0 {
			return
		}
		if len(rl.Implies()) != 0 {
			var sb strings.Builder
			for _, imp := range rl.Implies() {
				fmt.Fprintf(&sb, ",%d", imp)
			}
			implies = sb.String()[1:]
		}
		row.E("div").E("input type=checkbox class=s-check name=role value=%d label=%s", rl.ID(), rl.Name(),
			held[rl.ID()], "checked", implies != "", "data-implies=%s", implies)
	})
}

func handlePostRoles(r *request.Request, user, p *person.Person) {
	var held = make(map[role.ID]bool)

	for _, v := range r.Form["role"] {
		held[role.ID(util.ParseID(v))] = true
	}
	r.Transaction(func() {
		for _, org := range enum.AllOrgs {
			if org != enum.OrgAdmin && (user.HasPrivLevel(org, enum.PrivLeader) || user.IsAdminLeader()) {
				handlePostOrgRoles(r, p, held, org)
			}
		}
		if user.IsAdminLeader() {
			handlePostOrgRoles(r, p, held, enum.OrgAdmin)
		}
		recalc.Recalculate(r)
	})
	personview.Render(r, user, p, person.ViewFull, "")
}

func handlePostOrgRoles(r *request.Request, p *person.Person, held map[role.ID]bool, org enum.Org) {
	role.AllWithOrg(r, role.FID|role.FName, org, func(rl *role.Role) {
		if held[rl.ID()] {
			personrole.AddRole(r, p, rl)
		} else {
			personrole.RemoveRole(r, p, rl)
		}
	})
}
