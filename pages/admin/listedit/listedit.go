package listedit

import (
	"net/http"
	"regexp"
	"slices"
	"sort"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/listlist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listrole"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type roleData struct {
	rl       *role.Role
	sender   bool
	submodel listrole.SubscriptionModel
}

// Handle handles /admin/lists/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user       *person.Person
		l          *list.List
		ul         *list.Updater
		roles      []*roleData
		canDelete  bool
		nameError  string
		typeError  string
		rolesError string
		hasError   bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if idstr == "NEW" {
		ul = new(list.Updater)
	} else {
		if l = list.WithID(r, list.ID(util.ParseID(idstr))); l == nil {
			errpage.NotFound(r, user)
			return
		}
		ul = l.Updater()
		canDelete = true
		listrole.AllRolesForList(r, ul.ID, role.FID|role.FName, func(rl *role.Role, sender bool, submodel listrole.SubscriptionModel) {
			roles = append(roles, &roleData{rl.Clone(), sender, submodel})
		})
	}
	if r.Method == http.MethodPost {
		if canDelete && r.FormValue("delete") != "" {
			deleteList(r, l)
			listlist.Render(r, user)
			return
		}
		typeError = readType(r, ul)
		nameError = readName(r, ul)
		roles, rolesError = readRoles(r, roles)
		hasError = nameError != "" || typeError != "" || rolesError != "" || r.Request.Header.Get("X-Up-Validate") != ""
		if !hasError {
			saveList(r, l, ul, roles)
			listlist.Render(r, user)
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
	if l == nil {
		form.E("div class='formTitle formTitle-primary'>New List")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit List")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitType(form, ul, typeError)
	emitName(form, ul, nameError)
	emitRoles(form, ul, roles, rolesError)
	emitButtons(form, canDelete)
}

func emitType(form *htmlb.Element, ul *list.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=listeditEmail>Type")
	box := row.E("div class='formInput listeditType'")
	box.E("s-radio id=listeditEmail name=type value=email label=Email", ul.Type == list.Email, "checked")
	box.E("s-radio id=listeditSMS name=type value=sms label=SMS", ul.Type == list.SMS, "checked")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
func readType(r *request.Request, ul *list.Updater) string {
	switch r.FormValue("type") {
	case "email":
		ul.Type = list.Email
	case "sms":
		ul.Type = list.SMS
	case "":
		return "The list type is required."
	default:
		return "The list type is not valid."
	}
	return ""
}

var emailNameRE = regexp.MustCompile(`^[a-z][-a-z0-9]*$`)

func emitName(form *htmlb.Element, ul *list.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=listeditName>Name")
	row.E("input id=listeditName name=name class=formInput value=%s", ul.Name)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
func readName(r *request.Request, ul *list.Updater) string {
	ul.Name = strings.TrimSpace(r.FormValue("name"))
	if ul.Name == "" {
		return "The list name is required."
	} else if ul.DuplicateName(r) {
		return "Another list has this name."
	} else if ul.Type == list.Email && !emailNameRE.MatchString(ul.Name) {
		return "The list name is not valid as the first part of an @sunnyvaleserv.org email address."
	}
	return ""
}

func emitRoles(form *htmlb.Element, ul *list.Updater, roles []*roleData, err string) {
	row := form.E("div class=formRow")
	row.E("label>Roles")
	box := row.E("div class=formInput")
	for _, rd := range roles {
		if rd.submodel == 0 && !rd.sender {
			continue
		}
		div := box.E("div")
		div.E("a href=# class=listeditRoleEdit data-list=%d data-role=%d>%s", ul.ID, rd.rl.ID(), rd.rl.Name())
		div.R(" (")
		if rd.submodel != 0 {
			div.R(rd.submodel.String())
		}
		if rd.submodel != 0 && rd.sender {
			div.R(", ")
		}
		if rd.sender {
			div.R("sender")
		}
		div.R(")")
		div.E("input type=hidden id=listeditRole%d name=role%d value=%d:%v", rd.rl.ID(), rd.rl.ID(), rd.submodel, rd.sender)
	}
	box.E("a href=# class='sbtn sbtn-small sbtn-primary listeditRoleEdit' data-list=%d>Add", ul.ID)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
func readRoles(r *request.Request, oroles []*roleData) (roles []*roleData, err string) {
	for key := range r.Form {
		var rd roleData
		if !strings.HasPrefix(key, "role") {
			continue
		}
		if rd.rl = role.WithID(r, role.ID(util.ParseID(key[4:])), role.FID|role.FName); rd.rl == nil {
			continue
		}
		parts := strings.Split(r.FormValue(key), ":")
		if len(parts) != 2 {
			continue
		}
		if rd.submodel = listrole.SubscriptionModel(util.ParseID(parts[0])); !rd.submodel.Valid() {
			continue
		}
		rd.sender = parts[1] == "true"
		roles = append(roles, &rd)
	}
	if len(roles) == 0 {
		return nil, "At least one role must have privileges."
	}
	for _, or := range oroles {
		if !slices.ContainsFunc(roles, func(rd *roleData) bool { return rd.rl.ID() == or.rl.ID() }) {
			or.submodel, or.sender = 0, false
			roles = append(roles, or)
		}
	}
	sort.Slice(roles, func(i, j int) bool { return roles[i].rl.Name() < roles[j].rl.Name() })
	return roles, ""
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

func saveList(r *request.Request, l *list.List, ul *list.Updater, roles []*roleData) {
	r.Transaction(func() {
		if l == nil {
			l = list.Create(r, ul)
		} else {
			l.Update(r, ul)
		}
		for _, rd := range roles {
			listrole.SetListRole(r, l, rd.rl, rd.sender, rd.submodel)
		}
		role.Recalculate(r)
	})
}

func deleteList(r *request.Request, l *list.List) {
	r.Transaction(func() {
		l.Delete(r)
		role.Recalculate(r)
	})
}
