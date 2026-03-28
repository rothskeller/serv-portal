package classes

import (
	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/class"
	"sunnyvaleserv.org/portal/store/classreg"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func GetRegList(r *request.Request, cidstr string) {
	const classFields = class.FStart | class.FLimit | class.FReferrals | class.FType | class.FRegURL | class.FRole
	var (
		user *person.Person
		c    *class.Class
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if c = class.WithID(r, class.ID(util.ParseID(cidstr)), classFields); c == nil || c.RegURL() != "" {
		errpage.NotFound(r, user)
	}
	if !user.HasPrivLevel(c.Type().Org(), enum.PrivLeader) {
		errpage.Forbidden(r, user)
	}
	RenderRegList(r, user, c)
}

func RenderRegList(r *request.Request, user *person.Person, c *class.Class) {
	const classregFields = classreg.FID | classreg.FFirstName | classreg.FLastName | classreg.FEmail | classreg.FCellPhone | classreg.FRegisteredBy | classreg.FPerson | classreg.FWaitlist
	var (
		regs     []*classreg.ClassReg
		waitlist int
		opts     ui.PageOpts
	)
	classreg.AllForClass(r, c.ID(), classregFields, func(cr *classreg.ClassReg) {
		regs = append(regs, cr.Clone())
		if cr.Waitlist() {
			waitlist++
		}
	})
	opts = ui.PageOpts{
		Title:    "Class Registrations",
		MenuItem: "classes",
	}
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		var lastRB person.ID

		main.E("div class=reglistClass>%s", c.Type().String())
		main.E("div class=reglistStart>%s", c.Start())
		size := main.E("div class=reglistSize")
		if c.Limit() != 0 {
			size.E("div class=reglistSizeBar").E("div style=width:%d%%", (len(regs)-waitlist)*100/int(c.Limit()))
			size.TF("%d of %d registered", len(regs)-waitlist, c.Limit())
		} else {
			size.TF("%d registered", len(regs)-waitlist)
		}
		grid := main.E("div class=reglistGrid")
		grid.E("div").E("b>Name")
		grid.E("div").E("b>Email")
		grid.E("div").E("b>Cell Phone")
		grid.E("div").E("b>Status")
		grid.E("div")
		for _, reg := range regs {
			var proxy string
			if reg.RegisteredBy() == lastRB {
				proxy = "+ "
			}
			lastRB = reg.RegisteredBy()
			if reg.Person() != 0 {
				grid.E("div class=reglistName").R(proxy).E("a href=/people/%d up-target=.pageCanvas>%s %s", reg.Person(), reg.FirstName(), reg.LastName())
			} else {
				grid.E("div class=reglistName>%s%s %s", proxy, reg.FirstName(), reg.LastName())
			}
			grid.E("div class=reglistEmail").E("a href=mailto:%s target=_blank>%s", reg.Email(), reg.Email())
			grid.E("div class=reglistCellPhone>%s", reg.CellPhone())
			if reg.Person() != 0 && c.Role() != 0 && hasRole(r, reg.Person(), c.Role()) {
				grid.E("div>In Class")
			} else if !reg.Waitlist() {
				grid.E("div>Registered")
			} else {
				grid.E("div>Waitlist")
			}
			grid.E("div").E("a href=/classes/regedit/%d up-layer=new up-size=grow up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Edit", reg.ID())
		}
		if len(regs) != 0 {
			buttons := main.E("div class=reglistButtons")
			buttons.E("a href=/classes/%d/lists up-layer=new up-size=grow up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Email Lists", c.ID())
			main.E("div class=reglistReferralsHeading>Referred by:")
			grid := main.E("div class=reglistReferrals")
			for _, ref := range class.AllReferrals {
				grid.E("div>%d", c.Referrals()[ref])
				grid.E("div>%s", ref.String())
			}
		}
	})
}

func hasRole(r *request.Request, p person.ID, rl role.ID) bool {
	held, _ := personrole.PersonHasRole(r, p, rl)
	return held
}
