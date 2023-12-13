package personedit

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const namesPersonFields = person.FInformalName | person.FFormalName | person.FSortName | person.FCallSign | person.FPronouns | person.FBirthdate

// HandleNames handles requests for /people/$id/ednames.
func HandleNames(r *request.Request, idstr string) {
	var (
		user              *person.Person
		p                 *person.Person
		up                *person.Updater
		informalNameError string
		formalNameError   string
		sortNameError     string
		callSignError     string
		birthdateError    string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), namesPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() != p.ID() && !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	up = p.Updater()
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		informalNameError = readInformalName(r, up)
		formalNameError = readFormalName(r, up)
		sortNameError = readSortName(r, up)
		callSignError = readCallSign(r, up)
		birthdateError = readBirthdate(r, up)
		readPronouns(r, up)
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && informalNameError == "" && formalNameError == "" &&
			sortNameError == "" && callSignError == "" && birthdateError == "" {
			r.Transaction(func() {
				p.Update(r, up, namesPersonFields)
			})
			personview.Render(r, user, p, false, "names")
			return
		}
	}
	r.HTMLNoCache()
	if informalNameError != "" || formalNameError != "" || sortNameError != "" || callSignError != "" || birthdateError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewNames")
	form.E("div class='formTitle formTitle-primary'>Edit Names")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "informalName") {
		emitInformalName(form, up, informalNameError != "" || (formalNameError == "" && sortNameError == "" && callSignError == ""), informalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "formalName") {
		emitFormalName(form, up, formalNameError != "", formalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "sortName") {
		emitSortName(form, up, sortNameError != "", sortNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "callSign") {
		emitCallSign(form, up, callSignError != "", callSignError)
	}
	if len(validate) == 0 || slices.Contains(validate, "birthdate") {
		emitBirthdate(form, up, birthdateError != "", birthdateError)
	}
	if len(validate) == 0 {
		emitPronouns(form, up)
		emitButtons(form)
	}
}

func readInformalName(r *request.Request, up *person.Updater) string {
	if up.InformalName = strings.TrimSpace(r.FormValue("informalName")); up.InformalName == "" {
		return "The name is required."
	}
	return ""
}

func emitInformalName(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditInformalName>Name")
	row.E("input id=personeditInformalName name=informalName s-validate value=%s", up.InformalName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>What you like to be called, e.g. “Joe Banks”")
}

func readFormalName(r *request.Request, up *person.Updater) string {
	if up.FormalName = strings.TrimSpace(r.FormValue("formalName")); up.FormalName == "" {
		return "The formal name is required."
	}
	return ""
}

func emitFormalName(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditFormalName>Formal name")
	row.E("input id=personeditFormalName name=formalName s-validate value=%s", up.FormalName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>For formal documents, e.g. “Joseph A. Banks, Jr.”")
}

func readSortName(r *request.Request, up *person.Updater) string {
	if up.SortName = strings.TrimSpace(r.FormValue("sortName")); up.SortName == "" {
		return "The sort name is required."
	} else if up.DuplicateSortName(r) {
		return fmt.Sprintf("Another person has the sort name %q.", up.SortName)
	}
	return ""
}

func emitSortName(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditSortName>Sort name")
	row.E("input id=personeditSortName name=sortName s-validate value=%s", up.SortName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>For appearance in sorted lists, e.g. “Banks, Joe”")
}

var callsignRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)

func readCallSign(r *request.Request, up *person.Updater) string {
	if up.CallSign = strings.ToUpper(strings.TrimSpace(r.FormValue("callSign"))); up.CallSign == "" {
		return ""
	}
	if !callsignRE.MatchString(up.CallSign) {
		return fmt.Sprintf("%q is not a valid FCC amateur radio call sign.", up.CallSign)
	}
	if up.DuplicateCallSign(r) {
		return fmt.Sprintf("Another person has the call sign %q.", up.CallSign)
	}
	return ""
}

func emitCallSign(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditCallSign>Call sign")
	row.E("input id=personeditCallSign name=callSign s-validate value=%s", up.CallSign, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp>FCC amateur radio license (optional)")
}

func readBirthdate(r *request.Request, up *person.Updater) string {
	if up.Birthdate = strings.ToUpper(strings.TrimSpace(r.FormValue("birthdate"))); up.Birthdate == "" {
		return ""
	}
	if d, err := time.Parse("2006-01-02", up.Birthdate); err != nil || d.Format("2006-01-02") != up.Birthdate {
		return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", up.Birthdate)
	}
	return ""
}

func emitBirthdate(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditBirthdate>Birthdate")
	row.E("input type=date id=personeditBirthdate name=birthdate s-validate value=%s", up.Birthdate, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readPronouns(r *request.Request, up *person.Updater) {
	up.Pronouns = strings.TrimSpace(r.FormValue("pronouns"))
}

func emitPronouns(form *htmlb.Element, up *person.Updater) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditPronouns>Pronouns")
	row.E("input id=personeditPronouns name=pronouns list=personeditPronounChoices value=%s", up.Pronouns)
	row.E("datalist id=personeditPronounChoices").
		E("option value='he/him/his'").P().
		E("option value='she/her/hers'").P().
		E("option value='they/them/theirs'")
}

func emitButtons(form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
}
