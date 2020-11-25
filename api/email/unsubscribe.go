package email

import (
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// PostUnsubscribeList handles POST /unsubscribe/$token/$email requests.  These
// come from List-Unsubscribe headers in emails sent to the lists.
func PostUnsubscribeList(r *util.Request, token, email string) error {
	var person *model.Person

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
