package listperson

import (
	"slices"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/store/internal/phys"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/person"
)

var unmoderatedLists = []string{
	"admin", "cert", "listos", "outreach", "sares", "snap",
} // TODO this really shouldn't be hard-coded.  But then, the moderator email
// address for everything else shouldn't be either.

// ListData generates the JSON list data descriptor used by the mailing list
// software.  The JSON schema is
//
//	{
//	    "cert-63": {
//	        "senders": ["steve@rothskeller.net", ...],
//	        "moderators": ["sroth@sunnyvale.ca.gov", ...],
//	        "receivers": [
//	            {
//	                "addr": "steve@rothskeller.net",
//	                "name": "Steve Roth",
//	                "token": "...",
//	            },
//	            ...
//	        ],
//	    },
//	    ...
//	}
func ListData(storer phys.Storer) (by []byte) {
	var (
		jw    jwriter.Writer
		first = true
	)
	jw.RawByte('{')
	list.All(storer, func(l *list.List) {
		if l.Type != list.Email {
			return
		}
		if first {
			first = false
		} else {
			jw.RawByte(',')
		}
		listData(storer, &jw, l)
	})
	jw.RawByte('}')
	by, _ = jw.BuildBytes()
	return by
}
func listData(storer phys.Storer, jw *jwriter.Writer, l *list.List) {
	var (
		senders []string
		first   = true
	)
	jw.String(l.Name)
	jw.RawString(`:{"moderators":["sroth@sunnyvale.ca.gov","rothskeller@gmail.com"],"receivers":[`)
	All(storer, l.ID, person.FInformalName|person.FEmail|person.FEmail2|person.FFlags|person.FUnsubscribeToken, func(p *person.Person, sender, sub, unsub bool) {
		if sub && !unsub && p.Flags()&person.NoEmail == 0 {
			if p.Email() != "" {
				if first {
					first = false
				} else {
					jw.RawByte(',')
				}
				jw.RawString(`{"addr":`)
				jw.String(strings.ToLower(p.Email()))
				jw.RawString(`,"name":`)
				jw.String(p.InformalName())
				jw.RawString(`,"token":`)
				jw.String(p.UnsubscribeToken())
				jw.RawByte('}')
			}
			if p.Email2() != "" {
				if first {
					first = false
				} else {
					jw.RawByte(',')
				}
				jw.RawString(`{"addr":`)
				jw.String(strings.ToLower(p.Email2()))
				jw.RawString(`,"name":`)
				jw.String(p.InformalName())
				jw.RawString(`,"token":`)
				jw.String(p.UnsubscribeToken())
				jw.RawByte('}')
			}
		}
		if sender {
			if p.Email() != "" {
				senders = append(senders, p.Email())
			}
			if p.Email2() != "" {
				senders = append(senders, p.Email2())
			}
		}
	})
	jw.RawString(`],"senders":[`)
	if slices.Contains(unmoderatedLists, l.Name) {
		jw.RawString(`"*"`)
	} else {
		for i, sender := range senders {
			if i != 0 {
				jw.RawByte(',')
			}
			jw.String(sender)
		}
	}
	jw.RawString(`]}`)
}
