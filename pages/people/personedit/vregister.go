package personedit

import (
	"bytes"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/sendmail"
	"sunnyvaleserv.org/portal/util/volgistics"
)

// HandleVRegister handles requests for /people/$id/vregister.
func HandleVRegister(r *request.Request, idstr string) {
	const personFields = volgistics.NewVolunteerFields | person.FVolgisticsID | person.FNotes
	var (
		user              *person.Person
		p                 *person.Person
		up                *person.Updater
		informalNameError string
		formalNameError   string
		sortNameError     string
		callSignError     string
		emailError        string
		email2Error       string
		cellPhoneError    string
		homePhoneError    string
		workPhoneError    string
		homeAddressError  string
		workAddressError  string
		mailAddressError  string
		interests         []string
		agreementError    string
		haveErrors        bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), personview.PersonFields|personFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if (user.ID() != p.ID() && !user.IsAdminLeader()) || p.VolgisticsID() != 0 {
		errpage.Forbidden(r, user)
		return
	}
	if r.FormValue("confirm") != "" {
		personview.Render(r, user, p, person.ViewFull, "")
		return
	}
	up = p.Updater()
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		informalNameError = readInformalName(r, up)
		formalNameError = readFormalName(r, up)
		sortNameError = readSortName(r, up)
		callSignError = readCallSign(r, up)
		emailError = readEmail(r, up)
		email2Error = readEmail2(r, up)
		cellPhoneError = readCellPhone(r, up)
		homePhoneError = readVRegHomePhone(r, up)
		workPhoneError = readWorkPhone(r, up)
		homeAddressError = readVRegHomeAddress(r, up)
		workAddressError = readWorkAddress(r, up)
		mailAddressError = readMailAddress(r, up)
		interests = readInterests(r)
		agreementError = readAgreement(r)
		haveErrors = informalNameError != "" || formalNameError != "" || sortNameError != "" || callSignError != "" || emailError != "" || email2Error != "" || cellPhoneError != "" || homePhoneError != "" || workPhoneError != "" || homeAddressError != "" || workAddressError != "" || mailAddressError != "" || agreementError != ""
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !haveErrors {
			if vid, err := sendVolunteerRegistration(p); err != nil {
				r.LogEntry.Problems.Add("volgistics registration failed: " + err.Error())
			} else {
				up.VolgisticsID = uint(vid)
				if len(interests) != 0 {
					up.Notes = append(up.Notes, &person.Note{
						Note:       "Volunteer registration interests: " + strings.Join(interests, ", "),
						Date:       time.Now(),
						Visibility: person.NoteVisibleToAdmins,
					})
				}
			}
			r.Transaction(func() {
				p.Update(r, up, personFields)
			})
			confirmRegistration(r, p)
			return
		}
	}
	r.HTMLNoCache()
	if haveErrors {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class=form method=POST up-main up-layer=parent up-target=.personview")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Register as a City Volunteer"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	emitVRegIntro(r, form)
	if len(validate) == 0 || slices.Contains(validate, "informalName") {
		emitInformalName(r, form, up, informalNameError != "" || (formalNameError == "" && sortNameError == "" && callSignError == ""), informalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "formalName") {
		emitFormalName(r, form, up, formalNameError != "", formalNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "sortName") {
		emitSortName(r, form, up, sortNameError != "", sortNameError)
	}
	if len(validate) == 0 || slices.Contains(validate, "callSign") {
		emitCallSign(r, form, up, callSignError != "", callSignError)
	}
	if len(validate) == 0 || slices.Contains(validate, "email") {
		emitEmail(r, form, up, emailError != "" || !haveErrors, emailError)
	}
	if len(validate) == 0 || slices.Contains(validate, "email2") {
		emitEmail2(r, form, up, email2Error != "", email2Error)
	}
	if len(validate) == 0 || slices.Contains(validate, "cellPhone") {
		emitCellPhone(r, form, up, cellPhoneError != "", cellPhoneError)
	}
	if len(validate) == 0 || slices.Contains(validate, "homePhone") {
		emitHomePhone(r, form, up, homePhoneError != "", homePhoneError)
	}
	if len(validate) == 0 || slices.Contains(validate, "workPhone") {
		emitWorkPhone(r, form, up, workPhoneError != "", workPhoneError)
	}
	if len(validate) == 0 {
		emitHomeAddress(r, form, up, homeAddressError != "", homeAddressError)
		emitWorkAddress(r, form, up, workAddressError != "", workAddressError)
		emitMailAddress(r, form, up, mailAddressError != "", mailAddressError)
		emitInterests(r, form, interests)
		emitAgreement(r, form, agreementError)
		emitVRegisterButtons(r, form)
	}
}

func emitVRegIntro(r *request.Request, form *htmlb.Element) {
	form.E("div class='formRow-3col vregisterIntro'").R(r.Loc("Thank you for your interest in volunteering with the City of Sunnyvale, Office of Emergency Services.  Please complete this form to register as a City of Sunnyvale Volunteer.  (Please note: registering as a city volunteer is not required for taking one of our classes.  It is only required when joining one of our volunteer groups.)"))
}

func readVRegHomePhone(r *request.Request, up *person.Updater) string {
	if s := readHomePhone(r, up); s != "" {
		return s
	}
	if up.CellPhone == "" && up.HomePhone == "" {
		return r.Loc("A cell or home phone number is required.")
	}
	return ""
}

func readVRegHomeAddress(r *request.Request, up *person.Updater) string {
	if s := readHomeAddress(r, up); s != "" {
		return s
	}
	if up.Addresses.Home == nil || up.Addresses.Home.Address == "" {
		return r.Loc("Your home address is required.")
	}
	return ""
}

func readInterests(r *request.Request) (interests []string) {
	return r.Form["interests"]
}

func emitInterests(r *request.Request, form *htmlb.Element, interests []string) {
	row := form.E("div class=formRow")
	row.E("label for=vregisterCERTD").R(r.Loc("Interests"))
	box := row.E("div class=formInput")
	box.E("input type=checkbox id=vregisterCERTD name=interests value=CERT-D class=s-check label=%s",
		r.Loc("CERT Deployment Team"),
		slices.Contains(interests, "CERT-D"), "checked")
	box.E("input type=checkbox name=interests value=Outreach class=s-check label=%s",
		r.Loc("Community Outreach"),
		slices.Contains(interests, "Outreach"), "checked")
	box.E("input type=checkbox name=interests value=SARES class=s-check label=%s",
		r.Loc("Amateur Radio (SARES)"),
		slices.Contains(interests, "SARES"), "checked")
	box.E("input type=checkbox name=interests value=SNAP class=s-check label=%s",
		r.Loc("Neighborhood Preparedness Facilitator"),
		slices.Contains(interests, "SNAP"), "checked")
	box.E("input type=checkbox name=interests value=Listos class=s-check label=%s",
		r.Loc("Preparedness Class Instructor"),
		slices.Contains(interests, "Listos"), "checked")
	box.E("input type=checkbox name=interests value=CERT-T class=s-check label=%s",
		r.Loc("CERT Basic Training Instructor"),
		slices.Contains(interests, "CERT-T"), "checked")
}

func readAgreement(r *request.Request) (err string) {
	if r.FormValue("agreement") == "" {
		return r.Loc("Please check that you agree with the above statement in order to register.")
	}
	return ""
}

func emitAgreement(r *request.Request, form *htmlb.Element, err string) {
	form.E("div class='formRow-3col vregisterAgreement'").R(r.Loc("By submitting this application, I certify that all statements I have made on this application are true and correct and I hereby authorize the City of Sunnyvale to investigate the accuracy of this information.  I am aware that fingerprinting and a criminal records search is required for volunteers 18 years of age or older.  I understand that I am working at all times on a voluntary basis, without monetary compensation or benefits, and not as a paid employee.  I give the City of Sunnyvale permission to use any photographs or videos taken of me during my service without obligation or compensation to me.  I understand that the City of Sunnyvale reserves the right to terminate a volunteer's service at any time.  I understand that volunteers are covered under the City of Sunnyvale's Worker's Compensation Program for an injury or accident occurring while on duty."))
	row := form.E("div class=formRow")
	row.E("div class=formInput").E("input type=checkbox name=agreement class=s-check label=%s", r.Loc("I agree"))
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func emitVRegisterButtons(r *request.Request, form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss").R(r.Loc("Cancel"))
	buttons.E("input type=submit class='sbtn sbtn-primary' value=%s", r.Loc("Register"))
}

var addressSplitRE = regexp.MustCompile(`^(.*),\s*([^,]+),?\s+CA\s+(\d{5})$`)

func sendVolunteerRegistration(p *person.Person) (vid int, err error) {
	var c *volgistics.Client

	if c, err = volgistics.Login(); err != nil {
		return 0, err
	}
	if vid, err = c.NewVolunteer(p); err != nil {
		return 0, err
	}
	return vid, nil
}

func confirmRegistration(r *request.Request, p *person.Person) {
	r.HTMLNoCache()
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class=form method=POST up-main up-layer=parent up-target=.personview")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Register as a City Volunteer"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("div class='formRow-3col vregisterIntro'").R(r.Loc("Thank you for volunteering with the City of Sunnyvale, Office of Emergency Services.  One of our staff will contact you to schedule a fingerprinting appointment.  (Criminal history checks are required by city policy for all public-facing volunteers.)  If you have not heard from us within a few days, please email us at oes@sunnyvale.ca.gov to follow up.  We look forward to working with you!"))
	buttons := form.E("div class=formButtons")
	buttons.E("button type=submit name=confirm class='sbtn sbtn-primary' value=%s", r.Loc("OK"))
	// Also want to send an email to the admin.
	var body bytes.Buffer
	fmt.Fprintf(&body, "From: SunnyvaleSERV.org <admin@sunnyvaleserv.org>\r\nTo: admin@sunnyvaleserv.org\r\nSubject: New Volunteer Registration\r\n\r\n%s has submitted a volunteer registration.\r\n", p.InformalName())
	sendmail.SendMessage(r.Context(), config.Get("fromaddr"), []string{config.Get("adminEmail")}, body.Bytes())
}
