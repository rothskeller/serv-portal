package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"regexp"
	"strings"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/config"
	"sunnyvaleserv.org/portal/util/sendmail"
)

var addrRE = regexp.MustCompile(`([^"]+)_unsub`)
var letterRE = regexp.MustCompile(`(?s)\?letter=."><.*?\?letter=(.)`)
var csrfRE = regexp.MustCompile(`name="csrf_token" value="([^"]*)"`)

func main() {
	var (
		tx      *store.Tx
		err     error
		mail    bytes.Buffer
		changes bool
	)
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err := os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err := os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	store.Open("serv.db")
	tx = store.Begin(nil)
	fmt.Fprintf(&mail, "From: %s\r\nTo: <%s>\r\nSubject: SERV Mailing List Updates\r\n\r\n", config.Get("fromEmail"), config.Get("adminEmail"))
	for _, group := range tx.Authorizer().FetchGroups(tx.Authorizer().AllGroups()) {
		if group.Email == "" {
			continue
		}
		if updateMailingList(tx, &mail, group) {
			changes = true
		}
	}
	if !changes {
		return
	}
	err = sendmail.SendMessage(config.Get("fromAddr"), []string{config.Get("adminEmail")}, mail.Bytes())
	if err != nil {
		log.Fatalf("send mail: %s", err)
	}
}

func updateMailingList(tx *store.Tx, mail *bytes.Buffer, group *model.Group) (changed bool) {
	var auth = tx.Authorizer()
	var (
		disabled    model.GroupID
		listaddr    string
		client      http.Client
		resp        *http.Response
		body        []byte
		letter      string
		elist       strings.Builder
		csrf        string
		err         error
		desired     = map[string]bool{}
		actual      = map[string]bool{}
		noEmailPIDs = map[model.PersonID]bool{}
	)
	disabled = auth.FetchGroupByTag(model.GroupDisabled).ID
	for _, pid := range group.NoEmail {
		noEmailPIDs[pid] = true
	}
	for _, person := range tx.FetchPeople() {
		if auth.MemberPG(person.ID, disabled) || noEmailPIDs[person.ID] || person.NoEmail {
			continue
		}
		if auth.MemberPG(person.ID, group.ID) || auth.CanPAG(person.ID, model.PrivBCC, group.ID) {
			if person.Email != "" {
				desired[person.Email] = true
			}
			if person.Email2 != "" {
				desired[person.Email2] = true
			}
		}
	}
	if len(desired) == 0 {
		fmt.Fprintf(mail, "WARNING: no addresses found for list %s; skipping\n\n", group.Email)
		return true
	}
	client.Jar, _ = cookiejar.New(nil)
	listaddr = fmt.Sprintf("http://lists.sunnyvaleserv.org/admin.cgi/%s-sunnyvaleserv.org", group.Email)
	resp, err = client.PostForm(listaddr, url.Values{"adminpw": {config.Get("mailmainAdminPassword")}, "admlogin": {"Let me in..."}})
	if err != nil {
		fmt.Fprintf(mail, "ERROR: can't log in to %s admin console: %s\n\n", group.Email, err)
		return true
	}
	resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Fprintf(mail, "ERROR: can't log in to %s admin console: %s\n\n", group.Email, resp.Status)
		return true
	}
	for {
		resp, err = client.Get(listaddr + "/members" + letter)
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't get members of %s: %s\n\n", group.Email, err)
			return true
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Fprintf(mail, "ERROR: can't get members of %s: %s\n\n", group.Email, resp.Status)
			return true
		}
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't read members of %s: %s\n\n", group.Email, err)
			return true
		}
		for _, m := range addrRE.FindAllSubmatch(body, -1) {
			actual[strings.ToLower(strings.Replace(string(m[1]), "%40", "@", -1))] = true
		}
		if l := letterRE.FindSubmatch(body); l != nil {
			letter = "?letter=" + string(l[1])
		} else {
			break
		}
	}
	if len(actual) == 0 && len(os.Args) == 1 {
		fmt.Fprintf(mail, "WARNING: no addresses found on list %s; probable web scraping failure; skipping; override with any argument on synch-email command line\n\n",
			group.Email)
		return true
	}
	for a := range actual {
		if !desired[a] {
			elist.WriteString(a)
			elist.WriteByte('\n')
			fmt.Fprintf(mail, "%s: remove %s\n", group.Email, a)
		}
	}
	if elist.Len() != 0 {
		resp, err = client.Get(listaddr + "/members/remove")
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't get member remove page for %s: %s\n\n", group.Email, err)
			return true
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Fprintf(mail, "ERROR: can't get member remove page for %s: %s\n\n", group.Email, resp.Status)
			return true
		}
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't read member remove page for %s: %s\n\n", group.Email, err)
			return true
		}
		if c := csrfRE.FindSubmatch(body); c != nil {
			csrf = string(c[1])
		} else {
			fmt.Fprintf(mail, "ERROR: no CSRF token on member remove page for %s\n\n", group.Email)
			return true
		}
		resp, err = client.PostForm(listaddr+"/member/remove", url.Values{
			"setmemberopts_btn":                      {"Submit Your Changes"},
			"send_unsub_ack_to_this_batch":           {"0"},
			"send_unsub_notifications_to_list_owner": {"0"},
			"unsubscribees":                          {elist.String()},
			"csrf_token":                             {csrf},
		})
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't remove members from %s: %s\n\n", group.Email, err)
			return true
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(mail, "ERROR: can't remove members from %s: %s\n\n", group.Email, resp.Status)
			return true
		}
		changed = true
		elist.Reset()
	}
	for d := range desired {
		if !actual[d] {
			elist.WriteString(d)
			elist.WriteByte('\n')
			fmt.Fprintf(mail, "%s: add %s\n", group.Email, d)
		}
	}
	if elist.Len() != 0 {
		resp, err = client.Get(listaddr + "/members/add")
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't get member add page for %s: %s\n\n", group.Email, err)
			return true
		}
		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Fprintf(mail, "ERROR: can't get member add page for %s: %s\n\n", group.Email, resp.Status)
			return true
		}
		body, err = ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't read member add page for %s: %s\n\n", group.Email, err)
			return true
		}
		if c := csrfRE.FindSubmatch(body); c != nil {
			csrf = string(c[1])
		} else {
			fmt.Fprintf(mail, "ERROR: no CSRF token on member add page for %s\n\n", group.Email)
			return true
		}
		resp, err = client.PostForm(listaddr+"/member/add", url.Values{
			"setmemberopts_btn":                {"Submit Your Changes"},
			"subscribe_or_invite":              {"0"},
			"send_welcome_msg_to_this_batch":   {"0"},
			"send_notifications_to_list_owner": {"0"},
			"subscribees":                      {elist.String()},
			"csrf_token":                       {csrf},
		})
		if err != nil {
			fmt.Fprintf(mail, "ERROR: can't add members from %s: %s\n\n", group.Email, err)
			return true
		}
		resp.Body.Close()
		if resp.StatusCode != http.StatusOK {
			fmt.Fprintf(mail, "ERROR: can't add members from %s: %s\n\n", group.Email, resp.Status)
			return true
		}
		changed = true
	}
	if changed {
		mail.WriteByte('\n')
	}
	return changed
}
