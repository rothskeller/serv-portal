package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/event"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadEvents(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	var err error
	for {
		var e = new(model.Event)
		var ea = map[model.PersonID]model.AttendanceInfo{}
		var first = true

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
				e.ID = model.EventID(in.Int())
				if e.ID != 0 {
					if !first {
						fmt.Fprintf(os.Stderr, "ERROR: id must be first key in event\n")
						os.Exit(1)
					}
					eid := e.ID
					if e = tx.FetchEvent(e.ID); e == nil {
						fmt.Fprintf(os.Stderr, "ERROR: event %d does not exist\n", eid)
						os.Exit(1)
					}
					*e = model.Event{ID: eid}
				}
			case "name":
				e.Name = in.String()
			case "date":
				e.Date = in.String()
			case "start":
				e.Start = in.String()
			case "end":
				e.End = in.String()
			case "venue":
				if in.IsDelim('{') {
					var seen bool
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
							e.Venue = model.VenueID(in.Int())
							seen = true
						default:
							in.SkipRecursive()
						}
						in.WantComma()
					}
					in.Delim('}')
					if !seen {
						in.AddError(errors.New("missing venue.id"))
					}
				} else {
					e.Venue = model.VenueID(in.Int())
				}
			case "details":
				e.Details = in.String()
			case "type":
				etype := in.String()
				for _, et := range model.AllEventTypes {
					if etype == model.EventTypeNames[et] {
						e.Type = et
					}
				}
				if e.Type == 0 {
					in.AddError(errors.New("invalid type"))
				}
			case "renewsDSW":
				e.RenewsDSW = in.Bool()
			case "coveredByDSW":
				e.CoveredByDSW = in.Bool()
			case "org":
				e.Org, err = model.ParseOrg(in.String())
				in.AddError(err)
			case "roles":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var rid model.RoleID
						if in.IsDelim('{') {
							var seen bool
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
									seen = true
								default:
									in.SkipRecursive()
								}
								in.WantComma()
							}
							in.Delim('}')
							if !seen {
								in.AddError(errors.New("missing roles.id"))
							}
						} else {
							rid = model.RoleID(in.Int())
						}
						e.Roles = append(e.Roles, rid)
					}
					in.WantComma()
				}
				in.Delim(']')
			case "signupText":
				e.SignupText = in.String()
			case "shifts":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var s model.Shift
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
							case "start":
								s.Start = in.String()
							case "end":
								s.End = in.String()
							case "task":
								s.Task = in.String()
							case "min":
								s.Min = in.Int()
							case "max":
								s.Max = in.Int()
							case "announce":
								s.Announce = in.Bool()
							case "signedUp":
								in.Delim('[')
								for !in.IsDelim(']') {
									if in.IsNull() {
										in.Skip()
									} else {
										var pid model.PersonID
										if in.IsDelim('{') {
											var seen bool
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
													seen = true
												default:
													in.SkipRecursive()
												}
												in.WantComma()
											}
											in.Delim('}')
											if !seen {
												in.AddError(errors.New("missing assignments.shifts.signedUp.id"))
											}
										} else {
											pid = model.PersonID(in.Int())
										}
										s.SignedUp = append(s.SignedUp, pid)
									}
									in.WantComma()
								}
								in.Delim(']')
							case "declined":
								in.Delim('[')
								for !in.IsDelim(']') {
									if in.IsNull() {
										in.Skip()
									} else {
										var pid model.PersonID
										if in.IsDelim('{') {
											var seen bool
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
													seen = true
												default:
													in.SkipRecursive()
												}
												in.WantComma()
											}
											in.Delim('}')
											if !seen {
												in.AddError(errors.New("missing assignments.shifts.declined.id"))
											}
										} else {
											pid = model.PersonID(in.Int())
										}
										s.Declined = append(s.Declined, pid)
									}
									in.WantComma()
								}
								in.Delim(']')
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						e.Shifts = append(e.Shifts, &s)
					}
					in.WantComma()
				}
				in.Delim(']')
			case "attendance":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var ai model.AttendanceInfo
						var pid model.PersonID
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
							case "person":
								pid = model.PersonID(in.Int())
							case "type":
								ai.Type, err = model.ParseAttendanceType(in.String())
								in.AddError(err)
							case "minutes":
								ai.Minutes = in.Uint16()
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						if pid == 0 {
							in.AddError(errors.New("missing attendance.person"))
						}
						ea[pid] = ai
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
		if err := event.ValidateEvent(tx, e); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if e.ID == 0 {
			tx.CreateEvent(e)
		} else {
			tx.UpdateEvent(e)
		}
		tx.SaveEventAttendance(e, ea)
		record++
	}
}
