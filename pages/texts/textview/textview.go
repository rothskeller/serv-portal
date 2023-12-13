package textview

import (
	"fmt"
	"net/http"
	"slices"
	"time"

	"sunnyvaleserv.org/portal/pages/errpage"
	"sunnyvaleserv.org/portal/server/auth"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/store/textmsg"
	"sunnyvaleserv.org/portal/store/textrecip"
	"sunnyvaleserv.org/portal/ui"
	"sunnyvaleserv.org/portal/util"
	"sunnyvaleserv.org/portal/util/htmlb"
	"sunnyvaleserv.org/portal/util/request"
)

type recipData struct {
	name      string
	number    string
	status    string
	timestamp time.Time
	replies   []*replyData
}
type replyData struct {
	reply     string
	timestamp time.Time
}

// Get handles GET /texts/$id requests.
func Get(r *request.Request, idstr string) {
	const textmsgFields = textmsg.FID | textmsg.FLists | textmsg.FSender | textmsg.FTimestamp | textmsg.FMessage
	var (
		user    *person.Person
		tm      *textmsg.TextMessage
		visible bool
	)
	if user = auth.SessionUser(r, 0, true); user == nil {
		return
	}
	if tm = textmsg.WithID(r, textmsg.ID(util.ParseID(idstr)), textmsgFields); tm == nil {
		errpage.NotFound(r, user)
		return
	}
	for _, tml := range tm.Lists() {
		if listperson.CanSend(r, user.ID(), tml.ID) {
			visible = true
			break
		}
	}
	if !visible {
		errpage.Forbidden(r, user)
		return
	}
	Render(r, user, tm)
}

func Render(r *request.Request, user *person.Person, tm *textmsg.TextMessage) {
	var (
		recips []*recipData
		latest time.Time
		opts   ui.PageOpts
	)
	// Load all of the data.  We need it to calculate the last modified
	// header before we start rendering.
	textrecip.AllRecipientsOfText(r, tm.ID(), person.FID|person.FSortName, func(p *person.Person, number, status string, timestamp time.Time) {
		var rd = recipData{name: p.SortName(), number: number, status: status, timestamp: timestamp}
		if timestamp.After(latest) {
			latest = timestamp
		}
		textrecip.AllRepliesFromRecipient(r, tm.ID(), p.ID(), func(reply string, timestamp time.Time) bool {
			rd.replies = append(rd.replies, &replyData{reply, timestamp})
			if timestamp.After(latest) {
				latest = timestamp
			}
			return true
		})
		slices.Reverse(rd.replies)
		recips = append(recips, &rd)
	})
	latest = latest.Truncate(time.Second)
	r.Header().Set("Last-Modified", latest.Format(time.RFC1123))
	// This page gets reloaded every 5 seconds, so it would be nice not to
	// re-render it if nothing's changed.
	if ims := r.Request.Header.Get("If-Modified-Since"); ims != "" {
		if lastsent, err := time.Parse(time.RFC1123, ims); err == nil {
			if !latest.After(lastsent) {
				r.WriteHeader(http.StatusNotModified)
				return
			}
		}
	}
	// Looks like we do need to render it.
	opts = ui.PageOpts{
		Title:    "Text Message",
		MenuItem: "texts",
		Tabs: []ui.PageTab{
			{Name: "Messages", URL: "/texts", Target: ".pageCanvas"},
			{Name: "Delivery", URL: fmt.Sprintf("/texts/%d", tm.ID()), Target: "main", Active: true},
		},
	}
	r.HTMLNoCache()
	ui.Page(r, user, opts, func(main *htmlb.Element) {
		writeMetadata(r, main, tm)
		writeDeliveries(r, main, tm, recips)
	})
}

func writeMetadata(r *request.Request, main *htmlb.Element, tm *textmsg.TextMessage) {
	meta := main.E("div class=textviewMeta")
	meta.E("div class=textviewMetaLabel>Message sent")
	meta.E("div class=textviewMetaValue>%s", tm.Timestamp().Format("2006-01-02 15:04:05"))
	meta.E("div class=textviewMetaLabel>Sent by")
	meta.E("div class=textviewMetaValue>%s", person.WithID(r, tm.Sender(), person.FInformalName).InformalName())
	meta.E("div class=textviewMetaLabel>Sent to")
	list := meta.E("div class=textviewMetaValue")
	for _, tml := range tm.Lists() {
		list.E("div>%s", tml.Name)
	}
	meta.E("div class=textviewMetaLabel>Message text")
	meta.E("div class=textviewMetaValue>%s", tm.Message())
}

func writeDeliveries(r *request.Request, main *htmlb.Element, tm *textmsg.TextMessage, recips []*recipData) {
	grid := main.E("div class=textviewGrid up-poll up-interval=5000")
	row := grid.E("div class=textviewGridHeading")
	row.E("div class=textviewGridPerson>Recipient")
	row.E("div class=textviewGridStatus>Status")
	row.E("div class=textviewGridReply>Reply")
	for _, recip := range recips {
		var statusStyle, statusText = formatStatus(recip.status, len(recip.replies) != 0)
		row = grid.E("div class=textviewGridRow")
		pbox := row.E("div class=textviewGridPerson")
		pbox.E("div class=textviewGridName>%s", recip.name)
		pbox.E("div class=textviewGridNumber>%s", recip.number)
		sbox := row.E("div class=textviewGridStatus")
		sbox.E("div class='textviewGridState %s'>%s", statusStyle, statusText)
		sbox.E("div class=textviewGridTime>%s", formatTimestamp(recip.status, recip.timestamp, recip.replies))
		rbox := row.E("div class=textviewGridReply")
		for _, reply := range recip.replies {
			rdiv := rbox.E("div")
			rdiv.E("span class=textviewGridReplyTime>%s", reply.timestamp.Format("15:04:05"))
			rdiv.T(reply.reply)
		}
	}
}

func formatStatus(status string, haveReplies bool) (string, string) {
	if haveReplies {
		return "textviewStatusReplied", "Replied"
	}
	switch status {
	case "sent":
		return "textviewStatusSent", "Sent"
	case "delivered":
		return "textviewStatusDelivered", "Delivered"
	case "undelivered":
		return "textviewStatusFailed", "Not Delivered"
	case "failed":
		return "textviewStatusFailed", "Failed"
	case "No Cell Phone":
		return "textviewStatusFailed", "No Cell Phone"
	case "queued":
		return "textviewStatusPending", "Queued"
	case "sending":
		return "textviewStatusPending", "Sending"
	default:
		return "textviewStatusPending", "Pending"
	}
}

func formatTimestamp(status string, timestamp time.Time, replies []*replyData) string {
	if len(replies) != 0 {
		return replies[len(replies)-1].timestamp.Format("15:04:05")
	}
	if status == "No Cell Phone" {
		return ""
	}
	return timestamp.Format("15:04:05")
}
