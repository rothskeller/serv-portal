package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/api/role"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadRoles(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	var err error
	for {
		var r = &model.Role{
			Implies: make(map[model.RoleID]bool),
			Lists:   make(map[model.ListID]model.RoleToList),
		}
		var first = true

		in.Delim('{')
		if in.Error() == io.EOF {
			authz.UpdateAuthz(tx)
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
				if r.ID != 0 {
					rid := r.ID
					if r = tx.FetchRole(r.ID); r == nil {
						fmt.Fprintf(os.Stderr, "ERROR: role %d does not exist\n", rid)
						os.Exit(1)
					}
					tx.WillUpdateRole(r)
					*r = model.Role{
						ID:      r.ID,
						Implies: make(map[model.RoleID]bool),
						Lists:   make(map[model.ListID]model.RoleToList),
					}
				}
			case "name":
				r.Name = in.String()
			case "title":
				r.Title = in.String()
			case "org":
				r.Org, err = model.ParseOrg(in.String())
				in.AddError(err)
			case "privLevel":
				r.PrivLevel, err = model.ParsePrivLevel(in.String())
				in.AddError(err)
			case "showRoster":
				r.ShowRoster = in.Bool()
			case "implicitOnly":
				r.ImplicitOnly = in.Bool()
			case "priority":
				r.Priority = in.Int()
			case "implies":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var rid model.RoleID
						var direct bool

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
								rid = model.RoleID(in.Int())
							case "direct":
								direct = in.Bool()
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						if rid == 0 {
							in.AddError(errors.New("missing implies.id"))
						} else if direct {
							r.Implies[rid] = true
						}
					}
					in.WantComma()
				}
				in.Delim(']')
			case "lists":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var lid model.ListID
						var rtl model.RoleToList

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
								lid = model.ListID(in.Int())
							case "subModel":
								str := in.String()
								for v, s := range model.ListSubModelNames {
									if s == str {
										rtl.SetSubModel(v)
										break
									}
								}
								if rtl.SubModel() == model.ListNoSub {
									in.AddError(errors.New("invalid subModel"))
								}
							case "sender":
								rtl.SetSender(in.Bool())
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						if lid == 0 {
							in.AddError(errors.New("missing lists.id"))
						} else {
							r.Lists[lid] = rtl
						}
					}
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
			fmt.Fprintf(os.Stderr, "ERROR: id must be first key in role\n")
			os.Exit(1)
		}
		if err := role.ValidateRole(tx, r); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if r.ID == 0 {
			tx.CreateRole(r)
		} else {
			tx.UpdateRole(r)
		}
		record++
	}

}
