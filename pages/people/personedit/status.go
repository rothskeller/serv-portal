package personedit

import (
	"fmt"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const (
	setStatusPersonFields = person.FVolgisticsID | person.FIdentification | person.FBGChecks | person.FDSWRegistrations | person.FFlags
	getStatusPersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | setStatusPersonFields
)

// HandleStatus handles requests for /people/$id/edstatus.
func HandleStatus(r *request.Request, idstr string) {
	var (
		user            *person.Person
		p               *person.Person
		up              *person.Updater
		volgisticsError string
		dswCERTError    string
		dswCommError    string
		bgDOJError      string
		bgFBIError      string
		bgPHSError      string
		haveErrors      bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !user.IsAdminLeader() {
		errpage.Forbidden(r, user)
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), getStatusPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	up = p.Updater()
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		volgisticsError = readVolgistics(r, up)
		dswCERTError = readDSWCERT(r, up)
		dswCommError = readDSWComm(r, up)
		bgDOJError = readBGDOJ(r, up)
		bgFBIError = readBGFBI(r, up)
		bgPHSError = readBGPHS(r, up)
		readIdents(r, up)
		haveErrors = volgisticsError != "" || dswCERTError != "" || dswCommError != "" || bgDOJError != "" || bgFBIError != "" || bgPHSError != ""
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !haveErrors {
			r.Transaction(func() {
				p.Update(r, up, setStatusPersonFields)
			})
			personview.Render(r, user, p, person.ViewFull, "status")
			return
		}
	}
	r.HTMLNoCache()
	if haveErrors {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewStatus")
	form.E("div class='formTitle formTitle-primary'>Edit Volunteer Status")
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "volgistics") {
		emitVolgistics(form, up, volgisticsError != "" || !haveErrors, volgisticsError)
	}
	emitIdents(form, up)
	if len(validate) == 0 || slices.Contains(validate, "dswCERTReg") || slices.Contains(validate, "dswCERTExp") {
		emitDSWCERT(form, up, dswCERTError != "", dswCERTError)
	}
	if len(validate) == 0 || slices.Contains(validate, "dswCommReg") || slices.Contains(validate, "dswCommExp") {
		emitDSWComm(form, up, dswCommError != "", dswCommError)
	}
	if len(validate) == 0 || slices.Contains(validate, "bgDOJCleared") || slices.Contains(validate, "bgDOJNLI") {
		emitBGDOJ(form, up, bgDOJError != "", bgDOJError)
	}
	if len(validate) == 0 || slices.Contains(validate, "bgFBICleared") || slices.Contains(validate, "bgFBINLI") {
		emitBGFBI(form, up, bgFBIError != "", bgFBIError)
	}
	if len(validate) == 0 || slices.Contains(validate, "bgPHSCleared") || slices.Contains(validate, "bgPHSNLI") {
		emitBGPHS(form, up, bgPHSError != "", bgPHSError)
	}
	if len(validate) == 0 {
		emitButtons(r, form)
	}
}

func readVolgistics(r *request.Request, up *person.Updater) string {
	if vstr := r.FormValue("volgistics"); vstr != "" {
		if v, err := strconv.Atoi(vstr); err == nil && v > 0 {
			up.VolgisticsID = uint(v)
			if up.DuplicateVolgisticsID(r) {
				return fmt.Sprintf("Volgistics ID %d is in use by a different person.", v)
			}
		} else {
			return fmt.Sprintf("%q is not a valid Volgistics ID (a positive integer).", vstr)
		}
	} else {
		up.VolgisticsID = 0
	}
	return ""
}

func emitVolgistics(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow personeditVolgistics'")
	row.E("label for=personeditVolgistics>Volgistics ID")
	row.E("input id=personeditVolgistics name=volgistics s-validate=.personeditVolgistics",
		up.VolgisticsID != 0, "value=%d", up.VolgisticsID,
		focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readIdents(r *request.Request, up *person.Updater) {
	up.Identification = 0
	for _, v := range r.Form["ident"] {
		switch v {
		case "photo":
			up.Identification |= person.IDPhoto
		case "cardKey":
			up.Identification |= person.IDCardKey
		case "servShirt":
			up.Identification |= person.IDSERVShirt
		case "certShirtLS":
			up.Identification |= person.IDCERTShirtLS
		case "certShirtSS":
			up.Identification |= person.IDCERTShirtSS
		}
	}
}

func emitIdents(form *htmlb.Element, up *person.Updater) {
	row := form.E("div class=formRow")
	row.E("label for=personeditIDPhoto>IDs Issued")
	in := row.E("div class=formInput")
	in.E("div").E("input type=checkbox class=s-check id=personeditIDPhoto name=ident value=photo label='photo ID'",
		up.Identification&person.IDPhoto != 0, "checked")
	in.E("div").E("input type=checkbox class=s-check name=ident value=cardKey label='access card'",
		up.Identification&person.IDCardKey != 0, "checked")
	in.E("div").E("input type=checkbox class=s-check name=ident value=servShirt label='tan SERV shirt'",
		up.Identification&person.IDSERVShirt != 0, "checked")
	in.E("div").E("input type=checkbox class=s-check name=ident value=certShirtLS label='green CERT shirt (LS)'",
		up.Identification&person.IDCERTShirtLS != 0, "checked")
	in.E("div").E("input type=checkbox class=s-check name=ident value=certShirtSS label='green CERT shirt (SS)'",
		up.Identification&person.IDCERTShirtSS != 0, "checked")
}

func readDSWCERT(r *request.Request, up *person.Updater) string {
	return readDSW(r, &up.DSWRegistrations.CERT, "CERT")
}

func readDSWComm(r *request.Request, up *person.Updater) string {
	return readDSW(r, &up.DSWRegistrations.Communications, "Comm")
}

func readDSW(r *request.Request, dsw **person.DSWRegistration, name string) string {
	var (
		dstr string
		date time.Time
		err  error
	)
	dstr = r.FormValue("dsw" + name + "Reg")
	if dstr == "" {
		*dsw = nil
		return ""
	}
	if date, err = time.ParseInLocation("2006-01-02", dstr, time.Local); err != nil {
		return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", dstr)
	}
	if *dsw == nil {
		*dsw = &person.DSWRegistration{Registered: date}
	} else {
		**dsw = person.DSWRegistration{Registered: date}
	}
	dstr = r.FormValue("dsw" + name + "Exp")
	if dstr == "" {
		return ""
	}
	if date, err = time.ParseInLocation("2006-01-02", dstr, time.Local); err != nil {
		return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", dstr)
	}
	if date.Before((*dsw).Registered) {
		return "Expiration date must not be before registration date."
	}
	(*dsw).Expiration = date
	return ""
}

func emitDSWCERT(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitDSW(form, up.DSWRegistrations.CERT, "CERT", focus, err)
}

func emitDSWComm(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitDSW(form, up.DSWRegistrations.Communications, "Communications (SARES)", focus, err)
}

func emitDSW(form *htmlb.Element, dsw *person.DSWRegistration, name string, focus bool, err string) {
	var date string
	form.E("div class='formRow-3col personeditStatusHeading'>DSW for %s", name)
	fs := form.E("div class=personeditDSW%s", name[:4])
	row := fs.E("div class=formRow")
	row.E("label for=personeditDSW%sReg>Registered", name[:4])
	if dsw != nil && !dsw.Registered.IsZero() {
		date = dsw.Registered.Format("2006-01-02")
	}
	row.E("input type=date id=personeditDSW%sReg name=dsw%sReg s-validate=.personeditDSW%s value=%s", name[:4], name[:4], name[:4], date,
		focus, "autofocus")
	row = fs.E("div class=formRow")
	row.E("label for=personeditDSW%sExp>Expiration", name[:4])
	ea := row.E("div class='formInput personeditDSWExpireAdd'")
	if dsw != nil && !dsw.Expiration.IsZero() {
		date = dsw.Expiration.Format("2006-01-02")
	} else {
		date = ""
	}
	ea.E("input type=date id=personeditDSW%sExp name=dsw%sExp class=formInput s-validate=.personeditDSW%s value=%s", name[:4], name[:4], name[:4], date)
	ea.E("button type=button class='s-btn s-btn-small s-btn-secondary personeditDSWExpireAdd'>+")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readBGDOJ(r *request.Request, up *person.Updater) string {
	return readBG(r, &up.BGChecks.DOJ, "DOJ", "NLI")
}

func readBGFBI(r *request.Request, up *person.Updater) string {
	return readBG(r, &up.BGChecks.FBI, "FBI", "NLI")
}

func readBGPHS(r *request.Request, up *person.Updater) string {
	return readBG(r, &up.BGChecks.PHS, "PHS", "Rescinded")
}

func readBG(r *request.Request, bg **person.BGCheck, tag, nliLabel string) string {
	cleared := r.FormValue("bg" + tag + "Cleared")
	nli := r.FormValue("bg" + tag + "NLI")
	assumed := r.FormValue("bg"+tag+"Assumed") != ""
	if cleared == "" && !assumed {
		*bg = nil
		return ""
	}
	if *bg == nil {
		*bg = &person.BGCheck{Assumed: assumed}
	} else {
		**bg = person.BGCheck{Assumed: assumed}
	}
	if cleared != "" {
		if date, err := time.ParseInLocation("2006-01-02", cleared, time.Local); err == nil {
			(*bg).Cleared = date
		} else {
			return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", cleared)
		}
	}
	if nli != "" {
		if date, err := time.ParseInLocation("2006-01-02", nli, time.Local); err == nil {
			(*bg).NLI = date
		} else {
			return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", cleared)
		}
	}
	if !(*bg).NLI.IsZero() && !(*bg).Cleared.Before((*bg).NLI) {
		return nliLabel + " date must not be before clearance date."
	}
	return ""
}

func emitBGDOJ(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitBG(form, up.BGChecks.DOJ, "LiveScan with CA DOJ", "DOJ", "NLI", focus, err)
}

func emitBGFBI(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitBG(form, up.BGChecks.FBI, "LiveScan with FBI", "FBI", "NLI", focus, err)
}

func emitBGPHS(form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitBG(form, up.BGChecks.PHS, "Full Background Check with PHS", "PHS", "Rescinded", focus, err)
}

func emitBG(form *htmlb.Element, bg *person.BGCheck, name, tag, nliLabel string, focus bool, err string) {
	var date string
	form.E("div class='formRow-3col personeditStatusHeading'>%s", name)
	fs := form.E("div class=personeditBG%s", tag)
	row := fs.E("div class=formRow")
	row.E("label for=personeditBG%sCleared>Cleared on", tag)
	ca := row.E("div class=formInput")
	if bg != nil && !bg.Cleared.IsZero() {
		date = bg.Cleared.Format("2006-01-02")
	}
	ca.E("input type=date id=personeditBG%sCleared name=bg%sCleared class=formInput s-validate=.personeditBG%s value=%s", tag, tag, tag, date,
		focus, "autofocus")
	ca.E("div class=personeditBGAssumed").E("input type=checkbox class=s-check name=bg%sAssumed label='Assumed cleared but no paper trail'", tag,
		bg != nil && bg.Assumed, "checked")
	row = fs.E("div class=formRow")
	row.E("label for=personeditBG%sNLI>%s on", tag, nliLabel)
	if bg != nil && !bg.NLI.IsZero() {
		date = bg.NLI.Format("2006-01-02")
	} else {
		date = ""
	}
	row.E("input type=date id=personeditBG%sNLI name=bg%sNLI class=formInput s-validate=.personeditBG%s value=%s", tag, tag, tag, date)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}
