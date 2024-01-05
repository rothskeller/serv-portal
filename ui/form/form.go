package form

import (
	"net/http"

	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// A Form represents a fillable form on the website.
type Form struct {
	// Attrs are attributes for the form.  These typically include method,
	// action, and up-target.
	Attrs string
	// Dialog indicates whether the form appears in a dialog box.  This
	// changes the styling, adds a Cancel button, etc.
	Dialog bool
	// TwoCol forces the form to use a two-column layout even when there is
	// space for three.
	TwoCol bool
	// Centered centers the form in its container.  The default is for it to
	// be left-justified in its container.
	Centered bool
	// Title is the string to be displayed in a title bar for the form.
	// If it is empty, no title bar is displayed.  Usually there is a title
	// for a dialog form and not for non-dialog forms.
	Title string
	// TitleStyle is the style for the form title bar, if any.  It defaults
	// to "primary".  Other valid values are "secondary", "warning", and
	// "danger".  Usually this is the same style as is used for the first
	// button in the Buttons list.
	TitleStyle string
	// Rows is the set of rows to be displayed in the form.  Some of them
	// may be RowGroups with nested sets of rows.  (Multiple layers of
	// nesting are allowed, but have no effect.)
	Rows []Row
	// Buttons is the set of buttons to be displayed on the form; there must
	// be at least one.
	Buttons []*Button
}

// Handle handles a request for the form.
func (f *Form) Handle(r *request.Request) {
	var rg = RowGroup{Rows: f.Rows}
	var vl = r.ValidationList()

	if r.Method == http.MethodPost {
		var valid = rg.Read(r)
		if valid && !vl.Enabled() {
			if executeClickedButton(r, f.Buttons) {
				return
			}
		}
		r.HTMLNoCache()
		r.WriteHeader(http.StatusUnprocessableEntity)
	} else {
		rg.Get()
		r.HTMLNoCache()
	}
	html := htmlb.HTML(r)
	defer html.Close()
	form := html.E("form class=form up-main")
	if f.Attrs != "" {
		form.A(f.Attrs)
	}
	if f.TwoCol {
		form.A("class=form-2col")
	}
	if f.Centered {
		form.A("class=form-centered")
	}
	if f.Dialog {
		form.A("up-layer=parent")
	}
	if r.CSRF != "" {
		form.E("input type=hidden name=csrf value=%s", r.CSRF)
	}
	if f.Title != "" {
		style := f.TitleStyle
		if style == "" {
			style = "primary"
		}
		form.E("div class='formTitle formTitle-%s'", style).T(r.Loc(f.Title))
	}
	rg.ShouldEmit(vl)
	rg.Emit(r, form, true)
	emitButtons(r, form, f.Buttons, f.Dialog)
}
