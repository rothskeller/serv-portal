// This program reads the scc-ares-races.org/activities web site, gets the event
// calendar from it, and mirrors those events into our database.

package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
)

var eventDetailHrefRE = regexp.MustCompile(`eventdetail\.php\?id=(\d+)`)

type eventData struct {
	model.Event
	venueName string
}

func main() {
	var (
		eventsResponse *http.Response
		eventsBody     *html.Node
		eventIDs       map[string]model.EventType
		events         []*eventData
		err            error
	)
	if eventsResponse, err = http.Get("https://www.scc-ares-races.org/activities/events.php"); err != nil {
		panicf("get events.php: %s", err)
	}
	if eventsResponse.StatusCode != http.StatusOK {
		panicf("get events.php: %s", eventsResponse.Status)
	}
	if eventsBody, err = html.Parse(eventsResponse.Body); err != nil {
		panicf("parse events.php: %s", err)
	}
	db.Open("serv.db")
	if eventIDs = getEventIDs(eventsBody); len(eventIDs) == 0 {
		panicf("no events found in events.php")
	}
	events = make([]*eventData, 0, len(eventIDs))
	for id, typ := range eventIDs {
		events = append(events, getEvent(id, typ))
	}
	applyRewrites(events)
	saveEvents(events)
}

func getEventIDs(node *html.Node) (ids map[string]model.EventType) {
	tx := db.Begin()
	typemap := tx.FetchSccAresEventTypes()
	tx.Commit()
	ids = make(map[string]model.EventType)
	node = expectNode(node, html.DocumentNode)
	node = expectNode(node.FirstChild, html.DoctypeNode)
	node = expectElement(node.NextSibling, atom.Html)
	node = expectElement(node.FirstChild, atom.Head)
	node = expectElement(node.NextSibling, atom.Body)
	node = findChildElement(node, atom.Div, "id", "layout_3")
	node = findChildElement(node, atom.Div, "id", "layout_3_helper")
	node = findChildElement(node, atom.Div, "id", "content")
	node = findChildElement(node, atom.Div, "class", "currentEvents")
	for node = node.FirstChild; node != nil; node = node.NextSibling {
		var (
			n       *html.Node
			eventID string
		)
		if uninterestingNode(node) {
			continue
		}
		node = expectElement(node, atom.Table)
		n = expectElement(node.FirstChild, atom.Tbody)
		n = expectElement(n.FirstChild, atom.Tr)
		n = expectElement(n.FirstChild, atom.Td)
		n = expectElement(n.FirstChild, atom.H3)
		n = expectElement(n.FirstChild, atom.A)
		for _, attr := range n.Attr {
			if attr.Namespace == "" && attr.Key == "href" {
				if match := eventDetailHrefRE.FindStringSubmatch(attr.Val); match != nil {
					eventID = match[1]
				}
			}
		}
		if eventID == "" {
			panicf(`expected <a href="eventdetails.php?id=...">, not found`)
		}
		node = expectElement(node.NextSibling, atom.Table)
		n = expectElement(node.FirstChild, atom.Tbody)
		n = expectElement(n.FirstChild, atom.Tr)
		n = expectElement(n.NextSibling, atom.Tr)
		n = expectElement(n.FirstChild, atom.Td)
		n = expectElement(n.NextSibling, atom.Td)
		n = expectElement(n.NextSibling, atom.Td)
		n = expectNode(n.FirstChild, html.TextNode)
		if et, ok := typemap[n.Data]; ok {
			ids[eventID] = et
		} else {
			ids[eventID] = 0
			fmt.Printf("WARNING: unmapped event type %q\n", n.Data)
		}
		node = expectElement(node.NextSibling, atom.Div)
	}
	return ids
}

func getEvent(eventID string, eventType model.EventType) (event *eventData) {
	var (
		eventResponse *http.Response
		node          *html.Node
		n             *html.Node
		err           error
	)
	event = new(eventData)
	event.SccAresID = eventID
	event.Type = eventType
	event.Details = fmt.Sprintf(`For details and to register, visit <a href="https://www.scc-ares-races.org/activities/eventdetail.php?id=%s" target="_blank" rel="nofollow noopener">scc-ares-races.org</a>.`, eventID)
	if eventResponse, err = http.Get(fmt.Sprintf("https://scc-ares-races.org/activities/eventdetail.php?id=%s", eventID)); err != nil {
		panicf("get eventdetail.php?id=%s: %s", eventID, err)
	}
	if eventResponse.StatusCode != http.StatusOK {
		panicf("get eventdetail.php?id=%s: %s", eventID, eventResponse.Status)
	}
	if node, err = html.Parse(eventResponse.Body); err != nil {
		panicf("parse eventdetail.php?id=%s: %s", eventID, err)
	}
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("eventID: %s\n", eventID)
			panic(p)
		}
	}()
	node = expectNode(node, html.DocumentNode)
	node = expectNode(node.FirstChild, html.DoctypeNode)
	node = expectElement(node.NextSibling, atom.Html)
	node = expectElement(node.FirstChild, atom.Head)
	node = expectElement(node.NextSibling, atom.Body)
	node = findChildElement(node, atom.Div, "id", "layout_4")
	node = findChildElement(node, atom.Div, "id", "layout_4_helper")
	node = findChildElement(node, atom.Div, "id", "content")
	node = findChildElement(node, atom.Table, "", "")
	node = expectElement(node.FirstChild, atom.Tbody)
	node = expectElement(node.FirstChild, atom.Tr)
	n = expectElement(node.FirstChild, atom.Td)
	n = expectElement(n.FirstChild, atom.H3)
	n = expectNode(n.FirstChild, html.TextNode)
	event.Name = n.Data
	node = expectElement(node.NextSibling, atom.Tr)
	node = expectElement(node.FirstChild, atom.Td)
	n = expectElement(node.FirstChild, atom.Strong)
	n = expectNode(n.NextSibling, html.TextNode)
	if d, err := time.Parse("01/02/06", strings.TrimSpace(n.Data)); err == nil {
		event.Date = d.Format("2006-01-02")
	} else {
		panicf("date doesn't look like a date")
	}
	node = expectElement(node.NextSibling, atom.Td)
	n = expectElement(node.FirstChild, atom.Strong)
	n = expectNode(n.NextSibling, html.TextNode)
	if t, err := time.Parse("3:04 PM", strings.TrimSpace(n.Data)); err == nil {
		event.Start = t.Format("15:04")
	} else {
		panicf("time doesn't look like a time")
	}
	node = expectElement(node.NextSibling, atom.Td)
	n = expectElement(node.FirstChild, atom.Strong)
	n = expectNode(n.NextSibling, html.TextNode)
	if t, err := time.Parse("3:04 PM", strings.TrimSpace(n.Data)); err == nil {
		event.End = t.Format("15:04")
	} else {
		panicf("time doesn't look like a time")
	}
	node = expectElement(node.NextSibling, atom.Td)
	n = expectElement(node.FirstChild, atom.Strong)
	n = expectNode(n.NextSibling, html.TextNode)
	event.venueName = strings.TrimSpace(n.Data)
	return event
}

func applyRewrites(events []*eventData) {
	var (
		tx     *db.Tx
		groups []model.GroupID
		vmap   map[string]*model.Venue
		nmap   map[string]string
	)
	tx = db.Begin()
	nmap = tx.FetchSccAresEventNames()
	vmap = tx.FetchSccAresEventVenues()
	groups = []model.GroupID{tx.FetchGroupByTag(model.GroupSccAres).ID}
	for _, e := range events {
		if mapped, ok := vmap[e.venueName]; ok {
			if mapped != nil {
				e.Venue = mapped.ID
			} else {
				e.Venue = 0
			}
		} else if mapped, ok := vmap[""]; ok {
			fmt.Printf("WARNING: no mapping for venue %q, recording as \"See Event Detail Page\"\n", e.venueName)
			e.Venue = mapped.ID
		} else {
			panic("no fallback venue in database")
		}
		if rw, ok := nmap[e.Name]; ok {
			e.Name = rw
		}
		e.Groups = groups
	}
	tx.Commit()
}

func saveEvents(events []*eventData) {
	var (
		dbe     *model.Event
		futures []*model.Event
		tx      = db.Begin()
		emap    = map[string]bool{}
	)
	for _, e := range events {
		if dbe = tx.FetchEventBySccAresID(e.SccAresID); dbe != nil {
			if e.Name == dbe.Name && e.Date == dbe.Date && e.Start == dbe.Start && e.End == dbe.End && e.Venue == dbe.Venue && e.Details == dbe.Details && e.Type == dbe.Type {
				// Nothing's changed; don't save it so we don't audit.
				emap[e.SccAresID] = true
				continue
			}
			e.ID = dbe.ID
		} else if e.Name != "" {
			fmt.Printf("ADD: new event %s %s\n", e.Date, e.Name)
		}
		if e.Name != "" {
			tx.SaveEvent(&e.Event)
			emap[e.SccAresID] = true
		}
	}
	futures = tx.FetchEvents(time.Now().Add(24*time.Hour).Format("2006-01-02"),
		time.Now().Add(5*365*24*time.Hour).Format("2006-01-02"))
	for _, e := range futures {
		if e.SccAresID != "" && !emap[e.SccAresID] {
			fmt.Printf("DELETE: removed event %s %s\n", e.Date, e.Name)
			tx.DeleteEvent(e)
		}
	}
	tx.Commit()
}

func uninterestingNode(node *html.Node) bool {
	if node.Type == html.CommentNode {
		return true
	}
	if node.Type == html.TextNode && strings.TrimSpace(node.Data) == "" {
		return true
	}
	return false
}
func expectNode(node *html.Node, typ html.NodeType) *html.Node {
	for ; node != nil && node.Type != typ; node = node.NextSibling {
		if uninterestingNode(node) {
			continue
		}
		panicf("expected type %d, found %d", typ, node.Type)
	}
	if node == nil {
		panicf("expected type %d, found nothing", typ)
	}
	return node
}
func expectElement(node *html.Node, data atom.Atom) *html.Node {
	node = expectNode(node, html.ElementNode)
	if node.DataAtom != data {
		panicf("expected %s, found %s", data.String(), node.Data)
	}
	return node
}
func findChildElement(parent *html.Node, data atom.Atom, key, val string) *html.Node {
	for child := parent.FirstChild; child != nil; child = child.NextSibling {
		if child.Type == html.ElementNode && child.DataAtom == data {
			if key == "" {
				return child
			}
			for _, a := range child.Attr {
				if a.Namespace == "" && a.Key == key && a.Val == val {
					return child
				}
			}
		}
	}
	panicf(`expected <%s %s="%s">, not found`, data.String(), key, val)
	panic("not reached")
}

func panicf(f string, args ...interface{}) {
	panic(fmt.Sprintf(f, args...))
}
