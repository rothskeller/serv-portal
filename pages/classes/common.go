package classes

import (
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func getClassesCommon(r *request.Request, user *person.Person, main *htmlb.Element, ctype class.Type) {
	var (
		classes  *htmlb.Element
		langFlag = class.FEnDesc
	)
	if r.Language == "es" {
		langFlag = class.FEsDesc
	}
	class.AllFuture(r, ctype, class.FID|class.FLimit|class.FStart|langFlag, func(c *class.Class) {
		if classes == nil {
			classes = main.E("div class=classesRegisterGrid")
		}
		if r.Language == "es" {
			classes.E("div").R(c.EsDesc())
		} else {
			classes.E("div").R(c.EnDesc())
		}
		if user.HasPrivLevel(ctype.Org(), enum.PrivLeader) {
			classes.E("div").E("a href=/classes/%d/reglist up-target=main class='sbtn sbtn-primary sbtn-small'>Registrations", c.ID())
		} else if classreg.ClassIsFull(r, c.ID()) {
			classes.E("div class=classesFull").R(r.Loc("This session is full."))
		} else {
			classes.E("div").E("a href=/classes/%d/register up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-primary sbtn-small'", c.ID()).R(r.Loc("Sign Up"))
		}
	})
	if classes == nil {
		main.E("div").R("No sessions of this class are currently scheduled.")
	}
	main.E("div class=classesSERV").R(r.Loc("This class is presented by Sunnyvale Emergency Response Volunteers (SERV), the volunteer arm of the Sunnyvale Office of Emergency Services."))
}
