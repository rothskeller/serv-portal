package main

import (
	"fmt"
	"os"

	"rothskeller.net/serv/db"
)

func main() {
	db.Open("data/serv.db")
	tx := db.Begin()
	for _, p := range tx.FetchPeople() {
		fmt.Print(".")
		switch len(p.Phones) {
		case 0:
			continue
		case 1:
			switch p.Phones[0].Label {
			case "":
				p.HomePhone = p.Phones[0].Phone
				tx.SavePerson(p)
				continue
			case "Cell":
				p.CellPhone = p.Phones[0].Phone
				tx.SavePerson(p)
				continue
			case "Home":
				p.HomePhone = p.Phones[0].Phone
				tx.SavePerson(p)
				continue
			}
		case 2:
			switch p.Phones[0].Label {
			case "":
				switch p.Phones[1].Label {
				case "":
					p.HomePhone = p.Phones[0].Phone
					p.Archive = append(p.Archive, "extra-phone="+p.Phones[1].Phone)
					tx.SavePerson(p)
					continue
				case "Cell":
					tx.SavePerson(p)
					continue
				case "Home":
					// Assume the unlabeled one is a cell.
					p.CellPhone = p.Phones[0].Phone
					p.HomePhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				}
			case "Cell":
				switch p.Phones[1].Label {
				case "":
					p.CellPhone = p.Phones[0].Phone
					p.HomePhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				case "Cell":
					p.CellPhone = p.Phones[0].Phone
					p.Archive = append(p.Archive, "extra-phone="+p.Phones[1].Phone+" Cell")
					tx.SavePerson(p)
					continue
				case "Home":
					p.CellPhone = p.Phones[0].Phone
					p.HomePhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				case "Work":
					p.CellPhone = p.Phones[0].Phone
					p.WorkPhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				}
			case "Home":
				switch p.Phones[1].Label {
				case "Cell":
					p.HomePhone = p.Phones[0].Phone
					p.CellPhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				}
			case "Work":
				switch p.Phones[1].Label {
				case "Cell":
					p.WorkPhone = p.Phones[0].Phone
					p.CellPhone = p.Phones[1].Phone
					tx.SavePerson(p)
					continue
				}
			}
		case 3:
			switch p.Phones[0].Label {
			case "Cell":
				switch p.Phones[1].Label {
				case "":
					switch p.Phones[2].Label {
					case "Cell":
						p.CellPhone = p.Phones[0].Phone
						p.HomePhone = p.Phones[1].Phone
						p.Archive = append(p.Archive, "extra-phone="+p.Phones[2].Phone+" Cell")
						tx.SavePerson(p)
						continue
					}
				}
			case "Cell, Verizon":
				switch p.Phones[1].Label {
				case "Cell, Google":
					switch p.Phones[2].Label {
					case "":
						p.CellPhone = p.Phones[0].Phone
						p.Archive = append(p.Archive, "extra-phone="+p.Phones[1].Phone+" Cell")
						p.HomePhone = p.Phones[2].Phone
						tx.SavePerson(p)
						continue
					}
				}
			case "Home":
				switch p.Phones[1].Label {
				case "Cell":
					switch p.Phones[2].Label {
					case "Work":
						p.HomePhone = p.Phones[0].Phone
						p.CellPhone = p.Phones[1].Phone
						p.WorkPhone = p.Phones[2].Phone
						tx.SavePerson(p)
						continue
					}
				case "Work":
					switch p.Phones[2].Label {
					case "Cell":
						p.HomePhone = p.Phones[0].Phone
						p.WorkPhone = p.Phones[1].Phone
						p.CellPhone = p.Phones[2].Phone
						tx.SavePerson(p)
						continue
					}
				}
			case "Work":
				switch p.Phones[1].Label {
				case "Home":
					switch p.Phones[2].Label {
					case "Cell":
						p.WorkPhone = p.Phones[0].Phone
						p.HomePhone = p.Phones[1].Phone
						p.CellPhone = p.Phones[2].Phone
						tx.SavePerson(p)
						continue
					}
				}
			}
		}
		fmt.Fprintf(os.Stderr, "ERROR: unexpected pattern of phones for person %d.\n", p.ID)
		os.Exit(1)
	}
	tx.Commit()
	fmt.Println()
}
