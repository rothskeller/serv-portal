// Package server contains the implementation of the SERV portal server.
package server

import (
	"net/http"
	"path"
	"strings"

	"sunnyvaleserv.org/portal/pages/admin/listedit"
	"sunnyvaleserv.org/portal/pages/admin/listlist"
	"sunnyvaleserv.org/portal/pages/admin/listpeople"
	"sunnyvaleserv.org/portal/pages/admin/listrole"
	"sunnyvaleserv.org/portal/pages/admin/roleedit"
	"sunnyvaleserv.org/portal/pages/admin/rolelist"
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/events/eventattend"
	"sunnyvaleserv.org/portal/pages/events/eventcopy"
	"sunnyvaleserv.org/portal/pages/events/eventedit"
	"sunnyvaleserv.org/portal/pages/events/eventscal"
	"sunnyvaleserv.org/portal/pages/events/eventslist"
	"sunnyvaleserv.org/portal/pages/events/eventview"
	"sunnyvaleserv.org/portal/pages/events/proxysignup"
	"sunnyvaleserv.org/portal/pages/events/signups"
	"sunnyvaleserv.org/portal/pages/files"
	"sunnyvaleserv.org/portal/pages/files/docedit"
	"sunnyvaleserv.org/portal/pages/files/folderedit"
	"sunnyvaleserv.org/portal/pages/homepage"
	"sunnyvaleserv.org/portal/pages/login"
	"sunnyvaleserv.org/portal/pages/people/activity"
	"sunnyvaleserv.org/portal/pages/people/peoplelist"
	"sunnyvaleserv.org/portal/pages/people/peoplemap"
	"sunnyvaleserv.org/portal/pages/people/personedit"
	"sunnyvaleserv.org/portal/pages/people/personview"
	attrep "sunnyvaleserv.org/portal/pages/reports/attendance"
	clearrep "sunnyvaleserv.org/portal/pages/reports/clearance"
	"sunnyvaleserv.org/portal/pages/search"
	"sunnyvaleserv.org/portal/pages/static"
	"sunnyvaleserv.org/portal/pages/texts/textlist"
	"sunnyvaleserv.org/portal/pages/texts/textnew"
	"sunnyvaleserv.org/portal/pages/texts/textview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/log"
	"sunnyvaleserv.org/portal/util/request"
)

// server is the handler type for the SERV portal server.
type server struct{}

// Server is the singleton instance of the SERV portal server.  Before using
// this server, the current working directory must be the server data directory,
// and the database must have been opened.
var Server server

// ServeHTTP serves web requests to the SERV portal server.
func (s server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var (
		reqpath string
		err     error
	)
	reqpath = path.Clean("/" + r.URL.Path)
	// The only allowable methods are GET and POST.  Treat everything else
	// as a GET.
	if r.Method != http.MethodPost {
		r.Method = http.MethodGet
	}
	// Create the request structure.
	var req = &request.Request{
		Request:        r,
		ResponseWriter: request.NewUncompressedResponse(w),
		Path:           reqpath,
	}
	if strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
		req.ResponseWriter = request.NewCompressedResponse(req.ResponseWriter)
	}
	req.LogEntry = log.New(req.Method, req.Path)
	// Parse the form now, before we connect to the data store.  This avoids
	// chewing up a connection while we're still reading from the client.
	r.ParseMultipartForm(1048576)
	req.LogEntry.Params = req.Form
	// If we have a session token, save it.
	if c, err := req.Cookie("auth"); err == nil {
		req.LogEntry.Session = c.Value
		req.SessionToken = c.Value
	}
	// Disable caching of anything we send to the client.  In the few cases
	// where caching is OK, the specific request handler will undo this.
	w.Header().Set("Cache-Control", "no-store")
	// Allocate a connection to the data store and run the request in it.
	err = store.Connect(req.Context(), req.LogEntry, func(store *store.Store) {
		req.Store = store
		route(req)
	})
	// Flush anything not already sent to the client.
	if flush, ok := req.ResponseWriter.(http.Flusher); ok {
		flush.Flush()
	}
	// Add the response status, and any error, to the log entry.
	if err != nil {
		// The handler returned an error.  Log it and return an internal
		// server error to the client.
		req.LogEntry.Problems.AddError(err)
		req.LogEntry.Status = http.StatusInternalServerError
	} else if req.StatusCode() != 0 {
		req.LogEntry.Status = req.StatusCode()
	} else {
		req.LogEntry.Problems.Add("no output")
		req.LogEntry.Status = http.StatusInternalServerError
	}
	// If we haven't written anything to the client, send an internal server
	// error page.
	if req.StatusCode() == 0 {
		errpage.ServerError(req, auth.SessionUser(req, 0, false))
		if flush, ok := req.ResponseWriter.(http.Flusher); ok {
			flush.Flush()
		}
	}
	// Flush the log entry to disk.
	req.LogEntry.Log()
}

// router sends the request to the appropriate handler given its method and
// path.
func route(r *request.Request) {
	c := strings.Split(r.Path[1:], "/")
	for len(c) < 6 {
		c = append(c, "")
	}
	if r.Method != http.MethodPost {
		r.Method = http.MethodGet
	}
	switch {
	case c[0] == "":
		homepage.Serve(r)
	case c[0] == "assets" && c[1] != "" && c[2] == "":
		ui.ServeAsset(r, c[1])
	case c[0] == "about" && c[1] == "":
		static.AboutPage(r)
	case c[0] == "admin" && c[1] == "lists" && c[2] == "":
		listlist.Get(r)
	case c[0] == "admin" && c[1] == "lists" && c[2] != "" && c[3] == "":
		listedit.Handle(r, c[2])
	case c[0] == "admin" && c[1] == "lists" && c[2] != "" && c[3] != "" && c[4] == "":
		listpeople.Get(r, c[2], c[3])
	case c[0] == "admin" && c[1] == "lists" && c[2] != "" && c[3] == "roleedit" && c[4] != "" && c[5] == "":
		listrole.Get(r, c[2], c[4])
	case c[0] == "admin" && c[1] == "roles" && c[2] == "":
		rolelist.Get(r)
	case c[0] == "admin" && c[1] == "roles" && c[2] != "" && c[3] == "":
		roleedit.Handle(r, c[2])
	case c[0] == "docedit" && c[1] != "" && c[2] != "" && c[3] == "":
		docedit.Handle(r, c[1], c[2])
	case c[0] == "email-lists" && c[1] == "":
		static.EmailListsPage(r)
	case c[0] == "events" && c[1] == "attendance" && c[2] != "" && c[3] == "":
		eventattend.Handle(r, c[2])
	case c[0] == "events" && c[1] == "calendar" && c[2] != "" && c[3] == "":
		eventscal.Get(r, c[2])
	case c[0] == "events" && c[1] == "create" && c[2] == "":
		eventedit.HandleCreate(r)
	case c[0] == "events" && c[1] == "edshift" && c[2] != "" && c[3] == "":
		eventedit.HandleShift(r, c[2])
	case c[0] == "events" && c[1] == "edtask" && c[2] != "" && c[3] == "":
		eventedit.HandleTask(r, c[2])
	case c[0] == "events" && c[1] == "list" && c[2] != "" && c[3] == "":
		eventslist.Get(r, c[2])
	case c[0] == "events" && c[1] == "proxysignup" && c[2] != "" && c[3] == "":
		proxysignup.Handle(r, c[2])
	case c[0] == "events" && c[1] == "signups" && c[3] == "":
		signups.Handle(r, c[2])
	case c[0] == "events" && c[1] != "" && c[2] == "":
		eventview.Handle(r, c[1])
	case c[0] == "events" && c[1] != "" && c[2] == "copy" && c[3] == "":
		eventcopy.Handle(r, c[1])
	case c[0] == "events" && c[1] != "" && c[2] == "eddetails" && c[3] == "":
		eventedit.HandleDetails(r, c[1])
	// case c[0] == "events" && c[1] != "" && c[2] == "edit" && c[3] == "":
	// 	eventedit.Handle(r, c[1])
	case c[0] == "files":
		files.Handle(r)
	case c[0] == "folderedit" && c[1] != "" && c[2] == "":
		folderedit.Handle(r, c[1])
	case c[0] == "login":
		login.HandleLogin(r)
	case c[0] == "logout":
		login.HandleLogout(r)
	case c[0] == "password-reset" && c[1] == "":
		login.HandlePWReset(r)
	case c[0] == "password-reset" && c[1] != "" && c[2] == "":
		login.HandlePWResetToken(r, c[1])
	case c[0] == "people" && c[1] == "":
		peoplelist.Handle(r)
	case c[0] == "people" && c[1] == "map" && c[2] == "":
		peoplemap.Handle(r)
	case c[0] == "people" && c[1] != "" && c[2] == "":
		personview.Get(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "activity" && c[3] != "" && c[4] == "":
		activity.HandleActivity(r, c[1], c[3])
	case c[0] == "people" && c[1] != "" && c[2] == "edcontact" && c[3] == "":
		personedit.HandleContact(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "ednames" && c[3] == "":
		personedit.HandleNames(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "ednote" && c[4] == "":
		personedit.HandleNote(r, c[1], c[3])
	case c[0] == "people" && c[1] != "" && c[2] == "edpassword" && c[3] == "":
		personedit.HandlePassword(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "edroles" && c[3] == "":
		personedit.HandleRoles(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "edstatus" && c[3] == "":
		personedit.HandleStatus(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "edsubscriptions" && c[3] == "":
		personedit.HandleSubscriptions(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "pwreset" && c[3] == "":
		personedit.HandlePWReset(r, c[1])
	case c[0] == "people" && c[1] != "" && c[2] == "vregister" && c[3] == "":
		personedit.HandleVRegister(r, c[1])
	case c[0] == "reports" && c[1] == "attendance" && c[2] == "":
		attrep.Get(r)
	case c[0] == "reports" && c[1] == "clearance" && c[2] == "":
		clearrep.Get(r)
	case c[0] == "search" && c[1] == "":
		search.Handle(r)
	case c[0] == "subscribe-calendar" && c[1] == "":
		static.SubscribeCalendarPage(r)
	case c[0] == "texts" && c[1] == "":
		textlist.Get(r)
	case c[0] == "texts" && c[1] == "NEW" && c[2] == "":
		textnew.Handle(r)
	case c[0] == "texts" && c[1] != "" && c[2] == "":
		textview.Get(r, c[1])
	case c[0] == "volunteer-hours" && c[1] != "" && c[2] == "":
		activity.HandleVolunteerHours(r, c[1])
		/*
			case r.User == nil && c[0] == "api":
				http.Error(r, "401 Unauthorized", http.StatusUnauthorized)
			case r.User == nil:
				http.Redirect(r, r.Request, "/login"+r.Path, http.StatusSeeOther)
				/***** END OF PAGES THAT DON'T REQUIRE LOGIN *****/
	/*
		case c[0] == "api":
			http.Error(r, "404 Not Found", http.StatusNotFound)
	*/
	default:
		errpage.NotFound(r, auth.SessionUser(r, 0, false))
	}
}