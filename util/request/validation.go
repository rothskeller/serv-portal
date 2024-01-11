package request

type ValidationList map[string]struct{}

// ValidationList returns the list of fields being validated in the request.
func (r *Request) ValidationList() (vl ValidationList) {
	list := r.LogEntry.Validate
	if len(list) == 0 {
		return nil
	}
	vl = make(ValidationList)
	for _, v := range list {
		vl[v] = struct{}{}
	}
	return vl
}

// Enabled returns whether validation is enabled.
func (vl ValidationList) Enabled() bool { return vl != nil }

// Validating returns whether the named item is being validated.  It also
// returns true if no validation is occurring.
func (vl ValidationList) Validating(name string) bool {
	if vl == nil {
		return true
	}
	_, ok := vl[name]
	return ok
}

// ValidatingAny returns whether any of the named items are being validated.  It
// also returns true if no validation is occurring.
func (vl ValidationList) ValidatingAny(names ...string) bool {
	if vl == nil {
		return true
	}
	for _, name := range names {
		if _, ok := vl[name]; ok {
			return true
		}
	}
	return false
}
