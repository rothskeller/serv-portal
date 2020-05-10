package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"sunnyvaleserv.org/portal/model"
	"sunnyvaleserv.org/portal/store"
	"sunnyvaleserv.org/portal/util/log"
)

func main() {
	switch os.Getenv("HOME") {
	case "/home/snyserv":
		if err := os.Chdir("/home/snyserv/sunnyvaleserv.org/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	case "/Users/stever":
		if err := os.Chdir("/Users/stever/src/serv-portal/data"); err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: %s\n", err)
			os.Exit(1)
		}
	}
	store.Open("serv.db")
	entry := log.New("", "dsw-import")
	defer entry.Log()
	tx := store.Begin(entry)
	scan := bufio.NewScanner(os.Stdin)
	for scan.Scan() {
		var form model.DSWForm
		fields := strings.Split(scan.Text(), "\t")
		id, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(err)
		}
		from, err := time.ParseInLocation("2006-01-02", fields[2], time.Local)
		if err != nil {
			panic(err)
		}
		form.From = from
		if len(fields) > 3 {
			form.Invalid = fields[3]
		}
		if form.Invalid == "" {
			form.To = time.Date(form.From.Year()+1, form.From.Month(), form.From.Day(), 0, 0, 0, 0, time.Local)
		}
		form.For = "CERT"
		person := tx.FetchPerson(model.PersonID(id))
		tx.WillUpdatePerson(person)
		person.DSWForms = append(person.DSWForms, &form)
		tx.UpdatePerson(person)
	}
	tx.Commit()
}
