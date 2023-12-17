package personedit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const contactPersonFields = person.FInformalName | person.FCallSign | person.FEmail | person.FEmail2 | person.FCellPhone | person.FHomePhone | person.FWorkPhone | person.FAddresses | person.FEmContacts
const addressVerificationAPI = "https://addressvalidation.googleapis.com/v1:validateAddress?key="

// HandleContact handles requests for /people/$id/edcontact.
func HandleContact(r *request.Request, idstr string) {
	var (
		user                 *person.Person
		p                    *person.Person
		up                   *person.Updater
		canEditEmContacts    bool
		emailError           string
		email2Error          string
		cellPhoneError       string
		homePhoneError       string
		workPhoneError       string
		homeAddressError     string
		workAddressError     string
		mailAddressError     string
		ec1HomePhoneError    string
		ec1CellPhoneError    string
		ec1RelationshipError string
		ec2HomePhoneError    string
		ec2CellPhoneError    string
		ec2RelationshipError string
		haveErrors           bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), contactPersonFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if user.ID() != p.ID() && !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	canEditEmContacts = user.ID() == p.ID() || user.IsAdminLeader()
	up = p.Updater()
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		emailError = readEmail(r, up)
		email2Error = readEmail2(r, up)
		cellPhoneError = readCellPhone(r, up)
		homePhoneError = readHomePhone(r, up)
		workPhoneError = readWorkPhone(r, up)
		homeAddressError = readHomeAddress(r, up)
		workAddressError = readWorkAddress(r, up)
		mailAddressError = readMailAddress(r, up)
		if canEditEmContacts {
			readECName(r, up, 0)
			ec1HomePhoneError = readECHomePhone(r, up, 0)
			ec1CellPhoneError = readECCellPhone(r, up, 0)
			ec1RelationshipError = readECRelationship(r, up, 0)
			readECName(r, up, 1)
			ec2HomePhoneError = readECHomePhone(r, up, 1)
			ec2CellPhoneError = readECCellPhone(r, up, 1)
			ec2RelationshipError = readECRelationship(r, up, 1)
		}
		haveErrors = emailError != "" || email2Error != "" || cellPhoneError != "" || homePhoneError != "" || workPhoneError != "" || homeAddressError != "" || workAddressError != "" || mailAddressError != "" || ec1HomePhoneError != "" || ec1CellPhoneError != "" || ec1RelationshipError != "" || ec2HomePhoneError != "" || ec2CellPhoneError != "" || ec2RelationshipError != ""
		// If there were no errors *and* we're not validating, save the
		// data and return to the view page.
		if len(validate) == 0 && !haveErrors {
			up.EmContacts = slices.DeleteFunc(up.EmContacts, func(ec *person.EmContact) bool { return ec.Name == "" })
			r.Transaction(func() {
				p.Update(r, up, contactPersonFields)
			})
			personview.Render(r, user, p, person.ViewFull, "contact")
			return
		}
	}
	r.HTMLNoCache()
	if haveErrors {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewContact")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Edit Contact Information"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
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
		if canEditEmContacts {
			emitEmergencyContact(r, form, up, 0, "", ec1HomePhoneError, ec1CellPhoneError, ec1RelationshipError)
			emitEmergencyContact(r, form, up, 1, "", ec2HomePhoneError, ec2CellPhoneError, ec2RelationshipError)
		}
		emitButtons(r, form)
	}
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func readEmail(r *request.Request, up *person.Updater) string {
	up.Email = strings.ToLower(strings.TrimSpace(r.FormValue("email")))
	if up.Email == "" {
		up.Email = strings.ToLower(strings.TrimSpace(r.FormValue("email2")))
	}
	if up.Email != "" && !emailRE.MatchString(up.Email) {
		return fmt.Sprintf(r.Loc("%q is not a valid email address."), up.Email)
	}
	if up.DuplicateEmail(r) {
		return fmt.Sprintf(r.Loc("The email address %q is in use by another person."), up.Email)
	}
	return ""
}

func emitEmail(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditEmail").R(r.Loc("Email"))
	row.E("input id=personeditEmail name=email s-validate value=%s", up.Email, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
	row.E("div class=formHelp").R(r.Loc("This is the email address you log in with."))
}

func readEmail2(r *request.Request, up *person.Updater) string {
	up.Email2 = strings.ToLower(strings.TrimSpace(r.FormValue("email2")))
	if up.Email2 != "" && !emailRE.MatchString(up.Email2) {
		return fmt.Sprintf(r.Loc("%q is not a valid email address."), up.Email2)
	}
	if up.Email == up.Email2 {
		up.Email2 = ""
	}
	return ""
}

func emitEmail2(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditEmail2").R(r.Loc("Alt. Email"))
	row.E("input id=personeditEmail2 name=email2 s-validate value=%s", up.Email2, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readCellPhone(r *request.Request, up *person.Updater) string {
	up.CellPhone = strings.TrimSpace(r.FormValue("cellPhone"))
	if !fmtPhone(&up.CellPhone, false) {
		return fmt.Sprintf(r.Loc("%q is not a valid 10-digit phone number."), up.CellPhone)
	}
	// In theory, we could use a Twilio API to verify that it's really a
	// mobile phone number.  Probably not worth bothering.
	return ""
}

func emitCellPhone(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditCellPhone").R(r.Loc("Cell Phone"))
	row.E("input id=personeditCellPhone name=cellPhone s-validate value=%s", up.CellPhone, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readHomePhone(r *request.Request, up *person.Updater) string {
	up.HomePhone = strings.TrimSpace(r.FormValue("homePhone"))
	if !fmtPhone(&up.HomePhone, false) {
		return fmt.Sprintf(r.Loc("%q is not a valid 10-digit phone number."), up.HomePhone)
	}
	if up.HomePhone == up.CellPhone {
		up.HomePhone = ""
	}
	return ""
}

func emitHomePhone(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditHomePhone").R(r.Loc("Home Phone"))
	row.E("input id=personeditHomePhone name=homePhone s-validate value=%s", up.HomePhone, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readWorkPhone(r *request.Request, up *person.Updater) string {
	up.WorkPhone = strings.TrimSpace(r.FormValue("workPhone"))
	if !fmtPhone(&up.WorkPhone, true) {
		return fmt.Sprintf(r.Loc("%q is not a valid phone number."), up.WorkPhone)
	}
	if up.WorkPhone == up.CellPhone || up.WorkPhone == up.HomePhone {
		up.WorkPhone = ""
	}
	return ""
}

func emitWorkPhone(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	row := form.E("div class='formRow'")
	row.E("label for=personeditWorkPhone").R(r.Loc("Work Phone"))
	row.E("input id=personeditWorkPhone name=workPhone s-validate value=%s", up.WorkPhone, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

type (
	addressVerifyRequest struct {
		Address addressVerifyRequestAddress `json:"address"`
	}
	addressVerifyRequestAddress struct {
		AddressLines []string `json:"addressLines"`
	}
	addressVerifyResponse struct {
		Result addressVerifyResponseResult `json:"result"`
	}
	addressVerifyResponseResult struct {
		Verdict addressVerifyResponseVerdict `json:"verdict"`
		Address addressVerifyResponseAddress `json:"address"`
		Geocode addressVerifyResponseGeocode `json:"geocode"`
	}
	addressVerifyResponseVerdict struct {
		AddressComplete       bool   `json:"addressComplete"`
		ValidationGranularity string `json:"validationGranularity"`
		GeocodeGranularity    string `json:"geocodeGranularity"`
	}
	addressVerifyResponseAddress struct {
		FormattedAddress string `json:"formattedAddress"`
	}
	addressVerifyResponseGeocode struct {
		Location addressVerifyLatLng `json:"location"`
	}
	addressVerifyLatLng struct {
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	}
)

var zip4RE = regexp.MustCompile(`-\d\d\d\d, USA`)

func readHomeAddress(r *request.Request, up *person.Updater) string {
	return readAddress(r, &up.Addresses.Home, "home", false, nil, false)
}
func readWorkAddress(r *request.Request, up *person.Updater) string {
	return readAddress(r, &up.Addresses.Work, "work", true, up.Addresses.Home, true)
}
func readMailAddress(r *request.Request, up *person.Updater) string {
	return readAddress(r, &up.Addresses.Mail, "mail", true, up.Addresses.Home, false)
}
func readAddress(
	r *request.Request, addr **person.Address, name string, canBeSameAsHome bool, home *person.Address, canGeocode bool,
) string {
	// If the same-as-home flag is checked, that's all we need to save.
	if canBeSameAsHome && r.FormValue(name+"SameAsHome") != "" {
		*addr = &person.Address{SameAsHome: true}
		if home == nil {
			return r.Loc("This address cannot be marked “same as home” when there is no home address.")
		}
		return ""
	}
	// If the input fields are empty, clear the address.
	line1 := strings.TrimSpace(r.FormValue(name + "Address"))
	line2 := strings.TrimSpace(r.FormValue(name + "CSZ"))
	if line1 == "" && line2 == "" {
		*addr = nil
		return ""
	}
	// We have an address.  If that address matches what we already have
	// saved for this person, we're done.
	lines := line1 + ", " + line2
	if *addr != nil && lines == (*addr).Address {
		return ""
	}
	// Set up to store the new address.
	if *addr == nil {
		*addr = &person.Address{Address: lines}
	} else {
		**addr = person.Address{Address: lines}
	}
	// Send the address to the address verification service.
	var answer addressVerifyResponse
	body, _ := json.Marshal(addressVerifyRequest{addressVerifyRequestAddress{[]string{line1, line2}}})
	resp, err := http.Post(addressVerificationAPI+config.Get("addressVerificationKey"), "application/json", bytes.NewReader(body))
	if err == nil {
		defer resp.Body.Close()
	}
	if err == nil && resp.StatusCode == http.StatusOK {
		err = json.NewDecoder(resp.Body).Decode(&answer)
	}
	if err != nil || resp.StatusCode != http.StatusOK {
		if err != nil {
			r.LogEntry.Problems.AddF("address verification failure %s", err)
		} else {
			r.LogEntry.Problems.AddF("address verification failure %s", resp.Status)
		}
		return r.Loc("Address changes cannot be accepted right now because the address verification service is offline.")
	}
	// If we got back a match for the address, save the reformatted address.
	switch answer.Result.Verdict.ValidationGranularity {
	case "SUB_PREMISE", "PREMISE":
		(*addr).Address = zip4RE.ReplaceAllLiteralString(answer.Result.Address.FormattedAddress, "")
	default:
		return r.Loc("This is not a valid address.")
	}
	// If geocoding is needed, save the coordinates.
	if canGeocode {
		switch answer.Result.Verdict.GeocodeGranularity {
		case "SUB_PREMISE", "PREMISE", "PREMISE_PROXIMITY":
			(*addr).Latitude = answer.Result.Geocode.Location.Latitude
			(*addr).Longitude = answer.Result.Geocode.Location.Longitude
			(*addr).FireDistrict = person.FireDistrict(*addr)
		}
	}
	return ""
}

func emitHomeAddress(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitAddress(r, form, up.Addresses.Home, "Home", r.Loc("Home Address"), false, focus, err)
}
func emitWorkAddress(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitAddress(r, form, up.Addresses.Work, "Work", r.Loc("Work Address"), true, focus, err)
}
func emitMailAddress(r *request.Request, form *htmlb.Element, up *person.Updater, focus bool, err string) {
	emitAddress(r, form, up.Addresses.Mail, "Mail", r.Loc("Mailing Address"), true, focus, err)
}
func emitAddress(r *request.Request, form *htmlb.Element, addr *person.Address, name4, name string, canSameAsHome, focus bool, err string) {
	var line1, line2 string
	var lname = strings.ToLower(name4)

	if addr != nil {
		parts := strings.SplitN(addr.Address, ",", 2)
		line1 = strings.TrimSpace(parts[0])
		if len(parts) == 2 {
			line2 = strings.TrimSpace(parts[1])
		}
	}
	row := form.E("div class='formRow personeditAddress'")
	if canSameAsHome {
		row.E("label for=personedit%sSameAsHome>%s", name4, name)
	} else {
		row.E("label for=personedit%sAddress>%s", name4, name)
	}
	in := row.E("div class=formInput")
	if canSameAsHome {
		in.E("div").E("input type=checkbox class=s-check id=personedit%sSameAsHome name=%sSameAsHome", name4, lname,
			"label=%s up-switch=.personedit%sAddressInput", r.Loc("Same as home address"), name4,
			addr != nil && addr.SameAsHome, "checked")
	}
	in.E("div").E("input id=personedit%sAddress name=%sAddress", name4, lname,
		"class='formInput personedit%sAddressInput personeditAddressLine1' value=%s", name4, line1,
		canSameAsHome, "up-hide-for=:checked",
		focus, "autofocus")
	in.E("div").E("input name=%sCSZ class='formInput personedit%sAddressInput personeditAddressLine2' value=%s", lname, name4, line2,
		canSameAsHome, "up-hide-for=:checked s-validate")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

var emContactRelationships = []string{
	// These are defined in Volgistics and should not be changed unilaterally.
	"Co-worker", "Daughter", "Father", "Friend", "Mother", "Neighbor", "Other",
	"Relative", "Son", "Spouse", "Supervisor",
}

func readECName(r *request.Request, up *person.Updater, num int) {
	for len(up.EmContacts) <= num {
		up.EmContacts = append(up.EmContacts, new(person.EmContact))
	}
	var ec = up.EmContacts[num]
	ec.Name = strings.TrimSpace(r.FormValue(fmt.Sprintf("emcname%d", num)))
}
func readECHomePhone(r *request.Request, up *person.Updater, num int) string {
	var ec = up.EmContacts[num]
	ec.HomePhone = strings.TrimSpace(r.FormValue(fmt.Sprintf("emchome%d", num)))
	if ec.Name == "" {
		if ec.HomePhone != "" {
			return r.Loc("A phone number may not be specified without a name.")
		}
		return ""
	}
	if !fmtPhone(&ec.HomePhone, false) {
		return fmt.Sprintf(r.Loc("%q is not a valid 10-digit phone number."), ec.HomePhone)
	}
	return ""
}
func readECCellPhone(r *request.Request, up *person.Updater, num int) string {
	var ec = up.EmContacts[num]
	ec.CellPhone = strings.TrimSpace(r.FormValue(fmt.Sprintf("emccell%d", num)))
	if ec.Name == "" {
		if ec.CellPhone != "" {
			return r.Loc("A phone number may not be specified without a name.")
		}
		return ""
	}
	if !fmtPhone(&ec.CellPhone, false) {
		return fmt.Sprintf(r.Loc("%q is not a valid 10-digit phone number."), ec.CellPhone)
	}
	if ec.HomePhone == ec.CellPhone {
		ec.HomePhone = ""
	}
	if ec.HomePhone == "" && ec.CellPhone == "" {
		return r.Loc("At least one phone number is required.")
	}
	return ""
}
func readECRelationship(r *request.Request, up *person.Updater, num int) string {
	var ec = up.EmContacts[num]
	ec.Relationship = r.FormValue(fmt.Sprintf("emcrel%d", num))
	if ec.Name == "" {
		if ec.Relationship != "" {
			return r.Loc("A relationship may not be specified without a name.")
		}
		return ""
	}
	if ec.Relationship == "" {
		return r.Loc("The relationship is required.")
	}
	if idx := slices.IndexFunc(emContactRelationships, func(s string) bool {
		return r.Loc(s) == ec.Relationship
	}); idx >= 0 {
		ec.Relationship = emContactRelationships[idx]
	}
	if !slices.Contains(emContactRelationships, ec.Relationship) {
		return fmt.Sprintf(r.Loc("%q is not one of the relationship choices."), ec.Relationship)
	}
	return ""
}

func emitEmergencyContact(r *request.Request, form *htmlb.Element, up *person.Updater, num int, nameErr, homePhoneErr, cellPhoneErr, relationshipErr string) {
	var ec *person.EmContact
	if num >= len(up.EmContacts) {
		ec = new(person.EmContact)
	} else {
		ec = up.EmContacts[num]
	}
	form.E("div class='formRow-3col personeditEmContact'>%s %d", r.Loc("Emergency Contact"), num+1)
	row := form.E("div class=formRow")
	row.E("label for=personeditEmContactName%d", num).R(r.Loc("Name"))
	row.E("input id=personeditEmContactName%d name=emcname%d class=formInput value=%s", num, num, ec.Name)
	if nameErr != "" {
		row.E("div class=formError>%s", nameErr)
	}
	row = form.E("div class=formRow")
	row.E("label for=personeditEmContactHomePhone%d", num).R(r.Loc("Home Phone"))
	row.E("input id=personeditEmContactHomePhone%d name=emchome%d class=formInput value=%s", num, num, ec.HomePhone,
		homePhoneErr != "", "autofocus")
	if homePhoneErr != "" {
		row.E("div class=formError>%s", homePhoneErr)
	}
	row = form.E("div class=formRow")
	row.E("label for=personeditEmContactCellPhone%d", num).R(r.Loc("Cell Phone"))
	row.E("input id=personeditEmContactCellPhone%d name=emccell%d class=formInput value=%s", num, num, ec.CellPhone,
		cellPhoneErr != "", "autofocus")
	if cellPhoneErr != "" {
		row.E("div class=formError>%s", cellPhoneErr)
	}
	row = form.E("div class=formRow")
	row.E("label for=personeditEmContactRelationship%d", num).R(r.Loc("Relationship"))
	sel := row.E("select id=personeditEmContactRelationship%d name=emcrel%d class=formInput", num, num,
		relationshipErr != "", "autofocus")
	if ec.Relationship == "" {
		sel.E("option value='' selected").R(r.Loc("(select relationship)"))
	}
	for _, rel := range emContactRelationships {
		sel.E("option", rel == ec.Relationship, "selected").T(r.Loc(rel))
	}
	if relationshipErr != "" {
		row.E("div class=formError>%s", relationshipErr)
	}
}

func fmtPhone(p *string, extraOK bool) bool {
	digits := strings.Map(func(r rune) rune {
		if r < '0' || r > '9' {
			return -1
		}
		return r
	}, *p)
	if len(digits) == 11 && digits[0] == '1' {
		digits = digits[1:]
	}
	switch len(digits) {
	case 0:
		*p = ""
		return true
	case 10:
		*p = digits[0:3] + "-" + digits[3:6] + "-" + digits[6:10]
		return true
	}
	if len(digits) > 10 && extraOK {
		*p = digits[0:3] + "-" + digits[3:6] + "-" + digits[6:10] + "x" + digits[10:]
		return true
	}
	return false
}
