package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/role"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadRoles(tx *store.Tx, in *jlexer.Lexer) {
	auth := tx.Authorizer()
	var record = 1
	for {
		var r = new(model.Role)
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
					fmt.Fprintf(os.Stderr, "ERROR: id must be first key in role\n")
					os.Exit(1)
				}
				r.ID = model.RoleID(in.Int())
				if r.ID == 0 {
					r = auth.CreateRole()
				} else {
					rid := r.ID
					if r = auth.FetchRole(r.ID); r == nil {
						fmt.Fprintf(os.Stderr, "ERROR: role %d does not exist\n", rid)
						os.Exit(1)
					}
					auth.WillUpdateRole(r)
					r.Name = ""
					r.Individual = false
				}
			case "name":
				r.Name = in.String()
			case "individual":
				r.Individual = in.Bool()
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
			fmt.Fprintf(os.Stderr, "ERROR: id must be first key in role\n")
			os.Exit(1)
		}
		if err := role.ValidateRole(auth, r); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		auth.UpdateRole(r)
		record++
	}
}
