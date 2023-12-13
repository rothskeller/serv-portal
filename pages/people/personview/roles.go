package personview

import (
	"strings"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/personrole"
	"sunnyvaleserv.org/portal/store/role"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type badgedata struct {
	badge  string
	titles []string
}

func showRoles(r *request.Request, main *htmlb.Element, user, p *person.Person) {
	var (
		badges   []*badgedata
		held     = map[role.ID]bool{}
		editable = user.HasPrivLevel(0, enum.PrivLeader) && p.ID() != person.AdminID
	)
	personrole.RolesForPerson(r, p.ID(), role.FID, func(rl *role.Role, _ bool) {
		held[rl.ID()] = true
	})
	role.All(r, role.FID|role.FTitle|role.FOrg, func(rl *role.Role) {
		var (
			badge string
			found bool
		)
		if !held[rl.ID()] || rl.Title() == "" {
			return
		}
		switch rl.Org() {
		case enum.OrgAdmin:
			if strings.HasPrefix(rl.Title(), "OES") {
				badge = "dps"
			} else {
				badge = "serv"
			}
		case enum.OrgCERTD, enum.OrgCERTT:
			badge = "cert"
		case enum.OrgListos:
			badge = "listos"
		case enum.OrgSARES:
			badge = "sares"
		case enum.OrgSNAP:
			badge = "snap"
		}
		for _, bd := range badges {
			if bd.badge == badge {
				bd.titles = append(bd.titles, rl.Title())
				found = true
				break
			}
		}
		if !found {
			badges = append(badges, &badgedata{badge: badge, titles: []string{rl.Title()}})
		}
	})
	if len(badges) == 0 && !editable {
		return
	}
	section := main.E("div class=personviewSection")
	sheader := section.E("div class=personviewSectionHeader")
	if len(badges) > 1 || (len(badges) == 1 && len(badges[0].titles) > 1) {
		sheader.E("div class=personviewSectionHeaderText>SERV Roles")
	} else {
		sheader.E("div class=personviewSectionHeaderText>SERV Role")
	}
	if editable {
		sheader.E("div class=personviewSectionHeaderEdit").
			E("a href=/people/%d/edroles up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Edit", p.ID())
	}
	bdiv := section.E("div class=personviewRoles")
	for _, b := range badges {
		bdiv.E("img class=personviewRolesBadge src=%s", ui.AssetURL(b.badge+"-badge.png"))
		tdiv := bdiv.E("div class=personviewRolesTitles")
		for _, t := range b.titles {
			tdiv.E("div class=personviewRolesTitle>%s", t)
		}
	}
	if len(badges) == 0 {
		bdiv.E("div class=personviewRolesTitles").
			E("div class=personviewRolesTitle>No current role in any SERV org.")
	}
}
