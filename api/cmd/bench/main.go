package main

import (
	"fmt"
	"os"
	"time"

	"rothskeller.net/serv/db"
)

func main() {
	// Change working directory to the data subdirectory of the CGI script
	// location.  This directory should be mode 700 so that it not directly
	// readable by the web server.
	if err := os.Chdir("data"); err != nil {
		fmt.Printf("Status: 500 Internal Server Error\nContent-Type: text/plain\n\n%s\n", err)
		os.Exit(1)
	}
	db.Open("serv.db")
	start := time.Now()
	tx := db.Begin()
	fmt.Println(time.Since(start).String())
	tx.Rollback()
}
