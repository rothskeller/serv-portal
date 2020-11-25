package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/config"
)

type pinfo struct {
	Name         string
	VolgisticsID int
	Minutes      map[int]uint16
	Total        uint16
}

var client http.Client

func submitHours(tx *store.Tx) {
	var (
		mstr   string
		eatt   = make(map[model.EventID]map[model.PersonID]model.AttendanceInfo)
		people = make(map[model.PersonID]*pinfo)
	)
	mstr = time.Time(mflag).Format("2006-01")
	for _, e := range tx.FetchEvents(mstr+"-01", mstr+"-31") {
		assn := orgToAssignment[e.Org]
		if assn == 0 {
			continue
		}
		eatt[e.ID] = tx.FetchAttendanceByEvent(e)
		for pid, ai := range eatt[e.ID] {
			if len(pflags) != 0 && !pflags[pid] {
				continue
			}
			if ai.Minutes == 0 || ai.Type == model.AttendAsAuditor || ai.Type == model.AttendAsStudent {
				continue
			}
			if people[pid] == nil {
				p := tx.FetchPerson(pid)
				people[pid] = &pinfo{
					Name:         p.InformalName,
					VolgisticsID: p.VolgisticsID,
					Minutes:      make(map[int]uint16),
				}
			}
			people[pid].Minutes[assn] += ai.Minutes
		}
	}
	for pid, pi := range people {
		if pi.VolgisticsID == 0 {
			delete(people, pid)
		}
	}
	submitToVolgistics(people, time.Date(time.Time(mflag).Year(), time.Time(mflag).Month()+1, 1, 0, 0, 0, 0, time.Local).Add(-time.Second))
}

func submitToVolgistics(people map[model.PersonID]*pinfo, date time.Time) {
	var (
		doc *goquery.Document
		id  string
	)
	client.Jar = &cookiejar.Jar{}

	// First, log in to Volgistics.  For subsequent requests, we need the
	// session ID, which we extract from the URL parameters in the
	// redirect response sent by the Login request.
	var loginForm = url.Values{}
	loginForm.Add("vtOS", "")
	loginForm.Add("1-0", config.Get("volgisticsAccount"))
	loginForm.Add("1-1", config.Get("volgisticsEmail"))
	loginForm.Add("Password", config.Get("volgisticsPassword"))
	loginForm.Add("Submit1", "Login")
	delay()
	doc = checkResponse(client.PostForm("https://www.volgistics.com/ex/login.dll/?NavTo=Start", loginForm))
	if form := doc.Find(`form[action="https://www.volgistics.com/ex/login.dll/?NavTo=OVR"]`); form.Length() == 1 {
		// It's complaining that we're already logged in.  Tell it to
		// go ahead anyway.
		var form2 = url.Values{}
		form2.Add("ID-2", form.Find("#ID-2").AttrOr("value", ""))
		form2.Add("Continue", "Continue")
		delay()
		doc = checkResponse(client.PostForm("https://www.volgistics.com/ex/login.dll/?NavTo=OVR", form2))
	}
	if nav1 := doc.Find("#nav1"); nav1.Length() == 1 {
		if src := nav1.AttrOr("src", ""); src != "" {
			if u, err := url.Parse(src); err == nil {
				id = u.Query().Get("ID")
			}
		}
	}
	if id == "" {
		fmt.Fprintln(os.Stderr, "ERROR logging into Volgistics: no ID in response")
		os.Exit(1)
	}
	for _, pi := range people {
		if pi.VolgisticsID == 0 {
			continue
		}
		submitPersonToVolgistics(id, date, pi)
	}
}

func submitPersonToVolgistics(id string, date time.Time, pi *pinfo) {
	var (
		doc  *goquery.Document
		row  *goquery.Selection
		key  string
		name string
	)

	// Find the volunteer.
	var findForm = url.Values{}
	findForm.Add("ID", id)
	findForm.Add("FB", "0")
	findForm.Add("NA", "")
	findForm.Add("Iop", "")
	findForm.Add("TAGS", "")
	findForm.Add("BS", "0")
	findForm.Add("BT", "0")
	findForm.Add("BC", "0")
	findForm.Add("BX", "0")
	findForm.Add("BG", "0")
	findForm.Add("BK", "0")
	findForm.Add("findnormal", "")
	findForm.Add("Fbt1", "")
	findForm.Add("Fbt4", "")
	findForm.Add("Fbt2", strconv.Itoa(pi.VolgisticsID))
	findForm.Add("Find1", "Go")
	findForm.Add("Fbt3", "")
	findForm.Add("Fbt10", "")
	delay()
	doc = checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers", findForm))
	if row = doc.Find("#volTable tbody tr"); row.Length() != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: volunteer %s (%d) not found in Volgistics\n", pi.Name, pi.VolgisticsID)
		return
	}
	if row = row.Children(); row.Length() < 3 {
		fmt.Fprintln(os.Stderr, "ERROR: can't parse Volgistics volunteer search results: wrong row length")
		os.Exit(1)
	}
	if key = row.First().Text(); key == "" {
		fmt.Fprintln(os.Stderr, "ERROR: can't parse Volgistics volunteer search results: no key")
		os.Exit(1)
	}
	name = strings.SplitN(row.Eq(2).Text(), ",", 2)[0]
	if !strings.Contains(pi.Name, name) {
		fmt.Fprintf(os.Stderr, "WARNING: volunteer %s (%d) name mismatch: Volgistics has %s\n", pi.Name, pi.VolgisticsID, row.Eq(2).Text())
	}

	// Read the page with that volunteer's assignments and service record.
	var monthexpand = date.Format("200601") + "00000000"
	var volPage = url.Values{}
	volPage.Add("ID", id)
	volPage.Add("KEY", key)
	volPage.Add("FB", "0")
	volPage.Add("HNUM", "0")
	volPage.Add("F0", "")
	volPage.Add("T0", "")
	volPage.Add("H0", "")
	volPage.Add("RESET", monthexpand)
	volPage.Add("LRESET", "")
	volPage.Add("A0", "Blank")
	volPage.Add("H1FD", "")
	volPage.Add("MM0", "0.00")
	volPage.Add("scrollTo", "#C5_26")
	delay()
	doc = checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Hours", volPage))

	// Handle each assignment type.
	datefmt := date.Format("01-02-2006")
	rows := doc.Find("td.volgistics487").FilterFunction(func(_ int, node *goquery.Selection) bool {
		return node.Text() == datefmt
	}).Parent()
ASSN:
	for a, label := range assnToLabel {
		var (
			found       bool
			disposition string
			updateForm  = url.Values{}
		)
		updateForm.Add("ID", id)
		updateForm.Add("KEY", key)
		updateForm.Add("FB", "0")
		updateForm.Add("F0", "")
		updateForm.Add("T0", "")
		updateForm.Add("H0", "")
		updateForm.Add("RESET", "")
		updateForm.Add("LRESET", monthexpand)
		updateForm.Add("scrollTo", "#C5_26")
		for i := 0; i < rows.Length(); i++ {
			cols := rows.Eq(i).Children()
			if cols.Eq(1).Text() != label {
				continue
			}
			var hnum = strings.Split(cols.Eq(3).Children().AttrOr("name", ":"), ":")[1]
			if pi.Minutes[a] != 0 {
				var current float64
				fmt.Sscanf(cols.Eq(2).Text(), "%f", &current)
				if hours := float64(pi.Minutes[a]) / 60; hours != current {
					updateForm.Add("HNUM", hnum)
					updateForm.Add("ODATE", datefmt)
					updateForm.Add("A0", strconv.Itoa(a))
					updateForm.Add("H1FD", datefmt)
					updateForm.Add("MM0", fmt.Sprintf("%f", hours))
					updateForm.Add("Save", "Save")
					disposition = "updated"
				} else {
					fmt.Printf("%s - %s - no change\n", pi.Name, assnToName[a])
					continue ASSN
				}
			} else {
				dskey := cols.Eq(4).Children().AttrOr("name", "")
				updateForm.Add("HNUM", "0")
				updateForm.Add(dskey+".x", "1")
				updateForm.Add(dskey+".y", "1")
				updateForm.Add("A0", "Blank")
				updateForm.Add("H1FD", "")
				updateForm.Add("MM0", "0.00")
				disposition = "deleted (must verify manually)"
			}
			found = true
			break
		}
		if !found && pi.Minutes[a] == 0 {
			continue
		}
		if !found {
			updateForm.Add("A0", strconv.Itoa(a))
			updateForm.Add("H1FD", datefmt)
			updateForm.Add("MM0", fmt.Sprintf("%f", float64(pi.Minutes[a])/60))
			updateForm.Add("Save", "Save")
			disposition = "added"
		}
		delay()
		checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Hours", updateForm))
		fmt.Printf("%s - %s - %s\n", name, assnToName[a], disposition)
	}
}

var lastRequest time.Time

func delay() {
	time.Sleep(time.Until(lastRequest.Add(500 * time.Millisecond)))
	lastRequest = time.Now()
}

func checkResponse(resp *http.Response, err error) (doc *goquery.Document) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR contacting Volgistics: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		fmt.Fprintf(os.Stderr, "ERROR response from Volgistics: %s\n", resp.Status)
		os.Exit(1)
	}
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR parsing response from Volgistics: %s\n", err)
	}
	return doc
}
