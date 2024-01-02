package classes

import (
	"net/http"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

func handleRegisterNotLoggedIn(r *request.Request, cidstr string) (user *person.Person) {
	var (
		email string
		valid bool
	)
	if r.Method != http.MethodPost {
		r.HTMLNoCache()
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form class='form form-2col' method=POST up-main up-target=form")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class='formRow-3col classregLoginIntro'").T(r.Loc("To register for this class, please enter your email address."))
		row := form.E("div class=formRow")
		row.E("label for=classregLoginEmail").T(r.Loc("Email"))
		row.E("input id=classregLoginEmail name=email class=formInput")
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
		buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Submit"))
		return
	}
	email = strings.TrimSpace(r.FormValue("email"))
	if !emailRE.MatchString(email) {
		r.HTMLNoCache()
		r.WriteHeader(http.StatusUnprocessableEntity)
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form class='form form-2col' method=POST up-main up-target=form")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class='formRow-3col classregLoginIntro'").T(r.Loc("To register for this class, please enter your email address."))
		row := form.E("div class=formRow")
		row.E("label for=classregLoginEmail").T(r.Loc("Email"))
		row.E("input id=classregLoginEmail name=email class=formInput value=%s", email)
		if email == "" {
			row.E("div class=formError").T(r.Loc("Your email address is required."))
		} else {
			row.E("div class=formError").T(r.Loc("This is not a valid email address."))
		}
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
		buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Submit"))
		return
	}
	if user = person.WithEmail(r, email, registerPersonFields|person.FBadLoginCount|person.FBadLoginTime|person.FPassword); user == nil {
		return handleCreateAccount(r, email)
	}
	if _, ok := r.Form["password"]; !ok {
		r.HTMLNoCache()
		r.WriteHeader(http.StatusUnprocessableEntity)
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form class='form form-2col' method=POST up-main up-target=form")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class='formRow-3col classregLoginIntro'").T(r.Loc("To register for this class, please log in."))
		row := form.E("div class=formRow")
		row.E("label for=classregLoginEmail").T(r.Loc("Email"))
		row.E("input id=classregLoginEmail name=email class=formInput value=%s", email)
		row = form.E("div class=formRow")
		row.E("label for=classregLoginPassword").T(r.Loc("Password"))
		row.E("input type=password id=classregLoginPassword name=password autocomplete=password autocapitalize=none autofocus")
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
		buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Login"))
		return nil
	}
	valid = true
	if user.ID() != person.AdminID { // admin cannot be disabled or locked out
		if user.BadLoginCount() >= maxBadLogins && time.Now().Before(user.BadLoginTime().Add(badLoginThreshold)) {
			valid = false // locked out
		}
		if held, _ := personrole.PersonHasRole(r, user.ID(), role.Disabled); held {
			valid = false // person is disabled
		}
	}
	if valid && !auth.CheckPassword(r, user, r.FormValue("password")) {
		valid = false
	}
	if !valid {
		r.HTMLNoCache()
		r.WriteHeader(http.StatusUnprocessableEntity)
		html := htmlb.HTML(r)
		defer html.Close()
		form := html.E("form class='form form-2col' method=POST up-main up-target=form")
		form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
		form.E("div class='formRow-3col classregLoginIntro'").T(r.Loc("To register for this class, please log in."))
		row := form.E("div class=formRow")
		row.E("label for=classregLoginEmail").T(r.Loc("Email"))
		row.E("input id=classregLoginEmail name=email class=formInput value=%s", email)
		row = form.E("div class=formRow")
		row.E("label for=classregLoginPassword").T(r.Loc("Password"))
		row.E("input type=password id=classregLoginPassword name=password autocomplete=password autocapitalize=none autofocus")
		if r.FormValue("password") == "" {
			row.E("div class=formError").T(r.Loc("Your password is required."))
		} else {
			row.E("div class=formError").T(r.Loc("Login incorrect. Please try again."))
		}
		buttons := form.E("div class=formButtons")
		buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
		buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Login"))
		return nil
	}
	r.Transaction(func() {
		if user.BadLoginCount() > 0 {
			up := user.Updater()
			up.BadLoginCount = 0
			up.BadLoginTime = time.Time{}
			user.Update(r, up, person.FBadLoginCount|person.FBadLoginTime)
		}
		auth.CreateSession(r, user, false)
	})
	r.Form.Set("csrf", r.CSRF)
	return user
}

func handleCreateAccount(r *request.Request, email string) (user *person.Person) {
	var password = r.FormValue("password")
	var firstName = strings.TrimSpace(r.FormValue("firstName"))
	var lastName = strings.TrimSpace(r.FormValue("lastName"))
	var cellPhone = r.FormValue("cellPhone")
	var nameError, cellPhoneError string
	if _, ok := r.Form["password"]; ok {
		if password != "" && !auth.StrongPassword(nil, password) {
			password = "" // somehow a weak password got through
		}
		if firstName == "" || lastName == "" {
			nameError = r.Loc("Your name is required.")
		}
		if !fmtPhone(&cellPhone) {
			cellPhoneError = r.Loc("The cell phone number is not valid.")
		}
	}
	if password != "" && nameError == "" && cellPhoneError == "" {
		r.Transaction(func() {
			up := &person.Updater{
				InformalName: firstName + " " + lastName,
				FormalName:   firstName + " " + lastName,
				SortName:     lastName + ", " + firstName,
				Email:        email,
				CellPhone:    cellPhone,
			}
			user = person.Create(r, up)
			auth.SetPassword(r, user, password)
			auth.CreateSession(r, user, false)
		})
		r.Form.Set("csrf", r.CSRF)
		return user
	}
	r.HTMLNoCache()
	r.WriteHeader(http.StatusUnprocessableEntity)
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-target=form")
	form.E("div class='formTitle formTitle-primary'").R(r.Loc("Class Registration"))
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	form.E("div class='formRow-3col classregLoginIntro'").T(r.Loc("To register for this class, please create an account."))
	row := form.E("div class=formRow")
	row.E("label for=classregLoginEmail").T(r.Loc("Email"))
	row.E("input id=classregLoginEmail name=email class=formInput value=%s", email)
	form.E("div class=formRow-3col").T(r.Loc("We do not have an account with this email address.  To create a new account, please provide the following information."))
	row = form.E("div class=formRow")
	row.E("label for=classregNewFirstname>%s", r.Loc("Name"))
	names := row.E("div class='formInput classregNames'")
	names.E("input id=classregNewFirstname name=firstName class='formInput classregFirstname' placeholder=%s value=%s autofocus", r.Loc("First"), firstName)
	names.E("input id=classregNewLastname name=lastName class='formInput classregLastname' placeholder=%s value=%s", r.Loc("Last"), lastName)
	if nameError != "" {
		row.E("div class=formError>%s", nameError)
	}
	row = form.E("div class=formRow")
	row.E("label for=classregNewCellPhone>%s", r.Loc("Cell Phone"))
	row.E("input id=classregNewCellPhone name=cellPhone class='formInput classregCellPhone' value=%s", cellPhone)
	if cellPhoneError != "" {
		row.E("div class=formError>%s", cellPhoneError)
	}
	row.E("div class=formHelp").R(r.Loc("The cell phone is used only for urgent notifications, such as last-minute cancellation of a class.  It is optional."))
	row = form.E("div class=formRow")
	row.E("label for=classregNewPassword").T(r.Loc("Password"))
	row.E("s-password id=classregNewPassword name=password hints=%s", strings.Join(auth.StrongPasswordHints(nil), ","))
	buttons := form.E("div class=formButtons")
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>%s", r.Loc("Cancel"))
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=%s", r.Loc("Create Account"))
	return nil
}
