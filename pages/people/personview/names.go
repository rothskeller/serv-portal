package personview

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const namesPersonFields = person.FInformalName | person.FFormalName | person.FCallSign | person.FPronouns

func showNames(r *request.Request, main *htmlb.Element, user, p *person.Person) {
	names := main.E("div class=personviewNames")
	ifc := names.E("div class=personviewNamesIFC")
	informal := ifc.E("div class=personviewNamesIC").
		E("div class=personviewNamesInformal>%s", p.InformalName())
	if p.CallSign() != "" {
		informal.E("div class=personviewNamesCall>%s", p.CallSign())
	}
	if p.FormalName() != p.InformalName() && p.Pronouns() != "" {
		ifc.E("div class=personviewNamesFormal>%s (%s)", p.FormalName(), p.Pronouns())
	} else if p.FormalName() != p.InformalName() {
		ifc.E("div class=personviewNamesFormal>%s", p.FormalName())
	} else if p.Pronouns() != "" {
		ifc.E("div class=personviewNamesFormal>(%s)", p.Pronouns())
	}
	if user.ID() == p.ID() || user.HasPrivLevel(0, enum.PrivLeader) {
		names.E("div class=personviewNamesEdit").
			E("a href=/people/%d/ednames up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'", p.ID()).R(r.Loc("Edit"))
	}
}
