package authn

import (
	"net/http"
	"time"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// Maximum bad login attempts before lockout
const maxBadLogins = 3

// Threshold time for bad login attempts
const badLoginThreshold = 20 * time.Minute

// GetLogin handles GET /api/login requests.
func GetLogin(r *util.Request) error {
	var out jwriter.Writer
	var sender bool
	out.RawString(`{"id":`)
	out.Int(int(r.Person.ID))
	out.RawString(`,"informalName":`)
	out.String(r.Person.InformalName)
	out.RawString(`,"webmaster":`)
	out.Bool(r.Person.Roles[model.Webmaster])
	out.RawString(`,"canAddEvents":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader))
	out.RawString(`,"canAddPeople":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader))
	for _, l := range r.Tx.FetchLists() {
		if l.Type == model.ListSMS && l.People[r.Person.ID]&model.ListSender != 0 {
			sender = true
			break
		}
	}
	out.RawString(`,"canSendTextMessages":`)
	out.Bool(sender)
	out.RawString(`,"canViewReports":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivLeader))
	out.RawString(`,"canViewRosters":`)
	out.Bool(r.Person.HasPrivLevel(model.PrivMember))
	out.RawString(`,"csrf":`)
	out.String(string(r.Session.CSRF))
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json")
	out.DumpTo(r)
	return nil
}

// PostLogin handles POST /api/login requests.
func PostLogin(r *util.Request) error {
	var (
		person   *model.Person
		username = r.FormValue("username")
		password = r.FormValue("password")
		remember = r.FormValue("remember") != ""
	)
	// Check that the login is valid.
	if person = r.Tx.FetchPersonByUsername(username); person == nil {
		goto FAIL // no person with that username
	}
	if person.ID != model.AdminPersonID { // admin cannot be disabled or locked out
		if person.Roles[model.DisabledUser] {
			goto FAIL // person is disabled
		}
		if !person.HasPrivLevel(model.PrivStudent) {
			goto FAIL // person belongs to no orgs
		}
		if person.BadLoginCount >= maxBadLogins && time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
			goto FAIL // locked out
		}
	} else { // admin can not be remembered
		remember = false
	}
	if !CheckPassword(person, password) {
		goto FAIL // password mismatch
	}
	// The login is valid.  Record it and create a session.
	r.Person = person
	if person.BadLoginCount > 0 {
		r.Tx.WillUpdatePerson(person)
		person.BadLoginCount = 0
		r.Tx.UpdatePerson(person)
	}
	util.CreateSession(r, remember)
	return GetLogin(r)

FAIL:
	if person != nil {
		// Record the bad login attempt.
		r.Tx.WillUpdatePerson(person)
		if time.Now().Before(person.BadLoginTime.Add(badLoginThreshold)) {
			person.BadLoginCount++
		} else {
			person.BadLoginCount = 1
		}
		person.BadLoginTime = time.Now()
		r.Tx.UpdatePerson(person)
	}
	r.Tx.Commit()
	return util.HTTPError(http.StatusUnauthorized, "401 Unauthorized")
}

// PostLogout handles POST /api/logout requests.
func PostLogout(r *util.Request) error {
	if c, err := r.Cookie("auth"); err == nil {
		r.Tx.DeleteSession(&model.Session{Token: model.SessionToken(c.Value), Person: r.Person})
		http.SetCookie(r, &http.Cookie{
			Name:   "auth",
			Value:  c.Value,
			Path:   "/",
			MaxAge: -1,
		})
	}
	r.Tx.Commit()
	return nil
}
