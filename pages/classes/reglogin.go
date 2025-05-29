package classes

import (
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/login"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui/form"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

type registerLogIn struct {
	r         *request.Request
	notify    bool
	f         form.Form
	email     string
	haveEmail bool
	password  string
	user      *person.Person
	firstName string
	lastName  string
	cellPhone string
	newpwd1   string
	newpwd2   string
	loggedIn  bool
}

func handleRegisterNotLoggedIn(r *request.Request, notify bool) *person.Person {
	login := registerLogIn{r: r, notify: notify}
	login.createForm()
	login.f.Handle(r)
	if login.loggedIn {
		return login.user
	}
	return nil
}

func (li *registerLogIn) createForm() {
	li.f.Attrs = "method=POST up-target=form"
	li.f.Dialog = true
	if li.notify {
		li.f.Title = "Class Notifications"
	} else {
		li.f.Title = "Class Registration"
	}
	li.f.Buttons = []*form.Button{{Label: "Submit", OnClick: func() bool { return false }}}
	li.f.Rows = []form.Row{
		&introRow{login: li},
		&emailRow{
			TextInputRow: form.TextInputRow{
				LabeledRow: form.LabeledRow{
					RowID: "classregLoginEmail",
					Label: "Email",
				},
				Name:   "email",
				ValueP: &li.email,
			},
			login: li,
		},
		&passwordMessageRow{
			MessageRow: form.MessageRow{
				LabeledRow: form.LabeledRow{Label: " "},
				HTML:       `Welcome back!  We have an account on this system for that email address.  Please enter the corresponding password.`,
			},
			login: li,
		},
		&passwordRow{
			PasswordRow: form.PasswordRow{
				InputRow: form.InputRow{
					LabeledRow: form.LabeledRow{
						RowID: "classregLoginPassword",
						Label: "Password",
					},
					Name:   "password",
					ValueP: &li.password,
				},
				Autocomplete: "current-password",
			},
			login: li,
		},
		&passwordMessageRow{
			MessageRow: form.MessageRow{
				LabeledRow: form.LabeledRow{Label: " "},
				HTML:       `If you don’t remember your password, you can <a href="/password-reset">reset it</a> and we’ll email you a new one.`,
			},
			login: li,
		},
		&createMessageRow{login: li},
		&namesRow{
			LabeledRow: form.LabeledRow{
				RowID: "classregLoginNames",
				Label: "Name",
			},
			login: li,
		},
		&cellPhoneRow{
			TextInputRow: form.TextInputRow{
				LabeledRow: form.LabeledRow{
					RowID: "classregLoginCellphone",
					Label: "Cell Phone",
					Help:  "The cell phone is used only for urgent notifications, such as last-minute cancellation of a class.  It is optional.",
				},
				Name:   "cellPhone",
				ValueP: &li.cellPhone,
			},
			login: li,
		},
		&newPasswordRow{
			NewPasswordPairRow: login.NewPasswordPairRow{
				LabeledRow: form.LabeledRow{
					RowID: "classregLoginNewpwd",
					Label: "Password",
				},
				Name:    "newpwd",
				ValueP1: &li.newpwd1,
				ValueP2: &li.newpwd2,
			},
			login: li,
		},
	}
}

type introRow struct {
	form.BaseRow
	login *registerLogIn
}

func (ir *introRow) Read(r *request.Request) bool { return true }
func (ir *introRow) Emit(r *request.Request, parent *htmlb.Element, _ bool) {
	div := parent.E("div class=formRow-3col")
	if ir.login.notify {
		if !ir.login.haveEmail {
			div.T(r.Loc("To subscribe to notifications of new classes, please enter your email address."))
		} else if ir.login.user == nil {
			div.T(r.Loc("To subscribe to notifications of new classes, please create an account."))
		} else {
			div.T(r.Loc("To subscribe to notifications of new classes, please log in."))
		}
	} else {
		if !ir.login.haveEmail {
			div.T(r.Loc("To register for this class, please enter your email address."))
		} else if ir.login.user == nil {
			div.T(r.Loc("To register for this class, please create an account."))
		} else {
			div.T(r.Loc("To register for this class, please log in."))
		}
	}
}

type emailRow struct {
	form.TextInputRow
	login *registerLogIn
}

func (er *emailRow) ReadOrder() int { return -1 }
func (er *emailRow) Read(r *request.Request) bool {
	// First, make sure we have an email address.  If not, we're asking for
	// one (and nothing else).
	if !er.TextInputRow.Read(r) {
		// nothing
	} else if *er.ValueP == "" {
		er.Error = r.Loc("Your email address is required.")
	} else if !emailRE.MatchString(*er.ValueP) {
		er.Error = r.Loc("This is not a valid email address.")
	} else {
		er.login.haveEmail = true
	}
	if !er.login.haveEmail {
		er.login.f.Buttons[0].Label = "Submit"
		er.login.f.Buttons[0].OnClick = er.login.submitEmail
		return false
	}
	er.login.user = person.WithEmail(r, *er.ValueP, notifyPersonFields|person.FBadLoginCount|person.FBadLoginTime|person.FPassword)
	if er.login.user == nil {
		er.login.f.Buttons[0].Label = "Create Account"
		er.login.f.Buttons[0].OnClick = er.login.createAccount
	} else {
		er.login.f.Buttons[0].Label = "Login"
		er.login.f.Buttons[0].OnClick = er.login.logIn
	}
	return true
}

type passwordMessageRow struct {
	form.MessageRow
	login *registerLogIn
}

func (pmr *passwordMessageRow) ShouldEmit(_ request.ValidationList) bool {
	return pmr.login.user != nil
}

type passwordRow struct {
	form.PasswordRow
	login *registerLogIn
}

func (pr *passwordRow) ShouldEmit(vl request.ValidationList) bool {
	return pr.login.user != nil
}

func (pr *passwordRow) Read(r *request.Request) bool {
	valid := true

	if pr.login.user == nil {
		return true
	}
	if !pr.PasswordRow.Read(r) {
		return false
	}
	if *pr.ValueP == "" {
		if _, ok := r.Form["password"]; ok { // don't display error when password field first shown
			pr.Error = r.Loc("Your password is required.")
		}
		return false
	}
	if pr.login.user.ID() != person.AdminID { // admin cannot be disabled or locked out
		if pr.login.user.BadLoginCount() >= maxBadLogins && time.Now().Before(pr.login.user.BadLoginTime().Add(badLoginThreshold)) {
			valid = false // locked out
		} else if held, _ := personrole.PersonHasRole(r, pr.login.user.ID(), role.Disabled); held {
			valid = false // disabled user
		}
	}
	if valid && !auth.CheckPassword(r, pr.login.user, r.FormValue("password")) {
		valid = false // wrong password
	}
	if !valid {
		pr.Error = r.Loc("Login incorrect. Please try again.")
		return false
	}
	return true
}

type createMessageRow struct {
	form.BaseRow
	login *registerLogIn
}

func (cmr *createMessageRow) Read(r *request.Request) bool {
	return true
}

func (cmr *createMessageRow) ShouldEmit(_ request.ValidationList) bool {
	return cmr.login.haveEmail && cmr.login.user == nil
}

func (cmr *createMessageRow) Emit(r *request.Request, parent *htmlb.Element, _ bool) {
	parent.E("div class=formRow-3col").T(r.Loc("We do not have an account with this email address.  To create a new account, please provide the following information."))
}

type namesRow struct {
	form.LabeledRow
	login *registerLogIn
}

func (nr *namesRow) Read(r *request.Request) bool {
	nr.Error = ""
	if !nr.login.haveEmail || nr.login.user != nil {
		return true
	}
	nr.login.firstName = strings.TrimSpace(r.FormValue("firstName"))
	nr.login.lastName = strings.TrimSpace(r.FormValue("lastName"))
	if nr.login.firstName == "" || nr.login.lastName == "" {
		if _, ok := r.Form["firstName"]; ok {
			nr.Error = r.Loc("Your name is required.")
		}
		return false
	}
	return true
}

func (nr *namesRow) ShouldEmit(_ request.ValidationList) bool {
	return nr.login.haveEmail && nr.login.user == nil
}

func (nr *namesRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	var focusID string
	if nr.RowID != "" {
		focusID = nr.RowID + "-in"
	}
	row := nr.EmitPrefix(r, parent, focusID)
	names := row.E("div class='formInput classregNames'")
	names.E("input id=classregLoginFirstname name=firstName class=formInput placeholder=%s value=%s autofocus", r.Loc("First"), nr.login.firstName)
	names.E("input id=classregLoginLastname name=lastName class=formInput placeholder=%s value=%s", r.Loc("Last"), nr.login.lastName)
	nr.EmitSuffix(r, row)
}

type cellPhoneRow struct {
	form.TextInputRow
	login *registerLogIn
}

func (cpr *cellPhoneRow) Read(r *request.Request) bool {
	if !cpr.login.haveEmail || cpr.login.user != nil {
		return true
	}
	if !cpr.TextInputRow.Read(r) {
		return false
	}
	if !fmtPhone(cpr.ValueP) {
		cpr.Error = r.Loc("The cell phone number is not valid.")
		return false
	}
	return true
}

func (cpr *cellPhoneRow) ShouldEmit(_ request.ValidationList) bool {
	return cpr.login.haveEmail && cpr.login.user == nil
}

type newPasswordRow struct {
	login.NewPasswordPairRow
	login *registerLogIn
}

func (npr *newPasswordRow) Read(r *request.Request) bool {
	if !npr.login.haveEmail || npr.login.user != nil {
		return true
	}
	if _, ok := r.Form["newpwd-1"]; !ok {
		npr.Score = -1 // inhibit meter
		return false
	}
	return npr.NewPasswordPairRow.Read(r)
}

func (npr *newPasswordRow) ShouldEmit(_ request.ValidationList) bool {
	return npr.login.haveEmail && npr.login.user == nil
}

func (login *registerLogIn) submitEmail() bool {
	return false
}

func (login *registerLogIn) logIn() bool {
	login.r.Transaction(func() {
		if login.user.BadLoginCount() > 0 {
			up := login.user.Updater()
			up.BadLoginCount = 0
			up.BadLoginTime = time.Time{}
			login.user.Update(login.r, up, person.FBadLoginCount|person.FBadLoginTime)
		}
		auth.CreateSession(login.r, login.user, false)
	})
	login.r.Form.Set("csrf", login.r.CSRF)
	login.loggedIn = true
	return true
}

func (login *registerLogIn) createAccount() bool {
	login.r.Transaction(func() {
		up := &person.Updater{
			InformalName: login.firstName + " " + login.lastName,
			FormalName:   login.firstName + " " + login.lastName,
			SortName:     login.lastName + ", " + login.firstName,
			Email:        login.email,
			CellPhone:    login.cellPhone,
		}
		login.user = person.Create(login.r, up)
		auth.SetPassword(login.r, login.user, login.newpwd1)
		auth.CreateSession(login.r, login.user, false)
	})
	login.r.Form.Set("csrf", login.r.CSRF)
	login.loggedIn = true
	return true
}
