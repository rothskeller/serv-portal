package form

import (
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// A Button is a single button to appear at the bottom of the form.
type Button struct {
	// Name is the input name with which the button is submitted.  It can be
	// empty for the first button but must be set for all others, and must
	// be unique among all controls on the form.
	Name string
	// Label is the English label for the button.  It will be localized for
	// display.
	Label string
	// Style is the style ("primary", "secondary", "warning", "danger") for
	// the button.  The default is "primary".
	Style string
	// OnClick is the function to be called when the button is pressed (and
	// the form contents are valid).  The OnClick function for the first
	// button is called if the form is submitted without clicking any button
	// (e.g., with the Enter key).  The function should return true if the
	// page has been handled, or false if the form should be rendered.
	OnClick func() bool
}

// emitButtons writes out the set of buttons for the form.  If dialog is false,
// they are written left to right, aligned to the left side of the form.  If
// dialog is true, the first one is written all the way to the right, a Cancel
// button that dismisses the dialog is written to its left, and the rest of the
// buttons are written to the left of those with a gap, written left to right.
func emitButtons(r *request.Request, form *htmlb.Element, buttons []*Button, dialog bool) {
	brow := form.E("div class=formButtons2")
	bgrp := brow.E("div class=formButtonGroup")
	if dialog {
		bgrp.E("button type=button class='sbtn sbtn-secondary' up-dismiss").R(r.Loc("Cancel"))
		emitButton(r, bgrp, buttons[0])
		if len(buttons) > 1 {
			bgrp = brow.E("div class=formButtonGroup")
			for _, btn := range buttons[1:] {
				emitButton(r, bgrp, btn)
			}
		}
	} else {
		for _, btn := range buttons {
			emitButton(r, bgrp, btn)
		}
	}
}
func emitButton(r *request.Request, bgrp *htmlb.Element, btn *Button) {
	var style = btn.Style
	if style == "" {
		style = "primary"
	}
	bgrp.E("input type=submit class='sbtn sbtn-%s' value=%s", style, r.Loc(btn.Label),
		btn.Name != "", "name=%s", btn.Name)
}

// executeClickedButton determines which button was pressed and executes its
// OnClick handler.  If no button was pressed, the OnClick handler of the first
// button is executed.
func executeClickedButton(r *request.Request, buttons []*Button) bool {
	for _, btn := range buttons {
		if r.FormValue(btn.Name) != "" {
			return btn.OnClick()
		}
	}
	return buttons[0].OnClick()
}
