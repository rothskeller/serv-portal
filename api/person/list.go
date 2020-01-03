package person

import (
	"github.com/mailru/easyjson/jwriter"

	"rothskeller.net/serv/model"
	"rothskeller.net/serv/util"
)

// GetPeople handles GET /api/people requests.
func GetPeople(r *util.Request) error {
	var (
		focus *model.Team
		out   jwriter.Writer
		first = true
	)
	focus = r.Tx.FetchTeam(model.TeamID(util.ParseID(r.FormValue("team"))))
	out.RawString(`{"people":[`)
	for _, p := range r.Tx.FetchPeople() {
		if !r.Person.CanViewPerson(p) {
			continue
		}
		if focus != nil && !p.IsMember(focus) {
			continue
		}
		if first {
			first = false
		} else {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(p.ID))
		out.RawString(`,"firstName":`)
		out.String(p.FirstName)
		out.RawString(`,"lastName":`)
		out.String(p.LastName)
		out.RawString(`,"email":`)
		out.String(p.Email)
		out.RawString(`,"phone":`)
		out.String(p.Phone)
		out.RawString(`,"roles":[`)
		for i, r := range p.Roles {
			if i != 0 {
				out.RawByte(',')
			}
			out.RawString(`{"team":`)
			out.String(r.Team.Name)
			out.RawString(`,"role":`)
			out.String(r.Name)
			out.RawByte('}')
		}
		out.RawString(`]}`)
	}
	out.RawString(`],"viewableTeams":[`)
	for i, t := range r.Person.ViewableTeams() {
		if i != 0 {
			out.RawByte(',')
		}
		out.RawString(`{"id":`)
		out.Int(int(t.ID))
		out.RawString(`,"name":`)
		out.String(t.Name)
		out.RawByte('}')
	}
	out.RawString(`],"canAdd":`)
	out.Bool(r.Person.CanManageTeams())
	out.RawByte('}')
	r.Tx.Commit()
	r.Header().Set("Content-Type", "application/json; charset=utf-8")
	out.DumpTo(r)
	return nil
}
