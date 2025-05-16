package volgistics

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"

	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/config"
)

// Assignment identifies an assignment in Volgistics: roughly analogous to the
// SERV organizations, but both CERT-D and CERT-T have the same assignment.
type Assignment int

const (
	Admin  Assignment = 1052
	CERT   Assignment = 1047
	Listos Assignment = 1048
	SARES  Assignment = 399
	SNAP   Assignment = 373
)

func (a Assignment) String() string {
	switch a {
	case Admin:
		return "Admin"
	case CERT:
		return "CERT"
	case Listos:
		return "Listos"
	case SARES:
		return "SARES"
	case SNAP:
		return "SNAP"
	default:
		return ""
	}
}

func (a Assignment) label() string {
	switch a {
	case Admin:
		return "SERV Admin [EMERGENCY PREPAREDNESS]"
	case CERT:
		return "CERT [EMERGENCY PREPAREDNESS]"
	case Listos:
		return "LISTOS [EMERGENCY PREPAREDNESS]"
	case SARES:
		return "SARES Volunteer [EMERGENCY PREPAREDNESS]"
	case SNAP:
		return "SNAP Volunteer [EMERGENCY PREPAREDNESS]"
	default:
		return ""
	}
}

var allAssignments = []Assignment{Admin, CERT, Listos, SARES, SNAP}

var OrgToAssignment = map[enum.Org]Assignment{
	enum.OrgAdmin:  Admin,
	enum.OrgCERTD:  CERT,
	enum.OrgCERTT:  CERT,
	enum.OrgListos: Listos,
	enum.OrgSARES:  SARES,
	enum.OrgSNAP:   SNAP,
}

var httpClient http.Client

type Client struct {
	id string
}

func Login() (client *Client, err error) {
	var doc *goquery.Document
	httpClient.Jar = &cookiejar.Jar{}

	// Log in to Volgistics.  For subsequent requests, we need the session
	// ID, which we extract from the URL parameters in the redirect response
	// sent by the Login request.
	loginForm := url.Values{}
	loginForm.Add("vtOS", "")
	loginForm.Add("1-0", config.Get("volgisticsAccount"))
	loginForm.Add("1-1", config.Get("volgisticsEmail"))
	loginForm.Add("Password", config.Get("volgisticsPassword"))
	loginForm.Add("Submit1", "Login")
	delay()
	if doc, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/login.dll/?NavTo=Start", loginForm)); err != nil {
		return nil, err
	}
	if form := doc.Find(`form[action="https://www.volgistics.com/ex/login.dll/?NavTo=OVR"]`); form.Length() == 1 {
		// It's complaining that we're already logged in elsewhere.
		// Tell it to go ahead anyway.
		form2 := url.Values{}
		form2.Add("ID-2", form.Find("#ID-2").AttrOr("value", ""))
		form2.Add("Continue", "Continue")
		delay()
		if doc, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/login.dll/?NavTo=OVR", form2)); err != nil {
			return nil, err
		}
	}
	if nav1 := doc.Find("#nav1"); nav1.Length() == 1 {
		if src := nav1.AttrOr("src", ""); src != "" {
			if u, err := url.Parse(src); err == nil {
				client = &Client{id: u.Query().Get("ID")}
			}
		}
	}
	if client == nil {
		return nil, errors.New("Volgistics login: no ID in response")
	}
	return client, nil
}

// SubmitHours submits the hours for a particular volunteer to Volgistics.
// date specifies the month for which hours are being submitted (the date within
// the month and the time are unused).  name is the volunteer name.  vid is the
// volunteer ID.  minutes is a map from assignment type (i.e., org, more or
// less) to the number of minutes spent on that assignment during the month (a
// non-negative multiple of 30).
//
// If the function is successful, it returns a map from assignment type to a
// string indicating what happened with that assignment type ("added",
// "deleted", "updated", or "no change").  It may also return an "error", which
// is a warning indicating that the name doesn't match the volunteer ID.
// If an actual error occurs, the function returns nil and the error.
func (c *Client) SubmitHours(date time.Time, name string, vid uint, minutes map[Assignment]uint) (disposition map[Assignment]string, err error) {
	var (
		doc  *goquery.Document
		key  string
		warn error
	)

	// Find the volunteer.
	if key, err = c.findVolunteer(name, vid); key == "" {
		return nil, err
	} else {
		warn = err
	}

	// Read the page with that volunteer's assignments and service record.
	monthexpand := date.Format("200601") + "00000000"
	volPage := url.Values{}
	volPage.Add("ID", c.id)
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
	if doc, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Hours", volPage)); err != nil {
		return nil, err
	}

	// Handle each assignment type.
	disposition = make(map[Assignment]string)
	datefmt := date.Format("01-02-2006")
	rows := doc.Find("td.volgistics487").FilterFunction(func(_ int, node *goquery.Selection) bool {
		return node.Text() == datefmt
	}).Parent()
ASSN:
	for _, a := range allAssignments {
		var (
			found      bool
			updateForm = url.Values{}
		)
		updateForm.Add("ID", c.id)
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
			if cols.Eq(1).Text() != a.label() {
				continue
			}
			hnum := strings.Split(cols.Eq(3).Children().AttrOr("name", ":"), ":")[1]
			if minutes[a] != 0 {
				var current float64
				fmt.Sscanf(cols.Eq(2).Text(), "%f", &current)
				if hours := float64(minutes[a]) / 60; hours != current {
					updateForm.Add("HNUM", hnum)
					updateForm.Add("ODATE", datefmt)
					updateForm.Add("A0", strconv.Itoa(int(a)))
					updateForm.Add("H1FD", datefmt)
					updateForm.Add("MM0", fmt.Sprintf("%f", hours))
					updateForm.Add("Save", "Save")
					disposition[a] = "updated"
				} else {
					disposition[a] = "no change"
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
				disposition[a] = "deleted"
			}
			found = true
			break
		}
		if !found && minutes[a] == 0 {
			continue
		}
		if !found {
			updateForm.Add("A0", strconv.Itoa(int(a)))
			updateForm.Add("H1FD", datefmt)
			updateForm.Add("MM0", fmt.Sprintf("%f", float64(minutes[a])/60))
			updateForm.Add("Save", "Save")
			disposition[a] = "added"
		}
		delay()
		if _, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Hours", updateForm)); err != nil {
			return nil, err
		}
	}
	return disposition, warn
}

// findVolunteer returns the API key for the volunteer with the specified ID.
// If the volunteer is not found, it returns "" and an error.  If the
// volunteer's name in Volgistics does not match the supplied name,
func (c *Client) findVolunteer(name string, vid uint) (key string, err error) {
	var (
		doc   *goquery.Document
		row   *goquery.Selection
		vname string
	)
	findForm := url.Values{}
	findForm.Add("ID", c.id)
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
	findForm.Add("Fbt2", strconv.Itoa(int(vid)))
	findForm.Add("Find1", "Go")
	findForm.Add("Fbt3", "")
	findForm.Add("Fbt10", "")
	delay()
	if doc, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/core.dll/volunteers", findForm)); err != nil {
		return "", err
	}
	if row = doc.Find("#volTable tbody tr"); row.Length() != 1 {
		return "", fmt.Errorf("volunteer %s (%d) not found in Volgistics", name, vid)
	}
	if row = row.Children(); row.Length() < 3 {
		panic("can't parse Volgistics volunteer search results: wrong row length")
	}
	if key = row.First().Text(); key == "" {
		panic("can't parse Volgistics volunteer search results: no key")
	}
	vname = strings.SplitN(row.Eq(2).Text(), ",", 2)[0]
	if !strings.Contains(name, vname) {
		err = fmt.Errorf("volunteer %s (%d) name mismatch: Volgistics has %s", name, vid, row.Eq(2).Text())
	}
	return key, err
}

var zipCodeRE = regexp.MustCompile(` (\d{5}(?:-\d{4})?)$`)

const NewVolunteerFields = person.FAddresses | person.FBirthdate | person.FCellPhone | person.FEmContacts | person.FEmail | person.FHomePhone | person.FSortName | person.FWorkPhone

// NewVolunteer creates a new volunteer record and returns the new volunteer ID.
func (c *Client) NewVolunteer(p *person.Person) (vid int, err error) {
	var doc *goquery.Document
	var (
		lastName     string
		firstName    string
		addressLine1 string
		city         string
		stateCode    string
		zipCode      string
		email        string
		homePhone    string
		cellPhone    string
		workPhone    string
	)
	lastName, firstName, _ = strings.Cut(p.SortName(), ",")
	lastName = strings.TrimSpace(lastName)
	firstName = strings.TrimSpace(firstName)
	email = p.Email()
	if a := p.Addresses().Home; a != nil {
		var rest string
		var ok bool

		addressLine1, rest, _ = strings.Cut(a.Address, ",")
		addressLine1 = strings.TrimSpace(addressLine1)
		rest = strings.TrimSpace(rest)
		if match := zipCodeRE.FindStringSubmatch(rest); match != nil {
			zipCode = match[1]
			rest = strings.TrimSpace(strings.TrimSuffix(rest, zipCode))
		}
		if city, stateCode, ok = strings.Cut(rest, ","); !ok || len(strings.TrimSpace(stateCode)) != 2 {
			city, stateCode = rest, ""
		}
	}
	homePhone = formatPhone(p.HomePhone())
	cellPhone = formatPhone(p.CellPhone())
	workPhone = formatPhone(p.WorkPhone())
	newForm := url.Values{}
	newForm.Add("ID", c.id)
	newForm.Add("KEY", "0")
	newForm.Add("FB", "0")
	newForm.Add("TAGS", "")
	newForm.Add("Action", "")
	newForm.Add("zipDocs", "on")
	newForm.Add("Cg", "0")
	newForm.Add("F1", lastName)
	newForm.Add("F2", firstName)
	newForm.Add("F18", "Active")
	newForm.Add("F4", "")
	newForm.Add("G1", "0")
	newForm.Add("G2", "")
	newForm.Add("G3", "0")
	newForm.Add("G4", "0")
	newForm.Add("A1", addressLine1)
	newForm.Add("A2", "")
	newForm.Add("A3", "")
	newForm.Add("A4", city)
	newForm.Add("A5", stateCode)
	newForm.Add("A6", zipCode)
	newForm.Add("A19", email)
	newForm.Add("A7", homePhone)
	newForm.Add("A11", cellPhone) // "(xxx) xxx-xxxx"
	newForm.Add("A9", workPhone)
	newForm.Add("A13", "")
	newForm.Add("A15", "")
	newForm.Add("A17", "")
	newForm.Add("Save", "Save")
	if doc, err = checkResponse(httpClient.PostForm("https://www.volgistics.com/ex/core.dll/volunteers?TAB=Core", newForm)); err != nil {
		return 0, err
	}
	if input := doc.Find("input[name=F0]"); input.Length() != 1 {
		return 0, fmt.Errorf("new volunteer %s %s: ID not found in response", firstName, lastName)
	} else if vidstr, ok := input.Attr("value"); !ok {
		return 0, fmt.Errorf("new volunteer %s %s: ID not found in response", firstName, lastName)
	} else if vid, err = strconv.Atoi(vidstr); err != nil {
		return 0, fmt.Errorf("new volunteer %s %s: ID not parseable in response", firstName, lastName)
	}
	return vid, nil
}

func formatPhone(ph string) string {
	if len(ph) == 10 {
		return fmt.Sprintf("(%s) %s-%s", ph[0:3], ph[3:6], ph[6:10])
	} else {
		return ph
	}
}

var lastRequest time.Time

func delay() {
	time.Sleep(time.Until(lastRequest.Add(500 * time.Millisecond)))
	lastRequest = time.Now()
}

// checkResponse checks the response from an HTTP request to Volgistics.
func checkResponse(resp *http.Response, err error) (doc *goquery.Document, _ error) {
	if err != nil {
		return nil, fmt.Errorf("Volgistics: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("Volgistics: %s", resp.Status)
	}
	doc, err = goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Volgistics response: %w", err)
	}
	return doc, nil
}
