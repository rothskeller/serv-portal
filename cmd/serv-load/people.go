package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/db"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/person"
)

func loadPeople(tx *db.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var p model.Person

		in.Delim('{')
		if in.Error() == io.EOF {
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
				p.ID = model.PersonID(in.Int())
			case "username":
				p.Username = in.String()
			case "informalName":
				p.InformalName = in.String()
			case "formalName":
				p.FormalName = in.String()
			case "sortName":
				p.SortName = in.String()
			case "callSign":
				p.CallSign = in.String()
			case "email":
				p.Email = in.String()
			case "email2":
				p.Email2 = in.String()
			case "homeAddress":
				p.HomeAddress.UnmarshalEasyJSON(in)
			case "workAddress":
				p.WorkAddress.UnmarshalEasyJSON(in)
			case "mailAddress":
				p.MailAddress.UnmarshalEasyJSON(in)
			case "cellPhone":
				p.CellPhone = in.String()
			case "homePhone":
				p.HomePhone = in.String()
			case "workPhone":
				p.WorkPhone = in.String()
			case "password":
				if in.IsNull() {
					in.Skip()
					p.Password = nil
				} else {
					p.Password = []byte(in.String())
				}
			case "badLoginCount":
				p.BadLoginCount = in.Int()
			case "badLoginTime":
				if data := in.Raw(); in.Ok() {
					in.AddError(p.BadLoginTime.UnmarshalJSON(data))
				}
			case "pwresetToken":
				p.PWResetToken = in.String()
			case "pwresetTime":
				if data := in.Raw(); in.Ok() {
					in.AddError(p.PWResetTime.UnmarshalJSON(data))
				}
			case "roles":
				if in.IsNull() {
					in.Skip()
					p.Roles = nil
				} else {
					in.Delim('[')
					if p.Roles == nil {
						if !in.IsDelim(']') {
							p.Roles = make([]model.RoleID, 0, 8)
						} else {
							p.Roles = []model.RoleID{}
						}
					} else {
						p.Roles = (p.Roles)[:0]
					}
					for !in.IsDelim(']') {
						if in.IsDelim(('{')) {
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
									p.Roles = append(p.Roles, model.RoleID(in.Int()))
								default:
									in.SkipRecursive()
								}
								in.WantComma()
							}
							in.Delim('}')
						} else {
							p.Roles = append(p.Roles, model.RoleID(in.Int()))
						}
						in.WantComma()
					}
					in.Delim(']')
				}
			case "notes":
				if in.IsNull() {
					in.Skip()
					p.Notes = nil
				} else {
					in.Delim('[')
					if p.Notes == nil {
						if !in.IsDelim(']') {
							p.Notes = make([]*model.PersonNote, 0, 4)
						} else {
							p.Notes = []*model.PersonNote{}
						}
					} else {
						p.Notes = p.Notes[:0]
					}
					for !in.IsDelim(']') {
						if in.IsNull() {
							in.Skip()
						} else {
							var pn model.PersonNote
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
								case "note":
									pn.Note = in.String()
								case "date":
									pn.Date = in.String()
								case "privilege":
									pn.Privilege.UnmarshalEasyJSON(in)
								default:
									in.SkipRecursive()
								}
								in.WantComma()
							}
							in.Delim('}')
							p.Notes = append(p.Notes, &pn)
						}
						in.WantComma()
					}
					in.Delim(']')
				}
			default:
				in.SkipRecursive()
			}
			in.WantComma()
		}
		in.Delim('}')
		if !in.Ok() {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, in.Error())
			os.Exit(1)
		}
		if err := person.ValidatePerson(tx, &p); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		tx.SavePerson(&p)
		record++
	}
}
