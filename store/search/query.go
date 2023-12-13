package search

import (
	"sunnyvaleserv.org/portal/store/document"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/folder"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util"
)

// Search runs a search for the specified string.  The results will be of type
// *document.Document, *event.Event, *folder.Folder, *person.Person, *role.Role,
// *textmsg.TextMessage, or *venue.Venue.
func Search(storer phys.Storer, query string) (results []any, err error) {
	const eventFields = event.FID | event.FStart | event.FName
	const folderFields = folder.FID | folder.FName | folder.FViewer | folder.FURLName | folder.FParent
	const personFields = person.FID | person.FSortName | person.FPrivLevels
	const roleFields = role.FID | role.FName | role.FOrg
	const venueFields = venue.FID | venue.FName | venue.FURL

	intlres, err := phys.Search(storer, query)
	if err != nil {
		return nil, err
	}
	for _, res := range intlres {
		switch res.Type {
		case "Document":
			if d := document.WithID(storer, document.ID(util.ParseID(res.Key[1:]))); d != nil {
				results = append(results, d)
			}
		case "Event":
			if e := event.WithID(storer, event.ID(util.ParseID(res.Key[1:])), eventFields); e != nil {
				results = append(results, e)
			}
		case "Folder":
			if f := folder.WithID(storer, folder.ID(util.ParseID(res.Key[1:])), folderFields); f != nil {
				results = append(results, f)
			}
		case "Person":
			if p := person.WithID(storer, person.ID(util.ParseID(res.Key[1:])), personFields); p != nil {
				results = append(results, p)
			}
		case "Role":
			if r := role.WithID(storer, role.ID(util.ParseID(res.Key[1:])), roleFields); r != nil {
				results = append(results, r)
			}
		case "Venue":
			if v := venue.WithID(storer, venue.ID(util.ParseID(res.Key[1:])), venueFields); v != nil {
				results = append(results, v)
			}
		}
	}
	return results, nil
}
