package roleedit

import (
	"fmt"
	"net/http"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/rolelist"
	"sunnyvaleserv.org/portal/pages/admin/roleselect"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const priorityBeforeUnset = 99999

// Handle handles /admin/roles/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user       *person.Person
		rl         *role.Role
		ur         *role.Updater
		canDelete  bool
		prioBefore uint = priorityBeforeUnset
		nameError  string
		orgError   string
		privError  string
		hasError   bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if idstr == "NEW" {
		ur = new(role.Updater)
	} else {
		if rl = role.WithID(r, role.ID(util.ParseID(idstr)), role.UpdaterFields); rl == nil {
			errpage.NotFound(r, user)
			return
		}
		ur = rl.Updater()
		canDelete = personrole.PeopleCountForRole(r, rl.ID()) == 0
	}
	if r.Method == http.MethodPost {
		if canDelete && r.FormValue("delete") != "" {
			deleteRole(r, rl)
			rolelist.Render(r, user)
			return
		}
		nameError = readName(r, ur)
		readTitle(r, ur)
		orgError = readOrg(r, ur)
		privError = readPrivLevel(r, ur)
		prioBefore = readPriority(r, ur)
		readFlags(r, ur)
		readImplies(r, ur)
		hasError = nameError != "" || orgError != "" || privError != ""
		if !hasError {
			saveRole(r, rl, ur, prioBefore)
			rolelist.Render(r, user)
			return
		}
	}
	r.HTMLNoCache()
	if hasError {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	if rl == nil {
		form.E("div class='formTitle formTitle-primary'>New Role")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit Role")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitName(form, ur, nameError)
	emitTitle(form, ur)
	emitOrg(form, ur, orgError)
	emitPrivLevel(form, ur, privError)
	emitPriority(r, form, ur, prioBefore)
	emitFlags(form, ur)
	emitImplies(r, form, ur)
	emitLists(r, form, ur)
	emitButtons(form, canDelete)
}

func emitName(form *htmlb.Element, ur *role.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditName>Name")
	row.E("input id=roleeditName name=name class=formInput value=%s", ur.Name)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Collective name for people who hold this role")
}
func readName(r *request.Request, ur *role.Updater) string {
	ur.Name = strings.TrimSpace(r.FormValue("name"))
	if ur.Name == "" {
		return "The role name is required."
	} else if ur.DuplicateName(r) {
		return "Another role has this name."
	}
	return ""
}

func emitTitle(form *htmlb.Element, ur *role.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditTitle>Title")
	row.E("input id=roleeditTitle name=title class=formInput value=%s", ur.Title,
		ur.Title == ur.Name, "data-match=true")
	row.E("div class=formHelp>Title for a single person who holds this role, or empty if this role should not be called out in people lists")
}
func readTitle(r *request.Request, ur *role.Updater) {
	ur.Title = strings.TrimSpace(r.FormValue("title"))
}

func emitOrg(form *htmlb.Element, ur *role.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditOrg1>Organization")
	box := row.E("div class='formInput roleeditOrg'")
	for i := enum.Org(1); i < enum.NumOrgs; i++ {
		box.E("s-radio id=roleeditOrg%d name=org value=%d label=%s", i, i, i.String(), ur.Org == i, "checked")
	}
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Organization for which this role grants privileges")
}
func readOrg(r *request.Request, ur *role.Updater) string {
	orgstr := r.FormValue("org")
	if orgstr == "" {
		return "The organization is required."
	}
	ur.Org = enum.Org(util.ParseID(orgstr))
	if !ur.Org.Valid() {
		return "The organization is not valid."
	}
	return ""
}

func emitPrivLevel(form *htmlb.Element, ur *role.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditPriv0>Privilege")
	box := row.E("div class='formInput roleeditPriv'")
	for i := enum.PrivLevel(0); i < enum.PrivMaster; i++ {
		name := i.String()
		if name == "" {
			name = "(none)"
		}
		box.E("s-radio id=roleeditPriv%d name=priv value=%d label=%s", i, i, name, ur.PrivLevel == i, "checked")
	}
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Privilege level granted to people holding this role")
}
func readPrivLevel(r *request.Request, ur *role.Updater) string {
	switch privstr := r.FormValue("priv"); privstr {
	case "0", "1", "2", "3":
		ur.PrivLevel = enum.PrivLevel(util.ParseID(privstr))
		return ""
	case "":
		return "The privilege level is required."
	default:
		return "The privilege level is not valid."
	}
}

func emitPriority(r *request.Request, form *htmlb.Element, ur *role.Updater, prioBefore uint) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditPriority>Sort before")
	sel := row.E("select id=roleeditPriority name=priority class=formInput")
	shownSelected := false
	role.All(r, role.FID|role.FName|role.FPriority, func(rl *role.Role) {
		if rl.ID() != ur.ID {
			var selected bool
			if prioBefore == priorityBeforeUnset {
				selected = rl.Priority() > ur.Priority && !shownSelected
			} else {
				selected = rl.Priority() == prioBefore
			}
			sel.E("option value=%d", rl.Priority(), selected, "selected").T(rl.Name())
			shownSelected = shownSelected || selected
		}
	})
	sel.E("option value=0", !shownSelected, "selected").R("(at end)")
}
func readPriority(r *request.Request, ur *role.Updater) (prioBefore uint) {
	return uint(util.ParseID(r.FormValue("priority")))
}

func emitFlags(form *htmlb.Element, ur *role.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=roleeditFilter>Flags")
	box := row.E("div class='formInput roleeditFlags'")
	box.E("input type=checkbox id=roleeditFilter name=filter class=s-check label='Available choice on People list page'",
		ur.Flags&role.Filter != 0, "checked")
	box.E("input type=checkbox name=implicitOnly class=s-check label='Role can only be implied, not assigned'",
		ur.Flags&role.ImplicitOnly != 0, "checked")
	box.E("input type=checkbox name=archived class=s-check label='Archived (no longer in use)'",
		ur.Flags&role.Archived != 0, "checked")
}
func readFlags(r *request.Request, ur *role.Updater) {
	ur.Flags = 0
	if r.FormValue("filter") != "" {
		ur.Flags |= role.Filter
	}
	if r.FormValue("implicitOnly") != "" {
		ur.Flags |= role.ImplicitOnly
	}
	if r.FormValue("archived") != "" {
		ur.Flags |= role.Archived
	}
}

func emitImplies(r *request.Request, form *htmlb.Element, ur *role.Updater) {
	const roleFields = role.FID | role.FName
	var (
		impliesList string
		treedesc    string
		impliers    map[role.ID]struct{}
	)
	if len(ur.Implies) != 0 {
		var sb strings.Builder
		for _, imp := range ur.Implies {
			fmt.Fprintf(&sb, "%d ", imp)
		}
		impliesList = strings.TrimSpace(sb.String())
	}
	if ur.ID != 0 {
		impliers = role.AllThatImply(r, ur.ID)
	}
	treedesc = roleselect.MakeRoleTree(r, roleFields, func(rl *role.Role) bool {
		_, ok := impliers[rl.ID()]
		return !ok
	})
	row := form.E("div class=formRow")
	row.E("label>Implies")
	row.E("s-seltree name=implies class=formInput value=%s", impliesList).R(treedesc)
}
func readImplies(r *request.Request, ur *role.Updater) {
	idstrs := strings.Fields(r.FormValue("implies"))
	ur.Implies = ur.Implies[:0]
	for _, idstr := range idstrs {
		if rid := role.ID(util.ParseID(idstr)); rid > 0 {
			ur.Implies = append(ur.Implies, rid)
		}
	}
}

func emitLists(r *request.Request, form *htmlb.Element, ur *role.Updater) {
	if ur.ID == 0 {
		return
	}
	var found bool

	row := form.E("div class=formRow")
	row.E("label>Lists")
	box := row.E("div class=formInput")
	listrole.AllListsForRole(r, ur.ID, func(l *list.List, sender bool, submodel listrole.SubscriptionModel) {
		found = true
		var name string
		if l.Type == list.SMS {
			name = "SMS: " + l.Name
		} else {
			name = l.Name + "@sunnyvaleserv.org"
		}
		div := box.E("div")
		div.E("a href=/admin/lists/%d up-target=main>%s", l.ID, name)
		div.R(" (")
		if submodel != 0 {
			div.R(submodel.String())
		}
		if submodel != 0 && sender {
			div.R(", ")
		}
		if sender {
			div.R("sender")
		}
		div.R(")")
	})
	if !found {
		box.R("None")
	}
}

func emitButtons(form *htmlb.Element, canDelete bool) {
	buttons := form.E("div class=formButtons")
	buttons.E("div class=formButtonSpace")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	// This button must appear lexically after Save, even though it appears
	// visually before it, so that Save is the default button when the user
	// presses Enter.  The formButton-beforeAll class implements that.
	if canDelete {
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
}

func saveRole(r *request.Request, rl *role.Role, ur *role.Updater, prioBefore uint) {
	r.Transaction(func() {
		if rl == nil {
			rl = role.Create(r, ur)
		} else {
			rl.Update(r, ur)
		}
		switch {
		case prioBefore == 0:
			rl.Reorder(r, 0)
		case prioBefore < rl.Priority():
			rl.Reorder(r, prioBefore)
		case prioBefore > rl.Priority()+1:
			rl.Reorder(r, prioBefore-1)
		}
		role.Recalculate(r)
	})
}

func deleteRole(r *request.Request, rl *role.Role) {
	r.Transaction(func() {
		rl.Delete(r)
		role.Recalculate(r)
	})
}
