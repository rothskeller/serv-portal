package classes

import (
	"bytes"
	"fmt"
	"net/http"
	"net/mail"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/server/l10n"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
	"sunnyvaleserv.org/portal/util/sendmail"
)

const registerPersonFields = person.FID | person.FInformalName | person.FSortName | person.FEmail | person.FEmail2 | person.FCellPhone | person.FCallSign

// HandleRegister handles /classes/$id/register requests.
func HandleRegister(r *request.Request, cidstr string) {
	var (
		user     *person.Person
		c        *class.Class
		regs     []*classreg.ClassReg
		uregs    []*classreg.Updater
		errors   []string
		referral class.Referral
		others   uint
		forceGet bool
		max      = -1
	)
	// Get the user information.
	if user = auth.SessionUser(r, registerPersonFields, false); user == nil {
		if user = handleRegisterNotLoggedIn(r, cidstr); user == nil {
			return
		}
		forceGet = true
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	// Get the class information and the current registrations by this user.
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), class.UpdaterFields); c == nil {
		errpage.NotFound(r, user)
		return
	}
	classreg.AllForClass(r, c.ID(), classreg.UpdaterFields, func(cr *classreg.ClassReg) {
		if cr.RegisteredBy() != user.ID() {
			others++
			return
		}
		regs = append(regs, cr.Clone())
	})
	// Determine how many people this user can register.
	if c.Limit() != 0 {
		max = int(c.Limit() - others)
		if max < len(regs) {
			max = len(regs)
		}
	}
	// Determine what to display in the form.
	if r.Method == http.MethodPost && !forceGet {
		uregs, errors, referral = readForm(r, max)
		if len(errors) == 0 {
			applyForm(r, user, c, regs, uregs, referral)
			return
		}
	} else {
		uregs = make([]*classreg.Updater, len(regs))
		for i, reg := range regs {
			uregs[i] = reg.Updater(r, c, nil, user)
		}
		if len(regs) == 0 {
			uregs = append(uregs, &classreg.Updater{
				Class:        c,
				Person:       user,
				RegisteredBy: user,
				FirstName:    personFirstName(user),
				LastName:     personLastName(user),
				Email:        personEmail(user),
				CellPhone:    user.CellPhone(),
			})
			referral = class.Referral(99) // a nonzero invalid value
		}
		uregs = append(uregs, new(classreg.Updater))
	}
	r.HTMLNoCache()
	if len(errors) != 0 || forceGet {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=main")
	if max > 0 {
		form.A("data-max=%d", max)
	}
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if max == 0 {
		form.E("div").R(r.Loc("This class is now full.  You will be placed on a waiting list for the class and will be notified if space becomes available."))
	}
	for i := 0; i < len(uregs); i++ {
		var err string
		if i < len(errors) {
			err = errors[i]
		}
		emitRow(r, form, uregs[i], err, i)
	}
	if referral != 0 {
		emitReferral(r, form, referral)
	}
	emitButtons(r, form)
}

func emitRow(r *request.Request, form *htmlb.Element, reg *classreg.Updater, err string, idx int) {
	div := form.E("div class='formRow-3col classregDivider'", idx == 0, "class=first")
	div.E("div").TF(r.Loc("Student %d"), idx+1)
	div.E("button type=button class='sbtn sbtn-xsmall sbtn-danger classregClear' data-row=%d>%s", idx, r.Loc("Clear"))
	row := form.E("div class=formRow")
	row.E("label for=classregFirstname%d>%s", idx, r.Loc("Name"))
	names := row.E("div class='formInput classregNames'")
	names.E("input id=classregFirstname%d name=firstName class='formInput classregFirstname' placeholder=%s value=%s", idx, r.Loc("First"), reg.FirstName)
	names.E("input id=classregLastname%d name=lastName class='formInput classregLastname' placeholder=%s value=%s", idx, r.Loc("Last"), reg.LastName)
	row = form.E("div class=formRow")
	row.E("label for=classregEmail%d>%s", idx, r.Loc("Email"))
	row.E("input id=classregEmail%d name=email class='formInput classregEmail' value=%s", idx, reg.Email)
	row = form.E("div class=formRow")
	row.E("label for=classregCellPhone%d>%s", idx, r.Loc("Cell Phone"))
	row.E("input id=classregCellPhone%d name=cellPhone class='formInput classregCellPhone' value=%s", idx, reg.CellPhone)
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func personFirstName(p *person.Person) string {
	parts := strings.SplitN(p.SortName(), ",", 2)
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return ""
}

func personLastName(p *person.Person) string {
	parts := strings.Split(p.SortName(), ",")
	return parts[0]
}

func personEmail(p *person.Person) string {
	if p.Email() != "" {
		return p.Email()
	}
	return p.Email2()
}

func emitReferral(r *request.Request, form *htmlb.Element, referral class.Referral) {
	row := form.E("div class='formRow-3col classregReferral'")
	row.E("label for=classregReferral>%s", r.Loc("How did you find out about this class?"))
	sel := row.E("select id=classregReferral name=referral class=formInput")
	if !referral.Valid() {
		sel.E("option value='' selected>%s", r.Loc("(select one)"))
	}
	for _, ref := range class.AllReferrals {
		sel.E("option value=%d", ref, ref == referral, "selected").T(r.Loc(ref.String()))
	}
}

func emitButtons(r *request.Request, form *htmlb.Element) {
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Sign Up"))
}

var emailRE = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

func readForm(r *request.Request, max int) (uregs []*classreg.Updater, errors []string, referral class.Referral) {
	for i := range r.Form["firstName"] {
		if i >= len(r.Form["lastName"]) || i >= len(r.Form["email"]) || i >= len(r.Form["cellPhone"]) {
			break
		}
		first, last, email, cellPhone := r.Form["firstName"][i], r.Form["lastName"][i], r.Form["email"][i], r.Form["cellPhone"][i]
		if first == "" && last == "" && email == "" && cellPhone == "" {
			continue
		}
		var err string
		if first == "" || last == "" {
			err = r.Loc("Both first and last name are required. ")
		} else {
			for _, ur := range uregs {
				if ur.FirstName == first && ur.LastName == last {
					err = r.Loc("Each student must have a different name. ")
					break
				}
			}
		}
		if email != "" && !emailRE.MatchString(email) {
			err += r.Loc("The email address is not valid. ")
		}
		if !fmtPhone(&cellPhone) {
			err += r.Loc("The cell phone number is not valid.") + " "
		}
		if max > 0 && len(uregs) >= max {
			err += r.Loc("The class does not have this many spaces left.")
		}
		if err != "" {
			for len(errors) < len(uregs) {
				errors = append(errors, "")
			}
			errors = append(errors, err)
		}
		uregs = append(uregs, &classreg.Updater{FirstName: first, LastName: last, Email: email, CellPhone: cellPhone})
	}
	if _, ok := r.Form["referral"]; ok {
		if referral = class.Referral(util.ParseID(r.FormValue("referral"))); !referral.Valid() {
			referral = class.Referral(99) // a nonzero invalid value
		}
	}
	return uregs, errors, referral
}

func fmtPhone(p *string) bool {
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
	return false
}

func applyForm(
	r *request.Request, user *person.Person, c *class.Class, regs []*classreg.ClassReg,
	uregs []*classreg.Updater, referral class.Referral,
) {
	// Determine adds, cancels, and changes.
	adds, changesTo, changesFrom, cancels := splitForm(regs, uregs, user, c)
	// Save the changes.
	saveRegistrations(r, adds, changesTo, changesFrom, cancels, c, referral)
	// Send the confirmation emails.  (We don't send confirmation emails for
	// changes, only for adds and cancels.)
	sendUserConfirmation(r, user, c, uregs, cancels)
	sendAddConfirmations(r, user, c, adds)
	sendCancelConfirmations(r, user, c, cancels)
	// Show the confirmation dialog.
	showConfirmation(r, user, uregs, cancels)
}

// splitForm compares the registrations submitted in the form with those already
// present and returns lists of adds, changes, and cancels.
func splitForm(
	regs []*classreg.ClassReg, uregs []*classreg.Updater, user *person.Person, c *class.Class,
) (adds, changesTo []*classreg.Updater, changesFrom, cancels []*classreg.ClassReg) {
	for _, ur := range uregs {
		var found bool
		for i, r := range regs {
			if r != nil && r.FirstName() == ur.FirstName && r.LastName() == ur.LastName {
				if r.Email() != ur.Email || r.CellPhone() != ur.CellPhone {
					ur.ID, ur.RegisteredBy, ur.Class = r.ID(), user, c
					changesTo = append(changesTo, ur)
					changesFrom = append(changesFrom, r)
				}
				found = true
				regs[i] = nil
				break
			}
		}
		if !found {
			ur.RegisteredBy, ur.Class = user, c
			adds = append(adds, ur)
		}
	}
	for _, r := range regs {
		if r != nil {
			cancels = append(cancels, r)
		}
	}
	return
}

func saveRegistrations(
	r *request.Request, adds, changesTo []*classreg.Updater, changesFrom, cancels []*classreg.ClassReg, c *class.Class,
	referral class.Referral,
) {
	r.Transaction(func() {
		for i, to := range changesTo {
			changesFrom[i].Update(r, to)
		}
		for _, can := range cancels {
			can.Delete(r, c)
		}
		for _, add := range adds {
			classreg.Create(r, add)
		}
		if referral.Valid() {
			uc := c.Updater()
			uc.Referrals[referral]++
			c.Update(r, uc)
		}
	})
}

func sendUserConfirmation(r *request.Request, user *person.Person, c *class.Class, uregs []*classreg.Updater, cancels []*classreg.ClassReg) {
	var (
		toaddrs []string
		recips  []string
		body    bytes.Buffer
		desc    []string
	)
	if user.Email() != "" {
		recips = append(recips, user.Email())
		toaddrs = append(toaddrs, (&mail.Address{Name: user.InformalName(), Address: user.Email()}).String())
	}
	if user.Email2() != "" {
		recips = append(recips, user.Email2())
		toaddrs = append(toaddrs, (&mail.Address{Name: user.InformalName(), Address: user.Email2()}).String())
	}
	recips = append(recips, config.Get("adminEmail"))
	fmt.Fprintf(&body, "From: %s\r\nTo: %s\r\n", config.Get("fromEmail"), strings.Join(toaddrs, ", "))
	fmt.Fprintf(&body, "Bcc: %s\r\n", config.Get("adminEmail"))
	fmt.Fprintf(&body, "Subject: %s: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n", r.Loc(c.Type().String()), r.Loc("Class Registration"))
	fmt.Fprintf(&body, r.Loc("Greetings, %s,"), user.InformalName())
	fmt.Fprint(&body, "\r\n\r\n")
	fmt.Fprintf(&body, r.Loc("Thank you for your interest in our “%s” class:"), r.Loc(c.Type().String()))
	if r.Language == "es" {
		desc = strings.Split(c.EsDesc(), "\n")
	} else {
		desc = strings.Split(c.EnDesc(), "\n")
	}
	for _, line := range desc {
		fmt.Fprint(&body, "\r\n    ", line)
	}
	if len(uregs) != 0 {
		fmt.Fprint(&body, "\r\n\r\n")
		if len(uregs) > 1 {
			fmt.Fprint(&body, r.Loc("We confirm the registrations of:"))
		} else {
			fmt.Fprint(&body, r.Loc("We confirm the registration of:"))
		}
		for _, ur := range uregs {
			fmt.Fprintf(&body, "\r\n    %s %s", ur.FirstName, ur.LastName)
			if ur.Email != "" {
				fmt.Fprintf(&body, ", %s", ur.Email)
			}
			if ur.CellPhone != "" {
				fmt.Fprintf(&body, ", %s", ur.CellPhone)
			}
		}
	}
	if len(cancels) != 0 {
		fmt.Fprint(&body, "\r\n\r\n")
		if len(cancels) > 1 {
			fmt.Fprint(&body, r.Loc("You have canceled the registrations of:"))
		} else {
			fmt.Fprint(&body, r.Loc("You have canceled the registration of:"))
		}
		for _, can := range cancels {
			fmt.Fprintf(&body, "\r\n    %s %s", can.FirstName(), can.LastName())
			if can.Email() != "" {
				fmt.Fprintf(&body, ", %s", can.Email())
			}
			if can.CellPhone() != "" {
				fmt.Fprintf(&body, ", %s", can.CellPhone())
			}
		}
	}
	if len(uregs) != 0 {
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("If you need to withdraw from the class or make other changes, please return to SunnyvaleSERV.org.  You may also reply to this email."))
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("We look forward to seeing you!"))
		fmt.Fprint(&body, "\r\nSunnyvale SERV\r\nserv@sunnyvale.ca.gov\r\n")
	} else {
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("We hope to be able to accommodate you at some future class."))
		fmt.Fprint(&body, "\r\n\r\nSunnyvale SERV\r\nserv@sunnyvale.ca.gov\r\n")
	}
	if err := sendmail.SendMessage(r.Context(), config.Get("fromAddr"), recips, body.Bytes()); err != nil {
		r.LogEntry.Problems.AddError(err)
	}
}

func sendAddConfirmations(r *request.Request, user *person.Person, c *class.Class, adds []*classreg.Updater) {
	for _, add := range adds {
		var (
			body bytes.Buffer
			name string
			desc []string
		)
		if add.Email == "" || add.Email == user.Email() || add.Email == user.Email2() {
			continue
		}
		name = fmt.Sprintf("%s %s", add.FirstName, add.LastName)
		fmt.Fprintf(&body, "From: %s\r\nTo: %s\r\n", config.Get("fromEmail"), (&mail.Address{Name: name, Address: add.Email}).String())
		fmt.Fprintf(&body, "Subject: %s: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n", r.Loc(c.Type().String()), r.Loc("Class Registration"))
		fmt.Fprintf(&body, r.Loc("Greetings, %s,"), name)
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprintf(&body, r.Loc("%s has registered you for our “%s” class:"), user.InformalName(), r.Loc(c.Type().String()))
		if r.Language == "es" {
			desc = strings.Split(c.EsDesc(), "\n")
		} else {
			desc = strings.Split(c.EnDesc(), "\n")
		}
		for _, line := range desc {
			fmt.Fprint(&body, "\r\n    ", line)
		}
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("If this is incorrect, or you need to withdraw from the class, please reply to this email and let us know."))
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("We look forward to seeing you!"))
		fmt.Fprint(&body, "\r\nSunnyvale SERV\r\nserv@sunnyvale.ca.gov\r\n")
		if err := sendmail.SendMessage(r.Context(), config.Get("fromAddr"), []string{add.Email}, body.Bytes()); err != nil {
			r.LogEntry.Problems.AddError(err)
		}
	}
}

func sendCancelConfirmations(r *request.Request, user *person.Person, c *class.Class, cancels []*classreg.ClassReg) {
	for _, can := range cancels {
		var (
			body bytes.Buffer
			name string
			desc []string
		)
		if can.Email() == "" || can.Email() == user.Email() || can.Email() == user.Email2() {
			continue
		}
		name = fmt.Sprintf("%s %s", can.FirstName(), can.LastName())
		fmt.Fprintf(&body, "From: %s\r\nTo: %s\r\n", config.Get("fromEmail"), (&mail.Address{Name: name, Address: can.Email()}).String())
		fmt.Fprintf(&body, "Subject: %s: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n", r.Loc(c.Type().String()), r.Loc("Class Registration"))
		fmt.Fprintf(&body, r.Loc("Greetings, %s,"), name)
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprintf(&body, r.Loc("%s has canceled your registration for our “%s” class:"), user.InformalName(), r.Loc(c.Type().String()))
		if r.Language == "es" {
			desc = strings.Split(c.EsDesc(), "\n")
		} else {
			desc = strings.Split(c.EnDesc(), "\n")
		}
		for _, line := range desc {
			fmt.Fprint(&body, "\r\n    ", line)
		}
		fmt.Fprint(&body, "\r\n\r\n")
		fmt.Fprint(&body, r.Loc("If this is incorrect, please reply to this email and let us know."))
		fmt.Fprint(&body, "\r\n\r\nSunnyvale SERV\r\nserv@sunnyvale.ca.gov\r\n")
		if err := sendmail.SendMessage(r.Context(), config.Get("fromAddr"), []string{can.Email()}, body.Bytes()); err != nil {
			r.LogEntry.Problems.AddError(err)
		}
	}
}

func showConfirmation(r *request.Request, user *person.Person, uregs []*classreg.Updater, cancels []*classreg.ClassReg) {
	var (
		recips    []string
		reciplist string
	)

	if user.Email() != "" {
		recips = append(recips, user.Email())
	}
	if user.Email2() != "" {
		recips = append(recips, user.Email2())
	}
	reciplist = l10n.Conjoin(recips, "and", r.Language)
	r.HTMLNoCache()
	r.WriteHeader(http.StatusUnprocessableEntity)
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' up-main")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
	row := form.E("div class=formRow-3col")
	if len(uregs) > 1 && len(cancels) == 0 {
		row.E("p").E("b").T(r.Loc("Thank you!  Your class registrations are confirmed."))
	} else if len(uregs) == 1 && len(cancels) == 0 {
		row.E("p").E("b").T(r.Loc("Thank you!  Your class registration is confirmed."))
	} else if len(uregs) == 0 && len(cancels) > 1 {
		row.E("p").E("b").T(r.Loc("Thank you!  Your class registrations are canceled."))
	} else if len(uregs) == 0 && len(cancels) == 1 {
		row.E("p").E("b").T(r.Loc("Thank you!  Your class registration is canceled."))
	} else {
		row.E("p").E("b").T(r.Loc("Thank you!  Your changes have been saved."))
	}
	if len(uregs) != 0 {
		row.E("p").TF(r.Loc("A confirmation message has been sent to %s. If you don’t receive it promptly, look for it in your Junk Mail folder. Move it to your inbox so that future messages from us about the class are not marked as Junk Mail."), reciplist)
		row.E("p").T(r.Loc("If you need to withdraw from the class, please return to this website and remove your registration.  You may also send email to serv@sunnyvale.ca.gov."))
	} else {
		row.E("p").TF(r.Loc("A confirmation message has been sent to %s."), reciplist)
	}
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-primary' up-dismiss>OK")
}
