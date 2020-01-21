// serv-init creates a new serv.db and populates it with the bare minimum
// entries needed for admin to log in.
package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"

	_ "github.com/scholacantorum/go-sqlite3"

	"rothskeller.net/serv/model"
)

func main() {
	var (
		schema []byte
		dbh    *sql.DB
		role   model.Role
		group  model.Group
		authz  model.AuthzData
		person model.Person
		venues model.Venues
		data   []byte
		err    error
	)
	if len(os.Args) != 3 {
		fmt.Fprintf(os.Stderr, "usage: serv-init schema.sql serv.db\n")
		os.Exit(2)
	}
	if schema, err = ioutil.ReadFile(os.Args[1]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\nusage: serv-init schema.sql serv.db\n", err)
		os.Exit(2)
	}
	if _, err = os.Stat(os.Args[2]); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s already exists or has an error: %s\nusage: serv-init schema.sql serv.db\n", os.Args[2], err)
		os.Exit(2)
	}
	if dbh, err = sql.Open("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=1", os.Args[2])); err != nil {
		panic(err)
	}
	if _, err = dbh.Exec(string(schema)); err != nil {
		fmt.Fprintf(os.Stderr, "applying schema: %s\n", err)
		os.Exit(1)
	}
	role.ID = 1
	role.Tag = model.RoleWebmaster
	role.Name = "Site Administrator"
	authz.Roles = []*model.Role{&role}
	group.ID = 1
	group.Tag = model.GroupDisabled
	group.Name = "Disabled Users"
	authz.Groups = []*model.Group{&group}
	if data, err = authz.Marshal(); err != nil {
		fmt.Fprintf(os.Stderr, "marshaling authz: %s\n", err)
		os.Exit(1)
	}
	if _, err = dbh.Exec(`INSERT INTO authz VALUES (?)`, data); err != nil {
		fmt.Fprintf(os.Stderr, "saving authz: %s\n", err)
		os.Exit(1)
	}
	person.ID = 1
	person.Username = "admin"
	person.FormalName = "Administrator"
	person.SortName = "Administrator"
	person.Nickname = "Admin"
	person.Roles = []model.RoleID{1}
	if data, err = person.Marshal(); err != nil {
		fmt.Fprintf(os.Stderr, "marshaling person: %s\n", err)
		os.Exit(1)
	}
	if _, err = dbh.Exec(`INSERT INTO person (id, username, data) VALUES (?,?,?)`, 1, "admin", data); err != nil {
		fmt.Fprintf(os.Stderr, "saving person: %s\n", err)
		os.Exit(1)
	}
	if data, err = venues.Marshal(); err != nil {
		fmt.Fprintf(os.Stderr, "marshaling venues: %s\n", err)
		os.Exit(1)
	}
	if _, err = dbh.Exec(`INSERT INTO venue VALUES (?)`, data); err != nil {
		fmt.Fprintf(os.Stderr, "saving venues: %s\n", err)
		os.Exit(1)
	}
}
