// Package form provides a common infrastructure for forms on the website.  It
// helps ensure consistency in form appearance and behavior.
//
// Code that handles a specific form starts by constructing a Form structure
// giving its basic characteristics.  It then calls methods on the Form to
// handle requests to the Form page.
package form

import (
	"net/http"
	"net/url"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// A Form is a handle to a single form instance.  It has public fields set by
// the calling code to define the basic form characteristics, and methods used
// to respond to web requests.
//
// Calling code should set all of the public fields of the structure, and then
// call the Start method.  Then it should walk through each of the fields of
// the form, calling the AddRow method for each.
type Form struct {
	// Method is the form submission method.  If not set, it defaults to
	// http.MethodPost.
	Method string
	// Action is the URL to which the form should submit.  If not set, no
	// "action" attribute is emitted, so the form will submit to the same
	// URL than rendered it.
	Action string
	// Attrs contains additional attributes for the <form> element, in a
	// format suitable to be passed to htmlb.Element.Attr().  This is
	// usually used for Unpoly attributes such as up-main, up-target, etc.
	Attrs string
	// Flags is a bitmask of flags governing the format of the Form.
	Flags Flag
	// Title is the English string to be displayed in a title bar for the
	// form.  It will be localized before being displayed.  If it is empty,
	// no title bar is displayed.  Usually there is a title for a dialog
	// form and not for non-dialog forms.
	Title string
	// TitleStyle is the style for the form title bar, if any.  It defaults
	// to "primary".  Other valid values are "secondary", "warning", and
	// "danger".  Usually this is the same style as is used for the first
	// (default) button.
	TitleStyle string

	// mode is the operating mode: get, validate, or submit.
	mode string
	// validate is the set of field names being validated in a validate
	// request.
	validate sets.Set[string]
	// values is a copy of the form values.
	values url.Values
}

// Flag is a flag governing the format of a Form.
type Flag uint8

const (
	// Dialog indicates whether the form appears in a dialog box.  This
	// changes the styling, adds a Cancel button, etc.
	Dialog Flag = 1 << iota
	// NoSubmit indicates that the form is never submitted to the server;
	// its submit is trapped by Javascript and handled client-side.
	NoSubmit
	// TwoColumn forces the form to use a two-column layout even when there
	// is space for three.
	TwoColumn
	// Centered centers the form in its container.  The default is for it
	// to be left-justified in its container.
	Centered
)

// Start begins the processing of the receiver form to handle the supplied
// request.
func (f *Form) Start(r *request.Request) {
	// Examine the request to determine the operating mode:  get, validate,
	// or submit.
	if r.Method != http.MethodPost {
		f.mode = "get"
	} else if v := r.FormValue("X-Up-Validate"); v != "" {
		f.mode = "validate"
		f.validate = sets.New(strings.Split(v, " ")...)
		f.values = r.Form
	} else {
		f.mode = "submit"
		f.values = r.Form
	}
}

// AddRow adds a row to the form.  The row structure defines the row layout.
// The error string should be empty if the field is valid and non-empty if it
// is invalid.  The function will be called later to render the row, if it
// turns out the form is being rendered at all and if this row is needed.
func (f *Form) AddRow(row *Row, error string, fn RenderFunc) {}

// A RenderFunc renders the HTML for a row, appending it to the specified
// parent element (which it must not close).  If the autofocus flag is true,
// the function should put an autofocus attribute on an appropriate input field
// within the row.  Generally this will be the field that has an error, if any,
// or else the first field in the row.
type RenderFunc func(parent *htmlb.Element, autofocus bool)

// A Row defines the layout of a single form row.  Rows generally contain three
// logical columns:  label, input, and error/help.  (These may be displayed as
// two columns or even a single column on narrow devices.)
//
// There are three possible layouts.  The most common layout is selected when
// Wide is false.  It puts Label in the first logical column, the result of
// RenderFunc in the second, and the error and Help text in the third.
//
// If Wide is true and Label is set, the layout will have Label in the first
// logical column and the result of RenderFunc spanning the second and third
// logical columns.  The error and Help text, if any, will be below the result
// of RenderFunc, also spanning the second and thild logical columns.
//
// If Wide is true and Label is not set, the layout will have the result of
// RenderFunc spanning all three columns.  The error and Help text, if any,
// will be below the result of RenderFunc, also spanning all three columns.
type Row struct {
	// Label is the label for the form row.
	Label string
	// Help is the help text for the form row.
	Help string
	// Wide indicates whether the result of RenderFunc can span multiple
	// logical columns.
	Wide bool
	// IfValidate specifies which field(s) we need to be validating in
	// order for this row to be shown.  This applies only when the form is
	// in validate mode; in the other modes, all form rows are shown.
	IfValidate []string
}
