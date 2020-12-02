package main

import (
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/text"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadTextMessages(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var t = new(model.TextMessage)
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
				if !first {
					fmt.Fprintf(os.Stderr, "ERROR: id must be first key in role\n")
					os.Exit(1)
				}
				t.ID = model.TextMessageID(in.Int())
				if t.ID != 0 {
					tid := t.ID
					if t = tx.FetchTextMessage(tid); t == nil {
						fmt.Fprintf(os.Stderr, "ERROR: text message %d does not exist\n", tid)
						os.Exit(1)
					}
					*t = model.TextMessage{ID: tid}
				}
			case "sender":
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
							t.Sender = model.PersonID(in.Int())
						default:
							in.SkipRecursive()
						}
						in.WantComma()
					}
					in.Delim('}')
				} else {
					t.Sender = model.PersonID(in.Int())
				}
			case "lists":
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
									t.Lists = append(t.Lists, model.ListID(in.Int()))
								default:
									in.SkipRecursive()
								}
								in.WantComma()
							}
							in.Delim('}')
						} else {
							t.Lists = append(t.Lists, model.ListID(in.Int()))
						}
					}
					in.WantComma()
				}
				in.Delim(']')
			case "timestamp":
				if data := in.Raw(); in.Ok() {
					in.AddError(t.Timestamp.UnmarshalJSON(data))
				}
			case "message":
				t.Message = in.String()
			case "recipients":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var tr model.TextRecipient

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
								tr.Recipient = model.PersonID(in.Int())
							case "number":
								tr.Number = in.String()
							case "status":
								tr.Status = in.String()
							case "timestamp":
								if data := in.Raw(); in.Ok() {
									in.AddError(tr.Timestamp.UnmarshalJSON(data))
								}
							case "responses":
								in.Delim('[')
								for !in.IsDelim(']') {
									if in.IsNull() {
										in.Skip()
									} else {
										var trr model.TextResponse
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
											case "response":
												trr.Response = in.String()
											case "timestamp":
												if data := in.Raw(); in.Ok() {
													in.AddError(trr.Timestamp.UnmarshalJSON(data))
												}
											default:
												in.SkipRecursive()
											}
											in.WantComma()
										}
										in.Delim('}')
										tr.Responses = append(tr.Responses, &trr)
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
						t.Recipients = append(t.Recipients, &tr)
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
		if err := text.ValidateTextMessage(tx, t); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if t.ID != 0 {
			tx.UpdateTextMessage(t)
		} else {
			tx.CreateTextMessage(t)
		}
		record++
	}
}
