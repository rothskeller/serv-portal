package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mailru/easyjson/jwriter"

	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/store/enum"
	"sunnyvaleserv.org/portal/store/list"
	"sunnyvaleserv.org/portal/store/listperson"
	"sunnyvaleserv.org/portal/store/person"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	var (
		entry  *log.Entry
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
	entry = log.New("", "synch-email")
	store.Connect(context.Background(), entry, func(st *store.Store) {
		out.RawByte('[')
		person.All(st, person.FID|person.FInformalName|person.FEmail|person.FEmail2|person.FFlags|person.FPrivLevels|person.FUnsubscribeToken, func(p *person.Person) {
			var sender, receiver []string

			if p.Flags()&person.NoEmail != 0 || (p.Email() == "" && p.Email2() == "") || !p.HasPrivLevel(0, enum.PrivStudent) {
				return
			}
			listperson.SubscriptionsByPerson(st, p.ID(), func(l *list.List) {
				if l.Type == list.Email {
					receiver = append(receiver, l.Name+"@sunnyvaleserv.org")
				}
			})
			listperson.SendersByPerson(st, p.ID(), func(l *list.List) {
				if l.Type == list.Email {
					sender = append(sender, l.Name+"@sunnyvaleserv.org")
				}
			})
			if len(sender) == 0 && len(receiver) == 0 {
				return
			}
			for _, email := range []string{p.Email(), p.Email2()} {
				if email == "" {
					continue
				}
				if first {
					first = false
				} else {
					out.RawByte(',')
				}
				out.RawString(`{"id":`)
				out.IntStr(int(p.ID()))
				out.RawString(`,"name":`)
				out.String(p.InformalName())
				out.RawString(`,"email":`)
				out.String(email)
				out.RawString(`,"token":`)
				out.String(p.UnsubscribeToken())
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
		})
		out.RawByte(']')
	})
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
