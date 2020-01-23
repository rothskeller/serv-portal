package main

import (
	"fmt"

	"rothskeller.net/serv/db"
	"rothskeller.net/serv/model"
	"rothskeller.net/serv/person"
)

func main() {
	db.Open("data/serv.db")
	tx := db.Begin()
	for _, p := range tx.FetchPeople() {
		change := addDistrict(&p.HomeAddress) || addDistrict(&p.WorkAddress)
		if change {
			tx.SavePerson(p)
		}
	}
	tx.Commit()
}

func addDistrict(a *model.Address) bool {
	if a.Latitude == 0 || a.Longitude == 0 {
		return false
	}
	a.FireDistrict = person.FireDistrict(a)
	fmt.Printf("%d   %s\n", a.FireDistrict, a.Address)
	return a.FireDistrict != 0
}
