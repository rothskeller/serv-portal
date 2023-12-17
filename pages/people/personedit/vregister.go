package personedit

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// HandleVRegister handles requests for /people/$id/vregister.
func HandleVRegister(r *request.Request, idstr string) {
	const personFields = person.FID | person.FInformalName | person.FVolgisticsID | person.FFlags | person.FFormalName | person.FSortName | person.FCallSign | person.FBirthdate | person.FEmail | person.FEmail2 | person.FCellPhone | person.FHomePhone | person.FWorkPhone | person.FAddresses | person.FEmContacts
	var (
		user                 *person.Person
		p                    *person.Person
		up                   *person.Updater
		informalNameError    string
		formalNameError      string
		sortNameError        string
		callSignError        string
		birthdateError       string
		emailError           string
		email2Error          string
		cellPhoneError       string
		homePhoneError       string
		workPhoneError       string
		homeAddressError     string
		workAddressError     string
		mailAddressError     string
		interests            []string
		ec1NameError         string
		ec1HomePhoneError    string
		ec1CellPhoneError    string
		ec1RelationshipError string
		ec2NameError         string
		ec2HomePhoneError    string
		ec2CellPhoneError    string
		ec2RelationshipError string
		agreementError       string
		haveErrors           bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil || !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), personview.PersonFields|personFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if (user.ID() != p.ID() && !user.IsAdminLeader()) || p.VolgisticsID() != 0 || p.Flags()&person.VolgisticsPending != 0 {
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
		birthdateError = readVRegBirthdate(r, up)
		emailError = readEmail(r, up)
		email2Error = readEmail2(r, up)
		cellPhoneError = readCellPhone(r, up)
		homePhoneError = readVRegHomePhone(r, up)
		workPhoneError = readWorkPhone(r, up)
		homeAddressError = readVRegHomeAddress(r, up)
		workAddressError = readWorkAddress(r, up)
		mailAddressError = readMailAddress(r, up)
		interests = readInterests(r)
		ec1NameError = readVRegECName(r, up, 0)
		ec1HomePhoneError = readECHomePhone(r, up, 0)
		ec1CellPhoneError = readECCellPhone(r, up, 0)
		ec1RelationshipError = readECRelationship(r, up, 0)
		ec2NameError = readVRegECName(r, up, 1)
		ec2HomePhoneError = readECHomePhone(r, up, 1)
		ec2CellPhoneError = readECCellPhone(r, up, 1)
		ec2RelationshipError = readECRelationship(r, up, 1)
		agreementError = readAgreement(r)
		haveErrors = informalNameError != "" || formalNameError != "" || sortNameError != "" || callSignError != "" || birthdateError != "" || emailError != "" || email2Error != "" || cellPhoneError != "" || homePhoneError != "" || workPhoneError != "" || homeAddressError != "" || workAddressError != "" || mailAddressError != "" || ec1NameError != "" || ec1HomePhoneError != "" || ec1CellPhoneError != "" || ec1RelationshipError != "" || ec2NameError != "" || ec2HomePhoneError != "" || ec2CellPhoneError != "" || ec2RelationshipError != "" || agreementError != ""
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !haveErrors {
			r.Transaction(func() {
				p.Update(r, up, personFields)
			})
			if err := SendVolunteerRegistration(p, interests); err != nil {
				r.LogEntry.Problems.Add("volgistics registration failed: " + err.Error())
			}
			personview.Render(r, user, p, person.ViewFull, "")
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
	if len(validate) == 0 || slices.Contains(validate, "birthdate") {
		emitBirthdate(r, form, up, birthdateError != "", birthdateError)
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
		emitEmergencyContact(r, form, up, 0, ec1NameError, ec1HomePhoneError, ec1CellPhoneError, ec1RelationshipError)
		emitEmergencyContact(r, form, up, 1, ec2NameError, ec2HomePhoneError, ec2CellPhoneError, ec2RelationshipError)
		emitAgreement(r, form, agreementError)
		emitVRegisterButtons(r, form)
	}
}

func emitVRegIntro(r *request.Request, form *htmlb.Element) {
	form.E("div class='formRow-3col vregisterIntro'").R(r.Loc("Thank you for your interest in volunteering with the City of Sunnyvale, Office of Emergency Services.  Please complete this form to register as a City of Sunnyvale Volunteer.  Once we receive your registration (which usually takes a few days) we will contact you to schedule an appointment for your fingerprinting.  (Please note: registering as a city volunteer is not required for taking one of our classes.  It is only required when joining one of our volunteer groups.)"))
}

func readVRegBirthdate(r *request.Request, up *person.Updater) string {
	if s := readBirthdate(r, up); s != "" {
		return s
	}
	if up.Birthdate == "" {
		return r.Loc("Your birthdate is required.")
	}
	return ""
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
	if up.Addresses.Home.Address == "" {
		return r.Loc("Your home address is required.")
	}
	return ""
}

func readVRegECName(r *request.Request, up *person.Updater, num int) string {
	readECName(r, up, num)
	if up.EmContacts[num].Name == "" {
		return r.Loc("The emergency contacts are required.")
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
	row.E("div class=formInput").E("input type=checkbox name=agreement class=s-check label=%s").R(r.Loc("I agree"))
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

func SendVolunteerRegistration(p *person.Person, interests []string) (err error) {
	var out jwriter.Writer

	parts := strings.SplitN(p.SortName(), ", ", 2)
	out.RawString(`{"g1f2":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g1f1":`)
	out.String(parts[0])
	out.RawString(`,"g1f5":0,"g1f6":`)
	match := addressSplitRE.FindStringSubmatch(p.Addresses().Home.Address)
	if match == nil {
		return errors.New("can't parse home address")
	}
	out.String(match[1])
	out.RawString(`,"g1f8":`)
	out.String(match[2])
	out.RawString(`,"g1f9":104,"g1f10":`)
	out.String(match[3])
	if p.HomePhone() != "" {
		out.RawString(`,"g1f11":`)
		out.String(volgisticsPhoneFormat(p.HomePhone()))
		if p.CellPhone() != "" {
			out.RawString(`,"g1f13":`)
			out.String(volgisticsPhoneFormat(p.CellPhone()))
		}
	} else {
		out.RawString(`,"g1f11":`)
		out.String(volgisticsPhoneFormat(p.CellPhone()))
	}
	out.RawString(`,"g1f22":false,"g1f23":false,"g1f17":`)
	out.String(p.Email())
	out.RawString(`,"g2f99":`)
	if len(interests) != 0 {
		out.String("ATTN: Steve Roth\nSunnyvale DPS/OES\nInterested in " + strings.Join(interests, ", "))
	} else {
		out.String("ATTN: Steve Roth\nSunnyvale DPS/OES")
	}
	out.RawString(`,"g8f1":`)
	out.String(p.Birthdate())
	out.RawString(`,"g8f4":0,"g8f5":0,"g8f6":0,"g8f7":0,"g8f21":"None","g9f3c251":false,"g9f3c252":false,"g9f3c253":false,"g9f3c254":false,"g9f3c255":false,"g9f3c256":false,"g9f3c257":false,"g9f3c258":false,"g9f1c213":false,"g9f1c214":false,"g9f1c215":false,"g9f1c233":false,"g9f1c223":false,"g9f1c216":false,"g9f1c339":false,"g9f1c217":false,"g9f1c218":false,"g9f1c219":false,"g9f1c220":false,"g9f1c221":false,"g9f1c222":false,"g9f1c302":false,"g9f1c224":false,"g9f1c225":false,"g9f1c226":false,"g9f1c227":false,"g9f1c232":false,"g9f1c228":false,"g9f1c229":false,"g9f1c230":false,"g9f1c231":false,"g9f1c334":false,"g9f2c234":false,"g9f2c235":false,"g9f2c236":false,"g9f2c237":false,"g9f2c238":false,"g9f2c239":false,"g9f2c240":false,"g9f2c241":false,"g9f2c242":false,"g9f2c243":false,"g9f2c244":false,"g9f2c245":false,"g9f2c246":false,"g9f2c247":false,"g9f2c248":false,"g9f2c249":false,"g9f2c250":false,"g9f22c93":false,"g9f22c85":false,"g9f22c143":false,"g9f22c129":false,"g9f22c75":false,"g9f22c18":false,"g9f22c73":false,"g9f22c16":false,"g9f22c77":false,"g9f22c79":false,"g9f22c81":false,"g9f22c41":false,"g9f22c83":false,"g9f22c91":false,"g9f22c1":false,"g2f2c1l1":`)
	parts = strings.Fields(p.EmContacts()[0].Name)
	out.String(strings.Join(parts[:len(parts)-1], " "))
	out.RawString(`,"g2f1c1l1":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g2f0c1l1":0`)
	if p.EmContacts()[0].HomePhone != "" {
		out.RawString(`,"g2f11c1l1":`)
		out.String(volgisticsPhoneFormat(p.EmContacts()[0].HomePhone))
		if p.EmContacts()[0].CellPhone != "" {
			out.RawString(`,"g2f13c1l1":`)
			out.String(volgisticsPhoneFormat(p.EmContacts()[0].CellPhone))
		}
	} else {
		out.RawString(`,"g2f11c1l1":`)
		out.String(volgisticsPhoneFormat(p.EmContacts()[0].CellPhone))
	}
	out.RawString(`,"g2f25c1l1":`)
	out.Int(relationshipCode[p.EmContacts()[0].Relationship])
	out.RawString(`,"g2f2c1l2":`)
	parts = strings.Fields(p.EmContacts()[1].Name)
	out.String(strings.Join(parts[:len(parts)-1], " "))
	out.RawString(`,"g2f1c1l2":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g2f0c1l2":0`)
	if p.EmContacts()[1].HomePhone != "" {
		out.RawString(`,"g2f11c1l2":`)
		out.String(volgisticsPhoneFormat(p.EmContacts()[1].HomePhone))
		if p.EmContacts()[1].CellPhone != "" {
			out.RawString(`,"g2f13c1l2":`)
			out.String(volgisticsPhoneFormat(p.EmContacts()[1].CellPhone))
		}
	} else {
		out.RawString(`,"g2f11c1l2":`)
		out.String(volgisticsPhoneFormat(p.EmContacts()[1].CellPhone))
	}
	out.RawString(`,"g2f25c1l2":`)
	out.Int(relationshipCode[p.EmContacts()[1].Relationship])
	out.RawString(`,"g13f1":true,"g16f3c0l1":false,"g16f3c0l2":false,"g16f3c0l3":false,"g16f3c1l1":false,"g16f3c1l2":false,"g16f3c1l3":false,"g16f3c2l1":false,"g16f3c2l2":false,"g16f3c2l3":false,"g16f3c3l1":false,"g16f3c3l2":false,"g16f3c3l3":false,"g16f3c4l1":false,"g16f3c4l2":false,"g16f3c4l3":false,"g16f3c5l1":false,"g16f3c5l2":false,"g16f3c5l3":false,"g16f3c6l1":false,"g16f3c6l2":false,"g16f3c6l3":false,"appSignature":[],"id":"777028848"}`)
	var body, _ = out.BuildBytes()

	var req, _ = http.NewRequest(http.MethodPost, "https://www.volgistics.com/ex/apForm.dll/AFR?Action=saveApplication", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("volgistics registration failed with status %d", resp.StatusCode)
	}
	return nil
}

func volgisticsPhoneFormat(s string) string {
	return fmt.Sprintf("(%s) %s", s[0:3], s[4:])
}

var relationshipCode = map[string]int{
	"Co-worker":  125,
	"Daughter":   313,
	"Father":     120,
	"Friend":     123,
	"Mother":     119,
	"Neighbor":   126,
	"Other":      314,
	"Relative":   275,
	"Son":        121,
	"Spouse":     331,
	"Supervisor": 124,
}
