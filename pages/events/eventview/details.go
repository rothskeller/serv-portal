package eventview

import (
	"time"

	"sunnyvaleserv.org/portal/server/l10n"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/event"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/shift"
	"sunnyvaleserv.org/portal/store/task"
	"sunnyvaleserv.org/portal/store/venue"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

const (
	detailsEventFields = event.FID | event.FStart | event.FEnd | event.FVenueURL | event.FDetails
	detailsTaskFields  = task.FOrg
	detailsVenueFields = venue.FName | venue.FURL
)

func showDetails(r *request.Request, main *htmlb.Element, user *person.Person, e *event.Event, ts []*task.Task) {
	var hasShifts bool
	editable := user.HasPrivLevel(ts[0].Org(), enum.PrivLeader)
	for _, t := range ts {
		if !user.HasPrivLevel(t.Org(), enum.PrivLeader) {
			editable = false
		}
		if !hasShifts && shift.ExistsForTask(r, t.ID()) {
			hasShifts = true
		}
	}
	if e.Start() == "00:00" && e.End() == "00:00" && e.Venue() == 0 && e.Details() == "" && !editable {
		return
	}
	section := main.E("div class=eventviewSection")
	sheader := section.E("div class=eventviewSectionHeader")
	sheader.E("div class=eventviewSectionHeaderText>Details")
	if editable {
		sheader.E("div class=eventviewSectionHeaderEdit").
			E("a href=/events/%d/eddetails up-layer=new up-size=grow up-dismissable=key up-history=false class='sbtn sbtn-small sbtn-primary'>Edit", e.ID())
	}
	bdiv := section.E("div class=eventviewDetails")
	datet, _ := time.ParseInLocation("2006-01-02T15:04", e.Start(), time.Local)
	date := l10n.LocalizeDate(datet, r.Language)
	if e.Start()[11:] != "00:00" || e.End()[11:] != "00:00" {
		if e.Start() != e.End() {
			bdiv.E("div class=eventviewDetailsTime").R(date+" ").TF(r.Loc("from %s to %s"), e.Start()[11:], e.End()[11:])
		} else {
			bdiv.E("div class=eventviewDetailsTime").R(date+" ").TF(r.Loc("at %s"), e.Start()[11:])
		}
	} else {
		bdiv.E("div class=eventviewDetailsTime>%s", date)
	}
	if e.Venue() != 0 {
		v := venue.WithID(r, e.Venue(), detailsVenueFields)
		vdiv := bdiv.E("div class=eventviewDetailsVenue")
		if e.VenueURL() != "" {
			vdiv.E("a href=%s target=_blank>%s", e.VenueURL(), v.Name())
		} else if v.URL() != "" {
			vdiv.E("a href=%s target=_blank>%s", v.URL(), v.Name())
		} else {
			vdiv.T(v.Name())
		}
	} else {
		bdiv.E("div class=eventviewDetailsVenue").R(r.Loc("Location TBD"))
	}
	if e.Details() != "" {
		bdiv.E("div class=eventviewDetailsDetails").R(e.Details())
	}
	if editable {
		bdiv.E("div class=eventviewDetailsButtons").
			E("a href=/events/eventlists/%d up-layer=new up-size=grow up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Email Lists", e.ID())
	}
}

// showEventEmailLists displays the email lists for the event.
func showEventEmailLists(r *request.Request, body *htmlb.Element, e *event.Event, hasShifts bool) {
	heading := body.E("div class=eventviewTaskHeading").R(r.Loc("Email Lists"))
	heading.E("a href=/events/eventlists/%d up-layer=new up-size=grow up-dismissable=false up-history=false class='sbtn sbtn-xsmall sbtn-primary'>Explanation", e.ID())
	addrs := body.E("div class=eventviewEmails")
	addrs.E("div>event-%d-invited@SunnyvaleSERV.org", e.ID())
	if hasShifts {
		addrs.E("div>event-%d-signedup@SunnyvaleSERV.org", e.ID())
	}
	addrs.E("div>event-%d-signedin@SunnyvaleSERV.org", e.ID())
}
