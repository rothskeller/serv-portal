package form

import (
	"cmp"
	"slices"

	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

// A RowGroup is a group of related rows, with a label.
type RowGroup struct {
	// Internal note:  A RowGroup without a label is used as the top-level
	// row container for the form.

	// Label is the label for the row group.
	Label string
	// Rows is the set of rows in the group.
	Rows []Row
	// Order is the read order for the group.  It is usually zero.
	Order int

	// vl preserves the validation list from the call to ShouldEmit until
	// the call to Emit.
	vl request.ValidationList
	// firstfail preserves the identity of the first invalid row from the
	// call to Read until the call to Emit.
	firstfail Row
}

func (rg *RowGroup) Get() {
	for _, row := range rg.Rows {
		row.Get()
	}
}

func (rg *RowGroup) ReadOrder() int {
	return rg.Order
}

func (rg *RowGroup) Read(r *request.Request) (ok bool) {
	rows := slices.Clone(rg.Rows)
	slices.SortStableFunc(rows, func(a, b Row) int {
		return cmp.Compare(a.ReadOrder(), b.ReadOrder())
	})
	rg.firstfail, ok = nil, true
	for _, row := range rows {
		if !row.Read(r) {
			if rg.firstfail == nil {
				rg.firstfail = row
			}
			ok = false
		}
	}
	return ok
}

func (rg *RowGroup) ShouldEmit(vl request.ValidationList) bool {
	rg.vl = vl
	for _, row := range rg.Rows {
		if row.ShouldEmit(vl) {
			return true
		}
	}
	return false
}

func (rg *RowGroup) Emit(r *request.Request, parent *htmlb.Element, focus bool) {
	var group *htmlb.Element

	if rg.Label != "" {
		group = parent.E("div class=formGroup")
		group.E("div class='formRow-3col formGroupLabel'").T(r.Loc(rg.Label))
	} else {
		group = parent
	}
	for _, row := range rg.Rows {
		if row.ShouldEmit(rg.vl) {
			rowfocus := focus && (row == rg.firstfail || rg.firstfail == nil)
			row.Emit(r, group, rowfocus)
			if rowfocus {
				focus = false
			}
		}
	}
}
