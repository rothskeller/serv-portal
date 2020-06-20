package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadVenues(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	for {
		var v = new(model.Venue)
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
					fmt.Fprintf(os.Stderr, "ERROR: id must be first key in venue\n")
					os.Exit(1)
				}
				v.ID = model.VenueID(in.Int())
				if v.ID != 0 {
					fmt.Fprintf(os.Stderr, "ERROR: updating existing venues not supported\n")
					os.Exit(1)
				}
			case "name":
				v.Name = in.String()
			case "address":
				v.Address = in.String()
			case "city":
				v.City = in.String()
			case "url":
				v.URL = in.String()
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
			fmt.Fprintf(os.Stderr, "ERROR: id must be first key in venue\n")
			os.Exit(1)
		}
		if v.URL != "" && !strings.HasPrefix(v.URL, "https://www.google.com/maps/") {
			fmt.Fprintf(os.Stderr, "ERROR: invalid venue URL\n")
			os.Exit(1)
		}
		tx.CreateVenue(v)
		record++
	}
}
