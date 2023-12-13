// Package problem provides for the accumulation of user-visible problem
// messages inside a standard Go "error".
package problem

import (
	"fmt"
	"strings"
)

// A Problem is a single problem message.  It mainly exists to take the place of
// a List in code contexts where only one problem is possible.  (errors.New
// doesn't serve that case well because go vet complains about capitalization of
// error messages, but problem strings are intended for end users and should be
// capitalized.)
type Problem string

// New creates a new Problem with the specified problem string.
func New(p string) error { return Problem(p) }

// NewF creates a new Problem with the specified problem string and interpolated
// arguments.
func NewF(p string, args ...interface{}) error { return Problem(fmt.Sprintf(p, args...)) }

// Error implements the "error" interface for Problem.  However, its
// implementation panics.  Problems are expected and allowed to be passed around
// as errors, but they should never be rendered through their Error function.
func (p Problem) Error() string { panic("call to Problem.Error") }

// String returns the Problem in string form.  (Kind of pointless, but just for
// completeness' sake.)
func (p Problem) String() string { return string(p) }

// OK returns true if the problem string is empty.
func (p Problem) OK() bool { return p == "" }

// Problems returns the list of problem strings.
func (p Problem) Problems() []string {
	if p == "" {
		return nil
	}
	return []string{string(p)}
}

// OrNil returns the receiver problem if it is filled in, and nil otherwise.
func (p Problem) OrNil() error {
	if p == "" {
		return nil
	}
	return p
}

// A List is a list of zero or more problem messages.  The zero value is an
// empty list.
type List struct{ l []string }

// Error implements the "error" interface for List.  However, its implementation
// panics.  Pointers to Lists are expected and allowed to be passed around as
// errors, but they should never be rendered through their Error function.
func (l *List) Error() string { panic("call to List.Error") }

// String returns the List in string form.
func (l *List) String() string { return strings.Join(l.l, "\n") }

// Add adds a problem string to the List.
func (l *List) Add(p string) { l.l = append(l.l, p) }

// AddF adds a formatted problem string to the List.
func (l *List) AddF(p string, args ...interface{}) { l.l = append(l.l, fmt.Sprintf(p, args...)) }

// AddProblem adds a Problem to the List.
func (l *List) AddProblem(p Problem) {
	if p != "" {
		l.l = append(l.l, string(p))
	}
}

// AddList adds the contents of the argument List to the receiver List.
func (l *List) AddList(o *List) {
	l.l = append(l.l, o.l...)
}

// AddError adds an error to the List.  It has special handling for cases where
// the err is nil, a Problem, or a List.
func (l *List) AddError(err error) {
	switch err := err.(type) {
	case *List:
		l.AddList(err)
	case Problem:
		l.AddProblem(err)
	case nil:
		break
	default:
		l.l = append(l.l, err.Error())
	}
}

// OK returns true if the List is empty.
func (l *List) OK() bool { return len(l.l) == 0 }

// Problems returns the list of problem strings.
func (l *List) Problems() []string { return l.l }

// OrNil returns false if the List is empty, otherise the List itself.
func (l *List) OrNil() error {
	if len(l.l) == 0 {
		return nil
	}
	return l
}

// OK returns true if the argument represents the lack of any problem (a nil
// error or an empty Problem or List).
func OK(err error) bool {
	switch err := err.(type) {
	case nil:
		return true
	case *List:
		return err.OK()
	case Problem:
		return err.OK()
	default:
		return false
	}
}

// Problems returns the list of problem strings.
func Problems(err error) []string {
	switch err := err.(type) {
	case nil:
		return nil
	case *List:
		return err.Problems()
	case Problem:
		return err.Problems()
	default:
		return []string{err.Error()}
	}
}

// String returns the problem in string form.
func String(err error) string {
	switch err := err.(type) {
	case nil:
		return ""
	case *List:
		return err.String()
	case Problem:
		return err.String()
	default:
		return err.Error()
	}
}
