package main

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/person"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadPeople(tx *store.Tx, in *jlexer.Lexer) {
	var err error
	auth := tx.Authorizer()
	var record = 1
	for {
		var p = new(model.Person)
		var roles []model.RoleID
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
				p.ID = model.PersonID(in.Int())
				if p.ID != 0 {
					if !first {
						fmt.Fprintf(os.Stderr, "ERROR: id must be first key in person\n")
						os.Exit(1)
					}
					pid := p.ID
					if p = tx.FetchPerson(p.ID); p == nil {
						fmt.Fprintf(os.Stderr, "ERROR: group %d does not exist\n", pid)
						os.Exit(1)
					}
					tx.WillUpdatePerson(p)
					*p = model.Person{ID: pid}
				}
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
			case "noEmail":
				p.NoEmail = in.Bool()
			case "noText":
				p.NoText = in.Bool()
			case "unsubscribeToken":
				p.UnsubscribeToken = in.String()
			case "roles":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						if in.IsDelim('{') {
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
									roles = append(roles, model.RoleID(in.Int()))
								default:
									in.SkipRecursive()
								}
								in.WantComma()
							}
							in.Delim('}')
						} else {
							roles = append(roles, model.RoleID(in.Int()))
						}
					}
					in.WantComma()
				}
				in.Delim(']')
			case "dswForms":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var form model.DSWForm
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
							case "from":
								form.From, err = time.ParseInLocation("2006-01-02", in.String(), time.Local)
								in.AddError(err)
							case "to":
								form.To, err = time.ParseInLocation("2006-01-02", in.String(), time.Local)
								in.AddError(err)
							case "for":
								form.For = in.String()
							case "invalid":
								form.Invalid = in.String()
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						p.DSWForms = append(p.DSWForms, &form)
					}
					in.WantComma()
				}
				in.Delim(']')
			case "volgisticsID":
				p.VolgisticsID = in.Int()
			case "backgroundCheck":
				p.BackgroundCheck = in.String()
			case "hoursToken":
				p.HoursToken = in.String()
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
		if err := person.ValidatePerson(tx, p, roles); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if p.ID == 0 {
			tx.CreatePerson(p)
		} else {
			tx.UpdatePerson(p)
		}
		tx.Authorizer().SetPersonRoles(p.ID, roles)
		record++
	}
}
