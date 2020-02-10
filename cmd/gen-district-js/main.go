package main

import (
	"fmt"

	"github.com/paulmach/orb"

	"sunnyvaleserv.org/portal/api/person"
)

func main() {
	emit(person.District1, "district1")
	emit(person.District2, "district2")
	emit(person.District3, "district3")
	emit(person.District4, "district4")
	emit(person.District5, "district5")
	emit(person.District6, "district6")
}

func emit(region orb.Ring, label string) {
	fmt.Printf("export const %s=[", label)
	for i, p := range region {
		if i != 0 {
			fmt.Print(",")
		}
		fmt.Printf("{lat:%f,lng:%f}", p.Lat(), p.Lon())
	}
	fmt.Println("]")
}
