package peoplelist

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleNewUser handles requests for /people/newuser
func HandleNewUser(r *request.Request) {
	var (
		user              *person.Person
		p                 *person.Person
		up                person.Updater
		informalNameError string
		formalNameError   string
		sortNameError     string
		callSignError     string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		informalWas := up.InformalName
		informalNameError = readInformalName(r, &up)
		formalNameError = readFormalName(r, &up)
		sortNameError = readSortName(r, &up)
		callSignError = readCallSign(r, &up)
		readPronouns(r, &up)
		if slices.Contains(validate, "informalName") {
			if up.FormalName == "" || up.FormalName == informalWas {
				up.FormalName = up.InformalName
			}
			if up.SortName == "" || up.SortName == sortName(informalWas) {
				up.SortName = sortName(up.InformalName)
			}
		}
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && informalNameError == "" && formalNameError == "" &&
			sortNameError == "" && callSignError == "" {
			up.UnsubscribeToken = util.RandomToken()
			r.Transaction(func() {
				p = person.Create(r, &up)
			})
			http.Redirect(r, r.Request, fmt.Sprintf("/people/%d", p.ID()), http.StatusSeeOther)
			return
		}
	}
	r.HTMLNoCache()
	if informalNameError != "" || formalNameError != "" || sortNameError != "" || callSignError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Create User"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "informalName") {
		emitInformalName(r, form, &up, informalNameError != "" || (formalNameError == "" && sortNameError == "" && callSignError == ""), informalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "informalName") || slices.Contains(validate, "formalName") {
		emitFormalName(r, form, &up, formalNameError != "", formalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "informalName") || slices.Contains(validate, "sortName") {
		emitSortName(r, form, &up, sortNameError != "", sortNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "callSign") {
		emitCallSign(r, form, &up, callSignError != "", callSignError)
	}
	if len(validate) == 0 {
		emitPronouns(r, form, &up)
		emitButtons(r, form)
	}
}

func readInformalName(r *request.Request, up *person.Updater) string {
	if up.InformalName = strings.TrimSpace(r.FormValue("informalName")); up.InformalName == "" {
		return r.Loc("The name is required.")
	}
	return ""
}

func emitInformalName(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div id=personeditInformalNameRow class=formRow")
	row.E("label for=personeditInformalName").R(r.Loc("Name"))
	row.E("input id=personeditInformalName name=informalName s-validate=#personeditInformalName,#personeditFormalName,#personeditSortName value=%s", up.InformalName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp").R(r.Loc("What you like to be called, e.g. “Joe Banks”"))
}

func readFormalName(r *request.Request, up *person.Updater) string {
	if up.FormalName = strings.TrimSpace(r.FormValue("formalName")); up.FormalName == "" {
		return r.Loc("The formal name is required.")
	}
	return ""
}

func emitFormalName(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div id=personeditFormalNameRow class=formRow")
	row.E("label for=personeditFormalName").R(r.Loc("Formal name"))
	row.E("input id=personeditFormalName name=formalName s-validate value=%s", up.FormalName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp").R(r.Loc("For formal documents, e.g. “Joseph A. Banks, Jr.”"))
}

func readSortName(r *request.Request, up *person.Updater) string {
	if up.SortName = strings.TrimSpace(r.FormValue("sortName")); up.SortName == "" {
		return r.Loc("The sort name is required.")
	} else if up.DuplicateSortName(r) {
		return fmt.Sprintf(r.Loc("Another person has the sort name %q."), up.SortName)
	}
	return ""
}

func emitSortName(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div id=personeditSortNameRow class=formRow")
	row.E("label for=personeditSortName").R(r.Loc("Sort name"))
	row.E("input id=personeditSortName name=sortName s-validate value=%s", up.SortName, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp").R(r.Loc("For appearance in sorted lists, e.g. “Banks, Joe”"))
}

func sortName(s string) string {
	fields := strings.Fields(s)
	if len(fields) < 2 {
		return s
	}
	return fields[len(fields)-1] + ", " + strings.Join(fields[:len(fields)-1], " ")
}

var callsignRE = regexp.MustCompile(`^[AKNW][A-Z]?[0-9][A-Z]{1,3}$`)

func readCallSign(r *request.Request, up *person.Updater) string {
	if up.CallSign = strings.ToUpper(strings.TrimSpace(r.FormValue("callSign"))); up.CallSign == "" {
		return ""
	}
	if !callsignRE.MatchString(up.CallSign) {
		return fmt.Sprintf(r.Loc("%q is not a valid FCC amateur radio call sign."), up.CallSign)
	}
	if up.DuplicateCallSign(r) {
		return fmt.Sprintf(r.Loc("Another person has the call sign %q."), up.CallSign)
	}
	return ""
}

func emitCallSign(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div id=personeditCallSignRow class=formRow")
	row.E("label for=personeditCallSign").R(r.Loc("Call sign"))
	row.E("input id=personeditCallSign name=callSign s-validate value=%s", up.CallSign, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp").R(r.Loc("FCC amateur radio license (if any)"))
}

func readPronouns(r *request.Request, up *person.Updater) {
	up.Pronouns = strings.TrimSpace(r.FormValue("pronouns"))
}

func emitPronouns(r *request.Request, form *htmlb.Element, up *person.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=personeditPronouns").R(r.Loc("Pronouns"))
	row.E("input id=personeditPronouns name=pronouns list=personeditPronounChoices value=%s", up.Pronouns)
	list := row.E("datalist id=personeditPronounChoices")
	list.E("option value=%s", r.Loc("he/him/his"))
	list.E("option value=%s", r.Loc("she/her/hers"))
	list.E("option value=%s", r.Loc("they/them/theirs"))
}

func emitButtons(r *request.Request, form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss").R(r.Loc("Cancel"))
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Save"))
}
