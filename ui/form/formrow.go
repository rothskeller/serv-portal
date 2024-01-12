package form

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// A Row is a single row of a form.
type Row interface {

	// Get performs any needed initialization of the row for GET requests.
	// It is not called for POST requests.  It is usually a no-op.
	Get()

	// ReadOrder determines the order in which rows of the form should be
	// read.  It is normally zero, in which case rows are read in visual
	// order.  However, it can be overridden in cases where rows need to be
	// read in a different order (e.g., one row's validation is influenced
	// by a later row's value).
	ReadOrder() int

	// Read reads the values for the row from the submitted form.  It is
	// called only for POST requests.  It returns true if the values
	// provided for the row are acceptable.
	Read(r *request.Request) bool

	// ShouldEmit returns whether the row should be emitted.
	ShouldEmit(vl request.ValidationList) bool

	// Emit emits the row, as a child of the supplied parent.  If focus is
	// true, the (first) input in the row gets the initial focus of the
	// form.
	Emit(r *request.Request, parent *htmlb.Element, focus bool)
}

// A BaseRow is a base class for a form row.  It provides default
// implementations for some (but not all) of the Row methods.
type BaseRow struct{}

func (br BaseRow) Get()                                     {}
func (br BaseRow) ReadOrder() int                           { return 0 }
func (br BaseRow) ShouldEmit(_ request.ValidationList) bool { return true }

// A LabeledRow is a base class for a typical form row.  It does not implement
// the Row interface by itself, but it can be embedded in other types to provide
// shared code for rendering the row element, the label, the error message, and
// the help message.
type LabeledRow struct {
	BaseRow

	// RowID is the ID of the row element.  This is typically used as the
	// target element for field validation, so it's needed on all rowws
	// except those where validation is disabled.
	RowID string

	// Label is the English label of the form row.  It will be localized
	// before display.  It is optional.
	Label string

	// Error is the localized error message to be displayed in the form row.
	// It is optional.
	Error string

	// Help is the English help text for the form row.  It will be localized
	// before display.  It is optional.
	Help string
}

// EmitPrefix creates the labeled row, adds the label, and returns the row
// element.
func (lr *LabeledRow) EmitPrefix(r *request.Request, parent *htmlb.Element, focusID string) *htmlb.Element {
	row := parent.E("div class=formRow", lr.RowID != "", "id=%s", lr.RowID)
	if lr.Label != "" {
		row.E("label", focusID != "", "for=%s", focusID).T(r.Loc(lr.Label))
	}
	return row
}

// EmitSuffix adds the error and help text to the row element.
func (lr *LabeledRow) EmitSuffix(r *request.Request, row *htmlb.Element) {
	if lr.Error != "" {
		row.E("div class=formError").T(lr.Error)
		r.LogEntry.Problems.Add(lr.Error)
	}
	if lr.Help != "" {
		row.E("div class=formHelp").T(r.Loc(lr.Help))
	}
}

type CheckboxesRow struct {
	LabeledRow
	FocusID  string
	Name     string
	Boxes    []*Checkbox
	Validate string
}
type Checkbox struct {
	Name     string
	Value    string
	Label    string
	CheckedP *bool
}

var _ Row = (*CheckboxesRow)(nil) // interface check

func (cbr *CheckboxesRow) ShouldEmit(vl request.ValidationList) bool {
	if cbr.Name != "" && vl.Validating(cbr.Name) {
		return true
	}
	for _, cb := range cbr.Boxes {
		if cb.Name != "" && vl.Validating(cb.Name) {
			return true
		}
	}
	return false
}
func (cbr *CheckboxesRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if cbr.FocusID == "" && cbr.RowID != "" {
		cbr.FocusID = cbr.RowID + "-in"
	}
	row := cbr.EmitPrefix(r, parent, cbr.FocusID)
	box := row.E("div class=formInput")
	for i, cb := range cbr.Boxes {
		box.E("input type=checkbox class=s-check label=%s", r.Loc(cb.Label),
			cb.Name != "", "name=%s", cb.Name,
			cb.Name == "", "name=%s", cbr.Name,
			cb.Value != "", "value=%s", cb.Value,
			*cb.CheckedP, "checked",
			i == 0 && cbr.FocusID != "", "id=%s", cbr.FocusID,
			i == 0 && focus, "autofocus",
			cbr.Validate == "", "s-validate",
			cbr.Validate != "" && cbr.Validate != NoValidate, "s-validate=%s", cbr.Validate,
		)
	}
	cbr.EmitSuffix(r, row)
}
func (cbr *CheckboxesRow) Read(r *request.Request) bool {
	cbr.Error = ""
	for _, cb := range cbr.Boxes {
		if cb.Name != "" {
			*cb.CheckedP = r.FormValue(cb.Name) != ""
		} else if cbr.Name != "" && cb.Value != "" {
			*cb.CheckedP = slices.Contains(r.Form[cbr.Name], cb.Value)
		} else {
			panic("CheckboxesRow must have box name set, or row name and box value set")
		}
	}
	return true
}

type FlagsRow[T ~uint] struct {
	CheckboxesRow
	ValueP    *T
	Flags     []T
	LabelFunc func(r *request.Request, v T) string
	checkeds  []bool
}

func (fr *FlagsRow[T]) init(r *request.Request) {
	fr.Boxes = make([]*Checkbox, len(fr.Flags))
	fr.checkeds = make([]bool, len(fr.Flags))
	for i, f := range fr.Flags {
		fr.Boxes[i] = &Checkbox{
			Value:    strconv.Itoa(int(f)),
			Label:    optLabel(r, f, fr.LabelFunc),
			CheckedP: &fr.checkeds[i],
		}
	}
}
func (fr *FlagsRow[T]) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if fr.Boxes == nil {
		fr.init(r)
		for i, f := range fr.Flags {
			fr.checkeds[i] = *fr.ValueP&f != 0
		}
	}
	fr.CheckboxesRow.Emit(r, parent, focus)
}
func (fr *FlagsRow[T]) Read(r *request.Request) bool {
	fr.init(r)
	fr.CheckboxesRow.Read(r)
	for i, f := range fr.Flags {
		if fr.checkeds[i] {
			*fr.ValueP |= f
		} else {
			*fr.ValueP &^= f
		}
	}
	return true
}

type InputRow struct {
	LabeledRow
	FocusID  string
	Name     string
	ValueP   *string
	Validate string
}
type TextInputRow = InputRow

var _ Row = (*InputRow)(nil) // interface check

const NoValidate = "NO_VALIDATE"

func (ir *InputRow) ShouldEmit(vl request.ValidationList) bool {
	return vl.Validating(ir.Name)
}
func (ir *InputRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	ir.EmitSuffix(r, ir.EmitPrefix(r, parent, focus))
}
func (ir *InputRow) EmitPrefix(r *request.Request, parent *htmlb.Element, focus bool) *htmlb.Element {
	if ir.FocusID == "" && ir.RowID != "" {
		ir.FocusID = ir.RowID + "-in"
	}
	row := ir.LabeledRow.EmitPrefix(r, parent, ir.FocusID)
	return row.E("input class=formInput name=%s value=%s", ir.Name, *ir.ValueP,
		ir.FocusID != "", "id=%s", ir.FocusID,
		focus, "autofocus",
		ir.Validate == "", "s-validate",
		ir.Validate != "" && ir.Validate != NoValidate, "s-validate=%s", ir.Validate,
	)
}
func (ir *InputRow) EmitSuffix(r *request.Request, input *htmlb.Element) {
	ir.LabeledRow.EmitSuffix(r, input.Parent())
}
func (ir *InputRow) Read(r *request.Request) bool {
	*ir.ValueP = strings.TrimSpace(r.FormValue(ir.Name))
	ir.Error = ""
	return true
}

type DateRow struct {
	InputRow
}

func (dr *DateRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	// Unpoly validation doesn't work properly for date inputs, so we won't
	// validate them.  They'll get checked when the form is submitted.
	dr.Validate = NoValidate
	dr.EmitSuffix(r, dr.EmitPrefix(r, parent, focus).A("type=date"))
}
func (dr *DateRow) Read(r *request.Request) bool {
	dr.InputRow.Read(r)
	if *dr.ValueP == "" {
		return true
	}
	if t, err := time.Parse("2006-01-02", *dr.ValueP); err != nil || t.Format("2006-01-02") != *dr.ValueP {
		dr.Error = fmt.Sprintf(r.Loc("%q is not a valid YYYY-MM-DD date."), *dr.ValueP)
		return false
	}
	return true
}

type IntegerRow[T ~int | ~uint] struct {
	InputRow
	ValueP   *T
	Min      int
	HideZero bool
	ValueStr string
}

func (ir *IntegerRow[T]) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if ir.InputRow.ValueP == nil {
		ir.InputRow.ValueP = &ir.ValueStr
		var value = int(*ir.ValueP)
		if value == 0 && ir.HideZero {
			ir.ValueStr = ""
		} else {
			ir.ValueStr = strconv.Itoa(value)
		}
	}
	ir.EmitSuffix(r, ir.EmitPrefix(r, parent, focus).A("type=number min=%d", ir.Min))
}
func (ir *IntegerRow[T]) Read(r *request.Request) bool {
	var (
		value int
		err   error
	)
	ir.InputRow.ValueP = &ir.ValueStr
	ir.InputRow.Read(r)
	if ir.ValueStr == "" {
		value = 0
	} else if value, err = strconv.Atoi(ir.ValueStr); err != nil {
		ir.Error = fmt.Sprintf(r.Loc("%q is not a valid number."), ir.ValueStr)
		return false
	}
	*ir.ValueP = T(value)
	return true
}

type PasswordRow struct {
	InputRow
	Autocomplete string
}

func (pr *PasswordRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if pr.Autocomplete == "" {
		pr.Autocomplete = "password"
	}
	pr.EmitSuffix(r,
		pr.EmitPrefix(r, parent, focus).
			A("type=password autocomplete=%s autocapitalize=none", pr.Autocomplete))
}
func (pr *PasswordRow) Read(r *request.Request) bool {
	// Not deferring to InputRow.Read because we don't want
	// strings.TrimSpace.
	*pr.ValueP = r.FormValue(pr.Name)
	pr.Error = ""
	return true
}

type SearchComboRow struct {
	InputRow
	ValueKey    string
	Filter      string
	Placeholder string
}

func (scr *SearchComboRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	scr.EmitSuffix(r,
		scr.EmitPrefix(r, parent, focus).
			A("class=s-search s-filter=%s", scr.Filter,
				scr.Placeholder != "", "placeholder=%s", scr.Placeholder,
				scr.ValueKey != "", "s-value=%s", scr.ValueKey))
}
func (scr *SearchComboRow) Read(r *request.Request) bool {
	scr.InputRow.Read(r)
	if len(r.Form[scr.Name]) > 1 {
		scr.ValueKey = r.Form[scr.Name][1]
	} else {
		scr.ValueKey = ""
	}
	return true
}

type MessageRow struct {
	LabeledRow
	HTML string
}

var _ Row = (*MessageRow)(nil) // interface check

func (mr *MessageRow) ShouldEmit(vl request.ValidationList) bool {
	return !vl.Enabled()
}
func (mr *MessageRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	row := mr.EmitPrefix(r, parent, "")
	row.E("div class=formInput").R(mr.HTML)
	mr.EmitSuffix(r, row)
}
func (mr *MessageRow) Read(r *request.Request) bool { return true }

type RadioGroupRow[T comparable] struct {
	LabeledRow
	FocusID   string
	Name      string
	ValueP    *T
	Options   []T
	ValueFunc func(v T) string
	LabelFunc func(r *request.Request, v T) string
	Validate  string
}

var _ Row = (*RadioGroupRow[string])(nil) // interface check

func (rgr *RadioGroupRow[T]) ShouldEmit(vl request.ValidationList) bool {
	return vl.Validating(rgr.Name)
}
func (rgr *RadioGroupRow[T]) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if rgr.FocusID == "" && rgr.RowID != "" {
		rgr.FocusID = rgr.RowID + "-in"
	}
	row := rgr.EmitPrefix(r, parent, rgr.FocusID)
	box := row.E("div class=formInput",
		rgr.Validate == "", "s-validate",
		rgr.Validate != "" && rgr.Validate != NoValidate, "s-validate=%s", rgr.Validate,
	)
	for i, opt := range rgr.Options {
		box.E("s-radio name=%s value=%s label=%s", rgr.Name, optValue(opt, rgr.ValueFunc), optLabel(r, opt, rgr.LabelFunc),
			i == 0 && rgr.FocusID != "", "id=%s", rgr.FocusID,
			i == 0 && focus, "autofocus",
			*rgr.ValueP == opt, "checked")
	}
	rgr.EmitSuffix(r, row)
}
func (rgr *RadioGroupRow[T]) Read(r *request.Request) bool {
	var valuestr = r.FormValue(rgr.Name)
	rgr.Error = ""
	if valuestr == "" {
		*rgr.ValueP = *new(T)
		return true
	}
	for _, opt := range rgr.Options {
		if optValue(opt, rgr.ValueFunc) == valuestr {
			*rgr.ValueP = opt
			return true
		}
	}
	rgr.Error = fmt.Sprintf(r.Loc("%q is not a valid value for %s."), valuestr, r.Loc(rgr.Label))
	*rgr.ValueP = *new(T)
	return false
}

type SelectRow[T comparable] struct {
	LabeledRow
	FocusID     string
	Name        string
	ValueP      *T
	Options     []T
	ValueFunc   func(T) string
	LabelFunc   func(r *request.Request, v T) string
	Placeholder string
	Validate    string
}

var _ Row = (*SelectRow[string])(nil) // interface check

func (sr *SelectRow[T]) ShouldEmit(vl request.ValidationList) bool {
	return vl.Validating(sr.Name)
}
func (sr *SelectRow[T]) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if sr.FocusID == "" && sr.RowID != "" {
		sr.FocusID = sr.RowID + "-in"
	}
	row := sr.EmitPrefix(r, parent, sr.FocusID)
	sel := row.E("select class=formInput name=%s", sr.Name,
		sr.FocusID != "", "id=%s", sr.FocusID,
		focus, "autofocus",
		sr.Validate == "", "s-validate",
		sr.Validate != "" && sr.Validate != NoValidate, "s-validate=%s", sr.Validate,
	)
	if *sr.ValueP == *new(T) && sr.Placeholder != "" {
		sel.E("option value='' selected").T(r.Loc(sr.Placeholder))
	}
	for _, opt := range sr.Options {
		sel.E("option value=%s", optValue(opt, sr.ValueFunc), *sr.ValueP == opt, "selected").
			T(optLabel(r, opt, sr.LabelFunc))
	}
	sr.EmitSuffix(r, row)
}
func (sr *SelectRow[T]) Read(r *request.Request) bool {
	var valuestr = r.FormValue(sr.Name)
	sr.Error = ""
	if valuestr == "" {
		*sr.ValueP = *new(T)
		return true
	}
	for _, opt := range sr.Options {
		if valuestr == optValue(opt, sr.ValueFunc) {
			*sr.ValueP = opt
			return true
		}
	}
	sr.Error = fmt.Sprintf(r.Loc("%q is not a valid value for %s."), valuestr, r.Loc(sr.Label))
	return false
}

type TextAreaRow struct {
	LabeledRow
	FocusID  string
	Name     string
	ValueP   *string
	Wrap     string
	Validate string
}

var _ Row = (*TextAreaRow)(nil) // interface check

func (tar *TextAreaRow) ShouldEmit(vl request.ValidationList) bool {
	return vl.Validating(tar.Name)
}
func (tar *TextAreaRow) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	if tar.FocusID == "" && tar.RowID != "" {
		tar.FocusID = tar.RowID + "-in"
	}
	row := tar.EmitPrefix(r, parent, tar.FocusID)
	row.E("textarea class=formInput name=%s", tar.Name,
		tar.FocusID != "", "id=%s", tar.FocusID,
		focus, "autofocus",
		tar.Validate == "", "s-validate",
		tar.Validate != "" && tar.Validate != NoValidate, "s-validate=%s", tar.Validate,
		tar.Wrap != "", "wrap=%s", tar.Wrap,
	).T(*tar.ValueP)
	tar.EmitSuffix(r, row)
}
func (tar *TextAreaRow) Read(r *request.Request) bool {
	*tar.ValueP = r.FormValue(tar.Name)
	tar.Error = ""
	return true
}

/*
Additional types of rows that we might add:
	RoleSelectRow -> implement elsewhere to avoid import cycle
	TimeRangeRow
	FileUploadRow
	AddressRow
	DateButtonRow
*/

type Inter interface {
	Int() int
}

func optValue[T comparable](v T, fn func(T) string) string {
	if fn != nil {
		return fn(v)
	}
	switch v := any(v).(type) {
	case string:
		return v
	case Inter:
		return strconv.Itoa(v.Int())
	case fmt.Stringer:
		return v.String()
	default:
		panic(fmt.Sprintf("ValueFunc must be provided for type %T", v))
	}
}
func optLabel[T comparable](r *request.Request, v T, fn func(*request.Request, T) string) string {
	if fn != nil {
		return fn(r, v)
	}
	switch v := any(v).(type) {
	case string:
		return r.Loc(v)
	case fmt.Stringer:
		return r.Loc(v.String())
	default:
		panic(fmt.Sprintf("LabelFunc must be provided for type %T", v))
	}
}
