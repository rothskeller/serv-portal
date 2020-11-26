package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/mailru/easyjson/jlexer"

	"sunnyvaleserv.org/portal/api/folder"
	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
)

func loadFolders(tx *store.Tx, in *jlexer.Lexer) {
	var record = 1
	var err error
	for {
		var f = &model.FolderNode{Folder: &model.Folder{}}
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
				f.ID = model.FolderID(in.Int())
				if f.ID != 0 {
					if !first {
						fmt.Fprintf(os.Stderr, "ERROR: id must be first key in folder\n")
						os.Exit(1)
					}
					fid := f.ID
					if f = tx.FetchFolder(f.ID); f == nil {
						fmt.Fprintf(os.Stderr, "ERROR: folder %d does not exist\n", fid)
						os.Exit(1)
					}
					tx.WillUpdateFolder(f)
					*f.Folder = model.Folder{ID: fid}
				}
			case "parent":
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
							f.Parent = model.FolderID(in.Int())
							seen = true
						default:
							in.SkipRecursive()
						}
						in.WantComma()
					}
					in.Delim('}')
					if !seen {
						in.AddError(errors.New("missing parent.id"))
					}
				} else {
					f.Parent = model.FolderID(in.Int())
				}
			case "name":
				f.Name = in.String()
			case "group":
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
							f.Group = model.GroupID(in.Int())
							seen = true
						default:
							in.SkipRecursive()
						}
						in.WantComma()
					}
					in.Delim('}')
					if !seen {
						in.AddError(errors.New("missing group.id"))
					}
				} else {
					f.Group = model.GroupID(in.Int())
				}
			case "org":
				f.Org, err = model.ParseFolderOrg(in.String())
				in.AddError(err)
			case "documents":
				in.Delim('[')
				for !in.IsDelim(']') {
					if in.IsNull() {
						in.Skip()
					} else {
						var d model.Document
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
								d.ID = model.DocumentID(in.Int())
							case "name":
								d.Name = in.String()
							case "postedBy":
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
											d.PostedBy = model.PersonID(in.Int())
										default:
											in.SkipRecursive()
										}
										in.WantComma()
									}
									in.Delim('}')
								} else {
									d.PostedBy = model.PersonID(in.Int())
								}
							case "postedAt":
								if data := in.Raw(); in.Ok() {
									in.AddError(d.PostedAt.UnmarshalJSON(data))
								}
							case "needsApproval":
								d.NeedsApproval = in.Bool()
							default:
								in.SkipRecursive()
							}
							in.WantComma()
						}
						in.Delim('}')
						f.Documents = append(f.Documents, &d)
					}
					in.WantComma()
				}
				in.Delim(']')
			case "approvals":
				f.Approvals = in.Int()
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
		if err := folder.ValidateFolder(tx, f); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: record %d: %s\n", record, err)
			os.Exit(1)
		}
		if f.ID == 0 {
			tx.CreateFolder(f)
		} else {
			tx.UpdateFolder(f)
		}
		record++
	}
}
