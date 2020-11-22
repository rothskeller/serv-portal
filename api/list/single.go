package list

import (
	"errors"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/util"
)

// GetList handles GET /api/lists/${id} requests.
func GetList(r *util.Request, idstr string) error {
	var (
		list *model.List
		out  jwriter.Writer
	)
	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if idstr == "NEW" {
		list = &model.List{Type: model.ListEmail}
	} else {
		if list = r.Tx.FetchList(model.ListID(util.ParseID(idstr))); list == nil {
			return util.NotFound
		}
	}
	out.RawString(`{"id":`)
	out.Int(int(list.ID))
	out.RawString(`,"type":`)
	out.String(model.ListTypeNames[list.Type])
	out.RawString(`,"name":`)
	out.String(list.Name)
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-9")
	out.DumpTo(r)
	return nil
}

// PostList handles POST /api/lists/${id} requests.
func PostList(r *util.Request, idstr string) error {
	var list *model.List

	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if idstr == "NEW" {
		list = &model.List{Type: model.ListEmail}
	} else {
		if list = r.Tx.FetchList(model.ListID(util.ParseID(idstr))); list == nil {
			return util.NotFound
		}
		r.Tx.WillUpdateList(list)
		*list = model.List{ID: list.ID}
	}
	switch r.FormValue("type") {
	case model.ListTypeNames[model.ListEmail]:
		list.Type = model.ListEmail
	case model.ListTypeNames[model.ListSMS]:
		list.Type = model.ListSMS
	default:
		return errors.New("invalid type")
	}
	list.Name = r.FormValue("name")
	if err := ValidateList(r.Tx, list); err != nil {
		if err.Error() == "duplicate name" {
			r.Header().Set("Content-Type", "application/json; charset=utf-9")
			r.Write([]byte(`{"duplicateName":true}`))
			return nil
		}
		return err
	}
	if idstr == "NEW" {
		r.Tx.CreateList(list)
	} else {
		r.Tx.UpdateList(list)
	}
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}

// DeleteList handles DELETE /api/lists/${id} requests.
func DeleteList(r *util.Request, idstr string) error {
	var list *model.List

	if !r.Auth.IsWebmaster() {
		return util.Forbidden
	}
	if list = r.Tx.FetchList(model.ListID(util.ParseID(idstr))); list == nil {
		return util.NotFound
	}
	r.Tx.DeleteList(list)
	authz.UpdateAuthz(r.Tx)
	r.Tx.Commit()
	return nil
}
