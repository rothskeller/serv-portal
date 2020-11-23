package email

import (
	"strconv"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetUnsubscribe handles GET /unsubscribe/$token requests.  Note that these
// requests are unauthenticated; the token is the authentication.
func GetUnsubscribe(r *util.Request, token string) error {
	var (
		person *model.Person
		out    jwriter.Writer
	)
	if person = r.Tx.FetchPersonByUnsubscribe(token); person == nil {
		return util.NotFound
	}
	r.Auth.SetMe(person)
	out.RawString(`{"noEmail":`)
	out.Bool(person.NoEmail)
	out.RawString(`,"groups":[`)
	first := true
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivMember)) {
		if g.Email == "" {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(g.ID))
		out.RawString(`,"email":`)
		out.String(g.Email)
		out.RawString(`,"unsub":`)
		unsub := false
		for _, p := range g.NoEmail {
			if p == person.ID {
				unsub = true
				break
			}
		}
		out.Bool(unsub)
		out.RawByte('}')
	}
	out.RawString(`]}`)
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}

// PostUnsubscribe handles POST /unsubscribe/$token requests.
func PostUnsubscribe(r *util.Request, token string) error {
	var (
		person   *model.Person
		noEmail  bool
		needSave bool
	)
	if person = r.Tx.FetchPersonByUnsubscribe(token); person == nil {
		return util.NotFound
	}
	r.Auth.SetMe(person)
	noEmail, _ = strconv.ParseBool(r.FormValue("noEmail"))
	if noEmail != person.NoEmail {
		r.Tx.WillUpdatePerson(person)
		person.NoEmail = noEmail
		r.Tx.UpdatePerson(person)
	}
	for _, g := range r.Auth.FetchGroups(r.Auth.GroupsA(model.PrivMember)) {
		var unsub bool
		var changed bool
		if g.Email == "" {
			continue
		}
		unsub, _ = strconv.ParseBool(r.FormValue("unsub:" + strconv.Itoa(int(g.ID))))
		j := 0
		for _, p := range g.NoEmail {
			if p != person.ID || unsub {
				g.NoEmail[j] = p
				j++
			} else {
				r.Auth.WillUpdateGroup(g)
				changed = true
			}
		}
		if changed {
			g.NoEmail = g.NoEmail[:j]
		} else if unsub {
			r.Auth.WillUpdateGroup(g)
			g.NoEmail = append(g.NoEmail, person.ID)
			changed = true
		}
		if changed {
			r.Auth.UpdateGroup(g)
			needSave = true
		}
	}
	if needSave {
		r.Auth.Save()
	}
	r.Tx.Commit()
	return nil
}

// PostUnsubscribeList handles POST /unsubscribe/$token/$email requests.
func PostUnsubscribeList(r *util.Request, token, email string) error {
	var person *model.Person

	if r.FormValue("List-Unsubscribe") != "One-Click" {
		return PostUnsubscribe(r, token)
	}
	if person = r.Tx.FetchPersonByUnsubscribe(token); person == nil {
		return util.NotFound
	}
	for _, list := range r.Tx.FetchLists() {
		if list.Type != model.ListEmail || list.Name != email {
			continue
		}
		r.Tx.WillUpdateList(list)
		list.People[person.ID] &^= model.ListSubscribed
		list.People[person.ID] |= model.ListUnsubscribed
		r.Tx.UpdateList(list)
		r.Tx.Commit()
		return nil
	}
	return util.NotFound
}
