package main

import (
	"fmt"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
)

// These volunteers are in our database, but they also do non-SERV volunteer
// work for DPS, so even if they stop being active for SERV, this automation
// should not mark them as Inactive.  Norma McConnell will do that manually if
// it ever becomes appropriate.
var doNotDeactivate = map[string]bool{
	"11240433": true, // Asselin, Jim
	"1904515":  true, // Lea, Irene
	"15623965": true, // Macmillan, Brian
	"15618603": true, // Sanchez, Jr., Leopoldo
}

// markActive makes sure that all active volunteers in our database are marked
// active in Volgistics.
func markActive(st *store.Store, loginID string) {
	var inactiveKeys = make(map[string]bool)

	// The people we consider "active" are the ones who are members of one
	// of our orgs, and who have a Volgistics ID.
	person.All(st, person.FID|person.FVolgisticsID|person.FPrivLevels|person.FInformalName, func(p *person.Person) {
		if p.VolgisticsID() != 0 && p.HasPrivLevel(0, enum.PrivMember) {
			setStatusToActive(loginID, p.InformalName(), p.VolgisticsID())
		} else {
			key := findVolunteerInVolgistics(loginID, p.InformalName(), p.VolgisticsID(), false)
			inactiveKeys[key] = true
		}
	})
	markInactive(st, loginID, inactiveKeys)
}

func setStatusToActive(loginID, name string, volunteerID uint) {
	var (
		key string
		val url.Values
		doc *goquery.Document
		sel *goquery.Selection
		was string
	)
	if key = findVolunteerInVolgistics(loginID, name, volunteerID, true); key == "" {
		return
	}
	val = make(url.Values)
	val.Add("ID", loginID)
	val.Add("KEY", key)
	val.Add("FB", "0")
	delay()
	doc = checkResponse(client.Get("https://www.volgistics.com/ex/core.dll/volunteers?" + val.Encode()))
	if sel = doc.Find(`select[name="F18"] option[selected]`); sel.Length() != 1 {
		fmt.Fprintln(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no selected status")
		os.Exit(1)
	}
	if was = sel.Text(); was == "Active" {
		return // already marked active
	}
	val.Add("zipDocs", "on")
	val.Add("Save", "Save")
	val.Add("F18", "Active")
	for _, k := range []string{
		"TAGS", "Action", "F18or", "F20or", "A5or", "F5or", "Cg", "F1",
		"F2", "F4", "G1", "G2", "G3", "G4", "A1", "A2", "A3", "A4",
		"A6", "A19", "A7", "A9", "A11", "A13", "A15", "A17", "4-1",
	} {
		if sel = doc.Find(fmt.Sprintf(`input[name="%s"]`, k)); sel.Length() < 1 {
			fmt.Fprintf(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no field %s\n", k)
			os.Exit(1)
		}
		val.Add(k, sel.First().AttrOr("value", ""))
	}
	if sel = doc.Find(`select[name="A5"] option[selected]`); sel.Length() != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no field A5\n")
		os.Exit(1)
	}
	val.Add("A5", sel.AttrOr("value", ""))
	delay()
	checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Core", val))
	fmt.Printf("%s - %s => Active\n", name, was)
}

// markInactive ensures that volunteers who are no longer active with us are
// marked Inactive in Volgistics.
func markInactive(st *store.Store, loginID string, inactiveKeys map[string]bool) {
	var (
		val  url.Values
		doc  *goquery.Document
		rows *goquery.Selection
	)
	// Get a list of all active volunteers.
	val = make(url.Values)
	val.Add("ID", loginID)
	val.Add("FB", "!S=1!18-105^")
	val.Add("NA", "64")
	val.Add("Iop", "")
	val.Add("TAGS", "")
	val.Add("BS", "105")
	val.Add("BT", "0")
	val.Add("BC", "0")
	val.Add("BX", "0")
	val.Add("BG", "0")
	val.Add("BK", "0")
	val.Add("findnormal", "")
	val.Add("Fbt1", "")
	val.Add("Fbt4", "")
	val.Add("Fbt2", "")
	val.Add("Fbt3", "")
	val.Add("Fbt10", "")
	doc = checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers", val))
	rows = doc.Find("#volTable tbody tr")
	rows.Each(func(_ int, row *goquery.Selection) {
		key := row.Children().First().Text()
		if !inactiveKeys[key] {
			// This person is either active in SERV, or not in
			// SERV's database (which means they're a non-SERV DPS
			// volunteer).  In either case we don't touch them.
			return
		}
		if doNotDeactivate[key] {
			// This person is inactive in SERV but active in a
			// non-SERV DPS volunteer role.  We don't touch them.
			return
		}
		setStatusToInactive(loginID, key, row.Children().Get(2).FirstChild.Data)
	})
}

func setStatusToInactive(loginID, key, name string) {
	var (
		val url.Values
		doc *goquery.Document
		sel *goquery.Selection
		was string
	)
	val = make(url.Values)
	val.Add("ID", loginID)
	val.Add("KEY", key)
	val.Add("FB", "0")
	delay()
	doc = checkResponse(client.Get("https://www.volgistics.com/ex/core.dll/volunteers?" + val.Encode()))
	if sel = doc.Find(`select[name="F18"] option[selected]`); sel.Length() != 1 {
		fmt.Fprintln(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no selected status")
		os.Exit(1)
	}
	if was = sel.Text(); was == "Inactive" {
		return
	}
	val.Add("zipDocs", "on")
	val.Add("Save", "Save")
	val.Add("F18", "Inactive")
	for _, k := range []string{
		"TAGS", "Action", "F18or", "F20or", "A5or", "F5or", "Cg", "F1",
		"F2", "F4", "G1", "G2", "G3", "G4", "A1", "A2", "A3", "A4",
		"A6", "A19", "A7", "A9", "A11", "A13", "A15", "A17",
	} {
		if sel = doc.Find(fmt.Sprintf(`input[name="%s"]`, k)); sel.Length() < 1 {
			fmt.Fprintf(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no field %s\n", k)
			os.Exit(1)
		}
		val.Add(k, sel.First().AttrOr("value", ""))
	}
	for _, k := range []string{
		"4-1",
	} {
		if sel = doc.Find(fmt.Sprintf(`input[name="%s"]`, k)); sel.Length() < 1 {
			continue
		}
		val.Add(k, sel.First().AttrOr("value", ""))
	}
	if sel = doc.Find(`select[name="A5"] option[selected]`); sel.Length() != 1 {
		fmt.Fprintf(os.Stderr, "ERROR: can't parse Volgistics volunteer status: no field A5\n")
		os.Exit(1)
	}
	val.Add("A5", sel.AttrOr("value", ""))
	delay()
	checkResponse(client.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Core", val))
	fmt.Printf("%s - %s => Inactive\n", name, was)
}
