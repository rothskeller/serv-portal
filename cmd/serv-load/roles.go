package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/role"
)

func loadRoles(tx *db.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var r = new(model.Role)
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
				r.ID = model.RoleID(in.Int())
				if r.ID != 0 {
					if !first {
						fmt.Fprintf(os.Stderr, "ERROR: id must be first key in role\n")
						os.Exit(1)
					}
					rid := r.ID
					if r = tx.FetchRole(r.ID); r == nil {
						fmt.Fprintf(os.Stderr, "ERROR: role %d does not exist\n", rid)
						os.Exit(1)
					}
					r.Name = ""
					r.Individual = false
					r.Privileges.Clear()
				}
			case "name":
				r.Name = in.String()
			case "individual":
				r.Individual = in.Bool()
			case "privileges":
				if in.IsNull() {
					in.Skip()
				} else {
					in.Delim('[')
					for !in.IsDelim(']') {
						var group *model.Group
						var priv model.Privilege
						in.Delim('{')
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
								if group = tx.FetchGroup(model.GroupID(in.Int())); group == nil {
									in.AddError(errors.New("invalid group"))
								}
							case "privilege":
								priv.UnmarshalEasyJSON(in)
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						in.WantComma()
						if group == nil {
							in.AddError(errors.New("missing group in privilege"))
						} else if priv == 0 {
							in.AddError(errors.New("missing privilege in privilege"))
						} else {
							r.Privileges.Add(group, priv)
						}
					}
					in.Delim(']')
				}
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
		if err := role.ValidateRole(tx, r); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if r.ID == 0 {
			tx.CreateRole(r)
		}
		record++
	}
}
