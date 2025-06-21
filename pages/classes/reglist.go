package classes

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func GetRegList(r *request.Request, cidstr string) {
	const classFields = class.FStart | class.FLimit | class.FReferrals | class.FType
	const classregFields = classreg.FFirstName | classreg.FLastName | classreg.FEmail | classreg.FCellPhone | classreg.FRegisteredBy | classreg.FPerson
	var (
		user *person.Person
		c    *class.Class
		regs []*classreg.ClassReg
		opts ui.PageOpts
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), classFields); c == nil {
		errpage.NotFound(r, user)
	}
	if !user.HasPrivLevel(c.Type().Org(), enum.PrivLeader) {
		errpage.Forbidden(r, user)
	}
	classreg.AllForClass(r, c.ID(), classregFields, func(cr *classreg.ClassReg) {
		regs = append(regs, cr.Clone())
	})
	opts = ui.PageOpts{
		Title:    "Class Registrations",
		MenuItem: "classes",
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		var lastRB person.ID

		main.E("div class=reglistClass>%s", c.Type().String())
		if c.Start() == "2999-12-31" {
			main.E("div class=reglistStart>Waiting List")
		} else {
			main.E("div class=reglistStart>%s", c.Start())
		}
		size := main.E("div class=reglistSize")
		if c.Limit() != 0 {
			size.E("div class=reglistSizeBar").E("div style=width:%d%%", len(regs)*100/int(c.Limit()))
			size.TF("%d of %d registered", len(regs), c.Limit())
		} else {
			size.TF("%d registered", len(regs))
		}
		if len(regs) != 0 {
			grid := main.E("div class=reglistGrid")
			grid.E("div").E("b>Registered By")
			grid.E("div").E("b>Name")
			grid.E("div").E("b>Email")
			grid.E("div").E("b>Cell Phone")
			for _, reg := range regs {
				if reg.RegisteredBy() != lastRB {
					lastRB = reg.RegisteredBy()
					grid.E("div class=reglistPerson>%s", person.WithID(r, lastRB, person.FInformalName).InformalName())
				}
				if reg.Person() != 0 {
					grid.E("div class=reglistName").E("a href=/people/%d up-target=.pageCanvas>%s %s", reg.Person(), reg.FirstName(), reg.LastName())
				} else {
					grid.E("div class=reglistName>%s %s", reg.FirstName(), reg.LastName())
				}
				grid.E("div class=reglistEmail").E("a href=mailto:%s target=_blank>%s", reg.Email(), reg.Email())
				grid.E("div class=reglistCellPhone>%s", reg.CellPhone())
			}
			main.E("div class=reglistButtons").
				E("a href=/classes/%d/lists up-layer=new up-size=grow up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Email Lists", c.ID())
			main.E("div class=reglistReferralsHeading>Referred by:")
			grid = main.E("div class=reglistReferrals")
			for _, ref := range class.AllReferrals {
				grid.E("div>%d", c.Referrals()[ref])
				grid.E("div>%s", ref.String())
			}
		}
	})
}
