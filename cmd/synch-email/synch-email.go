package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func main() {
	var (
		tx     *store.Tx
		lists  []*model.List
		out    jwriter.Writer
		tempfn string
		permfn string
		tempfh *os.File
		err    error
		first  = true
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
	for _, l := range tx.FetchLists() {
		if l.Type == model.ListEmail {
			l.Name += "@sunnyvaleserv.org"
			lists = append(lists, l)
		}
	}
	out.RawByte('[')
	for _, p := range tx.FetchPeople() {
		var sender, receiver []string

		if p.NoEmail || (p.Email == "" && p.Email2 == "") || !p.HasPrivLevel(model.PrivStudent) {
			continue
		}
		for _, l := range lists {
			if l.People[p.ID]&model.ListSender != 0 {
				sender = append(sender, l.Name)
			}
			if l.People[p.ID]&model.ListSubscribed != 0 {
				receiver = append(receiver, l.Name)
			}
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
