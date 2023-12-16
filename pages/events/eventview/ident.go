package eventview

import (
	"time"

	"sunnyvaleserv.org/portal/pages/events/langdate"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/ui/orgdot"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const identEventFields = event.FStart | event.FName | event.FActivation
const identTaskFields = task.FOrg | task.FFlags

func showIdent(r *request.Request, main *htmlb.Element, e *event.Event, ts []*task.Task) {
	names := main.E("div class=eventviewIdent")
	left := names.E("div class=eventviewIdentLeft")
	line1 := left.E("div class=eventviewIdentL1")
	line1.E("span class=eventviewIdentName>%s", e.Name())
	if act := e.Activation(); act != "" {
		line1.E("span class=eventviewIdentActivation>%s", e.Activation())
	}
	var orgs = make([]bool, enum.NumOrgs)
	for _, t := range ts {
		orgs[t.Org()] = true
	}
	dots := line1.E("span class=eventviewIdentOrgs")
	for org, show := range orgs {
		if show {
			orgdot.OrgDot(r, dots, enum.Org(org))
		}
	}
	dsw := ts[0].Flags()&task.CoveredByDSW != 0
	for _, t := range ts {
		if t.Flags()&task.CoveredByDSW == 0 {
			dsw = false
		}
	}
	if dsw {
		line1.E("span class=eventviewIdentDSW>DSW")
	}
	date, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
	left.E("div class=eventviewIdentDate>%s", langdate.LangDate(r, date))
}
