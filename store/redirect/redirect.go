// Package redirect defines the Redirect type, which describes a URL entrypoint
// that should be redirected to a different URL.
package redirect

// ID uniquely identifies a redirect.
type ID int

// Redirect describes a URl entrypoint that should be redirected to a different
// URL.
type Redirect struct {
	// ID is the unique identifier of the List.
	ID ID
	// Entry is the URL entrypoint that will be redirected.
	Entry string
	// Target is the URL to which the entrypoint will be redirected.
	Target string
}

func (l *Redirect) Clone() (c *Redirect) {
	c = new(Redirect)
	*c = *l
	return c
}
