package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/authz"
)

func main() {
	var (
		tx       *store.Tx
		auth     *authz.Authorizer
		groups   []*model.Group
		out      jwriter.Writer
		disabled model.GroupID
		tempfn   string
		permfn   string
		tempfh   *os.File
		err      error
		first    = true
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
	auth = tx.Authorizer()
	for _, g := range auth.FetchGroups(auth.AllGroups()) {
		if g.Email != "" {
			g.Email += "@sunnyvaleserv.org"
			groups = append(groups, g)
		}
		if g.Tag == model.GroupDisabled {
			disabled = g.ID
		}
	}
	out.RawByte('[')
	for _, p := range tx.FetchPeople() {
		var sender, receiver []string

		if auth.MemberPG(p.ID, disabled) || p.NoEmail || (p.Email == "" && p.Email2 == "") {
			continue
		}
	GROUPS:
		for _, g := range groups {
			if auth.CanPAG(p.ID, model.PrivSendEmailMessages, g.ID) {
				sender = append(sender, g.Email)
			}
			if !auth.MemberPG(p.ID, g.ID) && !auth.CanPAG(p.ID, model.PrivBCC, g.ID) {
				continue
			}
			for _, pid := range g.NoEmail {
				if pid == p.ID {
					continue GROUPS
				}
			}
			receiver = append(receiver, g.Email)
		}
		if len(sender) == 0 && len(receiver) == 0 {
			continue
		}
		for _, email := range []string{p.Email, p.Email2} {
			if email == "" {
				continue
			}
			if first {
				first = false
			} else {
				out.RawByte(',')
			}
			out.RawString(`{"id":`)
			out.IntStr(int(p.ID))
			out.RawString(`,"name":`)
			out.String(p.InformalName)
			out.RawString(`,"email":`)
			out.String(email)
			out.RawString(`,"token":`)
			out.String(p.UnsubscribeToken)
			if len(sender) != 0 {
				out.RawString(`,"sender":[`)
				for i, s := range sender {
					if i != 0 {
						out.RawByte(',')
					}
					out.String(s)
				}
				out.RawByte(']')
			}
			if len(receiver) != 0 {
				out.RawString(`,"receiver":[`)
				for i, r := range receiver {
					if i != 0 {
						out.RawByte(',')
					}
					out.String(r)
				}
				out.RawByte(']')
			}
			out.RawByte('}')
		}
	}
	out.RawByte(']')
	permfn = filepath.Join(os.Getenv("HOME"), "maillist", "lists")
	tempfn = permfn + ".temp"
	if tempfh, err = os.Create(tempfn); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
		os.Exit(1)
	}
	if _, err = out.DumpTo(tempfh); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: write %s: %s\n", tempfn, err)
		os.Exit(1)
	}
	if err = tempfh.Close(); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: close %s: %s\n", tempfn, err)
		os.Exit(1)
	}
	if err = os.Rename(tempfn, permfn); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: move %s to %s: %s\n", tempfn, permfn, err)
		os.Exit(1)
	}
}
