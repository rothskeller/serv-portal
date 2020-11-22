package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/authz"
	"sunnyvaleserv.org/portal/api/list"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadLists(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var l = new(model.List)
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
					fmt.Fprintf(os.Stderr, "ERROR: id must be first key in list\n")
					os.Exit(1)
				}
				l.ID = model.ListID(in.Int())
				if l.ID != 0 {
					lid := l.ID
					if l = tx.FetchList(l.ID); l == nil {
						fmt.Fprintf(os.Stderr, "ERROR: list %d does not exist\n", lid)
						os.Exit(1)
					}
					tx.WillUpdateList(l)
					*l = model.List{
						ID:     l.ID,
						People: make(map[model.PersonID]model.ListPersonStatus),
					}
				}
			case "type":
				str := in.String()
				for v, s := range model.ListTypeNames {
					if s == str {
						l.Type = v
						break
					}
				}
				if l.Type == 0 {
					in.AddError(errors.New("invalid type"))
				}
			case "name":
				l.Name = in.String()
			case "people":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var pid model.PersonID
						var lps model.ListPersonStatus

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
								pid = model.PersonID(in.Int())
							case "subscribed":
								if in.Bool() {
									lps |= model.ListSubscribed
								}
							case "unsubscribed":
								if in.Bool() {
									lps |= model.ListUnsubscribed
								}
							case "sender":
								if in.Bool() {
									lps |= model.ListSender
								}
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						if pid == 0 {
							in.AddError(errors.New("missing people.id"))
						} else {
							l.People[pid] = lps
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
			fmt.Fprintf(os.Stderr, "ERROR: id must be first key in list\n")
			os.Exit(1)
		}
		if err := list.ValidateList(tx, l); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if l.ID == 0 {
			tx.CreateList(l)
		} else {
			tx.UpdateList(l)
		}
		record++
	}

}
