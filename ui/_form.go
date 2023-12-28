package ui

import (
	"fmt"
	"net/http"
	"regexp"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// FormOpts contains the configuration for a Form.
type FormOpts struct {
	// Dialog is a flag indicating that this form is in a dialog box.
	Dialog bool
	// TwoColumn is a flag indicating that the form should be limited to
	// two columns.
	TwoColumn bool
	// Target is the target for successful submission of the form.
	Target string
	// Title is the title of the form, or empty if it shouldn't have one.
	// The title gets localized before display.
	Title string
	// TitleColor is the color of the form title.  It defaults to "primary".
	TitleColor string
	// Fields is the list of fields of the form.  It should be in the order
	// that the fields need to be read for proper validation.  Usually that
	// is also the order of display, but if not, CSS order attributes can be
	// used to correct the display order.
	Fields []*FormField
	// SubmitLabel is the label of the submit button for the form, in
	// English.
	SubmitLabel string
	// SubmitColor is the color of the submit button.  It defaults to
	// "primary".
	SubmitColor string
	// ExtraButtons is the list of extra buttons on the form, beyond Save
	// and Cancel.
	ExtraButtons []*FormButton
}

// FormButton describes an extra button on a form.
type FormButton struct {
	// Name is the name submitted with the form when this button is pressed,
	// in English.
	Name string
	// Color is the color of the button.  It defaults to "primary".
	Color string
	// Label is the label of the button, in English.
	Label string
}

// FormField is the definition of a single form field.
type FormField struct {
	// InputColumns indicates how many columns the input spans.  It defaults
	// to 1, but can be 2 or 3.
	InputColumns int
	// ID is the ID of the field input.  If the field has more than one
	// input, ID is the name of the first one.
	ID string
	// Label is the label for the field, in English.
	Label string
	// ShouldEmit is a function that returns whether the field needs to be
	// emitted.  It is passed the set of targets being rendered.  If this
	// function is not defined, the field is always emitted.
	ShouldEmit func(targets sets.Set) bool
	// Emit is a function that emits the input control.
	Emit func(*request.Request, *htmlb.Element)
	// Error is the current error with the field, in English, or an empty
	// string if there is none.
	Error string
	// Help is the help text to be displayed for the field, in English, or
	// an empty string if there is none.
	Help string
}

// Form handles a request for a form.  It returns true if the form has been
// successfully submitted (in which case, nothing has been rendered), or false
// otherwise (in which case the form has been rendered).
func Form(r *request.Request, opts *FormOpts) bool {
	var hasError bool

	// If this is a POST, read the previous form contents.
	if r.Method == http.MethodPost {
		if hasError = readForm(r, opts); !hasError {
			return true
		}
	}
	// This is either a GET or we had an error, so display the form.
	r.HTMLNoCache()
	if hasError {
		r.WriteHeader(http.StatusUnprocessableEntity)
	}
	html := htmlb.HTML(r)
	defer html.Close()
	emitForm(r, html, opts)
	return false
}

var commaSplitRE = regexp.MustCompile(`\s*,\s*`)

func emitForm(r *request.Request, html *htmlb.Element, opts *FormOpts) {
	var targets sets.Set[string]

	form := html.E("form class=form method=POST up-main")
	if opts.TwoColumn {
		form.Attr("class=form-2col")
	}
	if opts.Dialog {
		form.Attr("up-layer=parent")
	}
	if opts.Target != "" {
		form.Attr("up-target=%s", opts.Target)
	}
	form.E("input type=hidden name=csrf value=%s", r.CSRF)
	if opts.Title != "" {
		if opts.TitleColor == "" {
			opts.TitleColor = "primary"
		}
		form.E("div class='formTitle formTitle-%s'>%s", opts.TitleColor, r.Loc(opts.Title))
	}
	// If we are validating, we can be more efficient by emitting only those
	// fields that Unpoly is trying to validate.
	if r.Request.Header.Get("X-Up-Validate") != "" {
		// X-Up-Validate contains the list of fields being validated,
		// but it's better for us to key off of X-Up-Target, which is
		// the list of elements Unpoly is looking for in the response.
		// The two may not be the same when one field's validation can
		// result in new errors for a different field.
		targets = sets.New(commaSplitRE.Split(r.Request.Header.Get("X-Up-Target"), -1)...)
	}
	for _, ff := range opts.Fields {
		if targets == nil || ff.ShouldEmit == nil || ff.ShouldEmit(targets) {
			emitField(r, form, ff)
		}
	}
	if targets == nil {
		emitButtons(r, form, opts)
	}
}

func emitField(r *request.Request, form *htmlb.Element, field *FormField) {
	var row *htmlb.Element

	if field.InputColumns == 3 {
		row = form.E("div class=formRow-3col")
	} else {
		row = form.E("div class=formRow")
	}
	if field.Label != "" {
		label := row.E("label")
		if field.ID != "" {
			label.Attr("for=%s", field.ID)
		}
		label.T(r.Loc(field.Label))
	}
	field.Emit(r, row)
	if field.Error != "" {
		row.E("div class=formError").T(r.Loc(field.Error))
	}
	if field.Help != "" {
		row.E("div class=formHelp").T(r.Loc(field.Help))
	}
}

func emitButtons(r *request.Request, form *htmlb.Element, opts *FormOpts) {
	buttons := form.E("div class=formButtons")
	if len(opts.ExtraButtons) != 0 && opts.Dialog {
		buttons.E("div class=formButtonSpace")
	}
	if opts.Dialog {
		buttons.E("button type=button class='sbtn sbtn-secondary' up-dismiss>Cancel")
	}
	if opts.SubmitLabel == "" {
		opts.SubmitLabel = "Save"
	}
	if opts.SubmitColor == "" {
		opts.SubmitColor = "primary"
	}
	buttons.E("input type=submit class='sbtn sbtn-%s' value=%s", opts.SubmitColor, r.Loc(opts.SubmitLabel))
	if len(opts.ExtraButtons) != 0 && !opts.Dialog {
		buttons.E("div class=formButtonSpace")
	}
	for _, eb := range opts.ExtraButtons {
		if eb.Color == "" {
			eb.Color = "primary"
		}
		button := buttons.E("input type=submit name=%s class='sbtn sbtn-%s' value=%s", eb.Name, eb.Color, r.Loc(eb.Label))
		if opts.Dialog {
			button.Attr("class=formButton-beforeAll")
		}
	}
}

// EmitString emits and returns a string input field.
func EmitString(row *htmlb.Element, id, name, value string) *htmlb.Element {
	input := row.E("input class=formInput name=%s value=%s", name, value)
	if id != "" {
		input.A("id=%s", id)
	}
	return input
}

// EmitSelectString emits a select input field with a set of string values.  The
// passed options should be in English and will be localized.
func EmitSelectString(r *request.Request, row *htmlb.Element, id, name, value, unset string, options []string) {
	sel := row.E("select class=formInput name=%s", name)
	if id != "" {
		sel.A("id=%s", id)
	}
	if unset != "" {
		var valueidx = -1

		for i, opt := range options {
			if value == opt {
				valueidx = i
				break
			}
		}
		if valueidx == -1 {
			sel.E("option value='' selected>%s", r.Loc(unset))
		}
	}
	for _, opt := range options {
		sel.E("option value=%s", opt, value == opt, "selected").T(r.Loc(opt))
	}
}

// EmitSelectEnum emits a select input field with a set of enumerated values.
func EmitSelect(r *request.Request, row *htmlb.Element, id, name string, value any, unset string, options []any) {
	var (
		valueidx = -1
		opts     = make([]string, len(options))
	)
	for i, opt := range options {
		if value == opt {
			valueidx = i
		}
		opts[i] = opt.(fmt.Stringer).String()
	}
	sel := row.E("select class=formInput name=%s", name)
	if id != "" {
		sel.A("id=%s", id)
	}
	if unset != "" && valueidx == -1 {
		sel.E("option value='' selected>%s", r.Loc(unset))
	}
	for _, opt := range opts {
		sel.E("option value=%s", opt, value == opt, "selected").T(r.Loc(opt))
	}
}
