package personedit

import (
	"fmt"
	"net/http"
	"slices"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/pages/people/personview"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const notePersonFields = person.FInformalName | person.FCallSign | person.FPrivLevels | person.FNotes

// HandleNote handles requests for /people/$id/ednote[/$index].
func HandleNote(r *request.Request, idstr, indexstr string) {
	var (
		user      *person.Person
		p         *person.Person
		nidx      int
		n         *person.Note
		dateError string
		noteError string
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if !auth.CheckCSRF(r, user) {
		return
	}
	if p = person.WithID(r, person.ID(util.ParseID(idstr)), notePersonFields|person.CanViewTargetFields); p == nil {
		errpage.NotFound(r, user)
		return
	}
	if !user.HasPrivLevel(0, enum.PrivLeader) {
		errpage.Forbidden(r, user)
		return
	}
	if indexstr != "" {
		if nidx = util.ParseID(indexstr); nidx >= 0 && nidx < len(p.Notes()) {
			n = p.Notes()[nidx]
		} else {
			errpage.NotFound(r, user)
			return
		}
		switch n.Visibility {
		case person.NoteVisibleToAdmins:
			if !user.IsAdminLeader() {
				errpage.Forbidden(r, user)
				return
			}
		case person.NoteVisibleToWebmaster:
			if !user.IsWebmaster() {
				errpage.Forbidden(r, user)
				return
			}
		}
	} else {
		nidx = -1
		n = &person.Note{Date: time.Now()}
		if user.IsWebmaster() {
			n.Visibility = person.NoteVisibleToWebmaster
		} else if user.IsAdminLeader() {
			n.Visibility = person.NoteVisibleToAdmins
		} else {
			n.Visibility = person.NoteVisibleToLeaders
		}
	}
	validate := strings.Fields(r.Request.Header.Get("X-Up-Validate"))
	if r.Method == http.MethodPost {
		if r.FormValue("delete") != "" && nidx >= 0 {
			r.Transaction(func() {
				up := p.Updater()
				up.Notes = append(up.Notes[:nidx], up.Notes[nidx+1:]...)
				p.Update(r, up, person.FNotes)
			})
			personview.Render(r, user, p, user.CanView(p), "notes")
			return
		}
		dateError = readNoteDate(r, n)
		noteError = readNoteText(r, n)
		readNoteVisibility(r, user, n)
		if len(validate) == 0 && dateError == "" && noteError == "" {
			r.Transaction(func() {
				up := p.Updater()
				if nidx >= 0 {
					up.Notes[nidx] = n
				} else {
					up.Notes = append(up.Notes, n)
				}
				p.Update(r, up, person.FNotes)
			})
			personview.Render(r, user, p, user.CanView(p), "notes")
			return
		}
	}
	r.HTMLNoCache()
	if dateError != "" || noteError != "" {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class='form form-2col' method=POST up-main up-layer=parent up-target=.personviewNotes")
	if nidx < 0 {
		form.E("div class='formTitle formTitle-primary'>Add Note")
	} else {
		form.E("div class='formTitle formTitle-primary'>Edit Note")
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if len(validate) == 0 || slices.Contains(validate, "date") {
		emitNoteDate(form, n, dateError != "" || noteError == "", dateError)
	}
	if len(validate) == 0 || slices.Contains(validate, "note") {
		emitNoteText(form, n, noteError != "", noteError)
	}
	if len(validate) == 0 {
		emitNoteVisibility(form, user, n)
		emitNoteButtons(form, nidx >= 0)
	}
}

func readNoteDate(r *request.Request, n *person.Note) string {
	dstr := r.FormValue("date")
	if dstr == "" {
		return "The note date is required."
	}
	date, err := time.ParseInLocation("2006-01-02", dstr, time.Local)
	if err != nil {
		return fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", dstr)
	}
	n.Date = date
	return ""
}

func emitNoteDate(form *htmlb.Element, n *person.Note, focus bool, err string) {
	row := form.E("div class='formRow personeditNoteDate'")
	row.E("label for=personeditNoteDate>Date")
	row.E("input type=date id=personeditNoteDate name=date s-validate=.personeditNoteDate value=%s", n.Date.Format("2006-01-02"),
		focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readNoteText(r *request.Request, n *person.Note) string {
	if n.Note = strings.TrimSpace(r.FormValue("note")); n.Note == "" {
		return "The note text is required."
	}
	return ""
}

func emitNoteText(form *htmlb.Element, n *person.Note, focus bool, err string) {
	row := form.E("div class='formRow personeditNoteText'")
	row.E("label for=personeditNoteText>Note Text")
	row.E("input id=personeditNoteText name=note s-validate=.personeditNoteText value=%s", n.Note, focus, "autofocus")
	if err != "" {
		row.E("div class=formError>%s", err)
	}
}

func readNoteVisibility(r *request.Request, user *person.Person, n *person.Note) {
	switch r.FormValue("visibility") {
	case "webmaster":
		if user.IsWebmaster() {
			n.Visibility = person.NoteVisibleToWebmaster
		}
	case "admin":
		if user.IsAdminLeader() {
			n.Visibility = person.NoteVisibleToAdmins
		}
	case "leader":
		n.Visibility = person.NoteVisibleToLeaders
	case "contact":
		n.Visibility = person.NoteVisibleWithContact
	case "person":
		n.Visibility = person.NoteVisibleWithPerson
	}
}

func emitNoteVisibility(form *htmlb.Element, user *person.Person, n *person.Note) {
	row := form.E("div class=formRow")
	row.E("label>Visibility")
	rbs := row.E("div class=formInput")
	if user.IsWebmaster() {
		rbs.E("div").E("s-radio name=visibility value=webmaster label='Webmasters only'",
			n.Visibility == person.NoteVisibleToWebmaster, "checked")
	}
	if user.IsAdminLeader() {
		rbs.E("div").E("s-radio name=visibility value=admin label='DPS staff only'",
			n.Visibility == person.NoteVisibleToAdmins, "checked")
	}
	rbs.E("div").E("s-radio name=visibility value=leader label='SERV leads only'",
		n.Visibility == person.NoteVisibleToLeaders, "checked")
	rbs.E("div").E("s-radio name=visibility value=contact label='Anyone who can see contact info'",
		n.Visibility == person.NoteVisibleWithContact, "checked")
	rbs.E("div").E("s-radio name=visibility value=person label=Anyone",
		n.Visibility == person.NoteVisibleWithPerson, "checked")
}

func emitNoteButtons(form *htmlb.Element, canDelete bool) {
	buttons := form.E("div class=formButtons")
	if canDelete {
		buttons.E("div class=formButtonSpace")
	}
	buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	buttons.E("input type=submit name=save class='sbtn sbtn-primary' value=Save")
	if canDelete {
		// This button comes last in the tree order so that it is not
		// the default.  But it comes first in the visual order because
		// of the formButton-beforeAll class.
		buttons.E("input type=submit name=delete class='sbtn sbtn-danger formButton-beforeAll' value=Delete")
	}
}

/*
func apply(r *request.Request, p *person.Person, nidx int, d *notedata) (ok bool) {
	var up *person.Updater

	up = p.Updater()
	if nidx >= 0 && d.delnote {
		up.Notes = append(up.Notes[:nidx], up.Notes[nidx+1:]...)
	} else {
		var n *person.Note

		if nidx >= 0 {
			n = up.Notes[nidx]
		} else {
			n = new(person.Note)
			up.Notes = append(up.Notes, n)
		}
		if date, err := time.ParseInLocation("2006-01-02", d.date, time.Local); err != nil {
			d.dateError = fmt.Sprintf("%q is not a valid YYYY-MM-DD date.", d.date)
		} else {
			n.Date = date
		}
		d.note = strings.TrimSpace(d.note)
		if d.note == "" {
			d.noteError = "Note text is required."
		} else {
			n.Note = d.note
		}
		if d.canWebmaster && d.visibility == person.NoteVisibleToWebmaster {
			n.Visibility = d.visibility
		} else if d.canAdmin && d.visibility == person.NoteVisibleToAdmins {
			n.Visibility = d.visibility
		} else if d.visibility == person.NoteVisibleToLeaders || d.visibility == person.NoteVisibleWithContact || d.visibility == person.NoteVisibleWithPerson {
			n.Visibility = d.visibility
		} else if n.Visibility == 0 && d.canWebmaster {
			n.Visibility = person.NoteVisibleToWebmaster
		} else if n.Visibility == 0 && d.canAdmin {
			n.Visibility = person.NoteVisibleToAdmins
		} else if n.Visibility == 0 {
			n.Visibility = person.NoteVisibleToLeaders
		}
		if d.dateError != "" || d.noteError != "" {
			return false
		}
	}
	r.Transaction(func() {
		p.Update(r, up, personFields)
	})
	return true
}
*/
