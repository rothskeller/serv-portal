package state

import (
	"net/http"
	"strconv"
	"time"

	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/request"
)

var (
	focusRole         role.ID
	peopleInitialized bool
)

// GetFocusRole returns the role ID of the last focus role selected on the
// people list or map page (extracted from a session cookie).  If no role has
// been focused, it returns zero.
func GetFocusRole(r *request.Request) role.ID {
	initializePeople(r)
	return focusRole
}

// SetFocusRole sets the focused role.
func SetFocusRole(r *request.Request, rid role.ID) {
	focusRole = rid
	http.SetCookie(r, &http.Cookie{
		Name:    "serv-people-role",
		Path:    "/",
		Value:   strconv.Itoa(int(rid)),
		Expires: time.Now().AddDate(1, 0, 0),
	})
}

// initializePeople ensures that focusRole has a value.
func initializePeople(r *request.Request) {
	if peopleInitialized {
		return
	}
	if c, err := r.Cookie("serv-people-role"); err == nil && c != nil {
		focusRole = role.ID(util.ParseID(c.Value))
	}
}
