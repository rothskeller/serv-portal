package personview

import (
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

func showPassword(r *request.Request, main *htmlb.Element, user, p *person.Person) {
	if p.Email() == "" {
		return
	}
	canChange := user.ID() == p.ID() || user.IsWebmaster()
	canReset := user.ID() != p.ID() && p.ID() != person.AdminID && user.HasPrivLevel(0, enum.PrivLeader)
	if p.Email() == "" || (!canChange && !canReset) {
		return
	}
	section := main.E("div class=personviewSection")
	sheader := section.E("div class=personviewSectionHeader")
	sheader.E("div class=personviewSectionHeaderText>Password")
	section = section.E("div class=personviewPassword")
	if canChange {
		section.E("a href=/people/%d/edpassword up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Change Password", p.ID())
	}
	if canReset {
		section.E("a href=/people/%d/pwreset up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Reset Password", p.ID())
	}
}
