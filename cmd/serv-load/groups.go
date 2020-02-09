package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/group"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadGroups(tx *store.Tx, in *jlexer.Lexer) {
	auth := tx.Authorizer()
	var record = 1
	for {
		var g = new(model.Group)
		var first = true

		in.Delim('{')
		if in.Error() == io.EOF {
			auth.Save()
			return
		}
		for !in.IsDelim('}') {
			key := in.UnsafeString()
			in.WantColon()
			if in.IsNull() {
				in.Skip()
				in.WantComma()
				continue
			}
			switch key {
			case "id":
				if !first {
					fmt.Fprintf(os.Stderr, "ERROR: id must be first key in group\n")
					os.Exit(1)
				}
				g.ID = model.GroupID(in.Int())
				if g.ID == 0 {
					g = auth.CreateGroup()
				} else {
					gid := g.ID
					if g = auth.FetchGroup(g.ID); g == nil {
						fmt.Fprintf(os.Stderr, "ERROR: group %d does not exist\n", gid)
						os.Exit(1)
					}
					auth.WillUpdateGroup(g)
					g.Name = ""
					g.Tag = ""
					g.NoEmail = nil
					g.NoText = nil
				}
			case "tag":
				g.Tag = model.GroupTag(in.String())
			case "name":
				g.Name = in.String()
			case "email":
				g.Email = in.String()
			case "noEmail":
				in.Delim('[')
				for !in.IsDelim(']') {
					pid := model.PersonID(in.Int())
					if tx.FetchPerson(pid) == nil {
						fmt.Fprintf(os.Stderr, "ERROR: person %d does not exist\n", pid)
						os.Exit(1)
					}
					g.NoEmail = append(g.NoEmail, pid)
					in.WantComma()
				}
				in.Delim(']')
			case "noText":
				in.Delim('[')
				for !in.IsDelim(']') {
					pid := model.PersonID(in.Int())
					if tx.FetchPerson(pid) == nil {
						fmt.Fprintf(os.Stderr, "ERROR: person %d does not exist\n", pid)
						os.Exit(1)
					}
					g.NoText = append(g.NoEmail, pid)
					in.WantComma()
				}
				in.Delim(']')
			default:
				in.SkipRecursive()
			}
			in.WantComma()
			first = false
		}
		in.Delim('}')
		if !in.Ok() {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, in.Error())
			os.Exit(1)
		}
		if first {
			fmt.Fprintf(os.Stderr, "ERROR: id must be first key in group\n")
			os.Exit(1)
		}
		if err := group.ValidateGroup(auth, g); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		auth.UpdateGroup(g)
		record++
	}
}
