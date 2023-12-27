package classedit

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/admin/classlist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Handle handles /admin/classes/$id requests, where $id may be "NEW".
func Handle(r *request.Request, idstr string) {
	var (
		user        *person.Person
		c           *class.Class
		uc          *class.Updater
		canDelete   bool
		typeError   string
		startError  string
		enDescError string
		esDescError string
		limitError  string
		hasError    bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if idstr == "NEW" {
		uc = new(class.Updater)
	} else {
		if c = class.WithID(r, class.ID(util.ParseID(idstr)), class.UpdaterFields); c == nil {
			errpage.NotFound(r, user)
			return
		}
		uc = c.Updater()
		canDelete = false // TODO
	}
	if r.Method == http.MethodPost {
		if canDelete && r.FormValue("delete") != "" {
			deleteClass(r, c)
			classlist.Render(r, user)
			return
		}
		typeError = readType(r, uc)
		startError = readStart(r, uc)
		enDescError = readEnDesc(r, uc)
		esDescError = readEsDesc(r, uc)
		limitError = readLimit(r, uc)
		hasError = typeError != "" || startError != "" || enDescError != "" || esDescError != "" || limitError != ""
		if !hasError {
			saveClass(r, c, uc)
			classlist.Render(r, user)
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
	if c == nil {
		form.E("div class='formTitle formTitle-primary'>New Class")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit Class")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitType(form, uc, typeError)
	emitStart(form, uc, startError)
	emitEnDesc(form, uc, enDescError)
	emitEsDesc(form, uc, esDescError)
	emitLimit(form, uc, limitError)
	emitReferrals(form, uc)
	emitButtons(form, canDelete)
}

func emitType(form *htmlb.Element, uc *class.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=classeditType>Type")
	sel := row.E("select id=classeditType name=type class=formInput")
	if uc.Type == 0 {
		sel.E("option value=''>(select type)")
	}
	for _, t := range class.AllTypes {
		sel.E("option value=%d", t, uc.Type == t, "selected").T(t.String())
	}
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
func readType(r *request.Request, uc *class.Updater) string {
	uc.Type = class.Type(util.ParseID(r.FormValue("type")))
	if uc.Type == 0 {
		return "The class type is required."
	}
	for _, t := range class.AllTypes {
		if uc.Type == t {
			return ""
		}
	}
	uc.Type = 0
	return "The selected type is not valid."
}

func emitStart(form *htmlb.Element, uc *class.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=classeditStart>Start Date")
	row.E("input type=date id=classeditStart name=start class=formInput value=%s", uc.Start)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Use 2999-12-31 for a waiting list placeholder.")
}
func readStart(r *request.Request, uc *class.Updater) string {
	uc.Start = r.FormValue("start")
	if uc.Start == "" {
		return "The class starting date is required."
	} else if d, err := time.Parse("2006-01-02", uc.Start); err != nil || d.Format("2006-01-02") != uc.Start {
		return "The date is not a valid YYYY-MM-DD date."
	} else if uc.DuplicateStart(r) {
		return "Another class has the same type and start date."
	} else {
		return ""
	}
}

func emitEnDesc(form *htmlb.Element, uc *class.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=classeditEnDesc>English Desc.")
	row.E("textarea id=classeditEnDesc name=enDesc class=formInput").T(uc.EnDesc)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>Include date(s), time(s), location(s), maybe language")
}
func readEnDesc(r *request.Request, uc *class.Updater) string {
	if uc.EnDesc = strings.TrimSpace(r.FormValue("enDesc")); uc.EnDesc == "" {
		return "The English description is required."
	} else {
		return ""
	}
}

func emitEsDesc(form *htmlb.Element, uc *class.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=classeditEsDesc>Spanish Desc.")
	row.E("textarea id=classeditEsDesc name=esDesc class=formInput").T(uc.EsDesc)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
func readEsDesc(r *request.Request, uc *class.Updater) string {
	if uc.EsDesc = strings.TrimSpace(r.FormValue("esDesc")); uc.EsDesc == "" {
		return "The Spanish description is required."
	} else {
		return ""
	}
}

func emitLimit(form *htmlb.Element, uc *class.Updater, err string) {
	row := form.E("div class=formRow")
	row.E("label for=classeditLimit>Enrollment Limit")
	row.E("input type=number id=classeditLimit name=limit class=formInput value=%d min=0", uc.Limit)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>0 = unlimited")
}
func readLimit(r *request.Request, uc *class.Updater) string {
	if lstr := r.FormValue("limit"); lstr == "" {
		uc.Limit = 0
		return ""
	} else if lint, err := strconv.Atoi(lstr); err == nil && lint >= 0 {
		uc.Limit = uint(lint)
		return ""
	} else {
		return "The specified limit is not a valid integer."
	}
}

func emitReferrals(form *htmlb.Element, uc *class.Updater) {
	if uc.Referrals == nil {
		return
	}
	row := form.E("div class=formRow")
	row.E("label>Referrals")
	grid := row.E("div class='classeditReferrals formInput'")
	for _, ref := range class.AllReferrals {
		grid.E("div>%d", uc.Referrals[ref])
		grid.E("div>%s", ref.String())
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

func saveClass(r *request.Request, c *class.Class, ur *class.Updater) {
	r.Transaction(func() {
		if c == nil {
			c = class.Create(r, ur)
		} else {
			c.Update(r, ur)
		}
	})
}

func deleteClass(r *request.Request, c *class.Class) {
	r.Transaction(func() {
		c.Delete(r)
	})
}
