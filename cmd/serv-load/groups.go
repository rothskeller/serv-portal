package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/group"
	"sunnyvaleserv.org/portal/model"
)

func loadGroups(tx *db.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var g = new(model.Group)
		var first = true

		in.Delim('{')
		if in.Error() == io.EOF {
			tx.SaveAuthz()
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
				g.ID = model.GroupID(in.Int())
				if g.ID != 0 {
					if !first {
						fmt.Fprintf(os.Stderr, "ERROR: id must be first key in group\n")
						os.Exit(1)
					}
					rid := g.ID
					if g = tx.FetchGroup(g.ID); g == nil {
						fmt.Fprintf(os.Stderr, "ERROR: group %d does not exist\n", rid)
						os.Exit(1)
					}
					g.Name = ""
					g.Tag = ""
					g.AllowTextMessages = false
				}
			case "tag":
				g.Tag = model.GroupTag(in.String())
			case "name":
				g.Name = in.String()
			case "allowTextMessages":
				g.AllowTextMessages = in.Bool()
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
		if err := group.ValidateGroup(tx, g); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if g.ID == 0 {
			tx.CreateGroup(g)
		}
		record++
	}
}
