package person

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/mailru/easyjson/jwriter"
	"sunnyvaleserv.org/portal/model"
)

var addressSplitRE = regexp.MustCompile(`^(.*),\s*([^,]+),?\s+CA\s+(\d{5})$`)

func SendVolunteerRegistration(p *model.Person, interests []string) (err error) {
	var out jwriter.Writer

	parts := strings.SplitN(p.SortName, ", ", 2)
	out.RawString(`{"g1f2":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g1f1":`)
	out.String(parts[0])
	out.RawString(`,"g1f5":0,"g1f6":`)
	match := addressSplitRE.FindStringSubmatch(p.HomeAddress.Address)
	if match == nil {
		return errors.New("can't parse home address")
	}
	out.String(match[1])
	out.RawString(`,"g1f8":`)
	out.String(match[2])
	out.RawString(`,"g1f9":104,"g1f10":`)
	out.String(match[3])
	if p.HomePhone != "" {
		out.RawString(`,"g1f11":`)
		out.String(volgisticsPhoneFormat(p.HomePhone))
		if p.CellPhone != "" {
			out.RawString(`,"g1f13":`)
			out.String(volgisticsPhoneFormat(p.CellPhone))
		}
	} else {
		out.RawString(`,"g1f11":`)
		out.String(volgisticsPhoneFormat(p.CellPhone))
	}
	out.RawString(`,"g1f22":false,"g1f23":false,"g1f17":`)
	out.String(p.Email)
	out.RawString(`,"g2f99":`)
	if len(interests) != 0 {
		out.String("ATTN: Steve Roth\nSunnyvale DPS/OES\nInterested in " + strings.Join(interests, ", "))
	} else {
		out.String("ATTN: Steve Roth\nSunnyvale DPS/OES")
	}
	out.RawString(`,"g8f1":`)
	out.String(p.Birthdate)
	out.RawString(`,"g8f4":0,"g8f5":0,"g8f6":0,"g8f7":0,"g8f21":"None","g9f3c251":false,"g9f3c252":false,"g9f3c253":false,"g9f3c254":false,"g9f3c255":false,"g9f3c256":false,"g9f3c257":false,"g9f3c258":false,"g9f1c213":false,"g9f1c214":false,"g9f1c215":false,"g9f1c233":false,"g9f1c223":false,"g9f1c216":false,"g9f1c339":false,"g9f1c217":false,"g9f1c218":false,"g9f1c219":false,"g9f1c220":false,"g9f1c221":false,"g9f1c222":false,"g9f1c302":false,"g9f1c224":false,"g9f1c225":false,"g9f1c226":false,"g9f1c227":false,"g9f1c232":false,"g9f1c228":false,"g9f1c229":false,"g9f1c230":false,"g9f1c231":false,"g9f1c334":false,"g9f2c234":false,"g9f2c235":false,"g9f2c236":false,"g9f2c237":false,"g9f2c238":false,"g9f2c239":false,"g9f2c240":false,"g9f2c241":false,"g9f2c242":false,"g9f2c243":false,"g9f2c244":false,"g9f2c245":false,"g9f2c246":false,"g9f2c247":false,"g9f2c248":false,"g9f2c249":false,"g9f2c250":false,"g9f22c93":false,"g9f22c85":false,"g9f22c143":false,"g9f22c129":false,"g9f22c75":false,"g9f22c18":false,"g9f22c73":false,"g9f22c16":false,"g9f22c77":false,"g9f22c79":false,"g9f22c81":false,"g9f22c41":false,"g9f22c83":false,"g9f22c91":false,"g9f22c1":false,"g2f2c1l1":`)
	parts = strings.Fields(p.EmContacts[0].Name)
	out.String(strings.Join(parts[:len(parts)-1], " "))
	out.RawString(`,"g2f1c1l1":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g2f0c1l1":0`)
	if p.EmContacts[0].HomePhone != "" {
		out.RawString(`,"g2f11c1l1":`)
		out.String(volgisticsPhoneFormat(p.EmContacts[0].HomePhone))
		if p.EmContacts[0].CellPhone != "" {
			out.RawString(`,"g2f13c1l1":`)
			out.String(volgisticsPhoneFormat(p.EmContacts[0].CellPhone))
		}
	} else {
		out.RawString(`,"g2f11c1l1":`)
		out.String(volgisticsPhoneFormat(p.EmContacts[0].CellPhone))
	}
	out.RawString(`,"g2f25c1l1":`)
	out.Int(relationshipCode[p.EmContacts[0].Relationship])
	out.RawString(`,"g2f2c1l2":`)
	parts = strings.Fields(p.EmContacts[1].Name)
	out.String(strings.Join(parts[:len(parts)-1], " "))
	out.RawString(`,"g2f1c1l2":`)
	out.String(parts[len(parts)-1])
	out.RawString(`,"g2f0c1l2":0`)
	if p.EmContacts[1].HomePhone != "" {
		out.RawString(`,"g2f11c1l2":`)
		out.String(volgisticsPhoneFormat(p.EmContacts[1].HomePhone))
		if p.EmContacts[1].CellPhone != "" {
			out.RawString(`,"g2f13c1l2":`)
			out.String(volgisticsPhoneFormat(p.EmContacts[1].CellPhone))
		}
	} else {
		out.RawString(`,"g2f11c1l2":`)
		out.String(volgisticsPhoneFormat(p.EmContacts[1].CellPhone))
	}
	out.RawString(`,"g2f25c1l2":`)
	out.Int(relationshipCode[p.EmContacts[1].Relationship])
	out.RawString(`,"g13f1":true,"g16f3c0l1":false,"g16f3c0l2":false,"g16f3c0l3":false,"g16f3c1l1":false,"g16f3c1l2":false,"g16f3c1l3":false,"g16f3c2l1":false,"g16f3c2l2":false,"g16f3c2l3":false,"g16f3c3l1":false,"g16f3c3l2":false,"g16f3c3l3":false,"g16f3c4l1":false,"g16f3c4l2":false,"g16f3c4l3":false,"g16f3c5l1":false,"g16f3c5l2":false,"g16f3c5l3":false,"g16f3c6l1":false,"g16f3c6l2":false,"g16f3c6l3":false,"appSignature":[],"id":"777028848"}`)
	var body, _ = out.BuildBytes()

	var req, _ = http.NewRequest(http.MethodPost, "https://www.volgistics.com/ex/apForm.dll/AFR?Action=saveApplication", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json, text/plain, */*")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("volgistics registration failed with status %d", resp.StatusCode)
	}
	return nil
}

func volgisticsPhoneFormat(s string) string {
	return fmt.Sprintf("(%s) %s", s[0:3], s[4:])
}

var relationshipCode = map[string]int{
	"Co-worker":  125,
	"Daughter":   313,
	"Father":     120,
	"Friend":     123,
	"Mother":     119,
	"Neighbor":   126,
	"Other":      314,
	"Relative":   275,
	"Son":        121,
	"Spouse":     331,
	"Supervisor": 124,
}
